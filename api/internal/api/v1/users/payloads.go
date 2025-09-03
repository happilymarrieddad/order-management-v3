package users

// CreateUserPayload defines the structure for creating a new user.
type CreateUserPayload struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	CompanyID       int64  `json:"company_id" validate:"required"`
	AddressID       int64  `json:"address_id" validate:"required"`
}

// UpdateUserPayload defines the structure for updating a user.
// At least one field must be provided.
type UpdateUserPayload struct {
	FirstName *string `json:"first_name,omitempty" validate:"required_without_all=LastName AddressID"`
	LastName  *string `json:"last_name,omitempty"  validate:"required_without_all=FirstName AddressID"`
	AddressID *int64  `json:"address_id,omitempty" validate:"required_without_all=FirstName LastName"`
}

// UpdateUserCompanyPayload defines the structure for updating a user's company.
type UpdateUserCompanyPayload struct {
	CompanyID int64 `json:"company_id" validate:"required"`
}