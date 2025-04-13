package position

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SpliceLocation(target, current *Location, unqoute bool) *Location {
	if unqoute {
		return spliceLocationUnquote(target, current)
	}

	return spliceLocation(target, current)
}

func spliceLocation(target, current *Location) *Location {
	if target == nil {
		return current
	}
	if current == nil {
		return target
	}

	result := target.Copy()
	result.Parent = current
	return result
}

func spliceLocationUnquote(target, current *Location) *Location {
	if current != nil {
		return current
	}

	return target
}

// Represents something that contains a location.
type LocationInterface interface {
	Location() *Location
	SetLocation(*Location)
}

// Describes an arbitrary source position
// in a particular file.
// Lines and columns must be > 0.
type Location struct {
	*Span
	FilePath string
	Parent   *Location
}

var DefaultLocation = NewLocation("<main>", DefaultSpan)
var ZeroLocation = NewLocation("", ZeroSpan)

// Create a new location with a given position.
func NewLocation(filename string, span *Span) *Location {
	return &Location{
		Span:     span,
		FilePath: filename,
	}
}

// Create a new location with a given position.
func NewLocationWithParent(filename string, span *Span, parent *Location) *Location {
	return &Location{
		Span:     span,
		FilePath: filename,
		Parent:   parent,
	}
}

func (l *Location) Copy() *Location {
	return &Location{
		Span:     l.Span,
		FilePath: l.FilePath,
		Parent:   l.Parent,
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
	return l.Span.Equal(other.Span) &&
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

// Join two locations into one.
// Works properly when the receiver is nil or the argument is nil.
func (left *Location) Join(right *Location) *Location {
	if right == nil {
		return left
	}
	if left == nil {
		return right
	}

	return &Location{
		Span:     left.Span.Join(right.Span),
		FilePath: left.FilePath,
	}
}

// Join two locations into one.
// Works properly when the receiver is nil or the argument is nil.
func (left *Location) JoinSpan(right *Span) *Location {
	if right == nil {
		return left
	}
	if left == nil {
		return nil
	}

	return &Location{
		Span:     left.Span.Join(right),
		FilePath: left.FilePath,
	}
}

// Retrieve the location of the last element of a collection.
func LocationOfLastElement[Element LocationInterface](collection []Element) *Location {
	if len(collection) > 0 {
		return collection[len(collection)-1].Location()
	}

	return nil
}

// Joins the given position with the last element of the given collection.
func JoinLocationOfLastElement[Element LocationInterface](left *Location, rightCollection []Element) *Location {
	if len(rightCollection) > 0 {
		return left.Join(rightCollection[len(rightCollection)-1].Location())
	}

	return left
}

// Join the position of the first element of a collection with the last one.
func JoinLocationOfCollection[Element LocationInterface](collection []Element) *Location {
	if len(collection) < 1 {
		return nil
	}

	left := collection[0].Location()
	if len(collection) == 1 {
		return left
	}

	return JoinLocationOfLastElement(left, collection)
}
