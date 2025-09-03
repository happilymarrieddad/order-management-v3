package companies

// CreateCompanyPayload defines the structure for creating a new company.
type CreateCompanyPayload struct {
	Name      string `json:"name" validate:"required"`
	AddressID int64  `json:"address_id" validate:"required"`
}

// UpdateCompanyPayload defines the structure for updating a company.
// At least one field must be provided.
type UpdateCompanyPayload struct {
	Name      *string `json:"name,omitempty" validate:"required_without_all=AddressID"`
	AddressID *int64  `json:"address_id,omitempty" validate:"required_without_all=Name"`
}
