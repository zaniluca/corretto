package corretto

import (
	"io"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("panics if schema has unknown field", func(t *testing.T) {
		// Discard panic logs since they are expected
		logger.SetOutput(io.Discard)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Parse() should have panicked")
			}

			logger.SetOutput(os.Stderr)
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

func TestNestedSchemas(t *testing.T) {
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
		// Discard panic logs since they are expected
		logger.SetOutput(io.Discard)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Parse() should have panicked")
			}

			logger.SetOutput(os.Stderr)
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
