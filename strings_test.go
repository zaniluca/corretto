package corretto

import (
	"io"
	"testing"
)

func TestPredefinedRegex(t *testing.T) {
	logger.SetOutput(io.Discard)

	t.Run("email", func(t *testing.T) {
		schema := Schema{
			"Email": Field().Email(),
		}

		testCases := []struct {
			name        string
			email       string
			expectError bool
		}{
			{"empty email", "", false},
			{"valid email", "foo@bar.com", false},
			{"email without a domain", "foo@bar", true},
			{"email with subdomain", "foo@sub.bar.com", false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := schema.Parse(&struct{ Email string }{Email: tc.email})
				if tc.expectError && err == nil {
					t.Errorf("Parse() should have returned an error")
				}
			})
		}
	})
}
