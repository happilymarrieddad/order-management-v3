package locations

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

// @Summary      Delete a location
// @Description  Deletes a location by its ID.
// @Tags         locations
// @Param        id  path      int  true  "Location ID"
// @Success      204 "No Content"
// @Failure      400 {object} middleware.ErrorResponse "Bad Request - Invalid ID"
// @Failure      500 {object} middleware.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router       /locations/{id} [delete]
func Delete(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid location ID")
		return
	}

	if err := repo.Locations().Delete(r.Context(), id); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
