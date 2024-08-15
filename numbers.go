package corretto

import (
	"math"
	"reflect"
	"slices"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

const (
	notANumberMsg            = "%v is not a number"
	notAPositiveNumberMsg    = "%v must be a positive number"
	notANegativeNumberMsg    = "%v must be a negative number"
	notANonNegativeNumberMsg = "%v must be a non-negative number"
	notANonPositiveNumberMsg = "%v must be a non-positive number"
	notAMultipleOfMsg        = "%v must be a multiple of %v"
	notAFiniteNumberMsg      = "%v must be a finite number"
	zeroNumberErrorMsg       = "%v is required"
	minNumberErrorMsg        = "%v must be at least %v"
	maxNumberErrorMsg        = "%v must be less than %v"
)

type NumberValidator struct {
	*BaseValidator
}

// Number checks if the field is a number, either an integer or a float
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

// Positive checks if the field is a positive number (> 0)
//
// To check if the field is a non-negative number (>= 0), use [NumberValidator.NonNegative]
func (v *NumberValidator) Positive(msg ...string) *NumberValidator {
	cmsg := optional(msg)
	if cmsg == "" {
		cmsg = notAPositiveNumberMsg
	}

	return v.Min(1, cmsg)
}

// Negative checks if the field is a negative number (< 0)
//
// To check if the field is a non-positive number (<= 0), use [NumberValidator.NonPositive]
func (v *NumberValidator) Negative(msg ...string) *NumberValidator {
	cmsg := optional(msg)
	if cmsg == "" {
		cmsg = notANegativeNumberMsg
	}

	return v.Max(-1, cmsg)
}

// NonNegative checks if the field is a non-negative number (>= 0)
func (v *NumberValidator) NonNegative(msg ...string) *NumberValidator {
	cmsg := optional(msg)
	if cmsg == "" {
		cmsg = notANonNegativeNumberMsg
	}

	return v.Min(0, cmsg)
}

// NonPositive checks if the field is a non-positive number (<= 0)
func (v *NumberValidator) NonPositive(msg ...string) *NumberValidator {
	cmsg := optional(msg)
	if cmsg == "" {
		cmsg = notANonPositiveNumberMsg
	}

	return v.Max(0, cmsg)
}

// Min checks if the field is greater than or equal to the provided value
func (v *NumberValidator) Min(min int, msg ...string) *NumberValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		switch v.field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.field.Int() < int64(min) {
				return newValidationError(minNumberErrorMsg, cmsg, v.fieldName, min)
			}
		case reflect.Float64, reflect.Float32:
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

// Max checks if the field is less than or equal to the provided value
func (v *NumberValidator) Max(max int, msg ...string) *NumberValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		switch v.field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.field.Int() > int64(max) {
				return newValidationError(maxNumberErrorMsg, cmsg, v.fieldName, max)
			}
		case reflect.Float64, reflect.Float32:
			if v.field.Float() > float64(max) {
				return newValidationError(maxNumberErrorMsg, cmsg, v.fieldName, max)
			}
		default:
			logger.Panicf("unsupported type %v for Max(), can only be used with int or float", v.field.Kind())
		}

		return nil
	})

	return v
}

// Test allows you to run a custom validation function
//
// The function should have the signature:
//
//	func(ctx corretto.Context) error
func (v *NumberValidator) Test(f CustomValidationFunc) *NumberValidator {
	v.validations = append(v.validations, func() error {
		return f(v.ctx)
	})
	return v
}

// OneOf checks if the field value contains one of the provided values
//
// NOTE: This validation can only be used with Integers. If the field is a float, it will be converted to an int before being checked
func (v *NumberValidator) OneOf(allowed []int, msg ...string) *NumberValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		var val int
		switch v.field.Kind() {
		case reflect.Float64, reflect.Float32:
			val = int(v.field.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val = int(v.field.Int())
		default:
			logger.Panicf("unsupported type %v for OneOf(), can only be used with int or float", v.field.Kind())
		}

		if !oneOf(val, allowed) {
			return newValidationError(oneOfErrorMsg, cmsg, v.fieldName, allowed)
		}
		return nil
	})

	return v
}

// MultipleOf checks if the field value is a multiple of the provided divisor
//
// NOTE: By definition, 0 is a multiple of any number, so if the field value is 0, this validation will always pass
// NOTE: This validation can only be used with Integers. If the field is a float, it will be converted to an int before being checked
func (v *NumberValidator) MultipleOf(divisor int, msg ...string) *NumberValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		var val int
		switch v.field.Kind() {
		case reflect.Float64, reflect.Float32:
			val = int(v.field.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val = int(v.field.Int())
		default:
			logger.Panicf("unsupported type %v for MultipleOf(), can only be used with int or float", v.field.Kind())
		}
		if val%divisor != 0 {
			return newValidationError(notAMultipleOfMsg, cmsg, v.fieldName, divisor)
		}
		return nil
	})

	return v
}

// Finite checks if the field value is a finite number, i.e., not infinite
func (v *NumberValidator) Finite(msg ...string) *NumberValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		switch v.field.Kind() {
		case reflect.Float64, reflect.Float32:
			if math.IsInf(v.field.Float(), 0) {
				return newValidationError(notAFiniteNumberMsg, cmsg, v.fieldName)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if math.IsInf(float64(v.field.Int()), 0) {
				return newValidationError(notAFiniteNumberMsg, cmsg, v.fieldName)
			}
		default:
			logger.Panicf("unsupported type %v for Finite(), can only be used with int or float", v.field.Kind())
		}
		return nil
	})

	return v
}
