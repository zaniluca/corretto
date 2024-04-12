package corretto

import (
	"fmt"
	"reflect"
	"regexp"
)

const (
	emailRegexString = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

var (
	emailRegex = regexp.MustCompile(emailRegexString)
)

func (v *Validator) Matches(regex string) *Validator {
	r := regexp.MustCompile(regex)

	v.validations = append(v.validations, func() error {
		if v.field.Kind() != reflect.String {
			logger.Panic("Matches() can only be used with strings")
		}
		if !r.MatchString(v.field.String()) {
			return fmt.Errorf("%s is not in the correct format", v.field.String())
		}
		return nil
	})
	return v
}

func (v *Validator) Email() *Validator {
	return v.Matches(emailRegex.String())
}
