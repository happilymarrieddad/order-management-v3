package commodityattributes

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

// @Summary      Find commodity attributes
// @Description  Finds commodity attributes with optional filters and pagination by sending query parameters.
// @Tags         commodity-attributes
// @Accept       json
// @Produce      json
// @Param        limit query int false "Number of records to return"
// @Param        offset query int false "Number of records to skip"
// @Param        id query []int false "Filter by Commodity Attribute IDs"
// @Param        commodity_types query []string false "Filter by Commodity Types (e.g., 'grain', 'fruit')"
// @Success      200  {object}  object{data=[]types.CommodityAttribute,total=int} "A list of commodity attributes"
// @Failure      400  {object}  middleware.ErrorResponse "Bad Request"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /commodity-attributes/find [get]
func Find(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the context (cached by AuthMiddleware).
	_, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

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

	var commodityTypes []types.CommodityType
	for _, ctStr := range r.URL.Query()["commodity_types"] {
		ct, err := types.ParseCommodityType(ctStr)
		if err == nil && ct != types.CommodityTypeUnknown {
			commodityTypes = append(commodityTypes, ct)
		}
	}

	opts := &repos.CommodityAttributeFindOpts{
		Limit:          limit,
		Offset:         offset,
		IDs:            ids,
		CommodityTypes: commodityTypes,
	}

	commodityAttributes, count, err := gr.CommodityAttributes().Find(r.Context(), opts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to find commodity attributes")
		return
	}

	response := types.NewFindResult(commodityAttributes, count)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
