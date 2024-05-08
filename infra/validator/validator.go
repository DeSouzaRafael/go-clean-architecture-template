package validator

import (
	"fmt"
	"log"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := &Validator{validate: validator.New()}
	v.registerCustomValidations()
	return v
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

func (v *Validator) registerCustomValidations() {
	if err := v.validate.RegisterValidation("phone_format", validatePhoneFormat); err != nil {
		log.Printf("failed to register custom validation: %v", err)
	}
}

// Custom validation for Brazilian phone numbers, with placeholder
func validatePhoneFormat(fl validator.FieldLevel) bool {
	phoneRegex := `^\+55\d{10,11}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(fl.Field().String())
}

func formatValidationErrors(errs validator.ValidationErrors) error {
	messages := ""
	for _, err := range errs {
		switch err.Tag() {
		case "phone_format":
			messages += fmt.Sprintf(
				"Invalid phone number '%s'. Expected format: '+55XXXXXXXXXX' or '+55XXXXXXXXXXX'.",
				err.Value(),
			)
		default:
			messages += fmt.Sprintf(
				"Field validation for '%s' failed on the '%s' tag.",
				err.Field(), err.Tag(),
			)
		}
	}
	return fmt.Errorf(messages)
}
