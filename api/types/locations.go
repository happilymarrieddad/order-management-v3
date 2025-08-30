package types

import "time"

// Location represents a company's physical location in the system.
type Location struct {
	ID        int64     `json:"id" xorm:"pk autoincr 'id'"`
	CompanyID int64     `json:"companyId" xorm:"notnull 'company_id'"`
	AddressID int64     `json:"addressId" xorm:"notnull 'address_id'"`
	Name      string    `json:"name" xorm:"notnull 'name'"`
	Visible   bool      `xorm:"'visible'" json:"-"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"updated"`

	// Relations (for API responses)
	Company *Company `json:"company,omitempty" xorm:"-"`
	Address *Address `json:"address,omitempty" xorm:"-"`
}

// TableName specifies the table name for this struct.
func (Location) TableName() string {
	return "locations"
}
