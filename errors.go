package corretto

import (
	"fmt"
	"strings"
)

type ValidationErr struct {
	Err error
}

func (v ValidationErr) Error() string {
	return v.Err.Error()
}

func newValidationError(msg string, opts ValidationOpts, args ...any) ValidationErr {
	if opts.Message != "" {
		msg = opts.Message
	}

	// Count the number of %v placeholders in the format string
	// to truncate the arguments slice if necessary
	numPlaceholders := strings.Count(msg, "%v")

	// If there are more arguments than placeholders, truncate the arguments slice
	if len(args) > numPlaceholders {
		args = args[:numPlaceholders]
	}

	return ValidationErr{Err: fmt.Errorf(msg, args...)}
}
