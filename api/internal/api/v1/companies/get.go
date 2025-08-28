package companies

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// @Summary      Get a company by ID
// @Description  Retrieves the details of a single company by its unique ID.
// @Tags         companies
// @Produce      json
// @Param        id  path      int                      true  "Company ID"
// @Success      200 {object}  types.Company            "Successfully retrieved company"
// @Failure      400 {object}  middleware.ErrorResponse "Bad Request - Invalid ID"
// @Failure      404 {object}  middleware.ErrorResponse "Not Found - Company not found"
// @Security     BearerAuth
// @Router       /companies/{id} [get]
// Get handles retrieving a single company by its ID.
func Get(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid company ID")
		return
	}

	company, found, err := repo.Companies().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
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
