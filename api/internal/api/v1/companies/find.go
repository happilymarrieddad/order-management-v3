package companies

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
)

// Find handles listing companies with pagination.
func Find(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	var opts repos.CompanyFindOpts
	// Ignore error, opts will be zero-valued if body is empty or malformed
	_ = json.NewDecoder(r.Body).Decode(&opts)

	// Apply defaults if not provided
	if opts.Limit <= 0 {
		opts.Limit = defaultLimit
	}
	if opts.Offset < 0 {
		opts.Offset = defaultOffset
	}

	companies, count, err := repo.Companies().Find(r.Context(), &opts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := types.NewFindResult(companies, count)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
