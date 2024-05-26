package corretto

import (
	"reflect"
	"slices"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

const (
	notANumberMessage         = "%v is not a number"
	notAPositiveNumberMessage = "%v must be a positive number"
	notANegativeNumberMessage = "%v must be a negative number"
)

type NumberValidator struct {
	*BaseValidator
}

func (v *BaseValidator) Number(msg ...string) *NumberValidator {
	cmsg := optional(msg)
	numbers := []reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64}

	v.validations = append(v.validations, func() error {
		if !slices.Contains(numbers, v.field.Kind()) {
			return newValidationError(notANumberMessage, cmsg, v.fieldName)
		}
		return nil
	})

	return &NumberValidator{v}
}

// Positive checks if the field is a positive number (greater than or equal to zero)
func (v *NumberValidator) Positive(msg ...string) *NumberValidator {
	cmsg := optional(msg)
	if cmsg == "" {
		cmsg = notAPositiveNumberMessage
	}

	return v.Min(0, cmsg)
}

// Min checks if the field is greater than or equal to the provided value
func (v *NumberValidator) Min(min int, msg ...string) *NumberValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		switch v.field.Kind() {
		case reflect.Int:
			if v.field.Int() < int64(min) {
				return newValidationError(minErrorMsg, cmsg, v.fieldName, min)
			}
		case reflect.Float64:
			if v.field.Float() < float64(min) {
				return newValidationError(minErrorMsg, cmsg, v.fieldName, min)
			}
		default:
			logger.Panicf("unsupported type %v for Min(), can only be used with int or float", v.field.Kind())
		}

		return nil
	})

	return v
}
