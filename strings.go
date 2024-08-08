package corretto

import (
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

const (
	notAStringErrorMsg      = "%v is not a string"
	mustIncludeErrorMsg     = "%v must include %v"
	stringMinLengthErrorMsg = "%v must be at least %v characters long"
	stringMaxLengthErrorMsg = "%v must be at most %v characters long"
	stringLengthErrorMsg    = "%v must be %v characters long"
	matchesErrorMsg         = "%v is not in the correct format"
	nonEmptyErrorMsg        = "%v cannot be empty"
	mustStartWithErrorMsg   = "%v must start with %v"
	mustEndWithErrorMsg     = "%v must end with %v"
	notAValidURLErrorMsg    = "%v is not a valid URL"
)

const (
	emailRegexString    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	uuidRegexString     = `^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`
	cuidRegexString     = `^c[^\s-]{8,}$`
	hexColorRegexString = `#[a-f\d]{3}(?:[a-f\d]?|(?:[a-f\d]{3}(?:[a-f\d]{2})?)?)\b`
)

var (
	emailRegex    = regexp.MustCompile(emailRegexString)
	uuidRegex     = regexp.MustCompile(uuidRegexString)
	cuidRegex     = regexp.MustCompile(cuidRegexString)
	hexColorRegex = regexp.MustCompile(hexColorRegexString)
)

type StringValidator struct {
	*BaseValidator
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

// NonEmpty checks if the field does not contain an empty string, it trims the string before checking
func (v *StringValidator) NonEmpty(msg ...string) *StringValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if strings.TrimSpace(v.field.String()) == "" {
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

// MaxLength checks if the field has a length less than or equal to the provided value
func (v *StringValidator) MaxLength(max int, msg ...string) *StringValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if len(v.field.String()) > max {
			return newValidationError(stringMaxLengthErrorMsg, cmsg, v.fieldName, max)
		}
		return nil
	})
	return v
}

// Length checks if the field has a length equal to the provided value
//
// if you want to check for a range of values, use [StringValidator.MinLength] and [StringValidator.MaxLength]
func (v *StringValidator) Length(l int, msg ...string) *StringValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if len(v.field.String()) != l {
			return newValidationError(stringLengthErrorMsg, cmsg, v.fieldName, l)
		}
		return nil
	})

	return v
}

// Test allows you to run a custom validation function
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

// Includes checks if the field value contains the provided substring
//
// NOTE: it is case sensitive
//
// If you want to check only for prefix or suffix, use [StringValidator.StartsWith] or [StringValidator.EndsWith]
func (v *StringValidator) Includes(substr string, msg ...string) *StringValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if !strings.Contains(v.field.String(), substr) {
			return newValidationError(mustIncludeErrorMsg, cmsg, v.fieldName, substr)
		}
		return nil
	})

	return v
}

// StartsWith checks if the field value starts with the provided prefix
//
// NOTE: it is case sensitive
func (v *StringValidator) StartsWith(prefix string, msg ...string) *StringValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if !strings.HasPrefix(v.field.String(), prefix) {
			return newValidationError(mustStartWithErrorMsg, cmsg, v.fieldName, prefix)
		}
		return nil
	})

	return v
}

// EndsWith checks if the field value ends with the provided suffix
//
// NOTE: it is case sensitive
func (v *StringValidator) EndsWith(suffix string, msg ...string) *StringValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if !strings.HasSuffix(v.field.String(), suffix) {
			return newValidationError(mustEndWithErrorMsg, cmsg, v.fieldName, suffix)
		}
		return nil
	})

	return v
}

// Url checks if the field is a valid URL format
func (v *StringValidator) Url(msg ...string) *StringValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		_, err := url.ParseRequestURI(v.field.String())
		if err != nil {
			return newValidationError(notAValidURLErrorMsg, cmsg, v.fieldName)
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

// Uuid checks if the field is a valid UUID v4 format
//
// if the string is empty, it will not return error, use [StringValidator.NonEmpty] to check for empty strings
func (v *StringValidator) Uuid(msg ...string) *StringValidator {
	return v.Matches(uuidRegex.String(), msg...)
}

// Cuid checks if the field is a valid CUID format (Collision-resistant ids)
// See: https://github.com/paralleldrive/cuid
//
// if the string is empty, it will not return error, use [StringValidator.NonEmpty] to check for empty strings
func (v *StringValidator) Cuid(msg ...string) *StringValidator {
	return v.Matches(cuidRegex.String(), msg...)
}

// HexColor checks if the field is a valid HEX color format
// It supports both 3 and 6 characters long hex color and alpha values
// Example: #fff, #ffffff, #fff0, #ffffff00
//
// NOTE: it is case insensitive and it accepts both #
//
// if the string is empty, it will not return error, use [StringValidator.NonEmpty] to check for empty strings
func (v *StringValidator) HexColor(msg ...string) *StringValidator {
	return v.Matches(hexColorRegex.String(), msg...)
}
