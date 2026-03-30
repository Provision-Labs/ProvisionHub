package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	authv1 "github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/gen/authv1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCProvider struct {
	client  authv1.AuthPluginClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

func NewGRPCProvider(addr string, timeout time.Duration, skipTLS bool) (*GRPCProvider, error) {
	if addr == "" {
		return nil, fmt.Errorf("auth plugin address is empty")
	}

	if timeout <= 0 {
		timeout = 2 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	transportCreds := credentials.TransportCredentials(insecure.NewCredentials())
	if !skipTLS {
		transportCreds = credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
	}

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(transportCreds), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("dial auth plugin %s: %w", addr, err)
	}

	return &GRPCProvider{
		client:  authv1.NewAuthPluginClient(conn),
		conn:    conn,
		timeout: timeout,
	}, nil
}

func (g *GRPCProvider) Close() error {
	if g == nil || g.conn == nil {
		return nil
	}
	return g.conn.Close()
}

func (g *GRPCProvider) StartLogin(ctx context.Context, req StartLoginRequest) (StartLoginResponse, error) {
	rpcCtx, cancel := g.withTimeout(ctx)
	defer cancel()

	resp, err := g.client.StartLogin(rpcCtx, &authv1.StartLoginRequest{
		Http:    toProtoHTTP(req.HTTP),
		Cookies: toProtoCookies(req.Cookies),
	})
	if err != nil {
		return StartLoginResponse{}, err
	}

	return StartLoginResponse{
		StatusCode:  resp.GetStatusCode(),
		RedirectURL: resp.GetRedirectUrl(),
		SetHeaders:  fromProtoHeaders(resp.GetSetHeaders()),
		SetCookies:  fromProtoCookies(resp.GetSetCookies()),
	}, nil
}

func (g *GRPCProvider) HandleCallback(ctx context.Context, req HandleCallbackRequest) (HandleCallbackResponse, error) {
	rpcCtx, cancel := g.withTimeout(ctx)
	defer cancel()

	resp, err := g.client.HandleCallback(rpcCtx, &authv1.HandleCallbackRequest{
		Http:    toProtoHTTP(req.HTTP),
		Cookies: toProtoCookies(req.Cookies),
	})
	if err != nil {
		return HandleCallbackResponse{}, err
	}

	return HandleCallbackResponse{
		StatusCode:  resp.GetStatusCode(),
		RedirectURL: resp.GetRedirectUrl(),
		SetHeaders:  fromProtoHeaders(resp.GetSetHeaders()),
		SetCookies:  fromProtoCookies(resp.GetSetCookies()),
		Identity:    fromProtoIdentity(resp.GetIdentity()),
	}, nil
}

func (g *GRPCProvider) Logout(ctx context.Context, req LogoutRequest) (LogoutResponse, error) {
	rpcCtx, cancel := g.withTimeout(ctx)
	defer cancel()

	resp, err := g.client.Logout(rpcCtx, &authv1.LogoutRequest{
		Http:    toProtoHTTP(req.HTTP),
		Cookies: toProtoCookies(req.Cookies),
	})
	if err != nil {
		return LogoutResponse{}, err
	}

	return LogoutResponse{
		StatusCode:  resp.GetStatusCode(),
		RedirectURL: resp.GetRedirectUrl(),
		SetHeaders:  fromProtoHeaders(resp.GetSetHeaders()),
		SetCookies:  fromProtoCookies(resp.GetSetCookies()),
	}, nil
}

func (g *GRPCProvider) Authenticate(ctx context.Context, req AuthenticateRequest) (AuthenticateResponse, error) {
	rpcCtx, cancel := g.withTimeout(ctx)
	defer cancel()

	resp, err := g.client.Authenticate(rpcCtx, &authv1.AuthenticateRequest{
		Http:    toProtoHTTP(req.HTTP),
		Cookies: toProtoCookies(req.Cookies),
	})
	if err != nil {
		return AuthenticateResponse{}, err
	}

	return AuthenticateResponse{
		Authenticated: resp.GetAuthenticated(),
		Identity:      fromProtoIdentity(resp.GetIdentity()),
		Reason:        resp.GetReason(),
	}, nil
}

func (g *GRPCProvider) Health(ctx context.Context) (HealthResponse, error) {
	rpcCtx, cancel := g.withTimeout(ctx)
	defer cancel()

	resp, err := g.client.Health(rpcCtx, &authv1.HealthRequest{})
	if err != nil {
		return HealthResponse{}, err
	}

	return HealthResponse{
		Ready:    resp.GetReady(),
		Version:  resp.GetVersion(),
		Provider: resp.GetProvider(),
		Message:  resp.GetMessage(),
	}, nil
}

func (g *GRPCProvider) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, g.timeout)
}

func toProtoHTTP(ctx HTTPContext) *authv1.HttpContext {
	headers := make([]*authv1.Header, 0, len(ctx.Headers))
	for _, h := range ctx.Headers {
		headers = append(headers, &authv1.Header{Key: h.Key, Values: h.Values})
	}

	return &authv1.HttpContext{
		Method:     ctx.Method,
		Path:       ctx.Path,
		Query:      ctx.Query,
		Headers:    headers,
		RemoteAddr: ctx.RemoteAddr,
		Body:       ctx.Body,
	}
}

func toProtoCookies(cookies []Cookie) []*authv1.CookieKV {
	out := make([]*authv1.CookieKV, 0, len(cookies))
	for _, c := range cookies {
		out = append(out, &authv1.CookieKV{
			Name:     c.Name,
			Value:    c.Value,
			MaxAge:   c.MaxAge,
			Path:     c.Path,
			Domain:   c.Domain,
			Secure:   c.Secure,
			HttpOnly: c.HTTPOnly,
			SameSite: c.SameSite,
		})
	}
	return out
}

func fromProtoHeaders(headers []*authv1.Header) []Header {
	out := make([]Header, 0, len(headers))
	for _, h := range headers {
		if h == nil {
			continue
		}
		out = append(out, Header{Key: h.GetKey(), Values: h.GetValues()})
	}
	return out
}

func fromProtoCookies(cookies []*authv1.CookieKV) []Cookie {
	out := make([]Cookie, 0, len(cookies))
	for _, c := range cookies {
		if c == nil {
			continue
		}
		out = append(out, Cookie{
			Name:     c.GetName(),
			Value:    c.GetValue(),
			MaxAge:   c.GetMaxAge(),
			Path:     c.GetPath(),
			Domain:   c.GetDomain(),
			Secure:   c.GetSecure(),
			HTTPOnly: c.GetHttpOnly(),
			SameSite: c.GetSameSite(),
		})
	}
	return out
}

func fromProtoIdentity(identity *authv1.UserIdentity) UserIdentity {
	if identity == nil {
		return UserIdentity{}
	}

	return UserIdentity{
		UserID:   identity.GetUserId(),
		Username: identity.GetUsername(),
		Email:    identity.GetEmail(),
		Groups:   identity.GetGroups(),
		Extra:    identity.GetExtra(),
	}
}
