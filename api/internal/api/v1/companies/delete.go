package companies

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid company ID")
		return
	}

	if err := repo.Companies().Delete(r.Context(), id); err != nil {
		if types.IsNotFoundError(err) {
			middleware.WriteError(w, http.StatusNotFound, "company not found")
		} else {
			middleware.WriteError(w, http.StatusInternalServerError, "unable to delete company")
		}
		return // Always return after handling an error
	}

	w.WriteHeader(http.StatusNoContent)
}
