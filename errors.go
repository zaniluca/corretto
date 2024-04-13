package corretto

import (
	"fmt"
)

type ValidationErr struct {
	Err error
}

func (v ValidationErr) Error() string {
	return v.Err.Error()
}

func NewValidationError(msg string) ValidationErr {
	return ValidationErr{Err: fmt.Errorf(msg)}
}
