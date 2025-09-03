package addresses

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
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
	// Get the authenticated user from the context (cached by AuthMiddleware).
	_, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

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
		middleware.WriteError(w, http.StatusBadRequest, middleware.FormatValidationErrors(err))
		return
	}

	// First, get the existing address to ensure it exists
	address, found, err := repo.Addresses().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get address")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "address not found")
		return
	}

	// Update fields
	if payload.Line1 != nil {
		address.Line1 = utils.Deref(payload.Line1)
	}
	if payload.Line2 != nil {
		address.Line2 = utils.Deref(payload.Line2)
	}
	if payload.City != nil {
		address.City = utils.Deref(payload.City)
	}
	if payload.State != nil {
		address.State = utils.Deref(payload.State)
	}
	if payload.PostalCode != nil {
		address.PostalCode = utils.Deref(payload.PostalCode)
	}
	if payload.Country != nil {
		address.Country = utils.Deref(payload.Country)
	}

	if err := repo.Addresses().Update(r.Context(), address); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to update address")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(address)
}
