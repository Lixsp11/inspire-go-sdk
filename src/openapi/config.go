package openapi

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Lixsp11/inspire-go-sdk/src/cli"
	"github.com/Lixsp11/inspire-go-sdk/src/config"
	sdk "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	"github.com/alibabacloud-go/tea/dara"
)

func readBodyFile(ctx *cli.Context, path string) ([]byte, error) {
	var reader io.Reader = ctx.Stdin
	if path != "-" {
		f, err := os.Open(path)
		if err != nil {
			return nil, cli.UsageError("InvalidParameter", "Cannot open --body-file.")
		}
		defer f.Close()
		reader = f
	}
	data, err := io.ReadAll(io.LimitReader(reader, 2<<20+1))
	if err != nil || len(data) > 2<<20 {
		return nil, cli.UsageError("InvalidParameter", "Cannot read request body within the 2 MiB limit.")
	}
	return data, nil
}

func configuration(ctx *cli.Context) (*config.Configuration, string, string, error) {
	path, _ := ctx.Flags().GetValue("config-path")
	if path == "" {
		path = ctx.Getenv("INSPIRE_CONFIG_FILE")
	}
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, "", "", cli.UsageError("InvalidConfiguration", "Set --config-path to a configuration file.")
		}
		path = filepath.Join(home, ".inspire", "config.json")
	}
	c, err := config.LoadConfigurationFromFile(path)
	if err != nil {
		return nil, "", "", cli.UsageError("InvalidConfiguration", err.Error())
	}
	name, explicit := ctx.Flags().GetValue("profile")
	if !explicit {
		name = ctx.Getenv("INSPIRE_PROFILE")
	}
	if name == "" {
		name = c.CurrentProfile
	}
	if name == "" {
		name = "default"
	}
	return c, path, name, nil
}

func loadConfig(ctx *cli.Context) (*sdk.Config, *dara.RuntimeOptions, error) {
	c, _, name, err := configuration(ctx)
	if err != nil {
		return nil, nil, err
	}
	p, found := c.GetProfile(name)
	_, explicit := ctx.Flags().GetValue("profile")
	if !found && (explicit || ctx.Getenv("INSPIRE_PROFILE") != "" || name != "default") {
		return nil, nil, cli.UsageError("InvalidConfiguration", "Named profile does not exist.")
	}
	if value := ctx.Getenv("INSPIRE_ENDPOINT"); value != "" {
		p.Endpoint = value
	}
	if value := ctx.Getenv("INSPIRE_PROXY"); value != "" {
		p.Socks5Proxy = value
	}
	if value := ctx.Getenv("INSPIRE_COOKIE_FILE"); value != "" {
		p.CookieFile = value
	}
	err = p.OverwriteWithFlags(ctx)
	if err != nil {
		return nil, nil, err
	}
	if p.ReadTimeout > int(^uint(0)>>1)/1000 || p.ConnectTimeout > int(^uint(0)>>1)/1000 || p.RetryCount > 99 {
		return nil, nil, cli.UsageError("InvalidConfiguration", "Profile timeout or retry count is out of range.")
	}
	cfg := (&sdk.Config{}).SetEndpoint(p.Endpoint).SetSocks5Proxy(p.Socks5Proxy).SetHttpProxy(p.HttpProxy).SetHttpsProxy(p.HttpsProxy).SetNoProxy(p.NoProxy)
	if p.ReadTimeout > 0 {
		cfg.ReadTimeout = dara.Int(p.ReadTimeout * 1000)
	}
	if p.ConnectTimeout > 0 {
		cfg.ConnectTimeout = dara.Int(p.ConnectTimeout * 1000)
	}
	cookie, token := ctx.Getenv("INSPIRE_COOKIE"), ctx.Getenv("INSPIRE_ACCESS_TOKEN")
	if p.CookieFile != "" {
		f, err := os.Open(p.CookieFile)
		if err != nil {
			return nil, nil, cli.UsageError("InvalidConfiguration", "Cannot open the configured cookie file.")
		}
		data, readErr := io.ReadAll(io.LimitReader(f, 16385))
		f.Close()
		if readErr != nil || len(data) > 16384 {
			return nil, nil, cli.UsageError("InvalidConfiguration", "Cannot read the configured cookie file.")
		}
		cookie = strings.TrimSpace(string(data))
		token = ""
		if cookie == "" {
			return nil, nil, cli.UsageError("InvalidConfiguration", "Configured cookie file is empty.")
		}
	}
	if cookie != "" && token != "" {
		return nil, nil, cli.UsageError("AmbiguousCredential", "Choose INSPIRE_COOKIE or INSPIRE_ACCESS_TOKEN.")
	}
	if cookie != "" {
		cfg.Cookie = dara.String(cookie)
	}
	if token != "" {
		cfg.BearerToken = dara.String(token)
	}
	if p.RetryCount > 0 {
		cfg.RetryOptions = &dara.RetryOptions{Retryable: true, MaxAttempts: p.RetryCount + 1, RetryCondition: []*dara.RetryCondition{{MaxAttempts: p.RetryCount + 1, Exception: []string{"TransportError", "ServerError", "ThrottlingError"}, MaxDelay: 120000, Backoff: &dara.FixedBackoffPolicy{Period: 100}}}}
	}
	runtime := &dara.RuntimeOptions{}
	return cfg, runtime, nil
}

func profileView(p config.Profile, current string) map[string]interface{} {
	proxy := func(raw string) string {
		if raw == "" {
			return ""
		}
		parsed, err := url.Parse(raw)
		if err != nil || parsed.Host == "" {
			return "<configured>"
		}
		return parsed.Scheme + "://" + parsed.Host
	}
	return map[string]interface{}{"name": p.Name, "current": p.Name == current, "endpoint": p.Endpoint, "cookie_file_configured": p.CookieFile != "", "socks5_proxy": proxy(p.Socks5Proxy), "http_proxy": proxy(p.HttpProxy), "https_proxy": proxy(p.HttpsProxy), "read_timeout": p.ReadTimeout, "connect_timeout": p.ConnectTimeout, "retry_count": p.RetryCount}
}

func newConfigureCommand() *cli.Command {
	command := &cli.Command{Name: "configure", Short: "Manage named connection profiles; session credentials stay in environment or files.", Usage: "inspire configure <set|get|list|switch|delete> [--profile NAME]"}
	for _, action := range []string{"set", "get", "list", "switch", "delete"} {
		action := action
		child := &cli.Command{Name: action, Short: "Configure profiles: " + action, Usage: "inspire configure " + action + " [--profile NAME]"}
		child.Run = func(ctx *cli.Context, _ []string) error {
			c, path, name, err := configuration(ctx)
			if err != nil {
				return err
			}
			p, found := c.GetProfile(name)
			switch action {
			case "list":
				views := []map[string]interface{}{}
				for _, profile := range c.Profiles {
					views = append(views, profileView(profile, c.CurrentProfile))
				}
				return cli.WriteJSON(ctx.Stdout(), views)
			case "get":
				if !found {
					return cli.UsageError("InvalidConfiguration", "Named profile does not exist.")
				}
				return cli.WriteJSON(ctx.Stdout(), profileView(p, c.CurrentProfile))
			case "set":
				err = p.OverwriteWithFlags(ctx)
				if err != nil {
					return err
				}
				if p.CookieFile != "" {
					p.CookieFile, err = filepath.Abs(p.CookieFile)
					if err != nil {
						return cli.UsageError("InvalidConfiguration", "Invalid cookie file path.")
					}
				}
				// Validate transport settings without loading or contacting the session.
				check := (&sdk.Config{}).SetEndpoint(p.Endpoint).SetSocks5Proxy(p.Socks5Proxy).SetHttpProxy(p.HttpProxy).SetHttpsProxy(p.HttpsProxy).SetNoProxy(p.NoProxy)
				client, err := sdk.NewClient(check)
				if err != nil {
					return err
				}
				client.Close()
				c.PutProfile(p)
			case "switch":
				if !found {
					return cli.UsageError("InvalidConfiguration", "Named profile does not exist.")
				}
				c.CurrentProfile = name
			case "delete":
				if !found {
					return cli.UsageError("InvalidConfiguration", "Named profile does not exist.")
				}
				profiles := []config.Profile{}
				for _, profile := range c.Profiles {
					if profile.Name != name {
						profiles = append(profiles, profile)
					}
				}
				c.Profiles = profiles
				if c.CurrentProfile == name {
					c.CurrentProfile = "default"
				}
			}
			if err = c.Save(path); err != nil {
				return cli.UsageError("InvalidConfiguration", err.Error())
			}
			return cli.WriteJSON(ctx.Stdout(), map[string]string{"profile": name, "operation": action})
		}
		command.AddSubCommand(child)
	}
	return command
}

func integerFlag(ctx *cli.Context, name string, fallback int) (int, error) {
	raw, ok := ctx.Flags().GetValue(name)
	if !ok {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 || int64(value) > 9223372036 {
		return 0, cli.UsageError("InvalidParameter", fmt.Sprintf("--%s needs a positive integer within the duration limit.", name))
	}
	return value, nil
}
