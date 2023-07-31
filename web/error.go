package web

import (
	"errors"
)

var (
	// ErrUnknownError is returned when something is not ok.
	ErrUnknownError = errors.New("unknown error")
	// ErrUnsupportedField is returned when a field is not supported.
	ErrUnsupportedField = errors.New("unsupported field")
)

var _ error = (*Error)(nil)

// Error represents an application error
// that will be logged by WithErrors middleware
// and could return to as Message with Fields.
type Error struct {
	Value   error   `json:"-"`
	Code    int     `json:"-"`
	Log     Log     `json:"-"`
	Message Message `json:"message,omitempty"`
}

type (
	Fields  map[string]any
	Message Fields
	Log     Fields
)

func (e *Error) Error() string {
	return e.Value.Error()
}

func (e *Error) StatusCode() int {
	return e.Code
}

func NewError(err error, code int, fields ...any) *Error {
	msg := make(Message)
	log := make(Log)

	for _, f := range fields {
		switch v := f.(type) {
		case Message:
			msg = v
		case Log:
			log = v
		default:
			err = errors.Join(err, ErrUnsupportedField)
		}
	}

	return &Error{
		Value:   err,
		Code:    code,
		Log:     log,
		Message: msg,
	}
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
