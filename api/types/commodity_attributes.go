package types

import "time"

// CommodityAttribute represents an attribute that can be associated with a commodity type.
//
// Examples of attributes might include "Color", "Size", "Weight", etc.
type CommodityAttribute struct {
	ID            int64         `json:"id" xorm:"pk autoincr 'id'"`
	Name          string        `json:"name" xorm:"unique 'name'"` // Assuming attribute names are unique
	CommodityType CommodityType `json:"commodityType" xorm:"index 'commodity_type_id'"`
	CreatedAt     time.Time     `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt     time.Time     `json:"updatedAt" xorm:"updated 'updated_at'"`
}

// TableName specifies the table name for the CommodityAttribute model.
func (CommodityAttribute) TableName() string {
	return "commodity_attributes"
}
