package addresses

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// @Summary      Delete an address
// @Description  Deletes an address by its ID.
// @Tags         addresses
// @Param        id  path      int  true  "Address ID"
// @Success      204 "No Content"
// @Failure      400 {object} middleware.ErrorResponse "Bad Request - Invalid ID"
// @Failure      500 {object} middleware.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router       /addresses/{id} [delete]
// Delete handles the deletion of an address.
func Delete(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid address ID")
		return
	}

	if err := repo.Addresses().Delete(r.Context(), id); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
