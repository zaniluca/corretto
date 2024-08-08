package corretto

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
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
	schema := Schema{
		"stringField": Field().String().NonEmpty(),
	}

	tests := []struct {
		name        string
		stringField string
		expectError bool
	}{
		{"empty string", "", true},
		{"whitespace string", " ", true},
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

	t.Run("uuid", func(t *testing.T) {
		schema := Schema{
			"UUID": Field().String().Uuid(),
		}

		tests := []struct {
			name        string
			uuid        string
			expectError bool
		}{
			{"empty uuid", "", false},
			{"valid uuid", "550e8400-e29b-41d4-a716-446655440000", false},
			{"invalid uuid", "foo", true},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := schema.Parse(&struct{ UUID string }{UUID: tc.uuid})
				if tc.expectError && err == nil {
					t.Errorf("Parse() should have returned an error")
				}
			})
		}
	})

	t.Run("cuid", func(t *testing.T) {
		schema := Schema{
			"CUID": Field().String().Cuid(),
		}

		tests := []struct {
			name        string
			cuid        string
			expectError bool
		}{
			{"empty cuid", "", false},
			{"valid cuid", "cikq1l4lk0001qzv9f5hyq0k7", false},
			{"invalid cuid", "foo", true},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := schema.Parse(&struct{ CUID string }{CUID: tc.cuid})
				if tc.expectError && err == nil {
					t.Errorf("Parse() should have returned an error")
				}
			})
		}
	})

	t.Run("hexcolor", func(t *testing.T) {
		schema := Schema{
			"HexColor": Field().String().HexColor(),
		}

		tests := []struct {
			name        string
			hexcolor    string
			expectError bool
		}{
			{"empty hexcolor", "", false},
			{"6 characters no alpha hex color", "#FF0000", false},
			{"6 characters with alpha hex color", "#FF0000FF", false},
			{"3 characters no alpha hex color", "#F00", false},
			{"3 characters with alpha hex color", "#F00F", false},
			{"case insensitiveness", "#ff0000ff", false},
			{"invalid hex color for length 1", "#FF000", true},
			{"invalid hex color for length 2", "#F", true},
			{"invalid hex color for length 3", "#FF000000000", true},
			{"invalid hex color for format", "#FF00G0", true},
			{"invalid hex color missing hash", "FF0000", true},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := schema.Parse(&struct{ HexColor string }{HexColor: tc.hexcolor})
				if tc.expectError && err == nil {
					t.Errorf("Parse() should have returned an error")
				}
			})
		}
	})

}

func TestStringIncludes(t *testing.T) {
	schema := Schema{
		"stringField": Field().String().Includes("foo"),
	}

	tests := []struct {
		name        string
		stringField string
		expectError bool
	}{
		{"empty string", "", true},
		{"string containing value", "foo", false},
		{"string without value", "bar", true},
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

func TestStringStartsWith(t *testing.T) {
	schema := Schema{
		"stringField": Field().String().StartsWith("foo"),
	}

	tests := []struct {
		name        string
		stringField string
		expectError bool
	}{
		{"empty string", "", true},
		{"string starting with prefix", "foobar", false},
		{"string not starting with prefix", "barfoo", true},
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

func TestStringEndsWith(t *testing.T) {
	schema := Schema{
		"stringField": Field().String().EndsWith("foo"),
	}

	tests := []struct {
		name        string
		stringField string
		expectError bool
	}{
		{"empty string", "", true},
		{"string ending with suffix", "barfoo", false},
		{"string not ending with suffix", "foobar", true},
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

func TestStringUrl(t *testing.T) {
	schema := Schema{
		"url": Field().String().Url(),
	}

	tests := []struct {
		name        string
		url         string
		expectError bool
	}{
		{"empty url", "", true},
		{"valid url", "http://example.com", false},
		{"url without protocol", "example.com", true},
		{"url with subdomain", "http://sub.example.com", false},
		{"url with path", "http://example.com/path", false},
		{"url with query", "http://example.com?query=foo", false},
		{"url with fragment", "http://example.com#fragment", false},
		{"url with port", "http://example.com:8080", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Parse(&struct{ url string }{url: tc.url})
			if tc.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}
