package users

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// Create handles the creation of a new user.
func Create(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	var payload CreateUserPayload
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

	// Admins can create users in any company.
	// Non-admins can only create users in their own company.
	if !authUser.HasRole(types.RoleAdmin) {
		if authUser.CompanyID != payload.CompanyID {
			middleware.WriteError(w, http.StatusForbidden, "user not authorized to create users for this company")
			return
		}
	}

	// Check if user with that email already exists
	_, found, err := repo.Users().GetByEmail(r.Context(), payload.Email)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to check for existing user")
		return
	}
	if found {
		middleware.WriteError(w, http.StatusBadRequest, "user with that email already exists")
		return
	}

	// Check if dependencies exist
	_, found, err = repo.Companies().Get(r.Context(), payload.CompanyID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to validate company")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusBadRequest, "company not found")
		return
	}

	_, found, err = repo.Addresses().Get(r.Context(), payload.AddressID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to validate address")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusBadRequest, "address not found")
		return
	}

	user := &types.User{
		CompanyID: payload.CompanyID,
		Email:     payload.Email,
		Password:  payload.Password,
		AddressID: payload.AddressID,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Roles:     types.Roles{types.RoleUser}, // Default to user role
	}

	if err := repo.Users().Create(r.Context(), user); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to create user")
		return
	}

	user.Password = "" // Never return password

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
