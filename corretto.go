package corretto

import (
	"encoding/json"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
)

var (
	logger = log.New(log.Writer(), "corretto: ", log.LstdFlags)
)

type Schema map[string]validator

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

type CustomValidationFunc func(ctx Context, field reflect.Value) error

type validator interface {
	getBaseValidator() *baseValidator
}

// Represents a validator for a field
type baseValidator struct {
	ctx         Context          // The context of the validation, usually the struct that contains the field
	fieldName   string           // The name of the field to be displayed in the error message, by default it uses the struct field name
	field       reflect.Value    // The value of the field to be validated
	validations []ValidationFunc // The list of validations to be performed
}

// Returns the baseValidator of the validator
func (v *baseValidator) getBaseValidator() *baseValidator {
	return v
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

// Field creates a new validator for a field.
// Pass the field name to be displayed in the error message or leave it empty to use the struct field name.
// Although you can pass fieldName as a variadic parameter this is done only to make it optional.
// If you pass more than one parameter, only the first one will be used.
//
// Example:
//
//	Field("Name")
func Field(fieldName ...string) *baseValidator {
	name := optional(fieldName)

	return &baseValidator{fieldName: name}
}

// Parse validates the struct fields based on the schema
// It returns an error if any of the validations fail or nil if all validations pass
//
// Example:
//
//		schema := Schema{
//			"FirstName": Field("Name").Min(3).Test(customValidation),
//			"Age":       Field().Min(18),
//			"Email":     Field().String().Email(),
//		}
//		user := User{
//			FirstName: "John",
//			Age:       17,
//			Email:     "john@doe.com",
//		}
//
//		err := schema.Parse(user) // ValidationError{Message: "Age must be at least 18"}
//	 	// you can pass a reference too
//		err := schema.Parse(&user) // ValidationError{Message: "Age must be at least 18"}
func (s Schema) Parse(value any) error {
	for key, validator := range s {
		var t reflect.Type
		var v reflect.Value

		// Check if the value is a pointer to a struct or a struct value
		if reflect.TypeOf(value).Kind() == reflect.Ptr {
			t = reflect.TypeOf(value).Elem()
			v = reflect.ValueOf(value).Elem()
		} else {
			t = reflect.TypeOf(value)
			v = reflect.ValueOf(value)
		}

		// Check if the field exists in the struct
		if _, ok := t.FieldByName(key); !ok {
			logger.Panicf("field %s not found in struct %s", key, t.Name())
		}

		baseValidator := validator.getBaseValidator()
		baseValidator.field = v.FieldByName(key)
		baseValidator.ctx = value
		// If no custom field name is provided, use the struct field name
		if baseValidator.fieldName == "" {
			baseValidator.fieldName = key
		}

		for _, checkValidation := range baseValidator.validations {
			err := checkValidation()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Unmarshal parses the JSON data into the struct and validates the fields based on the schema
func (s Schema) Unmarshal(data []byte, v any) error {
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	return s.Parse(v)
}

// MustParse behaves the same as [Schema.Parse] but panics if any of the validations fail
func (s Schema) MustParse(value any) {
	err := s.Parse(value)
	if err != nil {
		panic(err)
	}
}

// MustUnmarshal behaves the same as [Schema.Unmarshal] but panics if any of the validations fail or if the JSON data cannot be parsed
func (s Schema) MustUnmarshal(data []byte, v any) {
	err := s.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}

// Concat adds the fields from another [Schema] to the current schema
func (s Schema) Concat(other Schema) {
	for key, value := range other {
		s[key] = value
	}
}
