package position

import (
	"fmt"
	"os"
	"path/filepath"
)

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

// Returns a path to the file that is relative to the current working directory.
// if it's impossible to get the working directory or the file path cannot
// be transformed into a relative one, the original file path is returned instead.
func (l *Location) RelFilename() string {
	workingDir, err := os.Getwd()
	if err != nil {
		return l.Filename
	}

	relPath, err := filepath.Rel(workingDir, l.Filename)
	if err != nil {
		return l.Filename
	}

	return relPath
}

// String representation of the location.
func (l *Location) String() string {
	if l == nil {
		return ""
	}

	l.RelFilename()
	return fmt.Sprintf("%s:%s", l.RelFilename(), l.StartPos.String())
}
