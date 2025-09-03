package locations

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

func Update(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid location ID")
		return
	}

	var payload UpdateLocationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, middleware.FormatValidationErrors(err))
		return
	}

	authUser, found := middleware.GetAuthUserFromContext(r.Context())
	if !found {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Admins can update any location.
	// Non-admins can only update locations in their own company.
	var companyID int64
	if !authUser.HasRole(types.RoleAdmin) {
		companyID = authUser.CompanyID
	}

	loc, found, err := repo.Locations().Get(r.Context(), companyID, id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get location")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "location not found")
		return
	}

	// No manual validation check needed here, as it's handled by types.Validate(payload)

	if payload.AddressID != nil {
		// Validate new address if provided
		_, found, err := repo.Addresses().Get(r.Context(), *payload.AddressID)
		if err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "unable to validate address")
			return
		}
		if !found {
			middleware.WriteError(w, http.StatusBadRequest, "new address not found")
			return
		}
		loc.AddressID = *payload.AddressID
	}

	if payload.Name != nil {
		loc.Name = *payload.Name
	}

	if err := repo.Locations().Update(r.Context(), loc); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to update location")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loc)
}