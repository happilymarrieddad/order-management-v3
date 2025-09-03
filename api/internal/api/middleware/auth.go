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

// AuthUserKey is the key for the authenticated user object in the context.
const AuthUserKey contextKey = "ctx:authUser"

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

		claims, err := jwtpkg.ValidateToken(tokenString)
		if err != nil {
			WriteError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// OPTIMIZATION: Fetch the user object once and add it to the context.
		repo := GetRepo(r.Context())
		user, found, err := repo.Users().Get(r.Context(), claims.CompanyID, claims.UserID)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to retrieve user")
			return
		}
		if !found {
			WriteError(w, http.StatusUnauthorized, "user not found")
			return
		}

		// Add the entire user object to the request's context.
		ctx := context.WithValue(r.Context(), AuthUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetAuthUserFromContext retrieves the authenticated user object from the context.
func GetAuthUserFromContext(ctx context.Context) (*types.User, bool) {
	user, ok := ctx.Value(AuthUserKey).(*types.User)
	return user, ok
}

// CheckAdminAndWriteError checks if the user in the context has the admin role.
func CheckAdminAndWriteError(w http.ResponseWriter, r *http.Request) bool {
	authUser, found := GetAuthUserFromContext(r.Context())
	if !found {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return true
	}

	if !authUser.HasRole(types.RoleAdmin) {
		WriteError(w, http.StatusForbidden, "forbidden")
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
