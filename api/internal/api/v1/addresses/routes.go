package addresses

import (
	"github.com/gorilla/mux"
)

// AddRoutes configures the address-related routes on the given subrouter.
func AddRoutes(r *mux.Router) {
	s := r.PathPrefix("/addresses").Subrouter()

	// Routes for any authenticated user
	s.HandleFunc("/{id:[0-9]+}", Get).Methods("GET")
	s.HandleFunc("/find", Find).Methods("GET")
	s.HandleFunc("", Create).Methods("POST")
	s.HandleFunc("/{id:[0-9]+}", Update).Methods("PUT")
}
