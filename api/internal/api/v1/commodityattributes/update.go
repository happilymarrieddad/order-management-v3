package commodityattributes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

// Update handles updating an existing commodity attribute.
// @Summary      Update a commodity attribute
// @Description  Updates an existing commodity attribute with the provided details.
// @Tags         commodity-attributes
// @Accept       json
// @Produce      json
// @Param        id        path      int  true  "Commodity Attribute ID"
// @Param        attribute body      UpdateCommodityAttributePayload true  "Commodity Attribute Update Payload"
// @Success      200       {object}  types.CommodityAttribute        "Successfully updated commodity attribute"
// @Failure      400       {object}  middleware.ErrorResponse        "Bad Request - Invalid input"
// @Failure      401       {object}  middleware.ErrorResponse        "Unauthorized - Missing or invalid token"
// @Failure      403       {object}  middleware.ErrorResponse        "Forbidden - Insufficient permissions"
// @Failure      404       {object}  middleware.ErrorResponse        "Not Found - Commodity attribute not found"
// @Failure      409       {object}  middleware.ErrorResponse        "Conflict - Duplicate attribute name"
// @Failure      500       {object}  middleware.ErrorResponse        "Internal Server Error"
// @Router       /commodity-attributes/{id} [put]
func Update(w http.ResponseWriter, r *http.Request) {
	gr := middleware.GetRepo(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid commodity attribute ID")
		return
	}

	var payload UpdateCommodityAttributePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, middleware.FormatValidationErrors(err))
		return
	}

	// Get existing attribute to ensure it exists and to preserve original fields
	ca, found, err := gr.CommodityAttributes().Get(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "unable to get commodity attribute")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusNotFound, "commodity attribute not found")
		return
	}

	if payload.Name != nil {
				ca.Name = utils.Deref(payload.Name)
	}
	if payload.CommodityType != nil {
				ca.CommodityType = utils.Deref(payload.CommodityType)
	}

	if err := gr.CommodityAttributes().Update(r.Context(), ca); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			middleware.WriteError(w, http.StatusConflict, "Commodity attribute with this name already exists")
			return
		}
		middleware.WriteError(w, http.StatusInternalServerError, "unable to update commodity attribute")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ca)
}
