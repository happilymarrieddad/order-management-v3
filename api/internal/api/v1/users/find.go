package users

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

// @Summary      Find users
// @Description  Finds users with optional filters and pagination by sending query parameters.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        limit query int false "Number of records to return"
// @Param        offset query int false "Number of records to skip"
// @Param        company_id query int false "Filter by Company ID"
// @Param        id query []int false "Filter by User IDs"
// @Param        email query []string false "Filter by Emails"
// @Param        first_name query []string false "Filter by First Names"
// @Param        last_name query []string false "Filter by Last Names"
// @Success      200  {object}  types.FindResult{data=[]types.User} "A list of users"
// @Failure      400  {object}  middleware.ErrorResponse "Bad Request"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /users/find [get]
func Find(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	limit, err := utils.GetQueryInt(r, "limit")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid limit format")
		return
	}
	if limit == 0 {
		limit = 10
	}

	offset, err := utils.GetQueryInt(r, "offset")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid offset format")
		return
	}

	companyID, err := utils.GetQueryInt64(r, "company_id")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid company_id format")
		return
	}

	ids, err := utils.GetQueryInt64Slice(r, "id")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	opts := repos.UserFindOpts{
		Limit:      limit,
		Offset:     offset,
		CompanyID:  companyID,
		IDs:        ids,
		Emails:     r.URL.Query()["email"],
		FirstNames: r.URL.Query()["first_name"],
		LastNames:  r.URL.Query()["last_name"],
	}

	// Get the authenticated user from the context (cached by AuthMiddleware).
	authUser, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Admins can see all users. Non-admins can only see users in their own company.
	if !authUser.HasRole(types.RoleAdmin) {
		opts.CompanyID = authUser.CompanyID
	}

	users, count, err := repo.Users().Find(r.Context(), &opts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, u := range users {
		u.Password = "" // Never return the password
	}

	response := types.NewFindResult(users, count)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
