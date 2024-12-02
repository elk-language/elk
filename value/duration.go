package value

import (
	"fmt"
	"time"
)

// Represents the elapsed time between two Times as an int64 nanosecond count.
// The representation limits the largest representable duration to approximately 290 years.
// Wraps Go's time.Duration.
type Duration struct {
	Go time.Duration
}

var DurationClass *Class // ::Std::Duration

// Parses a duration string and creates a Duration value.
// A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func ParseDuration(s String) (result Duration, err Value) {
	dur, er := time.ParseDuration(s.String())
	if er != nil {
		return result, NewError(FormatErrorClass, er.Error())
	}

	return ToElkDuration(dur), nil
}

func ToElkDuration(dur time.Duration) Duration {
	return Duration{Go: dur}
}

func (d Duration) Copy() Value {
	return d
}

func (Duration) Class() *Class {
	return DurationClass
}

func (Duration) DirectClass() *Class {
	return DurationClass
}

func (Duration) SingletonClass() *Class {
	return nil
}

func (d Duration) Inspect() string {
	return fmt.Sprintf("Duration('%s')", d.String())
}

func (Duration) InstanceVariables() SymbolMap {
	return nil
}

func (d Duration) Error() string {
	return d.Inspect()
}

func (d Duration) String() string {
	return d.Go.String()
}

func (d Duration) ToString() String {
	return String(d.String())
}

func (d Duration) Add(other Duration) Duration {
	return ToElkDuration(d.Go + other.Go)
}

func (d Duration) Subtract(other Duration) Duration {
	return ToElkDuration(d.Go - other.Go)
}

func (d Duration) Nanoseconds() Int64 {
	return Int64(d.Go.Nanoseconds())
}

func (d Duration) Microseconds() Int64 {
	return Int64(d.Go.Microseconds())
}

func (d Duration) Milliseconds() Int64 {
	return Int64(d.Go.Milliseconds())
}

func (d Duration) Seconds() Float64 {
	return Float64(d.Go.Seconds())
}

func (d Duration) Minutes() Float64 {
	return Float64(d.Go.Minutes())
}

func (d Duration) Hours() Float64 {
	return Float64(d.Go.Hours())
}

func initDuration() {
	DurationClass = NewClass()
	StdModule.AddConstantString("Duration", DurationClass)
}
