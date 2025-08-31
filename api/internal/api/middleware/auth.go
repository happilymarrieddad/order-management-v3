package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	jwtpkg "github.com/happilymarrieddad/order-management-v3/api/internal/jwt"
	"github.com/happilymarrieddad/order-management-v3/api/types"
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
		ctx := AddUserIDToContext(r.Context(), claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext retrieves the user ID from the provided context.
// It returns the user ID as an int64 and a boolean indicating if the ID was found.
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

// AddUserIDToContext adds the user ID to the provided context.
func AddUserIDToContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// CheckAdminAndWriteError checks if the user in the context has the admin role.
// If not, it writes a 403 Forbidden error to the response writer and returns true.
// It also handles cases where the user ID is not found in the context or the user
// is not found in the repository, returning 401 Unauthorized errors.
// Returns true if an error was written, false otherwise.
func CheckAdminAndWriteError(w http.ResponseWriter, r *http.Request) bool {
	userID, found := GetUserIDFromContext(r.Context())
	if !found {
		WriteError(w, http.StatusUnauthorized, "user ID not found in context")
		return true
	}

	repo := GetRepo(r.Context())
	user, found, err := repo.Users().Get(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to get user information")
		return true
	}
	if !found {
		WriteError(w, http.StatusUnauthorized, "user not found")
		return true
	}

	if !user.Roles.HasRole(types.RoleAdmin) {
		WriteError(w, http.StatusForbidden, "only administrators can perform this action")
		return true
	}

	return false
}

// AuthUserAdminRequired is a middleware that checks if the authenticated user has the admin role.
// If the user is not an admin, it writes a 403 Forbidden error.
func AuthUserAdminRequired(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if CheckAdminAndWriteError(w, r) {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthUserAdminRequiredMuxMiddleware returns a mux.MiddlewareFunc that checks if the authenticated user has the admin role.
// If the user is not an admin, it writes a 403 Forbidden error.
func AuthUserAdminRequiredMuxMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return AuthUserAdminRequired(next.ServeHTTP)
	}
}
