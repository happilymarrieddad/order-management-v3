package commodities

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Create a new commodity
// @Description  Creates a new commodity.
// @Tags         commodities
// @Accept       json
// @Produce      json
// @Param        commodity body      CreateCommodityPayload    true  "Commodity Creation Payload"
// @Success      201      {object}  types.Commodity           "Successfully created commodity"
// @Failure      400      {object}  middleware.ErrorResponse "Bad Request - Invalid input or validation failed"
// @Failure      500      {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /commodities [post]
func Create(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommodityPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, middleware.FormatValidationErrors(err))
		return
	}

	gr := middleware.GetRepo(r.Context())

	commodity := &types.Commodity{
		Name:          payload.Name,
		CommodityType: payload.CommodityType,
	}

	if err := gr.Commodities().Create(r.Context(), commodity); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to create commodity")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(commodity)
}
