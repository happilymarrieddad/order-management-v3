package types

import "time"

// Product represents a specific product for a company in the system.
type Product struct {
	ID          int64     `json:"id" xorm:"pk autoincr 'id'"`
	CommodityID int64     `json:"commodityId" xorm:"notnull index 'commodity_id'"`
	CompanyID   int64     `json:"companyId" xorm:"notnull index 'company_id'"`
	Name        string    `json:"name" xorm:"'name'"` // Derived name for the product
	Visible     bool      `xorm:"'visible'" json:"-"`
	CreatedAt   time.Time `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt   time.Time `json:"updatedAt" xorm:"updated 'updated_at'"`

	ProductAttributeValues []*ProductAttributeValue `xorm:"-" json:"attributes,omitempty"`
}

// TableName specifies the table name for the Product model.
func (Product) TableName() string {
	return "products"
}
