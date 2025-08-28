package locations

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Create a new location
// @Description  Creates a new location for a company. The location name must be unique per company.
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        location body      CreateLocationPayload    true  "Location Creation Payload"
// @Success      201      {object}  types.Location           "Successfully created location"
// @Failure      400      {object}  middleware.ErrorResponse "Bad Request - Invalid input, duplicate name, or dependency not found"
// @Failure      500      {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router       /locations [post]
func Create(w http.ResponseWriter, r *http.Request) {
	var payload CreateLocationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}

	repo := middleware.GetRepo(r.Context())

	// Dependency checks
	_, found, err := repo.Companies().Get(r.Context(), payload.CompanyID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusBadRequest, "company not found")
		return
	}

	_, found, err = repo.Addresses().Get(r.Context(), payload.AddressID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusBadRequest, "address not found")
		return
	}

	location := &types.Location{
		CompanyID: payload.CompanyID,
		AddressID: payload.AddressID,
		Name:      payload.Name,
	}

	if err := repo.Locations().Create(r.Context(), location); err != nil {
		if errors.Is(err, repos.ErrLocationNameExists) {
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
		} else {
			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(location)
}
