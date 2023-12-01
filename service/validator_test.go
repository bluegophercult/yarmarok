package service

import (
	"testing"
)

func TestDefaultValidator(t *testing.T) {
	validate := defaultValidator()

	type testStruct struct {
		Name  string `validate:"charsValidation"`
		Phone string `validate:"phoneValidation"`
	}

	if err := validate.Struct(&testStruct{Name: "test", Phone: "+380123456789"}); err != nil {
		t.Error(err)
	}
}
