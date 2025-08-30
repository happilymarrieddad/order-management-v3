package types

import "time"

// ProductAttributeValue represents a specific value for a CommodityAttribute on a Product.
type ProductAttributeValue struct {
	ID                   int64     `json:"id" xorm:"pk autoincr 'id'"`
	ProductID            int64     `json:"productId" xorm:"notnull index 'product_id'"`
	CompanyID            int64     `json:"companyId" xorm:"notnull index 'company_id'"`
	CommodityAttributeID int64     `json:"commodityAttributeId" xorm:"notnull index 'commodity_attribute_id'"`
	Value                string    `json:"value" xorm:"'value'"`
	CreatedAt            time.Time `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt            time.Time `json:"updatedAt" xorm:"updated 'updated_at'"`
}

// TableName specifies the table name for the CommodityBatchAttributeValue model.
func (ProductAttributeValue) TableName() string {
	return "product_attribute_values"
}
