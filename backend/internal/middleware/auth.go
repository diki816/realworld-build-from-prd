package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/realworld/backend/internal/utils"
)

type contextKey string

const UserContextKey = contextKey("user")

// User represents the authenticated user data stored in context
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Auth returns a middleware that validates JWT tokens
func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			// Parse Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			tokenString := parts[1]
			if tokenString == "" {
				writeError(w, http.StatusUnauthorized, "Token is required")
				return
			}

			// Validate token
			claims, err := utils.ValidateToken(tokenString, secret)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Create user object and add to context
			user := &User{
				ID:       claims.UserID,
				Username: claims.Username,
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext extracts the authenticated user from the request context
func GetUserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(UserContextKey).(*User)
	return user, ok
}

// writeError is a helper function to write JSON error responses
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// Simple JSON error response following RealWorld spec
	w.Write([]byte(`{"errors":{"body":["` + message + `"]}}`))
}