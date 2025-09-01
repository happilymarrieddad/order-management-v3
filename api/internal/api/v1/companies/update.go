package companies

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

func Update(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

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

	company, found, err := repo.Companies().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get company")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "company not found")
		return
	}

	if payload.Name != nil {
		company.Name = *payload.Name
	}
	if payload.AddressID != nil {
		company.AddressID = *payload.AddressID
	}

	if err := repo.Companies().Update(r.Context(), company); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to update company")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(company)
}
