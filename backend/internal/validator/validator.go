package validator

import (
	"fmt"
	"reflect"
	"strings"

	"enterprise-order-management-api/backend/internal/util"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func New() *CustomValidator {
	validate := validator.New()
	validate.RegisterTagNameFunc(jsonFieldName)

	return &CustomValidator{
		validator: validate,
	}
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return util.BadRequest("Validation failed", map[string]string{
				"request": "Invalid request body",
			})
		}

		return util.BadRequest("Validation failed", buildValidationMessages(validationErrors))
	}

	return nil
}

func buildValidationMessages(validationErrors validator.ValidationErrors) map[string]string {
	messages := make(map[string]string, len(validationErrors))

	for _, fieldError := range validationErrors {
		fieldName := fieldError.Field()
		messages[fieldName] = validationMessage(fieldError)
	}

	return messages
}

func validationMessage(fieldError validator.FieldError) string {
	fieldName := fieldError.Field()

	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldName)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fieldName)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters or items", fieldName, fieldError.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fieldName, fieldError.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fieldName, fieldError.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", fieldName, fieldError.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fieldName)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fieldName, fieldError.Param())
	default:
		return fmt.Sprintf("%s is invalid", fieldName)
	}
}

func jsonFieldName(field reflect.StructField) string {
	name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	if name == "" {
		return field.Name
	}
	return name
}
