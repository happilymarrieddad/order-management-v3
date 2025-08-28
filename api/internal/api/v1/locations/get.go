package locations

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// @Summary      Get a location by ID
// @Description  Retrieves the details of a single location, including its company and address.
// @Tags         locations
// @Produce      json
// @Param        id  path      int                      true  "Location ID"
// @Success      200 {object}  types.Location           "Successfully retrieved location"
// @Failure      400 {object}  middleware.ErrorResponse "Bad Request - Invalid ID"
// @Failure      404 {object}  middleware.ErrorResponse "Not Found - Location not found"
// @Failure      500 {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router       /locations/{id} [get]
func Get(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid location ID")
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

	// Populate related entities to provide a richer API response.
	location.Company, _, _ = repo.Companies().Get(r.Context(), location.CompanyID)
	location.Address, _, _ = repo.Addresses().Get(r.Context(), location.AddressID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(location)
}
