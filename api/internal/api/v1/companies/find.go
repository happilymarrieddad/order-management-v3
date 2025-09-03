package companies

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

func Find(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	limit, err := utils.GetQueryInt(r, "limit")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid limit format")
		return
	}
	if limit == 0 {
		limit = 10
	}

	offset, err := utils.GetQueryInt(r, "offset")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid offset format")
		return
	}

	ids, err := utils.GetQueryInt64Slice(r, "id")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	opts := &repos.CompanyFindOpts{
		Limit:  limit,
		Offset: offset,
		IDs:    ids,
		Names:  r.URL.Query()["name"],
	}

	companies, count, err := repo.Companies().Find(r.Context(), opts)
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
