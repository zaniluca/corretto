package corretto

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
		if !v.field.CanInterface() {
			logger.Panicf("field `%v` must be exported to be validated", v.key)
		}
		return s.Parse(v.field.Interface())
	})
	return v
}

// Test is a custom validation function that can be used to add custom validation
func (v *BaseValidator) Test(f CustomValidationFunc) *BaseValidator {
	v.validations = append(v.validations, func() error {
		return f(v.ctx, v.field)
	})
	return v
}

func (v *BaseValidator) OneOf(allowed []any, msg ...string) *BaseValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if !oneOf(v.field.Interface(), allowed) {
			return newValidationError("%v is not one of %v", cmsg, v.fieldName, allowed)
		}
		return nil
	})

	return v
}

func oneOf[T comparable](v T, allowed []T) bool {
	for _, a := range allowed {
		if v == a {
			return true
		}
	}
	return false
}
