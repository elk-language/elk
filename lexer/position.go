package lexer

import (
	"fmt"
)

// Position describes an arbitrary source position.
// Lines and columns must be > 0.
type Position struct {
	StartByte  int // Index of the first byte of the lexeme
	ByteLength int // Number of bytes of the lexeme
	Line       int // Source line number where the lexeme starts
	Column     int // Source column number where the lexeme starts
}

// Retrieve the position, used in interfaces.
func (p Position) Pos() Position {
	return p
}

// String returns a string formatted like that:
//
//	line:column
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// Join two positions into one.
func (left Position) Join(right Position) Position {
	return Position{
		StartByte:  left.StartByte,
		ByteLength: right.StartByte - left.StartByte + right.ByteLength,
		Line:       left.Line,
		Column:     left.Column,
	}
}
