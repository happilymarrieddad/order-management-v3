package addresses

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// @Summary      Get an address by ID
// @Description  Retrieves the details of a single address by its unique ID.
// @Tags         addresses
// @Produce      json
// @Param        id  path      int                      true  "Address ID"
// @Success      200 {object}  types.Address            "Successfully retrieved address"
// @Failure      400 {object}  middleware.ErrorResponse "Bad Request - Invalid ID"
// @Failure      404 {object}  middleware.ErrorResponse "Not Found - Address not found"
// @Security     BearerAuth
// @Router       /addresses/{id} [get]
// Get handles retrieving a single address by its ID.
func Get(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid address ID")
		return
	}

	address, found, err := repo.Addresses().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !found {
		middleware.WriteError(w, http.StatusNotFound, "address not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(address)
}
