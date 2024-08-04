package corretto

import (
	"reflect"
	"slices"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

const (
	notANumberMsg         = "%v is not a number"
	notAPositiveNumberMsg = "%v must be a positive number"
	notANegativeNumberMsg = "%v must be a negative number"
	zeroNumberErrorMsg    = "%v is required"
	minNumberErrorMsg     = "%v must be at least %v"
)

type NumberValidator struct {
	*BaseValidator
}

func (v *BaseValidator) Number(msg ...string) *NumberValidator {
	cmsg := optional(msg)
	numbers := []reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64}

	v.validations = append(v.validations, func() error {
		if !slices.Contains(numbers, v.field.Kind()) {
			return newValidationError(notANumberMsg, cmsg, v.fieldName)
		}
		return nil
	})

	return &NumberValidator{v}
}

// NonZero checks if the field is not "0"
func (v *NumberValidator) NonZero(msg ...string) *NumberValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if v.field.IsZero() {
			return newValidationError(zeroNumberErrorMsg, cmsg, v.fieldName)
		}
		return nil
	})
	return v
}

// Positive checks if the field is a positive number (greater than or equal to zero)
func (v *NumberValidator) Positive(msg ...string) *NumberValidator {
	cmsg := optional(msg)
	if cmsg == "" {
		cmsg = notAPositiveNumberMsg
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
				return newValidationError(minNumberErrorMsg, cmsg, v.fieldName, min)
			}
		case reflect.Float64:
			if v.field.Float() < float64(min) {
				return newValidationError(minNumberErrorMsg, cmsg, v.fieldName, min)
			}
		default:
			logger.Panicf("unsupported type %v for Min(), can only be used with int or float", v.field.Kind())
		}

		return nil
	})

	return v
}

// Test is a custom validation function that can be used to add custom validation
func (v *NumberValidator) Test(f CustomValidationFunc[int64]) *NumberValidator {
	v.validations = append(v.validations, func() error {
		return f(v.ctx, v.field.Int())
	})
	return v
}

// OneOf checks if the field value contains one of the provided values
func (v *NumberValidator) OneOf(allowed []int, msg ...string) *NumberValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if !oneOf(int(v.field.Int()), allowed) {
			return newValidationError(oneOfErrorMsg, cmsg, v.fieldName, allowed)
		}
		return nil
	})

	return v
}
