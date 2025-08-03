package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-inventory-system/shared"
)

// AuthMiddleware validates JWT tokens in requests
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for certain paths
		if shouldSkipAuth(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from header
		token, err := shared.ExtractTokenFromHeader(r)
		if err != nil {
			shared.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid or missing token")
			return
		}

		// Validate token
		claims, err := shared.ValidateJWT(token)
		if err != nil {
			shared.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add user info to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// shouldSkipAuth determines if authentication should be skipped for a path
func shouldSkipAuth(path string) bool {
	skipPaths := []string{
		"/auth/login",
		"/auth/register",
		"/metrics",
		"/health",
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}
