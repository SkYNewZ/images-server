package internal

import (
	"fmt"
	"net/http"
)

var _ error = (*Error)(nil)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// Error implements error
func (e Error) Error() string {
	return fmt.Sprintf("code=%d message=%s", e.Code, e.Message)
}

func newNotFoundError(err error) *Error {
	return &Error{
		Code:    http.StatusNotFound,
		Message: err.Error(),
	}
}

func newBadRequestError(err error) *Error {
	return &Error{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}
}

func newUnsupportedMediaType(err error) *Error {
	return &Error{
		Code:    http.StatusUnsupportedMediaType,
		Message: err.Error(),
	}
}
