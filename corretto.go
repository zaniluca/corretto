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

const (
	oneOfErrorMsg = "%v must be one of %v"
)

type ValidationFunc func() error

// Context is the whole struct that contains the field to be validated
// It can be used to access other fields in the struct and perform validations based on them
// Although it is defined as any, it is actually the struct that contains the field to be validated
//
// It is recommended to use a type assertion to convert it to the correct type
// Example:
//
//	u, ok := ctx.(User)
//	if !ok {
//	   panic("Invalid context")
//	}
type Context any

type CustomValidationFunc func(ctx Context) error

type validator interface {
	getBaseValidator() *BaseValidator
	check() error
}

// getBaseValidator returns the underlying baseValidator
func (v *BaseValidator) getBaseValidator() *BaseValidator {
	return v
}

// Check if the field is valid by running all validations
// If any of the validations fail, return the error
func (v *BaseValidator) check() error {
	for _, checkValidation := range v.validations {
		err := checkValidation()
		if err != nil {
			return err
		}
	}

	return nil
}

// Represents a validator for a field
type BaseValidator struct {
	ctx         Context          // The context of the validation, usually the struct that contains the field
	fieldName   string           // The name of the field to be displayed in the error message, by default it uses the struct field name
	field       reflect.Value    // The value of the field to be validated
	validations []ValidationFunc // The list of validations to be performed
	key         string           // field name in the struct (and key in the Schema)
}

// Utility to return the first parameter of a variadic function and log a warning if more than one parameter is passed
// If no parameter is passed, it returns the zero value of the type
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

func oneOf[T comparable](v T, allowed []T) bool {
	for _, a := range allowed {
		if v == a {
			return true
		}
	}
	return false
}

// Field creates a new validator for a field.
// Pass the field name to be displayed in the error message or leave it empty to use the struct field name.
// Although you can pass fieldName as a variadic parameter this is done only to make it optional.
// If you pass more than one parameter, only the first one will be used.
//
// Example:
//
//	Field("Name")
//	Field("Name", "MyName") // This will log a warning and use "Name" as the field name
func Field(fieldName ...string) *BaseValidator {
	name := optional(fieldName)

	return &BaseValidator{fieldName: name}
}
