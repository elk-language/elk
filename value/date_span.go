package value

import (
	"fmt"
	"math"
	"math/big"
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
	return fmt.Sprintf("Std::Date::Span.parse(%s)", d.ToString().Inspect())
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
	if years != 0 {
		fmt.Fprintf(&buff, "%dY", years)
	}
	if months != 0 {
		fmt.Fprintf(&buff, "%dM", months)
	}
	if days != 0 || d.IsZero() {
		fmt.Fprintf(&buff, "%dD", days)
	}
	return buff.String()
}

func (d DateSpan) ToString() String {
	return String(d.String())
}

func (d DateSpan) Negate() DateSpan {
	return DateSpan{
		months: -d.months,
		days:   -d.days,
	}
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
	case DateSpan:
		return Ref(d.AddDateSpan(other)), Undefined
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

func (d DateSpan) Subtract(other Value) (Value, Value) {
	switch other.flag {
	case DATE_SPAN_FLAG:
		return d.SubtractDateSpan(other.AsDateSpan()).ToValue(), Undefined
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

func (d DateSpan) SubtractDateSpan(other DateSpan) DateSpan {
	return DateSpan{
		months: d.months - other.months,
		days:   d.days - other.days,
	}
}

func (d DateSpan) SubtractTimeSpan(other TimeSpan) *DateTimeSpan {
	return d.AddTimeSpan(-other)
}

func (d DateSpan) SubtractDateTimeSpan(other *DateTimeSpan) *DateTimeSpan {
	return NewDateTimeSpan(
		d.AddDateSpan(other.DateSpan.Negate()),
		-other.TimeSpan,
	)
}

func (d DateSpan) Multiply(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return d.MultiplyBigInt(o).ToValue(), Undefined
		case *BigFloat:
			return Ref(d.MultiplyBigFloat(o)), Undefined
		default:
			return Undefined, Ref(NewCoerceError(d.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return d.MultiplySmallInt(other.AsSmallInt()).ToValue(), Undefined
	case FLOAT_FLAG:
		return Ref(d.MultiplyFloat(other.AsFloat())), Undefined
	default:
		return Undefined, Ref(NewCoerceError(d.Class(), other.Class()))
	}
}

func (d DateSpan) MultiplyBigInt(other *BigInt) DateSpan {
	o := other.ToGoBigInt()

	newMonths := big.NewInt(int64(d.months))
	newMonths.Mul(newMonths, o)

	newDays := big.NewInt(int64(d.days))
	newDays.Mul(newDays, o)

	return DateSpan{
		months: int32(newMonths.Int64()),
		days:   int32(newDays.Int64()),
	}
}

func (d DateSpan) MultiplySmallInt(other SmallInt) DateSpan {
	return DateSpan{
		months: int32(SmallInt(d.months) * other),
		days:   int32(SmallInt(d.days) * other),
	}
}

func (d DateSpan) MultiplyInt(other Value) DateSpan {
	if other.IsSmallInt() {
		return d.MultiplySmallInt(other.AsSmallInt())
	}

	return d.MultiplyBigInt((*BigInt)(other.Pointer()))
}

func (d DateSpan) MultiplyFloat(other Float) *DateTimeSpan {
	monthsFloat := float64(d.months) * float64(other)
	fullMonthsFloat, fracMonth := math.Modf(monthsFloat)
	months := int32(fullMonthsFloat)

	daysFloat := float64(d.days) * float64(other)
	fullDaysFloat, fracDay := math.Modf(daysFloat + fracMonth*MonthDays)
	days := int32(fullDaysFloat)

	newDateSpan := DateSpan{
		months: months,
		days:   days,
	}

	newTimeSpan := TimeSpan(fracDay * float64(Day))

	return NewDateTimeSpan(
		newDateSpan,
		newTimeSpan,
	)
}

func (d DateSpan) MultiplyBigFloat(other *BigFloat) *DateTimeSpan {
	prec := max(other.Precision(), 64)
	o := other.AsGoBigFloat()

	monthsBigFloat := new(big.Float).SetPrec(prec).SetInt64(int64(d.months))
	monthsBigFloat.Mul(monthsBigFloat, o)
	months64, _ := monthsBigFloat.Int64()
	months := int32(months64)

	fracMonth := new(big.Float).SetPrec(prec)
	fracMonth.Sub(monthsBigFloat, big.NewFloat(float64(months)))

	daysBigFloat := new(big.Float).SetPrec(prec).SetInt64(int64(d.days))
	daysBigFloat.Mul(daysBigFloat, o)
	daysBigFloat.Add(daysBigFloat, fracMonth.Mul(fracMonth, big.NewFloat(MonthDays)))
	days64, _ := daysBigFloat.Int64()
	days := int32(days64)

	fracDay := new(big.Float).SetPrec(prec)
	fracDay.Sub(daysBigFloat, big.NewFloat(float64(days)))

	newDateSpan := DateSpan{
		months: months,
		days:   days,
	}

	timeSpanInt, _ := fracDay.Mul(fracDay, big.NewFloat(float64(Day))).Int64()
	newTimeSpan := TimeSpan(timeSpanInt)

	return NewDateTimeSpan(
		newDateSpan,
		newTimeSpan,
	)
}

func initDateSpan() {
	DateSpanClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	DateSpanClass.IncludeMixin(DurationMixin)
	DateClass.AddConstantString("Span", Ref(DateSpanClass))
}
