package web

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrReasonRequired         = "required"
	ErrReasonInternalError    = "internalError"
	ErrReasonForbidden        = "forbidden"
	ErrReasonInvalidArgument  = "invalidArgument"  // body
	ErrReasonInvalidParameter = "invalidParameter" //url
	ErrReasonConflict         = "conflict"
	ErrReasonNotFound         = "notFound"
)

func ErrHandler(err error) error {
	var sErr *serviceError

	statusCode := 0

	if errors.As(err, &sErr) {
		switch sErr.Reason {
		case ErrReasonForbidden:
			statusCode = http.StatusForbidden
		case ErrReasonNotFound:
			statusCode = http.StatusNotFound
		case ErrReasonRequired, ErrReasonInvalidParameter, ErrReasonInvalidArgument:
			statusCode = http.StatusBadRequest
		case ErrReasonConflict:
			statusCode = http.StatusConflict
		case ErrReasonInternalError:
			statusCode = http.StatusInternalServerError
		default:
			return err
		}
		return NewRequestError(sErr, statusCode, sErr.Reason, sErr.LocationType, sErr.Location)
	}

	return err
}

type serviceError struct {
	Message      string `json:"message"`
	Reason       string `json:"reason,omitempty"`
	LocationType string `json:"locationType,omitempty"` // argument, parameter, ...
	Location     string `json:"location,omitempty"`     // argumentName, argument parameter...
}

func (e *serviceError) Error() string {
	return e.Message
}

func (e serviceError) NewRequestError(code int) error {
	return &requestError{
		serviceError: e,
		Code:         code,
	}
}

func NewError(message, reason, locationType, location string) error {
	return &serviceError{
		Message:      message,
		Reason:       reason,
		LocationType: locationType,
		Location:     location,
	}
}

// RequestError ...
type requestError struct {
	Code int `json:"code"`
	serviceError
}

// NewRequestError ...
func NewRequestError(err error, code int, reason, locationType, location string) error {
	return &requestError{
		Code: code,
		serviceError: serviceError{
			Message:      err.Error(),
			Reason:       reason,
			LocationType: locationType,
			Location:     location,
		},
	}
}

func (e *requestError) Error() string {
	return e.Message
}

// IsRequestError checks if an error of type RequestError exists.
func IsRequestError(err error) bool {
	var re *requestError
	return errors.As(err, &re)
}

// GetRequestError returns a copy of the RequestError pointer.
func GetRequestError(err error) *requestError {
	var re *requestError
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
