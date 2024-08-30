package error

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/position"
	"github.com/fatih/color"
)

type Severity uint8

func (s Severity) String() string {
	return severityNames[s]
}

func (s Severity) Tag(style bool) string {
	if !style {
		return fmt.Sprintf("[%s]", s.String())
	}

	c := color.New(color.Bold, s.color())
	return c.Sprintf("[%s]", s.String())
}

func (s Severity) color() color.Attribute {
	return severityColor[s]
}

const (
	FAILURE Severity = iota
	WARNING
)

var severityNames = []string{
	FAILURE: "FAIL",
	WARNING: "WARN",
}

var severityColor = []color.Attribute{
	FAILURE: color.FgRed,
	WARNING: color.FgYellow,
}

// Represents a single error in a particular source location.
type Error struct {
	*position.Location
	Severity Severity
	Message  string
}

// Create a new error.
func NewError(loc *position.Location, msg string, severity Severity) *Error {
	return &Error{
		Location: loc,
		Message:  msg,
		Severity: severity,
	}
}

// Create a new warning.
func NewFailure(loc *position.Location, msg string) *Error {
	return &Error{
		Location: loc,
		Message:  msg,
		Severity: FAILURE,
	}
}

// Create a new warning.
func NewWarning(loc *position.Location, msg string) *Error {
	return &Error{
		Location: loc,
		Message:  msg,
		Severity: WARNING,
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
	severityColor := color.New(color.Bold, e.Severity.color())
	if !style {
		severityColor.DisableColor()
	}
	result.WriteString(severityColor.Sprint(e.Location.String()))
	result.WriteString(": ")

	result.WriteString(e.Severity.Tag(style))
	result.WriteByte(' ')
	result.WriteString(e.Message)
	result.WriteByte('\n')
	if len(source) == 0 {
		return result.String()
	}

	var startOffset int
	result.WriteString("\n  ")
	lineNumberStr := fmt.Sprint(e.Location.StartPos.Line)
	faintColor := color.New(color.Faint)
	if !style {
		faintColor.DisableColor()
	}
	result.WriteString(faintColor.Sprint(lineNumberStr))
	result.WriteString(faintColor.Sprint(" | "))
	startOffset += 5 + len(lineNumberStr)

	lineStartIndex := strings.LastIndexByte(source[:e.Location.StartPos.ByteOffset], '\n')
	if lineStartIndex == -1 {
		lineStartIndex = 0
	}
	lineEndIndex := strings.IndexByte(source[e.Location.StartPos.ByteOffset:], '\n')
	if lineEndIndex == -1 {
		lineEndIndex = len(source)
	} else {
		lineEndIndex = e.Location.StartPos.ByteOffset + lineEndIndex
	}
	errorSourceLength := utf8.RuneCountInString(source[e.Location.StartPos.ByteOffset : e.Location.EndPos.ByteOffset+1])
	var currentSourceLength int
	var currentErrorLength int
	var ellipsisStart bool
	var ellipsisEnd bool
	sourceFragmentStartIndex := e.Location.StartPos.ByteOffset
	sourceFragmentEndIndex := e.Location.EndPos.ByteOffset + 1

	if errorSourceLength < maxSourceExampleLength {
		leftLength := maxSourceExampleLength - errorSourceLength
		beforeSource := source[:e.Location.StartPos.ByteOffset]
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
		s := source[sourceFragmentStartIndex:e.Location.StartPos.ByteOffset]
		startOffset += utf8.RuneCountInString(s)
	}

	exampleEnd := e.Location.EndPos.ByteOffset
	if lineEndIndex < exampleEnd {
		exampleEnd = lineEndIndex
	}
	for i := range source[e.Location.StartPos.ByteOffset : exampleEnd+1] {
		if currentSourceLength >= maxSourceExampleLength {
			if i < lineEndIndex-1 {
				ellipsisEnd = true
			}
			break
		}
		currentSourceLength++
		currentErrorLength++
		sourceFragmentEndIndex = e.Location.StartPos.ByteOffset + i
	}
	for i := range source[exampleEnd:lineEndIndex] {
		if currentSourceLength >= maxSourceExampleLength {
			if i < lineEndIndex-1 {
				ellipsisEnd = true
			}
			break
		}

		currentSourceLength++
		sourceFragmentEndIndex = e.Location.EndPos.ByteOffset + i
	}
	if ellipsisStart {
		result.WriteString(faintColor.Sprint(ellipsis))
	}
	sourceFragment := source[sourceFragmentStartIndex : sourceFragmentEndIndex+1]
	if style {
		sourceFragment = lexer.Colorize(sourceFragment)
	}
	result.WriteString(sourceFragment)
	if ellipsisEnd {
		result.WriteString(faintColor.Sprint(ellipsis))
	}
	result.WriteByte('\n')
	result.WriteString(strings.Repeat(" ", startOffset))

	if currentErrorLength <= 1 {
		result.WriteString(severityColor.Sprint("│"))
	} else {
		result.WriteString(severityColor.Sprint("└"))
		result.WriteString(severityColor.Sprint(strings.Repeat("─", currentErrorLength-2)))
		result.WriteString(severityColor.Sprint("┤"))
	}
	result.WriteByte('\n')
	var spaceCount int
	if currentErrorLength == 0 {
		spaceCount = startOffset
	} else {
		spaceCount = startOffset + currentErrorLength - 1
	}
	result.WriteString(strings.Repeat(" ", spaceCount))
	result.WriteString(severityColor.Sprint("└ Here\n"))
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

// Add a new error.
func (e *ErrorList) Append(err *Error) {
	*e = append(*e, err)
}

// Create and add a new error.
func (e *ErrorList) Add(message string, loc *position.Location, severity Severity) {
	e.Append(NewError(loc, message, severity))
}

// Create and add a new failure.
func (e *ErrorList) AddFailure(message string, loc *position.Location) {
	e.Append(NewError(loc, message, FAILURE))
}

// Create and add a new warning.
func (e *ErrorList) AddWarning(message string, loc *position.Location) {
	e.Append(NewError(loc, message, WARNING))
}

// Implements the error interface.
func (e ErrorList) Error() string {
	switch len(e) {
	case 0:
		return "<empty ErrorList>"
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
	maxSourceExampleLength = 80
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
		result.WriteRune('\n')
	}
	return result.String(), nil
}

// Return a string representation of this error list
// that can be presented to humans.
func (e ErrorList) HumanStringWithSource(source string, style bool) string {
	var result strings.Builder
	for _, err := range e {
		result.WriteString(err.HumanStringWithSource(source, style))
		result.WriteRune('\n')
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

func (e ErrorList) IsFailure() bool {
	for _, err := range e {
		if err.Severity == FAILURE {
			return true
		}
	}

	return false
}

// A thread-safe list of errors.
type SyncErrorList struct {
	ErrorList ErrorList
	mutex     sync.Mutex
}

// Create and add a new error.
func (e *SyncErrorList) Add(message string, loc *position.Location, severity Severity) {
	e.Append(NewError(loc, message, severity))
}

// Create and add a new failure.
func (e *SyncErrorList) AddFailure(message string, loc *position.Location) {
	e.Append(NewFailure(loc, message))
}

// Create and add a new warning.
func (e *SyncErrorList) AddWarning(message string, loc *position.Location) {
	e.Append(NewWarning(loc, message))
}

// Add a new error.
func (e *SyncErrorList) Append(err *Error) {
	e.mutex.Lock()
	e.ErrorList = append(e.ErrorList, err)
	e.mutex.Unlock()
}

func (e *SyncErrorList) Join(other *SyncErrorList) {
	e.mutex.Lock()
	other.mutex.Lock()

	e.ErrorList = append(e.ErrorList, other.ErrorList...)

	other.mutex.Unlock()
	e.mutex.Unlock()
}

func (e *SyncErrorList) JoinErrList(other ErrorList) {
	e.mutex.Lock()
	e.ErrorList = append(e.ErrorList, other...)
	e.mutex.Unlock()
}

func (e *SyncErrorList) IsFailure() bool {
	return e.ErrorList.IsFailure()
}

func (e *SyncErrorList) Error() string {
	return e.ErrorList.Error()
}
