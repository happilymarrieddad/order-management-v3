package commodityattributes

import (
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// AddRoutes configures the commodity attribute-related routes on the given subrouter.
func AddRoutes(r *mux.Router) {
	// Create a subrouter for the /commodity-attributes resource.
	commodityAttributesRouter := r.PathPrefix("/commodity-attributes").Subrouter()
	commodityAttributesRouter.Use(middleware.AuthUserAdminRequiredMuxMiddleware())
	commodityAttributesRouter.HandleFunc("", Create).Methods("POST")
	commodityAttributesRouter.HandleFunc("/find", Find).Methods("POST")
	commodityAttributesRouter.HandleFunc("/{id:[0-9]+}", Get).Methods("GET")
	commodityAttributesRouter.HandleFunc("/{id:[0-9]+}", Update).Methods("PUT")
}
