package addresses

import (
	"github.com/gorilla/mux"
)

// AddRoutes configures the address-related routes on the given subrouter.
func AddRoutes(r *mux.Router) {
	// Create a subrouter for the /addresses resource.
	addressesRouter := r.PathPrefix("/addresses").Subrouter()
	addressesRouter.HandleFunc("", Create).Methods("POST")
	addressesRouter.HandleFunc("/find", Find).Methods("GET")
	addressesRouter.HandleFunc("/{id:[0-9]+}", Get).Methods("GET")
	addressesRouter.HandleFunc("/{id:[0-9]+}", Update).Methods("PUT")
}
