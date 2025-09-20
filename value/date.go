package value

import (
	"fmt"
	"strings"
	"time"

	"github.com/elk-language/elk/value/timescanner"
)

const DefaultDateFormat = "%Y-%m-%d"

const dateYearBias = 1 << 22

const DateMaxYear = (1 << 22) - 1
const DateMinYear = -(1 << 22)

const DateMaxMonth = 12
const DateMinMonth = 1

const DateMaxDay = 31
const DateMinDay = 1

var DateClass *Class                  // ::Std::Date
var DateErrorClass *Class             // ::Std::Date::Error
var DateInvalidYearErrorClass *Class  // ::Std::Date::InvalidYearError
var DateInvalidMonthErrorClass *Class // ::Std::Date::InvalidMonthError
var DateInvalidDayErrorClass *Class   // ::Std::Date::InvalidDayError

// Represents a calendar date (year, month, day).
// It is an inline value on both 32bit and 64bit systems.
// The year range is from `-4_194_304` to `4_194_303`
type Date struct {
	bits uint32
}

func DateNow() Date {
	t := time.Now()
	return MakeDate(t.Year(), int(t.Month()), t.Day())
}

// Make and validate a new date
func MakeValidatedDate(year, month, day int) (Date, Value) {
	if year > DateMaxYear || year < DateMinYear {
		return Date{}, Ref(
			Errorf(
				DateInvalidYearErrorClass,
				"year %d is out of range %d...%d",
				year,
				DateMinYear,
				DateMaxYear,
			),
		)
	}
	if month > DateMaxMonth || month < DateMinMonth {
		return Date{}, Ref(
			Errorf(
				DateInvalidMonthErrorClass,
				"month %d is out of range %d...%d",
				month,
				DateMinMonth,
				DateMaxMonth,
			),
		)
	}
	if day > DateMaxDay || day < DateMinDay {
		return Date{}, Ref(
			Errorf(
				DateInvalidDayErrorClass,
				"day %d is out of range %d...%d",
				day,
				DateMinDay,
				DateMaxDay,
			),
		)
	}

	return MakeDate(year, month, day), Undefined
}

// Make a new date
func MakeDate(year, month, day int) Date {
	y := uint32(year + dateYearBias)
	m := uint32(month)
	d := uint32(day)
	bits := (y << 9) | (m << 5) | d

	return Date{
		bits: bits,
	}
}

func (d Date) ToValue() Value {
	return Value{
		flag: DATE_FLAG,
		data: uintptr(d.bits),
	}
}

func (d Date) Copy() Reference {
	return d
}

func (Date) Class() *Class {
	return DateClass
}

func (Date) DirectClass() *Class {
	return DateClass
}

func (Date) SingletonClass() *Class {
	return nil
}

func (d Date) Inspect() string {
	return fmt.Sprintf("Std::Date(%d, %d, %d)", d.Year(), d.Month(), d.Day())
}

func (Date) InstanceVariables() *InstanceVariables {
	return nil
}

func (d Date) Error() string {
	return d.Inspect()
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year(), d.Month(), d.Day())
}

func (d Date) ToString() String {
	return String(d.String())
}

func (d Date) Year() int {
	y := int32(d.bits >> 9)      // extract top 23 bits
	return int(y - dateYearBias) // remove bias
}

func (d Date) Month() int {
	m := int32((d.bits >> 5) & 0b1111)
	return int(m)
}

func (d Date) Day() int {
	m := int32(d.bits & 0b11111)
	return int(m)
}

func (d Date) ToGoTime() time.Time {
	return time.Date(d.Year(), time.Month(d.Month()), d.Day(), 0, 0, 0, 0, time.UTC)
}

func (d Date) ToDateTime() *DateTime {
	return ToElkDateTime(d.ToGoTime())
}

func (d Date) ToDateTimeValue() DateTime {
	return ToElkDateTimeValue(d.ToGoTime())
}

func (d Date) ISOYear() int {
	return d.ToDateTimeValue().ISOYear()
}

func (d Date) ISOYearLastTwo() int {
	return d.ToDateTimeValue().ISOYearLastTwo()
}

func (d Date) YearLastTwo() int {
	return d.Year() % 100
}

func (d Date) Century() int {
	return d.Year() / 100
}

func (d Date) MonthName() string {
	return time.Month(d.Month()).String()
}

func (d Date) AbbreviatedMonthName() string {
	return d.MonthName()[0:3]
}

// Day of the year.
func (d Date) YearDay() int {
	return d.ToDateTimeValue().YearDay()
}

func (d Date) WeekdayName() string {
	return d.ToDateTimeValue().WeekdayName()
}

func (d Date) AbbreviatedWeekdayName() string {
	return d.ToDateTimeValue().AbbreviatedWeekdayName()
}

// Specifies the day of the week (Monday = 1, ...).
func (d Date) WeekdayFromMonday() int {
	return d.ToDateTimeValue().WeekdayFromMonday()
}

func (d Date) ISOWeek() int {
	return d.ToDateTimeValue().ISOWeek()
}

// Specifies the day of the week (Sunday = 0, ...).
func (d Date) WeekdayFromSunday() int {
	return d.ToDateTimeValue().WeekdayFromSunday()
}

// The week number of the current year as a decimal number,
// range 00 to 53, starting with the first Monday
// as the first day of week 01.
func (d Date) WeekFromMonday() int {
	return d.ToDateTimeValue().WeekFromMonday()
}

// The week number of the current year as a decimal number,
// range 00 to 53, starting with the first Sunday
// as the first day of week 01.
func (d Date) WeekFromSunday() int {
	return d.ToDateTimeValue().WeekFromSunday()
}

func (d Date) IsSunday() bool {
	return d.ToDateTimeValue().IsSunday()
}

func (d Date) IsMonday() bool {
	return d.ToDateTimeValue().IsMonday()
}

func (d Date) IsTuesday() bool {
	return d.ToDateTimeValue().IsTuesday()
}

func (d Date) IsWednesday() bool {
	return d.ToDateTimeValue().IsWednesday()
}

func (d Date) IsThursday() bool {
	return d.ToDateTimeValue().IsThursday()
}

func (d Date) IsFriday() bool {
	return d.ToDateTimeValue().IsFriday()
}

func (d Date) IsSaturday() bool {
	return d.ToDateTimeValue().IsSaturday()
}

func (d Date) MustFormat(formatString string) string {
	result, err := d.Format(formatString)
	if !err.IsUndefined() {
		panic(err)
	}

	return result
}

// Create a string formatted according to the given format string.
func (d Date) Format(formatString string) (_ string, err Value) {
	scanner := timescanner.New(formatString)
	var buffer strings.Builder

tokenLoop:
	for {
		token, value := scanner.Next()
		switch token {
		case timescanner.END_OF_FILE:
			break tokenLoop
		case timescanner.INVALID_FORMAT_DIRECTIVE:
			return "", Ref(Errorf(
				FormatErrorClass,
				"invalid format directive: %s",
				value,
			))
		case timescanner.PERCENT:
			buffer.WriteByte('%')
		case timescanner.NEWLINE:
			buffer.WriteByte('\n')
		case timescanner.TAB:
			buffer.WriteByte('\t')
		case timescanner.TEXT:
			buffer.WriteString(value)
		case timescanner.FULL_YEAR_WEEK_BASED:
			fmt.Fprintf(&buffer, "%d", d.ISOYear())
		case timescanner.FULL_YEAR_WEEK_BASED_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%4d", d.ISOYear())
		case timescanner.FULL_YEAR_WEEK_BASED_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%04d", d.ISOYear())
		case timescanner.FULL_YEAR:
			fmt.Fprintf(&buffer, "%d", d.Year())
		case timescanner.FULL_YEAR_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%4d", d.Year())
		case timescanner.FULL_YEAR_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%04d", d.Year())
		case timescanner.CENTURY:
			fmt.Fprintf(&buffer, "%d", d.Century())
		case timescanner.CENTURY_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", d.Century())
		case timescanner.CENTURY_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", d.Century())
		case timescanner.YEAR_LAST_TWO_WEEK_BASED:
			fmt.Fprintf(&buffer, "%d", d.ISOYearLastTwo())
		case timescanner.YEAR_LAST_TWO_WEEK_BASED_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", d.ISOYearLastTwo())
		case timescanner.YEAR_LAST_TWO_WEEK_BASED_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", d.ISOYearLastTwo())
		case timescanner.YEAR_LAST_TWO:
			fmt.Fprintf(&buffer, "%d", d.YearLastTwo())
		case timescanner.YEAR_LAST_TWO_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", d.YearLastTwo())
		case timescanner.YEAR_LAST_TWO_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", d.YearLastTwo())
		case timescanner.MONTH:
			fmt.Fprintf(&buffer, "%d", d.Month())
		case timescanner.MONTH_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", d.Month())
		case timescanner.MONTH_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", d.Month())
		case timescanner.MONTH_FULL_NAME:
			buffer.WriteString(d.MonthName())
		case timescanner.MONTH_FULL_NAME_UPPERCASE:
			buffer.WriteString(strings.ToUpper(d.MonthName()))
		case timescanner.MONTH_ABBREVIATED_NAME:
			buffer.WriteString(d.AbbreviatedMonthName())
		case timescanner.MONTH_ABBREVIATED_NAME_UPPERCASE:
			buffer.WriteString(strings.ToUpper(d.AbbreviatedMonthName()))
		case timescanner.DAY_OF_MONTH:
			fmt.Fprintf(&buffer, "%d", d.Day())
		case timescanner.DAY_OF_MONTH_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", d.Day())
		case timescanner.DAY_OF_MONTH_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", d.Day())
		case timescanner.DAY_OF_YEAR:
			fmt.Fprintf(&buffer, "%d", d.YearDay())
		case timescanner.DAY_OF_YEAR_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%3d", d.YearDay())
		case timescanner.DAY_OF_YEAR_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%03d", d.YearDay())
		case timescanner.DAY_OF_WEEK_FULL_NAME:
			buffer.WriteString(d.WeekdayName())
		case timescanner.DAY_OF_WEEK_FULL_NAME_UPPERCASE:
			buffer.WriteString(strings.ToUpper(d.WeekdayName()))
		case timescanner.DAY_OF_WEEK_ABBREVIATED_NAME:
			buffer.WriteString(d.AbbreviatedWeekdayName())
		case timescanner.DAY_OF_WEEK_ABBREVIATED_NAME_UPPERCASE:
			buffer.WriteString(strings.ToUpper(d.AbbreviatedWeekdayName()))
		case timescanner.DAY_OF_WEEK_NUMBER:
			fmt.Fprintf(&buffer, "%d", d.WeekdayFromMonday())
		case timescanner.DAY_OF_WEEK_NUMBER_ALT:
			fmt.Fprintf(&buffer, "%d", d.WeekdayFromSunday())
		case timescanner.WEEK_OF_WEEK_BASED_YEAR:
			fmt.Fprintf(&buffer, "%d", d.ISOWeek())
		case timescanner.WEEK_OF_WEEK_BASED_YEAR_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", d.ISOWeek())
		case timescanner.WEEK_OF_WEEK_BASED_YEAR_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", d.ISOWeek())
		case timescanner.WEEK_OF_YEAR:
			fmt.Fprintf(&buffer, "%d", d.WeekFromMonday())
		case timescanner.WEEK_OF_YEAR_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", d.WeekFromMonday())
		case timescanner.WEEK_OF_YEAR_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", d.WeekFromMonday())
		case timescanner.WEEK_OF_YEAR_ALT:
			fmt.Fprintf(&buffer, "%d", d.WeekFromSunday())
		case timescanner.WEEK_OF_YEAR_ALT_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", d.WeekFromSunday())
		case timescanner.WEEK_OF_YEAR_ALT_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", d.WeekFromSunday())
		case timescanner.DATE:
			fmt.Fprintf(&buffer, "%02d/%02d/%02d", d.Month(), d.Day(), d.YearLastTwo())
		case timescanner.ISO8601_DATE:
			fmt.Fprintf(&buffer, "%04d-%02d-%02d", d.Year(), d.Month(), d.Day())
		default:
			return "", Ref(Errorf(
				FormatErrorClass,
				"unsupported format directive: %s",
				token.String(),
			))
		}

	}

	return buffer.String(), err
}

func initDate() {
	DateClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	StdModule.AddConstantString("Date", Ref(DateClass))
	DateClass.AddConstantString("DEFAULT_FORMAT", Ref(String(DefaultDateFormat)))

	DateErrorClass = NewClassWithOptions(ClassWithSuperclass(ErrorClass))
	DateClass.AddConstantString("Error", Ref(DateErrorClass))

	DateInvalidYearErrorClass = NewClassWithOptions(ClassWithSuperclass(DateErrorClass))
	DateClass.AddConstantString("InvalidYearError", Ref(DateInvalidYearErrorClass))

	DateInvalidMonthErrorClass = NewClassWithOptions(ClassWithSuperclass(DateErrorClass))
	DateClass.AddConstantString("InvalidMonthError", Ref(DateInvalidMonthErrorClass))

	DateInvalidDayErrorClass = NewClassWithOptions(ClassWithSuperclass(DateErrorClass))
	DateClass.AddConstantString("InvalidDayError", Ref(DateInvalidDayErrorClass))
}
