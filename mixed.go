package corretto

import (
	"fmt"
	"reflect"
)

func (v *Validator) Required() *Validator {
	v.validations = append(v.validations, func() error {
		// Check if the value is at its zero value
		if v.field.IsZero() {
			return fmt.Errorf("Required")
		}
		return nil
	})
	return v
}

func (v *Validator) Min(min int) *Validator {
	v.validations = append(v.validations, func() error {
		switch v.field.Type().Kind() {
		case reflect.Int:
			if v.field.Int() < int64(min) {
				return fmt.Errorf("Min: %d", min)
			}
		case reflect.Float64:
			if v.field.Float() < float64(min) {
				return fmt.Errorf("Min: %d", min)
			}
		case reflect.String:
			if len(v.field.String()) < min {
				return fmt.Errorf("Min: %d", min)
			}
		default:
			logger.Panicf("unsopported type %s for Min(), can only be used with int, float or string", v.field.Type().Kind())
		}

		return nil
	})
	return v
}
