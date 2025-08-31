package locations

import (
	"github.com/gorilla/mux"
)

// AddRoutes configures the location-related routes on the given subrouter.
func AddRoutes(r *mux.Router) {
	// Create a subrouter for the /locations resource.
	locationsRouter := r.PathPrefix("/locations").Subrouter()
	locationsRouter.HandleFunc("", Create).Methods("POST")
	locationsRouter.HandleFunc("/find", Find).Methods("GET")
	locationsRouter.HandleFunc("/{id:[0-9]+}", Get).Methods("GET")
	locationsRouter.HandleFunc("/{id:[0-9]+}", Update).Methods("PUT")
	locationsRouter.HandleFunc("/{id:[0-9]+}", Delete).Methods("DELETE")
}
