package companies

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// @Summary      Delete a company
// @Description  Deletes a company by its ID.
// @Tags         companies
// @Param        id  path      int  true  "Company ID"
// @Success      204 "No Content"
// @Failure      400 {object} middleware.ErrorResponse "Bad Request - Invalid ID"
// @Failure      500 {object} middleware.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router       /companies/{id} [delete]
// Delete handles the deletion of a company.
func Delete(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid company ID")
		return
	}

	if err := repo.Companies().Delete(r.Context(), id); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
