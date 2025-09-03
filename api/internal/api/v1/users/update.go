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

	// Get the authenticated user from the context (cached by AuthMiddleware).
	authUser, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get the user to be updated
	targetUser, found, err := repo.Users().Get(r.Context(), authUser.CompanyID, id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get user")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	// Admins can update any user.
	// Non-admins can only update themselves.
	if !authUser.Roles.HasRole(types.RoleAdmin) {
		if authUser.ID != targetUser.ID {
			middleware.WriteError(w, http.StatusForbidden, "user not authorized to update this user")
			return
		}
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

	targetUser.FirstName = payload.FirstName
	targetUser.LastName = payload.LastName
	targetUser.AddressID = payload.AddressID

	if err := repo.Users().Update(r.Context(), targetUser); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to update user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(targetUser)
}
