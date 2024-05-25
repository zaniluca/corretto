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

func (v *NumberValidator) Positive(msg ...string) *NumberValidator {
	cmsg := optional(msg)
	if cmsg == "" {
		cmsg = notAPositiveNumberMessage
	}

	return v // TODO
}
