package position

import "fmt"

// Describes an arbitrary source position
// in a particular file.
// Lines and columns must be > 0.
type Location struct {
	Position
	Filename string
}

// Create a new location struct.
func NewLocation(filename string, start, length, line, column int) *Location {
	return &Location{
		Position: Position{
			StartByte:  start,
			ByteLength: length,
			Line:       line,
			Column:     column,
		},
		Filename: filename,
	}
}

// String representation of the location.
func (l *Location) String() string {
	if l == nil {
		return ""
	}

	return fmt.Sprintf("%s:%s", l.Filename, l.Position.HumanString())
}
