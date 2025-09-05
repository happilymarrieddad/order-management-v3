package products

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Create a new product
// @Description  Creates a new product with the provided details and attributes.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product body      CreateProductPayload    true  "Product Creation Payload"
// @Success      201     {object}  types.Product           "Successfully created product"
// @Failure      400     {object}  middleware.ErrorResponse "Bad Request - Invalid input or validation failed"
// @Failure      401     {object}  middleware.ErrorResponse "Unauthorized"
// @Failure      403     {object}  middleware.ErrorResponse "Forbidden"
// @Failure      500     {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /products [post]
func Create(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())

	var payload CreateProductPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, middleware.FormatValidationErrors(err))
		return
	}

	authUser, found := middleware.GetAuthUserFromContext(r.Context())
	if !found {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Only admins can create products for other companies.
	// Non-admins can only create products for their own company.
	if !authUser.HasRole(types.RoleAdmin) && authUser.CompanyID != payload.CompanyID {
		middleware.WriteError(w, http.StatusForbidden, "user not authorized to create products for this company")
		return
	}

	// Validate commodity exists
	_, found, err := gr.Commodities().Get(r.Context(), payload.CommodityID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to validate commodity")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusBadRequest, "commodity not found")
		return
	}

	product := &types.Product{
		CompanyID:   payload.CompanyID,
		CommodityID: payload.CommodityID,
	}

	if err := gr.Products().Create(r.Context(), product, payload.Attributes); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to create product")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}
