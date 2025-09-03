package types

import "time"

// User represents a user in the system.
type User struct {
	ID        int64     `json:"id" xorm:"pk autoincr 'id'"`
	FirstName string    `validate:"required,min=2,max=50" json:"firstName" xorm:"'first_name'"`
	LastName  string    `validate:"required,min=2,max=50" json:"lastName" xorm:"'last_name'"`
	Email     string    `validate:"required,email" json:"email" xorm:"unique 'email'"`
	Password  string    `validate:"required,min=8" json:"-" xorm:"'password'"`
	CompanyID int64     `validate:"required" json:"companyId" xorm:"notnull index 'company_id'"`
	AddressID int64     `validate:"required" json:"addressId" xorm:"notnull index 'address_id'"`
	Visible   bool      `xorm:"'visible'" json:"-"`
	Roles     Roles     `json:"roles" xorm:"'roles'"`
	CreatedAt time.Time `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"updated 'updated_at'"`

	// Relations
	Address *Address `json:"address,omitempty" xorm:"-"`
}

// TableName specifies the table name for the User model.
// This is used by the XORM to map this struct to the 'users' database table.
func (User) TableName() string {
	return "users"
}

// HasRole checks if a specific role exists in the user Roles.
func (u User) HasRole(role Role) bool {
	for _, existingRole := range u.Roles {
		if existingRole == role {
			return true
		}
	}
	return false
}
