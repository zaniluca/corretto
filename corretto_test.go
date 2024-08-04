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
			"UnexistingField": Field().String().NonEmpty(),
		}

		_ = schema.Parse(&struct{ Name string }{Name: "John"})
	})

	t.Run("accepts both pointers and values", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Parse() should not have panicked")
			}
		}()

		schema := Schema{
			"Name": Field().String().NonEmpty(),
		}

		_ = schema.Parse(&struct{ Name string }{Name: "John"})
		_ = schema.Parse(struct{ Name string }{Name: "John"})
	})
}

func TestUnmarshal(t *testing.T) {
	logger.SetOutput(io.Discard)

	tests := []struct {
		name        string
		bytes       []byte
		dst         any
		expectError bool
	}{
		{
			name:  "valid json and valid data",
			bytes: []byte(`{"Name": "John", "Age": 30}`),
			dst: &struct {
				Name string
				Age  int
			}{},
			expectError: false,
		},
		{
			name:  "valid json and invalid data",
			bytes: []byte(`{"Name": "John", "Age": 12}`),
			dst: &struct {
				Name string
				Age  int
			}{},
			expectError: true,
		},
		{
			name:  "invalid json",
			bytes: []byte(`{"Name": "John", "Age": `),
			dst: &struct {
				Name string
				Age  int
			}{},
			expectError: true,
		},
	}

	schema := Schema{
		"Name": Field().String().NonEmpty(),
		"Age":  Field().Number().Min(18),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Unmarshal(tt.bytes, tt.dst)
			if tt.expectError && err == nil {
				t.Errorf("Unmarshal() should have returned an error")
			}
		})
	}
}

func TestConcat(t *testing.T) {
	logger.SetOutput(io.Discard)

	s := &struct{ Field1, Field2 string }{Field1: "John"}
	schema := Schema{
		"Field1": Field().String().NonEmpty(),
	}
	otherSchema := Schema{
		"Field2": Field().String().NonEmpty(),
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

func TestNestedSchemas(t *testing.T) {
	logger.SetOutput(io.Discard)

	t.Run("validates nested field", func(t *testing.T) {
		type Nested struct {
			NestedField1 int
		}

		type Struct struct {
			Field1      string
			NestedField *Nested
		}

		v := &Struct{
			Field1: "John",
			NestedField: &Nested{
				NestedField1: 3,
			},
		}

		s2 := Schema{
			"NestedField1": Field().Number().Min(4),
		}

		s1 := Schema{
			"Field1":      Field().String().NonEmpty(),
			"NestedField": Field().Schema(s2),
		}

		err := s1.Parse(v)

		if err == nil {
			t.Errorf("Parse() should have returned an error because NestedField1 is not valid")
		}
		if err.Error() != "NestedField1 must be at least 4" {
			t.Errorf("Parse() should have returned the correct error message")
		}
	})

	t.Run("panics if field is not exported", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Parse() should have panicked")
			}
		}()

		type Nested struct {
			NestedField1 int
		}

		type Struct struct {
			Field1      string
			nestedField *Nested // not exported
		}

		s2 := Schema{
			"NestedField1": Field().Number().Min(4),
		}

		s1 := Schema{
			"Field1":      Field().String().NonEmpty(),
			"nestedField": Field().Schema(s2),
		}

		v := &Struct{
			Field1: "John",
			nestedField: &Nested{
				NestedField1: 4,
			},
		}

		_ = s1.Parse(v)
	})
}

func TestMin(t *testing.T) {
	logger.SetOutput(io.Discard)

	tests := []struct {
		name        string
		intValue    int
		stringValue string
		floatValue  float64
		min         int
		expectError bool
	}{
		{
			name:        "value is less than min",
			intValue:    5,
			stringValue: "hello",
			floatValue:  5.0,
			min:         10,
			expectError: true,
		},
		{
			name:        "value is equal to min",
			intValue:    10,
			stringValue: "helloworld",
			floatValue:  10.0,
			min:         10,
			expectError: false,
		},
		{
			name:        "value is greater than min",
			intValue:    15,
			stringValue: "helloworld",
			floatValue:  15.0,
			min:         10,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := Schema{
				"Field1": Field().Number().Min(tt.min),
			}

			err := schema.Parse(&struct{ Field1 int }{Field1: tt.intValue})
			err = schema.Parse(&struct{ Field1 string }{Field1: tt.stringValue})
			err = schema.Parse(&struct{ Field1 float64 }{Field1: tt.floatValue})
			if tt.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}
		})
	}
}

func TestOneOf(t *testing.T) {
	logger.SetOutput(io.Discard)

	s := Schema{
		"Field1": Field().OneOf([]any{1, 2, 3}),
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
