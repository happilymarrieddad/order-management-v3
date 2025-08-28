package companies

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Create a new company
// @Description  Creates a new company with the provided details.
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        company body      CreateCompanyPayload     true  "Company Creation Payload"
// @Success      201     {object}  types.Company            "Successfully created company"
// @Failure      400     {object}  middleware.ErrorResponse "Bad Request - Invalid input"
// @Failure      500     {object}  middleware.ErrorResponse "Internal Server Error"
// @Router       /companies [post]
// Create handles the creation of a new company.
func Create(w http.ResponseWriter, r *http.Request) {
	var payload CreateCompanyPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	repo := middleware.GetRepo(r.Context())

	company := &types.Company{
		Name:      payload.Name,
		AddressID: payload.AddressID,
	}

	if err := repo.Companies().Create(r.Context(), company); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(company)
}
