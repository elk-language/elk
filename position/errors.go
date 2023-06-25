package position

import (
	"fmt"
	"strings"
)

// Represents a single error in a particular source location.
type Error struct {
	*Location
	Message string
}

// Create a new error.
func NewError(loc *Location, msg string) *Error {
	return &Error{
		Location: loc,
		Message:  msg,
	}
}

// Implements the error interface.
func (e *Error) Error() string {
	return e.String()
}

// Implements the fmt.Stringer interface
func (e *Error) String() string {
	return fmt.Sprintf("%s: %s", e.Location.String(), e.Message)
}

// ErrorList is a list of *Errors.
// The zero value for an ErrorList is an empty ErrorList ready to use.
type ErrorList []*Error

// Create a new slice containing the elements of the
// two given error lists.
func (e ErrorList) Join(other ErrorList) ErrorList {
	n := len(e)
	return append(e[:n:n], other...)
}

// Add a new syntax error.
func (e *ErrorList) Add(message string, loc *Location) {
	*e = append(*e, NewError(loc, message))
}

// Implements the error interface.
func (e ErrorList) Error() string {
	switch len(e) {
	case 0:
		return "no errors"
	case 1:
		return e[0].Error()
	}

	var result strings.Builder
	for _, err := range e {
		result.WriteString(err.String())
		result.WriteByte('\n')
	}
	return result.String()
}

// Err returns an error equivalent to this error list.
// If the list is empty, Err returns nil.
func (e ErrorList) Err() error {
	if len(e) == 0 {
		return nil
	}
	return e
}
