package value

import (
	"fmt"
	"math"
	"math/big"
	"time"
)

// Represents the elapsed time between two Times as an int64 nanosecond count.
// The representation limits the largest representable duration to approximately 290 years.
// Wraps Go's time.Duration.
type TimeSpan time.Duration

const MonthDays = 30.4375
const YearDays = 365.25

const (
	Nanosecond  TimeSpan = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
	Day                  = 24 * Hour
	Week                 = 7 * Day
	Month                = TimeSpan(MonthDays * float64(Day))
	Year                 = TimeSpan(YearDays * float64(Day))
)

var TimeSpanClass *Class // ::Std::Time::Span

// Parses a time span string and creates a time span value.
// A time span string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func ParsTimeSpan(s String) (result TimeSpan, err Value) {
	dur, er := time.ParseDuration(s.String())
	if er != nil {
		return result, Ref(NewError(FormatErrorClass, er.Error()))
	}

	return TimeSpan(dur), Undefined
}

func DurationSince(t DateTime) TimeSpan {
	return TimeSpan(time.Since(t.Go))
}

func DurationUntil(t DateTime) TimeSpan {
	return TimeSpan(time.Until(t.Go))
}

func (t TimeSpan) Go() time.Duration {
	return time.Duration(t)
}

func (t TimeSpan) Copy() Reference {
	return t
}

func (TimeSpan) Class() *Class {
	return TimeSpanClass
}

func (TimeSpan) DirectClass() *Class {
	return TimeSpanClass
}

func (TimeSpan) SingletonClass() *Class {
	return nil
}

func (t TimeSpan) Inspect() string {
	return fmt.Sprintf("Std::Time::Span('%s')", t.String())
}

func (TimeSpan) InstanceVariables() *InstanceVariables {
	return nil
}

func (t TimeSpan) Error() string {
	return t.Inspect()
}

func (t TimeSpan) String() string {
	return t.Go().String()
}

func (t TimeSpan) ToString() String {
	return String(t.String())
}

func (t TimeSpan) Add(other Value) (Value, Value) {
	switch other.flag {
	case DATE_SPAN_FLAG:
		return Ref(t.AddDateSpan(other.AsDateSpan())), Undefined
	case TIME_SPAN_FLAG:
		return t.AddTimeSpan(other.AsTimeSpan()).ToValue(), Undefined
	case REFERENCE_FLAG:
	default:
		return Undefined, Ref(NewArgumentTypeError("other", other.Class().Inspect(), durationUnionType))
	}

	switch other := other.AsReference().(type) {
	case TimeSpan:
		return t.AddTimeSpan(other).ToValue(), Undefined
	case *DateTimeSpan:
		return Ref(t.AddDateTimeSpan(other)), Undefined
	default:

		return Undefined, Ref(NewArgumentTypeError("other", other.Class().Inspect(), durationUnionType))
	}
}

func (t TimeSpan) AddTimeSpan(other TimeSpan) TimeSpan {
	return TimeSpan(t.Go() + other.Go())
}

func (t TimeSpan) AddDateSpan(other DateSpan) *DateTimeSpan {
	return NewDateTimeSpan(other, t)
}

func (t TimeSpan) AddDateTimeSpan(other *DateTimeSpan) *DateTimeSpan {
	return NewDateTimeSpan(
		other.DateSpan,
		t.AddTimeSpan(other.TimeSpan),
	)
}

func (t TimeSpan) Subtract(other TimeSpan) TimeSpan {
	return TimeSpan(t.Go() - other.Go())
}

func (t TimeSpan) Multiply(other Value) (TimeSpan, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			newBig := big.NewInt(int64(t))
			result := ToElkBigInt(newBig.Mul(newBig, o.ToGoBigInt()))
			return TimeSpan(result.ToSmallInt()), Undefined
		case *BigFloat:
			prec := max(o.Precision(), 64)
			iBigFloat := (&BigFloat{}).SetPrecision(prec).SetInt64(Int64(t))
			iBigFloat.MulBigFloat(iBigFloat, o)
			return TimeSpan(iBigFloat.ToInt64()), Undefined
		default:
			return 0, Ref(NewCoerceError(t.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return t * TimeSpan(other.AsSmallInt()), Undefined
	case FLOAT_FLAG:
		return TimeSpan(Float(t) * other.AsFloat()), Undefined
	default:
		return 0, Ref(NewCoerceError(t.Class(), other.Class()))
	}
}

func (t TimeSpan) Divide(other Value) (TimeSpan, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			if o.IsZero() {
				return 0, Ref(NewZeroDivisionError())
			}
			newBig := big.NewInt(int64(t))
			result := ToElkBigInt(newBig.Div(newBig, o.ToGoBigInt()))
			return TimeSpan(result.ToSmallInt()), Undefined
		case *BigFloat:
			if o.IsZero() {
				return 0, Ref(NewZeroDivisionError())
			}
			prec := max(o.Precision(), 64)
			iBigFloat := (&BigFloat{}).SetPrecision(prec).SetInt64(Int64(t))
			iBigFloat.DivBigFloat(iBigFloat, o)
			return TimeSpan(iBigFloat.ToInt64()), Undefined
		default:
			return 0, Ref(NewCoerceError(t.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		if o == 0 {
			return 0, Ref(NewZeroDivisionError())
		}
		return t / TimeSpan(o), Undefined
	case FLOAT_FLAG:
		o := other.AsFloat()
		if o == 0 {
			return 0, Ref(NewZeroDivisionError())
		}
		return TimeSpan(Float(t) / o), Undefined
	default:
		return 0, Ref(NewCoerceError(t.Class(), other.Class()))
	}
}

func (t TimeSpan) TotalNanoseconds() Value {
	return ToElkInt(t.Go().Nanoseconds())
}

func (t TimeSpan) Nanoseconds() int {
	return int(t.Go().Nanoseconds() % 1000)
}

func (t TimeSpan) InNanoseconds() Float {
	return Float(t.Go().Nanoseconds())
}

func (t TimeSpan) TotalMicroseconds() Value {
	return ToElkInt(t.Go().Microseconds())
}

func (t TimeSpan) Microseconds() int {
	return int(t.Go().Microseconds() % 1000)
}

func (t TimeSpan) InMicroseconds() Float {
	return Float(t.Go().Microseconds())
}

func (t TimeSpan) TotalMilliseconds() Value {
	return ToElkInt(t.Go().Milliseconds())
}

func (t TimeSpan) Milliseconds() int {
	return int(t.Go().Milliseconds() % 1000)
}

func (t TimeSpan) InMilliseconds() Float {
	return Float(t.Go().Milliseconds())
}

func (t TimeSpan) TotalSeconds() Value {
	return ToElkInt(int64(t / Second))
}

func (t TimeSpan) Seconds() int {
	return int(int64(t/Second) % 60)
}

func (t TimeSpan) InSeconds() Float {
	return Float(t.Go().Seconds())
}

func (t TimeSpan) TotalMinutes() Value {
	return ToElkInt(int64(t / Minute))
}

func (t TimeSpan) Minutes() int {
	return int(t / Minute % 60)
}

func (t TimeSpan) InMinutes() Float {
	return Float(t.Go().Minutes())
}

func (t TimeSpan) TotalHours() Value {
	return ToElkInt(int64(t / Hour))
}

func (t TimeSpan) Hours() int {
	return int(t / Hour % 24)
}

func (t TimeSpan) InHours() Float {
	return Float(t.Go().Hours())
}

func (t TimeSpan) TotalDays() Value {
	return ToElkInt(int64(t / Day))
}

func (t TimeSpan) InDays() Float {
	day := t / Day
	nsec := t % Day
	return Float(day) + Float(nsec)/(24*60*60*1e9)
}

func (t TimeSpan) Days() int {
	return int(math.Mod(float64(t)/float64(Day), MonthDays))
}

func (t TimeSpan) TotalWeeks() Value {
	return ToElkInt(int64(t / Week))
}

func (t TimeSpan) InWeeks() Float {
	week := t / Week
	nsec := t % Week
	return Float(week) + Float(nsec)/(7*24*60*60*1e9)
}

func (t TimeSpan) TotalMonths() Value {
	return ToElkInt(int64(t / Month))
}

func (t TimeSpan) Months() int {
	return int(t / Month % 12)
}

func (t TimeSpan) InMonths() Float {
	day := t / Day
	nsec := t % Day
	return Float(day) + Float(nsec)/(MonthDays*24*60*60*1e9)
}

func (t TimeSpan) TotalYears() Value {
	return ToElkInt(int64(t / Year))
}

func (t TimeSpan) InYears() Float {
	week := t / Year
	nsec := t % Year
	return Float(week) + Float(nsec)/(YearDays*7*24*60*60*1e9)
}

func initTimeSpan() {
	TimeSpanClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	TimeSpanClass.IncludeMixin(DurationMixin)
	TimeClass.AddConstantString("Span", Ref(TimeSpanClass))
}
