package companies

// CreateCompanyPayload defines the structure for creating a new company.
type CreateCompanyPayload struct {
	Name      string `json:"name" validate:"required" example:"Awesome Inc."`
	AddressID int64  `json:"address_id" validate:"required" example:"1"`
}

// UpdateCompanyPayload defines the structure for updating a company.
type UpdateCompanyPayload struct {
	Name      string `json:"name" validate:"required" example:"Even Better Inc."`
	AddressID int64  `json:"address_id" validate:"required" example:"2"`
}
