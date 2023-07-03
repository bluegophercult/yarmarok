// Package resperr provides an extended error functionality
// to be used for structured error responses.
package resperr

import "encoding/json"

// ResponseError is an error that can be used in a response.
type ResponseError interface {
	error
	Message() string
	Code() string
}

// JSONError is an error that can be marshaled to JSON.
type JSONError interface {
	error
	json.Marshaler
}

// Error is a structured error.
type Error struct {
	code    string
	message string
}

// NewError creates a new error.
func NewError(code, message string) *Error {
	return &Error{
		code:    code,
		message: message,
	}
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.message
}

// Code returns the error code.
func (e *Error) Code() string {
	return e.code
}

// MarshalJSON marshals the error to JSON.
func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"code":    e.code,
		"message": e.message,
	})
}

// Is reports whether the error is of the same type as target.
// It matches the pointer of the error so it will not match
// if the target has same code and message but
// the underlying error instance is not the original one.
func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}

	if targetError, ok := target.(*Error); ok {
		return targetError == e
	}

	return false
}

var _ error = (*Error)(nil)
