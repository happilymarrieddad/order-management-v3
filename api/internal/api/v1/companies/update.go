package companies

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Update a company
// @Description  Updates an existing company with new details.
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        id      path      int                      true  "Company ID"
// @Param        company body      UpdateCompanyPayload     true  "Company Update Payload"
// @Success      200     {object}  types.Company            "Successfully updated company"
// @Failure      400     {object}  middleware.ErrorResponse "Bad Request - Invalid input or ID"
// @Failure      404     {object}  middleware.ErrorResponse "Not Found - Company not found"
// @Security     BearerAuth
// @Router       /companies/{id} [put]
// Update handles updating an existing company.
func Update(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid company ID")
		return
	}

	var payload UpdateCompanyPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// First, get the existing company to ensure it exists
	company, found, err := repo.Companies().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "company not found")
		return
	}

	// Update fields
	company.Name = payload.Name
	company.AddressID = payload.AddressID

	if err := repo.Companies().Update(r.Context(), company); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(company)
}
