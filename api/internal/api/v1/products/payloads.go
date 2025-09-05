package products

import (
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// CreateProductPayload represents the request body for creating a new product.
type CreateProductPayload struct {
	CompanyID   int64  `json:"company_id" validate:"required"`
	CommodityID int64  `json:"commodity_id" validate:"required,gt=0"`
	Attributes  []*types.ProductAttributeValue `json:"attributes"`
}

// UpdateProductPayload represents the request body for updating an existing product.
type UpdateProductPayload struct {
	CommodityID *int64  `json:"commodity_id" validate:"required_without_all=Attributes,gt=0"`
	Attributes  []*types.ProductAttributeValue `json:"attributes" validate:"required_without_all=CommodityID"`
}
