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
