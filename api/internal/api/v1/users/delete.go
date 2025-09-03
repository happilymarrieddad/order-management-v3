package users

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Delete a user
// @Description  Deletes a user by their ID.
// @Tags         users
// @Param        id  path      int  true  "User ID"
// @Success      204 "No Content"
// @Failure      400 {object} middleware.ErrorResponse "Bad Request - Invalid ID"
// @Failure      500 {object} middleware.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/{id} [delete]
// Delete handles the HTTP request to delete a user by their ID.
func Delete(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Get the authenticated user from the context (cached by AuthMiddleware).
	authUser, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Only admins can delete users.
	if !authUser.HasRole(types.RoleAdmin) {
		middleware.WriteError(w, http.StatusForbidden, "unauthorized")
		return
	}

	if err := repo.Users().Delete(r.Context(), id); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
