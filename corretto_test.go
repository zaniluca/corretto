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
			"UnexistingField": Field().Min(10),
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
