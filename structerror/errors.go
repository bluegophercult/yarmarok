// Package structerror provides extended errors functionality.
// It allows to attach a code and labels to an error.
package structerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	DefaultCode    = "UnexpectedError"
	DefaultMessage = "An unexpected error occurred"
)

// AsJSON returns the JSON representation of the error.
func AsJSON(err error) ([]byte, error) {
	var marshaler json.Marshaler
	if !errors.As(err, &marshaler) {
		return nil, fmt.Errorf("error does not implement json.Marshaler")
	}

	return marshaler.MarshalJSON()
}

// APIError is a type of error that handles its API representation.
type APIError interface {
	error
	Code() string
	Message() string
}

// Response is an API friendly error response.
type Response struct {
	Code    string `json:"code"`
	Message string `json:"error"`
}

// DefaultResponse returns a new response with the default values.
func DefaultResponse() Response {
	return Response{
		Code:    DefaultCode,
		Message: DefaultMessage,
	}
}

// Coder Is an interface that errors should implement
// to make AsResponse function unpack the code from the error.
type Coder interface {
	Code() string
}

// Messager Is an interface that errors should implement
// to make AsResponse function unpack the message from the error.
type Messager interface {
	Message() string
}

// AsResponse returns the API representation of the error.
// It tries to unpack the code and message from the error.
// If the error does not implement the Coder or Messager interfaces,
// it returns the default response.
func AsResponse(err error) Response {
	response := DefaultResponse()

	var coder Coder
	if errors.As(err, &coder) {
		response.Code = coder.Code()
	}

	var messager Messager
	if errors.As(err, &messager) {
		response.Message = messager.Message()
	}

	return response
}

type base struct {
	error
}

func newBase(err error) *base {
	return &base{error: err}
}

func (b *base) Unwrap() error {
	return b.error
}

// CodeError is an error that can have a code attached to it.
type CodeError struct {
	*base
	code string
}

// NewCodeError returns a new error with the given code and message.
func NewCodeError(code, format string, args ...any) error {
	return &CodeError{
		code: code,
		base: newBase(
			fmt.Errorf(format, args...),
		),
	}
}

// Error returns the error message with the attached code.
func (e *CodeError) Error() string {
	if e.code == "" {
		return e.base.Error()
	}

	if e.base.Error() == "" {
		return e.code
	}

	return e.code + ": " + e.base.Error()
}

// MarshalJSON returns the JSON representation of the error.
func (e *CodeError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"code":  e.code,
		"error": e.base.Error(),
	})
}

// Code returns the attached code.
func (e *CodeError) Code() string {
	return e.code
}

// WithCode returns a new error with the given code attached to it.
func WithCode(code string, err error) error {
	if err == nil {
		return nil
	}
	return &CodeError{code: code, base: newBase(err)}
}

// LabeledError is an error that can have labels attached to it.
type LabeledError struct {
	*base
	labels []Label
}

// WithLabel returns a new error with the given label attached to it.
func WithLabel(err error, key, value string) *LabeledError {
	if err == nil {
		return nil
	}

	labels := []Label{{Key: key, Value: value}}

	return &LabeledError{base: newBase(err), labels: labels}
}

// Error returns the error message with the attached labels.
func (e *LabeledError) Error() string {
	if len(e.labels) == 0 {
		return e.base.Error()
	}

	return fmt.Sprintf("%s: %s", e.formatLabels(), e.base.Error())
}

func (e *LabeledError) formatLabels() string {
	strs := make([]string, 0, len(e.labels))
	for _, label := range e.labels {
		strs = append(strs, label.String())
	}
	return strings.Join(strs, ", ")
}

// Label is a key-value pair that can be attached to an error.
type Label struct {
	Key   string
	Value any
}

// String returns the string representation of the label.
func (l Label) String() string {
	return fmt.Sprintf("%s=%s", l.Key, l.Value)
}

// KV returns a new label with the given key and value.
func KV(key string, value any) Label {
	return Label{Key: key, Value: value}
}

// WithLabels returns a new error with the given labels attached to it.
func WithLabels(err error, labels ...Label) *LabeledError {
	if err == nil {
		return nil
	}

	return &LabeledError{base: newBase(err), labels: labels}
}
