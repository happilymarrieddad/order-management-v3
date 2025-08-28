package addresses

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// Create handles the creation of a new address.
func Create(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Line1      string  `json:"line_1" validate:"required"`
		Line2      *string `json:"line_2"`
		City       string  `json:"city" validate:"required"`
		State      string  `json:"state" validate:"required"`
		PostalCode string  `json:"postal_code" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	repo := middleware.GetRepo(r.Context())

	address := &types.Address{
		Line1:      payload.Line1,
		City:       payload.City,
		State:      payload.State,
		PostalCode: payload.PostalCode,
	}

	if payload.Line2 != nil {
		address.Line2 = *payload.Line2
	}

	newAddress, err := repo.Addresses().Create(r.Context(), address)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAddress)
}
