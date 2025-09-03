package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// UpdateUserCompany handles the request to update a user's company.
// @Summary      Update a user's company
// @Description  Moves a user to a different company. This is an admin-only operation.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        payload body      UpdateUserCompanyPayload true "Company ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  middleware.ErrorResponse
// @Failure      401  {object}  middleware.ErrorResponse
// @Failure      403  {object}  middleware.ErrorResponse
// @Failure      404  {object}  middleware.ErrorResponse
// @Failure      500  {object}  middleware.ErrorResponse
// @Router       /v1/users/{id}/company [put]
// @Security     ApiKeyAuth
func UpdateUserCompany(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	usersRepo := repo.Users()
	companiesRepo := repo.Companies()

	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var p UpdateUserCompanyPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := types.Validate(p); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate that the target company exists
	_, found, err := companiesRepo.Get(r.Context(), p.CompanyID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to verify company")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "company not found")
		return
	}

	// Validate that the user exists
	_, found, err = usersRepo.GetIncludeInvisible(r.Context(), userID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to verify user")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	if err := usersRepo.UpdateUserCompany(r.Context(), userID, p.CompanyID); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to update user company")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
