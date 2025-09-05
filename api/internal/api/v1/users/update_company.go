package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// UpdateUserCompany handles updating a user's company.
//
//	@Summary	Update a user's company (Admin only)
//	@Description	Moves a user to a new company. This is an admin-only endpoint.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Param			body	body	UpdateUserCompanyPayload	true	"Update Company Payload"
//	@Success		204	"No Content"
//	@Failure		400	{object}	middleware.ErrorResponse	"Invalid request body or Company not found"
//	@Failure		403	{object}	middleware.ErrorResponse	"Forbidden"
//	@Failure		404	{object}	middleware.ErrorResponse	"User not found"
//	@Failure		500	{object}	middleware.ErrorResponse	"Internal Server Error"
//	@Router			/users/{id}/company	[put]
func UpdateUserCompany(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())

	_, ok := middleware.GetAuthUserFromContext(r.Context())
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

	var p UpdateUserCompanyPayload
	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err = validator.New().Struct(p); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get the user to be updated
	// The route is admin-only, so we can fetch the user without a companyID scope.
	_, has, err := gr.Users().Get(r.Context(), 0, id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get user")
		return
	}
	if !has {
		middleware.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	// Check if the new company exists
	_, has, err = gr.Companies().Get(r.Context(), p.CompanyID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get company")
		return
	}
	if !has {
		middleware.WriteError(w, http.StatusBadRequest, "company not found") // 400 because the payload is invalid
		return
	}

	if err = gr.Users().UpdateUserCompany(r.Context(), id, p.CompanyID); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to update user company")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
