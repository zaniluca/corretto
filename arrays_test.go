package corretto

import (
	"io"
	"testing"
)

func TestArray(t *testing.T) {
	logger.SetOutput(io.Discard)

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

func TestArrayOf(t *testing.T) {
	logger.SetOutput(io.Discard)

	type innerStruct struct {
		innerField int
	}

	type testStruct struct {
		WithSchema    []innerStruct // Since we need to validate the inner field with a schema it must be exported
		withValidator []int
	}

	innerSchema := Schema{
		"innerField": Field().Min(3),
	}

	schema := Schema{
		"WithSchema":    Field().Array().Of(Field().Schema(innerSchema)),
		"withValidator": Field().Array().Of(Field().Min(3)),
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
