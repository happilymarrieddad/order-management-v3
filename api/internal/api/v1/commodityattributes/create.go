package commodityattributes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// Create handles the creation of a new commodity attribute.
// @Summary      Create a new commodity attribute
// @Description  Creates a new commodity attribute with the provided details.
// @Tags         commodity-attributes
// @Accept       json
// @Produce      json
// @Param        attribute body      CreateCommodityAttributePayload true  "Commodity Attribute Creation Payload"
// @Success      201       {object}  types.CommodityAttribute        "Successfully created commodity attribute"
// @Failure      400       {object}  middleware.ErrorResponse        "Bad Request - Invalid input"
// @Failure      401       {object}  middleware.ErrorResponse        "Unauthorized - Missing or invalid token"
// @Failure      403       {object}  middleware.ErrorResponse        "Forbidden - Insufficient permissions"
// @Failure      409       {object}  middleware.ErrorResponse        "Conflict - Duplicate attribute name"
// @Failure      500       {object}  middleware.ErrorResponse        "Internal Server Error"
// @Router       /commodity-attributes [post]
func Create(w http.ResponseWriter, r *http.Request) {
	repo := middleware.GetRepo(r.Context())
	var payload CreateCommodityAttributePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	ca := &types.CommodityAttribute{
		Name:          payload.Name,
		CommodityType: payload.CommodityType,
	}

	if err := repo.CommodityAttributes().Create(r.Context(), ca); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			middleware.WriteError(w, http.StatusConflict, "Commodity attribute with this name already exists")
			return
		}
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ca)
}
