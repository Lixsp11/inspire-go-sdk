package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	notebook "github.com/Lixsp11/inspire-go-sdk/src/notebookv2/client"
	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	project "github.com/Lixsp11/inspire-go-sdk/src/projectv2/client"
	train "github.com/Lixsp11/inspire-go-sdk/src/trainv2/client"
	user "github.com/Lixsp11/inspire-go-sdk/src/userv2/client"
	workspace "github.com/Lixsp11/inspire-go-sdk/src/workspacev2/client"
	"github.com/alibabacloud-go/tea/dara"
)

type fixture struct {
	Path   string                 `json:"path"`
	Body   map[string]interface{} `json:"body"`
	Result map[string]interface{} `json:"result"`
}

func fixtures(t *testing.T) map[string]fixture {
	t.Helper()
	data, err := os.ReadFile("testdata/contracts.json")
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]fixture
	if err = json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	return out
}
func envelope(result interface{}) []byte {
	data, _ := json.Marshal(map[string]interface{}{"ResponseMetadata": map[string]interface{}{"RequestId": "fixture-request"}, "Result": result})
	return data
}
func equalJSON(t *testing.T, a, b []byte) {
	t.Helper()
	var x, y interface{}
	if err := json.Unmarshal(a, &x); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &y); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(x, y) {
		t.Fatalf("JSON mismatch:\n%s\n%s", a, b)
	}
}
func config(endpoint string) *openapi.Config {
	return (&openapi.Config{}).SetEndpoint(endpoint).SetCookie("inspire-session=fixture")
}
func sdkError(t *testing.T, err error, name string) *openapi.SDKError {
	t.Helper()
	if err == nil {
		t.Fatalf("expected %s", name)
	}
	sdk, ok := err.(*openapi.SDKError)
	if !ok || dara.StringValue(sdk.Name) != name {
		t.Fatalf("want %s; got %T %v", name, err, err)
	}
	return sdk
}
func userClient(t *testing.T, cfg *openapi.Config) *user.Client {
	t.Helper()
	client, err := user.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close() })
	return client
}

func TestNativeServiceContracts(t *testing.T) {
	reference := fixtures(t)
	type item struct {
		service, action string
		request         interface{}
		body, result    map[string]interface{}
	}
	cases := []item{
		{"train", "ListJobs", &train.ListJobsRequest{}, reference["list_jobs"].Body, reference["list_jobs"].Result},
		{"train", "GetJob", &train.GetJobRequest{}, map[string]interface{}{"job_id": "job-1"}, map[string]interface{}{"job_id": "job-1", "status": "RUNNING"}},
		{"train", "CreateJobConsole", &train.CreateJobConsoleRequest{}, reference["create_job"].Body, reference["create_job"].Result},
		{"train", "StopJob", &train.StopJobRequest{}, map[string]interface{}{"job_id": "job-1"}, map[string]interface{}{}},
		{"project", "ListProjects", &project.ListProjectsRequest{}, reference["list_projects"].Body, reference["list_projects"].Result},
		{"notebook", "ListNotebooks", &notebook.ListNotebooksRequest{}, reference["list_notebooks"].Body, reference["list_notebooks"].Result},
		{"notebook", "GetNotebook", &notebook.GetNotebookRequest{}, map[string]interface{}{"notebook_id": "nb-1"}, map[string]interface{}{"notebook_id": "nb-1"}},
		{"workspace", "ListLogicComputeGroups", &workspace.ListLogicComputeGroupsRequest{}, map[string]interface{}{"page_size": -1, "page_num": 1, "filter": map[string]interface{}{"workspace_id": "ws-x"}}, map[string]interface{}{"logic_compute_groups": []interface{}{}}},
		{"user", "GetUserDetail", &user.GetUserDetailRequest{}, map[string]interface{}{}, map[string]interface{}{"id": "u-fixture", "name_en": "fixture"}},
	}
	for _, test := range cases {
		t.Run(test.action, func(t *testing.T) {
			var count atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				count.Add(1)
				if r.Method != "POST" || r.URL.Path != "/api/v2/"+test.service || r.URL.Query().Get("Action") != test.action {
					t.Errorf("wrong native routing")
				}
				if r.Header.Get("Cookie") != "inspire-session=fixture" {
					t.Errorf("wrong credential")
				}
				data, _ := io.ReadAll(r.Body)
				expected, _ := json.Marshal(test.body)
				equalJSON(t, data, expected)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(envelope(test.result))
			}))
			defer server.Close()
			var client interface{}
			var err error
			switch test.service {
			case "train":
				client, err = train.NewClient(config(server.URL))
			case "user":
				client, err = user.NewClient(config(server.URL))
			case "project":
				client, err = project.NewClient(config(server.URL))
			case "notebook":
				client, err = notebook.NewClient(config(server.URL))
			case "workspace":
				client, err = workspace.NewClient(config(server.URL))
			}
			if err != nil {
				t.Fatal(err)
			}
			defer client.(interface{ Close() error }).Close()
			data, _ := json.Marshal(test.body)
			if err = json.Unmarshal(data, test.request); err != nil {
				t.Fatal(err)
			}
			before, _ := json.Marshal(test.request)
			for _, suffix := range []string{"", "WithOptions"} {
				args := []reflect.Value{reflect.ValueOf(test.request)}
				if suffix != "" {
					args = append(args, reflect.ValueOf(&dara.RuntimeOptions{}))
				}
				result := reflect.ValueOf(client).MethodByName(test.action + suffix).Call(args)
				if !result[1].IsNil() {
					t.Fatal(result[1].Interface())
				}
				out := result[0].Interface()
				raw, _ := json.Marshal(out)
				var object struct {
					Headers    map[string]string `json:"headers"`
					StatusCode int               `json:"statusCode"`
					Body       struct {
						Result json.RawMessage `json:"Result"`
					} `json:"body"`
				}
				if err = json.Unmarshal(raw, &object); err != nil {
					t.Fatal(err)
				}
				if object.StatusCode != 200 || object.Headers["content-type"] != "application/json" {
					t.Fatalf("wrong V2 response: %s", raw)
				}
				expected, _ := json.Marshal(test.result)
				equalJSON(t, object.Body.Result, expected)
				ctx := reflect.ValueOf(out).Elem().FieldByName("SdkContext").Interface().(*openapi.ResponseContext)
				if ctx.Attempts != 1 || ctx.Action != test.action || ctx.RequestId != "fixture-request" {
					t.Fatal("missing native diagnostics")
				}
			}
			after, _ := json.Marshal(test.request)
			equalJSON(t, before, after)
			if count.Load() != 2 {
				t.Fatal(count.Load())
			}
		})
	}
}
func TestDiscoveryAnonymous(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.RequestURI() != "/discovery" || r.Header.Get("Cookie") != "" {
			t.Error("discovery must be anonymous GET")
		}
		_, _ = w.Write(envelope(map[string]interface{}{"Version": "fixture", "Services": []interface{}{}}))
	}))
	defer server.Close()
	c, err := openapi.NewClient(config(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	out, err := c.GetDiscovery()
	if err != nil {
		t.Fatal(err)
	}
	if out.Body.Result["Version"] != "fixture" {
		t.Fatal("bad discovery body")
	}
}
func TestModelUnknownFieldsAndLargeNumbers(t *testing.T) {
	raw := []byte(`{"workspace_id":"ws","new_flag":false,"new_number":9007199254740993}`)
	request := &train.ListJobsRequest{}
	if err := json.Unmarshal(raw, request); err != nil {
		t.Fatal(err)
	}
	out, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "9007199254740993") {
		t.Fatal("unknown integer lost precision")
	}
	value, err := request.ToMap()
	if err != nil {
		t.Fatal(err)
	}
	if value["new_number"].(json.Number).String() != "9007199254740993" {
		t.Fatal("ToMap lost precision")
	}
	request.SetPageSize(0)
	if err = request.Validate(); err == nil {
		t.Fatal("zero page size accepted")
	}
	request.SetPageSize(2)
	if err = request.Validate(); err != nil {
		t.Fatal(err)
	}
	worker := (&train.CreateJobConsoleRequestFrameworkConfig{}).SetGpuCount(0)
	data, _ := json.Marshal(worker)
	if !strings.Contains(string(data), `"gpu_count":0`) {
		t.Fatal("pointer zero omitted")
	}
	req := (&train.CreateJobConsoleRequest{}).SetEnableNotification(false)
	data, _ = json.Marshal(req)
	if !strings.Contains(string(data), `"enable_notification":false`) {
		t.Fatal("pointer false omitted")
	}
}
func TestInvalidRequestsNeverSend(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := train.NewClient(config(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	for _, request := range []*train.ListJobsRequest{nil, {}, {WorkspaceId: dara.String("ws"), PageSize: dara.Int32(1000)}} {
		if _, err = c.ListJobs(request); err == nil {
			t.Fatal("invalid request accepted")
		}
	}
	body := fixtures(t)["create_job"].Body
	data, _ := json.Marshal(body)
	request := &train.CreateJobConsoleRequest{}
	_ = json.Unmarshal(data, request)
	request.TaskPriority = nil
	if _, err = c.CreateJobConsole(request); err == nil {
		t.Fatal("missing priority accepted")
	}
	if calls.Load() != 0 {
		t.Fatal("invalid input was sent")
	}
}
