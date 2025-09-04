package value

import (
	"fmt"
	"strings"
)

// Represents the difference between two dates
// as a 32 bit number of months and 32 bit number of days.
// Can store up to 2 billion days and 2 billion months (178 million years).
type DateSpan struct {
	months int32
	days   int32
}

func MakeDateSpan(years, months, days int) DateSpan {
	months += years * 12

	return DateSpan{
		months: int32(months),
		days:   int32(days),
	}
}

var DateSpanClass *Class // ::Std::Date::Span

func (d DateSpan) IsZero() bool {
	return d.months == 0 && d.days == 0
}

func (d DateSpan) Copy() Reference {
	return d
}

func (DateSpan) Class() *Class {
	return DateSpanClass
}

func (DateSpan) DirectClass() *Class {
	return DateSpanClass
}

func (DateSpan) SingletonClass() *Class {
	return nil
}

func (d DateSpan) Inspect() string {
	return fmt.Sprintf("Std::Date::Span(%d, %d, %d)", d.Years(), d.Months(), d.Days())
}

func (DateSpan) InstanceVariables() *InstanceVariables {
	return nil
}

func (d DateSpan) Error() string {
	return d.Inspect()
}

func (d DateSpan) String() string {
	years := d.months / 12
	months := d.Months()
	days := d.days

	var buff strings.Builder
	if years > 0 {
		fmt.Fprintf(&buff, "%dY", years)
	}
	if months > 0 {
		fmt.Fprintf(&buff, "%dM", months)
	}
	if days > 0 || d.IsZero() {
		fmt.Fprintf(&buff, "%dD", days)
	}
	return buff.String()
}

func (d DateSpan) ToString() String {
	return String(d.String())
}

func (d DateSpan) TotalNanoseconds() Value {
	result, err := MultiplyVal(d.TotalMicroseconds(), SmallInt(1000).ToValue())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not calculate total nanoseconds in Date::Span: %s", err.Error()))
	}

	return result
}

func (d DateSpan) InNanoseconds() Float {
	return d.InMicroseconds() * 1000
}

func (d DateSpan) TotalMicroseconds() Value {
	result, err := MultiplyVal(d.TotalMilliseconds(), SmallInt(1000).ToValue())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not calculate total microseconds in Date::Span: %s", err.Error()))
	}

	return result
}

func (d DateSpan) InMicroseconds() Float {
	return d.InMilliseconds() * 1000
}

func (d DateSpan) TotalMilliseconds() Value {
	result, err := MultiplyVal(d.TotalSeconds(), SmallInt(1000).ToValue())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not calculate total milliseconds in Date::Span: %s", err.Error()))
	}

	return result
}

func (d DateSpan) InMilliseconds() Float {
	return d.InSeconds() * 1000
}

func (d DateSpan) TotalSeconds() Value {
	result, err := MultiplyVal(d.TotalMinutes(), SmallInt(60).ToValue())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not calculate total seconds in Date::Span: %s", err.Error()))
	}

	return result
}

func (d DateSpan) InSeconds() Float {
	return d.InMinutes() * 60
}

func (d DateSpan) TotalMinutes() Value {
	result, err := MultiplyVal(d.TotalHours(), SmallInt(60).ToValue())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not calculate total minutes in Date::Span: %s", err.Error()))
	}

	return result
}

func (d DateSpan) InMinutes() Float {
	return d.InHours() * 60
}

func (d DateSpan) TotalHours() Value {
	result, err := MultiplyVal(d.TotalDays(), SmallInt(24).ToValue())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not calculate total hours in Date::Span: %s", err.Error()))
	}

	return result
}

func (d DateSpan) InHours() Float {
	return d.InDays() * 24
}

func (d DateSpan) Days() int {
	return int(d.days)
}

func (d DateSpan) TotalDays() Value {
	return ToElkInt(int64(float64(d.months)*MonthDays) + int64(d.days))
}

func (d DateSpan) InDays() Float {
	return Float(d.months)*MonthDays + Float(d.days)
}

func (d DateSpan) Months() int {
	return int(d.months % 12)
}

func (d DateSpan) TotalWeeks() Value {
	return ToElkInt(int64((float64(d.months)*MonthDays + float64(d.days)) / 7))
}

func (d DateSpan) InWeeks() Float {
	return (Float(d.months)*MonthDays + Float(d.days)) / 7
}

func (d DateSpan) TotalMonths() Value {
	return ToElkInt(int64(d.months) + int64(float64(d.days)/MonthDays))
}

func (d DateSpan) InMonths() Float {
	return Float(d.months) + Float(d.days)/MonthDays
}

func (d DateSpan) Years() int {
	return int(d.months / 12)
}

func (d DateSpan) TotalYears() Value {
	return ToElkInt(int64(d.months/12) + int64(float64(d.days)/YearDays))
}

func (d DateSpan) InYears() Float {
	return Float(d.months)/12 + Float(d.days)/YearDays
}

const durationUnionType = "Std::Time::Span | Std::Date::Span | Std::DateTime::Span"

func (d DateSpan) Add(other Value) (Value, Value) {
	switch other.flag {
	case DATE_SPAN_FLAG:
		return d.AddDateSpan(other.AsDateSpan()).ToValue(), Undefined
	case TIME_SPAN_FLAG:
		return Ref(d.AddTimeSpan(other.AsTimeSpan())), Undefined
	case REFERENCE_FLAG:
	default:
		return Undefined, Ref(NewArgumentTypeError("other", other.Class().Inspect(), durationUnionType))
	}

	switch other := other.AsReference().(type) {
	case TimeSpan:
		return Ref(d.AddTimeSpan(other)), Undefined
	case *DateTimeSpan:
		return Ref(d.AddDateTimeSpan(other)), Undefined
	default:

		return Undefined, Ref(NewArgumentTypeError("other", other.Class().Inspect(), durationUnionType))
	}
}

func (d DateSpan) AddDateSpan(other DateSpan) DateSpan {
	return DateSpan{
		months: d.months + other.months,
		days:   d.days + other.days,
	}
}

func (d DateSpan) AddTimeSpan(other TimeSpan) *DateTimeSpan {
	return NewDateTimeSpan(
		d,
		other,
	)
}

func (d DateSpan) AddDateTimeSpan(other *DateTimeSpan) *DateTimeSpan {
	return NewDateTimeSpan(
		d.AddDateSpan(other.DateSpan),
		other.TimeSpan,
	)
}

func initDateSpan() {
	DateSpanClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	DateSpanClass.IncludeMixin(DurationMixin)
	DateClass.AddConstantString("Span", Ref(DateSpanClass))
}
