package corretto

import (
	"fmt"
	"io"
	"testing"
)

func TestNumber(t *testing.T) {
	logger.SetOutput(io.Discard)

	schema := Schema{
		"Field1": Field().Number(),
	}

	tests := []struct {
		name        string
		value       any
		expectError bool
	}{
		{
			name:        "zero value",
			value:       0,
			expectError: false,
		},
		{
			name:        "positive value",
			value:       42,
			expectError: false,
		},
		{
			name:        "negative value",
			value:       -42,
			expectError: false,
		},
		{
			name:        "float value",
			value:       42.0,
			expectError: false,
		},
		{
			name:        "string value",
			value:       "42",
			expectError: true,
		},
		{
			name:        "complex value",
			value:       42i,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Parse(&struct{ Field1 any }{Field1: tt.value})
			if tt.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestNumberMin(t *testing.T) {
	logger.SetOutput(io.Discard)

	tests := []struct {
		name        string
		value       int
		min         int
		expectError bool
	}{
		{
			name:        "value is less than min",
			value:       5,
			min:         10,
			expectError: true,
		},
		{
			name:        "value is equal to min",
			value:       10,
			min:         10,
			expectError: false,
		},
		{
			name:        "value is greater than min",
			value:       15,
			min:         10,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := Schema{
				"Field1": Field().Number().Min(tt.min),
			}

			err := schema.Parse(&struct{ Field1 int }{Field1: tt.value})
			if tt.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestNumberCustomValidation(t *testing.T) {
	logger.SetOutput(io.Discard)

	schema := Schema{
		"stringField": Field().Number().Test(func(ctx Context, field int) error {
			if field == 41 {
				return nil
			}
			return fmt.Errorf("field must be 41")
		}),
	}

	tests := []struct {
		name        string
		value       any
		expectError bool
	}{
		{
			name:        "value is not 41",
			value:       42,
			expectError: true,
		},
		{
			name:        "value is 41",
			value:       41,
			expectError: false,
		},
		{
			name:        "value is 41 but as a float",
			value:       41.0,
			expectError: true,
		},
		{
			name:        "value is 41 but as a string",
			value:       "41",
			expectError: true,
		},
		{
			name:        "value is not exactly 41 but as a float",
			value:       41.3,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Parse(&struct{ stringField any }{stringField: tt.value})
			if tt.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestNumberOneOf(t *testing.T) {
	logger.SetOutput(io.Discard)

	s := Schema{
		"Field1": Field().Number().OneOf([]int{1, 2, 3}),
	}

	tests := []struct {
		name        string
		value       int
		expectError bool
	}{
		{
			name:        "value is not in the list",
			value:       4,
			expectError: true,
		},
		{
			name:        "value is in the list",
			value:       2,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Parse(&struct{ Field1 int }{Field1: tt.value})
			if tt.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}
