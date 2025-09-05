package value

import (
	"fmt"
	"strings"
)

// Represents the difference between two datetimes.
// It is made up of a `Date::Span` and a `Time::Span`
type DateTimeSpan struct {
	DateSpan DateSpan
	TimeSpan TimeSpan
}

func NewDateTimeSpan(datespan DateSpan, timespan TimeSpan) *DateTimeSpan {
	span := &DateTimeSpan{
		DateSpan: datespan,
		TimeSpan: timespan,
	}
	span.Normalise()
	return span
}

var DateTimeSpanClass *Class // ::Std::DateTime::Span

func (d *DateTimeSpan) Normalise() {
	days := d.TimeSpan / Day
	if days != 0 {
		d.DateSpan.days += int32(days)
		d.TimeSpan %= Day
	}

	if d.DateSpan.days > 0 && d.TimeSpan < 0 {
		d.DateSpan.days -= 1
		d.TimeSpan += Day
	} else if d.DateSpan.days < 0 && d.TimeSpan > 0 {
		d.DateSpan.days += 1
		d.TimeSpan -= Day
	}
}

func (d *DateTimeSpan) Copy() Reference {
	return &DateTimeSpan{
		DateSpan: d.DateSpan,
		TimeSpan: d.TimeSpan,
	}
}

func (*DateTimeSpan) Class() *Class {
	return DateTimeSpanClass
}

func (*DateTimeSpan) DirectClass() *Class {
	return DateTimeSpanClass
}

func (*DateTimeSpan) SingletonClass() *Class {
	return nil
}

func (d *DateTimeSpan) Inspect() string {
	return fmt.Sprintf("Std::DateTime::Span.parse('%s')", d.String())
}

func (*DateTimeSpan) InstanceVariables() *InstanceVariables {
	return nil
}

func (d *DateTimeSpan) Error() string {
	return d.Inspect()
}

func (d *DateTimeSpan) IsZero() bool {
	return d.DateSpan.IsZero() && d.TimeSpan == 0
}

func (d *DateTimeSpan) String() string {
	if d.IsZero() {
		return "0D"
	}

	var buff strings.Builder
	if !d.DateSpan.IsZero() {
		buff.WriteString(d.DateSpan.String())
	}
	if d.TimeSpan != 0 {
		buff.WriteString(d.TimeSpan.String())
	}

	return buff.String()
}

func (d *DateTimeSpan) ToString() String {
	return String(d.String())
}

func (d *DateTimeSpan) Negate() *DateTimeSpan {
	return NewDateTimeSpan(
		d.DateSpan.Negate(),
		-d.TimeSpan,
	)
}

func (d *DateTimeSpan) Nanoseconds() int {
	return d.TimeSpan.Nanoseconds()
}

func (d *DateTimeSpan) TotalNanoseconds() Value {
	result, err := AddVal(d.DateSpan.TotalNanoseconds(), d.TimeSpan.TotalNanoseconds())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not add total nanoseconds in DateTime::Span: %s", err.Error()))
	}

	return result
}

func (d *DateTimeSpan) InNanoseconds() Float {
	return d.DateSpan.InNanoseconds() + d.TimeSpan.InNanoseconds()
}

func (d *DateTimeSpan) Microseconds() int {
	return d.TimeSpan.Microseconds()
}

func (d *DateTimeSpan) TotalMicroseconds() Value {
	result, err := AddVal(d.DateSpan.TotalMicroseconds(), d.TimeSpan.TotalMicroseconds())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not add total microseconds in DateTime::Span: %s", err.Error()))
	}

	return result
}

func (d *DateTimeSpan) InMicroseconds() Float {
	return d.DateSpan.InMicroseconds() + d.TimeSpan.InMicroseconds()
}

func (d *DateTimeSpan) Milliseconds() int {
	return d.TimeSpan.Milliseconds()
}

func (d *DateTimeSpan) TotalMilliseconds() Value {
	result, err := AddVal(d.DateSpan.TotalMilliseconds(), d.TimeSpan.TotalMilliseconds())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not add total milliseconds in DateTime::Span: %s", err.Error()))
	}

	return result
}

func (d *DateTimeSpan) InMilliseconds() Float {
	return d.DateSpan.InMilliseconds() + d.TimeSpan.InMilliseconds()
}

func (d *DateTimeSpan) Seconds() int {
	return d.TimeSpan.Seconds()
}

func (d *DateTimeSpan) TotalSeconds() Value {
	result, err := AddVal(d.DateSpan.TotalSeconds(), d.TimeSpan.TotalSeconds())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not add total seconds in DateTime::Span: %s", err.Error()))
	}

	return result
}

func (d *DateTimeSpan) InSeconds() Float {
	return d.DateSpan.InSeconds() + d.TimeSpan.InSeconds()
}

func (d *DateTimeSpan) Minutes() int {
	return d.TimeSpan.Minutes()
}

func (d *DateTimeSpan) TotalMinutes() Value {
	result, err := AddVal(d.DateSpan.TotalMinutes(), d.TimeSpan.TotalMinutes())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not add total minutes in DateTime::Span: %s", err.Error()))
	}

	return result
}

func (d *DateTimeSpan) InMinutes() Float {
	return d.DateSpan.InMinutes() + d.TimeSpan.InMinutes()
}

func (d *DateTimeSpan) Hours() int {
	return d.TimeSpan.Hours()
}

func (d *DateTimeSpan) TotalHours() Value {
	result, err := AddVal(d.DateSpan.TotalHours(), d.TimeSpan.TotalHours())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not add total hours in DateTime::Span: %s", err.Error()))
	}

	return result
}

func (d *DateTimeSpan) InHours() Float {
	return d.DateSpan.InHours() + d.TimeSpan.InHours()
}

func (d *DateTimeSpan) Days() int {
	return d.DateSpan.Days()
}

func (d *DateTimeSpan) TotalDays() Value {
	result, err := AddVal(d.DateSpan.TotalDays(), d.TimeSpan.TotalDays())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not add total days in DateTime::Span: %s", err.Error()))
	}

	return result
}

func (d *DateTimeSpan) InDays() Float {
	return d.DateSpan.InDays() + d.TimeSpan.InDays()
}

func (d *DateTimeSpan) TotalWeeks() Value {
	result, err := AddVal(d.DateSpan.TotalWeeks(), d.TimeSpan.TotalWeeks())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not add total weeks in DateTime::Span: %s", err.Error()))
	}

	return result
}

func (d *DateTimeSpan) InWeeks() Float {
	return d.DateSpan.InWeeks() + d.TimeSpan.InWeeks()
}

func (d *DateTimeSpan) Months() int {
	return d.DateSpan.Months()
}

func (d *DateTimeSpan) TotalMonths() Value {
	result, err := AddVal(d.DateSpan.TotalMonths(), d.TimeSpan.TotalMonths())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not add total months in DateTime::Span: %s", err.Error()))
	}

	return result
}

func (d *DateTimeSpan) InMonths() Float {
	return d.DateSpan.InMonths() + d.TimeSpan.InMonths()
}

func (d *DateTimeSpan) Years() int {
	return d.DateSpan.Years()
}

func (d *DateTimeSpan) TotalYears() Value {
	result, err := AddVal(d.DateSpan.TotalYears(), d.TimeSpan.TotalYears())
	if !err.IsUndefined() {
		panic(fmt.Sprintf("could not add total years in DateTime::Span: %s", err.Error()))
	}

	return result
}

func (d *DateTimeSpan) InYears() Float {
	return d.DateSpan.InYears() + d.TimeSpan.InYears()
}

func (d *DateTimeSpan) Add(other Value) (Value, Value) {
	switch other.flag {
	case DATE_SPAN_FLAG:
		return Ref(d.AddDateSpan(other.AsDateSpan())), Undefined
	case TIME_SPAN_FLAG:
		return Ref(d.AddTimeSpan(other.AsTimeSpan())), Undefined
	case REFERENCE_FLAG:
	default:
		return Undefined, Ref(NewArgumentTypeError("other", other.Class().Inspect(), durationUnionType))
	}

	switch other := other.AsReference().(type) {
	case TimeSpan:
		return Ref(d.AddTimeSpan(other)), Undefined
	case DateSpan:
		return Ref(d.AddDateSpan(other)), Undefined
	case *DateTimeSpan:
		return Ref(d.AddDateTimeSpan(other)), Undefined
	default:

		return Undefined, Ref(NewArgumentTypeError("other", other.Class().Inspect(), durationUnionType))
	}
}

func (d *DateTimeSpan) AddTimeSpan(other TimeSpan) *DateTimeSpan {
	return NewDateTimeSpan(
		d.DateSpan,
		d.TimeSpan.AddTimeSpan(other),
	)
}

func (d *DateTimeSpan) AddDateSpan(other DateSpan) *DateTimeSpan {
	return NewDateTimeSpan(
		d.DateSpan.AddDateSpan(other),
		d.TimeSpan,
	)
}

func (d *DateTimeSpan) AddDateTimeSpan(other *DateTimeSpan) *DateTimeSpan {
	return NewDateTimeSpan(
		d.DateSpan.AddDateSpan(other.DateSpan),
		d.TimeSpan.AddTimeSpan(other.TimeSpan),
	)
}

func (d *DateTimeSpan) Subtract(other Value) (Value, Value) {
	switch other.flag {
	case DATE_SPAN_FLAG:
		return Ref(d.SubtractDateSpan(other.AsDateSpan())), Undefined
	case TIME_SPAN_FLAG:
		return Ref(d.SubtractTimeSpan(other.AsTimeSpan())), Undefined
	case REFERENCE_FLAG:
	default:
		return Undefined, Ref(NewArgumentTypeError("other", other.Class().Inspect(), durationUnionType))
	}

	switch other := other.AsReference().(type) {
	case TimeSpan:
		return Ref(d.SubtractTimeSpan(other)), Undefined
	case DateSpan:
		return Ref(d.SubtractDateSpan(other)), Undefined
	case *DateTimeSpan:
		return Ref(d.SubtractDateTimeSpan(other)), Undefined
	default:

		return Undefined, Ref(NewArgumentTypeError("other", other.Class().Inspect(), durationUnionType))
	}
}

func (d *DateTimeSpan) SubtractTimeSpan(other TimeSpan) *DateTimeSpan {
	return NewDateTimeSpan(
		d.DateSpan,
		d.TimeSpan.AddTimeSpan(-other),
	)
}

func (d *DateTimeSpan) SubtractDateSpan(other DateSpan) *DateTimeSpan {
	return NewDateTimeSpan(
		d.DateSpan.AddDateSpan(other.Negate()),
		d.TimeSpan,
	)
}

func (d *DateTimeSpan) SubtractDateTimeSpan(other *DateTimeSpan) *DateTimeSpan {
	return NewDateTimeSpan(
		d.DateSpan.AddDateSpan(other.DateSpan.Negate()),
		d.TimeSpan.AddTimeSpan(-other.TimeSpan),
	)
}

func initDateTimeSpan() {
	DateTimeSpanClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	DateTimeSpanClass.IncludeMixin(DurationMixin)
	DateTimeClass.AddConstantString("Span", Ref(DateTimeSpanClass))
}
