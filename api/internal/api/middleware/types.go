package middleware

// ErrorResponse is the standard structure for API error responses.
type ErrorResponse struct {
	Error string `json:"error" example:"a description of the error"`
}
