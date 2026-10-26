package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	user "github.com/Lixsp11/inspire-go-sdk/src/userv2/client"
	"github.com/alibabacloud-go/tea/dara"
)

func TestRefreshingProviderSingleFlightAndCanceledWaiter(t *testing.T) {
	var loads atomic.Int32
	entered, release := make(chan struct{}), make(chan struct{})
	provider, err := openapi.NewRefreshingCredentialProvider(func(ctx context.Context) (*openapi.Credential, error) {
		loads.Add(1)
		close(entered)
		select {
		case <-release:
			return openapi.NewCredential("cookie", "session=fixture", time.Time{})
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	first := make(chan error, 1)
	go func() { _, e := provider.GetCredential(context.Background()); first <- e }()
	<-entered
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err = provider.GetCredential(ctx); !errors.Is(err, context.Canceled) {
		t.Fatal("waiting context ignored")
	}
	var group sync.WaitGroup
	for i := 0; i < 12; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			if _, e := provider.GetCredential(context.Background()); e != nil {
				t.Error(e)
			}
		}()
	}
	close(release)
	if err = <-first; err != nil {
		t.Fatal(err)
	}
	group.Wait()
	if loads.Load() != 1 {
		t.Fatal("credential load was not coalesced")
	}
}
func TestLoaderPanicDoesNotPoisonFutureLoads(t *testing.T) {
	var loads atomic.Int32
	provider, err := openapi.NewRefreshingCredentialProvider(func(context.Context) (*openapi.Credential, error) {
		if loads.Add(1) == 1 {
			panic("private-loader-details")
		}
		return openapi.NewCredential("cookie", "session=fixture", time.Time{})
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = provider.GetCredential(context.Background())
	sdkError(t, err, "CredentialError")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err = provider.GetCredential(ctx); err != nil {
		t.Fatal(err)
	}
}
func TestExpiredCredentialIsReloadedBeforeRetry(t *testing.T) {
	var loads, calls atomic.Int32
	provider, _ := openapi.NewRefreshingCredentialProvider(func(context.Context) (*openapi.Credential, error) {
		n := loads.Add(1)
		if n == 1 {
			return openapi.NewCredential("cookie", "old", time.Now().Add(200*time.Millisecond))
		}
		return openapi.NewCredential("cookie", "new", time.Time{})
	})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(503)
			return
		}
		if r.Header.Get("Cookie") != "new" {
			t.Error("expired credential was reused")
		}
		_, _ = w.Write(envelope(map[string]interface{}{"id": "u"}))
	}))
	defer server.Close()
	cfg := (&openapi.Config{}).SetEndpoint(server.URL).SetCredential(provider)
	client := userClient(t, cfg)
	runtime := &dara.RuntimeOptions{Autoretry: dara.Bool(true), MaxAttempts: dara.Int(2), BackoffPolicy: dara.String("fix"), BackoffPeriod: dara.Int(250)}
	if _, err := client.GetUserDetailWithOptions(&user.GetUserDetailRequest{}, runtime); err != nil {
		t.Fatal(err)
	}
	if loads.Load() != 2 {
		t.Fatal("expired provider cache was reused")
	}
}
func TestCredentialValidationAndAmbiguousEnvironment(t *testing.T) {
	for _, value := range []string{"", " leading", "session=bad\nheader"} {
		if _, err := openapi.NewCredential("cookie", value, time.Time{}); err == nil {
			t.Fatal("invalid credential accepted")
		}
	}
	t.Setenv("INSPIRE_COOKIE", "session=fixture")
	t.Setenv("INSPIRE_ACCESS_TOKEN", "token-fixture")
	provider := &openapi.EnvironmentCredentialProvider{}
	_, err := provider.GetCredential(context.Background())
	sdkError(t, err, "CredentialError")
}
