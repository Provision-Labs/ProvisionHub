package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func InitStore(secret string) {
	store = sessions.NewCookieStore([]byte(secret))
}

var publicRoutes = map[string]bool{
	"/auth/login":    true,
	"/auth/callback": true,
	"/auth/logout":   true,
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if publicRoutes[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		session, err := store.Get(r, "auth-session")
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}
		token, ok := session.Values["access-token"].(string)
		if !ok || token == "" {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
