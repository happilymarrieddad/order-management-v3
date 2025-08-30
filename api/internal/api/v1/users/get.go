package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// @Summary      Get a user by ID
// @Description  Retrieves the details of a single user by their unique ID.
// @Tags         users
// @Produce      json
// @Param        id  path      int                      true  "User ID"
// @Success      200 {object}  types.User               "Successfully retrieved user"
// @Failure      400 {object}  middleware.ErrorResponse "Bad Request - Invalid ID"
// @Failure      404 {object}  middleware.ErrorResponse "Not Found - User not found"
// @Failure      500 {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /users/{id} [get]
func Get(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid user ID")
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

	user.Password = "" // Never return the password

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
