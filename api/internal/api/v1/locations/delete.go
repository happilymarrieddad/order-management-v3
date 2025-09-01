package locations

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid location ID")
		return
	}

	// Ensure the location exists before attempting to delete
	_, found, err := repo.Locations().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get location")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "location not found")
		return
	}

	if err := repo.Locations().Delete(r.Context(), id); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to delete location")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
