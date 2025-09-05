package products

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

// @Summary      Find products
// @Description  Finds products with optional filters and pagination by sending query parameters.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        limit query int false "Number of records to return"
// @Param        offset query int false "Number of records to skip"
// @Param        name query string false "Product name filter"
// @Success      200  {object}  object{data=[]types.Product,total=int} "A list of products"
// @Failure      400  {object}  middleware.ErrorResponse "Bad Request"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /products/find [get]
func Find(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())

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

	opts := repos.ProductFindOpts{
		Limit:  limit,
		Offset: offset,
		Name:  r.URL.Query().Get("name"),
	}

	// Get the authenticated user from the context (cached by AuthMiddleware).
	authUser, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// For now we force everyone to only see products in their own company.
	opts.CompanyID = authUser.CompanyID

	products, count, err := gr.Products().Find(r.Context(), &opts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to find products")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.NewFindResult(products, count))
}
