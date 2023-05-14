// Package position implements a struct
// that describes where a sequence of characters
// is located in a file.
package position

import (
	"fmt"
)

// Represents something that can be positioned.
type Interface interface {
	Pos() *Position
}

// Position describes an arbitrary source position.
// Lines and columns must be > 0.
type Position struct {
	StartByte  int // Index of the first byte of the lexeme
	ByteLength int // Number of bytes of the lexeme
	Line       int // Source line number where the lexeme starts
	Column     int // Source column number where the lexeme starts
}

// Create a new source position struct.
func New(start, length, line, column int) *Position {
	return &Position{
		StartByte:  start,
		ByteLength: length,
		Line:       line,
		Column:     column,
	}
}

// Retrieve the position, used in interfaces.
func (p *Position) Pos() *Position {
	return p
}

// String returns a string formatted like that:
//
//	line:column
func (p *Position) HumanString() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// Check if the position is valid.
func (p *Position) Valid() bool {
	if p.ByteLength > 0 && p.StartByte >= 0 && p.Column > 0 && p.Line > 0 {
		return true
	}

	return false
}

// Join two positions into one.
func (left *Position) Join(right *Position) *Position {
	return &Position{
		StartByte:  left.StartByte,
		ByteLength: right.StartByte - left.StartByte + right.ByteLength,
		Line:       left.Line,
		Column:     left.Column,
	}
}