package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/dara"
)

const Version = "0.1.0"

type GlobalParameters struct {
	Headers map[string]*string `json:"headers,omitempty"`
	Queries map[string]*string `json:"queries,omitempty"`
}
type Config struct {
	Endpoint         *string            `json:"endpoint,omitempty"`
	Protocol         *string            `json:"protocol,omitempty"`
	ReadTimeout      *int               `json:"readTimeout,omitempty"`
	ConnectTimeout   *int               `json:"connectTimeout,omitempty"`
	IdleTimeout      *int               `json:"idleTimeout,omitempty"`
	HttpProxy        *string            `json:"httpProxy,omitempty"`
	HttpsProxy       *string            `json:"httpsProxy,omitempty"`
	Socks5Proxy      *string            `json:"socks5Proxy,omitempty"`
	NoProxy          *string            `json:"noProxy,omitempty"`
	MaxIdleConns     *int               `json:"maxIdleConns,omitempty"`
	UserAgent        *string            `json:"userAgent,omitempty"`
	Ca               *string            `json:"ca,omitempty"`
	Key              *string            `json:"key,omitempty"`
	Cert             *string            `json:"cert,omitempty"`
	Cookie           *string            `json:"-"`
	BearerToken      *string            `json:"-"`
	Credential       CredentialProvider `json:"-"`
	RetryOptions     *dara.RetryOptions `json:"-"`
	GlobalParameters *GlobalParameters  `json:"globalParameters,omitempty"`
	OperationTimeout *int               `json:"operationTimeout,omitempty"`
	MaxResponseBytes *int64             `json:"maxResponseBytes,omitempty"`
	Transport        http.RoundTripper  `json:"-"`
	Observer         func(Event)        `json:"-"`
}

func (c Config) String() string   { return "<Inspire Config: credential and proxy values redacted>" }
func (c Config) GoString() string { return c.String() }
func (c *Config) UnmarshalJSON(data []byte) error {
	type alias Config
	var value alias
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&value); err != nil {
		return configError("Unsupported or invalid Config field.")
	}
	*c = Config(value)
	return nil
}

type Params struct {
	Action      *string `json:"action,omitempty"`
	Version     *string `json:"version,omitempty"`
	Protocol    *string `json:"protocol,omitempty"`
	Pathname    *string `json:"pathname,omitempty"`
	Method      *string `json:"method,omitempty"`
	AuthType    *string `json:"authType,omitempty"`
	Style       *string `json:"style,omitempty"`
	ReqBodyType *string `json:"reqBodyType,omitempty"`
	BodyType    *string `json:"bodyType,omitempty"`
	ReadOnly    *bool   `json:"readOnly,omitempty"`
	Source      *string `json:"source,omitempty"`
	Referer     *string `json:"-"`
}
type OpenApiRequest struct {
	Headers map[string]*string `json:"headers,omitempty"`
	Query   map[string]*string `json:"query,omitempty"`
	Body    interface{}        `json:"body,omitempty"`
}
type ResponseContext struct {
	RequestId string
	Attempts  int
	Action    string
	Source    string
}
type Event struct {
	Action     string
	Attempt    int
	StatusCode int
	RequestId  string
	Code       string
	Elapsed    time.Duration
}
type resolved struct {
	endpoint                    string
	connect, read, budget, idle time.Duration
	maxBytes                    int64
	network                     networkKey
	retry                       *dara.RetryOptions
	headers                     http.Header
	queries                     url.Values
}
type networkKey struct {
	proxy, noProxy, ca, cert, key string
	connect, read, idle           time.Duration
	maxIdle                       int
	ignoreSSL, disableKeepAlive   bool
}

func snapshot(c *Config) *Config {
	out := *c
	source, target := reflect.ValueOf(c).Elem(), reflect.ValueOf(&out).Elem()
	for i := 0; i < source.NumField(); i++ {
		v := source.Field(i)
		if v.Kind() == reflect.Pointer && !v.IsNil() && v.Elem().Kind() != reflect.Struct {
			p := reflect.New(v.Elem().Type())
			p.Elem().Set(v.Elem())
			target.Field(i).Set(p)
		}
	}
	if c.GlobalParameters != nil {
		out.GlobalParameters = &GlobalParameters{cloneMap(c.GlobalParameters.Headers), cloneMap(c.GlobalParameters.Queries)}
	}
	out.RetryOptions = cloneRetry(c.RetryOptions)
	return &out
}
func cloneMap(source map[string]*string) map[string]*string {
	out := map[string]*string{}
	for k, v := range source {
		if v != nil {
			out[k] = dara.String(*v)
		} else {
			out[k] = nil
		}
	}
	return out
}
func pick(value, fallback *string) string {
	if value != nil {
		return *value
	}
	return dara.StringValue(fallback)
}
func millis(value, fallback *int, defaultValue int, allowZero bool) (time.Duration, error) {
	n := defaultValue
	if fallback != nil {
		n = *fallback
	}
	if value != nil {
		n = *value
	}
	if n < 0 || (!allowZero && n == 0) || int64(n) > int64((time.Duration(1<<63-1))/time.Millisecond) {
		return 0, configError("Timeouts must be valid positive integer milliseconds.")
	}
	return time.Duration(n) * time.Millisecond, nil
}
func resolve(c *Config, r *dara.RuntimeOptions) (*resolved, error) {
	if r == nil {
		r = &dara.RuntimeOptions{}
	}
	rv := reflect.ValueOf(r).Elem()
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		name := rt.Field(i).Name
		if name == "LocalAddr" || strings.HasPrefix(name, "WebSocket") {
			if !rv.Field(i).IsNil() {
				return nil, NewError("ConfigurationError", "UnsupportedRuntimeOption", "Runtime option is not implemented for native Action calls.")
			}
		}
	}
	if r.Socks5NetWork != nil && *r.Socks5NetWork != "tcp" {
		return nil, configError("Only TCP SOCKS proxies are supported.")
	}
	endpoint := dara.StringValue(c.Endpoint)
	if endpoint == "" {
		endpoint = "qz.sii.edu.cn"
	}
	if !strings.Contains(endpoint, "://") {
		scheme := dara.StringValue(c.Protocol)
		if scheme == "" {
			scheme = "https"
		}
		endpoint = strings.ToLower(scheme) + "://" + endpoint
	}
	u, err := url.Parse(endpoint)
	if err != nil || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return nil, configError("Endpoint must be an origin without credentials, path or query.")
	}
	if u.Scheme != "https" && !(u.Scheme == "http" && (u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1" || u.Hostname() == "::1")) {
		return nil, configError("HTTPS is required except for literal loopback HTTP endpoints.")
	}
	endpoint = strings.TrimSuffix(endpoint, "/")
	read, err := millis(r.ReadTimeout, c.ReadTimeout, 10000, false)
	if err != nil {
		return nil, err
	}
	connect, err := millis(r.ConnectTimeout, c.ConnectTimeout, 5000, false)
	if err != nil {
		return nil, err
	}
	budget, err := millis(nil, c.OperationTimeout, 60000, false)
	if err != nil {
		return nil, err
	}
	if connect > time.Duration(1<<63-1)-read {
		return nil, configError("Combined timeout is too large.")
	}
	idle, err := millis(r.IdleTimeout, c.IdleTimeout, 0, true)
	if err != nil {
		return nil, err
	}
	proxy := pick(r.Socks5Proxy, c.Socks5Proxy)
	if proxy == "" {
		if u.Scheme == "https" {
			proxy = pick(r.HttpsProxy, c.HttpsProxy)
		} else {
			proxy = pick(r.HttpProxy, c.HttpProxy)
		}
	}
	if proxy != "" {
		p, e := url.Parse(proxy)
		if e != nil || p.Hostname() == "" || (p.Scheme != "http" && p.Scheme != "https" && p.Scheme != "socks5" && p.Scheme != "socks5h") || (p.Path != "" && p.Path != "/") || p.RawQuery != "" || p.Fragment != "" {
			return nil, configError("Invalid explicit proxy URL.")
		}
	}
	maxIdle := 0
	if c.MaxIdleConns != nil {
		maxIdle = *c.MaxIdleConns
	}
	if r.MaxIdleConns != nil {
		maxIdle = *r.MaxIdleConns
	}
	if maxIdle < 0 {
		return nil, configError("max_idle_conns cannot be negative.")
	}
	cap := int64(16 * 1024 * 1024)
	if c.MaxResponseBytes != nil {
		cap = *c.MaxResponseBytes
	}
	if cap <= 0 || cap == int64(1<<63-1) {
		return nil, configError("max_response_bytes must be positive and bounded.")
	}
	retry, err := resolveRetry(c.RetryOptions, r)
	if err != nil {
		return nil, err
	}
	headers := http.Header{}
	queries := url.Values{}
	layers := []*GlobalParameters{c.GlobalParameters}
	if r.ExtendsParameters != nil {
		layers = append(layers, &GlobalParameters{r.ExtendsParameters.Headers, r.ExtendsParameters.Queries})
	}
	for _, layer := range layers {
		if layer == nil {
			continue
		}
		if err = mergeHeaders(headers, layer.Headers); err != nil {
			return nil, err
		}
		if err = mergeQuery(queries, layer.Queries); err != nil {
			return nil, err
		}
	}
	if c.UserAgent != nil && !validHeader(*c.UserAgent) {
		return nil, configError("Invalid user agent.")
	}
	return &resolved{endpoint, connect, read, budget, idle, cap, networkKey{proxy, pick(r.NoProxy, c.NoProxy), pick(r.Ca, c.Ca), pick(r.Cert, c.Cert), pick(r.Key, c.Key), connect, read, idle, maxIdle, dara.BoolValue(r.IgnoreSSL), r.KeepAlive != nil && !*r.KeepAlive}, retry, headers, queries}, nil
}
