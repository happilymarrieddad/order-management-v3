package locations

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

// @Summary      Find locations
// @Description  Finds locations with optional filters and pagination by sending query parameters.
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        limit query int false "Number of records to return"
// @Param        offset query int false "Number of records to skip"
// @Param        name query string false "Location name filter"
// @Success      200  {object}  object{data=[]types.Location,total=int} "A list of locations"
// @Failure      400  {object}  middleware.ErrorResponse "Bad Request"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /locations/find [get]
func Find(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())

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

	opts := repos.LocationFindOpts{
		Limit:  limit,
		Offset: offset,
		Names:  r.URL.Query()["name"],
	}

	// Get the authenticated user from the context (cached by AuthMiddleware).
	authUser, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// For now we force everyone to only see locations in their own company.
	opts.CompanyIDs = []int64{authUser.CompanyID}

	locations, count, err := gr.Locations().Find(r.Context(), &opts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to find locations")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.NewFindResult(locations, count))
}