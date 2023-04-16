package parser

import (
	"strings"

	"github.com/elk-language/elk/lexer"
)

// Represents a single syntax error.
// Position points to the invalid token.
type Error struct {
	lexer.Position
	Message string
}

// Implements the error interface.
func (e Error) Error() string {
	return e.String()
}

// Implements the fmt.Stringer interface
func (e Error) String() string {
	return e.Position.String() + ": " + e.Message
}

// ErrorList is a list of *Errors.
// The zero value for an ErrorList is an empty ErrorList ready to use.
type ErrorList []*Error

// Add a new syntax error.
func (e *ErrorList) Add(message string, pos lexer.Position) {
	*e = append(*e, &Error{pos, message})
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
