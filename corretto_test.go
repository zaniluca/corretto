package corretto

import (
	"io"
	"testing"
)

func TestParse(t *testing.T) {
	logger.SetOutput(io.Discard)

	t.Run("panics if schema has unknown field", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Parse() should have panicked")
			}
		}()

		schema := Schema{
			"UnexistingField": Field().Required(),
		}

		_ = schema.Parse(&struct{ Name string }{Name: "John"})
	})
}

func TestConcat(t *testing.T) {
	logger.SetOutput(io.Discard)

	t.Run("schema concatenation works", func(t *testing.T) {
		s := &struct{ Field1, Field2 string }{Field1: "John"}
		schema := Schema{
			"Field1": Field().Required(),
		}
		otherSchema := Schema{
			"Field2": Field().Required(),
		}
		err := schema.Parse(s)
		if err != nil {
			t.Errorf("Parse() should have returned nil because Field2 is not required")
		}

		schema.Concat(otherSchema)

		err = schema.Parse(s)
		if err == nil {
			t.Errorf("Parse() should have returned an error because Field2 is required")
		}
	})
}

func TestValidationOpts(t *testing.T) {
	logger.SetOutput(io.Discard)

	t.Run("custom message for validation", func(t *testing.T) {
		tests := []struct {
			name          string
			customMessage string
			expectedError string
		}{
			{
				name:          "default message",
				customMessage: "",
				expectedError: "Field1 must be at least 10",
			},
			{
				name:          "without placeholders",
				customMessage: "Field1 is supposed to be a minimum of 10",
				expectedError: "Field1 is supposed to be a minimum of 10",
			},
			{
				name:          "with placeholders",
				customMessage: "%v is supposed to be a minimum of %v",
				expectedError: "Field1 is supposed to be a minimum of 10",
			},
			{
				name:          "with extra placeholders",
				customMessage: "%v is supposed to be a minimum of %v and %v",
				expectedError: "Field1 is supposed to be a minimum of 10 and %!v(MISSING)",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				schema := Schema{
					"Field1": Field().Min(10, ValidationOpts{Message: tt.customMessage}),
				}

				err := schema.Parse(&struct{ Field1 int }{Field1: 5})
				if err != nil && err.Error() != tt.expectedError {
					t.Errorf("Parse() should have returned custom message")
				}
			})
		}
	})
}
