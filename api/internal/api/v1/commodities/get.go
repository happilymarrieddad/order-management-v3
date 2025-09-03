package commodities

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// @Summary      Get a commodity by ID
// @Description  Retrieves the details of a single commodity.
// @Tags         commodities
// @Produce      json
// @Param        id  path      int                      true  "Commodity ID"
// @Success      200 {object}  types.Commodity           "Successfully retrieved commodity"
// @Failure      400 {object}  middleware.ErrorResponse "Bad Request - Invalid ID"
// @Failure      404 {object}  middleware.ErrorResponse "Not Found - Commodity not found"
// @Failure      500 {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /commodities/{id} [get]
func Get(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the context (cached by AuthMiddleware).
	_, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid commodity ID")
		return
	}

	commodity, found, err := repo.Commodities().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get commodity")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "commodity not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(commodity)
}
