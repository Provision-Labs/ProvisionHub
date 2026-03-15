package auth

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/login", handleLogin)
	mux.HandleFunc("/auth/callback", handleCallback)
	mux.HandleFunc("/auth/logout", handleLogout)
}
