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
	// DefaultCode is the default code for an error.
	DefaultCode = "UnexpectedError"
)

// AsJSON returns the JSON representation of the error.
func AsJSON(err error) ([]byte, error) {
	response := AsResponse(err)
	return json.Marshal(response)
}

// Response is an API friendly error response.
type Response struct {
	Code string `json:"errorCode"`
}

// DefaultResponse returns a new response with the default values.
func DefaultResponse() Response {
	return Response{
		Code: DefaultCode,
	}
}

// Coder Is an interface that errors should implement
// to make AsResponse function unpack the code from the error.
type Coder interface {
	Code() string
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

	return response
}

// Error is an error that is represented by a code.
type Error struct {
	error
}

// New returns a new error with the given code.
func New(code string) error {
	return &Error{
		error: errors.New(code),
	}
}

// Code returns the attached code.
func (e *Error) Code() string {
	return e.Error()
}

// WithCode returns a new error with the given code attached to it.
func WithCode(code string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", New(code), err)
}

func (e *LabeledError) Unwrap() error {
	return e.error
}

// LabeledError is an error that can have labels attached to it.
type LabeledError struct {
	error
	labels []Label
}

// WithLabel returns a new error with the given label attached to it.
func WithLabel(err error, key, value string) *LabeledError {
	if err == nil {
		return nil
	}

	labels := []Label{{Key: key, Value: value}}

	return &LabeledError{error: err, labels: labels}
}

// Error returns the error message with the attached labels.
func (e *LabeledError) Error() string {
	if len(e.labels) == 0 {
		return e.error.Error()
	}

	return fmt.Sprintf("%s: %s", e.formatLabels(), e.error.Error())
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

	return &LabeledError{error: err, labels: labels}
}
