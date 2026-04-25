package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{validate: validator.New()}
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		if validationErrors, exist := err.(validator.ValidationErrors); exist {
			return formatValidationErrors(validationErrors)
		}
		return err
	}
	return nil
}

func formatValidationErrors(errs validator.ValidationErrors) error {
	messages := ""
	for _, err := range errs {
		messages += fmt.Sprintf(
			"Field validation for '%s' failed on the '%s' tag.",
			err.Field(), err.Tag(),
		)
	}
	return fmt.Errorf(messages)
}
