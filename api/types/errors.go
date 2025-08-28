package types

import (
	"fmt"
	"strings"
)

const (
	notFoundErr       string = "Not Found:"
	unauthorizedErr   string = "Unauthorized:"
	notImplementedErr string = "Not Implemented:"
	badRequest        string = "Bad Request:"
	intErr            string = "Internal Error:"
)

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), notFoundErr)
}

func NewNotFoundError(msg string) error {
	return fmt.Errorf("%s %s", notFoundErr, msg)
}

func IsUnauthorizedError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), unauthorizedErr)
}

func NewUnauthorizedError(msg string) error {
	return fmt.Errorf("%s %s", unauthorizedErr, msg)
}

func IsNotImplementedError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), notImplementedErr)
}

func NewNotImplementedError() error {
	return fmt.Errorf("%s", notImplementedErr)
}

func IsBadRequestError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), badRequest)
}

func NewBadRequestError(msg string) error {
	return fmt.Errorf("%s %s", badRequest, msg)
}

func IsInternalServerErrorError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), intErr)
}

func NewInternalServerError(msg string) error {
	return fmt.Errorf("%s %s", intErr, msg)
}
