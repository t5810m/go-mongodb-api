package helpers

import (
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

var validate = validator.New()

// ValidateStruct validates a struct and returns structured error messages
func ValidateStruct(data interface{}) []ValidationError {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	var errors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   err.Field(),
			Message: getErrorMessage(err),
		})
	}
	return errors
}

// getErrorMessage returns a user-friendly error message based on the validation tag
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short (minimum: " + err.Param() + ")"
	case "max":
		return "Value is too long (maximum: " + err.Param() + ")"
	case "gt":
		return "Value must be greater than " + err.Param()
	case "gte":
		return "Value must be greater than or equal to " + err.Param()
	case "oneof":
		return "Invalid value. Allowed values: " + err.Param()
	default:
		return "Validation failed"
	}
}
