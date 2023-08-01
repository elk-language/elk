package errors

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/position"
	"github.com/fatih/color"
)

// Represents a single error in a particular source location.
type Error struct {
	*position.Location
	Message string
}

// Create a new error.
func NewError(loc *position.Location, msg string) *Error {
	return &Error{
		Location: loc,
		Message:  msg,
	}
}

// Implements the error interface.
func (e *Error) Error() string {
	return e.String()
}

// Implements the fmt.Stringer interface
func (e *Error) String() string {
	return fmt.Sprintf("%s: %s", e.Location.String(), e.Message)
}

// Return a string representation of this error
// that can be presented to humans.
func (e *Error) HumanString(style bool) (string, error) {
	source, err := os.ReadFile(e.Location.Filename)
	if err != nil {
		return "", err
	}
	return e.HumanStringWithSource(string(source), style), nil
}

// Return a string representation of this error
// that can be presented to humans.
func (e *Error) HumanStringWithSource(source string, style bool) string {
	var result strings.Builder
	c := color.New(color.Bold, color.FgRed)
	if !style {
		c.DisableColor()
	}
	result.WriteString(c.Sprint(e.Location.String()))
	result.WriteString(": ")
	result.WriteString(e.Message)
	result.WriteByte('\n')
	if len(source) == 0 {
		return result.String()
	}

	var startOffset int
	result.WriteString("\n  ")
	lineNumberStr := fmt.Sprint(e.Location.Line)
	result.WriteString(lineNumberStr)
	result.WriteString(" | ")
	startOffset += 5 + len(lineNumberStr)

	lineStartIndex := strings.LastIndexByte(source[:e.Location.StartByte], '\n')
	if lineStartIndex == -1 {
		lineStartIndex = 0
	}
	lineEndIndex := strings.IndexByte(source[e.Location.StartByte:], '\n')
	if lineEndIndex == -1 {
		lineEndIndex = len(source)
	} else {
		lineEndIndex = e.Location.StartByte + lineEndIndex
	}
	errorSourceLength := utf8.RuneCountInString(source[e.Location.StartByte : e.Location.StartByte+e.Location.ByteLength])
	var currentSourceLength int
	var currentErrorLength int
	var ellipsisStart bool
	var ellipsisEnd bool
	sourceFragmentStartIndex := e.Location.StartByte
	sourceFragmentEndIndex := e.Location.StartByte + e.Location.ByteLength
	if errorSourceLength < maxSourceExampleLength {
		leftLength := maxSourceExampleLength - errorSourceLength
		beforeSource := source[:e.Location.StartByte]
		for {
			if leftLength == 0 {
				break
			}
			char, size := utf8.DecodeLastRuneInString(beforeSource)
			if char == utf8.RuneError {
				break
			}

			beforeSource = beforeSource[:len(beforeSource)-size]
			currentSourceLength++
			leftLength--
		}
		if len(beforeSource) > lineStartIndex {
			ellipsisStart = true
			startOffset += len(ellipsis)
		}
		sourceFragmentStartIndex = len(beforeSource)
		s := source[sourceFragmentStartIndex:e.Location.StartByte]
		startOffset += utf8.RuneCountInString(s)
	}

	exampleEnd := e.Location.StartByte + e.Location.ByteLength
	if lineEndIndex < exampleEnd {
		exampleEnd = lineEndIndex
	}
	for i := range source[e.Location.StartByte:exampleEnd] {
		if currentSourceLength >= maxSourceExampleLength {
			if i < lineEndIndex-1 {
				ellipsisEnd = true
			}
			break
		}
		currentSourceLength++
		currentErrorLength++
		sourceFragmentEndIndex = e.Location.StartByte + i
	}
	for i := range source[exampleEnd:lineEndIndex] {
		if currentSourceLength >= maxSourceExampleLength {
			if i < lineEndIndex-1 {
				ellipsisEnd = true
			}
			break
		}

		currentSourceLength++
		sourceFragmentEndIndex = e.Location.StartByte + e.Location.ByteLength + i
	}
	ellipsisColor := color.New(color.FgHiBlack, color.FgHiBlack)
	if !style {
		ellipsisColor.DisableColor()
	}
	if ellipsisStart {
		result.WriteString(ellipsisColor.Sprint(ellipsis))
	}
	sourceFragment := source[sourceFragmentStartIndex : sourceFragmentEndIndex+1]
	if style {
		sourceFragment = lexer.Colorize(sourceFragment)
	}
	result.WriteString(sourceFragment)
	if ellipsisEnd {
		result.WriteString(ellipsisColor.Sprint(ellipsis))
	}
	result.WriteByte('\n')
	result.WriteString(strings.Repeat(" ", startOffset))
	lineColor := color.New(color.FgHiRed, color.Bold)
	if !style {
		lineColor.DisableColor()
	}

	if currentErrorLength <= 1 {
		result.WriteString(lineColor.Sprint("│"))
	} else {
		result.WriteString(lineColor.Sprint("└"))
		result.WriteString(lineColor.Sprint(strings.Repeat("─", currentErrorLength-2)))
		result.WriteString(lineColor.Sprint("┤"))
	}
	result.WriteByte('\n')
	var spaceCount int
	if currentErrorLength == 0 {
		spaceCount = startOffset
	} else {
		spaceCount = startOffset + currentErrorLength - 1
	}
	result.WriteString(strings.Repeat(" ", spaceCount))
	result.WriteString(lineColor.Sprint("└ Here\n"))
	return result.String()
}

// ErrorList is a list of *Errors.
// The zero value for an ErrorList is an empty ErrorList ready to use.
type ErrorList []*Error

// Create a new slice containing the elements of the
// two given error lists.
func (e ErrorList) Join(other ErrorList) ErrorList {
	n := len(e)
	return append(e[:n:n], other...)
}

// Add a new syntax error.
func (e *ErrorList) Add(message string, loc *position.Location) {
	*e = append(*e, NewError(loc, message))
}

// Implements the error interface.
func (e ErrorList) Error() string {
	switch len(e) {
	case 0:
		return "no errors"
	case 1:
		return e[0].Error()
	}

	var result strings.Builder
	for _, err := range e {
		result.WriteString(err.String())
		result.WriteByte('\n')
	}
	return result.String()
}

const (
	maxSourceExampleLength = 32
	ellipsis               = "..."
)

// Return a string representation of this error list
// that can be presented to humans.
func (el ErrorList) HumanString(style bool) (string, error) {
	var result strings.Builder
	for _, e := range el {
		msg, err := e.HumanString(style)
		if err != nil {
			return "", err
		}
		result.WriteString(msg)
	}
	return result.String(), nil
}

// Return a string representation of this error list
// that can be presented to humans.
func (e ErrorList) HumanStringWithSource(source string, style bool) string {
	var result strings.Builder
	for _, err := range e {
		result.WriteString(err.HumanStringWithSource(source, style))
	}
	return result.String()
}

// Err returns an error equivalent to this error list.
// If the list is empty, Err returns nil.
func (e ErrorList) Err() error {
	if len(e) == 0 {
		return nil
	}
	return e
}
