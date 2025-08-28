package locations

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Update a location
// @Description  Updates an existing location's name and/or address.
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        id       path      int                      true  "Location ID"
// @Param        location body      UpdateLocationPayload    true  "Location Update Payload"
// @Success      200      {object}  types.Location           "Successfully updated location"
// @Failure      400      {object}  middleware.ErrorResponse "Bad Request - Invalid input or ID"
// @Failure      404      {object}  middleware.ErrorResponse "Not Found - Location or new address not found"
// @Failure      500      {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router       /locations/{id} [put]
func Update(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
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
		middleware.WriteError(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}

	location, found, err := repo.Locations().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "location not found")
		return
	}

	_, found, err = repo.Addresses().Get(r.Context(), payload.AddressID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusBadRequest, "new address not found")
		return
	}

	location.Name = payload.Name
	location.AddressID = payload.AddressID

	if err := repo.Locations().Update(r.Context(), location); err != nil {
		if errors.Is(err, repos.ErrLocationNameExists) {
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
		} else {
			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(location)
}
