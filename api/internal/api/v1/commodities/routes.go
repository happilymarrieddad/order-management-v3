package commodities

import (
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// AddRoutes configures the commodity-related routes on the given subrouter.
// All routes require authentication. POST, PUT, and /find require admin privileges.
func AddRoutes(r *mux.Router) {
	s := r.PathPrefix("/commodities").Subrouter()

	// Routes for any authenticated user
	s.HandleFunc("/{id:[0-9]+}", Get).Methods("GET")
	s.HandleFunc("/find", Find).Methods("GET")

	// Routes for admin users only
	adminRouter := s.NewRoute().Subrouter()
	adminRouter.Use(middleware.AuthUserAdminRequiredMuxMiddleware())
	adminRouter.HandleFunc("", Create).Methods("POST")
	adminRouter.HandleFunc("/{id:[0-9]+}", Update).Methods("PUT")
}
