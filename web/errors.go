package web

import (
	"errors"
)

var ErrUnknownError = errors.New("unknown error")

type Error struct {
	Value   error   `json:"-"`
	Code    int     `json:"-"`
	Message Details `json:"message,omitempty"`
}

type Details map[string]any

func (e *Error) With(key string, val any) *Error {
	e.Message[key] = val
	return e
}

func NewError(err error, code int) *Error {
	return &Error{
		Value:   err,
		Code:    code,
		Message: make(Details),
	}
}

func (e *Error) Error() string {
	return e.Value.Error()
}

func (e *Error) StatusCode() int {
	return e.Code
}

func ErrorIs(err error) bool {
	_, ok := ErrorAs[*Error](err)
	return ok
}

func ErrorAs[E error](err error) (errv E, ok bool) {
	return errv, errors.As(err, &errv)
}

var _ error = (*Error)(nil)
