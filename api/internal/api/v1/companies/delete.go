package companies

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
		middleware.WriteError(w, http.StatusBadRequest, "invalid company ID")
		return
	}

	// Ensure the company exists before attempting to delete
	_, found, err := repo.Companies().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get company")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "company not found")
		return
	}

	if err := repo.Companies().Delete(r.Context(), id); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to delete company")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
