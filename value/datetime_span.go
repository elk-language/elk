package value

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/elk-language/elk/value/durationscanner"
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

func (d *DateTimeSpan) Dup() *DateTimeSpan {
	return &DateTimeSpan{
		DateSpan: d.DateSpan,
		TimeSpan: d.TimeSpan,
	}
}

func (d *DateTimeSpan) Copy() Reference {
	return d.Dup()
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
		if buff.Len() != 0 {
			buff.WriteByte(' ')
		}
		buff.WriteString(d.TimeSpan.String())
	}

	return buff.String()
}

func (d *DateTimeSpan) ToDateTime() *DateTime {
	return d.ToDateTimeWithZone(LocalTimezone)
}

func (d *DateTimeSpan) ToDateTimeWithZone(zone *Timezone) *DateTime {
	date := d.DateSpan.ToDate()
	time := d.TimeSpan.ToTime()

	return NewDateTime(
		date.Year(),
		date.Month(),
		date.Day(),
		time.Hour(),
		time.Minute(),
		time.Second(),
		time.Millisecond(),
		time.Microsecond(),
		time.Nanosecond(),
		zone,
	)
}

func (x *DateTimeSpan) Cmp(y *DateTimeSpan) int {
	dtx := MakeZeroDateTime()
	dtx = *dtx.AddDateTimeSpan(x)

	dty := MakeZeroDateTime()
	dty = *dty.AddDateTimeSpan(y)

	return dtx.Cmp(&dty)
}

func (d *DateTimeSpan) CompareVal(other Value) (Value, Value) {
	switch other.flag {
	case DATE_SPAN_FLAG:
		return SmallInt(d.Cmp(other.AsDateSpan().ToDateTimeSpan())).ToValue(), Undefined
	case REFERENCE_FLAG:
	default:
		return Undefined, Ref(NewCoerceError(d.Class(), other.Class()))
	}

	if !other.IsReference() {
		return Undefined, Ref(NewCoerceError(d.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case DateSpan:
		return SmallInt(d.Cmp(o.ToDateTimeSpan())).ToValue(), Undefined
	case *DateTimeSpan:
		return SmallInt(d.Cmp(o)).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(d.Class(), other.Class()))
	}
}

// Check whether d is greater than other and return an error
// if something went wrong.
func (d *DateTimeSpan) GreaterThan(other Value) (result bool, err Value) {
	switch other.flag {
	case DATE_SPAN_FLAG:
		return d.Cmp(other.AsDateSpan().ToDateTimeSpan()) == 1, err
	case REFERENCE_FLAG:
	default:
		return result, Ref(NewCoerceError(d.Class(), other.Class()))
	}

	if !other.IsReference() {
		return result, Ref(NewCoerceError(d.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case DateSpan:
		return d.Cmp(o.ToDateTimeSpan()) == 1, err
	case *DateTimeSpan:
		return d.Cmp(o) == 1, err
	default:
		return result, Ref(NewCoerceError(d.Class(), other.Class()))
	}
}

func (d *DateTimeSpan) GreaterThanVal(other Value) (Value, Value) {
	result, err := d.GreaterThan(other)
	return ToElkBool(result), err
}

// Check whether d is greater than or equal to other and return an error
// if something went wrong.
func (d *DateTimeSpan) GreaterThanEqual(other Value) (result bool, err Value) {
	switch other.flag {
	case DATE_SPAN_FLAG:
		return d.Cmp(other.AsDateSpan().ToDateTimeSpan()) >= 0, err
	case REFERENCE_FLAG:
	default:
		return result, Ref(NewCoerceError(d.Class(), other.Class()))
	}

	if !other.IsReference() {
		return result, Ref(NewCoerceError(d.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case DateSpan:
		return d.Cmp(o.ToDateTimeSpan()) >= 0, err
	case *DateTimeSpan:
		return d.Cmp(o) >= 0, err
	default:
		return result, Ref(NewCoerceError(d.Class(), other.Class()))
	}
}

func (d *DateTimeSpan) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := d.GreaterThanEqual(other)
	return ToElkBool(result), err
}

// Check whether d is less than other and return an error
// if something went wrong.
func (d *DateTimeSpan) LessThan(other Value) (result bool, err Value) {
	switch other.flag {
	case DATE_SPAN_FLAG:
		return d.Cmp(other.AsDateSpan().ToDateTimeSpan()) == -1, err
	case REFERENCE_FLAG:
	default:
		return result, Ref(NewCoerceError(d.Class(), other.Class()))
	}

	if !other.IsReference() {
		return result, Ref(NewCoerceError(d.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case DateSpan:
		return d.Cmp(o.ToDateTimeSpan()) == -1, err
	case *DateTimeSpan:
		return d.Cmp(o) == -1, err
	default:
		return result, Ref(NewCoerceError(d.Class(), other.Class()))
	}
}

func (d *DateTimeSpan) LessThanVal(other Value) (Value, Value) {
	result, err := d.LessThan(other)
	return ToElkBool(result), err
}

// Check whether d is less than or equal to other and return an error
// if something went wrong.
func (d *DateTimeSpan) LessThanEqual(other Value) (result bool, err Value) {
	switch other.flag {
	case DATE_SPAN_FLAG:
		return d.Cmp(other.AsDateSpan().ToDateTimeSpan()) <= 0, err
	case REFERENCE_FLAG:
	default:
		return result, Ref(NewCoerceError(d.Class(), other.Class()))
	}

	if !other.IsReference() {
		return result, Ref(NewCoerceError(d.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case DateSpan:
		return d.Cmp(o.ToDateTimeSpan()) <= 0, err
	case *DateTimeSpan:
		return d.Cmp(o) <= 0, err
	default:
		return result, Ref(NewCoerceError(d.Class(), other.Class()))
	}
}

func (d *DateTimeSpan) Equal(other Value) bool {
	if !other.IsReference() {
		return false
	}

	o, ok := other.AsReference().(*DateTimeSpan)
	if !ok {
		return false
	}

	return d.DateSpan.Equal(o.DateSpan.ToValue()) &&
		d.TimeSpan.Equal(o.TimeSpan.ToValue())
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

func (d *DateTimeSpan) AddMutTimeSpan(other TimeSpan) *DateTimeSpan {
	d.TimeSpan = d.TimeSpan.AddTimeSpan(other)
	return d
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

func (d *DateTimeSpan) AddMutDateTimeSpan(other *DateTimeSpan) *DateTimeSpan {
	d.DateSpan = d.DateSpan.AddDateSpan(other.DateSpan)
	d.TimeSpan = d.TimeSpan.AddTimeSpan(other.TimeSpan)
	return d
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

func (d *DateTimeSpan) Multiply(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return Ref(d.MultiplyBigInt(o)), Undefined
		case *BigFloat:
			return Ref(d.MultiplyBigFloat(o)), Undefined
		default:
			return Undefined, Ref(NewCoerceError(d.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return Ref(d.MultiplySmallInt(other.AsSmallInt())), Undefined
	case FLOAT_FLAG:
		return Ref(d.MultiplyFloat(other.AsFloat())), Undefined
	default:
		return Undefined, Ref(NewCoerceError(d.Class(), other.Class()))
	}
}

func (d *DateTimeSpan) MultiplyBigInt(other *BigInt) *DateTimeSpan {
	o := other.ToGoBigInt()

	num := new(big.Int)

	newMonths := num.SetInt64(int64(d.DateSpan.months))
	newMonths.Mul(newMonths, o)
	months := int32(newMonths.Int64())

	newDays := num.SetInt64(int64(d.DateSpan.days))
	newDays.Mul(newDays, o)
	days := int32(newDays.Int64())

	dateSpan := DateSpan{
		months: months,
		days:   days,
	}

	newNanoseconds := num.SetInt64(int64(d.TimeSpan))
	newNanoseconds.Mul(newNanoseconds, o)
	timeSpan := TimeSpan(newNanoseconds.Int64())

	return NewDateTimeSpan(
		dateSpan,
		timeSpan,
	)
}

func (d *DateTimeSpan) MultiplySmallInt(other SmallInt) *DateTimeSpan {
	dateSpan := DateSpan{
		months: int32(SmallInt(d.DateSpan.months) * other),
		days:   int32(SmallInt(d.DateSpan.days) * other),
	}

	timeSpan := TimeSpan(int64(d.TimeSpan) * int64(other))

	return NewDateTimeSpan(
		dateSpan,
		timeSpan,
	)
}

func (d *DateTimeSpan) MultiplyInt(other Value) *DateTimeSpan {
	if other.IsSmallInt() {
		return d.MultiplySmallInt(other.AsSmallInt())
	}

	return d.MultiplyBigInt((*BigInt)(other.Pointer()))
}

func (d *DateTimeSpan) MultiplyFloat(other Float) *DateTimeSpan {
	monthsFloat := float64(d.DateSpan.months) * float64(other)
	fullMonthsFloat, fracMonth := math.Modf(monthsFloat)
	months := int32(fullMonthsFloat)

	daysFloat := float64(d.DateSpan.days) * float64(other)
	fullDaysFloat, fracDay := math.Modf(daysFloat + fracMonth*MonthDays)
	days := int32(fullDaysFloat)

	newDateSpan := DateSpan{
		months: months,
		days:   days,
	}

	newTimeSpan := TimeSpan(float64(d.TimeSpan)*float64(other) + fracDay*float64(Day))

	return NewDateTimeSpan(
		newDateSpan,
		newTimeSpan,
	)
}

func (d *DateTimeSpan) MultiplyBigFloat(other *BigFloat) *DateTimeSpan {
	prec := max(other.Precision(), 64)

	monthsBigFloat := new(BigFloat).SetPrecision(prec).SetInt64(int64(d.DateSpan.months))
	monthsBigFloat.MulBigFloat(monthsBigFloat, other)
	months := int32(monthsBigFloat.Int64())

	fracMonth := new(BigFloat).SetPrecision(prec)
	fracMonth.SubBigFloat(monthsBigFloat, NewBigFloat(float64(months)))

	daysBigFloat := new(BigFloat).SetPrecision(prec).SetInt64(int64(d.DateSpan.days))
	daysBigFloat.MulBigFloat(daysBigFloat, other)
	daysBigFloat.AddBigFloat(daysBigFloat, fracMonth.MulBigFloat(fracMonth, NewBigFloat(MonthDays)))
	days := int32(daysBigFloat.Int64())

	fracDay := new(BigFloat).SetPrecision(prec)
	fracDay.SubBigFloat(daysBigFloat, NewBigFloat(float64(days)))

	newDateSpan := DateSpan{
		months: months,
		days:   days,
	}

	timeSpanBigFloat := fracDay.MulBigFloat(fracDay, NewBigFloat(float64(Day)))
	nanoseconds := daysBigFloat.SetInt64(int64(d.TimeSpan))
	nanoseconds.MulBigFloat(nanoseconds, other)
	timeSpanBigFloat.AddBigFloat(timeSpanBigFloat, nanoseconds)
	newTimeSpan := TimeSpan(timeSpanBigFloat.Int64())

	return NewDateTimeSpan(
		newDateSpan,
		newTimeSpan,
	)
}

func (d *DateTimeSpan) Divide(other Value) (Value, Value) {
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

func (d *DateTimeSpan) DivideBigInt(other *BigInt) *DateTimeSpan {
	prec := max(uint(other.BitSize()), 64)
	o := new(BigFloat).SetBigInt(other)

	monthsBigFloat := new(BigFloat).SetPrecision(prec).SetInt64(int64(d.DateSpan.months))
	monthsBigFloat.DivBigFloat(monthsBigFloat, o)
	months := int32(monthsBigFloat.Int64())

	fracMonth := new(BigFloat).SetPrecision(prec)
	fracMonth.SubBigFloat(monthsBigFloat, NewBigFloat(float64(months)))

	daysBigFloat := new(BigFloat).SetPrecision(prec).SetInt64(int64(d.DateSpan.days))
	daysBigFloat.MulBigFloat(daysBigFloat, o)
	daysBigFloat.AddBigFloat(daysBigFloat, fracMonth.MulBigFloat(fracMonth, NewBigFloat(MonthDays)))
	days := int32(daysBigFloat.Int64())

	fracDay := new(BigFloat).SetPrecision(prec)
	fracDay.SubBigFloat(daysBigFloat, NewBigFloat(float64(days)))

	newDateSpan := DateSpan{
		months: months,
		days:   days,
	}

	timeSpanBigFloat := fracDay.MulBigFloat(fracDay, NewBigFloat(float64(Day)))
	nanoseconds := daysBigFloat.SetInt64(int64(d.TimeSpan))
	nanoseconds.DivBigFloat(nanoseconds, o)
	timeSpanBigFloat.AddBigFloat(timeSpanBigFloat, nanoseconds)
	newTimeSpan := TimeSpan(timeSpanBigFloat.Int64())

	return NewDateTimeSpan(
		newDateSpan,
		newTimeSpan,
	)
}

func (d *DateTimeSpan) DivideSmallInt(other SmallInt) *DateTimeSpan {
	o := float64(other)

	monthsFloat := float64(d.DateSpan.months) / o
	fullMonthsFloat, fracMonth := math.Modf(monthsFloat)
	months := int32(fullMonthsFloat)

	daysFloat := float64(d.DateSpan.days) / o
	fullDaysFloat, fracDay := math.Modf(daysFloat + fracMonth*MonthDays)
	days := int32(fullDaysFloat)

	newDateSpan := DateSpan{
		months: months,
		days:   days,
	}

	newTimeSpan := TimeSpan(float64(d.TimeSpan)/o + fracDay*float64(Day))

	return NewDateTimeSpan(
		newDateSpan,
		newTimeSpan,
	)
}

func (d *DateTimeSpan) DivideInt(other Value) *DateTimeSpan {
	if other.IsSmallInt() {
		return d.DivideSmallInt(other.AsSmallInt())
	}

	return d.DivideBigInt((*BigInt)(other.Pointer()))
}

func (d *DateTimeSpan) DivideFloat(other Float) *DateTimeSpan {
	o := float64(other)
	monthsFloat := float64(d.DateSpan.months) / o
	fullMonthsFloat, fracMonth := math.Modf(monthsFloat)
	months := int32(fullMonthsFloat)

	daysFloat := float64(d.DateSpan.days) / o
	fullDaysFloat, fracDay := math.Modf(daysFloat + fracMonth*MonthDays)
	days := int32(fullDaysFloat)

	newDateSpan := DateSpan{
		months: months,
		days:   days,
	}

	newTimeSpan := TimeSpan(float64(d.TimeSpan)/o + fracDay*float64(Day))

	return NewDateTimeSpan(
		newDateSpan,
		newTimeSpan,
	)
}

func (d *DateTimeSpan) DivideBigFloat(other *BigFloat) *DateTimeSpan {
	prec := max(other.Precision(), 64)

	monthsBigFloat := new(BigFloat).SetPrecision(prec).SetInt64(int64(d.DateSpan.months))
	monthsBigFloat.DivBigFloat(monthsBigFloat, other)
	months := int32(monthsBigFloat.Int64())

	fracMonth := new(BigFloat).SetPrecision(prec)
	fracMonth.SubBigFloat(monthsBigFloat, NewBigFloat(float64(months)))

	daysBigFloat := new(BigFloat).SetPrecision(prec).SetInt64(int64(d.DateSpan.days))
	daysBigFloat.DivBigFloat(daysBigFloat, other)
	daysBigFloat.AddBigFloat(daysBigFloat, fracMonth.MulBigFloat(fracMonth, NewBigFloat(MonthDays)))
	days := int32(daysBigFloat.Int64())

	fracDay := new(BigFloat).SetPrecision(prec)
	fracDay.SubBigFloat(daysBigFloat, NewBigFloat(float64(days)))

	newDateSpan := DateSpan{
		months: months,
		days:   days,
	}

	timeSpanBigFloat := fracDay.MulBigFloat(fracDay, NewBigFloat(float64(Day)))
	nanoseconds := daysBigFloat.SetInt64(int64(d.TimeSpan))
	nanoseconds.DivBigFloat(nanoseconds, other)
	timeSpanBigFloat.AddBigFloat(timeSpanBigFloat, nanoseconds)
	newTimeSpan := TimeSpan(timeSpanBigFloat.Int64())

	return NewDateTimeSpan(
		newDateSpan,
		newTimeSpan,
	)
}

// Parses a time span string and creates a datetime span value.
// A datetime span string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "Y", "M", "D", "h", "m", "s", "ms", "us" (or "µs"), "ns".
func ParseDateTimeSpan(str string) (result *DateTimeSpan, err Value) {
	scanner := durationscanner.New(str)
	result = &DateTimeSpan{}

tokenLoop:
	for {
		token, value := scanner.Next()
		switch token {
		case durationscanner.END_OF_FILE:
			break tokenLoop
		case durationscanner.ERROR:
			return nil, Ref(Errorf(
				FormatErrorClass,
				"invalid date time span string: %s",
				value,
			))
		case durationscanner.YEARS_INT:
			bigYears, err := ParseBigInt(value, 10)
			if !err.IsUndefined() {
				return nil, err
			}

			years := int32(bigYears.ToSmallInt())
			result.DateSpan.months += years * 12
		case durationscanner.YEARS_FLOAT:
			years, er := strconv.ParseFloat(value, 64)
			if er != nil {
				return nil, Ref(Errorf(
					FormatErrorClass,
					"invalid float in datetime span string: %s",
					er.Error(),
				))
			}
			yearsSpan := MakeDateSpan(1, 0, 0).MultiplyFloat(Float(years))
			result.AddMutDateTimeSpan(yearsSpan)
		case durationscanner.MONTHS_INT:
			bigMonths, err := ParseBigInt(value, 10)
			if !err.IsUndefined() {
				return nil, err
			}

			months := int32(bigMonths.ToSmallInt())
			result.DateSpan.months += months
		case durationscanner.MONTHS_FLOAT:
			months, er := strconv.ParseFloat(value, 64)
			if er != nil {
				return nil, Ref(Errorf(
					FormatErrorClass,
					"invalid float in datetime span string: %s",
					er.Error(),
				))
			}
			monthsSpan := MakeDateSpan(0, 1, 0).MultiplyFloat(Float(months))
			result.AddMutDateTimeSpan(monthsSpan)
		case durationscanner.DAYS_INT:
			bigDays, err := ParseBigInt(value, 10)
			if !err.IsUndefined() {
				return nil, err
			}

			days := int32(bigDays.ToSmallInt())
			result.DateSpan.days += days
		case durationscanner.DAYS_FLOAT:
			days, er := strconv.ParseFloat(value, 64)
			if er != nil {
				return nil, Ref(Errorf(
					FormatErrorClass,
					"invalid float in datetime span string: %s",
					er.Error(),
				))
			}
			daysSpan := MakeDateSpan(0, 0, 1).MultiplyFloat(Float(days))
			result.AddMutDateTimeSpan(daysSpan)
		case durationscanner.HOURS_INT:
			bigHours, err := ParseBigInt(value, 10)
			if !err.IsUndefined() {
				return nil, err
			}

			hours := TimeSpan(bigHours.ToSmallInt())
			result.TimeSpan += hours * Hour
		case durationscanner.HOURS_FLOAT:
			hours, er := strconv.ParseFloat(value, 64)
			if er != nil {
				return nil, Ref(Errorf(
					FormatErrorClass,
					"invalid float in datetime span string: %s",
					er.Error(),
				))
			}
			hoursSpan := Hour.MultiplyFloat(Float(hours))
			result.AddMutTimeSpan(hoursSpan)
		case durationscanner.MINUTES_INT:
			bigMinutes, err := ParseBigInt(value, 10)
			if !err.IsUndefined() {
				return nil, err
			}

			minutes := TimeSpan(bigMinutes.ToSmallInt())
			result.TimeSpan += minutes * Minute
		case durationscanner.MINUTES_FLOAT:
			minutes, er := strconv.ParseFloat(value, 64)
			if er != nil {
				return nil, Ref(Errorf(
					FormatErrorClass,
					"invalid float in datetime span string: %s",
					er.Error(),
				))
			}
			minutesSpan := Minute.MultiplyFloat(Float(minutes))
			result.AddMutTimeSpan(minutesSpan)
		case durationscanner.SECONDS_INT:
			bigSeconds, err := ParseBigInt(value, 10)
			if !err.IsUndefined() {
				return nil, err
			}

			seconds := TimeSpan(bigSeconds.ToSmallInt())
			result.TimeSpan += seconds * Second
		case durationscanner.SECONDS_FLOAT:
			seconds, er := strconv.ParseFloat(value, 64)
			if er != nil {
				return nil, Ref(Errorf(
					FormatErrorClass,
					"invalid float in datetime span string: %s",
					er.Error(),
				))
			}
			secondsSpan := Second.MultiplyFloat(Float(seconds))
			result.AddMutTimeSpan(secondsSpan)
		case durationscanner.MILLISECONDS_INT:
			bigMilliseconds, err := ParseBigInt(value, 10)
			if !err.IsUndefined() {
				return nil, err
			}

			milliseconds := TimeSpan(bigMilliseconds.ToSmallInt())
			result.TimeSpan += milliseconds * Millisecond
		case durationscanner.MILLISECONDS_FLOAT:
			milliseconds, er := strconv.ParseFloat(value, 64)
			if er != nil {
				return nil, Ref(Errorf(
					FormatErrorClass,
					"invalid float in datetime span string: %s",
					er.Error(),
				))
			}
			millisecondsSpan := Millisecond.MultiplyFloat(Float(milliseconds))
			result.AddMutTimeSpan(millisecondsSpan)
		case durationscanner.MICROSECONDS_INT:
			bigMicroseconds, err := ParseBigInt(value, 10)
			if !err.IsUndefined() {
				return nil, err
			}

			microseconds := TimeSpan(bigMicroseconds.ToSmallInt())
			result.TimeSpan += microseconds * Microsecond
		case durationscanner.MICROSECONDS_FLOAT:
			microseconds, er := strconv.ParseFloat(value, 64)
			if er != nil {
				return nil, Ref(Errorf(
					FormatErrorClass,
					"invalid float in datetime span string: %s",
					er.Error(),
				))
			}
			microsecondsSpan := Microsecond.MultiplyFloat(Float(microseconds))
			result.AddMutTimeSpan(microsecondsSpan)
		case durationscanner.NANOSECONDS_INT:
			bigNanoseconds, err := ParseBigInt(value, 10)
			if !err.IsUndefined() {
				return nil, err
			}

			nanoseconds := TimeSpan(bigNanoseconds.ToSmallInt())
			result.TimeSpan += nanoseconds * Nanosecond
		case durationscanner.NANOSECONDS_FLOAT:
			nanoseconds, er := strconv.ParseFloat(value, 64)
			if er != nil {
				return nil, Ref(Errorf(
					FormatErrorClass,
					"invalid float in datetime span string: %s",
					er.Error(),
				))
			}
			nanosecondsSpan := Nanosecond.MultiplyFloat(Float(nanoseconds))
			result.AddMutTimeSpan(nanosecondsSpan)
		default:
			return nil, Ref(Errorf(
				FormatErrorClass,
				"undefined datetime span token: %s",
				token.String(),
			))
		}

	}

	result.Normalise()
	return result, Undefined
}

func initDateTimeSpan() {
	DateTimeSpanClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	DateTimeSpanClass.IncludeMixin(DurationMixin)
	DateTimeClass.AddConstantString("Span", Ref(DateTimeSpanClass))
}
