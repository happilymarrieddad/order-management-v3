package locations

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Get a location
// @Description  Gets a location by its ID.
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Location ID"
// @Success      200  {object}  types.Location
// @Failure      400  {object}  middleware.ErrorResponse "Invalid Location ID"
// @Failure      401  {object}  middleware.ErrorResponse "Unauthorized"
// @Failure      403  {object}  middleware.ErrorResponse "Forbidden"
// @Failure      404  {object}  middleware.ErrorResponse "Location not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /locations/{id} [get]
func Get(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid location ID")
		return
	}

	authUser, found := middleware.GetAuthUserFromContext(r.Context())
	if !found {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Admins can access any location, non-admins can only access locations for their company
	var companyID int64
	if !authUser.HasRole(types.RoleAdmin) {
		companyID = authUser.CompanyID
	}

	loc, found, err := gr.Locations().Get(r.Context(), companyID, id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get location")
		return
	}

	if !found {
		middleware.WriteError(w, http.StatusNotFound, "location not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loc)
}