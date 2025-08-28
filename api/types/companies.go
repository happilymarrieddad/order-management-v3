package types

import "time"

// Company represents a company in the system.
type Company struct {
	ID        int64     `json:"id" xorm:"pk autoincr 'id'"`
	Name      string    `json:"name" xorm:"notnull unique 'name'" validate:"required"`
	AddressID int64     `json:"addressId" xorm:"notnull index 'address_id'" validate:"required"`
	CreatedAt time.Time `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"updated 'updated_at'"`

	// Relations
	Address *Address `json:"address,omitempty" xorm:"-"`
}

// TableName specifies the table name for the Company model.
// This is used by XORM to map this struct to the 'companies' database table.
func (Company) TableName() string {
	return "companies"
}
