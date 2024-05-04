package corretto

import (
	"reflect"
)

const (
	requiredErrorMsg = "%v is a required field"
	minErrorMsg      = "%v must be at least %v"
)

// Required checks if the field is not at its zero value
// It can be used with any type that has a zero value
//
// If the field is not supported, it will panic
func (v *BaseValidator) Required(msg ...string) *BaseValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		// Check if the value is at its zero value
		if v.field.IsZero() {
			return newValidationError(requiredErrorMsg, cmsg, v.fieldName)
		}
		return nil
	})
	return v
}

// Min checks if the field is greater than or equal to the provided value
// It can be used with int, float, string or slice
//
// If the field is not supported, it will panic
func (v *BaseValidator) Min(min int, msg ...string) *BaseValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		switch v.field.Type().Kind() {
		case reflect.Int:
			if v.field.Int() < int64(min) {
				return newValidationError(minErrorMsg, cmsg, v.fieldName, min)
			}
		case reflect.Float64:
			if v.field.Float() < float64(min) {
				return newValidationError(minErrorMsg, cmsg, v.fieldName, min)
			}
		case reflect.String:
			if len(v.field.String()) < min {
				return newValidationError(minErrorMsg+" characters long", cmsg, v.fieldName, min)
			}
		case reflect.Slice:
			if v.field.Len() < min {
				return newValidationError(minErrorMsg+" elements long", cmsg, v.fieldName, min)
			}
		default:
			logger.Panicf("unsopported type %v for Min(), can only be used with int, float or string", v.field.Type().Kind())
		}

		return nil
	})
	return v
}

// Schema checks if the field can be parsed by the provided schema
// Use it to validate nested structs
//
// NOTE: the field associated with the schema must be exported
//
//	type Parent struct {
//		Son  *Son // Field.Schema() works fine
//		daughter *Daughter // Field.Schema() will panic
//	}
func (v *BaseValidator) Schema(s Schema) *BaseValidator {
	v.validations = append(v.validations, func() error {
		// TODO: rewrite the panic message to be more informative
		return s.Parse(v.field.Interface())
	})
	return v
}
