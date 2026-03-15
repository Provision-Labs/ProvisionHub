package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	provider *oidc.Provider
	oauthCfg *oauth2.Config
	verifier *oidc.IDTokenVerifier
)

var store *sessions.CookieStore

var cfg *config.Config

func InitStore(secret string) {
	store = sessions.NewCookieStore([]byte(secret))
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func Init(c *config.Config) {
	cfg = c

	provider, oauthCfg = setupOAuth(cfg)
	verifier = provider.Verifier(&oidc.Config{ClientID: cfg.ClientId})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := generateState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, oauthCfg.AuthCodeURL(state), http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	expectedState, ok := session.Values["state"].(string)
	if !ok || expectedState != r.URL.Query().Get("state") {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}
	delete(session.Values, "state")

	token, err := oauthCfg.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth token", http.StatusInternalServerError)
		return
	}

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Token verify failed", http.StatusInternalServerError)
		return
	}

	var claims struct {
		Sub           string   `json:"sub"`
		Email         string   `json:"email"`
		PreferredName string   `json:"preferred_username"`
		RealmRoles    []string `json:"realm_roles.roles"`
	}
	err = idToken.Claims(&claims)
	if err != nil {
		http.Error(w, "Token claims failed", http.StatusInternalServerError)
		return
	}

	session.Values["access-token"] = token.AccessToken
	session.Values["id_token"] = rawIDToken
	session.Values["username"] = claims.PreferredName
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	idToken, _ := session.Values["id_token"].(string)

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	logout := fmt.Sprintf(
		"%s/protocol/openid-connect/logout?id_token_hint=%s&post_logout_redirect_uri=%s",
		cfg.Issuer,
		idToken,
		url.QueryEscape(cfg.LogoutRedirect),
	)
	http.Redirect(w, r, logout, http.StatusFound)
}
