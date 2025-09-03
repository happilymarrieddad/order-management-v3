package types

import "time"

// CompanyAttribute represents the link between a company and a commodity attribute,
// indicating that the company uses this attribute for its products.
type CompanyAttribute struct {
	ID                   int64     `json:"id" xorm:"pk autoincr 'id'"`
	Position             int       `json:"position" xorm:"notnull 'position'"`
	CommodityAttributeID int64     `json:"commodity_attribute_id" xorm:"notnull index 'commodity_attribute_id'"`
	CompanyID            int64     `json:"company_id" xorm:"notnull index 'company_id'"`
	CreatedAt            time.Time `json:"created_at" xorm:"created 'created_at'"`
	UpdatedAt            time.Time `json:"updated_at" xorm:"updated 'updated_at'"`
}

// TableName specifies the table name for the CompanyAttribute model.
func (CompanyAttribute) TableName() string {
	return "company_attributes"
}
