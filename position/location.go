package position

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"unicode/utf8"

	"github.com/elk-language/elk/lexer/colorizer"
	"github.com/fatih/color"
)

func SpliceLocation(target, current *Location, unqoute bool) *Location {
	if unqoute {
		return spliceLocationUnquote(target, current)
	}

	return spliceLocation(target, current)
}

func spliceLocationUnquote(target, current *Location) *Location {
	if current != nil {
		return current
	}

	return target
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

// Return a string representation of this error
// that can be presented to humans.
//
// Reads the content of the file using the OS.
func (l *Location) HumanString(style bool, colorizer colorizer.Colorizer, accentColor *color.Color) (string, error) {
	fmt.Printf("PIPA!!!")
	debug.PrintStack()
	source, err := os.ReadFile(l.FilePath)
	if err != nil {
		return "", err
	}
	return l.HumanStringWithSource(string(source), style, colorizer, accentColor)
}

// Return a string representation of this error
// that can be presented to humans.
//
// Tries to fetch the source from the provided map (keys are filepaths, values are source strings).
// If there is no entry for the required filepath it tries to read the file through the OS.
func (l *Location) HumanStringWithSourceMap(style bool, colorizer colorizer.Colorizer, accentColor *color.Color, sourceMap map[string]string) (string, error) {
	if sourceMap == nil {
		return l.HumanString(style, colorizer, accentColor)
	}

	source, ok := sourceMap[l.FilePath]
	if !ok {
		sourceBytes, err := os.ReadFile(l.FilePath)
		if err != nil {
			return "", err
		}
		source = string(sourceBytes)
	}
	return l.HumanStringWithSource(source, style, colorizer, accentColor)
}

const (
	maxSourceExampleLength = 80
	ellipsis               = "..."
)

// Return a string representation of this error
// that can be presented to humans.
//
// Uses the given source string instead of reading the content of the file.
func (l *Location) HumanStringWithSource(source string, style bool, colorizer colorizer.Colorizer, accentColor *color.Color) (string, error) {
	var result strings.Builder

	result.WriteString(accentColor.Sprint(l.String()))

	if len(source) == 0 {
		return result.String(), nil
	}

	var startOffset int
	result.WriteString("\n  ")
	lineNumberStr := fmt.Sprint(l.StartPos.Line)
	faintColor := color.New(color.Faint)
	if !style {
		faintColor.DisableColor()
	}
	result.WriteString(faintColor.Sprint(lineNumberStr))
	result.WriteString(faintColor.Sprint(" | "))
	startOffset += 5 + len(lineNumberStr)

	lineStartIndex := strings.LastIndexByte(source[:l.StartPos.ByteOffset], '\n')
	if lineStartIndex == -1 {
		lineStartIndex = 0
	}
	lineEndIndex := strings.IndexByte(source[l.StartPos.ByteOffset:], '\n')
	if lineEndIndex == -1 {
		lineEndIndex = len(source)
	} else {
		lineEndIndex = l.StartPos.ByteOffset + lineEndIndex
	}
	errorSourceLength := utf8.RuneCountInString(source[l.StartPos.ByteOffset : l.EndPos.ByteOffset+1])
	var currentSourceLength int
	var currentErrorLength int
	var ellipsisStart bool
	var ellipsisEnd bool
	sourceFragmentStartIndex := l.StartPos.ByteOffset
	sourceFragmentEndIndex := l.EndPos.ByteOffset + 1

	if errorSourceLength < maxSourceExampleLength {
		leftLength := maxSourceExampleLength - errorSourceLength
		beforeSource := source[:l.StartPos.ByteOffset]
	backtrackLoop:
		for {
			if leftLength == 0 {
				break
			}
			char, size := utf8.DecodeLastRuneInString(beforeSource)
			switch char {
			case utf8.RuneError, '\r', '\n':
				break backtrackLoop
			}

			beforeSource = beforeSource[:len(beforeSource)-size]
			currentSourceLength++
			leftLength--
		}
		if len(beforeSource) > lineStartIndex+1 {
			ellipsisStart = true
			startOffset += len(ellipsis)
		}
		sourceFragmentStartIndex = len(beforeSource)
		s := source[sourceFragmentStartIndex:l.StartPos.ByteOffset]
		startOffset += utf8.RuneCountInString(s)
	}

	exampleEnd := l.EndPos.ByteOffset
	if lineEndIndex < exampleEnd {
		exampleEnd = lineEndIndex
	}
	for i := range source[l.StartPos.ByteOffset : exampleEnd+1] {
		if currentSourceLength >= maxSourceExampleLength {
			if i < lineEndIndex-1 {
				ellipsisEnd = true
			}
			break
		}
		currentSourceLength++
		currentErrorLength++
		sourceFragmentEndIndex = l.StartPos.ByteOffset + i
	}
	for i := range source[exampleEnd:lineEndIndex] {
		if currentSourceLength >= maxSourceExampleLength {
			if i < lineEndIndex-1 {
				ellipsisEnd = true
			}
			break
		}

		currentSourceLength++
		sourceFragmentEndIndex = l.EndPos.ByteOffset + i
	}
	if ellipsisStart {
		result.WriteString(faintColor.Sprint(ellipsis))
	}
	sourceFragment := source[sourceFragmentStartIndex : sourceFragmentEndIndex+1]
	var sourceFragmentBuff strings.Builder
	// replace tabs with spaces
	for _, char := range sourceFragment {
		if char == '\t' {
			sourceFragmentBuff.WriteByte(' ')
			continue
		}
		sourceFragmentBuff.WriteRune(char)
	}
	sourceFragment = sourceFragmentBuff.String()

	if style && colorizer != nil {
		var err error
		sourceFragment, err = colorizer.Colorize(sourceFragment)
		if err != nil {
			return "", err
		}
	}
	result.WriteString(sourceFragment)
	if ellipsisEnd {
		result.WriteString(faintColor.Sprint(ellipsis))
	}
	result.WriteByte('\n')
	result.WriteString(strings.Repeat(" ", startOffset))

	if currentErrorLength <= 1 {
		result.WriteString(accentColor.Sprint("│"))
	} else {
		result.WriteString(accentColor.Sprint("└"))
		result.WriteString(accentColor.Sprint(strings.Repeat("─", currentErrorLength-2)))
		result.WriteString(accentColor.Sprint("┤"))
	}
	result.WriteByte('\n')
	var spaceCount int
	if currentErrorLength == 0 {
		spaceCount = startOffset
	} else {
		spaceCount = startOffset + currentErrorLength - 1
	}
	result.WriteString(strings.Repeat(" ", spaceCount))
	result.WriteString(accentColor.Sprint("└ Here\n"))
	return result.String(), nil
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
