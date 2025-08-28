package locations

// CreateLocationPayload defines the structure for creating a new location.
type CreateLocationPayload struct {
	CompanyID int64  `json:"company_id" validate:"required" example:"1"`
	AddressID int64  `json:"address_id" validate:"required" example:"1"`
	Name      string `json:"name" validate:"required" example:"Main Warehouse"`
}

// UpdateLocationPayload defines the structure for updating a location.
type UpdateLocationPayload struct {
	AddressID int64  `json:"address_id" validate:"required" example:"2"`
	Name      string `json:"name" validate:"required" example:"Downtown Office"`
}
