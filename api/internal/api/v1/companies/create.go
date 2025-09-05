package companies

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

func Create(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())

	var payload CreateCompanyPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, middleware.FormatValidationErrors(err))
		return
	}

	// Validate that the address exists
	_, found, err := gr.Addresses().Get(r.Context(), payload.AddressID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to validate address")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusBadRequest, "address not found")
		return
	}

	company := &types.Company{
		Name:      payload.Name,
		AddressID: payload.AddressID,
	}

	if err := gr.Companies().Create(r.Context(), company); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to create company")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(company)
}
