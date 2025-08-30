package commodityattributes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// Get handles retrieving a commodity attribute by ID.
// @Summary      Get a commodity attribute by ID
// @Description  Retrieves a single commodity attribute by its ID.
// @Tags         commodity-attributes
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Commodity Attribute ID"
// @Success      200  {object}  types.CommodityAttribute "Successfully retrieved commodity attribute"
// @Failure      400  {object}  middleware.ErrorResponse "Bad Request - Invalid ID"
// @Failure      401  {object}  middleware.ErrorResponse "Unauthorized - Missing or invalid token"
// @Failure      403  {object}  middleware.ErrorResponse "Forbidden - Insufficient permissions"
// @Failure      404  {object}  middleware.ErrorResponse "Not Found - Commodity attribute not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Router       /commodity-attributes/{id} [get]
func Get(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid commodity attribute ID")
		return
	}

	ca, found, err := repo.CommodityAttributes().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "commodity attribute not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ca)
}
