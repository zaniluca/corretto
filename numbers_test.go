package corretto

import (
	"fmt"
	"math"
	"testing"
)

func TestNumber(t *testing.T) {
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
	schema := Schema{
		"stringField": Field().Number().Test(func(ctx Context) error {
			v := ctx.(struct{ stringField any })
			if v.stringField == 41 {
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
	schema := Schema{
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
			err := schema.Parse(&struct{ Field1 int }{Field1: tt.value})
			if tt.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

// Test against the number validator alias methods (methods that rely on other methods like Positive() or Negative())
func TestNumberAliases(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		schema := Schema{
			"Field1": Field().Number().Positive(),
		}

		tests := []struct {
			name        string
			value       any
			expectError bool
		}{
			{
				name:        "zero",
				value:       0,
				expectError: true,
			},
			{
				name:        "positive",
				value:       42,
				expectError: false,
			},
			{
				name:        "negative",
				value:       -42,
				expectError: true,
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
	})

	t.Run("negative", func(t *testing.T) {
		schema := Schema{
			"Field1": Field().Number().Negative(),
		}

		tests := []struct {
			name        string
			value       any
			expectError bool
		}{
			{
				name:        "zero",
				value:       0,
				expectError: true,
			},
			{
				name:        "positive",
				value:       42,
				expectError: true,
			},
			{
				name:        "negative",
				value:       -42,
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
	})

	t.Run("non-positive", func(t *testing.T) {
		schema := Schema{
			"Field1": Field().Number().NonPositive(),
		}

		tests := []struct {
			name        string
			value       any
			expectError bool
		}{
			{
				name:        "zero",
				value:       0,
				expectError: false,
			},
			{
				name:        "positive",
				value:       42,
				expectError: true,
			},
			{
				name:        "negative",
				value:       -42,
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
	})

	t.Run("non-negative", func(t *testing.T) {
		schema := Schema{
			"Field1": Field().Number().NonNegative(),
		}

		tests := []struct {
			name        string
			value       any
			expectError bool
		}{
			{
				name:        "zero",
				value:       0,
				expectError: false,
			},
			{
				name:        "positive",
				value:       42,
				expectError: false,
			},
			{
				name:        "negative",
				value:       -42,
				expectError: true,
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
	})
}

func TestNumberMultipleOf(t *testing.T) {
	schema := Schema{
		"Field1": Field().Number().MultipleOf(3),
	}

	tests := []struct {
		name        string
		value       any
		expectError bool
	}{
		{
			name:        "multiple of",
			value:       6,
			expectError: false,
		},
		{
			name:        "not multiple of",
			value:       5,
			expectError: true,
		},
		{
			name:        "zero",
			value:       0,
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

func TestNumberFinite(t *testing.T) {
	schema := Schema{
		"Field1": Field().Number().Finite(),
	}

	tests := []struct {
		name        string
		value       any
		expectError bool
	}{
		{
			name:        "finite",
			value:       42,
			expectError: false,
		},
		{
			name:        "positive infinity",
			value:       math.Inf(1),
			expectError: true,
		},
		{
			name:        "negative infinity",
			value:       math.Inf(-1),
			expectError: true,
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
