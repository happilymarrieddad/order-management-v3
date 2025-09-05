package types

import "time"

// CompanyAttributeSetting defines the display order of a commodity attribute for a specific company.
type CompanyAttributeSetting struct {
	ID                   int64     `json:"id" xorm:"pk autoincr 'id'"`
	CompanyID            int64     `json:"companyId" xorm:"notnull 'company_id'"`
	CommodityAttributeID int64     `json:"commodityAttributeId" xorm:"notnull 'commodity_attribute_id'"`
	DisplayOrder         int       `json:"displayOrder" xorm:"'display_order'"`
	CreatedAt            time.Time `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt            time.Time `json:"updatedAt" xorm:"updated 'updated_at'"`
}

func (CompanyAttributeSetting) TableName() string {
	return "company_attribute_settings"
}
