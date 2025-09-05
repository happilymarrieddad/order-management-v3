package commodityattributes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	_ "github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Get a commodity attribute
// @Description  Gets a commodity attribute by its ID.
// @Tags         commodity-attributes
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Commodity Attribute ID"
// @Success      200  {object}  types.CommodityAttribute
// @Failure      400  {object}  middleware.ErrorResponse "Invalid Commodity Attribute ID"
// @Failure      401  {object}  middleware.ErrorResponse "Unauthorized"
// @Failure      404  {object}  middleware.ErrorResponse "Commodity Attribute not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /commodity-attributes/{id} [get]
func Get(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the context (cached by AuthMiddleware).
	_, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	gr := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid commodity attribute ID")
		return
	}

	attr, found, err := gr.CommodityAttributes().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get commodity attribute")
		return
	}

	if !found {
		middleware.WriteError(w, http.StatusNotFound, "commodity attribute not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(attr)
}
