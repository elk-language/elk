package position

import "fmt"

// Describes an arbitrary source position
// in a particular file.
// Lines and columns must be > 0.
type Location struct {
	Span
	Filename string
}

// Create a new location struct.
func NewLocation(filename string, startPos *Position, endPos *Position) *Location {
	return &Location{
		Span: Span{
			StartPos: startPos,
			EndPos:   endPos,
		},
		Filename: filename,
	}
}

// Create a new location with a given position.
func NewLocationWithSpan(filename string, span *Span) *Location {
	return &Location{
		Span:     *span,
		Filename: filename,
	}
}

// String representation of the location.
func (l *Location) String() string {
	if l == nil {
		return ""
	}

	return fmt.Sprintf("%s:%s", l.Filename, l.StartPos.String())
}
