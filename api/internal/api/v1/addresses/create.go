package addresses

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// Create handles the creation of a new address.
// @Summary      Create a new address
// @Description  Creates a new address.
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        address body      CreateAddressPayload    true  "Address Creation Payload"
// @Success      201      {object}  types.Address           "Successfully created address"
// @Failure      400      {object}  middleware.ErrorResponse "Bad Request - Invalid input or validation failed"
// @Failure      401      {object}  middleware.ErrorResponse "Unauthorized"
// @Failure      500      {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /addresses [post]
func Create(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the context (cached by AuthMiddleware).
	_, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var payload CreateAddressPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, middleware.FormatValidationErrors(err))
		return
	}

	gr := middleware.GetRepo(r.Context())

	address := &types.Address{
		Line1:      payload.Line1,
		Line2:      payload.Line2,
		City:       payload.City,
		State:      payload.State,
		PostalCode: payload.PostalCode,
		Country:    payload.Country,
	}

	newAddress, err := gr.Addresses().Create(r.Context(), address)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to create address")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAddress)
}