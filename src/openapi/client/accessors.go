package client

import (
	"github.com/alibabacloud-go/tea/dara"
	"net/http"
)

func (s *Config) GetEndpoint() *string                                    { return s.Endpoint }
func (s *Config) SetEndpoint(v string) *Config                            { s.Endpoint = &v; return s }
func (s *Config) GetProtocol() *string                                    { return s.Protocol }
func (s *Config) SetProtocol(v string) *Config                            { s.Protocol = &v; return s }
func (s *Config) GetReadTimeout() *int                                    { return s.ReadTimeout }
func (s *Config) SetReadTimeout(v int) *Config                            { s.ReadTimeout = &v; return s }
func (s *Config) GetConnectTimeout() *int                                 { return s.ConnectTimeout }
func (s *Config) SetConnectTimeout(v int) *Config                         { s.ConnectTimeout = &v; return s }
func (s *Config) GetIdleTimeout() *int                                    { return s.IdleTimeout }
func (s *Config) SetIdleTimeout(v int) *Config                            { s.IdleTimeout = &v; return s }
func (s *Config) GetHttpProxy() *string                                   { return s.HttpProxy }
func (s *Config) SetHttpProxy(v string) *Config                           { s.HttpProxy = &v; return s }
func (s *Config) GetHttpsProxy() *string                                  { return s.HttpsProxy }
func (s *Config) SetHttpsProxy(v string) *Config                          { s.HttpsProxy = &v; return s }
func (s *Config) GetSocks5Proxy() *string                                 { return s.Socks5Proxy }
func (s *Config) SetSocks5Proxy(v string) *Config                         { s.Socks5Proxy = &v; return s }
func (s *Config) GetNoProxy() *string                                     { return s.NoProxy }
func (s *Config) SetNoProxy(v string) *Config                             { s.NoProxy = &v; return s }
func (s *Config) GetMaxIdleConns() *int                                   { return s.MaxIdleConns }
func (s *Config) SetMaxIdleConns(v int) *Config                           { s.MaxIdleConns = &v; return s }
func (s *Config) GetUserAgent() *string                                   { return s.UserAgent }
func (s *Config) SetUserAgent(v string) *Config                           { s.UserAgent = &v; return s }
func (s *Config) GetCa() *string                                          { return s.Ca }
func (s *Config) SetCa(v string) *Config                                  { s.Ca = &v; return s }
func (s *Config) GetKey() *string                                         { return s.Key }
func (s *Config) SetKey(v string) *Config                                 { s.Key = &v; return s }
func (s *Config) GetCert() *string                                        { return s.Cert }
func (s *Config) SetCert(v string) *Config                                { s.Cert = &v; return s }
func (s *Config) GetCookie() *string                                      { return s.Cookie }
func (s *Config) SetCookie(v string) *Config                              { s.Cookie = &v; return s }
func (s *Config) GetBearerToken() *string                                 { return s.BearerToken }
func (s *Config) SetBearerToken(v string) *Config                         { s.BearerToken = &v; return s }
func (s *Config) GetCredential() CredentialProvider                       { return s.Credential }
func (s *Config) SetCredential(v CredentialProvider) *Config              { s.Credential = v; return s }
func (s *Config) GetRetryOptions() *dara.RetryOptions                     { return s.RetryOptions }
func (s *Config) SetRetryOptions(v *dara.RetryOptions) *Config            { s.RetryOptions = v; return s }
func (s *Config) GetGlobalParameters() *GlobalParameters                  { return s.GlobalParameters }
func (s *Config) SetGlobalParameters(v *GlobalParameters) *Config         { s.GlobalParameters = v; return s }
func (s *Config) GetOperationTimeout() *int                               { return s.OperationTimeout }
func (s *Config) SetOperationTimeout(v int) *Config                       { s.OperationTimeout = &v; return s }
func (s *Config) GetMaxResponseBytes() *int64                             { return s.MaxResponseBytes }
func (s *Config) SetMaxResponseBytes(v int64) *Config                     { s.MaxResponseBytes = &v; return s }
func (s *Config) GetTransport() http.RoundTripper                         { return s.Transport }
func (s *Config) SetTransport(v http.RoundTripper) *Config                { s.Transport = v; return s }
func (s *Config) GetObserver() func(Event)                                { return s.Observer }
func (s *Config) SetObserver(v func(Event)) *Config                       { s.Observer = v; return s }
func (s *Params) GetAction() *string                                      { return s.Action }
func (s *Params) SetAction(v string) *Params                              { s.Action = &v; return s }
func (s *Params) GetVersion() *string                                     { return s.Version }
func (s *Params) SetVersion(v string) *Params                             { s.Version = &v; return s }
func (s *Params) GetProtocol() *string                                    { return s.Protocol }
func (s *Params) SetProtocol(v string) *Params                            { s.Protocol = &v; return s }
func (s *Params) GetPathname() *string                                    { return s.Pathname }
func (s *Params) SetPathname(v string) *Params                            { s.Pathname = &v; return s }
func (s *Params) GetMethod() *string                                      { return s.Method }
func (s *Params) SetMethod(v string) *Params                              { s.Method = &v; return s }
func (s *Params) GetAuthType() *string                                    { return s.AuthType }
func (s *Params) SetAuthType(v string) *Params                            { s.AuthType = &v; return s }
func (s *Params) GetStyle() *string                                       { return s.Style }
func (s *Params) SetStyle(v string) *Params                               { s.Style = &v; return s }
func (s *Params) GetReqBodyType() *string                                 { return s.ReqBodyType }
func (s *Params) SetReqBodyType(v string) *Params                         { s.ReqBodyType = &v; return s }
func (s *Params) GetBodyType() *string                                    { return s.BodyType }
func (s *Params) SetBodyType(v string) *Params                            { s.BodyType = &v; return s }
func (s *Params) GetReadOnly() *bool                                      { return s.ReadOnly }
func (s *Params) SetReadOnly(v bool) *Params                              { s.ReadOnly = &v; return s }
func (s *Params) GetSource() *string                                      { return s.Source }
func (s *Params) SetSource(v string) *Params                              { s.Source = &v; return s }
func (s *Params) GetReferer() *string                                     { return s.Referer }
func (s *Params) SetReferer(v string) *Params                             { s.Referer = &v; return s }
func (s *OpenApiRequest) GetHeaders() map[string]*string                  { return s.Headers }
func (s *OpenApiRequest) SetHeaders(v map[string]*string) *OpenApiRequest { s.Headers = v; return s }
func (s *OpenApiRequest) GetQuery() map[string]*string                    { return s.Query }
func (s *OpenApiRequest) SetQuery(v map[string]*string) *OpenApiRequest   { s.Query = v; return s }
func (s *OpenApiRequest) GetBody() interface{}                            { return s.Body }
func (s *OpenApiRequest) SetBody(v interface{}) *OpenApiRequest           { s.Body = v; return s }
func (s *GlobalParameters) GetHeaders() map[string]*string                { return s.Headers }
func (s *GlobalParameters) SetHeaders(v map[string]*string) *GlobalParameters {
	s.Headers = v
	return s
}
func (s *GlobalParameters) GetQueries() map[string]*string { return s.Queries }
func (s *GlobalParameters) SetQueries(v map[string]*string) *GlobalParameters {
	s.Queries = v
	return s
}
func (c *Config) Validate() error {
	if c == nil {
		return configError("Config is required.")
	}
	_, err := resolve(c, nil)
	return err
}
func (p *Params) Validate() error { return validateParams(p) }
