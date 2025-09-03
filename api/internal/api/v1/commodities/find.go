package commodities

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

// @Summary      Find commodities
// @Description  Finds commodities with optional filters and pagination using query parameters.
// @Tags         commodities
// @Accept       json
// @Produce      json
// @Param        id             query    []int     false  "Commodity IDs"
// @Param        name           query    []string  false  "Commodity names"
// @Param        commodity_type query    int       false  "Commodity type"
// @Param        limit          query    int       false  "Limit"
// @Param        offset         query    int       false  "Offset"
// @Success      200            {object}  types.FindResult{data=[]types.Commodity} "A list of commodities"
// @Failure      400            {object}  middleware.ErrorResponse "Bad Request"
// @Failure      500            {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /commodities/find [get]
func Find(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the context (cached by AuthMiddleware).
	_, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	repo := middleware.GetRepo(r.Context())

	ids, err := utils.GetQueryInt64Slice(r, "id")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	commodityTypeInt, err := utils.GetQueryInt(r, "commodity_type")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid commodity_type format")
		return
	}

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

	opts := &repos.FindCommoditiesOpts{
		IDs:           ids,
		Names:         r.URL.Query()["name"],
		CommodityType: types.CommodityType(commodityTypeInt),
		Limit:         limit,
		Offset:        offset,
	}

	commodities, count, err := repo.Commodities().Find(r.Context(), opts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to find commodities")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.FindResult[*types.Commodity]{
		Data:  commodities,
		Total: count,
	})
}
