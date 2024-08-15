package corretto

import (
	"fmt"
	"reflect"
)

const (
	notAnArrayErrorMsg     = "%v is not an array"
	arrayMinLengthErrorMsg = "%v must be at least %v elements long"
	emptyArrayErrorMsg     = "%v cannot be empty"
)

const (
	arrayElementFieldName = "%v's elements"
)

type ArrayValidator struct {
	*BaseValidator
}

// NonEmpty checks if the field does not contain an empty array
func (v *ArrayValidator) NonEmpty(msg ...string) *ArrayValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if v.field.Len() == 0 {
			return newValidationError(emptyArrayErrorMsg, cmsg, v.fieldName)
		}
		return nil
	})
	return v
}

// MinLength checks if the array has a length greater than or equal to the provided value
func (v *ArrayValidator) MinLength(min int, msg ...string) *ArrayValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if v.field.Len() < min {
			return newValidationError(arrayMinLengthErrorMsg, cmsg, v.fieldName, min)
		}
		return nil
	})

	return v
}

// MaxLength checks if the array has a length less than or equal to the provided value
func (v *ArrayValidator) MaxLength(max int, msg ...string) *ArrayValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if v.field.Len() > max {
			return newValidationError(arrayMinLengthErrorMsg, cmsg, v.fieldName, max)
		}
		return nil
	})

	return v
}

// Length checks if the array has a length equal to the provided value
func (v *ArrayValidator) Length(length int, msg ...string) *ArrayValidator {
	cmsg := optional(msg)

	v.validations = append(v.validations, func() error {
		if v.field.Len() != length {
			return newValidationError(arrayMinLengthErrorMsg, cmsg, v.fieldName, length)
		}
		return nil
	})

	return v
}

// Test allows you to run a custom validation function on the array
//
// The function should have the signature:
//
//	func (ctx corretto.Context) error
func (v *ArrayValidator) Test(f CustomValidationFunc) *ArrayValidator {
	v.validations = append(v.validations, func() error {
		return f(v.ctx)
	})
	return v
}

// Array checks if the field is an array (slice)
//
// It doesn't check if the array is empty, use [ArrayValidator.NonEmpty] to check for empty arrays
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
//	"Hobbies": corretto.Field().Array().Of(corretto.Field("Hobby").String().MinLength(3)),
//
// You can also use [BaseValidator.Schema] to validate the array elements with a schema
//
//	 s := corretto.Schema{
//		"Name": corretto.Field("Name").String().NonEmpty(),
//		"Age": corretto.Field("Age").Number().Min(18),
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
