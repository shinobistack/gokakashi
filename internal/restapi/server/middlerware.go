package server

import (
	"net/http"
	"strings"
)

type BearerTokenAuth struct {
	AuthToken string
}

func (b *BearerTokenAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		providedToken := strings.TrimPrefix(token, "Bearer ")
		if providedToken != b.AuthToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
