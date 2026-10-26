package config

import "github.com/Lixsp11/inspire-go-sdk/src/cli"

const (
	ProfileFlagName        = "profile"
	ConfigurePathFlagName  = "config-path"
	EndpointFlagName       = "endpoint"
	CookieFileFlagName     = "cookie-file"
	Socks5ProxyFlagName    = "socks5-proxy"
	HttpProxyFlagName      = "http-proxy"
	HttpsProxyFlagName     = "https-proxy"
	NoProxyFlagName        = "no-proxy"
	ReadTimeoutFlagName    = "read-timeout"
	ConnectTimeoutFlagName = "connect-timeout"
	RetryCountFlagName     = "retry-count"
)

func AddFlags(fs *cli.FlagSet) {
	fs.Add(NewProfileFlag())
	fs.Add(NewConfigurePathFlag())
	fs.Add(NewEndpointFlag())
	fs.Add(NewCookieFileFlag())
	fs.Add(NewSocks5ProxyFlag())
	fs.Add(NewHttpProxyFlag())
	fs.Add(NewHttpsProxyFlag())
	fs.Add(NewNoProxyFlag())
	fs.Add(NewReadTimeoutFlag())
	fs.Add(NewConnectTimeoutFlag())
	fs.Add(NewRetryCountFlag())
}

func NewProfileFlag() *cli.Flag {
	return &cli.Flag{Name: ProfileFlagName, Short: "Named connection profile.", AssignedMode: cli.AssignedOnce, Persistent: true, Category: "caller"}
}
func NewConfigurePathFlag() *cli.Flag {
	return &cli.Flag{Name: ConfigurePathFlagName, Short: "CLI configuration file path.", AssignedMode: cli.AssignedOnce, Persistent: true, Category: "caller"}
}
func NewEndpointFlag() *cli.Flag {
	return &cli.Flag{Name: EndpointFlagName, Short: "Platform origin.", AssignedMode: cli.AssignedOnce, Persistent: true, Category: "caller"}
}
func NewCookieFileFlag() *cli.Flag {
	return &cli.Flag{Name: CookieFileFlagName, Short: "File containing an existing inspire-session cookie.", AssignedMode: cli.AssignedOnce, Persistent: true, Category: "caller"}
}
func NewSocks5ProxyFlag() *cli.Flag {
	return &cli.Flag{Name: Socks5ProxyFlagName, Short: "Explicit SOCKS5 proxy URL.", AssignedMode: cli.AssignedOnce, Persistent: true, Category: "caller"}
}
func NewHttpProxyFlag() *cli.Flag {
	return &cli.Flag{Name: HttpProxyFlagName, Short: "HTTP proxy URL.", AssignedMode: cli.AssignedOnce, Persistent: true, Category: "caller"}
}
func NewHttpsProxyFlag() *cli.Flag {
	return &cli.Flag{Name: HttpsProxyFlagName, Short: "HTTPS proxy URL.", AssignedMode: cli.AssignedOnce, Persistent: true, Category: "caller"}
}
func NewNoProxyFlag() *cli.Flag {
	return &cli.Flag{Name: NoProxyFlagName, Short: "Proxy bypass hosts.", AssignedMode: cli.AssignedOnce, Persistent: true, Category: "caller"}
}
func NewReadTimeoutFlag() *cli.Flag {
	return &cli.Flag{Name: ReadTimeoutFlagName, Short: "Read timeout in seconds.", AssignedMode: cli.AssignedOnce, Persistent: true, Category: "caller"}
}
func NewConnectTimeoutFlag() *cli.Flag {
	return &cli.Flag{Name: ConnectTimeoutFlagName, Short: "Connect timeout in seconds.", AssignedMode: cli.AssignedOnce, Persistent: true, Category: "caller"}
}
func NewRetryCountFlag() *cli.Flag {
	return &cli.Flag{Name: RetryCountFlagName, Short: "Extra read-only attempts (0..99); writes are never replayed.", AssignedMode: cli.AssignedOnce, Persistent: true, Category: "caller"}
}
