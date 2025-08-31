package commodities

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Update a commodity
// @Description  Updates an existing commodity by ID.
// @Tags         commodities
// @Accept       json
// @Produce      json
// @Param        id        path      int                      true  "Commodity ID"
// @Param        commodity body      UpdateCommodityPayload   true  "Commodity Update Payload"
// @Success      200       {object}  types.Commodity          "Successfully updated commodity"
// @Failure      400       {object}  middleware.ErrorResponse "Bad Request - Invalid input or validation failed" 
// @Failure      404       {object}  middleware.ErrorResponse "Not Found - Commodity not found"
// @Failure      500       {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /commodities/{id} [put]
func Update(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid commodity ID")
		return
	}

	var payload UpdateCommodityPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}

	commodity, found, err := repo.Commodities().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "commodity not found")
		return
	}

	commodity.Name = payload.Name
	commodity.CommodityType = payload.CommodityType

	if err := repo.Commodities().Update(r.Context(), commodity); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(commodity)
}
