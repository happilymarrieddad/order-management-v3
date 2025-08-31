package commodities

import (
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// AddRoutes configures the routes for the commodities resource.
func AddRoutes(r *mux.Router) {
	// Create a subrouter for the /commodities resource.
	commoditiesRouter := r.PathPrefix("/commodities").Subrouter()
	commoditiesRouter.Use(middleware.AuthUserAdminRequiredMuxMiddleware())
	commoditiesRouter.HandleFunc("", Create).Methods("POST")
	commoditiesRouter.HandleFunc("/find", Find).Methods("POST")
	commoditiesRouter.HandleFunc("/{id:[0-9]+}", Get).Methods("GET")
	commoditiesRouter.HandleFunc("/{id:[0-9]+}", Update).Methods("PUT")
}
