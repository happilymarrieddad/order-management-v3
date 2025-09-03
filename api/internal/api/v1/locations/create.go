package locations

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

func Create(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	var payload CreateLocationPayload
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

	if !authUser.HasRole(types.RoleAdmin) && authUser.CompanyID != payload.CompanyID {
		middleware.WriteError(w, http.StatusForbidden, "user not authorized to create locations for this company")
		return
	}

	// Validate dependencies
	_, found, err := repo.Companies().Get(r.Context(), payload.CompanyID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to validate company")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusBadRequest, "company not found")
		return
	}

	_, found, err = repo.Addresses().Get(r.Context(), payload.AddressID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to validate address")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusBadRequest, "address not found")
		return
	}

	loc := &types.Location{
		Name:      payload.Name,
		CompanyID: payload.CompanyID,
		AddressID: payload.AddressID,
	}

	if err := repo.Locations().Create(r.Context(), loc); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to create location")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loc)
}
