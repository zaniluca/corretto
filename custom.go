package corretto

// Test is a custom validation function that can be used to add custom validation
func (v *BaseValidator) Test(f CustomValidationFunc) *BaseValidator {
	v.validations = append(v.validations, func() error {
		return f(v.ctx, v.field)
	})
	return v
}
