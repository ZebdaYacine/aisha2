package httpapi

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type RequestValidator struct{ validate *validator.Validate }

func NewRequestValidator() *RequestValidator {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &RequestValidator{validate: validate}
}

func (v *RequestValidator) Validate(value any) error {
	if err := v.validate.Struct(value); err != nil {
		fields := make(FieldErrors)
		var validationErrors validator.ValidationErrors
		if errors, ok := err.(validator.ValidationErrors); ok {
			validationErrors = errors
		} else {
			return WrapAPIError(err, CodeInternalError, "An unexpected error occurred.")
		}
		for _, fieldError := range validationErrors {
			fields[fieldError.Field()] = validationCode(fieldError.Tag(), fieldError.Param())
		}
		return NewAPIError(CodeValidationError, "The request is invalid.", fields)
	}
	return nil
}

func validationCode(tag, parameter string) string {
	switch tag {
	case "required":
		return "REQUIRED"
	case "email":
		return "INVALID_EMAIL"
	case "min":
		return "MIN_LENGTH_" + parameter
	case "max":
		return "MAX_LENGTH_" + parameter
	default:
		return "INVALID"
	}
}
