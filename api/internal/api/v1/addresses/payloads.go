package addresses

// CreateAddressPayload defines the expected JSON structure for creating a new address.
type CreateAddressPayload struct {
	Line1      string `json:"line_1" validate:"required"`
	Line2      string `json:"line_2"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
	Country    string `json:"country" validate:"required"`
}

// UpdateAddressPayload defines the expected JSON structure for updating an address.
// Pointers are used to distinguish between a field not being provided and a field being set to its zero value.
// At least one field must be provided.
type UpdateAddressPayload struct {
	Line1      *string `json:"line_1,omitempty" validate:"required_without_all=Line2 City State PostalCode Country"`
	Line2      *string `json:"line_2,omitempty"`
	City       *string `json:"city,omitempty" validate:"required_without_all=Line1 Line2 State PostalCode Country"`
	State      *string `json:"state,omitempty" validate:"required_without_all=Line1 Line2 City PostalCode Country"`
	PostalCode *string `json:"postal_code,omitempty" validate:"required_without_all=Line1 Line2 City State Country"`
	Country    *string `json:"country,omitempty" validate:"required_without_all=Line1 Line2 City State PostalCode"`
}
