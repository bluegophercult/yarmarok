package service

import (
	"regexp"

	"github.com/go-playground/validator"
	"github.com/kaznasho/yarmarok/structerror"
)

var (
	ErrParticipantPhoneOnlyDigits = structerror.New("PhoneContainsLetters")
	ErrNameTooShort               = structerror.New("NameTooShort")
	ErrInvalidRequest             = structerror.New("InvalidRequest")
)

// Acceptable characters are the English alphabet, numbers,
// symbols from the symbols variable, the Ukrainian alphabet
func charsValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	regex := regexp.MustCompile(`^[a-zA-Z0-9 !@#$%^&*()_{}\[\]:;<>,.?~абвгґдеєжзиіїйклмнопрстуфхцчшщьюяАБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯ]*$`)
	return regex.MatchString(value)
}

func phoneValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	regex := regexp.MustCompile(`^\+380\d{9,10}$`)
	return regex.MatchString(value)
}

func defaultValidator() *validator.Validate {
	validate := validator.New()

	if err := validate.RegisterValidation("charsValidation", charsValidation); err != nil {
		panic(err)
	}

	if err := validate.RegisterValidation("phoneValidation", phoneValidation); err != nil {
		panic(err)
	}

	return validate
}
