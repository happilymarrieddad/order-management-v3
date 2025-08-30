package public

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// GetRoles godoc
// @Summary      Get all roles
// @Description  Retrieves a list of all available user roles.
// @Tags         public
// @Accept       json
// @Produce      json
// @Success      200  {array}   string  "Successfully retrieved roles"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Router       /roles [get]
func GetRoles(w http.ResponseWriter, r *http.Request) {
	roles := types.AllRoles()
	stringRoles := make([]string, len(roles))
	for i, role := range roles {
		stringRoles[i] = role.String()
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stringRoles); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to encode roles")
	}
}

// GetCommodityTypes godoc
// @Summary      Get all commodity types
// @Description  Retrieves a list of all available commodity types.
// @Tags         public
// @Accept       json
// @Produce      json
// @Success      200  {array}   string  "Successfully retrieved commodity types"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Router       /commodity-types [get]
func GetCommodityTypes(w http.ResponseWriter, r *http.Request) {
	commodityTypes := types.AllCommodityTypes()
	stringCommodityTypes := make([]string, len(commodityTypes))
	for i, ct := range commodityTypes {
		stringCommodityTypes[i] = ct.String()
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stringCommodityTypes); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to encode commodity types")
	}
}

// GetOrderStatuses godoc
// @Summary      Get all order statuses
// @Description  Retrieves a list of all available order statuses.
// @Tags         public
// @Accept       json
// @Produce      json
// @Success      200  {array}   string  "Successfully retrieved order statuses"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Router       /order-statuses [get]
func GetOrderStatuses(w http.ResponseWriter, r *http.Request) {
	orderStatuses := types.AllOrderStatuses()
	stringOrderStatuses := make([]string, len(orderStatuses))
	for i, os := range orderStatuses {
		stringOrderStatuses[i] = string(os)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stringOrderStatuses); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to encode order statuses")
	}
}

// AddPublicRoutes registers public (unauthenticated) API routes.
func AddPublicRoutes(r *mux.Router) {
	r.HandleFunc("/roles", GetRoles).Methods("GET")
	r.HandleFunc("/commodity-types", GetCommodityTypes).Methods("GET")
	r.HandleFunc("/order-statuses", GetOrderStatuses).Methods("GET")
}
