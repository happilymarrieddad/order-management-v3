package addresses

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// Create handles the creation of a new address.
func Create(w http.ResponseWriter, r *http.Request) {
	var payload CreateAddressPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, middleware.FormatValidationErrors(err))
		return
	}

	repo := middleware.GetRepo(r.Context())

	address := &types.Address{
		Line1:      payload.Line1,
		Line2:      payload.Line2,
		City:       payload.City,
		State:      payload.State,
		PostalCode: payload.PostalCode,
		Country:    payload.Country,
	}

	newAddress, err := repo.Addresses().Create(r.Context(), address)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to create address")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAddress)
}
