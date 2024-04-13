package corretto_test

import (
	"corretto"
	"testing"
)

type T1 struct {
	Name string
}

func (t T1) ValidationSchema() corretto.Schema {
	return corretto.Schema{
		"Name":      corretto.Field().Min(3),
		"WrongName": corretto.Field().Min(10),
	}
}

func TestValidate(t *testing.T) {
	t.Run("panics if schema has unknown field", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Validate() should have panicked")
			}
		}()

		corretto.Validate(&T1{Name: "Name"})
	})
}
