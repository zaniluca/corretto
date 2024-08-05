package corretto

import (
	"fmt"
	"io"
	"testing"
)

func TestString(t *testing.T) {
	logger.SetOutput(io.Discard)

	schema := Schema{
		"stringField": Field().String(),
	}

	tests := []struct {
		name        string
		stringField any
		expectError bool
	}{
		{"empty string", "", false},
		{"zero value non string", 0, true},
		{"valid string", "hello world", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Parse(&struct{ stringField any }{stringField: tc.stringField})
			if tc.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestStringNonEmpty(t *testing.T) {
	logger.SetOutput(io.Discard)

	schema := Schema{
		"stringField": Field().String().NonEmpty(),
	}

	tests := []struct {
		name        string
		stringField string
		expectError bool
	}{
		{"empty string", "", true},
		{"non-empty string", "hello world", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Parse(&struct{ stringField string }{stringField: tc.stringField})
			if tc.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestStringMinLength(t *testing.T) {
	logger.SetOutput(io.Discard)

	schema := Schema{
		"stringField": Field().String().MinLength(5),
	}

	tests := []struct {
		name        string
		stringField any
		expectError bool
	}{
		{"empty string", "", true},
		{"valid string", "hello world", false},
		{"string with less than 5 characters", "foo", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Parse(&struct{ stringField any }{stringField: tc.stringField})
			if tc.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestStringCustomValidation(t *testing.T) {
	logger.SetOutput(io.Discard)

	schema := Schema{
		"stringField": Field().String().Test(func(ctx Context, field string) error {
			if field == "foo" {
				return nil
			}
			return fmt.Errorf("field must be 'foo'")
		}),
	}

	tests := []struct {
		name        string
		stringField any
		expectError bool
	}{
		{"empty string", "", true},
		{"valid string", "foo", false},
		{"invalid string", "bar", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Parse(&struct{ stringField any }{stringField: tc.stringField})
			if tc.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestStringOneOf(t *testing.T) {
	logger.SetOutput(io.Discard)

	schema := Schema{
		"stringField": Field().String().OneOf([]string{"foo", "bar"}),
	}

	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{"empty string", "", true},
		{"valid string", "foo", false},
		{"invalid string", "baz", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Parse(&struct{ stringField string }{stringField: tc.value})
			if tc.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestStringPredefinedRegex(t *testing.T) {
	logger.SetOutput(io.Discard)

	t.Run("email", func(t *testing.T) {
		schema := Schema{
			"Email": Field().String().Email(),
		}

		tests := []struct {
			name        string
			email       string
			expectError bool
		}{
			{"empty email", "", false},
			{"valid email", "foo@bar.com", false},
			{"email without a domain", "foo@bar", true},
			{"email with subdomain", "foo@sub.bar.com", false},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := schema.Parse(&struct{ Email string }{Email: tc.email})
				if tc.expectError && err == nil {
					t.Errorf("Parse() should have returned an error")
				}
			})
		}
	})
}
