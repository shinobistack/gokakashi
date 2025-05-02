package middleware

import (
	"log"
	"net/http"
	"time"
)

// RequestLogger middleware logs details of every incoming request
type RequestLogger struct{}

// NewRequestLogger creates a new RequestLogger instance
func NewRequestLogger() *RequestLogger {
	return &RequestLogger{}
}

// Middleware implements the middleware interface
func (l *RequestLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Log request details
		log.Printf("[REQUEST] %s %s %s %s",
			r.Method,
			r.URL.Path,
			r.Proto,
			r.RemoteAddr,
		)

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log response details
		log.Printf("[RESPONSE] %s %s completed in %v",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}
