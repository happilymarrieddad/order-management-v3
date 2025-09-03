package locations

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
		middleware.WriteError(w, http.StatusBadRequest, "invalid location ID")
		return
	}

	authUser, found := middleware.GetAuthUserFromContext(r.Context())
	if !found {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Before deleting, we must ensure the user has ownership of the location.
	// Admins can delete any location.
	if !authUser.HasRole(types.RoleAdmin) {
		loc, found, err := repo.Locations().Get(r.Context(), authUser.CompanyID, id) // Corrected call
		if err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "unable to verify location ownership")
			return
		}
		if !found {
			middleware.WriteError(w, http.StatusNotFound, "location not found")
			return
		}
		if loc.CompanyID != authUser.CompanyID {
			middleware.WriteError(w, http.StatusForbidden, "user not authorized to delete this location")
			return
		}
	}

	if err := repo.Locations().Delete(r.Context(), id); err != nil {
		if types.IsNotFoundError(err) {
			middleware.WriteError(w, http.StatusNotFound, "location not found")
		} else {
			middleware.WriteError(w, http.StatusInternalServerError, "unable to delete location")
		}
		return // Always return after handling an error
	}

	w.WriteHeader(http.StatusNoContent)
}