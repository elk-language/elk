package value

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/elk-language/elk/value/durationscanner"
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
		if buff.Len() != 0 {
			buff.WriteByte(' ')
		}
		fmt.Fprintf(&buff, "%dM", months)
	}
	if days != 0 || d.IsZero() {
		if buff.Len() != 0 {
			buff.WriteByte(' ')
		}
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
	months := int32(newMonths.Int64())

	newDays := newMonths.SetInt64(int64(d.days))
	newDays.Mul(newDays, o)
	days := int32(newDays.Int64())

	return DateSpan{
		months: months,
		days:   days,
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

	monthsBigFloat := new(BigFloat).SetPrecision(prec).SetInt64(int64(d.months))
	monthsBigFloat.MulBigFloat(monthsBigFloat, other)
	months := int32(monthsBigFloat.Int64())

	fracMonth := new(BigFloat).SetPrecision(prec)
	fracMonth.SubBigFloat(monthsBigFloat, NewBigFloat(float64(months)))

	daysBigFloat := new(BigFloat).SetPrecision(prec).SetInt64(int64(d.days))
	daysBigFloat.MulBigFloat(daysBigFloat, other)
	daysBigFloat.AddBigFloat(daysBigFloat, fracMonth.MulBigFloat(fracMonth, NewBigFloat(MonthDays)))
	days := int32(daysBigFloat.Int64())

	fracDay := new(BigFloat).SetPrecision(prec)
	fracDay.SubBigFloat(daysBigFloat, NewBigFloat(float64(days)))

	newDateSpan := DateSpan{
		months: months,
		days:   days,
	}

	timeSpanInt := fracDay.MulBigFloat(fracDay, NewBigFloat(float64(Day))).Int64()
	newTimeSpan := TimeSpan(timeSpanInt)

	return NewDateTimeSpan(
		newDateSpan,
		newTimeSpan,
	)
}

func (d DateSpan) Divide(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return Ref(d.DivideBigInt(o)), Undefined
		case *BigFloat:
			return Ref(d.DivideBigFloat(o)), Undefined
		default:
			return Undefined, Ref(NewCoerceError(d.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return Ref(d.DivideSmallInt(other.AsSmallInt())), Undefined
	case FLOAT_FLAG:
		return Ref(d.DivideFloat(other.AsFloat())), Undefined
	default:
		return Undefined, Ref(NewCoerceError(d.Class(), other.Class()))
	}
}

func (d DateSpan) DivideBigInt(other *BigInt) *DateTimeSpan {
	prec := max(uint(other.BitSize()), 64)
	o := new(BigFloat).SetBigInt(other)

	monthsBigFloat := new(BigFloat).SetPrecision(prec).SetInt64(int64(d.months))
	monthsBigFloat.DivBigFloat(monthsBigFloat, o)
	months := int32(monthsBigFloat.Int64())

	fracMonth := new(BigFloat).SetPrecision(prec)
	fracMonth.SubBigFloat(monthsBigFloat, NewBigFloat(float64(months)))

	daysBigFloat := new(BigFloat).SetPrecision(prec).SetInt64(int64(d.days))
	daysBigFloat.MulBigFloat(daysBigFloat, o)
	daysBigFloat.AddBigFloat(daysBigFloat, fracMonth.MulBigFloat(fracMonth, NewBigFloat(MonthDays)))
	days := int32(daysBigFloat.Int64())

	fracDay := new(BigFloat).SetPrecision(prec)
	fracDay.SubBigFloat(daysBigFloat, NewBigFloat(float64(days)))

	newDateSpan := DateSpan{
		months: months,
		days:   days,
	}

	timeSpanInt := fracDay.MulBigFloat(fracDay, NewBigFloat(float64(Day))).Int64()
	newTimeSpan := TimeSpan(timeSpanInt)

	return NewDateTimeSpan(
		newDateSpan,
		newTimeSpan,
	)
}

func (d DateSpan) DivideSmallInt(other SmallInt) *DateTimeSpan {
	o := float64(other)

	monthsFloat := float64(d.months) / o
	fullMonthsFloat, fracMonth := math.Modf(monthsFloat)
	months := int32(fullMonthsFloat)

	daysFloat := float64(d.days) / o
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

func (d DateSpan) DivideInt(other Value) *DateTimeSpan {
	if other.IsSmallInt() {
		return d.DivideSmallInt(other.AsSmallInt())
	}

	return d.DivideBigInt((*BigInt)(other.Pointer()))
}

func (d DateSpan) DivideFloat(other Float) *DateTimeSpan {
	monthsFloat := float64(d.months) / float64(other)
	fullMonthsFloat, fracMonth := math.Modf(monthsFloat)
	months := int32(fullMonthsFloat)

	daysFloat := float64(d.days) / float64(other)
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

func (d DateSpan) DivideBigFloat(other *BigFloat) *DateTimeSpan {
	prec := max(other.Precision(), 64)

	monthsBigFloat := new(BigFloat).SetPrecision(prec).SetInt64(int64(d.months))
	monthsBigFloat.DivBigFloat(monthsBigFloat, other)
	months := int32(monthsBigFloat.Int64())

	fracMonth := new(BigFloat).SetPrecision(prec)
	fracMonth.SubBigFloat(monthsBigFloat, NewBigFloat(float64(months)))

	daysBigFloat := new(BigFloat).SetPrecision(prec).SetInt64(int64(d.days))
	daysBigFloat.DivBigFloat(daysBigFloat, other)
	daysBigFloat.AddBigFloat(daysBigFloat, fracMonth.MulBigFloat(fracMonth, NewBigFloat(MonthDays)))
	days := int32(daysBigFloat.Int64())

	fracDay := new(BigFloat).SetPrecision(prec)
	fracDay.SubBigFloat(daysBigFloat, NewBigFloat(float64(days)))

	newDateSpan := DateSpan{
		months: months,
		days:   days,
	}

	timeSpanInt := fracDay.MulBigFloat(fracDay, NewBigFloat(float64(Day))).Int64()
	newTimeSpan := TimeSpan(timeSpanInt)

	return NewDateTimeSpan(
		newDateSpan,
		newTimeSpan,
	)
}

// Create a string formatted according to the given format string.
func ParseDateSpan(str string) (result DateSpan, err Value) {
	scanner := durationscanner.New(str)
	result = DateSpan{}

tokenLoop:
	for {
		token, value := scanner.Next()
		switch token {
		case durationscanner.END_OF_FILE:
			break tokenLoop
		case durationscanner.ERROR:
			return DateSpan{}, Ref(Errorf(
				FormatErrorClass,
				"invalid date span string: %s",
				value,
			))
		case durationscanner.YEARS_INT:
			bigYears, err := ParseBigInt(value, 10)
			if !err.IsUndefined() {
				return DateSpan{}, err
			}

			years := int32(bigYears.ToSmallInt())
			result.months += years * 12
		case durationscanner.MONTHS_INT:
			bigMonths, err := ParseBigInt(value, 10)
			if !err.IsUndefined() {
				return DateSpan{}, err
			}

			months := int32(bigMonths.ToSmallInt())
			result.months += months
		case durationscanner.DAYS_INT:
			bigDays, err := ParseBigInt(value, 10)
			if !err.IsUndefined() {
				return DateSpan{}, err
			}

			days := int32(bigDays.ToSmallInt())
			result.days += days
		default:
			return DateSpan{}, Ref(Errorf(
				FormatErrorClass,
				"undefined date span token: %s",
				token.String(),
			))
		}

	}

	return result, Undefined
}

func initDateSpan() {
	DateSpanClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	DateSpanClass.IncludeMixin(DurationMixin)
	DateClass.AddConstantString("Span", Ref(DateSpanClass))
}
