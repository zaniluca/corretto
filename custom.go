package corretto

// Test is a custom validation function that can be used to add custom validation
func (v *Validator) Test(f CustomValidationFunc) *Validator {
	v.validations = append(v.validations, func() error {
		return f(v.ctx, v.field)
	})
	return v
}
