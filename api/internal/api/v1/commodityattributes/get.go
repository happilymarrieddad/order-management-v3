package commodityattributes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// Get handles retrieving a single commodity attribute by its ID.
func Get(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the context (cached by AuthMiddleware).
	_, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	repo := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid commodity attribute ID")
		return
	}

	attr, found, err := repo.CommodityAttributes().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get commodity attribute")
		return
	}

	if !found {
		middleware.WriteError(w, http.StatusNotFound, "commodity attribute not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(attr)
}
