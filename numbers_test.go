package corretto

import (
	"io"
	"testing"
)

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
