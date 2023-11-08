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
	errParticipantPhoneOnlyDigits = errors.New("phone should contain only digits")
	errNameTooShort               = errors.New("name is too short")
)

// Acceptable characters are the English alphabet, numbers,
// symbols from the symbols variable, the Ukrainian alphabet
func charsValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	regex := regexp.MustCompile("^[a-zA-Z0-9 !@#$%^&*()_{}\\[\\]:;<>,.?~абвгґдеєжзиіїйклмнопрстуфхцчшщьюяАБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯ]*$")
	return regex.MatchString(value)
}

func validateStruct(s interface{}) error {
	validate := validator.New()

	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return errors.New(fmt.Sprintf("Validation error in field %s: %s", validationErrors[0].Field(), validationErrors[0].Tag()))
		}
		return err
	}

	return nil
}

func validateRaffle(raf *RaffleRequest) error {
	if err := validateStruct(raf); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.RegisterValidation("charsValidation", charsValidation); err != nil {
		return err
	}

	if err := validate.Var(raf.Name, "charsValidation"); err != nil {
		return errors.New("name contains invalid characters")
	}

	if err := validate.Var(raf.Note, "charsValidation"); err != nil {
		return errors.New("note contains invalid characters")
	}

	return nil
}

func validatePrize(p *PrizeRequest) error {
	if err := validateStruct(p); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.RegisterValidation("charsValidation", charsValidation); err != nil {
		return err
	}

	if err := validate.Var(p.Name, "charsValidation"); err != nil {
		return errors.New("name contains invalid characters")
	}

	if err := validate.Var(p.Description, "charsValidation"); err != nil {
		return errors.New("note contains invalid characters")
	}

	return nil
}

func validateParticipant(p *ParticipantRequest) error {
	phoneRegex := regexp.MustCompile(`^\+380\d{9,10}$`)

	validate := validator.New()
	if err := validate.RegisterValidation("charsValidation", charsValidation); err != nil {
		return err
	}

	if err := validate.Var(p.Name, "required,min=2,max=50,charsValidation"); err != nil {
		if strings.Contains(err.Error(), "min") {
			return errNameTooShort
		}
		return err
	}

	if err := validate.Var(p.Note, "charsValidation,lte=1000"); err != nil {
		return errors.New("note contains invalid characters")
	}

	if !phoneRegex.MatchString(p.Phone) {
		return fmt.Errorf("phone should be between %d and %d digits long: %w", minParticipantPhoneLength, maxParticipantPhoneLength, errParticipantPhoneOnlyDigits)
	}

	return nil
}
