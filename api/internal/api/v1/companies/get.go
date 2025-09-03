package companies

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/happilymarrieddad/order-management-v3/api/types"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

func Get(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid company ID")
		return
	}

	// Get the authenticated user to check for permissions
	authUserID, found := middleware.GetAuthUserFromContext(r.Context())
	if !found {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if !authUserID.HasRole(types.RoleAdmin) && authUserID.CompanyID != id {
		middleware.WriteError(w, http.StatusForbidden, "user not authorized to view this company")
		return
	}

	company, found, err := repo.Companies().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get company")
		return
	}

	if !found {
		middleware.WriteError(w, http.StatusNotFound, "company not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(company)
}
