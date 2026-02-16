package utils

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// ValidateStruct validates a struct and returns errors
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		details := make([]map[string]string, 0)
		
		for _, e := range validationErrors {
			details = append(details, map[string]string{
				"field":   e.Field(),
				"message": getValidationMessage(e),
			})
		}
		
		return fmt.Errorf("validation failed: %d errors", len(details))
	}
	
	return nil
}

// getValidationMessage returns a user-friendly validation message
func getValidationMessage(e validator.FieldError) string {
	field := e.Field()
	
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "len":
		return fmt.Sprintf("%s must be %s characters long", field, e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of %s", field, e.Param())
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// GetFieldTag returns a tag value from a struct field
func GetFieldTag(s interface{}, field, tag string) string {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		fieldInfo := typ.Field(i)
		if fieldInfo.Name == field {
			return fieldInfo.Tag.Get(tag)
		}
	}
	
	return ""
}
