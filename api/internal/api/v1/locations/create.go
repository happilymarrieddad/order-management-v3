package locations

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Create a new location
// @Description  Creates a new location for a company.
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        location body      CreateLocationPayload    true  "Location Creation Payload"
// @Success      201      {object}  types.Location           "Successfully created location"
// @Failure      400      {object}  middleware.ErrorResponse "Bad Request - Invalid input or validation failed"
// @Failure      401      {object}  middleware.ErrorResponse "Unauthorized"
// @Failure      403      {object}  middleware.ErrorResponse "Forbidden"
// @Failure      500      {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /locations [post]
func Create(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())

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
	_, found, err := gr.Companies().Get(r.Context(), payload.CompanyID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to validate company")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusBadRequest, "company not found")
		return
	}

	_, found, err = gr.Addresses().Get(r.Context(), payload.AddressID)
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

	if err := gr.Locations().Create(r.Context(), loc); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to create location")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loc)
}