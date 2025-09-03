package locations

// CreateLocationPayload defines the structure for creating a new location.
type CreateLocationPayload struct {
	Name      string `json:"name" validate:"required"`
	CompanyID int64  `json:"company_id" validate:"required"`
	AddressID int64  `json:"address_id" validate:"required"`
}

// UpdateLocationPayload defines the structure for updating a location.
type UpdateLocationPayload struct {
	Name      *string `json:"name,omitempty" validate:"required_without_all=AddressID"`
	AddressID *int64  `json:"address_id,omitempty" validate:"required_without_all=Name"`
}