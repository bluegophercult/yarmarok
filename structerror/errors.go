package structerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// AsJSON returns the JSON representation of the error.
func AsJSON(err error) ([]byte, error) {
	var marshaler json.Marshaler
	if !errors.As(err, &marshaler) {
		return nil, fmt.Errorf("error does not implement json.Marshaler")
	}

	return marshaler.MarshalJSON()
}

// CodeError is an error that can have a code attached to it.
type CodeError struct {
	error
	code string
}

// New returns a new error with the given code and message.
func New(code, format string, args ...any) error {
	return &CodeError{code: code, error: fmt.Errorf(format, args...)}
}

// Error returns the error message with the attached code.
func (e *CodeError) Error() string {
	if e.code == "" {
		return e.error.Error()
	}

	if e.error.Error() == "" {
		return e.code
	}

	return e.code + ": " + e.error.Error()
}

// MarshalJSON returns the JSON representation of the error.
func (e *CodeError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"code":  e.code,
		"error": e.error.Error(),
	})
}

// Unwrap returns the wrapped error.
func (e *CodeError) Unwrap() error {
	return e.error
}

func (e *CodeError) Code() string {
	return e.code
}

// WithCode returns a new error with the given code attached to it.
func WithCode(code string, err error) error {
	if err == nil {
		return nil
	}
	return &CodeError{code: code, error: err}
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

// Unwrap returns the wrapped error.
func (e *LabeledError) Unwrap() error {
	return e.error
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
