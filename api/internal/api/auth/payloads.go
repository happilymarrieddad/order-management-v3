package auth

// LoginPayload defines the structure for a login request.
type LoginPayload struct {
	Email    string `json:"email" validate:"required,email" example:"test@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// LoginResponse defines the structure for a successful login response.
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
