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

	// Get the authenticated user to check for permissions
	authUserID, found := middleware.GetUserIDFromContext(r.Context())
	if !found {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	authUser, found, err := repo.Users().Get(r.Context(), authUserID)
	if err != nil || !found {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get authenticated user")
		return
	}

	loc, found, err := repo.Locations().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get location")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "location not found")
		return
	}

	// Admins can update any location.
	// Non-admins can only update locations in their own company.
	if !authUser.Roles.HasRole(types.RoleAdmin) {
		if authUser.CompanyID != loc.CompanyID {
			middleware.WriteError(w, http.StatusForbidden, "user not authorized to update this location")
			return
		}
	}

	// Validate new address if provided
	if payload.AddressID != nil {
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
