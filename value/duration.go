package value

import (
	"fmt"
	"math/big"
	"time"
)

// Represents the elapsed time between two Times as an int64 nanosecond count.
// The representation limits the largest representable duration to approximately 290 years.
// Wraps Go's time.Duration.
type Duration time.Duration

const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
	Day                  = 24 * Hour
	Week                 = 7 * Day
	Year                 = Duration(365.25 * float64(Week))
)

var DurationClass *Class // ::Std::Duration

// Parses a duration string and creates a Duration value.
// A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func ParseDuration(s String) (result Duration, err Value) {
	dur, er := time.ParseDuration(s.String())
	if er != nil {
		return result, Ref(NewError(FormatErrorClass, er.Error()))
	}

	return Duration(dur), Undefined
}

func DurationSince(t Time) Duration {
	return Duration(time.Since(t.Go))
}

func DurationUntil(t Time) Duration {
	return Duration(time.Until(t.Go))
}

func (d Duration) Go() time.Duration {
	return time.Duration(d)
}

func (d Duration) Copy() Reference {
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
	return fmt.Sprintf("Std::Duration('%s')", d.String())
}

func (Duration) InstanceVariables() *InstanceVariables {
	return nil
}

func (d Duration) Error() string {
	return d.Inspect()
}

func (d Duration) String() string {
	return d.Go().String()
}

func (d Duration) ToString() String {
	return String(d.String())
}

func (d Duration) Add(other Duration) Duration {
	return Duration(d.Go() + other.Go())
}

func (d Duration) Subtract(other Duration) Duration {
	return Duration(d.Go() - other.Go())
}

func (d Duration) Multiply(other Value) (Duration, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			newBig := big.NewInt(int64(d))
			result := ToElkBigInt(newBig.Mul(newBig, o.ToGoBigInt()))
			return Duration(result.ToSmallInt()), Undefined
		case *BigFloat:
			prec := max(o.Precision(), 64)
			iBigFloat := (&BigFloat{}).SetPrecision(prec).SetInt64(Int64(d))
			iBigFloat.MulBigFloat(iBigFloat, o)
			return Duration(iBigFloat.ToInt64()), Undefined
		default:
			return 0, Ref(NewCoerceError(d.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return d * Duration(other.AsSmallInt()), Undefined
	case FLOAT_FLAG:
		return Duration(Float(d) * other.AsFloat()), Undefined
	default:
		return 0, Ref(NewCoerceError(d.Class(), other.Class()))
	}
}

func (d Duration) Divide(other Value) (Duration, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			if o.IsZero() {
				return 0, Ref(NewZeroDivisionError())
			}
			newBig := big.NewInt(int64(d))
			result := ToElkBigInt(newBig.Div(newBig, o.ToGoBigInt()))
			return Duration(result.ToSmallInt()), Undefined
		case *BigFloat:
			if o.IsZero() {
				return 0, Ref(NewZeroDivisionError())
			}
			prec := max(o.Precision(), 64)
			iBigFloat := (&BigFloat{}).SetPrecision(prec).SetInt64(Int64(d))
			iBigFloat.DivBigFloat(iBigFloat, o)
			return Duration(iBigFloat.ToInt64()), Undefined
		default:
			return 0, Ref(NewCoerceError(d.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		if o == 0 {
			return 0, Ref(NewZeroDivisionError())
		}
		return d / Duration(o), Undefined
	case FLOAT_FLAG:
		o := other.AsFloat()
		if o == 0 {
			return 0, Ref(NewZeroDivisionError())
		}
		return Duration(Float(d) / o), Undefined
	default:
		return 0, Ref(NewCoerceError(d.Class(), other.Class()))
	}
}

func (d Duration) Nanoseconds() Value {
	return ToElkInt(d.Go().Nanoseconds())
}

func (d Duration) InNanoseconds() Float {
	return Float(d.Go().Nanoseconds())
}

func (d Duration) Microseconds() Value {
	return ToElkInt(d.Go().Microseconds())
}

func (d Duration) InMicroseconds() Float {
	return Float(d.Go().Microseconds())
}

func (d Duration) Milliseconds() Value {
	return ToElkInt(d.Go().Milliseconds())
}

func (d Duration) InMilliseconds() Float {
	return Float(d.Go().Milliseconds())
}

func (d Duration) Seconds() Value {
	return ToElkInt(int64(d / Second))
}

func (d Duration) InSeconds() Float {
	return Float(d.Go().Seconds())
}

func (d Duration) Minutes() Value {
	return ToElkInt(int64(d / Minute))
}

func (d Duration) InMinutes() Float {
	return Float(d.Go().Minutes())
}

func (d Duration) Hours() Value {
	return ToElkInt(int64(d / Hour))
}

func (d Duration) InHours() Float {
	return Float(d.Go().Hours())
}

func (d Duration) Days() Value {
	return ToElkInt(int64(d / Day))
}

func (d Duration) InDays() Float {
	day := d / Day
	nsec := d % Day
	return Float(day) + Float(nsec)/(24*60*60*1e9)
}

func (d Duration) Weeks() Value {
	return ToElkInt(int64(d / Week))
}

func (d Duration) InWeeks() Float {
	week := d / Week
	nsec := d % Week
	return Float(week) + Float(nsec)/(7*24*60*60*1e9)
}

func (d Duration) Years() Value {
	return ToElkInt(int64(d / Year))
}

func (d Duration) InYears() Float {
	week := d / Year
	nsec := d % Year
	return Float(week) + Float(nsec)/(365.25*7*24*60*60*1e9)
}

func initDuration() {
	DurationClass = NewClassWithOptions(ClassWithParent(ValueClass))
	StdModule.AddConstantString("Duration", Ref(DurationClass))
}
