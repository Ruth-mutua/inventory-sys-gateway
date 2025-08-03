package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimiter holds rate limiters for different clients
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rps      rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rps,
		burst:    burst,
	}
}

// getLimiter gets or creates a rate limiter for a client
func (rl *RateLimiter) getLimiter(clientID string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[clientID]
	if !exists {
		limiter = rate.NewLimiter(rl.rps, rl.burst)
		rl.limiters[clientID] = limiter
	}

	return limiter
}

// RateLimitingMiddleware applies rate limiting to requests
func RateLimitingMiddleware(next http.Handler) http.Handler {
	limiter := NewRateLimiter(rate.Limit(10), 30) // 10 requests per second, burst of 30

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := getClientID(r)
		limiter := limiter.getLimiter(clientID)

		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getClientID extracts client identifier from request
func getClientID(r *http.Request) string {
	// Use IP address as client identifier
	// In production, you might want to use API keys or other identifiers
	return r.RemoteAddr
}
