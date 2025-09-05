package products

import (
	"net/http"

	"github.com/gorilla/mux"
)

// AddRoutes configures the product-related routes on the given subrouter.
func AddRoutes(r *mux.Router) {
	// Create a subrouter for the /products resource.
	s := r.PathPrefix("/products").Subrouter()

	// Routes accessible to any authenticated user
	s.HandleFunc("", Create).Methods(http.MethodPost)
	s.HandleFunc("/{id:[0-9]+}", Get).Methods(http.MethodGet)
	s.HandleFunc("/{id:[0-9]+}", Update).Methods(http.MethodPut)
	s.HandleFunc("/find", Find).Methods(http.MethodGet)
	s.HandleFunc("/{id:[0-9]+}", Delete).Methods(http.MethodDelete)

	// Admin-only routes (if any, add here)
	// adminRouter := s.NewRoute().Subrouter()
	// adminRouter.Use(middleware.AuthUserAdminRequiredMuxMiddleware())
}
