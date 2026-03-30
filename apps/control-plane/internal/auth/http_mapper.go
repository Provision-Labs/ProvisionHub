package auth

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

func RequestToHTTPContext(r *http.Request) (HTTPContext, error) {
	ctx := HTTPContext{
		Method:     r.Method,
		Path:       r.URL.Path,
		Query:      r.URL.RawQuery,
		Headers:    make([]Header, 0, len(r.Header)),
		RemoteAddr: r.RemoteAddr,
	}

	for key, values := range r.Header {
		copied := append([]string(nil), values...)
		ctx.Headers = append(ctx.Headers, Header{Key: key, Values: copied})
	}

	if r.Body == nil {
		return ctx, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return HTTPContext{}, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	ctx.Body = body

	return ctx, nil
}

func RequestCookies(r *http.Request) []Cookie {
	raw := r.Cookies()
	out := make([]Cookie, 0, len(raw))

	for _, c := range raw {
		out = append(out, Cookie{
			Name:     c.Name,
			Value:    c.Value,
			MaxAge:   int32(c.MaxAge),
			Path:     c.Path,
			Domain:   c.Domain,
			Secure:   c.Secure,
			HTTPOnly: c.HttpOnly,
			SameSite: sameSiteToString(c.SameSite),
		})
	}

	return out
}

func ApplyHeaders(w http.ResponseWriter, headers []Header) {
	for _, h := range headers {
		for _, value := range h.Values {
			w.Header().Add(h.Key, value)
		}
	}
}

func ApplyCookies(w http.ResponseWriter, cookies []Cookie) {
	for _, c := range cookies {
		http.SetCookie(w, &http.Cookie{
			Name:     c.Name,
			Value:    c.Value,
			MaxAge:   int(c.MaxAge),
			Path:     c.Path,
			Domain:   c.Domain,
			Secure:   c.Secure,
			HttpOnly: c.HTTPOnly,
			SameSite: parseSameSite(c.SameSite),
		})
	}
}

func sameSiteToString(s http.SameSite) string {
	switch s {
	case http.SameSiteStrictMode:
		return "Strict"
	case http.SameSiteNoneMode:
		return "None"
	case http.SameSiteLaxMode:
		return "Lax"
	default:
		return ""
	}
}

func parseSameSite(raw string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	case "lax":
		return http.SameSiteLaxMode
	default:
		return http.SameSiteDefaultMode
	}
}
