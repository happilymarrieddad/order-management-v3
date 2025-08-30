package types

import "time"

// Commodity represents a commodity in the system.
//
// Examples of commodities might include "Potato", "Watermelon", etc.
type Commodity struct {
	ID            int64         `json:"id" xorm:"pk autoincr 'id'"`
	Name          string        `json:"name" xorm:"unique 'name'"`
	CommodityType CommodityType `json:"commodityType" xorm:"index 'commodity_type'"`
	Visible       bool          `json:"visible" xorm:"'visible'"`
	CreatedAt     time.Time     `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt     time.Time     `json:"updatedAt" xorm:"updated 'updated_at'"`
}

// TableName specifies the table name for the Commodity model.
func (Commodity) TableName() string {
	return "commodities"
}
