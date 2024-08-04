package corretto

import (
	"reflect"
	"regexp"
)

const (
	notAStringErrorMsg      = "%v is not a string"
	stringMinLengthErrorMsg = "%v must be at least %v characters long"
	matchesErrorMsg         = "%v is not in the correct format"
	nonEmptyErrorMsg        = "%v cannot be empty"
)

const (
	emailRegexString = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

var (
	emailRegex = regexp.MustCompile(emailRegexString)
)

type StringValidator struct {
	*BaseValidator
}

// NonEmpty checks if the field does not contain an empty string
func (v *StringValidator) NonEmpty(msg ...string) *StringValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if v.field.IsZero() {
			return newValidationError(nonEmptyErrorMsg, cmsg, v.fieldName)
		}
		return nil
	})
	return v
}

// MinLength checks if the field has a length greater than or equal to the provided value
func (v *StringValidator) MinLength(min int, msg ...string) *StringValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if len(v.field.String()) < min {
			return newValidationError(stringMinLengthErrorMsg, cmsg, v.fieldName, min)
		}
		return nil
	})
	return v
}

// Test is a custom validation function that can be used to add custom validation
func (v *StringValidator) Test(f CustomValidationFunc[string]) *StringValidator {
	v.validations = append(v.validations, func() error {
		return f(v.ctx, v.field.String())
	})
	return v
}

// Matches checks if the field matches the provided regex pattern
//
// if the string is empty, it will not return error, use [StringValidator.NonEmpty] to check for empty strings
//
// it uses the [regexp] package to match the regex, if the regex is invalid, it will panic
func (v *StringValidator) Matches(regex string, msg ...string) *StringValidator {
	cmsg := optional(msg)
	r := regexp.MustCompile(regex)

	v.validations = append(v.validations, func() error {
		if v.field.String() != "" && !r.MatchString(v.field.String()) {
			return newValidationError(matchesErrorMsg, cmsg, v.fieldName)
		}
		return nil
	})
	return v
}

// Email checks if the field is a valid email address format
//
// if the string is empty, it will not return error, use [StringValidator.NonEmpty] to check for empty strings
func (v *StringValidator) Email(msg ...string) *StringValidator {
	return v.Matches(emailRegex.String(), msg...)
}

func (v *BaseValidator) String(msg ...string) *StringValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if v.field.Kind() != reflect.String {
			return newValidationError(notAStringErrorMsg, cmsg, v.fieldName)
		}
		return nil
	})

	return &StringValidator{v}
}

// OneOf checks if the field value contains one of the provided values
//
// NOTE: it is case sensitive
func (v *StringValidator) OneOf(allowed []string, msg ...string) *StringValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if !oneOf(v.field.String(), allowed) {
			return newValidationError(oneOfErrorMsg, cmsg, v.fieldName, allowed)
		}
		return nil
	})

	return v
}
