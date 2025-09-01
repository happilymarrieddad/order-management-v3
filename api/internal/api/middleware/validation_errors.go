package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FormatValidationErrors takes an error and returns a user-friendly string
// if the error is a validator.ValidationErrors type. Otherwise, it returns
// a generic error message.
func FormatValidationErrors(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var errorMessages []string
		for _, fe := range ve {
			switch fe.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' is required.", strings.ToLower(fe.Field())))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' has an invalid value.", strings.ToLower(fe.Field())))
			}
		}
		return strings.Join(errorMessages, ", ")
	}
	return "An unexpected validation error occurred."
}
