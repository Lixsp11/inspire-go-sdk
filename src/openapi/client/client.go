// Package client implements the native Inspire OpenAPI transport shared by service clients.
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alibabacloud-go/tea/dara"
	"golang.org/x/net/http/httpproxy"
)

var identifier = regexp.MustCompile(`^[A-Za-z0-9_.:-]{1,128}$`)
var actionName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
var apiPath = regexp.MustCompile(`^/api/v2/[A-Za-z0-9_-]+$`)
var headerName = regexp.MustCompile("^[!#$%&'*+.^_`|~0-9A-Za-z-]+$")

type Client struct {
	mu         sync.Mutex
	config     *Config
	credential CredentialProvider
	pools      map[networkKey]*http.Transport
	closed     bool
}

func NewClient(config *Config) (*Client, error) {
	client := &Client{}
	if err := client.Init(config); err != nil {
		return nil, err
	}
	return client, nil
}
func (c *Client) Init(config *Config) error {
	if config == nil {
		return configError("Config is required.")
	}
	value := snapshot(config)
	if value.Transport != nil && nilValue(value.Transport) {
		return configError("Transport cannot be a typed nil.")
	}
	if _, err := resolve(value, nil); err != nil {
		return err
	}
	n := 0
	if value.Cookie != nil {
		n++
	}
	if value.BearerToken != nil {
		n++
	}
	if value.Credential != nil {
		n++
	}
	if n > 1 {
		return configError("Choose credential provider, Cookie or BearerToken.")
	}
	provider := value.Credential
	if provider != nil {
		if nilValue(provider) {
			return configError("Credential provider cannot be nil.")
		}
	}
	if value.Cookie != nil || value.BearerToken != nil {
		kind, secret := "cookie", dara.StringValue(value.Cookie)
		if value.BearerToken != nil {
			kind, secret = "bearer", *value.BearerToken
		}
		cred, err := NewCredential(kind, secret, time.Time{})
		if err != nil {
			return err
		}
		provider = NewStaticCredentialProvider(cred)
	}
	if provider == nil {
		provider = &EnvironmentCredentialProvider{}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.config != nil {
		return configError("Client is already initialized.")
	}
	c.config, c.credential, c.pools = value, provider, map[networkKey]*http.Transport{}
	return nil
}

// Close prevents new operations and closes owned idle connections. In-flight calls retain their contexts.
// A Transport supplied by the caller is never closed by this client.
func (c *Client) Close() error {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	for _, transport := range c.pools {
		transport.CloseIdleConnections()
	}
	return nil
}
func (c *Client) state() (*Config, CredentialProvider, error) {
	if c == nil {
		return nil, nil, configError("Client is nil.")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil, nil, NewError("ClientClosedError", "ClientClosed", "Client is closed.")
	}
	if c.config == nil {
		return nil, nil, configError("Use NewClient or Init before calling an API.")
	}
	return c.config, c.credential, nil
}
func (c *Client) transport(config *Config, key networkKey) (http.RoundTripper, error) {
	if config.Transport != nil {
		return config.Transport, nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil, NewError("ClientClosedError", "ClientClosed", "Client is closed.")
	}
	if value := c.pools[key]; value != nil {
		return value, nil
	}
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: key.ignoreSSL} // Explicit RuntimeOptions.IgnoreSSL only.
	if key.ca != "" {
		pool, err := x509.SystemCertPool()
		if err != nil {
			pool = x509.NewCertPool()
		}
		if !pool.AppendCertsFromPEM([]byte(key.ca)) {
			return nil, configError("Ca contains no valid PEM certificates.")
		}
		tlsConfig.RootCAs = pool
	}
	if key.key != "" && key.cert == "" {
		return nil, configError("Key requires Cert.")
	}
	if key.cert != "" {
		cert, err := tls.X509KeyPair([]byte(key.cert), []byte(key.key))
		if err != nil {
			return nil, configError("Client certificate/key PEM cannot be loaded.")
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	transport := &http.Transport{DialContext: (&net.Dialer{Timeout: key.connect, KeepAlive: 30 * time.Second}).DialContext,
		TLSClientConfig: tlsConfig, TLSHandshakeTimeout: key.connect, ResponseHeaderTimeout: key.read,
		MaxIdleConns: key.maxIdle, MaxIdleConnsPerHost: key.maxIdle, IdleConnTimeout: key.idle, DisableKeepAlives: key.disableKeepAlive}
	if key.proxy != "" {
		proxyURL, _ := url.Parse(key.proxy)
		transport.Proxy = http.ProxyURL(proxyURL)
		if key.noProxy != "" {
			selector := (&httpproxy.Config{HTTPProxy: key.proxy, HTTPSProxy: key.proxy, NoProxy: key.noProxy}).ProxyFunc()
			transport.Proxy = func(req *http.Request) (*url.URL, error) { return selector(req.URL) }
		}
	}
	c.pools[key] = transport
	return transport, nil
}
func validHeader(value string) bool {
	if value != "" && (value[0] == ' ' || value[0] == '\t') {
		return false
	}
	for _, r := range value {
		if r < 32 || r == 127 || r > 255 {
			return false
		}
	}
	return true
}
func mergeHeaders(target http.Header, values map[string]*string) error {
	seen := map[string]bool{}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := values[key]
		canonical := http.CanonicalHeaderKey(key)
		if !headerName.MatchString(key) || value == nil || !validHeader(*value) {
			return configError("Invalid custom HTTP header.")
		}
		if seen[canonical] {
			return configError("Duplicate custom header names differ only by casing.")
		}
		seen[canonical] = true
		switch strings.ToLower(key) {
		case "host", "content-length", "transfer-encoding", "connection", "proxy-authorization", "proxy-connection", "idempotency-key", "x-idempotency-key":
			return configError("Custom header conflicts with native routing or replay protection.")
		}
		target.Set(canonical, *value)
	}
	return nil
}
func mergeQuery(target url.Values, values map[string]*string) error {
	for key, value := range values {
		if strings.EqualFold(key, "Action") || value == nil {
			return configError("Custom query cannot override Action and must contain string values.")
		}
		target.Set(key, *value)
	}
	return nil
}
func validateParams(p *Params) error {
	if p == nil || !actionName.MatchString(dara.StringValue(p.Action)) {
		return ValidationError("A valid Action is required.")
	}
	public := dara.StringValue(p.Pathname) == "/discovery" && dara.StringValue(p.Method) == "GET" && dara.StringValue(p.AuthType) == "Anonymous"
	if !public && (!apiPath.MatchString(dara.StringValue(p.Pathname)) || dara.StringValue(p.Method) != "POST") {
		return ValidationError("Native Action calls require POST and a platform /api/v2 service path.")
	}
	if p.AuthType != nil && *p.AuthType != "SSO" && *p.AuthType != "Anonymous" {
		return ValidationError("Unsupported native authentication type.")
	}
	for _, kind := range []*string{p.ReqBodyType, p.BodyType} {
		if kind != nil && *kind != "json" {
			return ValidationError("Native Action calls use JSON bodies.")
		}
	}
	if p.Style != nil && *p.Style != "ROA" {
		return ValidationError("Native Action calls do not implement cloud RPC signing.")
	}
	if p.Version != nil && *p.Version != "v2" {
		return ValidationError("Unsupported native API version.")
	}
	if p.Protocol != nil && *p.Protocol != "HTTPS" && *p.Protocol != "HTTP" {
		return ValidationError("Unsupported protocol.")
	}
	if p.Referer != nil && (!strings.HasPrefix(*p.Referer, "/") || strings.HasPrefix(*p.Referer, "//") || !validHeader(*p.Referer)) {
		return ValidationError("Invalid native Referer path.")
	}
	return nil
}

type validatedModel interface{ Validate() error }
type diagnosticResponse interface {
	SetDiagnostics([]byte, *ResponseContext)
}

func (c *Client) CallApi(params *Params, request *OpenApiRequest, runtime *dara.RuntimeOptions) (map[string]interface{}, error) {
	return c.CallApiWithContext(context.Background(), params, request, runtime)
}
func (c *Client) CallApiWithContext(ctx context.Context, params *Params, request *OpenApiRequest, runtime *dara.RuntimeOptions) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.call(ctx, params, request, runtime, func(data []byte, headers map[string]*string, status int32) error {
		var body interface{}
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.UseNumber()
		if err := decoder.Decode(&body); err != nil {
			return responseError("InvalidJson")
		}
		result = map[string]interface{}{"headers": headers, "statusCode": status, "body": body}
		return nil
	}, nil)
	return result, err
}

// CallTyped is the shared execution path used by service Clients. Parsing happens inside the attempt.
func nilValue(value interface{}) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Func, reflect.Map, reflect.Slice, reflect.Chan:
		return v.IsNil()
	}
	return false
}

func (c *Client) CallTyped(ctx context.Context, params *Params, request validatedModel, runtime *dara.RuntimeOptions, response validatedModel) error {
	if nilValue(request) || nilValue(response) || reflect.TypeOf(response).Kind() != reflect.Pointer || reflect.TypeOf(response).Elem().Kind() != reflect.Struct {
		return ValidationError("Typed request and response models are required.")
	}
	if err := request.Validate(); err != nil {
		return localError(err)
	}
	return c.call(ctx, params, &OpenApiRequest{Body: request}, runtime, func(data []byte, headers map[string]*string, status int32) error {
		value, err := json.Marshal(map[string]interface{}{"headers": headers, "statusCode": status, "body": json.RawMessage(data)})
		if err != nil {
			return responseError("InvalidJson")
		}
		if err = json.Unmarshal(value, response); err != nil {
			return responseError("InvalidResponseModel")
		}
		if err = response.Validate(); err != nil {
			var sdk *SDKError
			if errors.As(err, &sdk) {
				return sdk
			}
			return responseError("InvalidResponseModel")
		}
		return nil
	}, func(data []byte, details *ResponseContext) {
		if target, ok := response.(diagnosticResponse); ok {
			target.SetDiagnostics(data, details)
		}
	})
}
func loadCredential(ctx context.Context, provider CredentialProvider) (*Credential, error) {
	value, err := provider.GetCredential(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return nil, timeError(ctx)
		}
		var sdk *SDKError
		if errors.As(err, &sdk) {
			return nil, sdk
		}
		return nil, NewError("CredentialError", "CredentialProviderFailed", "Management credential provider failed.")
	}
	if !value.valid() {
		return nil, NewError("CredentialError", "InvalidCredential", "Management credential is missing or expired.")
	}
	return value, nil
}

func (c *Client) call(parent context.Context, p *Params, request *OpenApiRequest, runtime *dara.RuntimeOptions, parse func([]byte, map[string]*string, int32) error, attach func([]byte, *ResponseContext)) error {
	if parent == nil {
		return configError("Context cannot be nil.")
	}
	config, provider, err := c.state()
	if err != nil {
		return err
	}
	if err = validateParams(p); err != nil {
		return err
	}
	if request == nil {
		return ValidationError("OpenApiRequest is required.")
	}
	resolved, err := resolve(config, runtime)
	if err != nil {
		return err
	}
	if config.Protocol == nil && config.Endpoint != nil && !strings.Contains(*config.Endpoint, "://") && p.Protocol != nil {
		copyConfig := *config
		copyConfig.Protocol = dara.String(strings.ToLower(*p.Protocol))
		resolved, err = resolve(&copyConfig, runtime)
		if err != nil {
			return err
		}
	}
	if err = mergeHeaders(resolved.headers, request.Headers); err != nil {
		return err
	}
	if err = mergeQuery(resolved.queries, request.Query); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(parent, resolved.budget)
	defer cancel()
	if ctx.Err() != nil {
		return timeError(ctx)
	}
	body := request.Body
	if body == nil {
		body = map[string]interface{}{}
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return ValidationError("Request body must be finite native JSON.")
	}
	method := dara.StringValue(p.Method)
	authenticated := dara.StringValue(p.AuthType) != "Anonymous"
	readOnly := dara.BoolValue(p.ReadOnly)
	var credential *Credential
	if authenticated {
		credential, err = loadCredential(ctx, provider)
		if err != nil {
			return err
		}
	}
	if ctx.Err() != nil {
		return timeError(ctx)
	}
	transport, err := c.transport(config, resolved.network)
	if err != nil {
		return err
	}
	client := &http.Client{Transport: transport, Timeout: resolved.connect + resolved.read, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	queries := resolved.queries
	if method == "POST" {
		queries.Set("Action", *p.Action)
	}
	target := resolved.endpoint + *p.Pathname
	if len(queries) > 0 {
		target += "?" + queries.Encode()
	}
	refreshed := false
	start := time.Now()
	for attempt := 1; ; attempt++ {
		if ctx.Err() != nil {
			return timeError(ctx)
		}
		if _, _, err = c.state(); err != nil {
			return err
		}
		if credential != nil && !credential.valid() {
			credential, err = loadCredential(ctx, provider)
			if err != nil {
				return err
			}
		}
		headers := http.Header{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}, "User-Agent": []string{"github.com/Lixsp11/inspire-go-sdk/" + Version}, "Origin": []string{resolved.endpoint}}
		if config.UserAgent != nil {
			headers.Set("User-Agent", headers.Get("User-Agent")+" "+*config.UserAgent)
		}
		for key, value := range resolved.headers {
			headers[key] = append([]string(nil), value...)
		}
		if credential != nil {
			headers.Del("Cookie")
			headers.Del("Authorization")
			if credential.kind == "cookie" {
				headers.Set("Cookie", credential.value)
			} else {
				headers.Set("Authorization", "Bearer "+credential.value)
			}
		}
		if p.Referer != nil {
			headers.Set("Referer", resolved.endpoint+*p.Referer)
		}
		var reader io.Reader
		if method == "POST" {
			reader = bytes.NewReader(encoded)
		}
		httpRequest, e := http.NewRequestWithContext(ctx, method, target, reader)
		if e != nil {
			return configError("HTTP request configuration is invalid.")
		}
		httpRequest.GetBody = nil // Prevent transport-level replay of native write bodies.
		httpRequest.Header = headers
		response, sendErr := client.Do(httpRequest)
		status, requestId := 0, ""
		var raw []byte
		var failure *SDKError
		serverDelay := time.Duration(0)
		if sendErr != nil {
			if ctx.Err() != nil {
				failure = timeError(ctx)
			} else {
				failure = NewError("TransportError", "HttpTransportError", "HTTP connection or response read failed.")
				if isTLSError(sendErr) {
					failure = NewError("TlsError", "TlsError", "TLS negotiation or certificate verification failed.")
				}
			}
			if response != nil && response.Body != nil {
				response.Body.Close()
			}
		} else {
			status = response.StatusCode
			requestId = safeIdentifier(response.Header.Get("X-Request-Id"))
			serverDelay = retryAfter(response.Header.Get("Retry-After"))
			raw, e = io.ReadAll(io.LimitReader(response.Body, resolved.maxBytes+1))
			response.Body.Close()
			if e != nil {
				if ctx.Err() != nil {
					failure = timeError(ctx)
				} else {
					failure = NewError("TransportError", "HttpTransportError", "Response read failed.")
				}
			} else if int64(len(raw)) > resolved.maxBytes {
				failure = responseError("ResponseTooLarge")
			} else {
				var decoded map[string]json.RawMessage
				failure, requestId = decodeResponse(raw, status, requestId, &decoded)
				if failure == nil {
					if e = parse(raw, lowerHeaders(response.Header), int32(status)); e != nil {
						failure = localError(e)
					}
				}
			}
		}
		if failure == nil && ctx.Err() != nil {
			failure = timeError(ctx)
		}
		if failure != nil {
			errorDiagnostic(failure, *p.Action, attempt, status, requestId)
			if serverDelay > 0 {
				failure.retryAfter = dara.Int64(serverDelay.Milliseconds())
			}
		}
		if config.Observer != nil {
			code := ""
			if failure != nil {
				code = dara.StringValue(failure.Code)
			}
			notify(config.Observer, Event{*p.Action, attempt, status, requestId, code, time.Since(start)})
		}
		if failure == nil {
			if attach != nil {
				attach(append([]byte(nil), raw...), &ResponseContext{requestId, attempt, *p.Action, dara.StringValue(p.Source)})
			}
			return nil
		}
		kind := dara.StringValue(failure.Name)
		if !readOnly {
			switch kind {
			case "TransportError", "ServerError", "ThrottlingError", "InvalidResponseError", "OperationTimeoutError", "CanceledError":
				return uncertain(failure)
			}
			return failure
		}
		eligible := kind == "TransportError" || kind == "ServerError" || kind == "ThrottlingError" || kind == "AuthenticationError"
		if !eligible {
			return failure
		}
		retry, delay, retryErr := retryDecision(resolved.retry, failure, attempt)
		if retryErr != nil {
			return policyError(retryErr)
		}
		if !retry {
			return failure
		}
		if kind == "AuthenticationError" {
			refresher, ok := provider.(RefreshingProvider)
			if !ok || refreshed {
				return failure
			}
			refreshed = true
			updated, refreshErr := refresher.Refresh(ctx, credential)
			if refreshErr != nil || !updated.valid() {
				if ctx.Err() != nil {
					return timeError(ctx)
				}
				return NewError("CredentialError", "CredentialRefreshFailed", "Management session refresh failed.")
			}
			if updated.kind == credential.kind && updated.value == credential.value {
				return failure
			}
			credential = updated
		}
		if serverDelay > delay {
			delay = serverDelay
		}
		if e = wait(ctx, delay); e != nil {
			out := timeError(ctx)
			errorDiagnostic(out, *p.Action, attempt, status, requestId)
			out.retryAfter = failure.retryAfter
			return out
		}
	}
}

func isTLSError(err error) bool {
	var certificate *tls.CertificateVerificationError
	var record tls.RecordHeaderError
	var alert tls.AlertError
	if errors.As(err, &certificate) || errors.As(err, &record) || errors.As(err, &alert) {
		return true
	}
	// crypto/tls also uses private alert and plain error types for negotiation failures.
	for current := err; current != nil; current = errors.Unwrap(current) {
		message := current.Error()
		if strings.HasPrefix(message, "tls:") || strings.HasPrefix(message, "remote error: tls:") {
			return true
		}
	}
	return false
}
func notify(observer func(Event), event Event) { defer func() { _ = recover() }(); observer(event) }
func safeIdentifier(value string) string {
	if identifier.MatchString(value) {
		return value
	}
	return ""
}
func lowerHeaders(headers http.Header) map[string]*string {
	out := map[string]*string{}
	for key, value := range headers {
		if len(value) > 0 {
			out[strings.ToLower(key)] = dara.String(value[0])
		}
	}
	return out
}
func decodeResponse(raw []byte, status int, requestId string, out *map[string]json.RawMessage) (*SDKError, string) {
	decodeErr := json.Unmarshal(raw, out)
	var metadata struct {
		RequestId string `json:"RequestId"`
		Error     *struct {
			Code string `json:"Code"`
		} `json:"Error"`
	}
	if decodeErr == nil {
		_ = json.Unmarshal((*out)["ResponseMetadata"], &metadata)
	}
	if value := safeIdentifier(metadata.RequestId); value != "" {
		requestId = value
	}
	code := ""
	if metadata.Error != nil {
		code = safeIdentifier(metadata.Error.Code)
	}
	normalized := strings.ReplaceAll(strings.ToLower(code), "_", "")
	kind := ""
	switch {
	case status == 401 || normalized == "unauthorized" || normalized == "invalidtoken" || normalized == "tokenexpired" || normalized == "unauthorizedexception":
		kind = "AuthenticationError"
	case status >= 300 && status < 400:
		kind = "RedirectError"
	case status == 403 || normalized == "forbidden" || normalized == "accessdenied" || normalized == "accessforbidden":
		kind = "PermissionDeniedError"
	case status == 429 || normalized == "slowdown" || normalized == "throttling" || normalized == "throttlingexception" || normalized == "toomanyrequests" || normalized == "toomanyrequestsexception":
		kind = "ThrottlingError"
	case status == 500 || status == 502 || status == 503 || status == 504 || normalized == "internalerror" || normalized == "internalfailure" || normalized == "internalservererror" || normalized == "requesttimeout" || normalized == "serviceunavailable":
		kind = "ServerError"
	case status < 200 || status >= 300 || metadata.Error != nil:
		kind = "PlatformError"
	}
	if kind != "" {
		if code == "" {
			code = "Http" + dara.Stringify(status)
		}
		return NewError(kind, code, "Platform request was rejected."), requestId
	}
	if decodeErr != nil || *out == nil {
		return responseError("InvalidJson"), requestId
	}
	if _, ok := (*out)["Result"]; !ok {
		return responseError("MissingResult"), requestId
	}
	return nil, requestId
}
