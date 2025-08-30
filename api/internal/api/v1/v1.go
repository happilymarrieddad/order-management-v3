package v1

import (
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/addresses"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/companies"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/locations"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/users"
)

// AddAuthRoutes configures the v1 API routes on the given subrouter.
func AddAuthRoutes(r *mux.Router) {
	// Create a subrouter for the /users resource.
	// This groups all user-related endpoints under /api/v1/users.
	usersRouter := r.PathPrefix("/users").Subrouter()

	// Define the routes for the users resource in a standard RESTful fashion.
	usersRouter.HandleFunc("", users.Create).Methods("POST")
	usersRouter.HandleFunc("/find", users.Find).Methods("GET") // A dedicated find endpoint for complex queries.
	usersRouter.HandleFunc("/{id:[0-9]+}", users.Get).Methods("GET")
	usersRouter.HandleFunc("/{id:[0-9]+}", users.Update).Methods("PUT")
	usersRouter.HandleFunc("/{id:[0-9]+}", users.Delete).Methods("DELETE")

	// Create a subrouter for the /companies resource.
	companiesRouter := r.PathPrefix("/companies").Subrouter()
	companiesRouter.HandleFunc("", companies.Create).Methods("POST")
	companiesRouter.HandleFunc("/find", companies.Find).Methods("GET")
	companiesRouter.HandleFunc("/{id:[0-9]+}", companies.Get).Methods("GET")
	companiesRouter.HandleFunc("/{id:[0-9]+}", companies.Update).Methods("PUT")
	companiesRouter.HandleFunc("/{id:[0-9]+}", companies.Delete).Methods("DELETE")

	// Create a subrouter for the /addresses resource.
	addressesRouter := r.PathPrefix("/addresses").Subrouter()
	addressesRouter.HandleFunc("", addresses.Create).Methods("POST")
	addressesRouter.HandleFunc("/find", addresses.Find).Methods("GET")
	addressesRouter.HandleFunc("/{id:[0-9]+}", addresses.Get).Methods("GET")
	addressesRouter.HandleFunc("/{id:[0-9]+}", addresses.Update).Methods("PUT")

	// Create a subrouter for the /locations resource.
	locationsRouter := r.PathPrefix("/locations").Subrouter()
	locationsRouter.HandleFunc("", locations.Create).Methods("POST")
	locationsRouter.HandleFunc("/find", locations.Find).Methods("GET")
	locationsRouter.HandleFunc("/{id:[0-9]+}", locations.Get).Methods("GET")
	locationsRouter.HandleFunc("/{id:[0-9]+}", locations.Update).Methods("PUT")
	locationsRouter.HandleFunc("/{id:[0-9]+}", locations.Delete).Methods("DELETE")
}
