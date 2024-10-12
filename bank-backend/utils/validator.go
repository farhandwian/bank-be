package utils

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateIndonesianPhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	indonesianPhoneRegex := `^(\+62|62|0)[\s-]?8[1-9]{1}[0-9]{8,10}$`
	return regexp.MustCompile(indonesianPhoneRegex).MatchString(phoneNumber)
}

func FormatValidationErrors(err error) map[string]string {
	errorMessages := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			errorMessages[err.Field()] = "This field is required"
		case "email":
			errorMessages[err.Field()] = "Invalid email format"
		case "min":
			errorMessages[err.Field()] = fmt.Sprintf("Must be at least %s characters long", err.Param())
		case "max":
			errorMessages[err.Field()] = fmt.Sprintf("Must not be longer than %s characters", err.Param())
		case "alphanum":
			errorMessages[err.Field()] = "Must contain only alphanumeric characters"
		case "len":
			errorMessages[err.Field()] = fmt.Sprintf("Must be exactly %s characters long", err.Param())
		case "numeric":
			errorMessages[err.Field()] = "Must contain only numeric characters"
		case "strongpassword":
			errorMessages[err.Field()] = "Must contain at least one uppercase letter, one lowercase letter, one number, and one special character"
		case "indonesianphone":
			errorMessages[err.Field()] = "Must be a valid Indonesian phone number : start with +628121.."
		default:
			errorMessages[err.Field()] = fmt.Sprintf("Failed on %s validation", err.Tag())
		}
	}

	return errorMessages
}
