package locations

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Find locations
// @Description  Finds locations with optional filters and pagination by sending query parameters.
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        limit query int false "Number of records to return"
// @Param        offset query int false "Number of records to skip"
// @Param        name query string false "Location name filter"
// @Success      200  {object}  types.FindResult{data=[]types.Location} "A list of locations"
// @Failure      400  {object}  middleware.ErrorResponse "Bad Request"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /locations/find [get]
func Find(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	var opts repos.LocationFindOpts
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

	if name := query.Get("name"); name != "" {
		opts.Names = append(opts.Names, name)
	}

	// Get the authenticated user from the context (cached by AuthMiddleware).
	authUser, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// For now we force everyone to only see locations in their own company.
	opts.CompanyIDs = []int64{authUser.CompanyID}

	locations, count, err := repo.Locations().Find(r.Context(), &opts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to find locations")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.NewFindResult(locations, count))
}