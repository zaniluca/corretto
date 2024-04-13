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
			"Name":            Field().Min(3),
			"UnexistingField": Field().Min(10),
		}

		_ = schema.Parse(&struct{ Name string }{Name: "John"})
	})
}
