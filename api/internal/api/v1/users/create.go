package users

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Create a new user
// @Description  Creates a new user account with the provided details.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body      CreateUserPayload        true  "User Creation Payload"
// @Success      201  {object}  types.User               "Successfully created user"
// @Failure      400  {object}  middleware.ErrorResponse "Bad Request - Invalid input or user already exists"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Router       /users [post]
// Create handles the HTTP request for creating a new user.
func Create(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}

	repo := middleware.GetRepo(r.Context())

	_, found, err := repo.Users().GetByEmail(r.Context(), payload.Email)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if found {
		middleware.WriteError(w, http.StatusBadRequest, "user with that email already exists")
		return
	}

	_, found, err = repo.Companies().Get(r.Context(), payload.CompanyID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusBadRequest, "company not found")
		return
	}

	_, found, err = repo.Addresses().Get(r.Context(), payload.AddressID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusBadRequest, "address not found")
		return
	}

	user := &types.User{
		Email:     payload.Email,
		Password:  payload.Password,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		CompanyID: payload.CompanyID,
		AddressID: payload.AddressID,
		Roles:     types.Roles{types.RoleUser},
	}

	if err := repo.Users().Create(r.Context(), user); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
