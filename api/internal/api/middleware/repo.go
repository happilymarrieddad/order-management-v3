package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
)

// A private key for context that only this package can access. This is good practice to prevent
// collisions with keys from other packages.
type repoContextKey string

// RepoKey is the key for the GlobalRepo in the context.
const RepoKey repoContextKey = "repo"

// RepoMiddleware creates a new middleware that injects the GlobalRepo into the request context.
func RepoMiddleware(repo repos.GlobalRepo) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), RepoKey, repo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetRepo retrieves the GlobalRepo from the context.
// It panics if the repo is not found, as this indicates a server configuration error.
func GetRepo(ctx context.Context) repos.GlobalRepo {
	return ctx.Value(RepoKey).(repos.GlobalRepo)
}
