package position

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Describes an arbitrary source position
// in a particular file.
// Lines and columns must be > 0.
type Location struct {
	Span
	FilePath string
}

var DefaultLocation = NewLocationWithSpan("<main>", DefaultSpan)

// Create a new location struct.
func NewLocation(filename string, startPos *Position, endPos *Position) *Location {
	return &Location{
		Span: Span{
			StartPos: startPos,
			EndPos:   endPos,
		},
		FilePath: filename,
	}
}

// Create a new location with a given position.
func NewLocationWithSpan(filename string, span *Span) *Location {
	return &Location{
		Span:     *span,
		FilePath: filename,
	}
}

// Returns a path to the file that is relative to the current working directory.
// if it's impossible to get the working directory or the file path cannot
// be transformed into a relative one, the original file path is returned instead.
func (l *Location) RelFilename() string {
	workingDir, err := os.Getwd()
	if err != nil {
		return l.FilePath
	}

	relPath, err := filepath.Rel(workingDir, l.FilePath)
	if err != nil {
		return l.FilePath
	}
	if strings.HasPrefix(relPath, "..") {
		return l.FilePath
	}

	return relPath
}

func (l *Location) Equal(other *Location) bool {
	if l == other {
		return true
	}
	return l.Span.Equal(&other.Span) &&
		l.FilePath == other.FilePath
}

// String representation of the location.
func (l *Location) String() string {
	if l == nil {
		return ""
	}

	l.RelFilename()
	return fmt.Sprintf("%s:%s", l.RelFilename(), l.StartPos.String())
}
