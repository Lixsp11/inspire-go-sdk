package tests

import (
	"context"
	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	user "github.com/Lixsp11/inspire-go-sdk/src/userv2/client"
	"github.com/alibabacloud-go/tea/dara"
	"os"
	"testing"
	"time"
)

// TestLiveReadonly is opt-in and only reads the current identity through an explicit proxy.
func TestLiveReadonly(t *testing.T) {
	if os.Getenv("INSPIRE_RUN_LIVE") != "1" {
		t.Skip("INSPIRE_RUN_LIVE=1 required")
	}
	if os.Getenv("INSPIRE_COOKIE") == "" {
		t.Skip("existing management Cookie required")
	}
	if os.Getenv("INSPIRE_PROXY") == "" {
		t.Fatal("explicit INSPIRE_PROXY required for this probe")
	}
	cfg := (&openapi.Config{}).SetEndpoint("qz.sii.edu.cn").SetSocks5Proxy(os.Getenv("INSPIRE_PROXY")).SetOperationTimeout(45000)
	client, err := user.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	result, err := client.GetUserDetailWithContext(ctx, &user.GetUserDetailRequest{}, &dara.RuntimeOptions{Autoretry: dara.Bool(true), MaxAttempts: dara.Int(3), BackoffPolicy: dara.String("fix"), BackoffPeriod: dara.Int(1000)})
	if err != nil {
		t.Fatal(err)
	}
	if result.Body.Result.Id == nil && result.Body.Result.UserId == nil {
		t.Fatal("identity response omitted user ID")
	}
	t.Logf("GetUserDetail HTTP=%d attempts=%d authenticated_id_present=true", dara.Int32Value(result.StatusCode), result.SdkContext.Attempts)
}
