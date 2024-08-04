package corretto

import (
	"encoding/json"
	"reflect"
)

type Schema map[string]validator

// Schema checks if the field can be parsed by the provided schema
// Use it to validate nested structs
//
// NOTE: the field associated with the schema must be exported
//
//	type Parent struct {
//		Son  *Son // Field.Schema() works fine
//		daughter *Daughter // Field.Schema() will panic
//	}
func (v *BaseValidator) Schema(s Schema) *BaseValidator {
	v.validations = append(v.validations, func() error {
		if !v.field.CanInterface() {
			logger.Panicf("field `%v` must be exported to be validated", v.key)
		}
		return s.Parse(v.field.Interface())
	})
	return v
}

// Parse validates the struct fields based on the schema
// It returns an error if any of the validations fail or nil if all validations pass
//
// Example:
//
//		schema := Schema{
//			"FirstName": Field("Name").String().MinLength(3).Test(customValidation),
//			"Age":       Field().Number().Min(18),
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
		baseValidator.key = key
		// If no custom field name is provided, use the struct field name
		if baseValidator.fieldName == "" {
			baseValidator.fieldName = key
		}

		// If any of the validations fail, return the error
		if err := baseValidator.check(); err != nil {
			return err
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
// If the field already exists, it will be overwritten
func (s Schema) Concat(other Schema) {
	for key, value := range other {
		s[key] = value
	}
}
