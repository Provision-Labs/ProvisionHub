package auth

import (
	"errors"
	"net/http"
)

// handleLogin delegates to the auth provider (plugin via gRPC)
func handleLogin(w http.ResponseWriter, r *http.Request) {
	p, err := getProvider()
	if err != nil {
		writeProviderError(w, err)
		return
	}

	reqCtx, err := RequestToHTTPContext(r)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusInternalServerError)
		return
	}

	resp, err := p.StartLogin(r.Context(), StartLoginRequest{
		HTTP:    reqCtx,
		Cookies: RequestCookies(r),
	})
	if err != nil {
		writeProviderError(w, err)
		return
	}

	writeAuthResponse(w, r, int(resp.StatusCode), resp.RedirectURL, resp.SetHeaders, resp.SetCookies)
}

// handleCallback delegates to the auth provider (plugin via gRPC)
func handleCallback(w http.ResponseWriter, r *http.Request) {
	p, err := getProvider()
	if err != nil {
		writeProviderError(w, err)
		return
	}

	reqCtx, err := RequestToHTTPContext(r)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusInternalServerError)
		return
	}

	resp, err := p.HandleCallback(r.Context(), HandleCallbackRequest{
		HTTP:    reqCtx,
		Cookies: RequestCookies(r),
	})
	if err != nil {
		writeProviderError(w, err)
		return
	}

	writeAuthResponse(w, r, int(resp.StatusCode), resp.RedirectURL, resp.SetHeaders, resp.SetCookies)
}

// handleLogout delegates to the auth provider (plugin via gRPC)
func handleLogout(w http.ResponseWriter, r *http.Request) {
	p, err := getProvider()
	if err != nil {
		writeProviderError(w, err)
		return
	}

	reqCtx, err := RequestToHTTPContext(r)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusInternalServerError)
		return
	}

	resp, err := p.Logout(r.Context(), LogoutRequest{
		HTTP:    reqCtx,
		Cookies: RequestCookies(r),
	})
	if err != nil {
		writeProviderError(w, err)
		return
	}

	writeAuthResponse(w, r, int(resp.StatusCode), resp.RedirectURL, resp.SetHeaders, resp.SetCookies)
}

// writeAuthResponse applies provider response to HTTP response
func writeAuthResponse(w http.ResponseWriter, r *http.Request, statusCode int, redirectURL string, headers []Header, cookies []Cookie) {
	ApplyHeaders(w, headers)
	ApplyCookies(w, cookies)

	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	if redirectURL != "" {
		if statusCode < 300 || statusCode > 399 {
			statusCode = http.StatusFound
		}
		http.Redirect(w, r, redirectURL, statusCode)
		return
	}

	w.WriteHeader(statusCode)
}

// writeProviderError handles provider errors
func writeProviderError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrProviderNotConfigured) {
		http.Error(w, "Auth provider not configured", http.StatusServiceUnavailable)
		return
	}

	http.Error(w, "Auth provider request failed", http.StatusBadGateway)
}
