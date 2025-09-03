package commodities

import "github.com/happilymarrieddad/order-management-v3/api/types"

// CreateCommodityPayload defines the structure for creating a new commodity.
type CreateCommodityPayload struct {
	Name          string              `json:"name" validate:"required" example:"Apple"`
	CommodityType types.CommodityType `json:"commodity_type" validate:"required" example:"1"`
}

// UpdateCommodityPayload defines the structure for updating a commodity.
type UpdateCommodityPayload struct {
	Name          *string              `json:"name,omitempty" validate:"required_without_all=CommodityType"`
	CommodityType *types.CommodityType `json:"commodity_type,omitempty" validate:"required_without_all=Name"`
}
