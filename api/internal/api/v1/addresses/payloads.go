package addresses

// CreateAddressPayload defines the structure for creating a new address.
type CreateAddressPayload struct {
	Line1      string  `json:"line_1" validate:"required" example:"123 Main St"`
	Line2      *string `json:"line_2" example:"Apt 4B"`
	City       string  `json:"city" validate:"required" example:"Anytown"`
	State      string  `json:"state" validate:"required" example:"CA"`
	PostalCode string  `json:"postal_code" validate:"required" example:"12345"`
}

// UpdateAddressPayload defines the structure for updating an address.
type UpdateAddressPayload struct {
	Line1      string  `json:"line_1" validate:"required" example:"456 Market St"`
	Line2      *string `json:"line_2" example:"Suite 200"`
	City       string  `json:"city" validate:"required" example:"Newville"`
	State      string  `json:"state" validate:"required" example:"NY"`
	PostalCode string  `json:"postal_code" validate:"required" example:"54321"`
}
