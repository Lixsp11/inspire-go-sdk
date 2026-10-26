package openapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/Lixsp11/inspire-go-sdk/src/cli"
	"github.com/Lixsp11/inspire-go-sdk/src/openapi"
)

func execute(t *testing.T, input, endpoint string, args ...string) (int, string, string) {
	t.Helper()
	out, diagnostics := &bytes.Buffer{}, &bytes.Buffer{}
	ctx := cli.NewCommandContext(out, diagnostics)
	ctx.Context = context.Background()
	ctx.Stdin = strings.NewReader(input)
	configuration := filepath.Join(t.TempDir(), "config.json")
	ctx.Getenv = func(name string) string {
		switch name {
		case "INSPIRE_COOKIE":
			return "inspire-session=cli-fixture"
		case "INSPIRE_CONFIG_FILE":
			return configuration
		}
		return ""
	}
	if endpoint != "" {
		args = append([]string{"--endpoint", endpoint}, args...)
	}
	code := openapi.Run(ctx, args)
	return code, out.String(), diagnostics.String()
}
func response(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"ResponseMetadata": map[string]interface{}{"RequestId": "cli-test"}, "Result": result})
}

func TestNativeAPIRequests(t *testing.T) {
	data, err := os.ReadFile("../tests/testdata/contracts.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]struct {
		Body   json.RawMessage `json:"body"`
		Result interface{}     `json:"result"`
	}
	if err = json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		fixture, product, action string
		args                     []string
	}{
		{"list_jobs", "train", "ListJobs", []string{"--WorkspaceId", "ws-x", "--PageNum", "2", "--PageSize", "7", "--CreatedBy", "user-x", "--Keyword", "demo"}},
		{"list_projects", "project", "ListProjects", []string{"--Page", "1", "--PageSize", "100", "--Filter.WorkspaceId", "ws-x", "--Filter.CheckAdmin", "true"}},
		{"list_notebooks", "notebook", "ListNotebooks", nil},
		{"create_job", "train", "CreateJobConsole", nil},
	}
	for _, test := range cases {
		t.Run(test.action, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != "POST" || r.URL.Query().Get("Action") != test.action || r.Header.Get("Cookie") != "inspire-session=cli-fixture" {
					t.Error("incorrect action or authentication")
				}
				var got, want interface{}
				if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
					t.Error(err)
				}
				_ = json.Unmarshal(fixtures[test.fixture].Body, &want)
				if !reflect.DeepEqual(got, want) {
					t.Errorf("wire body mismatch: %v != %v", got, want)
				}
				response(w, fixtures[test.fixture].Result)
			}))
			defer server.Close()
			args := append([]string{test.product, test.action}, test.args...)
			input := ""
			if test.args == nil {
				args = append(args, "--body-file", "-")
				input = string(fixtures[test.fixture].Body)
			}
			code, out, diagnostics := execute(t, input, server.URL, args...)
			if code != 0 || diagnostics != "" || !json.Valid([]byte(out)) || calls.Load() != 1 {
				t.Fatalf("code=%d stderr=%s calls=%d output=%s", code, diagnostics, calls.Load(), out)
			}
		})
	}
}

func TestNumberedParametersAndDryRun(t *testing.T) {
	args := []string{"train", "CreateJobConsole", "--Name", "demo", "--WorkspaceId", "ws-x", "--ProjectId", "project-1", "--LogicComputeGroupId", "lcg-1", "--Command", "python train.py", "--Framework", "pytorch", "--TaskPriority", "1", "--EnableNotification", "false", "--FrameworkConfig.1.Cpu", "4", "--FrameworkConfig.1.GpuCount", "0", "--FrameworkConfig.1.MemGi", "16", "--FrameworkConfig.1.ShmGi", "2", "--FrameworkConfig.1.Image", "registry.example/train:v1", "--FrameworkConfig.1.ImageType", "SOURCE_PRIVATE", "--FrameworkConfig.1.InstanceCount", "1", "--FrameworkConfig.1.ResourceSpecPrice.CpuType", "cpu", "--FrameworkConfig.1.ResourceSpecPrice.CpuCount", "4", "--FrameworkConfig.1.ResourceSpecPrice.MemorySizeGib", "16", "--FrameworkConfig.1.ResourceSpecPrice.LogicComputeGroupId", "lcg-1", "--FrameworkConfig.1.ResourceSpecPrice.QuotaId", "quota-1", "--cli-dry-run-json"}
	code, out, diagnostics := execute(t, "", "https://invalid.example", args...)
	if code != 0 || diagnostics != "" {
		t.Fatalf("%d %s", code, diagnostics)
	}
	var result struct {
		DryRun     bool
		Parameters map[string]interface{}
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatal(err)
	}
	if !result.DryRun || result.Parameters["enable_notification"] != false || len(result.Parameters["framework_config"].([]interface{})) != 1 {
		t.Fatal(out)
	}
}

func TestInvalidInputDoesNotContactPlatform(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); response(w, map[string]interface{}{}) }))
	defer server.Close()
	cases := [][]string{
		{"train", "ListJobs", "--WorkspaceId", "ws", "--WorkspaceId", "other"},
		{"train", "ListJobs", "--WorkspaceId", "ws", "--PageSize", "no"},
		{"train", "ListJobs", "--WorkspaceId", "ws", "--PageSize", "2147483648"},
		{"train", "ListJobs", "--WorkspaceID", "ws"},
		{"train", "ListJobs", "--body", `{"workspace_id":"ws","workspace_id":"other"}`},
		{"train", "ListJobs", "--body", `{"workspace_id":"ws"} {}`},
		{"train", "ListJobs", "--body", `{"workspace_id":"ws","extra":1}`},
		{"train", "ListJobs", "--body", `{"workspace_id":"ws"}`, "--PageSize", "1"},
		{"train", "CreateJobConsole", "--FrameworkConfig.0.Cpu", "1"},
		{"notebook", "ListNotebooks", "--WorkspaceId", "ws", "--FilterBy.UserId.1", "alpha", "--FilterBy.UserId.01", "beta"},
		{"notebook", "ListNotebooks", "--WorkspaceId", "ws", "--FilterBy.UserId.1", "alpha", "--FilterBy.UserId.+1", "beta"},
		{"train", "CreateJobConsole", "--FrameworkConfig.01.Cpu", "1"},
		{"train", "CreateJobConsole", "--FrameworkConfig", "[]", "--FrameworkConfig.1.Cpu", "1"},
		{"train", "GetJob", "--JobId", "job-1", "--waiter", "expr=wrong", "to=RUNNING"},
		{"train", "ListJobs", "--WorkspaceId", "ws", "--max-pages", "2"},
	}
	for _, args := range cases {
		code, out, diagnostics := execute(t, "", server.URL, args...)
		if code != 2 || out != "" || !json.Valid([]byte(diagnostics)) {
			t.Errorf("args=%v code=%d stdout=%s stderr=%s", args, code, out, diagnostics)
		}
	}
	if calls.Load() != 0 {
		t.Fatal("invalid input contacted platform")
	}
}

func TestPagerAndWaiter(t *testing.T) {
	t.Run("pager", func(t *testing.T) {
		var calls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n := calls.Add(1)
			response(w, map[string]interface{}{"jobs": []interface{}{map[string]interface{}{"job_id": "job-" + string(rune('0'+n))}}, "total": 2})
		}))
		defer server.Close()
		code, out, diagnostics := execute(t, "", server.URL, "train", "ListJobs", "--WorkspaceId", "ws", "--PageSize", "1", "--pager")
		if code != 0 || calls.Load() != 2 || !strings.Contains(out, "job-2") {
			t.Fatalf("%d %s %s", code, out, diagnostics)
		}
	})
	t.Run("limit", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response(w, map[string]interface{}{"jobs": []interface{}{map[string]interface{}{"job_id": "job-1"}}, "total": 2})
		}))
		defer server.Close()
		code, out, diagnostics := execute(t, "", server.URL, "train", "ListJobs", "--WorkspaceId", "ws", "--PageSize", "1", "--pager", "--max-pages", "1")
		if code == 0 || out != "" || !strings.Contains(diagnostics, "PageLimitReached") {
			t.Fatalf("%d %s %s", code, out, diagnostics)
		}
	})
	t.Run("waiter", func(t *testing.T) {
		var calls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			status := "RUNNING"
			if calls.Add(1) > 1 {
				status = "STOPPED"
			}
			response(w, map[string]interface{}{"job_id": "job-1", "status": status})
		}))
		defer server.Close()
		code, out, diagnostics := execute(t, "", server.URL, "train", "GetJob", "--JobId", "job-1", "--waiter", "expr=Result.status", "to=STOPPED", "timeout=1", "interval=0.001")
		if code != 0 || calls.Load() != 2 || !strings.Contains(out, "STOPPED") {
			t.Fatalf("%d %s %s", code, out, diagnostics)
		}
	})
}

func TestWritesAreNotRetried(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(503)
		_, _ = w.Write([]byte(`{"ResponseMetadata":{"Error":{"Code":"InternalError","Message":"hidden-server-value"}}}`))
	}))
	defer server.Close()
	code, out, diagnostics := execute(t, "", server.URL, "train", "StopJob", "--JobId", "job-1", "--retry-count", "3")
	if code != 1 || out != "" || calls.Load() != 1 || !strings.Contains(diagnostics, "MutationUncertain") || !strings.Contains(diagnostics, "reconcile") || strings.Contains(diagnostics, "hidden-server-value") {
		t.Fatalf("%d %s %s calls=%d", code, out, diagnostics, calls.Load())
	}
}

func TestMachineHelpAndProfile(t *testing.T) {
	code, out, diagnostics := execute(t, "", "", "train", "CreateJobConsole", "--help", "--cli-output", "json")
	if code != 0 || diagnostics != "" || !strings.Contains(out, `"cli_name": "QuotaId"`) || !strings.Contains(out, `"required"`) {
		t.Fatalf("%d %s %s", code, out, diagnostics)
	}
	path := filepath.Join(t.TempDir(), "config.json")
	secret := filepath.Join(t.TempDir(), "session")
	if err := os.WriteFile(secret, []byte("inspire-session=profile-secret"), 0600); err != nil {
		t.Fatal(err)
	}
	code, _, diagnostics = execute(t, "", "", "configure", "set", "--config-path", path, "--profile", "test", "--cookie-file", secret, "--socks5-proxy", "socks5://operator:proxy-secret@127.0.0.1:1080")
	if code != 0 {
		t.Fatal(diagnostics)
	}
	saved, err := os.ReadFile(path)
	if err != nil || bytes.Contains(saved, []byte("inspire-session=")) {
		t.Fatal("session was persisted in configuration")
	}
	code, out, diagnostics = execute(t, "", "", "configure", "get", "--config-path", path, "--profile", "test")
	if code != 0 || strings.Contains(out, "proxy-secret") || strings.Contains(out, "operator") || strings.Contains(out, "profile-secret") {
		t.Fatalf("%d %s %s", code, out, diagnostics)
	}
}

func TestMissingCurrentProfileDoesNotUseDefaultConnection(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"current":"missing","profiles":[]}`), 0600); err != nil {
		t.Fatal(err)
	}
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); response(w, map[string]interface{}{}) }))
	defer server.Close()
	code, out, diagnostics := execute(t, "", server.URL, "train", "StopJob", "--JobId", "job-1", "--config-path", path)
	if code != 2 || out != "" || calls.Load() != 0 || !strings.Contains(diagnostics, "InvalidConfiguration") {
		t.Fatalf("%d %s %s calls=%d", code, out, diagnostics, calls.Load())
	}
}

func TestDryRunValidatesCallerOptions(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); response(w, map[string]interface{}{}) }))
	defer server.Close()
	cases := [][]string{
		{"train", "GetJob", "--JobId", "job-1", "--waiter", "expr=wrong", "to=RUNNING"},
		{"train", "GetJob", "--JobId", "job-1", "--timeout", "invalid"},
		{"train", "GetJob", "--JobId", "job-1", "--retry-count", "invalid"},
		{"train", "GetJob", "--JobId", "job-1", "--read-timeout", "0"},
		{"train", "GetJob", "--JobId", "job-1", "--profile="},
		{"train", "GetJob", "--JobId", "job-1", "--socks5-proxy", "invalid"},
		{"train", "ListJobs", "--WorkspaceId", "ws", "--pager", "--PageNum", "2"},
		{"train", "ListJobs", "--WorkspaceId", "ws", "--pager", "--max-pages", "0"},
	}
	for _, args := range cases {
		for _, dry := range []bool{false, true} {
			input := append([]string(nil), args...)
			if dry {
				input = append(input, "--cli-dry-run-json")
			}
			code, out, diagnostics := execute(t, "", server.URL, input...)
			if code != 2 || out != "" || !json.Valid([]byte(diagnostics)) {
				t.Errorf("dry=%v args=%v code=%d out=%s stderr=%s", dry, args, code, out, diagnostics)
			}
		}
	}
	if calls.Load() != 0 {
		t.Fatal("invalid caller options contacted platform")
	}
}

func TestAuthenticationAndReadRetries(t *testing.T) {
	t.Run("expired session at HTTP 200", func(t *testing.T) {
		var calls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ResponseMetadata":{"Error":{"Code":"TokenExpired","Message":"hidden-session-value"}}}`))
		}))
		defer server.Close()
		code, out, diagnostics := execute(t, "", server.URL, "user", "GetUserDetail", "--retry-count", "3")
		if code != 1 || out != "" || calls.Load() != 1 || !strings.Contains(diagnostics, `"action": "configure"`) || strings.Contains(diagnostics, "hidden-session-value") {
			t.Fatalf("%d %s %s calls=%d", code, out, diagnostics, calls.Load())
		}
	})
	t.Run("waiter inherits retries", func(t *testing.T) {
		var calls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if calls.Add(1) == 1 {
				w.WriteHeader(503)
				_, _ = w.Write([]byte(`{}`))
				return
			}
			response(w, map[string]interface{}{"job_id": "job-1", "status": "RUNNING"})
		}))
		defer server.Close()
		code, out, diagnostics := execute(t, "", server.URL, "train", "GetJob", "--JobId", "job-1", "--retry-count", "1", "--waiter", "expr=Result.status", "--waiter", "to=RUNNING", "timeout=2")
		if code != 0 || calls.Load() != 2 || !strings.Contains(out, "RUNNING") {
			t.Fatalf("%d %s %s calls=%d", code, out, diagnostics, calls.Load())
		}
	})
}

func TestRemainingSDKRoutes(t *testing.T) {
	cases := []struct {
		product, action, path, body string
		args                        []string
		result                      interface{}
	}{
		{"user", "GetUserDetail", "/api/v2/user", `{}`, nil, map[string]interface{}{"id": "user-1"}},
		{"train", "GetJob", "/api/v2/train", `{"job_id":"job-1"}`, []string{"--JobId", "job-1"}, map[string]interface{}{"job_id": "job-1", "status": "RUNNING"}},
		{"train", "StopJob", "/api/v2/train", `{"job_id":"job-1"}`, []string{"--JobId", "job-1"}, map[string]interface{}{}},
		{"notebook", "GetNotebook", "/api/v2/notebook", `{"notebook_id":"nb-1"}`, []string{"--NotebookId", "nb-1"}, map[string]interface{}{"id": "nb-1"}},
		{"workspace", "ListLogicComputeGroups", "/api/v2/workspace", `{"filter":{"workspace_id":"ws-x"},"page_num":1,"page_size":-1}`, []string{"--Filter.WorkspaceId", "ws-x", "--PageNum", "1", "--PageSize", "-1"}, map[string]interface{}{"logic_compute_groups": []interface{}{}, "total": 0}},
		{"openapi", "GetDiscovery", "/discovery", "", nil, map[string]interface{}{"Version": "test", "Services": []interface{}{}}},
	}
	for _, test := range cases {
		t.Run(test.action, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.URL.Path != test.path {
					t.Errorf("wrong path %s", r.URL.Path)
				}
				if test.action == "GetDiscovery" {
					if r.Method != "GET" || r.URL.RawQuery != "" || r.Header.Get("Cookie") != "" {
						t.Error("discovery was not anonymous GET")
					}
				} else {
					if r.Method != "POST" || r.URL.Query().Get("Action") != test.action {
						t.Error("wrong action")
					}
					var got, want interface{}
					_ = json.NewDecoder(r.Body).Decode(&got)
					_ = json.Unmarshal([]byte(test.body), &want)
					if !reflect.DeepEqual(got, want) {
						t.Errorf("wrong native body %v", got)
					}
				}
				response(w, test.result)
			}))
			defer server.Close()
			args := append([]string{test.product, test.action}, test.args...)
			code, out, diagnostics := execute(t, "", server.URL, args...)
			if code != 0 || diagnostics != "" || !json.Valid([]byte(out)) || calls.Load() != 1 {
				t.Fatalf("%d %s %s", code, out, diagnostics)
			}
		})
	}
}
