package companies

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

func Find(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	var opts repos.CompanyFindOpts
	query := r.URL.Query()

	// Parse query parameters
	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			opts.Limit = int(limit)
		}
	}
	if opts.Limit == 0 {
		opts.Limit = 10
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		if offset, err := strconv.ParseInt(offsetStr, 10, 64); err == nil {
			opts.Offset = int(offset)
		}
	}

	// Handle multiple IDs
	for _, idStr := range query["id"] {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			opts.IDs = append(opts.IDs, id)
		}
	}

	// Handle multiple names
	for _, name := range query["name"] {
		opts.Names = append(opts.Names, name)
	}

	companies, count, err := repo.Companies().Find(r.Context(), &opts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to find companies")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.FindResult[*types.Company]{
		Data:  companies,
		Total: count,
	})
}
