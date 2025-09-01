package companies

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

func Find(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	var opts repos.CompanyFindOpts
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if opts.Limit == 0 {
		opts.Limit = 10
	}

	companies, count, err := repo.Companies().Find(r.Context(), &opts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to find companies")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.FindResult{
		Data:  companies,
		Total: count,
	})
}
