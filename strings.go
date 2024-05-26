package corretto

import (
	"reflect"
	"regexp"
)

const (
	notAStringErrorMsg = "%v is not a string"
	minLengthErrorMsg  = "%v must be at least %v characters long"
	matchesErrorMsg    = "%v is not in the correct format"
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

// Min checks if the field has a length greater than or equal to the provided value
func (v *StringValidator) Min(min int, msg ...string) *StringValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if len(v.field.String()) < min {
			return newValidationError(minLengthErrorMsg, cmsg, v.fieldName, min)
		}
		return nil
	})
	return v
}

// Matches checks if the field matches the provided regex pattern
//
// if the string is empty, it will not return error, use [BaseValidator.Required] to check for empty strings
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
// if the string is empty, it will not return error, use [BaseValidator.Required] to check for empty strings
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
