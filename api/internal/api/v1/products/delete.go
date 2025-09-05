package products

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Delete a product
// @Description  Deletes a product by its ID.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      204  "No Content"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid Product ID"
// @Failure      404  {object}  middleware.ErrorResponse "Product not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /products/{id} [delete]
func Delete(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the context (cached by AuthMiddleware).
	var found bool
	var err error
	_, found = middleware.GetAuthUserFromContext(r.Context())
	if !found { // Should be caught by middleware, but good practice to check
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	gr := middleware.GetRepo(r.Context())

	var id int64 // Declare id here
	id, err = strconv.ParseInt(mux.Vars(r)["id"], 10, 64) // Assign value to id
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	// Check if product exists
	_, found, err = gr.Products().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get product")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "product not found")
		return
	}

	if err := gr.Products().Delete(r.Context(), id); err != nil {
		if types.IsNotFoundError(err) {
			middleware.WriteError(w, http.StatusNotFound, "product not found")
		} else {
			middleware.WriteError(w, http.StatusInternalServerError, "unable to delete product")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}