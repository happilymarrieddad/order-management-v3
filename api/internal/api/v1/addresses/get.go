package addresses

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// @Summary      Get an address
// @Description  Gets an address by its ID.
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Address ID"
// @Success      200  {object}  types.Address
// @Failure      400  {object}  middleware.ErrorResponse "Invalid Address ID"
// @Failure      401  {object}  middleware.ErrorResponse "Unauthorized"
// @Failure      404  {object}  middleware.ErrorResponse "Address not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /addresses/{id} [get]
func Get(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the context (cached by AuthMiddleware).
	_, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	gr := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid address ID")
		return
	}

	addr, found, err := gr.Addresses().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get address")
		return
	}

	if !found {
		middleware.WriteError(w, http.StatusNotFound, "address not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(addr)
}