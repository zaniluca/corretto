package corretto

import (
	"reflect"
	"regexp"
)

const (
	notAStringErrorMsg = "%v is not a string"
	matchesErrorMsg    = "%v is not in the correct format"
)

const (
	emailRegexString = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

var (
	emailRegex = regexp.MustCompile(emailRegexString)
)

type stringValidator struct {
	*baseValidator
}

// Matches checks if the field matches the provided regex pattern
// It can only be used with string
//
// if the string is empty, it will return true, use [baseValidator.Required] to check for empty strings
//
// If the field is not a string, it will panic
func (v *stringValidator) Matches(regex string, msg ...string) *stringValidator {
	cmsg := optional(msg)
	r := regexp.MustCompile(regex)

	v.validations = append(v.validations, func() error {
		if v.field.Kind() != reflect.String {
			logger.Panic("Matches() can only be used with strings")
		}
		if !r.MatchString(v.field.String()) && v.field.String() != "" {
			return newValidationError(matchesErrorMsg, cmsg, v.fieldName)
		}
		return nil
	})
	return v
}

// Email checks if the field is a valid email address format
// It can only be used with string
//
// if the string is empty, it will return true, use [baseValidator.Required] to check for empty strings
//
// If the field is not a string, it will panic
func (v *stringValidator) Email(msg ...string) *stringValidator {
	return v.Matches(emailRegex.String(), msg...)
}

func (v *baseValidator) String(msg ...string) *stringValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if v.field.Kind() != reflect.String {
			return newValidationError(notAStringErrorMsg, cmsg, v.fieldName)
		}
		return nil
	})

	return &stringValidator{v}
}
