package v1

import (
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/addresses"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodities"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodityattributes"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/companies"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/locations"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/users"
)

// AddAuthRoutes configures the v1 API routes on the given subrouter.
func AddAuthRoutes(r *mux.Router) {
	// Gemini order all routes
	addresses.AddRoutes(r)
	commodities.AddRoutes(r)
	commodityattributes.AddRoutes(r)
	companies.AddRoutes(r)
	locations.AddRoutes(r)
	users.AddRoutes(r)
}
