package products

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Get a product by ID
// @Description  Retrieves the details of a single product.
// @Tags         products
// @Produce      json
// @Param        id  path      int                      true  "Product ID"
// @Success      200 {object}  types.Product           "Successfully retrieved product"
// @Failure      400 {object}  middleware.ErrorResponse "Bad Request - Invalid ID"
// @Failure      401 {object}  middleware.ErrorResponse "Unauthorized"
// @Failure      404 {object}  middleware.ErrorResponse "Not Found - Product not found"
// @Failure      500 {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /products/{id} [get]
func Get(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the context (cached by AuthMiddleware).
	authUser, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	gr := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	product, found, err := gr.Products().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get product")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "product not found")
		return
	}

	// Authorization check: Normal users can only get products from their own company.
	// Admins can get any product.
	if !authUser.HasRole(types.RoleAdmin) && product.CompanyID != authUser.CompanyID {
		middleware.WriteError(w, http.StatusForbidden, "user not authorized to view this product")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}