package types

import "time"

// Address represents a physical address, mirroring the 'addresses' table schema.
type Address struct {
	ID         int64     `validate:"-" json:"id" xorm:"pk autoincr 'id'"`
	Line1      string    `validate:"required" json:"line1" xorm:"notnull 'line_1'"`
	Line2      string    `validate:"-" json:"line2,omitempty" xorm:"'line_2'"`
	City       string    `validate:"required" json:"city" xorm:"notnull 'city'"`
	State      string    `validate:"required" json:"state" xorm:"notnull 'state'"`
	PostalCode string    `validate:"required" json:"postalCode" xorm:"notnull 'postal_code'"`
	Country    string    `validate:"required" json:"country" xorm:"notnull 'country'"`
	GlobalCode string    `validate:"-" json:"globalCode" xorm:"'global_code'"`
	CreatedAt  time.Time `validate:"-" json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt  time.Time `validate:"-" json:"updatedAt" xorm:"updated 'updated_at'"`
}

// TableName specifies the table name for the Address model.
// This is used by XORM to map this struct to the 'addresses' database table.
func (Address) TableName() string {
	return "addresses"
}
