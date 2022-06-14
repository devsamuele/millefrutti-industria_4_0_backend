package web

import (
	"github.com/pkg/errors"
)

type ErrorResponse struct {
	Error RequestError `json:"error"`
}

// RequestError ...
type RequestError struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	Reason       string `json:"reason,omitempty"`
	LocationType string `json:"locationType,omitempty"` // argument, parameter, ...
	Location     string `json:"location,omitempty"`     // argumentName, argument parameter...
}

// NewRequestError ...
func NewRequestError(err error, code int, reason, locationType, location string) error {
	return &RequestError{
		Message:      err.Error(),
		Code:         code,
		Reason:       reason,
		LocationType: locationType,
		Location:     location,
	}
}

func (re *RequestError) Error() string {
	return re.Message
}

// IsRequestError checks if an error of type RequestError exists.
func IsRequestError(err error) bool {
	var re *RequestError
	return errors.As(err, &re)
}

// GetRequestError returns a copy of the RequestError pointer.
func GetRequestError(err error) *RequestError {
	var re *RequestError
	if !errors.As(err, &re) {
		return nil
	}
	return re
}

// shutdown ...
type shutdown struct {
	Message string
}

// NewShutdownError ...
func NewShutdownError(message string) error {
	return &shutdown{Message: message}
}

// Error ...
func (e *shutdown) Error() string {
	return e.Message
}

// IsShutdown ...
func IsShutdown(err error) bool {
	if _, ok := errors.Cause(err).(*shutdown); ok {
		return true
	}
	return false
}
