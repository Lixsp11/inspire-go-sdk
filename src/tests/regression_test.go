package tests

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
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

func TestTypedNilResponseRejectedBeforeAnyWrite(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		_, _ = w.Write(envelope(map[string]interface{}{}))
	}))
	defer server.Close()
	c, err := openapi.NewClient(config(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	p := (&openapi.Params{}).SetAction("StopJob").SetVersion("v2").SetPathname("/api/v2/train").SetMethod("POST").SetAuthType("SSO")
	var invalid *train.StopJobResponse
	err = c.CallTyped(context.Background(), p, (&train.StopJobRequest{}).SetJobId("job"), nil, invalid)
	sdkError(t, err, "ValidationError")
	if calls.Load() != 0 {
		t.Fatal("typed nil caused a write")
	}
}
func TestWaiterCancellationIsStructuredAndMatchesContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(envelope(map[string]interface{}{"job_id": "job", "status": "RUNNING"}))
	}))
	defer server.Close()
	cfg := config(server.URL)
	cfg.Observer = func(openapi.Event) { cancel() }
	c, err := train.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	_, err = c.WaitForJob(ctx, "job", &train.WaiterOptions{Interval: time.Second})
	sdkError(t, err, "CanceledError")
	if !errors.Is(err, context.Canceled) {
		t.Fatal("context cause missing")
	}
}
func TestPermanentTLSNegotiationFailureIsNotRetried(t *testing.T) {
	var connections atomic.Int32
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { t.Error("unexpected HTTP request") }))
	server.TLS = &tls.Config{MinVersion: tls.VersionTLS10, MaxVersion: tls.VersionTLS11}
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			connections.Add(1)
		}
	}
	server.StartTLS()
	defer server.Close()
	client := userClient(t, config(server.URL))
	_, err := client.GetUserDetailWithOptions(&user.GetUserDetailRequest{}, retries(3))
	sdkError(t, err, "TlsError")
	if connections.Load() != 1 {
		t.Fatalf("permanent TLS failure retried %d times", connections.Load())
	}
}
func TestRetryAfterIsExposedInOfficialMilliseconds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(429)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()
	client := userClient(t, config(server.URL))
	_, err := client.GetUserDetail(&user.GetUserDetailRequest{})
	failure := sdkError(t, err, "ThrottlingError")
	if failure.GetRetryAfter() == nil || *failure.GetRetryAfter() != 1000 {
		t.Fatal("Retry-After not exposed in milliseconds")
	}
	var base *dara.SDKError
	if !errors.As(err, &base) {
		t.Fatal("official SDK error not discoverable")
	}
	if dara.StringValue(base.GetName()) != "ThrottlingError" {
		t.Fatal("official SDK error name delegate is incomplete")
	}
}
func TestExplicitEmptyNativeCollectionsRoundTrip(t *testing.T) {
	query := &train.CreateJobConsoleRequest{}
	query.ExtraFields = map[string]json.RawMessage{"future_empty": json.RawMessage(`[]`)}
	query.Envs = []*train.CreateJobConsoleRequestEnvs{}
	data, err := json.Marshal(query)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"envs":[]`) || !strings.Contains(string(data), `"future_empty":[]`) {
		t.Fatal("empty native fields were dropped")
	}
}
