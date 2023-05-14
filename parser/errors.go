package parser

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
)

// Represents a single syntax error.
// Position points to the invalid token.
type Error struct {
	*position.Position
	Message string
}

// Create a new parser error.
func NewError(pos *position.Position, msg string) *Error {
	return &Error{
		Position: pos,
		Message:  msg,
	}
}

// Implements the error interface.
func (e *Error) Error() string {
	return e.String()
}

// Implements the fmt.Stringer interface
func (e *Error) String() string {
	return fmt.Sprintf("%s: %s", e.Position.HumanString(), e.Message)
}

// Same as [String] but prepends the result with the specified path.
func (e *Error) StringWithPath(path string) string {
	return fmt.Sprintf("%s:%s", path, e.String())
}

// ErrorList is a list of *Errors.
// The zero value for an ErrorList is an empty ErrorList ready to use.
type ErrorList []*Error

// Add a new syntax error.
func (e *ErrorList) Add(message string, pos *position.Position) {
	*e = append(*e, NewError(pos, message))
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
