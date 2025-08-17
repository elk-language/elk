package diagnostic

import (
	"fmt"
	"strings"
	"sync"

	"github.com/elk-language/elk/indent"
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

func writeDiagnosticLocationToBuffer(location *position.Location, buffer *strings.Builder, style bool, colorizer colorizer.Colorizer, accentColor *color.Color, sourceMap map[string]string) error {
	result, err := location.HumanStringWithSourceMap(style, colorizer, accentColor, sourceMap)
	if err != nil {
		return err
	}

	buffer.WriteString("\n\n")
	indent.IndentString(buffer, result, 1)

	if location.Parent != nil {
		writeDiagnosticLocationToBuffer(location.Parent, buffer, style, colorizer, accentColor, sourceMap)
	}

	return nil
}

func (d *Diagnostic) SeverityColor() *color.Color {
	return color.New(color.Bold, d.Severity.color())
}

// Return a string representation of this error
// that can be presented to humans.
func (d *Diagnostic) HumanString(style bool, colorizer colorizer.Colorizer) (string, error) {
	return d.HumanStringWithSourceMap(style, colorizer, nil)
}

// Return a string representation of this error
// that can be presented to humans.
func (d *Diagnostic) HumanStringWithSourceMap(style bool, colorizer colorizer.Colorizer, sourceMap map[string]string) (string, error) {
	var buffer strings.Builder
	severityColor := d.SeverityColor()
	if !style {
		severityColor.DisableColor()
	}

	buffer.WriteString(d.Severity.Tag(style))
	buffer.WriteByte(' ')
	buffer.WriteString(d.Message)

	err := writeDiagnosticLocationToBuffer(d.Location, &buffer, style, colorizer, severityColor, sourceMap)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
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
func (dl *DiagnosticList) Append(d ...*Diagnostic) {
	*dl = append(*dl, d...)
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

// Return a string representation of this diagnostic list
// that can be presented to humans.
func (dl DiagnosticList) HumanString(style bool, colorizer colorizer.Colorizer) (string, error) {
	var buffer strings.Builder
	for _, d := range dl {
		msg, err := d.HumanString(style, colorizer)
		if err != nil {
			return "", err
		}
		buffer.WriteString(msg)
		buffer.WriteString("\n\n")
	}
	return buffer.String(), nil
}

// Return a string representation of this diagnostic list
// that can be presented to humans.
func (dl DiagnosticList) HumanStringWithSourceMap(style bool, colorizer colorizer.Colorizer, sourceMap map[string]string) (string, error) {
	var buffer strings.Builder
	for _, d := range dl {
		msg, err := d.HumanStringWithSourceMap(style, colorizer, sourceMap)
		if err != nil {
			return "", err
		}
		buffer.WriteString(msg)
		buffer.WriteString("\n\n")
	}
	return buffer.String(), nil
}

// Err returns an error equivalent to this diagnostic list.
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
