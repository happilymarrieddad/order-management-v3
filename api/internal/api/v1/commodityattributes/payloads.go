package commodityattributes

import "github.com/happilymarrieddad/order-management-v3/api/types"

// CreateCommodityAttributePayload defines the request body for creating a new commodity attribute.
type CreateCommodityAttributePayload struct {
	Name          string            `json:"name" validate:"required,min=2,max=255"`
	CommodityType types.CommodityType `json:"commodityType" validate:"required,oneof=1"`
}

// UpdateCommodityAttributePayload defines the request body for updating an existing commodity attribute.
type UpdateCommodityAttributePayload struct {
	Name          string            `json:"name" validate:"required,min=2,max=255"`
	CommodityType types.CommodityType `json:"commodityType" validate:"required,oneof=1"`
}
