package tests

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	user "github.com/Lixsp11/inspire-go-sdk/src/userv2/client"
	"github.com/alibabacloud-go/tea/dara"
)

func TestExplicitHTTPProxyAndNoProxy(t *testing.T) {
	var proxies, targets atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targets.Add(1)
		_, _ = w.Write(envelope(map[string]interface{}{"id": "direct"}))
	}))
	defer target.Close()
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxies.Add(1)
		if !r.URL.IsAbs() {
			t.Error("HTTP proxy did not receive absolute target URI")
		}
		_, _ = w.Write(envelope(map[string]interface{}{"id": "proxy"}))
	}))
	defer proxy.Close()
	cfg := config(target.URL).SetHttpProxy(proxy.URL)
	client := userClient(t, cfg)
	out, err := client.GetUserDetail(&user.GetUserDetailRequest{})
	if err != nil || dara.StringValue(out.Body.Result.Id) != "proxy" {
		t.Fatalf("proxy failed: %v", err)
	}
	out, err = client.GetUserDetailWithOptions(&user.GetUserDetailRequest{}, &dara.RuntimeOptions{NoProxy: dara.String("127.0.0.1")})
	if err != nil || dara.StringValue(out.Body.Result.Id) != "direct" {
		t.Fatalf("no_proxy failed: %v", err)
	}
	if proxies.Load() != 1 || targets.Load() != 1 {
		t.Fatal("proxy path mismatch")
	}
}

type socksFixture struct {
	listener           net.Listener
	target             string
	username, password string
	refuse             bool
	connections        atomic.Int32
	mu                 sync.Mutex
	domains            []string
	clients            map[net.Conn]bool
	wg                 sync.WaitGroup
}

func newSOCKS(t *testing.T, target string, refuse bool) *socksFixture {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	fixture := &socksFixture{listener: listener, target: target, username: "proxy-user", password: "proxy-password", refuse: refuse, clients: map[net.Conn]bool{}}
	fixture.wg.Add(1)
	go func() {
		defer fixture.wg.Done()
		for {
			client, err := listener.Accept()
			if err != nil {
				return
			}
			fixture.connections.Add(1)
			fixture.mu.Lock()
			fixture.clients[client] = true
			fixture.mu.Unlock()
			fixture.wg.Add(1)
			go func() {
				defer fixture.wg.Done()
				defer func() { fixture.mu.Lock(); delete(fixture.clients, client); fixture.mu.Unlock(); client.Close() }()
				fixture.serve(client)
			}()
		}
	}()
	t.Cleanup(func() {
		listener.Close()
		fixture.mu.Lock()
		for client := range fixture.clients {
			client.Close()
		}
		fixture.mu.Unlock()
		fixture.wg.Wait()
	})
	return fixture
}
func (f *socksFixture) url() string {
	return "socks5h://" + f.username + ":" + f.password + "@" + f.listener.Addr().String()
}
func readN(connection net.Conn, count int) ([]byte, error) {
	data := make([]byte, count)
	_, err := io.ReadFull(connection, data)
	return data, err
}
func (f *socksFixture) serve(client net.Conn) {
	hello, err := readN(client, 2)
	if err != nil {
		return
	}
	if _, err = readN(client, int(hello[1])); err != nil {
		return
	}
	if _, err = client.Write([]byte{5, 2}); err != nil {
		return
	}
	auth, err := readN(client, 2)
	if err != nil {
		return
	}
	name, err := readN(client, int(auth[1]))
	if err != nil {
		return
	}
	size, err := readN(client, 1)
	if err != nil {
		return
	}
	password, err := readN(client, int(size[0]))
	if err != nil {
		return
	}
	if string(name) != f.username || string(password) != f.password || f.refuse {
		_, _ = client.Write([]byte{1, 1})
		return
	}
	_, _ = client.Write([]byte{1, 0})
	request, err := readN(client, 4)
	if err != nil {
		return
	}
	if request[3] != 3 {
		return
	}
	length, err := readN(client, 1)
	if err != nil {
		return
	}
	domain, err := readN(client, int(length[0]))
	if err != nil {
		return
	}
	if _, err = readN(client, 2); err != nil {
		return
	}
	f.mu.Lock()
	f.domains = append(f.domains, string(domain))
	f.mu.Unlock()
	upstream, err := net.Dial("tcp", f.target)
	if err != nil {
		return
	}
	defer upstream.Close()
	_, _ = client.Write([]byte{5, 0, 0, 1, 127, 0, 0, 1, 0, 0})
	done := make(chan struct{})
	go func() { _, _ = io.Copy(upstream, client); upstream.Close(); close(done) }()
	_, _ = io.Copy(client, upstream)
	client.Close()
	<-done
}
func TestSOCKSAuthenticationRemoteDNSAndReuse(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Cookie") != "inspire-session=fixture" {
			t.Error("session lost")
		}
		_, _ = w.Write(envelope(map[string]interface{}{"id": "socks"}))
	}))
	defer target.Close()
	u, _ := url.Parse(target.URL)
	proxy := newSOCKS(t, u.Host, false)
	client := userClient(t, config("http://localhost:9").SetSocks5Proxy(proxy.url()))
	for i := 0; i < 2; i++ {
		result, err := client.GetUserDetail(&user.GetUserDetailRequest{})
		if err != nil {
			t.Fatal(err)
		}
		if dara.StringValue(result.Body.Result.Id) != "socks" {
			t.Fatal("wrong upstream")
		}
	}
	proxy.mu.Lock()
	domains := append([]string(nil), proxy.domains...)
	proxy.mu.Unlock()
	if fmt.Sprint(domains) != "[localhost]" || proxy.connections.Load() != 1 {
		t.Fatalf("DNS/reuse mismatch: %v, connections=%d", domains, proxy.connections.Load())
	}
}
func TestSOCKSRefusalIsStructuredAndBounded(t *testing.T) {
	proxy := newSOCKS(t, "127.0.0.1:9", true)
	client := userClient(t, config("http://localhost:9").SetSocks5Proxy(proxy.url()))
	_, err := client.GetUserDetailWithContext(context.Background(), &user.GetUserDetailRequest{}, retries(2))
	failure := sdkError(t, err, "TransportError")
	if failure.Attempts != 2 || proxy.connections.Load() != 2 {
		t.Fatal("SOCKS failure was not bounded")
	}
	if strings.Contains(err.Error(), proxy.password) || strings.Contains(err.Error(), proxy.username) {
		t.Fatal("proxy credentials leaked")
	}
}
