package corretto

import (
	"fmt"
	"reflect"
)

const (
	notAnArrayErrorMsg = "%v is not an array"
)

const (
	arrayElementFieldName = "%v's elements"
)

type ArrayValidator struct {
	*BaseValidator
}

// Array checks if the field is an array (slice)
//
// It doesn't check if the array is empty, use [BaseValidator.Required] to check for empty arrays
func (v *BaseValidator) Array(msg ...string) *ArrayValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if v.field.Kind() != reflect.Slice {
			return newValidationError(notAnArrayErrorMsg, cmsg, v.fieldName)
		}
		return nil
	})

	return &ArrayValidator{v}
}

// Of checks if all elements of the array are valid according to the provided validator
//
//	"Hobbies": corretto.Field().Array().Of(corretto.Field("Hobby").Min(3)),
//
// You can also use [BaseValidator.Schema] to validate the array elements with a schema
//
//	 s := corretto.Schema{
//		"Name": corretto.Field("Name").Required(),
//		"Age": corretto.Field("Age").Min(18),
//	 }
//
//	"Users": corretto.Field().Array().Of(corretto.Field("User").Schema(s)),
func (v *ArrayValidator) Of(validator validator) *ArrayValidator {
	v.validations = append(v.validations, func() error {
		for i := 0; i < v.field.Len(); i++ {
			bv := validator.getBaseValidator()
			bv.field = v.field.Index(i)
			// If no custom field name is provided, use the struct field name formatted accordingly
			if bv.fieldName == "" {
				bv.fieldName = fmt.Sprintf(arrayElementFieldName, v.fieldName)
			}
			bv.ctx = v.ctx

			// If any of the elements fail the validation, return the error
			if err := bv.check(); err != nil {
				return err
			}
		}
		return nil
	})

	return v
}

func (v *ArrayValidator) Min(min int, msg ...string) *ArrayValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if v.field.Len() < min {
			return newValidationError(minErrorMsg+" elements long", cmsg, v.fieldName, min)
		}
		return nil
	})

	return v
}
