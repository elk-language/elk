package value

import (
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

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

	return Date{bits: bits}.Normalize()
}

func (d Date) Normalize() Date {
	return d.ToDateTime().Date()
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

func (d *Date) SetYear(v int) {
	*d = MakeDate(v, d.Month(), d.Day())
}

func (d Date) Month() int {
	m := int32((d.bits >> 5) & 0b1111)
	return int(m)
}

func (d *Date) SetMonth(v int) {
	*d = MakeDate(d.Year(), v, d.Day())
}

func (d Date) Day() int {
	m := int32(d.bits & 0b11111)
	return int(m)
}

func (d *Date) SetDay(v int) {
	*d = MakeDate(d.Year(), d.Month(), v)
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

func ParseDate(formatString, input string) (result Date, err Value) {
	scanner := timescanner.New(formatString)
	currentInput := input

tokenLoop:
	for {
		token, value := scanner.Next()
		switch token {
		case timescanner.END_OF_FILE:
			if len(currentInput) > 0 {
				return Date{}, Ref(NewIncompatibleDateFormatError(formatString, input))
			}
			break tokenLoop
		case timescanner.INVALID_FORMAT_DIRECTIVE:
			return Date{}, Ref(Errorf(
				FormatErrorClass,
				"invalid date format directive: %s",
				value,
			))
		case timescanner.PERCENT:
			err = parseDateMatchText(formatString, input, &currentInput, "%")
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.NEWLINE:
			err = parseDateMatchText(formatString, input, &currentInput, "\n")
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.TAB:
			err = parseDateMatchText(formatString, input, &currentInput, "\t")
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.TEXT:
			err = parseDateMatchText(formatString, input, &currentInput, value)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.FULL_YEAR_WEEK_BASED, timescanner.FULL_YEAR_WEEK_BASED_ZERO_PADDED:
			err = parseDateYearWeekBased(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.FULL_YEAR_WEEK_BASED_SPACE_PADDED:
			err = parseDateYearWeekBased(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.FULL_YEAR, timescanner.FULL_YEAR_ZERO_PADDED:
			err = parseDateYear(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.FULL_YEAR_SPACE_PADDED:
			err = parseDateYear(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.CENTURY, timescanner.CENTURY_ZERO_PADDED:
			err = parseDateCentury(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.CENTURY_SPACE_PADDED:
			err = parseDateCentury(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.YEAR_LAST_TWO, timescanner.YEAR_LAST_TWO_ZERO_PADDED:
			err = parseDateYearLastTwo(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.YEAR_LAST_TWO_SPACE_PADDED:
			err = parseDateYearLastTwo(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.YEAR_LAST_TWO_WEEK_BASED, timescanner.YEAR_LAST_TWO_WEEK_BASED_ZERO_PADDED:
			err = parseDateYearLastTwoWeekBased(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.YEAR_LAST_TWO_WEEK_BASED_SPACE_PADDED:
			err = parseDateYearLastTwoWeekBased(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.MONTH, timescanner.MONTH_ZERO_PADDED:
			err = parseDateMonth(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.MONTH_SPACE_PADDED:
			err = parseDateMonth(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.MONTH_FULL_NAME, timescanner.MONTH_FULL_NAME_UPPERCASE:
			err = parseDateMonthName(&currentInput, &result)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.MONTH_ABBREVIATED_NAME, timescanner.MONTH_ABBREVIATED_NAME_UPPERCASE:
			err = parseDateAbbreviatedMonthName(&currentInput, &result)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.DAY_OF_MONTH, timescanner.DAY_OF_MONTH_ZERO_PADDED:
			err = parseDateDayOfMonth(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.DAY_OF_MONTH_SPACE_PADDED:
			err = parseDateDayOfMonth(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.DAY_OF_YEAR, timescanner.DAY_OF_YEAR_ZERO_PADDED:
			err = parseDateDayOfYear(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.DAY_OF_YEAR_SPACE_PADDED:
			err = parseDateDayOfYear(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.DAY_OF_WEEK_FULL_NAME, timescanner.DAY_OF_WEEK_FULL_NAME_UPPERCASE:
			err = parseDateDayOfWeek(&currentInput, &result)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.DAY_OF_WEEK_ABBREVIATED_NAME, timescanner.DAY_OF_WEEK_ABBREVIATED_NAME_UPPERCASE:
			err = parseDateAbbreviatedDayOfWeek(&currentInput, &result)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.DAY_OF_WEEK_NUMBER:
			err = parseDateDayOfWeekNumber(formatString, input, &currentInput, &result)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.DAY_OF_WEEK_NUMBER_ALT:
			err = parseDateDayOfWeekNumberAlt(formatString, input, &currentInput, &result)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.WEEK_OF_WEEK_BASED_YEAR, timescanner.WEEK_OF_WEEK_BASED_YEAR_ZERO_PADDED:
			err = parseWeekOfWeekBasedYear(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.WEEK_OF_WEEK_BASED_YEAR_SPACE_PADDED:
			err = parseWeekOfWeekBasedYear(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.WEEK_OF_YEAR, timescanner.WEEK_OF_YEAR_ZERO_PADDED:
			err = parseWeekOfYear(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.WEEK_OF_YEAR_SPACE_PADDED:
			err = parseWeekOfYear(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.WEEK_OF_YEAR_ALT, timescanner.WEEK_OF_YEAR_ALT_ZERO_PADDED:
			err = parseWeekOfYearAlt(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.WEEK_OF_YEAR_ALT_SPACE_PADDED:
			err = parseWeekOfYearAlt(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.DATE:
			err = parseDateMonth(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
			err = parseDateMatchText(formatString, input, &currentInput, "/")
			if !err.IsUndefined() {
				return Date{}, err
			}
			err = parseDateDayOfMonth(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
			err = parseDateMatchText(formatString, input, &currentInput, "/")
			if !err.IsUndefined() {
				return Date{}, err
			}
			err = parseDateYearLastTwo(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
		case timescanner.ISO8601_DATE:
			err = parseDateYear(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
			err = parseDateMatchText(formatString, input, &currentInput, "-")
			if !err.IsUndefined() {
				return Date{}, err
			}
			err = parseDateMonth(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
			err = parseDateMatchText(formatString, input, &currentInput, "-")
			if !err.IsUndefined() {
				return Date{}, err
			}
			err = parseDateDayOfMonth(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Date{}, err
			}
		default:
			return Date{}, Ref(Errorf(
				FormatErrorClass,
				"unsupported date format directive: %s",
				token.String(),
			))
		}

	}

	return result, err
}

var months = map[string]int{
	"january":   1,
	"february":  2,
	"march":     3,
	"april":     4,
	"may":       5,
	"june":      6,
	"july":      7,
	"august":    8,
	"september": 9,
	"october":   10,
	"november":  11,
	"december":  12,
}

var abbreviatedMonths = map[string]int{
	"jan": 1,
	"feb": 2,
	"mar": 3,
	"apr": 4,
	"may": 5,
	"jun": 6,
	"jul": 7,
	"aug": 8,
	"sep": 9,
	"oct": 10,
	"nov": 11,
	"dec": 12,
}

var days = map[string]int{
	"monday":    1,
	"tuesday":   2,
	"wednesday": 3,
	"thursday":  4,
	"friday":    5,
	"saturday":  6,
	"sunday":    7,
}

var abbreviatedDays = map[string]int{
	"mon": 1,
	"tue": 2,
	"wed": 3,
	"thu": 4,
	"fri": 5,
	"sat": 6,
	"sun": 7,
}

func datetimeMondayOfFirstGregorianWeek(year int) DateTime {
	yearStart := MakeDateTime(year, 1, 0, 0, 0, 0, 0, nil)
	firstDay := yearStart.WeekdayFromMonday()

	if firstDay == Monday {
		return yearStart
	}

	return *yearStart.Add(TimeSpan(8-firstDay) * Day)
}

func datetimeSundayOfFirstGregorianWeek(year int) DateTime {
	yearStart := MakeDateTime(year, 1, 0, 0, 0, 0, 0, nil)
	firstDay := yearStart.WeekdayFromSunday()

	if firstDay == SundayAlt {
		return yearStart
	}

	return *yearStart.Add(TimeSpan(7-firstDay) * Day)
}

func parseWeekOfYearAlt(formatString, input string, currentInput *string, result *Date, spacePadded bool) Value {
	var n int
	var ok bool
	n, *currentInput, ok = parseTemporalDigitsOk(*currentInput, 2, spacePadded)
	if !ok {
		return Ref(NewIncompatibleDateFormatError(formatString, input))
	}

	if n > 53 || n < 0 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for week out of range: %d",
			n,
		))
	}

	if n == 0 {
		*result = MakeDate(result.Year(), 1, 1)
		return Undefined
	}

	firstWeek := datetimeSundayOfFirstGregorianWeek(result.Year())
	datetime := firstWeek.Add(TimeSpan(n) * Week)
	*result = datetime.Date()

	return Undefined
}

func parseWeekOfYear(formatString, input string, currentInput *string, result *Date, spacePadded bool) Value {
	var n int
	var ok bool
	n, *currentInput, ok = parseTemporalDigitsOk(*currentInput, 2, spacePadded)
	if !ok {
		return Ref(NewIncompatibleDateFormatError(formatString, input))
	}

	if n > 53 || n < 0 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for week out of range: %d",
			n,
		))
	}

	if n == 0 {
		*result = MakeDate(result.Year(), 1, 1)
		return Undefined
	}

	firstWeek := datetimeMondayOfFirstGregorianWeek(result.Year())
	datetime := firstWeek.Add(TimeSpan(n) * Week)
	*result = datetime.Date()

	return Undefined
}

func parseWeekOfWeekBasedYear(formatString, input string, currentInput *string, result *Date, spacePadded bool) Value {
	var n int
	var ok bool
	n, *currentInput, ok = parseTemporalDigitsOk(*currentInput, 2, spacePadded)
	if !ok {
		return Ref(NewIncompatibleDateFormatError(formatString, input))
	}

	if n > 53 || n < 1 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for ISO week out of range: %d",
			n,
		))
	}

	yearStart := datetimeISOYearStart(result.Year())
	datetime := yearStart.Add(TimeSpan(n) * Week)
	*result = datetime.Date()

	return Undefined
}

func parseDateDayOfWeekNumberAlt(formatString, input string, currentInput *string, result *Date) Value {
	var n int
	var ok bool
	n, *currentInput, ok = parseTemporalDigitsOk(*currentInput, 1, false)
	if !ok {
		return Ref(NewIncompatibleDateFormatError(formatString, input))
	}

	if n > 6 || n < 0 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for date alt day of week number out of range: %d",
			n,
		))
	}

	currentDay := result.WeekdayFromSunday()
	diff := n - currentDay
	datetime := result.ToDateTimeValue()
	*result = datetime.Add(TimeSpan(diff) * Day).Date()

	return Undefined
}

func parseDateDayOfWeekNumber(formatString, input string, currentInput *string, result *Date) Value {
	var n int
	var ok bool
	n, *currentInput, ok = parseTemporalDigitsOk(*currentInput, 1, false)
	if !ok {
		return Ref(NewIncompatibleDateFormatError(formatString, input))
	}

	if n > 7 || n < 1 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for date day of week number out of range: %d",
			n,
		))
	}

	currentDay := result.WeekdayFromMonday()
	diff := n - currentDay
	datetime := result.ToDateTimeValue()
	*result = datetime.Add(TimeSpan(diff) * Day).Date()

	return Undefined
}

func parseDateAbbreviatedDayOfWeek(currentInput *string, result *Date) Value {
	var buffer strings.Builder

	for len(*currentInput) > 0 {
		char, size := utf8.DecodeRuneInString(*currentInput)
		*currentInput = (*currentInput)[size:]
		if !unicode.IsLetter(char) {
			break
		}
		buffer.WriteRune(unicode.ToLower(char))
	}

	dayName := buffer.String()
	_, ok := abbreviatedDays[dayName]
	if !ok {
		return Ref(Errorf(
			FormatErrorClass,
			"invalid date abbreviated day of week name: %s",
			dayName,
		))
	}

	return Undefined
}

func parseDateDayOfWeek(currentInput *string, result *Date) Value {
	var buffer strings.Builder

	for len(*currentInput) > 0 {
		char, size := utf8.DecodeRuneInString(*currentInput)
		*currentInput = (*currentInput)[size:]
		if !unicode.IsLetter(char) {
			break
		}
		buffer.WriteRune(unicode.ToLower(char))
	}

	dayName := buffer.String()
	_, ok := days[dayName]
	if !ok {
		return Ref(Errorf(
			FormatErrorClass,
			"invalid date day of week name: %s",
			dayName,
		))
	}

	return Undefined
}

func parseDateDayOfYear(formatString, input string, currentInput *string, result *Date, spacePadded bool) Value {
	var n int
	var ok bool
	n, *currentInput, ok = parseTemporalDigitsOk(*currentInput, 3, spacePadded)
	if !ok {
		return Ref(NewIncompatibleDateFormatError(formatString, input))
	}

	if n > 366 || n < 0 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for date day of year out of range: %d",
			n,
		))
	}

	year := result.Year()
	startOfYear := MakeDateTime(year, 1, 1, 0, 0, 0, 0, nil)
	dateTime := startOfYear.Add(TimeSpan(n) * Day)
	*result = dateTime.Date()

	return Undefined
}

func parseDateDayOfMonth(formatString, input string, currentInput *string, result *Date, spacePadded bool) Value {
	var n int
	var ok bool
	n, *currentInput, ok = parseTemporalDigitsOk(*currentInput, 2, spacePadded)
	if !ok {
		return Ref(NewIncompatibleDateFormatError(formatString, input))
	}

	if n > 31 || n < 0 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for date day out of range: %d",
			n,
		))
	}

	result.SetDay(n)
	return Undefined
}

func parseDateAbbreviatedMonthName(currentInput *string, result *Date) Value {
	var buffer strings.Builder

	for len(*currentInput) > 0 {
		char, size := utf8.DecodeRuneInString(*currentInput)
		*currentInput = (*currentInput)[size:]
		if !unicode.IsLetter(char) {
			break
		}
		buffer.WriteRune(unicode.ToLower(char))
	}

	monthName := buffer.String()
	monthNumber, ok := abbreviatedMonths[monthName]
	if !ok {
		return Ref(Errorf(
			FormatErrorClass,
			"invalid date abbreviated month name: %s",
			monthName,
		))
	}

	result.SetMonth(monthNumber)
	return Undefined
}

func parseDateMonthName(currentInput *string, result *Date) Value {
	var buffer strings.Builder

	for len(*currentInput) > 0 {
		char, size := utf8.DecodeRuneInString(*currentInput)
		*currentInput = (*currentInput)[size:]
		if !unicode.IsLetter(char) {
			break
		}
		buffer.WriteRune(unicode.ToLower(char))
	}

	monthName := buffer.String()
	monthNumber, ok := months[monthName]
	if !ok {
		return Ref(Errorf(
			FormatErrorClass,
			"invalid date month name: %s",
			monthName,
		))
	}

	result.SetMonth(monthNumber)
	return Undefined
}

func parseDateMonth(formatString, input string, currentInput *string, result *Date, spacePadded bool) Value {
	var n int
	var ok bool
	n, *currentInput, ok = parseTemporalDigitsOk(*currentInput, 2, spacePadded)
	if !ok {
		return Ref(NewIncompatibleDateFormatError(formatString, input))
	}
	if n < 1 || n > 12 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for date month out of range: %d",
			n,
		))
	}

	result.SetMonth(n)
	return Undefined
}

func parseDateYearLastTwoWeekBased(formatString, input string, currentInput *string, result *Date, spacePadded bool) Value {
	var n int
	var ok bool
	n, *currentInput, ok = parseTemporalDigitsOk(*currentInput, 2, spacePadded)
	if !ok {
		return Ref(NewIncompatibleDateFormatError(formatString, input))
	}

	var century int
	if result.Year() == 0 {
		century = DateNow().Century()
	} else {
		century = result.Century()
	}
	yearStart := datetimeISOYearStart(century*100 + n)
	*result = yearStart.Date()
	return Undefined
}

func parseDateYearLastTwo(formatString, input string, currentInput *string, result *Date, spacePadded bool) Value {
	var n int
	var ok bool
	n, *currentInput, ok = parseTemporalDigitsOk(*currentInput, 2, spacePadded)
	if !ok {
		return Ref(NewIncompatibleDateFormatError(formatString, input))
	}

	if result.Year() == 0 {
		result.SetYear(DateNow().Century()*100 + n)
	} else {
		result.SetYear(result.Century()*100 + n)
	}
	return Undefined
}

func parseDateCentury(formatString, input string, currentInput *string, result *Date, spacePadded bool) Value {
	var n int
	var ok bool
	n, *currentInput, ok = parseTemporalDigitsOk(*currentInput, 2, spacePadded)
	if !ok {
		return Ref(NewIncompatibleDateFormatError(formatString, input))
	}

	result.SetYear((result.Year() % 100) + n*100)
	return Undefined
}

func parseDateYear(formatString, input string, currentInput *string, result *Date, spacePadded bool) Value {
	var n int
	var ok bool
	n, *currentInput, ok = parseTemporalDigitsOk(*currentInput, 4, spacePadded)
	if !ok {
		return Ref(NewIncompatibleDateFormatError(formatString, input))
	}

	result.SetYear(n)
	return Undefined
}

func parseDateYearWeekBased(formatString, input string, currentInput *string, result *Date, spacePadded bool) Value {
	var n int
	var ok bool
	n, *currentInput, ok = parseTemporalDigitsOk(*currentInput, 4, spacePadded)
	if !ok {
		return Ref(NewIncompatibleDateFormatError(formatString, input))
	}

	yearStart := datetimeISOYearStart(n)
	*result = yearStart.Date()
	return Undefined
}

func datetimeISOYearStart(year int) DateTime {
	// ISO week 1 is the week with the year's first Thursday in it.
	// So we start from Jan 4th (always in week 1) and adjust.
	jan4 := MakeDateTime(year, 1, 4, 0, 0, 0, 0, nil)

	weekday := jan4.WeekdayFromMonday()
	diff := -weekday + 1
	yearStart := jan4.Subtract(TimeSpan(diff) * Day)

	return *yearStart
}

func parseDateMatchText(formatString, input string, currentInput *string, text string) Value {
	return parseTemporalMatchText("date", formatString, input, currentInput, text)
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
