package client

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"
)

type Credential struct {
	kind, value string
	expiresAt   time.Time
}

func NewCredential(kind, value string, expiresAt time.Time) (*Credential, error) {
	if kind != "cookie" && kind != "bearer" {
		return nil, NewError("CredentialError", "InvalidCredential", "Unsupported management credential kind.")
	}
	if value == "" || strings.TrimSpace(value) != value {
		return nil, NewError("CredentialError", "InvalidCredential", "Empty or invalid management credential.")
	}
	for _, c := range value {
		if c < 32 || c == 127 || c > 255 || (kind == "bearer" && c > 127) {
			return nil, NewError("CredentialError", "InvalidCredential", "Invalid management credential encoding.")
		}
	}
	return &Credential{kind, value, expiresAt}, nil
}
func (c *Credential) String() string   { return "<management credential>" }
func (c *Credential) GoString() string { return c.String() }
func (c *Credential) valid() bool {
	return c != nil && c.value != "" && (c.expiresAt.IsZero() || time.Now().Before(c.expiresAt))
}

type CredentialProvider interface {
	GetCredential(context.Context) (*Credential, error)
}
type RefreshingProvider interface {
	CredentialProvider
	Refresh(context.Context, *Credential) (*Credential, error)
}
type StaticCredentialProvider struct{ credential *Credential }

func NewStaticCredentialProvider(credential *Credential) *StaticCredentialProvider {
	return &StaticCredentialProvider{credential}
}
func (p *StaticCredentialProvider) GetCredential(ctx context.Context) (*Credential, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return p.credential, nil
}

type EnvironmentCredentialProvider struct{}

func (p *EnvironmentCredentialProvider) GetCredential(ctx context.Context) (*Credential, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	cookie, bearer := os.Getenv("INSPIRE_COOKIE"), os.Getenv("INSPIRE_ACCESS_TOKEN")
	if cookie != "" && bearer != "" {
		return nil, NewError("CredentialError", "AmbiguousCredential", "Choose one management credential source.")
	}
	if cookie != "" {
		return NewCredential("cookie", cookie, time.Time{})
	}
	if bearer != "" {
		return NewCredential("bearer", bearer, time.Time{})
	}
	return nil, NewError("CredentialError", "MissingCredential", "An existing management session is required.")
}

// RefreshingCredentialProvider coalesces concurrent loads. The loader must honor ctx.
type RefreshingCredentialProvider struct {
	mu      sync.Mutex
	cached  *Credential
	loading chan struct{}
	loader  func(context.Context) (*Credential, error)
}

func NewRefreshingCredentialProvider(loader func(context.Context) (*Credential, error)) (*RefreshingCredentialProvider, error) {
	if loader == nil {
		return nil, configError("Credential loader is required.")
	}
	return &RefreshingCredentialProvider{loader: loader}, nil
}
func (p *RefreshingCredentialProvider) GetCredential(ctx context.Context) (*Credential, error) {
	return p.load(ctx, nil, false)
}
func (p *RefreshingCredentialProvider) Refresh(ctx context.Context, rejected *Credential) (*Credential, error) {
	return p.load(ctx, rejected, true)
}
func (p *RefreshingCredentialProvider) load(ctx context.Context, rejected *Credential, force bool) (*Credential, error) {
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		p.mu.Lock()
		if p.cached.valid() && (!force || (rejected != nil && (p.cached.kind != rejected.kind || p.cached.value != rejected.value))) {
			out := p.cached
			p.mu.Unlock()
			return out, nil
		}
		if pending := p.loading; pending != nil {
			p.mu.Unlock()
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-pending:
				continue
			}
		}
		pending := make(chan struct{})
		p.loading = pending
		p.mu.Unlock()
		value, err := invokeLoader(ctx, p.loader)
		if err == nil && !value.valid() {
			err = NewError("CredentialError", "InvalidCredential", "Credential loader returned no valid management session.")
		}
		p.mu.Lock()
		if err == nil {
			p.cached = value
		}
		p.loading = nil
		close(pending)
		p.mu.Unlock()
		return value, err
	}
}

func invokeLoader(ctx context.Context, loader func(context.Context) (*Credential, error)) (value *Credential, err error) {
	defer func() {
		if recover() != nil {
			value = nil
			err = NewError("CredentialError", "CredentialProviderFailed", "Management credential loader failed.")
		}
	}()
	value, err = loader(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, NewError("CredentialError", "CredentialProviderFailed", "Management credential loader failed.")
	}
	return value, nil
}
