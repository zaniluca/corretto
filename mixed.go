package corretto

const (
	oneOfErrorMsg = "%v must be one of %v"
)

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

func (v *BaseValidator) OneOf(allowed []any, msg ...string) *BaseValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if !oneOf(v.field.Interface(), allowed) {
			return newValidationError(oneOfErrorMsg, cmsg, v.fieldName, allowed)
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
