package companies

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

// @Summary      Find companies
// @Description  Finds companies with optional filters and pagination using query parameters.
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        id     query    []int  false  "Company IDs"
// @Param        name   query    []string false "Company names"
// @Param        limit  query    int    false  "Limit"
// @Param        offset query    int    false  "Offset"
// @Success      200    {object}  object{data=[]types.Company,total=int} "A list of companies"
// @Failure      400    {object}  middleware.ErrorResponse "Bad Request"
// @Failure      500    {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /companies/find [get]
func Find(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())

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

	companies, count, err := gr.Companies().Find(r.Context(), opts)
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