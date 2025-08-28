package users

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Find users
// @Description  Finds users with optional filters and pagination.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        opts body      repos.UserFindOpts         true "Find Options"
// @Success      200  {object}  types.FindResult{data=[]types.User} "A list of users"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/find [post]
// Find handles the HTTP request to search for users based on criteria.
func Find(w http.ResponseWriter, r *http.Request) {
	var opts repos.UserFindOpts
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	repo := middleware.GetRepo(r.Context())

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
	json.NewEncoder(w).Encode(response)
}
