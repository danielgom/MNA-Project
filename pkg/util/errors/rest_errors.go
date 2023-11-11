// Package errors contains a common way to declare API errors
package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

const (
	ErrNotFound = "no rows in result set"

	errDatabase = "Database error"
)

var errInvalidJSON = errors.New("invalid json")

// CommonError is a common interface where we can create different kind of errors.
type CommonError interface {
	Message() string
	Status() int
	Error() string
	Cause() interface{}
}

type restError struct {
	M string      `json:"message"`
	S int         `json:"status"`
	E string      `json:"error"`
	C interface{} `json:"causes,omitempty"`
}

func (r *restError) Message() string {
	return r.M
}

func (r *restError) Status() int {
	return r.S
}

func (r *restError) Error() string {
	return r.E
}

func (r *restError) Cause() interface{} {
	return r.C
}

// NewRestError creates an error based on user input.
func NewRestError(message string, status int, errStr string, err error) CommonError {
	r := &restError{
		M: message,
		S: status,
		E: errStr,
		C: message,
	}
	if err != nil {
		r.C = err.Error()
	}

	return r
}

// NewRestErrorFromBytes is based on upcoming bytes.
func NewRestErrorFromBytes(b []byte) (CommonError, error) {
	var r restError
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errInvalidJSON
	}
	return &r, nil
}

// NewBadRequestError returns a bad request CommonError.
func NewBadRequestError(m string, err error) CommonError {
	r := &restError{
		M: m,
		S: http.StatusBadRequest,
		E: "Bad Request",
		C: m,
	}

	if err != nil {
		r.C = err.Error()
	}

	return r
}

// NewNotFoundError returns a bad request CommonError
// We normally error not found errors based on database response.
func NewNotFoundError(m string) CommonError {
	return &restError{
		M: m,
		S: http.StatusNotFound,
		E: "Not found",
	}
}

// NewInternalServerError returns an internal server error (such as DB calls or service-to-service calls).
func NewInternalServerError(m string, err error) CommonError {
	r := &restError{
		M: m,
		S: http.StatusInternalServerError,
		E: "Internal server error",
		C: m,
	}

	if err != nil {
		r.C = err.Error()
	}
	return r
}

// NewUnauthorisedError returns an Unauthorised error (no access to the api).
func NewUnauthorisedError(m string) CommonError {
	return &restError{
		M: m,
		S: http.StatusUnauthorized,
		E: "Unauthorised",
	}
}

func WrapDatabaseError(err error) CommonError {
	return NewInternalServerError(errDatabase, err)
}
