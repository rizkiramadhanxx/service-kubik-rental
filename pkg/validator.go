package pkg

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// TYPE VALIDATOR
var Validate = validator.New()

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			var msg string

			switch e.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", field)
			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
			case "max":
				msg = fmt.Sprintf("%s must be at most %s characters", field, e.Param())
			case "alphanum":
				msg = fmt.Sprintf("%s must only contain letters and numbers", field)
			default:
				msg = fmt.Sprintf("%s is not valid", field)
			}

			errors[field] = msg
		}
	}

	return errors
}
