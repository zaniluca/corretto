package corretto

import (
	"fmt"
	"testing"
)

func TestArray(t *testing.T) {
	schema := Schema{
		"arrayField": Field().Array(),
	}

	tests := []struct {
		name        string
		arrayField  any
		expectError bool
	}{
		{"empty array", []int{}, false},
		{"zero value non array", 0, true},
		{"valid array", []int{1, 2, 3}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Parse(&struct{ arrayField any }{arrayField: tc.arrayField})
			if tc.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestArrayLength(t *testing.T) {
	schema := Schema{
		"arrayField": Field().Array().Length(3),
	}

	tests := []struct {
		name        string
		arrayField  []int
		expectError bool
	}{
		{"empty array", []int{}, true},
		{"array with 2 elements", []int{1, 2}, true},
		{"array with 3 elements", []int{1, 2, 3}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Parse(&struct{ arrayField []int }{arrayField: tc.arrayField})
			if tc.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestArrayCustomValidation(t *testing.T) {
	onlyPositiveIntegers := func(ctx Context) error {
		// if flag is set to true, the array must only contain positive integers

		// Testing cross-field validation
		v := ctx.(*struct {
			arrayField           []int
			onlyPositiveIntegers bool
		})
		if !v.onlyPositiveIntegers {
			return nil
		}

		for i := 0; i < len(v.arrayField); i++ {
			elem := v.arrayField[i]
			if elem <= 0 {
				return fmt.Errorf("Array must only contain positive integers")
			}
		}

		return nil
	}

	schema := Schema{
		"arrayField":           Field().Array().Test(onlyPositiveIntegers),
		"onlyPositiveIntegers": Field().Bool(),
	}

	tests := []struct {
		name                 string
		arrayField           []int
		onlyPositiveIntegers bool
		expectError          bool
	}{
		{"some negative integers but flag not set", []int{1, 2, 3, 0, -1}, false, false},
		{"some negative integers and flag set", []int{1, 2, 3, 0, -1}, true, true},
		{"all positive integers and flag set", []int{1, 2, 3}, true, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Parse(&struct {
				arrayField           []int
				onlyPositiveIntegers bool
			}{arrayField: tc.arrayField, onlyPositiveIntegers: tc.onlyPositiveIntegers})
			if tc.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestArrayMinLength(t *testing.T) {
	schema := Schema{
		"arrayField": Field().Array().MinLength(3),
	}

	tests := []struct {
		name        string
		arrayField  []int
		expectError bool
	}{
		{"empty array", []int{}, true},
		{"array with 2 elements", []int{1, 2}, true},
		{"array with 3 elements", []int{1, 2, 3}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Parse(&struct{ arrayField []int }{arrayField: tc.arrayField})
			if tc.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestArrayMaxLength(t *testing.T) {
	schema := Schema{
		"arrayField": Field().Array().MaxLength(3),
	}

	tests := []struct {
		name        string
		arrayField  []int
		expectError bool
	}{
		{"empty array", []int{}, false},
		{"array with 4 elements", []int{1, 2, 3, 4}, true},
		{"array with 3 elements", []int{1, 2, 3}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Parse(&struct{ arrayField []int }{arrayField: tc.arrayField})
			if tc.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestArrayOf(t *testing.T) {
	type innerStruct struct {
		innerField int
	}

	type testStruct struct {
		WithSchema    []innerStruct // Since we need to validate the inner field with a schema it must be exported
		withValidator []int
	}

	innerSchema := Schema{
		"innerField": Field().Number().Min(3),
	}

	schema := Schema{
		"WithSchema":    Field().Array().Of(Field().Schema(innerSchema)),
		"withValidator": Field().Array().Of(Field().Number().Min(3)),
	}

	testsWithValidator := []struct {
		name          string
		withValidator []int
		expectError   bool
	}{
		{"invalid array element", []int{1, 2, 1, 2}, true},
		{"valid array element", []int{3}, false},
	}

	for _, tc := range testsWithValidator {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Parse(&testStruct{withValidator: tc.withValidator})
			if tc.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}

	testsWithSchema := []struct {
		name        string
		WithSchema  []innerStruct
		expectError bool
	}{
		{"invalid array element with schema", []innerStruct{{innerField: 1}}, true},
		{"valid array element with schema", []innerStruct{{innerField: 3}}, false},
	}

	for _, tc := range testsWithSchema {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Parse(&testStruct{WithSchema: tc.WithSchema})
			if tc.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}
