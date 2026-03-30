package auth

import (
	"context"
	"errors"
	"sync"
)

var ErrProviderNotConfigured = errors.New("auth provider is not configured")

type Header struct {
	Key    string
	Values []string
}

type HTTPContext struct {
	Method     string
	Path       string
	Query      string
	Headers    []Header
	RemoteAddr string
	Body       []byte
}

type Cookie struct {
	Name     string
	Value    string
	MaxAge   int32
	Path     string
	Domain   string
	Secure   bool
	HTTPOnly bool
	SameSite string
}

type UserIdentity struct {
	UserID   string
	Username string
	Email    string
	Groups   []string
	Extra    map[string]string
}

type StartLoginRequest struct {
	HTTP    HTTPContext
	Cookies []Cookie
}

type StartLoginResponse struct {
	StatusCode  int32
	RedirectURL string
	SetHeaders  []Header
	SetCookies  []Cookie
}

type HandleCallbackRequest struct {
	HTTP    HTTPContext
	Cookies []Cookie
}

type HandleCallbackResponse struct {
	StatusCode  int32
	RedirectURL string
	SetHeaders  []Header
	SetCookies  []Cookie
	Identity    UserIdentity
}

type LogoutRequest struct {
	HTTP    HTTPContext
	Cookies []Cookie
}

type LogoutResponse struct {
	StatusCode  int32
	RedirectURL string
	SetHeaders  []Header
	SetCookies  []Cookie
}

type AuthenticateRequest struct {
	HTTP    HTTPContext
	Cookies []Cookie
}

type AuthenticateResponse struct {
	Authenticated bool
	Identity      UserIdentity
	Reason        string
}

type HealthResponse struct {
	Ready    bool
	Version  string
	Provider string
	Message  string
}

type AuthProvider interface {
	StartLogin(ctx context.Context, req StartLoginRequest) (StartLoginResponse, error)
	HandleCallback(ctx context.Context, req HandleCallbackRequest) (HandleCallbackResponse, error)
	Logout(ctx context.Context, req LogoutRequest) (LogoutResponse, error)
	Authenticate(ctx context.Context, req AuthenticateRequest) (AuthenticateResponse, error)
	Health(ctx context.Context) (HealthResponse, error)
}

var (
	providerMu sync.RWMutex
	provider   AuthProvider
)

func SetProvider(p AuthProvider) {
	providerMu.Lock()
	defer providerMu.Unlock()
	provider = p
}

func GetProvider() (AuthProvider, error) {
	providerMu.RLock()
	defer providerMu.RUnlock()

	if provider == nil {
		return nil, ErrProviderNotConfigured
	}

	return provider, nil
}

func getProvider() (AuthProvider, error) {
	providerMu.RLock()
	defer providerMu.RUnlock()

	if provider == nil {
		return nil, ErrProviderNotConfigured
	}

	return provider, nil
}
