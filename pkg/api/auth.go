package api

import (
	"net/http"
	"strings"
)

func BearerTokenAuth(next http.Handler, authToken string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		providedToken := strings.TrimPrefix(token, "Bearer ")
		if providedToken != authToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// # Start a scan
//curl -X POST "http://localhost:8000/api/v0/scan?image=nginx:latest&visibility=private" \
//     -H "Authorization: Bearer your_api_token_here"
//# Check scan status
//curl -X GET "http://localhost:8000/api/v0/scan/scan_20230101120000/status" \
//     -H "Authorization: Bearer your_api_token_here"
