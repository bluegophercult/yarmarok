package web

import (
	"errors"
)

// ErrUnknownError is returned when something is not ok.
var ErrUnknownError = errors.New("unknown error")

var _ error = (*Error)(nil)

// Error represents an application error
// that will be logged by WithErrors middleware
// and could return to the client as Message with Fields.
type Error struct {
	Value   error  `json:"-"`
	Code    int    `json:"-"`
	Message Fields `json:"message,omitempty"`
}

type Fields map[string]any

func NewError(err error, code int, fields ...Fields) error {
	newErr := Error{
		Value:   err,
		Code:    code,
		Message: make(Fields),
	}

	for i := range fields {
		for k, v := range fields[i] {
			newErr.Message[k] = v
		}
	}

	return &newErr
}

func (e *Error) Error() string {
	return e.Value.Error()
}

func (e *Error) StatusCode() int {
	return e.Code
}

func ErrorIs(err error) bool {
	_, ok := errorAs[*Error](err)
	return ok
}

func ErrorAs(err error) (*Error, bool) {
	return errorAs[*Error](err)
}

func errorAs[E error](err error) (errv E, ok bool) {
	return errv, errors.As(err, &errv)
}
