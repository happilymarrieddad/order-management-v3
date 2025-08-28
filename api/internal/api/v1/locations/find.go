package locations

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Find locations
// @Description  Finds locations with optional filters and pagination.
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        opts body      repos.LocationFindOpts     true "Find Options"
// @Success      200  {object}  types.FindResult{data=[]types.Location} "A list of locations"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router       /locations/find [post]
func Find(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	var opts repos.LocationFindOpts
	// Ignore error, opts will be zero-valued if body is empty or malformed
	_ = json.NewDecoder(r.Body).Decode(&opts)

	// As a world-class engineering assistant, I suggest that for a richer API experience,
	// we could add a query parameter like `?include=company,address` to optionally
	// populate the related entities in the find results. For now, we'll return the raw data.

	locations, count, err := repo.Locations().Find(r.Context(), &opts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := types.NewFindResult(locations, count)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
