package locations

import (
	"net/http" // Added

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// AddRoutes configures the location-related routes on the given subrouter.
// All routes require authentication. POST, PUT, and /find require admin privileges.
func AddRoutes(r *mux.Router) {
	s := r.PathPrefix("/locations").Subrouter()

	// Routes for any authenticated user
	s.HandleFunc("", Create).Methods(http.MethodPost)
	s.HandleFunc("/find", Find).Methods(http.MethodGet)
	s.HandleFunc("/{id:[0-9]+}", Get).Methods(http.MethodGet)
	s.HandleFunc("/{id:[0-9]+}", Update).Methods(http.MethodPut)

	// Routes for admin users only
	adminRouter := s.NewRoute().Subrouter()
	adminRouter.Use(middleware.AuthUserAdminRequiredMuxMiddleware())
	adminRouter.HandleFunc("/{id:[0-9]+}", Delete).Methods(http.MethodDelete)
}
