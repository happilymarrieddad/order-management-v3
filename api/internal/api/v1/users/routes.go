package users

import (
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// AddRoutes configures the user-related routes on the given subrouter.
func AddRoutes(r *mux.Router) {
	// Create a subrouter for the /users resource.
	s := r.PathPrefix("/users").Subrouter()

	// Routes accessible to any authenticated user
	s.HandleFunc("", Create).Methods("POST")
	s.HandleFunc("/{id:[0-9]+}", Get).Methods("GET")

	// Routes for admin users only
	adminRouter := s.NewRoute().Subrouter()
	adminRouter.Use(middleware.AuthUserAdminRequiredMuxMiddleware())
	adminRouter.HandleFunc("/find", Find).Methods("POST")
	adminRouter.HandleFunc("/{id:[0-9]+}", Update).Methods("PUT")
}
