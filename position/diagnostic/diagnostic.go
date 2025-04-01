package diagnostic

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/elk-language/elk/lexer/colorizer"
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
	FAIL Severity = iota
	WARN
	INFO
)

var severityNames = []string{
	FAIL: "FAIL",
	WARN: "WARN",
	INFO: "INFO",
}

var severityColor = []color.Attribute{
	FAIL: color.FgRed,
	WARN: color.FgYellow,
	INFO: color.FgBlue,
}

// Represents a single diagnostic
// at a particular source location (an error, a warning, an info message)
type Diagnostic struct {
	*position.Location
	Severity Severity
	Message  string
}

// Create a new diagnostic.
func NewDiagnostic(loc *position.Location, msg string, severity Severity) *Diagnostic {
	return &Diagnostic{
		Location: loc,
		Message:  msg,
		Severity: severity,
	}
}

// Create a new info message.
func NewInfo(loc *position.Location, msg string) *Diagnostic {
	return &Diagnostic{
		Location: loc,
		Message:  msg,
		Severity: INFO,
	}
}

// Create a new failure.
func NewFailure(loc *position.Location, msg string) *Diagnostic {
	return &Diagnostic{
		Location: loc,
		Message:  msg,
		Severity: FAIL,
	}
}

// Create a new warning.
func NewWarning(loc *position.Location, msg string) *Diagnostic {
	return &Diagnostic{
		Location: loc,
		Message:  msg,
		Severity: WARN,
	}
}

// Implements the error interface.
func (e *Diagnostic) Error() string {
	return e.String()
}

// Implements the fmt.Stringer interface
func (e *Diagnostic) String() string {
	return fmt.Sprintf("%s: %s", e.Location.String(), e.Message)
}

// Return a string representation of this error
// that can be presented to humans.
func (e *Diagnostic) HumanString(style bool, colorizer colorizer.Colorizer) (string, error) {
	source, err := os.ReadFile(e.Location.FilePath)
	if err != nil {
		return "", err
	}
	return e.HumanStringWithSource(string(source), style, colorizer)
}

// Return a string representation of this error
// that can be presented to humans.
func (e *Diagnostic) HumanStringWithSource(source string, style bool, colorizer colorizer.Colorizer) (string, error) {
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
		return result.String(), nil
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
	return result.String(), nil
}

// DiagnosticList is a list of *Errors.
// The zero value for an DiagnosticList is an empty DiagnosticList ready to use.
type DiagnosticList []*Diagnostic

// Create a new slice containing the elements of the
// two given diagnostic lists.
func (e DiagnosticList) Join(other DiagnosticList) DiagnosticList {
	n := len(e)
	return append(e[:n:n], other...)
}

// Add a new diagnostic.
func (dl *DiagnosticList) Append(d *Diagnostic) {
	*dl = append(*dl, d)
}

// Create and add a new diagnostic.
func (e *DiagnosticList) Add(message string, loc *position.Location, severity Severity) {
	e.Append(NewDiagnostic(loc, message, severity))
}

// Create and add a new failure.
func (e *DiagnosticList) AddFailure(message string, loc *position.Location) {
	e.Append(NewDiagnostic(loc, message, FAIL))
}

// Create and add a new warning.
func (e *DiagnosticList) AddWarning(message string, loc *position.Location) {
	e.Append(NewDiagnostic(loc, message, WARN))
}

// Create and add a new info.
func (e *DiagnosticList) AddInfo(message string, loc *position.Location) {
	e.Append(NewDiagnostic(loc, message, INFO))
}

// Implements the diagnostic interface.
func (e DiagnosticList) Error() string {
	switch len(e) {
	case 0:
		return "<empty DiagnosticList>"
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
func (el DiagnosticList) HumanString(style bool, colorizer colorizer.Colorizer) (string, error) {
	var result strings.Builder
	for _, e := range el {
		msg, err := e.HumanString(style, colorizer)
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
func (dl DiagnosticList) HumanStringWithSource(source string, style bool, colorizer colorizer.Colorizer) (string, error) {
	var result strings.Builder
	for _, d := range dl {
		str, err := d.HumanStringWithSource(source, style, colorizer)
		if err != nil {
			return "", err
		}
		result.WriteString(str)
		result.WriteRune('\n')
	}
	return result.String(), nil
}

// Err returns an error equivalent to this error list.
// If the list is empty, Err returns nil.
func (e DiagnosticList) Err() error {
	if len(e) == 0 {
		return nil
	}
	return e
}

func (e DiagnosticList) IsFailure() bool {
	for _, err := range e {
		if err.Severity == FAIL {
			return true
		}
	}

	return false
}

// A thread-safe list of errors.
type SyncDiagnosticList struct {
	DiagnosticList DiagnosticList
	Mutex          sync.Mutex
}

func NewSyncDiagnosticList() *SyncDiagnosticList {
	return &SyncDiagnosticList{}
}

// Create and add a new error.
func (e *SyncDiagnosticList) Add(message string, loc *position.Location, severity Severity) {
	e.Append(NewDiagnostic(loc, message, severity))
}

// Create and add a new failure.
func (e *SyncDiagnosticList) AddFailure(message string, loc *position.Location) {
	e.Append(NewFailure(loc, message))
}

// Create and add a new warning.
func (e *SyncDiagnosticList) AddWarning(message string, loc *position.Location) {
	e.Append(NewWarning(loc, message))
}

// Create and add a new warning.
func (e *SyncDiagnosticList) AddInfo(message string, loc *position.Location) {
	e.Append(NewInfo(loc, message))
}

// Add a new diagnostic.
func (e *SyncDiagnosticList) Append(err *Diagnostic) {
	e.Mutex.Lock()
	e.DiagnosticList = append(e.DiagnosticList, err)
	e.Mutex.Unlock()
}

func (e *SyncDiagnosticList) Join(other *SyncDiagnosticList) {
	e.Mutex.Lock()
	other.Mutex.Lock()

	e.DiagnosticList = append(e.DiagnosticList, other.DiagnosticList...)

	other.Mutex.Unlock()
	e.Mutex.Unlock()
}

func (e *SyncDiagnosticList) JoinErrList(other DiagnosticList) {
	e.Mutex.Lock()
	e.DiagnosticList = append(e.DiagnosticList, other...)
	e.Mutex.Unlock()
}

func (e *SyncDiagnosticList) IsFailure() bool {
	return e.DiagnosticList.IsFailure()
}

func (e *SyncDiagnosticList) Error() string {
	return e.DiagnosticList.Error()
}
