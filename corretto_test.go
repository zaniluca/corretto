package corretto

import (
	"io"
	"testing"
)

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
					"Field1": Field().Number().Min(10, tt.customMessage),
				}

				err := schema.Parse(&struct{ Field1 int }{Field1: 5})
				if err != nil && err.Error() != tt.expectedError {
					t.Errorf("Parse() should have returned custom message")
					t.Errorf("expected: %s, got: %s", tt.expectedError, err.Error())
				}
			})
		}
	})
}
