package companies

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// @Summary      Delete a company
// @Description  Deletes a company by its ID.
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Company ID"
// @Success      204  "No Content"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid Company ID"
// @Failure      404  {object}  middleware.ErrorResponse "Company not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal Server Error"
// @Security     AppTokenAuth
// @Router       /companies/{id} [delete]
func Delete(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid company ID")
		return
	}

	if err := gr.Companies().Delete(r.Context(), id); err != nil {
		if types.IsNotFoundError(err) {
			middleware.WriteError(w, http.StatusNotFound, "company not found")
		} else {
			middleware.WriteError(w, http.StatusInternalServerError, "unable to delete company")
		}
		return // Always return after handling an error
	}

	w.WriteHeader(http.StatusNoContent)
}