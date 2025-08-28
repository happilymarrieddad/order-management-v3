package users

// CreateUserPayload defines the structure for creating a new user.
type CreateUserPayload struct {
	Email           string `json:"email" validate:"required,email" example:"test@example.com"`
	Password        string `json:"password" validate:"required,min=8" example:"password123"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password" example:"password123"`
	FirstName       string `json:"first_name" validate:"required" example:"John"`
	LastName        string `json:"last_name" validate:"required" example:"Doe"`
	CompanyID       int64  `json:"company_id" validate:"required" example:"1"`
	AddressID       int64  `json:"address_id" validate:"required" example:"1"`
}

// UpdateUserPayload defines the structure for updating a user.
type UpdateUserPayload struct {
	FirstName string `json:"first_name" validate:"required" example:"Jane"`
	LastName  string `json:"last_name" validate:"required" example:"Doe"`
	AddressID int64  `json:"address_id" validate:"required" example:"2"`
}
