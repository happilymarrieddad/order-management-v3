package users

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// Delete handles the deletion of a user.
//	@Summary	Delete a user
//	@Description	Deletes a user by their ID. A user can delete themselves, or an admin can delete any user within the same company.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	middleware.ErrorResponse	"Invalid User ID"
//	@Failure		403	{object}	middleware.ErrorResponse	"Forbidden"
//	@Failure		404	{object}	middleware.ErrorResponse	"User not found"
//	@Failure		500	{object}	middleware.ErrorResponse	"Internal Server Error"
//	@Router			/users/{id}	[delete]
func Delete(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())

	authUser, ok := middleware.GetAuthUserFromContext(r.Context())
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unable to get authenticated user")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	// Get the user to be deleted
	// When getting a user, the company id must be passed in
	userToDelete, has, err := gr.Users().Get(r.Context(), authUser.CompanyID, id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get user")
		return
	}
	if !has {
		middleware.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	// Authorization check:
	// 1. An admin can delete any user in their own company.
	// 2. A non-admin user can only delete themselves.
	isSameCompany := authUser.CompanyID == userToDelete.CompanyID
	isAdmin := authUser.HasRole(types.RoleAdmin)
	isSelf := authUser.ID == userToDelete.ID

	if !isSelf && (!isAdmin || !isSameCompany) {
		middleware.WriteError(w, http.StatusForbidden, "you are not authorized to delete this user")
		return
	}

	if err = gr.Users().Delete(r.Context(), id); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
