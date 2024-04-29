package corretto

import (
	"fmt"
	"strings"
)

type validationErr struct {
	Err error
}

func (v validationErr) Error() string {
	return v.Err.Error()
}

func newValidationError(msg string, cmsg string, args ...any) validationErr {
	if cmsg != "" {
		msg = cmsg
	}

	// Count the number of %v placeholders in the format string
	// to truncate the arguments slice if necessary
	numPlaceholders := strings.Count(msg, "%v")

	// If there are more arguments than placeholders, truncate the arguments slice
	if len(args) > numPlaceholders {
		args = args[:numPlaceholders]
	}

	return validationErr{Err: fmt.Errorf(msg, args...)}
}
