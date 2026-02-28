package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyPreview string

		if r.Method != "GET" && r.Method != "HEAD" && r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				bodyPreview = "[error reading body]"
			} else {
				bodyPreview = string(bodyBytes)
			}
			// restore body so the next handler can read it
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		if bodyPreview != "" {
			log.Printf("%s %s body=%s", r.Method, r.URL.Path, bodyPreview)
		} else {
			log.Printf("%s %s", r.Method, r.URL.Path)
		}

		next.ServeHTTP(w, r)
	})
}
