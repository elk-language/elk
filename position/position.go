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
	SetPos(*Position)
}

// Position describes an arbitrary source code position.
// Lines and columns must be > 0.
type Position struct {
	ByteOffset int // Index of the first byte of the source code fragment
	Line       int // Source line number where the fragment starts
	Column     int // Source column number where the fragment starts
}

// Create a new source position struct.
func New(byteOffset, line, column int) *Position {
	return &Position{
		ByteOffset: byteOffset,
		Line:       line,
		Column:     column,
	}
}

// Retrieve the position, used in interfaces.
func (p *Position) Pos() *Position {
	return p
}

func (p *Position) SetPos(pos *Position) {
	p.ByteOffset = pos.ByteOffset
	p.Line = pos.Line
	p.Column = pos.Column
}

// String returns a string formatted like that:
//
//	line:column
func (p *Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// Check if the position is valid.
func (p *Position) Valid() bool {
	if p.ByteOffset >= 0 && p.Column > 0 && p.Line > 0 {
		return true
	}

	return false
}
