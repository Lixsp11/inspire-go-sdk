package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	train "github.com/Lixsp11/inspire-go-sdk/src/trainv2/client"
	user "github.com/Lixsp11/inspire-go-sdk/src/userv2/client"
	"github.com/alibabacloud-go/tea/dara"
)

func TestPaginatorSnapshotsAndRetainsObservedMetadata(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		page := calls.Add(1)
		if body["page_num"] != float64(page) || body["workspace_id"] != "ws" {
			t.Error("wrong requested page")
		}
		jobs := []map[string]interface{}{{"job_id": "one"}, {"job_id": "two"}}
		if page == 2 {
			jobs = []map[string]interface{}{{"job_id": "three"}}
		}
		_, _ = w.Write(envelope(map[string]interface{}{"jobs": jobs, "total": "3"}))
	}))
	defer server.Close()
	client, err := train.NewClient(config(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	request := (&train.ListJobsRequest{}).SetWorkspaceId("ws").SetPageSize(2)
	pager, err := train.NewListJobsPaginator(client, request, nil)
	if err != nil {
		t.Fatal(err)
	}
	request.SetWorkspaceId("changed")
	count := 0
	for pager.HasMorePages() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		count += len(page.Body.Result.Jobs)
		if page.SdkContext.RequestId != "fixture-request" {
			t.Fatal("metadata lost")
		}
	}
	total, present := pager.ReportedTotal()
	if count != 3 || calls.Load() != 2 || pager.PagesRead() != 2 || pager.LastPageNum() != 2 || total != 3 || !present {
		t.Fatal("pagination metadata incorrect")
	}
}
func TestPaginatorRejectsRepeatedChangedAndTruncatedPages(t *testing.T) {
	for _, mode := range []string{"repeated", "changed", "truncated", "limit"} {
		t.Run(mode, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				page := calls.Add(1)
				jobs := []map[string]interface{}{{"job_id": "one"}}
				total := 3
				if page > 1 {
					switch mode {
					case "changed":
						total = 5
						jobs = []map[string]interface{}{{"job_id": "two"}}
					case "truncated":
						jobs = []map[string]interface{}{}
					}
				}
				_, _ = w.Write(envelope(map[string]interface{}{"jobs": jobs, "total": total}))
			}))
			defer server.Close()
			client, err := train.NewClient(config(server.URL))
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()
			pager, err := train.NewListJobsPaginator(client, (&train.ListJobsRequest{}).SetWorkspaceId("ws").SetPageSize(1), nil)
			if err != nil {
				t.Fatal(err)
			}
			if mode == "limit" {
				_ = pager.SetMaxPages(1)
			}
			_, err = pager.NextPage(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			_, err = pager.NextPage(context.Background())
			sdkError(t, err, "PaginationError")
			if pager.HasMorePages() {
				t.Fatal("failed paginator remains runnable")
			}
		})
	}
}
func TestWaiterSuccessFailureAndTimeout(t *testing.T) {
	for _, mode := range []string{"success", "failure", "timeout"} {
		t.Run(mode, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("Action") != "GetJob" {
					t.Error("waiter made non-read call")
				}
				n := calls.Add(1)
				status := "RUNNING"
				if mode == "failure" {
					status = "FAILED"
				}
				if mode == "success" && n > 1 {
					status = "SUCCEEDED"
				}
				_, _ = w.Write(envelope(map[string]interface{}{"job_id": "job", "status": status}))
			}))
			defer server.Close()
			client, err := train.NewClient(config(server.URL))
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()
			result, err := client.WaitForJob(context.Background(), "job", &train.WaiterOptions{Timeout: 30 * time.Millisecond, Interval: time.Millisecond})
			if mode == "success" {
				if err != nil || *result.Body.Result.Status != "SUCCEEDED" {
					t.Fatal(err)
				}
			} else {
				sdkError(t, err, "WaiterError")
			}
		})
	}
}
func TestAuthenticationRefreshUsesSameAttemptsAndBackoff(t *testing.T) {
	var loads, calls atomic.Int32
	provider, err := openapi.NewRefreshingCredentialProvider(func(context.Context) (*openapi.Credential, error) {
		n := loads.Add(1)
		value := "old"
		if n > 1 {
			value = "new"
		}
		return openapi.NewCredential("cookie", value, time.Time{})
	})
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if n == 1 {
			if r.Header.Get("Cookie") != "old" {
				t.Error("wrong initial credential")
			}
			w.WriteHeader(401)
			return
		}
		if r.Header.Get("Cookie") != "new" {
			t.Error("credential did not refresh")
		}
		_, _ = w.Write(envelope(map[string]interface{}{"id": "u"}))
	}))
	defer server.Close()
	cfg := (&openapi.Config{}).SetEndpoint(server.URL).SetCredential(provider).SetRetryOptions(&dara.RetryOptions{Retryable: true, MaxAttempts: 2, RetryCondition: []*dara.RetryCondition{{MaxAttempts: 2, Exception: []string{"AuthenticationError"}, MaxDelay: 120000, Backoff: &dara.FixedBackoffPolicy{Period: 30}}}})
	client := userClient(t, cfg)
	started := time.Now()
	result, err := client.GetUserDetail(&user.GetUserDetailRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if time.Since(started) < 25*time.Millisecond || result.SdkContext.Attempts != 2 || loads.Load() != 2 {
		t.Fatal("auth refresh ignored backoff or count")
	}
}
func TestConfigSnapshotAndInputErrors(t *testing.T) {
	for _, raw := range []string{`{"tlsMinVersion":"TLSv1.3"}`, `{"disableHttp2":false}`, `{"endpoint":42}`} {
		var cfg openapi.Config
		if err := json.Unmarshal([]byte(raw), &cfg); err == nil {
			t.Fatal("unsupported config accepted")
		}
	}
	cfg := config("http://127.0.0.1:9")
	cfg.ReadTimeout = dara.Int(0)
	if _, err := openapi.NewClient(cfg); err == nil {
		t.Fatal("zero timeout accepted")
	}
	cfg = config("https://not-used.example")
	cfg.GlobalParameters = &openapi.GlobalParameters{Headers: map[string]*string{"Cookie": dara.String("one"), "cookie": dara.String("two")}}
	if _, err := openapi.NewClient(cfg); err == nil {
		t.Fatal("ambiguous aliases accepted")
	}
	cfg = config("https://not-used.example")
	if strings.Contains(cfg.String(), "fixture") || strings.Contains(cfg.GoString(), "fixture") {
		t.Fatal("config string leaked session")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-App") != "original" {
			t.Error("config was not snapshotted")
		}
		_, _ = w.Write(envelope(map[string]interface{}{"id": "u"}))
	}))
	defer server.Close()
	cfg = config(server.URL)
	cfg.GlobalParameters = &openapi.GlobalParameters{Headers: map[string]*string{"X-App": dara.String("original")}}
	client := userClient(t, cfg)
	*cfg.Endpoint = "https://wrong.example"
	*cfg.GlobalParameters.Headers["X-App"] = "changed"
	if _, err := client.GetUserDetail(&user.GetUserDetailRequest{}); err != nil {
		t.Fatal(err)
	}
}
