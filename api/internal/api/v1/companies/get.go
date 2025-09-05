package companies

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/happilymarrieddad/order-management-v3/api/types"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// @Summary      Get a company
// @Description  Gets a company by its ID.
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Company ID"
// @Success      200  {object}  types.Company
// @Failure      400  {object}  middleware.ErrorResponse "Invalid Company ID"
// @Failure      401  {object}  middleware.ErrorResponse "Unauthorized"
// @Failure      403  {object}  middleware.ErrorResponse "Forbidden"
// @Failure      404  {object}  middleware.ErrorResponse "Company not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /companies/{id} [get]
func Get(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid company ID")
		return
	}

	// Get the authenticated user to check for permissions
	authUser, found := middleware.GetAuthUserFromContext(r.Context())
	if !found {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if !authUser.HasRole(types.RoleAdmin) && authUser.CompanyID != id {
		middleware.WriteError(w, http.StatusForbidden, "user not authorized to view this company")
		return
	}

	company, found, err := gr.Companies().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get company")
		return
	}

	if !found {
		middleware.WriteError(w, http.StatusNotFound, "company not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(company)
}