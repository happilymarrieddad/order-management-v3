package addresses

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Find addresses
// @Description  Finds addresses with optional filters and pagination by sending a JSON body.
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        opts body      repos.AddressFindOpts true "Find options"
// @Success      200  {object}  types.FindResult{data=[]types.Address} "A list of addresses"
// @Failure      400  {object}  middleware.ErrorResponse "Bad Request"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /addresses/find [post]
func Find(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	var opts repos.AddressFindOpts
	// Decode the request body into the options struct
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Set default limit if none is provided or if it's invalid.
	if opts.Limit <= 0 {
		opts.Limit = 10
	}

	addresses, count, err := repo.Addresses().Find(r.Context(), &opts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := types.NewFindResult(addresses, count)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}