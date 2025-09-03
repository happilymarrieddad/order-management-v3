package commodityattributes

import (
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// AddRoutes configures the commodity attribute-related routes on the given subrouter.
// All routes require authentication. POST, PUT, and /find require admin privileges.
func AddRoutes(r *mux.Router) {
	// Create a subrouter for the /commodity-attributes resource.
	s := r.PathPrefix("/commodity-attributes").Subrouter()

	// Routes that require authentication but not admin role.
	// The parent router already applies the AuthMiddleware.
	s.HandleFunc("/{id:[0-9]+}", Get).Methods("GET")
	s.HandleFunc("/find", Find).Methods("GET")

	// Create a subrouter for routes that require admin privileges.
	adminRouter := s.NewRoute().Subrouter()
	adminRouter.Use(middleware.AuthUserAdminRequiredMuxMiddleware())
	adminRouter.HandleFunc("", Create).Methods("POST")
	adminRouter.HandleFunc("/{id:[0-9]+}", Update).Methods("PUT")
}
