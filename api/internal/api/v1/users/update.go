package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// Update handles updating an existing user.
func Update(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var payload UpdateUserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, middleware.FormatValidationErrors(err))
		return
	}

	user, found, err := repo.Users().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get user")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	// Validate the new address exists before updating the user
	_, found, err = repo.Addresses().Get(r.Context(), payload.AddressID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to validate address")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusBadRequest, "address not found")
		return
	}

	user.FirstName = payload.FirstName
	user.LastName = payload.LastName
	user.AddressID = payload.AddressID

	if err := repo.Users().Update(r.Context(), user); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to update user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
