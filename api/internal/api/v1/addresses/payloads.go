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
type UpdateAddressPayload struct {
	Line1      *string `json:"line_1"`
	Line2      *string `json:"line_2"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	PostalCode *string `json:"postal_code"`
	Country    *string `json:"country"`
}
