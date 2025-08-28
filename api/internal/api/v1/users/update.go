package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Update a user
// @Description  Updates an existing user's details. This does not update the password.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int                      true  "User ID"
// @Param        user body      UpdateUserPayload        true  "User Update Payload"
// @Success      200  {object}  types.User               "Successfully updated user"
// @Failure      400  {object}  middleware.ErrorResponse "Bad Request - Invalid input or ID"
// @Failure      404  {object}  middleware.ErrorResponse "Not Found - User not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/{id} [put]
// Update handles the HTTP request to update an existing user.
func Update(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
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
		middleware.WriteError(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}

	user, found, err := repo.Users().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "user not found")
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

	user.FirstName = payload.FirstName
	user.LastName = payload.LastName
	user.AddressID = payload.AddressID

	if err := repo.Users().Update(r.Context(), user); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
