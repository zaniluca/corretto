package corretto

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"
)

var (
	logger = log.New(log.Writer(), "corretto: ", log.LstdFlags)
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

// Utility to return the first parameter of a variadic function and log a warning if more than one parameter is passed
func optional[T any](params []T) T {
	if len(params) == 1 {
		return params[0]
	} else if len(params) > 1 {
		// Get the caller function name
		pc, _, _, _ := runtime.Caller(1)
		caller := runtime.FuncForPC(pc)
		// Get the file and line number that the caller function was called
		// This is done to show the correct file and line number in the log message
		_, file, line, _ := runtime.Caller(2)
		filename := filepath.Base(file)

		logger.Printf("WARN %s (%s:%d): calling with more than one parameter, only the first one will be used", caller.Name(), filename, line)
	}

	// Returns the zero value of the type
	return *new(T)
}

// Field creates a new validator for a field.
// Pass the field name to be displayed in the error message or leave it empty to use the struct field name.
// Although you can pass fieldName as a variadic parameter this is done only to make it optional.
// If you pass more than one parameter, only the first one will be used.
//
// Example: Field("Name")
func Field(fieldName ...string) *Validator {
	name := optional(fieldName)

	return &Validator{fieldName: name}
}

func Validate[T validable](v *T) error {
	schema := (*v).ValidationSchema()

	for key, validator := range schema {
		// Check if the field exists in the struct
		t := reflect.TypeOf(v).Elem()
		if _, ok := t.FieldByName(key); !ok {
			logger.Panicf("field %s not found in struct %s", key, t.Name())
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
