package locations

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Update a location
// @Description  Updates an existing location's details.
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        id       path      int                      true  "Location ID"
// @Param        location body      UpdateLocationPayload    true  "Location Update Payload"
// @Success      200      {object}  types.Location           "Successfully updated location"
// @Failure      400      {object}  middleware.ErrorResponse "Bad Request - Invalid input or ID"
// @Failure      401      {object}  middleware.ErrorResponse "Unauthorized"
// @Failure      403      {object}  middleware.ErrorResponse "Forbidden"
// @Failure      404      {object}  middleware.ErrorResponse "Not Found - Location not found"
// @Failure      500      {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /locations/{id} [put]
func Update(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())

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

	loc, found, err := gr.Locations().Get(r.Context(), companyID, id)
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
		_, found, err := gr.Addresses().Get(r.Context(), *payload.AddressID)
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

	if err := gr.Locations().Update(r.Context(), loc); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to update location")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loc)
}
