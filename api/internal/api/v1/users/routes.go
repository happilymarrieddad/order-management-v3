package users

import (
	"github.com/gorilla/mux"
)

// AddRoutes configures the user-related routes on the given subrouter.
func AddRoutes(r *mux.Router) {
	// Create a subrouter for the /users resource.
	usersRouter := r.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("", Create).Methods("POST")
	usersRouter.HandleFunc("/find", Find).Methods("GET")
	usersRouter.HandleFunc("/{id:[0-9]+}", Get).Methods("GET")
	usersRouter.HandleFunc("/{id:[0-9]+}", Update).Methods("PUT")
	usersRouter.HandleFunc("/{id:[0-9]+}", Delete).Methods("DELETE")
}
