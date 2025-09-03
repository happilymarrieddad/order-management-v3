package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
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

	var opts repos.UserFindOpts
	query := r.URL.Query()

	// Parse query parameters
	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			opts.Limit = int(limit)
		}
	}
	if opts.Limit <= 0 {
		opts.Limit = 10
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		if offset, err := strconv.ParseInt(offsetStr, 10, 64); err == nil {
			opts.Offset = int(offset)
		}
	}

	if companyIDStr := query.Get("company_id"); companyIDStr != "" {
		if companyID, err := strconv.ParseInt(companyIDStr, 10, 64); err == nil {
			opts.CompanyID = companyID
		}
	}

	// Handle multiple IDs
	for _, idStr := range query["id"] {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			opts.IDs = append(opts.IDs, id)
		}
	}

	// Handle multiple emails
	for _, email := range query["email"] {
		opts.Emails = append(opts.Emails, email)
	}

	// Handle multiple first names
	for _, firstName := range query["first_name"] {
		opts.FirstNames = append(opts.FirstNames, firstName)
	}

	// Handle multiple last names
	for _, lastName := range query["last_name"] {
		opts.LastNames = append(opts.LastNames, lastName)
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
