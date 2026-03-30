package middleware

import (
	"context"
	"net/http"

	"github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/auth"
)

type contextUserKey string

const userIdentityKey contextUserKey = "auth.user_identity"

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

		p, err := auth.GetProvider()
		if err != nil {
			http.Error(w, "Auth provider unavailable", http.StatusServiceUnavailable)
			return
		}

		reqCtx, err := auth.RequestToHTTPContext(r)
		if err != nil {
			http.Error(w, "Failed to read request", http.StatusInternalServerError)
			return
		}

		result, err := p.Authenticate(r.Context(), auth.AuthenticateRequest{
			HTTP:    reqCtx,
			Cookies: auth.RequestCookies(r),
		})
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}

		if !result.Authenticated {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), userIdentityKey, result.Identity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIdentityFromContext(ctx context.Context) (auth.UserIdentity, bool) {
	id, ok := ctx.Value(userIdentityKey).(auth.UserIdentity)
	return id, ok
}
