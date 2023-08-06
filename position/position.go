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
	StartByte  int // Index of the first byte of the source code fragment
	ByteLength int // Number of bytes of the source code fragment
	Line       int // Source line number where the fragment starts
	Column     int // Source column number where the fragment starts
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

func (p *Position) SetPos(pos *Position) {
	p.StartByte = pos.StartByte
	p.ByteLength = pos.ByteLength
	p.Line = pos.Line
	p.Column = pos.Column
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
// Works properly when the receiver is nil or the argument is nil.
func (left *Position) Join(right *Position) *Position {
	if right == nil {
		return left
	}
	if left == nil {
		return right
	}

	return &Position{
		StartByte:  left.StartByte,
		ByteLength: right.StartByte - left.StartByte + right.ByteLength,
		Line:       left.Line,
		Column:     left.Column,
	}
}

// Joins the given position with the last element of the given collection.
func JoinLastElement[Element Interface](left *Position, rightCollection []Element) *Position {
	if len(rightCollection) > 0 {
		return left.Join(rightCollection[len(rightCollection)-1].Pos())
	}

	return left
}

// Retrieve the position of the last element of a collection.
func OfLastElement[Element Interface](collection []Element) *Position {
	if len(collection) > 0 {
		return collection[len(collection)-1].Pos()
	}

	return nil
}
