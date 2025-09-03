package addresses

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

// @Summary      Find addresses
// @Description  Finds addresses with optional filters and pagination using query parameters.
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        id     query    []int  false  "Address IDs"
// @Param        limit  query    int    false  "Limit"
// @Param        offset query    int    false  "Offset"
// @Success      200    {object}  types.FindResult{data=[]types.Address} "A list of addresses"
// @Failure      400    {object}  middleware.ErrorResponse "Bad Request"
// @Failure      500    {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /addresses/find [get]
func Find(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the context (cached by AuthMiddleware).
	_, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	repo := middleware.GetRepo(r.Context())

	ids, err := utils.GetQueryInt64Slice(r, "id")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	limit, err := utils.GetQueryInt(r, "limit")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid limit format")
		return
	}
	if limit == 0 {
		limit = 10
	}

	offset, err := utils.GetQueryInt(r, "offset")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid offset format")
		return
	}

	opts := &repos.AddressFindOpts{
		IDs:    ids,
		Limit:  limit,
		Offset: offset,
	}

	addrs, count, err := repo.Addresses().Find(r.Context(), opts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to find addresses")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types.FindResult[*types.Address]{
		Data:  addrs,
		Total: count,
	})
}
