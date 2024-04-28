package corretto

import (
	"reflect"
	"regexp"
)

const (
	matchesErrorMsg = "%v is not in the correct format"
)

const (
	emailRegexString = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

var (
	emailRegex = regexp.MustCompile(emailRegexString)
)

// Matches checks if the field matches the provided regex pattern
// It can only be used with string
//
// if the string is empty, it will return true, use [Validator.Required] to check for empty strings
//
// If the field is not a string, it will panic
func (v *Validator) Matches(regex string, msg ...string) *Validator {
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
// if the string is empty, it will return true, use [Validator.Required] to check for empty strings
//
// If the field is not a string, it will panic
func (v *Validator) Email(msg ...string) *Validator {
	return v.Matches(emailRegex.String(), msg...)
}
