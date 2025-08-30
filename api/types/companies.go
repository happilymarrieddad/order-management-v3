package types

import "time"

// Company represents a company in the system.
type Company struct {
	ID                 int64     `xorm:"pk autoincr 'id'" json:"id"`
	Name               string    `xorm:"varchar(255) notnull 'name'" json:"name" validate:"required"`
	AddressID          int64     `xorm:"notnull 'address_id'" json:"address_id" validate:"required"`
	OrderPrefix        string    `xorm:"'order_prefix'" json:"order_prefix"`
	OrderPostfix       string    `xorm:"'order_postfix'" json:"order_postfix"`
	DefaultOrderNumber int       `xorm:"'default_order_number'" json:"default_order_number"`
	Visible            bool      `xorm:"'visible'" json:"-"`
	CreatedAt          time.Time `xorm:"created" json:"created_at"`
	UpdatedAt          time.Time `xorm:"updated" json:"updated_at"`

	// Relations
	Address *Address `json:"address,omitempty" xorm:"-"`
}

// TableName specifies the table name for the Company model.
// This is used by XORM to map this struct to the 'companies' database table.
func (Company) TableName() string {
	return "companies"
}
