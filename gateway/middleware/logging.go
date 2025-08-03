package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs all requests and responses
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture status code
		responseWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process request
		next.ServeHTTP(responseWriter, r)

		// Calculate duration
		duration := time.Since(start)

		// Log request details
		log.Printf(
			"%s %s %s %d %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			responseWriter.statusCode,
			duration,
		)
	})
}
