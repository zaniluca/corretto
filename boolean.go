package corretto

import "reflect"

const notABooleanErrorMsg = "field %s is not a boolean"

type BooleanValidator struct {
	*BaseValidator
}

// Boolean checks if the field is a boolean
func (v *BaseValidator) Boolean(msg ...string) *BooleanValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if v.field.Kind() != reflect.Bool {
			return newValidationError(notABooleanErrorMsg, cmsg, v.fieldName)
		}
		return nil
	})

	return &BooleanValidator{v}
}

// Test allows you to run a custom validation function
//
// The function should have the signature:
//
//	func(ctx corretto.Context, value bool) error
func (v *BooleanValidator) Test(f CustomValidationFunc[bool]) *BooleanValidator {
	v.validations = append(v.validations, func() error {
		return f(v.ctx, v.field.Bool())
	})
	return v
}
