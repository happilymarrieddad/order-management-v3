package middleware

import (
	"context"
	"net/http"

	jwtpkg "github.com/happilymarrieddad/order-management-v3/api/internal/jwt"
)

// A private key for context that only this package can access. This is good practice to prevent
// collisions with keys from other packages.
type contextKey string

// UserIDKey is the key for the user ID in the context.
const UserIDKey contextKey = "userID"

// AuthMiddleware checks for a valid JWT in the X-App-Token header.
// If the token is missing or invalid, it returns a 401 Unauthorized error.
// If valid, it adds the user ID from the token's claims to the request context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("X-App-Token")
		if tokenString == "" {
			WriteError(w, http.StatusUnauthorized, "missing token")
			return
		}

		// Validate the token using the dedicated jwt package.
		// This function is expected to parse the token, validate its signature and claims,
		// and return the user ID (from the 'sub' claim) on success.
		claims, err := jwtpkg.ValidateToken(tokenString)
		if err != nil {
			WriteError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// Add the userID to the request's context so it can be accessed by downstream handlers.
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
