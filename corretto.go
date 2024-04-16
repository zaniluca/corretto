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

type ValidationFunc func() error

// Context is the whole struct that contains the field to be validated
// It can be used to access other fields in the struct and perform validations based on them
// Although it is defined as any, it is actually a POINTER to the struct
//
// It is recommended to use a type assertion to convert it to the correct type
// Example:
//
//	u, ok := ctx.(*User)
//	if !ok {
//	   panic("Invalid context")
//	}
type Context any

type CustomValidationFunc func(ctx Context, field reflect.Value) error

// Represents a validator for a field
type Validator struct {
	// The context of the validation, usually the struct that contains the field
	ctx         Context
	fieldName   string
	field       reflect.Value
	validations []ValidationFunc
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
// Example:
//
//	Field("Name")
func Field(fieldName ...string) *Validator {
	name := optional(fieldName)

	return &Validator{fieldName: name}
}

// Parse validates the struct fields based on the schema
// It returns an error if any of the validations fail or nil if all validations pass
//
// Example:
//
//	schema := Schema{
//		"FirstName": Field("Name").Min(3).Test(customValidation),
//		"Age":       Field().Min(18),
//		"Email":     Field().Email(),
//	}
//	user := &User{
//		FirstName: "John",
//		Age:       17,
//		Email:     "john@doe.com",
//	}
//
//	// Remember that the struct MUST BE A POINTER
//	err := schema.Parse(user) // ValidationError{Message: "Age must be at least 18"}
func (s Schema) Parse(value any) error {
	for key, validator := range s {

		var t reflect.Type
		if reflect.TypeOf(value).Kind() == reflect.Ptr {
			t = reflect.TypeOf(value).Elem()
		} else {
			logger.Panicf("value must be a pointer to a struct")
		}

		// Check if the field exists in the struct
		if _, ok := t.FieldByName(key); !ok {
			logger.Panicf("field %s not found in struct %s", key, t.Name())
		}

		validator.field = reflect.ValueOf(value).Elem().FieldByName(key)
		validator.ctx = value
		// If no custom field name is provided, use the struct field name
		if validator.fieldName == "" {
			validator.fieldName = key
		}

		for _, checkValidation := range validator.validations {
			err := checkValidation()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Concat adds the fields from another schema to the current schema
func (s Schema) Concat(other Schema) {
	for key, value := range other {
		s[key] = value
	}
}

type ValidationOpts struct {
	// The error message to be displayed if the validation fails
	Message string
}
