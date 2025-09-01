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
// @Param        opts  body      repos.AddressFindOpts  true  "Find options"
// @Success      200   {object}  types.FindResult{data=[]types.Address} "A list of addresses"
// @Failure      400   {object}  middleware.ErrorResponse "Bad Request"
// @Failure      500   {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /addresses/find [post]
func Find(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	var opts repos.AddressFindOpts
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Apply default limit if not provided
	if opts.Limit == 0 {
		opts.Limit = 10
	}

	addrs, count, err := repo.Addresses().Find(r.Context(), &opts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to find addresses")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types.FindResult{
		Data:  addrs,
		Total: count,
	})
}
