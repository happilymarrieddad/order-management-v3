package companies

import (
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// AddRoutes configures the company-related routes on the given subrouter.
func AddRoutes(r *mux.Router) {
	// Create a subrouter for the /companies resource.
	companiesRouter := r.PathPrefix("/companies").Subrouter()
	companiesRouter.Use(middleware.AuthUserAdminRequiredMuxMiddleware())
	companiesRouter.HandleFunc("", Create).Methods("POST")
	companiesRouter.HandleFunc("/find", Find).Methods("GET")
	companiesRouter.HandleFunc("/{id:[0-9]+}", Get).Methods("GET")
	companiesRouter.HandleFunc("/{id:[0-9]+}", Update).Methods("PUT")
	companiesRouter.HandleFunc("/{id:[0-9]+}", Delete).Methods("DELETE")
}
