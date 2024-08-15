package corretto

import "reflect"

const notABoolErrorMsg = "field %s is not a boolean"

type BoolValidator struct {
	*BaseValidator
}

// Bool checks if the field is a boolean
func (v *BaseValidator) Bool(msg ...string) *BoolValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if v.field.Kind() != reflect.Bool {
			return newValidationError(notABoolErrorMsg, cmsg, v.fieldName)
		}
		return nil
	})

	return &BoolValidator{v}
}

// Test allows you to run a custom validation function
//
// The function should have the signature:
//
//	func(ctx corretto.Context, value bool) error
func (v *BoolValidator) Test(f CustomValidationFunc[bool]) *BoolValidator {
	v.validations = append(v.validations, func() error {
		return f(v.ctx, v.field.Bool())
	})
	return v
}
