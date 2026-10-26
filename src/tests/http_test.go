package tests

import (
	"context"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	train "github.com/Lixsp11/inspire-go-sdk/src/trainv2/client"
	user "github.com/Lixsp11/inspire-go-sdk/src/userv2/client"
	"github.com/alibabacloud-go/tea/dara"
)

func retries(n int) *dara.RuntimeOptions {
	return &dara.RuntimeOptions{Autoretry: dara.Bool(true), MaxAttempts: dara.Int(n), BackoffPolicy: dara.String("no")}
}
func TestHTTPReadRetryAndDefaultSingleSend(t *testing.T) {
	for _, on := range []bool{false, true} {
		t.Run(map[bool]string{false: "default", true: "explicit"}[on], func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if calls.Add(1) == 1 {
					w.WriteHeader(503)
					_, _ = w.Write([]byte(`{}`))
					return
				}
				_, _ = w.Write(envelope(map[string]interface{}{"id": "u"}))
			}))
			defer server.Close()
			client := userClient(t, config(server.URL))
			runtime := &dara.RuntimeOptions{}
			if on {
				runtime = retries(3)
			}
			result, err := client.GetUserDetailWithOptions(&user.GetUserDetailRequest{}, runtime)
			if on {
				if err != nil || calls.Load() != 2 || result.SdkContext.Attempts != 2 {
					t.Fatalf("retry failed %v", err)
				}
			} else {
				sdkError(t, err, "ServerError")
				if calls.Load() != 1 {
					t.Fatal("default unexpectedly retried")
				}
			}
		})
	}
}
func TestHTTP200BusinessAndTransientErrors(t *testing.T) {
	for _, code := range []string{"InvalidParameter", "internal_failure", "TooManyRequestsException"} {
		t.Run(code, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				_, _ = io.WriteString(w, `{"ResponseMetadata":{"RequestId":"req-safe","Error":{"Code":"`+code+`","Message":"remote-secret"}}}`)
			}))
			defer server.Close()
			client := userClient(t, config(server.URL))
			_, err := client.GetUserDetailWithOptions(&user.GetUserDetailRequest{}, retries(2))
			if err == nil || strings.Contains(err.Error(), "remote-secret") {
				t.Fatal("unsafe or missing error")
			}
			expected := int32(2)
			if code == "InvalidParameter" {
				expected = 1
			}
			if calls.Load() != expected {
				t.Fatalf("count %d", calls.Load())
			}
		})
	}
}
func TestWritesNeverReplayUncertainOutcomes(t *testing.T) {
	for _, mode := range []string{"drop", "server", "missing-ack", "invalid-json"} {
		t.Run(mode, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				_, _ = io.ReadAll(r.Body)
				switch mode {
				case "drop":
					conn, _, _ := w.(http.Hijacker).Hijack()
					_ = conn.Close()
				case "server":
					w.WriteHeader(503)
					_, _ = w.Write([]byte(`{}`))
				case "missing-ack":
					_, _ = w.Write(envelope(map[string]interface{}{}))
				case "invalid-json":
					_, _ = w.Write([]byte(`<html>`))
				}
			}))
			defer server.Close()
			client, err := train.NewClient(config(server.URL))
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()
			request := &train.CreateJobConsoleRequest{}
			if err = request.FromMap(fixtures(t)["create_job"].Body); err != nil {
				t.Fatal(err)
			}
			_, err = client.CreateJobConsoleWithOptions(request, retries(3))
			sdkError(t, err, "MutationUncertainError")
			if calls.Load() != 1 {
				t.Fatalf("write replayed %d times", calls.Load())
			}
		})
	}
}
func TestHTTPRetryAfterBudget(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Retry-After", "20")
		w.WriteHeader(429)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()
	cfg := config(server.URL).SetOperationTimeout(100)
	client := userClient(t, cfg)
	start := time.Now()
	_, err := client.GetUserDetailWithOptions(&user.GetUserDetailRequest{}, retries(3))
	sdkError(t, err, "OperationTimeoutError")
	if calls.Load() != 1 || time.Since(start) > time.Second {
		t.Fatal("retry-after ignored budget")
	}
}
func TestHTTPReadTimeoutAndContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		_, _ = w.Write(envelope(map[string]interface{}{"id": "u"}))
	}))
	defer server.Close()
	client := userClient(t, config(server.URL))
	_, err := client.GetUserDetailWithOptions(&user.GetUserDetailRequest{}, &dara.RuntimeOptions{ReadTimeout: dara.Int(10), ConnectTimeout: dara.Int(5)})
	sdkError(t, err, "TransportError")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = client.GetUserDetailWithContext(ctx, &user.GetUserDetailRequest{}, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("context cancellation lost: %v", err)
	}
}
func TestWriteCancellationReportsUncertainty(t *testing.T) {
	entered := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.ReadAll(r.Body)
		close(entered)
		<-r.Context().Done()
	}))
	defer server.Close()
	client, err := train.NewClient(config(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	request := &train.CreateJobConsoleRequest{}
	_ = request.FromMap(fixtures(t)["create_job"].Body)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { _, err := client.CreateJobConsoleWithContext(ctx, request, nil); done <- err }()
	<-entered
	cancel()
	err = <-done
	sdkError(t, err, "MutationUncertainError")
	if !errors.Is(err, context.Canceled) {
		t.Fatal("lost cancellation cause")
	}
}
func TestProviderWinsHeaderAliasesAndCookiesAreNotInherited(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Header.Get("Cookie") != "inspire-session=fixture" || r.Header.Get("Authorization") != "" || r.Header.Get("X-App") != "runtime" {
			t.Error("header precedence failed")
		}
		w.Header().Set("Set-Cookie", "inspire-session=server-side")
		_, _ = w.Write(envelope(map[string]interface{}{"id": "u"}))
	}))
	defer server.Close()
	cfg := config(server.URL).SetGlobalParameters(&openapi.GlobalParameters{Headers: map[string]*string{"cookie": dara.String("wrong"), "X-App": dara.String("global")}})
	client := userClient(t, cfg)
	runtime := &dara.RuntimeOptions{ExtendsParameters: &dara.ExtendsParameters{Headers: map[string]*string{"COOKIE": dara.String("wrong-runtime"), "Authorization": dara.String("wrong-auth"), "x-app": dara.String("runtime")}}}
	for i := 0; i < 2; i++ {
		if _, err := client.GetUserDetailWithOptions(&user.GetUserDetailRequest{}, runtime); err != nil {
			t.Fatal(err)
		}
	}
	if calls.Load() != 2 {
		t.Fatal("unexpected sends")
	}
}
func TestRedirectDoesNotLeakCredentialToOtherOrigin(t *testing.T) {
	var leaked atomic.Int32
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { leaked.Add(1) }))
	defer other.Close()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, other.URL, 302) }))
	defer server.Close()
	client := userClient(t, config(server.URL))
	_, err := client.GetUserDetail(&user.GetUserDetailRequest{})
	sdkError(t, err, "RedirectError")
	if leaked.Load() != 0 {
		t.Fatal("followed cross-origin redirect")
	}
}
func TestTLSCAIsPEMAndVerificationIsDefault(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(envelope(map[string]interface{}{"id": "u"}))
	}))
	defer server.Close()
	client := userClient(t, config(server.URL))
	_, err := client.GetUserDetail(&user.GetUserDetailRequest{})
	sdkError(t, err, "TlsError")
	ca := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw}))
	trusted := userClient(t, config(server.URL).SetCa(ca))
	if _, err = trusted.GetUserDetail(&user.GetUserDetailRequest{}); err != nil {
		t.Fatal(err)
	}
}
func TestConnectionReuseAndConcurrentCalls(t *testing.T) {
	var mu sync.Mutex
	ports := map[string]bool{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		ports[r.RemoteAddr] = true
		mu.Unlock()
		_, _ = w.Write(envelope(map[string]interface{}{"id": "u"}))
	}))
	defer server.Close()
	client := userClient(t, config(server.URL))
	for i := 0; i < 3; i++ {
		if _, err := client.GetUserDetail(&user.GetUserDetailRequest{}); err != nil {
			t.Fatal(err)
		}
	}
	mu.Lock()
	if len(ports) != 1 {
		t.Error("sequential calls did not reuse a connection")
	}
	mu.Unlock()
	var group sync.WaitGroup
	for i := 0; i < 16; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			if _, err := client.GetUserDetail(&user.GetUserDetailRequest{}); err != nil {
				t.Error(err)
			}
		}()
	}
	group.Wait()
	_ = client.Close()
	_ = client.Close()
	_, err := client.GetUserDetail(&user.GetUserDetailRequest{})
	sdkError(t, err, "ClientClosedError")
}
func TestResponseCapFailsWithoutTruncation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(envelope(map[string]interface{}{"id": strings.Repeat("x", 500)}))
	}))
	defer server.Close()
	cfg := config(server.URL)
	cfg.MaxResponseBytes = dara.Int64(100)
	client := userClient(t, cfg)
	_, err := client.GetUserDetail(&user.GetUserDetailRequest{})
	got := sdkError(t, err, "InvalidResponseError")
	if dara.StringValue(got.Code) != "ResponseTooLarge" {
		t.Fatal(got)
	}
}
