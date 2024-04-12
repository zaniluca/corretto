package corretto

import (
	"log"
	"reflect"
)

type Schema map[string]*Validator

// Interface for the validator
type validable interface {
	ValidationSchema() Schema
}

// Represents a validator for a field
type Validator struct {
	// Field name to be displayed in the error message
	// Can be customized by passing a string to the `Field()` function
	fieldName   string
	field       reflect.Value
	validations []func() error
}

// Field creates a new validator for a field.
// Pass the field name to be displayed in the error message or leave it empty to use the struct field name.
// Although you can pass fieldName as a variadic parameter this is done only to make it optional.
// If you pass more than one parameter, only the first one will be used.
//
// Example: Field("Name")
func Field(fieldName ...string) *Validator {
	name := ""
	if len(fieldName) > 0 {
		name = fieldName[0]
	}
	return &Validator{fieldName: name}
}

func Validate[T validable](v *T) error {
	schema := (*v).ValidationSchema()

	for key, validator := range schema {
		// Check if the field exists in the struct
		t := reflect.TypeOf(v).Elem()
		if _, ok := t.FieldByName(key); !ok {
			log.Panicf("corretto: field %s not found in struct %s", key, t.Name())
		}

		validator.field = reflect.ValueOf(v).Elem().FieldByName(key)

		for _, checkValidation := range validator.validations {
			err := checkValidation()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
