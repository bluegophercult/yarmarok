package service

import (
	"errors"
	"fmt"
	"regexp"

	"strings"

	"github.com/go-playground/validator"
)

const (
	minParticipantPhoneLength = 10
	maxParticipantPhoneLength = 12
)

var (
	ErrParticipantPhoneOnlyDigits = errors.New("phone should contain only digits")
	ErrNameTooShort               = errors.New("name is too short")
	ErrInvalidRequest             = errors.New("invalid request")
)

// Acceptable characters are the English alphabet, numbers,
// symbols from the symbols variable, the Ukrainian alphabet
func charsValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	regex := regexp.MustCompile(`^[a-zA-Z0-9 !@#$%^&*()_{}\[\]:;<>,.?~абвгґдеєжзиіїйклмнопрстуфхцчшщьюяАБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯ]*$`)
	return regex.MatchString(value)
}

func defaultValidator() *validator.Validate {
	validate := validator.New()

	if err := validate.RegisterValidation("charsValidation", charsValidation); err != nil {
		panic(err)
	}

	return validate
}

func validateParticipant(p *ParticipantRequest) error {
	phoneRegex := regexp.MustCompile(`^\+380\d{9,10}$`)

	validate := validator.New()
	if err := validate.RegisterValidation("charsValidation", charsValidation); err != nil {
		return err
	}

	if err := validate.Var(p.Name, "required,min=2,max=50,charsValidation"); err != nil {
		if strings.Contains(err.Error(), "min") {
			return ErrNameTooShort
		}
		return err
	}

	if err := validate.Var(p.Note, "charsValidation,lte=1000"); err != nil {
		return errors.New("note contains invalid characters")
	}

	if !phoneRegex.MatchString(p.Phone) {
		return fmt.Errorf("phone should be between %d and %d digits long: %w", minParticipantPhoneLength, maxParticipantPhoneLength, ErrParticipantPhoneOnlyDigits)
	}

	return nil
}
