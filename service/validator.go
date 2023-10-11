package service

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"strconv"
	"strings"
)

const (
	minParticipantPhoneLength = 10
	maxParticipantPhoneLength = 12
	allowedChars              = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_{}[]:;<>,.?~абвгґдеєжзиіїйклмнопрстуфхцчшщьюяАБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯ"
)

var (
	errParticipantPhoneOnlyDigits = errors.New("phone should contain only digits")
)

func customValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	for _, char := range value {
		if !strings.ContainsRune(allowedChars, char) {
			return false
		}
	}

	return true
}

func validateStruct(s interface{}) error {
	validate := validator.New()
	validate.RegisterValidation("allowedChars", customValidation)

	if err := validate.Struct(s); err != nil {
		var validationErrs []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrs = append(validationErrs, fmt.Sprintf("Validation error in field %s: %s", err.Field(), err.Tag()))
		}
		return errors.New(strings.Join(validationErrs, "\n"))
	}

	return nil
}

func validateRaffle(raf *RaffleRequest) error {
	return validateStruct(raf)
}

func validatePrize(p *PrizeRequest) error {
	return validateStruct(p)
}

func validateParticipant(p *ParticipantRequest) error {
	if err := validateStruct(p); err != nil {
		return err
	}

	p.Phone = strings.ReplaceAll(p.Phone, " ", "")
	p.Phone = strings.ReplaceAll(p.Phone, "+", "")
	runePhone := []rune(p.Phone)

	_, err := strconv.Atoi(p.Phone)
	if err != nil {
		return errParticipantPhoneOnlyDigits
	}

	if len(runePhone) < minParticipantPhoneLength || len(runePhone) > maxParticipantPhoneLength {
		return fmt.Errorf("phone should be between %d and %d digits long", minParticipantPhoneLength, maxParticipantPhoneLength)
	}

	return nil
}
