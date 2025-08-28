package addresses

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Update an address
// @Description  Updates an existing address with new details.
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        id      path      int                      true  "Address ID"
// @Param        address body      UpdateAddressPayload     true  "Address Update Payload"
// @Success      200     {object}  types.Address            "Successfully updated address"
// @Failure      400     {object}  middleware.ErrorResponse "Bad Request - Invalid input or ID"
// @Failure      404     {object}  middleware.ErrorResponse "Not Found - Address not found"
// @Security     BearerAuth
// @Router       /addresses/{id} [put]
// Update handles updating an existing address.
func Update(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid address ID")
		return
	}

	var payload UpdateAddressPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// First, get the existing address to ensure it exists
	address, found, err := repo.Addresses().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "address not found")
		return
	}

	// Update fields
	address.Line1 = payload.Line1
	if payload.Line2 != nil {
		address.Line2 = *payload.Line2
	} else {
		address.Line2 = "" // Clear the field if null is passed
	}
	address.City = payload.City
	address.State = payload.State
	address.PostalCode = payload.PostalCode

	if err := repo.Addresses().Update(r.Context(), address); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(address)
}
