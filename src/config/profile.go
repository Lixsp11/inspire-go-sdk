package config

import (
	"github.com/Lixsp11/inspire-go-sdk/src/cli"
	"strconv"
)

func (cp *Profile) OverwriteWithFlags(ctx *cli.Context) error {
	stringsToSet := []struct {
		name  string
		value *string
	}{
		{"endpoint", &cp.Endpoint}, {"cookie-file", &cp.CookieFile}, {"socks5-proxy", &cp.Socks5Proxy}, {"http-proxy", &cp.HttpProxy}, {"https-proxy", &cp.HttpsProxy}, {"no-proxy", &cp.NoProxy},
	}
	for _, f := range stringsToSet {
		if value, ok := ctx.Flags().GetValue(f.name); ok {
			*f.value = value
		}
	}
	for _, f := range []struct {
		name  string
		value *int
	}{{"read-timeout", &cp.ReadTimeout}, {"connect-timeout", &cp.ConnectTimeout}, {"retry-count", &cp.RetryCount}} {
		if raw, ok := ctx.Flags().GetValue(f.name); ok {
			value, err := strconv.Atoi(raw)
			if err != nil || value < 0 || (f.name != "retry-count" && (value == 0 || int64(value) > int64(^uint(0)>>1)/1000)) || (f.name == "retry-count" && value > 99) {
				return cli.UsageError("InvalidParameter", "Invalid --"+f.name+".")
			}
			*f.value = value
		}
	}
	return nil
}
