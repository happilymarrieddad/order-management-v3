package companies

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

// @Summary      Update a company
// @Description  Updates an existing company's details.
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        id      path      int                      true  "Company ID"
// @Param        company body      UpdateCompanyPayload     true  "Company Update Payload"
// @Success      200     {object}  types.Company            "Successfully updated company"
// @Failure      400     {object}  middleware.ErrorResponse "Bad Request - Invalid input or ID"
// @Failure      404     {object}  middleware.ErrorResponse "Not Found - Company not found"
// @Security     AppTokenAuth
// @Router       /companies/{id} [put]
func Update(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
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
		middleware.WriteError(w, http.StatusBadRequest, middleware.FormatValidationErrors(err))
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

	// Only admins can update companies.
	// This check is done in the middleware layer.

	if payload.Name != nil {
		company.Name = utils.Deref(payload.Name)
	}

	if payload.AddressID != nil {
		// Validate the new address exists before updating the company
		_, found, err := gr.Addresses().Get(r.Context(), utils.Deref(payload.AddressID))
		if err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "unable to validate address")
			return
		}
		if !found {
			middleware.WriteError(w, http.StatusBadRequest, "address not found")
			return
		}
		company.AddressID = utils.Deref(payload.AddressID)
	}

	if err := gr.Companies().Update(r.Context(), company); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to update company")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(company)
}