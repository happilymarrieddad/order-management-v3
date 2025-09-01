package addresses

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
)

func Get(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid address ID")
		return
	}

	addr, found, err := repo.Addresses().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get address")
		return
	}

	if !found {
		middleware.WriteError(w, http.StatusNotFound, "address not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(addr)
}
