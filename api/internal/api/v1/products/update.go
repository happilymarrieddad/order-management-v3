package products

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Update a product
// @Description  Updates an existing product by ID.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id        path      int                      true  "Product ID"
// @Param        product body      UpdateProductPayload   true  "Product Update Payload"
// @Success      200       {object}  types.Product          "Successfully updated product"
// @Failure      400       {object}  middleware.ErrorResponse "Bad Request - Invalid input or validation failed"
// @Failure      404       {object}  middleware.ErrorResponse "Not Found - Product not found"
// @Failure      500       {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /products/{id} [put]
func Update(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the context (cached by AuthMiddleware).
	_, found := middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	gr := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	var payload UpdateProductPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, middleware.FormatValidationErrors(err))
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

	if payload.CommodityID != nil {
		// Validate commodity exists
		_, found, err := gr.Commodities().Get(r.Context(), *payload.CommodityID)
		if err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "unable to validate commodity")
			return
		}
		if !found {
			middleware.WriteError(w, http.StatusBadRequest, "commodity not found")
			return
		}
		product.CommodityID = *payload.CommodityID
	}

	if err := gr.Products().Update(r.Context(), product, payload.Attributes); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to update product")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
