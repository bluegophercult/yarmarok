package web

import (
	"errors"
	"net/http"
	"strings"
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
	value   error
	code    int
	log     Log
	Message Message
}

type (
	Fields  map[string]any
	Message Fields
	Log     Fields
)

func (e *Error) Error() string {
	return e.value.Error()
}

func (e *Error) StatusCode() int {
	return e.code
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
		value:   err,
		code:    code,
		log:     log,
		Message: msg,
	}
}

func ErrorAs(err error) (*Error, bool) {
	return errorAs[*Error](err)
}

func errorAs[E error](err error) (errv E, ok bool) {
	return errv, errors.As(err, &errv)
}

func statusText(code int) string {
	return strings.ToLower(http.StatusText(code))
}
