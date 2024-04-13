package corretto

// Test is a custom validation function that can be used to add custom validation
func (v *Validator) Test(f ValidationFunc) *Validator {
	v.validations = append(v.validations, f)
	return v
}
