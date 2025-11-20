package value

import (
	"fmt"
	"math/big"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/elk/value/timescanner"
)

var DateTimeClass *Class // ::Std::DateTime

const (
	SundayAlt = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// Elk's DateTime value
type DateTime struct {
	Go time.Time
}

func (t *DateTime) Copy() Reference {
	newT := *t
	return &newT
}

func (DateTime) Class() *Class {
	return DateTimeClass
}

func (DateTime) DirectClass() *Class {
	return DateTimeClass
}

func (DateTime) SingletonClass() *Class {
	return nil
}

const DefaultDateTimeFormat = "%Y-%m-%d %H:%M:%S.%9N %:z"

func (t DateTime) Inspect() string {
	return fmt.Sprintf("Std::DateTime.parse('%s')", t.ToString().String())
}

func (t DateTime) Error() string {
	return t.Inspect()
}

func (t DateTime) InstanceVariables() *InstanceVariables {
	return nil
}

func ToElkDateTime(time time.Time) *DateTime {
	return &DateTime{Go: time}
}

func ToElkDateTimeValue(time time.Time) DateTime {
	return DateTime{Go: time}
}

func (t DateTime) ToString() String {
	return String(t.String())
}

func (t DateTime) String() string {
	return t.MustFormat(DefaultDateTimeFormat)
}

// Create a new DateTime object.
func NewDateTime(year, month, day, hour, min, sec, millisec, microsec, nsec int, zone *Timezone) *DateTime {
	t := MakeDateTime(year, month, day, hour, min, sec, millisec, microsec, nsec, zone)
	return &t
}

func MakeDateTimeFromDateAndTime(date Date, time Time, zone *Timezone) DateTime {
	return MakeDateTime(
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

func MakeZeroDateTime() DateTime {
	return MakeDateTime(0, 1, 1, 0, 0, 0, 0, 0, 0, nil)
}

// Create a new DateTime value.
func MakeDateTime(year, month, day, hour, min, sec, millisec, microsec, nsec int, zone *Timezone) DateTime {
	var location *time.Location
	if zone == nil {
		location = time.Local
	} else {
		location = zone.ToGoLocation()
	}

	return DateTime{
		Go: time.Date(
			year,
			time.Month(month),
			day,
			hour,
			min,
			sec,
			millisec*int(Millisecond)+microsec*int(Microsecond)+nsec,
			location,
		),
	}
}

func DateTimeNow() *DateTime {
	return ToElkDateTime(time.Now())
}

func (t *DateTime) Date() Date {
	return MakeDate(
		t.Year(),
		t.Month(),
		t.Day(),
	)
}

func (t *DateTime) Time() Time {
	return MakeTime(
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Millisecond(),
		t.Microsecond(),
		t.Nanosecond(),
	)
}

func (t *DateTime) InZone(zone *Timezone) *DateTime {
	time := t.Go.In(zone.ToGoLocation())
	return &DateTime{Go: time}
}

func (t *DateTime) WithZone(zone *Timezone) *DateTime {
	return NewDateTime(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Millisecond(),
		t.Microsecond(),
		t.Nanosecond(),
		zone,
	)
}

func (t *DateTime) Add(val Value) (Value, Value) {
	switch val.flag {
	case DATE_SPAN_FLAG:
		return Ref(t.AddTimeSpan(val.AsInlineTimeSpan())), Undefined
	case TIME_SPAN_FLAG:
		return Ref(t.AddTimeSpan(val.AsInlineTimeSpan())), Undefined
	case REFERENCE_FLAG:
	default:
		return Undefined, Ref(NewArgumentTypeError("other", val.Class().Inspect(), durationUnionType))
	}

	switch v := val.AsReference().(type) {
	case TimeSpan:
		return Ref(t.AddTimeSpan(v)), Undefined
	case DateSpan:
		return Ref(t.AddDateSpan(v)), Undefined
	case *DateTimeSpan:
		return Ref(t.AddDateTimeSpan(v)), Undefined
	default:
		return Undefined, Ref(NewArgumentTypeError("other", val.Class().Inspect(), durationUnionType))
	}
}

func (t *DateTime) AddDateTimeSpan(val *DateTimeSpan) *DateTime {
	t = t.AddTimeSpan(val.TimeSpan)
	t = t.AddDateSpan(val.DateSpan)
	return t
}

func (t *DateTime) AddTimeSpan(val TimeSpan) *DateTime {
	return ToElkDateTime(t.Go.Add(val.Go()))
}

func daysOfMonth(year, month int) int {
	d := MakeDateTime(year, month+1, 1, 0, 0, 0, 0, 0, 0, nil)
	lastDay := d.SubtractTimeSpan(Day)
	return lastDay.Day()
}

func (t *DateTime) AddDateSpan(val DateSpan) *DateTime {
	result := t.AddTimeSpan(TimeSpan(val.Days()) * Day)
	oldDay := result.Day()

	month := result.Month() + int(val.months)
	year := result.Year() + month/12
	month %= 12

	daysOfNewMonth := daysOfMonth(year, month)
	newDay := min(oldDay, daysOfNewMonth)

	return NewDateTime(
		year,
		month,
		newDay,
		result.Hour(),
		result.Minute(),
		result.Second(),
		result.Millisecond(),
		result.Microsecond(),
		result.Nanosecond(),
		result.Zone(),
	)
}

func (t *DateTime) Subtract(val Value) (Value, Value) {
	switch val.flag {
	case DATE_SPAN_FLAG:
		return Ref(t.SubtractTimeSpan(val.AsInlineTimeSpan())), Undefined
	case TIME_SPAN_FLAG:
		return Ref(t.SubtractTimeSpan(val.AsInlineTimeSpan())), Undefined
	case REFERENCE_FLAG:
	default:
		return Undefined, Ref(NewArgumentTypeError("other", val.Class().Inspect(), durationUnionType))
	}

	switch v := val.AsReference().(type) {
	case TimeSpan:
		return Ref(t.SubtractTimeSpan(v)), Undefined
	case DateSpan:
		return Ref(t.SubtractDateSpan(v)), Undefined
	case *DateTimeSpan:
		return Ref(t.SubtractDateTimeSpan(v)), Undefined
	default:
		return Undefined, Ref(NewArgumentTypeError("other", val.Class().Inspect(), durationUnionType))
	}
}

func (t *DateTime) SubtractDateTimeSpan(val *DateTimeSpan) *DateTime {
	t = t.SubtractTimeSpan(val.TimeSpan)
	t = t.SubtractDateSpan(val.DateSpan)
	return t
}

func (t *DateTime) SubtractDateSpan(val DateSpan) *DateTime {
	return t.ToDateTimeSpan().SubtractDateSpan(val).ToDateTime()
}

func (t *DateTime) SubtractTimeSpan(val TimeSpan) *DateTime {
	return ToElkDateTime(t.Go.Add(-val.Go()))
}

func (t *DateTime) Diff(val Value) (Value, Value) {
	switch val.flag {
	case DATE_SPAN_FLAG:
		return Ref(t.DiffDate(val.AsDate())), Undefined
	case REFERENCE_FLAG:
	default:
		return Undefined, Ref(NewArgumentTypeError("other", val.Class().Inspect(), DateTimeClass.Inspect()))
	}

	switch v := val.AsReference().(type) {
	case *DateTime:
		return Ref(t.DiffDateTime(v)), Undefined
	default:
		return Undefined, Ref(NewArgumentTypeError("other", val.Class().Inspect(), DateTimeClass.Inspect()))
	}
}

// Calculates the difference between two datetime objects.
// Returns a span.
func (t *DateTime) DiffDate(val Date) *DateTimeSpan {
	return t.ToDateTimeSpan().SubtractDateSpan(val.ToDateSpan())
}

// Calculates the difference between two datetime objects.
// Returns a span.
func (t *DateTime) DiffDateTime(val *DateTime) *DateTimeSpan {
	return t.ToDateTimeSpan().SubtractDateTimeSpan(val.ToDateTimeSpan())
}

func (t DateTime) ToGoTime() time.Time {
	return t.Go
}

func (t DateTime) ToDateTimeSpan() *DateTimeSpan {
	return NewDateTimeSpan(
		t.Date().ToDateSpan(),
		t.Time().ToTimeSpan(),
	)
}

func (t DateTime) Zone() *Timezone {
	return NewTimezone(t.Go.Location())
}

func (t DateTime) ISOYear() int {
	year, _ := t.Go.ISOWeek()
	return year
}

func (t DateTime) ISOYearLastTwo() int {
	return t.ISOYear() % 100
}

func (t DateTime) YearLastTwo() int {
	return t.Go.Year() % 100
}

func (t DateTime) Year() int {
	return t.Go.Year()
}

func (t DateTime) Century() int {
	return t.Go.Year() / 100
}

func (t DateTime) Millenium() int {
	return t.Go.Year() / 1000
}

func (t DateTime) Month() int {
	return int(t.Go.Month())
}

func (t DateTime) MonthName() string {
	return t.Go.Month().String()
}

func (t DateTime) AbbreviatedMonthName() string {
	return t.MonthName()[0:3]
}

func (t DateTime) Day() int {
	return t.Go.Day()
}

// Day of the year.
func (t DateTime) YearDay() int {
	return t.Go.YearDay()
}

func (t *DateTime) ToDate() Date {
	return MakeDate(t.Year(), t.Month(), t.Day())
}

// Day of the ISO year.
func (t DateTime) ISOYearDay() int {
	yearStart := datetimeISOYearStart(t.ISOYear())
	yearStartDate := yearStart.ToDate()
	span := t.ToDate().DiffDate(yearStartDate)
	return int(span.TotalDays().AsSmallInt())
}

// Hour in a 24 hour clock.
func (t DateTime) Hour() int {
	return t.Go.Hour()
}

// Whether the current hour is AM.
func (t DateTime) IsAM() bool {
	hour := t.Hour()

	return hour < 12
}

// Whether the current hour is PM.
func (t DateTime) IsPM() bool {
	hour := t.Hour()

	return hour >= 12
}

func (t DateTime) Meridiem() string {
	if t.IsAM() {
		return "AM"
	}

	return "PM"
}

func (t DateTime) MeridiemLowercase() string {
	if t.IsAM() {
		return "am"
	}

	return "pm"
}

// Hour in a twelve hour clock.
func (t DateTime) Hour12() int {
	hour := t.Hour()
	if hour == 0 {
		return 12
	}

	if hour <= 12 {
		return hour
	}

	return hour - 12
}

func (t DateTime) Minute() int {
	return t.Go.Minute()
}

func (t DateTime) Second() int {
	return t.Go.Second()
}

func (t DateTime) Millisecond() int {
	return t.NanosecondsInSecond() / 1000_000
}

func (t DateTime) Microsecond() int {
	return t.MicrosecondsInSecond() % 1000
}

func (t DateTime) MicrosecondsInSecond() int {
	return t.NanosecondsInSecond() / 1000
}

func (t DateTime) Nanosecond() int {
	return t.NanosecondsInSecond() % 1000
}

func (t DateTime) NanosecondsInSecond() int {
	return t.Go.Nanosecond()
}

func (t DateTime) PicosecondsInSecond() int64 {
	return int64(t.NanosecondsInSecond()) * 1000
}

func (t DateTime) FemtosecondsInSecond() int64 {
	return int64(t.NanosecondsInSecond()) * 1000_000
}

func (t DateTime) AttosecondsInSecond() int64 {
	return int64(t.NanosecondsInSecond()) * 1000_000_000
}

func (t DateTime) ZeptosecondsInSecond() *big.Int {
	i := big.NewInt(int64(t.NanosecondsInSecond()))
	i.Mul(i, big.NewInt(1000_000_000_000))
	return i
}

func (t DateTime) YoctosecondsInSecond() *big.Int {
	i := big.NewInt(int64(t.NanosecondsInSecond()))
	i.Mul(i, big.NewInt(1000_000_000_000_000))
	return i
}

func (t DateTime) ZoneName() string {
	return t.Go.Location().String()
}

func (t DateTime) ZoneAbbreviatedName() string {
	name, _ := t.Go.Zone()
	return name
}

func (t DateTime) ZoneOffset() TimeSpan {
	_, offset := t.Go.Zone()
	return TimeSpan(offset) * Second
}

func (t DateTime) ZoneOffsetSeconds() int {
	_, offset := t.Go.Zone()
	return offset
}

func (t DateTime) ZoneOffsetHours() int {
	_, offset := t.Go.Zone()
	return offset / 3600
}

func (t DateTime) ZoneOffsetHourMinutes() int {
	_, offset := t.Go.Zone()
	return (offset % 3600) / 60
}

func (t DateTime) WeekdayName() string {
	return t.Go.Weekday().String()
}

func (t DateTime) AbbreviatedWeekdayName() string {
	return t.WeekdayName()[0:3]
}

// Specifies the day of the week (Monday = 1, ...).
func (t DateTime) WeekdayFromMonday() int {
	weekday := int(t.Go.Weekday())
	if weekday == 0 {
		return 7
	}

	return weekday
}

// Specifies the day of the week (Sunday = 0, ...).
func (t DateTime) WeekdayFromSunday() int {
	return int(t.Go.Weekday())
}

func (t DateTime) UnixSeconds() int64 {
	return t.Go.Unix()
}

func (t DateTime) UnixMilliseconds() int64 {
	return t.Go.UnixMilli()
}

func (t DateTime) UnixMicroseconds() *big.Int {
	i := big.NewInt(t.UnixSeconds())
	i = i.Mul(i, big.NewInt(1000_000))
	return i.Add(i, big.NewInt(int64(t.MicrosecondsInSecond())))
}

func (t DateTime) UnixNanoseconds() *big.Int {
	i := big.NewInt(t.UnixSeconds())
	i = i.Mul(i, big.NewInt(1000_000_000))
	return i.Add(i, big.NewInt(int64(t.NanosecondsInSecond())))
}

func (t DateTime) UnixPicoseconds() *big.Int {
	i := big.NewInt(t.UnixSeconds())
	i = i.Mul(i, big.NewInt(1000_000_000_000))
	return i.Add(i, big.NewInt(t.PicosecondsInSecond()))
}

func (t DateTime) UnixFemtoseconds() *big.Int {
	i := big.NewInt(t.UnixSeconds())
	i = i.Mul(i, big.NewInt(1000_000_000_000_000))
	return i.Add(i, big.NewInt(t.FemtosecondsInSecond()))
}

func (t DateTime) UnixAttoseconds() *big.Int {
	i := big.NewInt(t.UnixSeconds())
	i = i.Mul(i, big.NewInt(1000_000_000_000_000_000))
	return i.Add(i, big.NewInt(t.AttosecondsInSecond()))
}

func (t DateTime) UnixZeptoseconds() *big.Int {
	i := t.UnixAttoseconds()
	return i.Mul(i, big.NewInt(1000))
}

func (t DateTime) UnixYoctoseconds() *big.Int {
	i := t.UnixAttoseconds()
	return i.Mul(i, big.NewInt(1000_000))
}

func (t DateTime) ISOWeek() int {
	_, week := t.Go.ISOWeek()
	return week
}

func (t DateTime) weekNumber(firstWeekday int) int {
	yday := t.YearDay()
	wday := t.WeekdayFromSunday()

	if firstWeekday == 1 {
		if wday == 0 { // sunday
			wday = 6
		} else {
			wday--
		}
	}
	ret := ((yday + 7 - wday) / 7)
	if ret < 0 {
		return 0
	}
	return ret
}

// The week number of the current year as a decimal number,
// range 00 to 53, starting with the first Monday
// as the first day of week 01.
func (t DateTime) WeekFromMonday() int {
	return t.weekNumber(1)
}

// The week number of the current year as a decimal number,
// range 00 to 53, starting with the first Sunday
// as the first day of week 01.
func (t DateTime) WeekFromSunday() int {
	return t.weekNumber(0)
}

func (t DateTime) IsSunday() bool {
	return t.WeekdayFromSunday() == 0
}

func (t DateTime) IsMonday() bool {
	return t.WeekdayFromSunday() == 1
}

func (t DateTime) IsTuesday() bool {
	return t.WeekdayFromSunday() == 2
}

func (t DateTime) IsWednesday() bool {
	return t.WeekdayFromSunday() == 3
}

func (t DateTime) IsThursday() bool {
	return t.WeekdayFromSunday() == 4
}

func (t DateTime) IsFriday() bool {
	return t.WeekdayFromSunday() == 5
}

func (t DateTime) IsSaturday() bool {
	return t.WeekdayFromSunday() == 6
}

func (t DateTime) IsUTC() bool {
	return t.Zone().IsUTC()
}

func (t DateTime) IsLocal() bool {
	return t.Zone().IsLocal()
}

// Convert the time to the UTC zone.
func (t *DateTime) UTC() *DateTime {
	return ToElkDateTime(t.Go.UTC())
}

// Convert the time to the local timezone.
func (t *DateTime) Local() *DateTime {
	return ToElkDateTime(t.Go.Local())
}

// Cmp compares x and y and returns:
//
//	  -1 if x <  y
//		 0 if x == y
//	  +1 if x >  y
func (x *DateTime) Cmp(y *DateTime) int {
	return x.Go.Compare(y.Go)
}

func (d *DateTime) CompareVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *DateTime:
			return SmallInt(d.Cmp(o)).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(d.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case DATE_FLAG:
		return SmallInt(d.Cmp(other.AsDate().ToDateTime())).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(d.Class(), other.Class()))
	}
}

func (t DateTime) MustFormat(formatString string) string {
	result, err := t.Format(formatString)
	if !err.IsUndefined() {
		panic(err)
	}

	return result
}

func parseDateTimeMatchText(formatString, input string, currentInput *string, text string) Value {
	return parseTemporalMatchText("datetime", formatString, input, currentInput, text)
}

type tmpDateTime struct {
	date tmpDate
	time tmpTime

	zone *Timezone
}

func ParseDateTime(formatString, input string) (result *DateTime, err Value) {
	scanner := timescanner.New(formatString)
	currentInput := input

	var tmp tmpDateTime

tokenLoop:
	for {
		token, value := scanner.Next()
		switch token {
		case timescanner.END_OF_FILE:
			if len(currentInput) > 0 {
				return nil, Ref(NewIncompatibleDateTimeFormatError(formatString, input))
			}
			break tokenLoop
		case timescanner.INVALID_FORMAT_DIRECTIVE:
			return nil, Ref(Errorf(
				FormatErrorClass,
				"invalid date format directive: %s",
				value,
			))
		case timescanner.PERCENT:
			err = parseDateTimeMatchText(formatString, input, &currentInput, "%")
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.NEWLINE:
			err = parseDateTimeMatchText(formatString, input, &currentInput, "\n")
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.TAB:
			err = parseDateTimeMatchText(formatString, input, &currentInput, "\t")
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.TEXT:
			err = parseDateTimeMatchText(formatString, input, &currentInput, value)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.FULL_ISO_YEAR, timescanner.FULL_ISO_YEAR_ZERO_PADDED:
			err = parseDateISOYear(formatString, input, &currentInput, &tmp.date, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.FULL_ISO_YEAR_SPACE_PADDED:
			err = parseDateISOYear(formatString, input, &currentInput, &tmp.date, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.FULL_YEAR, timescanner.FULL_YEAR_ZERO_PADDED:
			err = parseDateYear(formatString, input, &currentInput, &tmp.date, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.FULL_YEAR_SPACE_PADDED:
			err = parseDateYear(formatString, input, &currentInput, &tmp.date, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.CENTURY, timescanner.CENTURY_ZERO_PADDED:
			err = parseDateCentury(formatString, input, &currentInput, &tmp.date, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.CENTURY_SPACE_PADDED:
			err = parseDateCentury(formatString, input, &currentInput, &tmp.date, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.YEAR_LAST_TWO, timescanner.YEAR_LAST_TWO_ZERO_PADDED:
			err = parseDateYearLastTwo(formatString, input, &currentInput, &tmp.date, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.YEAR_LAST_TWO_SPACE_PADDED:
			err = parseDateYearLastTwo(formatString, input, &currentInput, &tmp.date, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.YEAR_LAST_TWO_WEEK_BASED, timescanner.YEAR_LAST_TWO_WEEK_BASED_ZERO_PADDED:
			err = parseDateYearLastTwoWeekBased(formatString, input, &currentInput, &tmp.date, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.YEAR_LAST_TWO_WEEK_BASED_SPACE_PADDED:
			err = parseDateYearLastTwoWeekBased(formatString, input, &currentInput, &tmp.date, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.MONTH, timescanner.MONTH_ZERO_PADDED:
			err = parseDateMonth(formatString, input, &currentInput, &tmp.date, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.MONTH_SPACE_PADDED:
			err = parseDateMonth(formatString, input, &currentInput, &tmp.date, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.MONTH_FULL_NAME, timescanner.MONTH_FULL_NAME_UPPERCASE:
			err = parseDateMonthName(&currentInput, &tmp.date)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.MONTH_ABBREVIATED_NAME, timescanner.MONTH_ABBREVIATED_NAME_UPPERCASE:
			err = parseDateAbbreviatedMonthName(&currentInput, &tmp.date)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.DAY_OF_MONTH, timescanner.DAY_OF_MONTH_ZERO_PADDED:
			err = parseDateDayOfMonth(formatString, input, &currentInput, &tmp.date, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.DAY_OF_MONTH_SPACE_PADDED:
			err = parseDateDayOfMonth(formatString, input, &currentInput, &tmp.date, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.DAY_OF_YEAR, timescanner.DAY_OF_YEAR_ZERO_PADDED:
			err = parseDateDayOfYear(formatString, input, &currentInput, &tmp.date, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.DAY_OF_YEAR_SPACE_PADDED:
			err = parseDateDayOfYear(formatString, input, &currentInput, &tmp.date, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.DAY_OF_WEEK_FULL_NAME, timescanner.DAY_OF_WEEK_FULL_NAME_UPPERCASE:
			err = parseDateDayOfWeekName(&currentInput, &tmp.date)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.DAY_OF_WEEK_ABBREVIATED_NAME, timescanner.DAY_OF_WEEK_ABBREVIATED_NAME_UPPERCASE:
			err = parseDateAbbreviatedDayOfWeekName(&currentInput, &tmp.date)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.DAY_OF_WEEK_NUMBER:
			err = parseDateDayOfWeekNumber(formatString, input, &currentInput, &tmp.date)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.DAY_OF_WEEK_NUMBER_ALT:
			err = parseDateDayOfWeekNumberAlt(formatString, input, &currentInput, &tmp.date)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.ISO_WEEK, timescanner.ISO_WEEK_ZERO_PADDED:
			err = parseISOWeek(formatString, input, &currentInput, &tmp.date, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.ISO_WEEK_SPACE_PADDED:
			err = parseISOWeek(formatString, input, &currentInput, &tmp.date, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.WEEK_OF_YEAR, timescanner.WEEK_OF_YEAR_ZERO_PADDED:
			err = parseWeekOfYear(formatString, input, &currentInput, &tmp.date, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.WEEK_OF_YEAR_SPACE_PADDED:
			err = parseWeekOfYear(formatString, input, &currentInput, &tmp.date, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.WEEK_OF_YEAR_ALT, timescanner.WEEK_OF_YEAR_ALT_ZERO_PADDED:
			err = parseWeekOfYearAlt(formatString, input, &currentInput, &tmp.date, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.WEEK_OF_YEAR_ALT_SPACE_PADDED:
			err = parseWeekOfYearAlt(formatString, input, &currentInput, &tmp.date, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.DATE:
			err = parseDate(formatString, input, &currentInput, &tmp.date)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.ISO8601_DATE:
			err = parseISO8601Date(formatString, input, &currentInput, &tmp.date)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.HOUR_OF_DAY, timescanner.HOUR_OF_DAY_ZERO_PADDED:
			err = parseTimeHour(formatString, input, &currentInput, &tmp.time, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.HOUR_OF_DAY_SPACE_PADDED:
			err = parseTimeHour(formatString, input, &currentInput, &tmp.time, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.HOUR_OF_DAY12, timescanner.HOUR_OF_DAY12_ZERO_PADDED:
			err = parseTime12Hour(formatString, input, &currentInput, &tmp.time, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.HOUR_OF_DAY12_SPACE_PADDED:
			err = parseTime12Hour(formatString, input, &currentInput, &tmp.time, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.MERIDIEM_INDICATOR_LOWERCASE, timescanner.MERIDIEM_INDICATOR_UPPERCASE:
			err = parseTimeMeridiem(formatString, input, &currentInput, &tmp.time)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.MINUTE_OF_HOUR, timescanner.MINUTE_OF_HOUR_ZERO_PADDED:
			err = parseTimeMinute(formatString, input, &currentInput, &tmp.time, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.MINUTE_OF_HOUR_SPACE_PADDED:
			err = parseTimeMinute(formatString, input, &currentInput, &tmp.time, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.SECOND_OF_MINUTE, timescanner.SECOND_OF_MINUTE_ZERO_PADDED:
			err = parseTimeSecond(formatString, input, &currentInput, &tmp.time, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.SECOND_OF_MINUTE_SPACE_PADDED:
			err = parseTimeSecond(formatString, input, &currentInput, &tmp.time, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.MILLISECOND_OF_SECOND, timescanner.MILLISECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeMillisecond(formatString, input, &currentInput, &tmp.time, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.MILLISECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeMillisecond(formatString, input, &currentInput, &tmp.time, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.MICROSECOND_OF_SECOND, timescanner.MICROSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeMicrosecond(formatString, input, &currentInput, &tmp.time, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.MICROSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeMicrosecond(formatString, input, &currentInput, &tmp.time, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.NANOSECOND_OF_SECOND, timescanner.NANOSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeNanosecond(formatString, input, &currentInput, &tmp.time, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.NANOSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeNanosecond(formatString, input, &currentInput, &tmp.time, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.PICOSECOND_OF_SECOND, timescanner.PICOSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeSubNanosecond(12, "picosecond", formatString, input, &currentInput, &tmp.time, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.PICOSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeSubNanosecond(12, "picosecond", formatString, input, &currentInput, &tmp.time, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.FEMTOSECOND_OF_SECOND, timescanner.FEMTOSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeSubNanosecond(15, "femtosecond", formatString, input, &currentInput, &tmp.time, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.FEMTOSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeSubNanosecond(15, "femtosecond", formatString, input, &currentInput, &tmp.time, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.ATTOSECOND_OF_SECOND, timescanner.ATTOSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeSubNanosecond(18, "attosecond", formatString, input, &currentInput, &tmp.time, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.ATTOSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeSubNanosecond(18, "attosecond", formatString, input, &currentInput, &tmp.time, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.ZEPTOSECOND_OF_SECOND, timescanner.ZEPTOSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeSubNanosecond(21, "zeptosecond", formatString, input, &currentInput, &tmp.time, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.ZEPTOSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeSubNanosecond(21, "zeptosecond", formatString, input, &currentInput, &tmp.time, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.YOCTOSECOND_OF_SECOND, timescanner.YOCTOSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeSubNanosecond(24, "yoctosecond", formatString, input, &currentInput, &tmp.time, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.YOCTOSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeSubNanosecond(24, "yoctosecond", formatString, input, &currentInput, &tmp.time, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.TIMEZONE_NAME:
			err = parseDateTimeTimezoneName(formatString, input, &currentInput, &tmp)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.TIMEZONE_IANA_NAME:
			err = parseDateTimeTimezoneIANAName(formatString, input, &currentInput, &tmp)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.TIMEZONE_OFFSET:
			err = parseDateTimeTimezoneOffset(formatString, input, &currentInput, &tmp, false)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.TIMEZONE_OFFSET_COLON:
			err = parseDateTimeTimezoneOffset(formatString, input, &currentInput, &tmp, true)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.TIME12:
			err = parseTime12(formatString, input, &currentInput, &tmp.time)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.TIME24:
			err = parseTime24(formatString, input, &currentInput, &tmp.time)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.TIME24_SECONDS:
			err = parseTime24Seconds(formatString, input, &currentInput, &tmp.time)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.DATE_AND_TIME, timescanner.DATE_AND_TIME_UPPERCASE:
			err = parseDateAndTime(formatString, input, &currentInput, &tmp)
			if !err.IsUndefined() {
				return nil, err
			}
		case timescanner.DATE1_FORMAT, timescanner.DATE1_FORMAT_UPPERCASE:
			err = parseDate1(formatString, input, &currentInput, &tmp)
			if !err.IsUndefined() {
				return nil, err
			}
		default:
			return nil, Ref(Errorf(
				FormatErrorClass,
				"unsupported datetime format directive: %s",
				token.String(),
			))
		}

	}

	date := constructDateFromTmp(tmp.date)
	time := constructTimeFromTmp(tmp.time)
	var zone *Timezone
	if tmp.zone != nil {
		zone = tmp.zone
	} else {
		zone = LocalTimezone
	}

	r := MakeDateTimeFromDateAndTime(date, time, zone)
	return &r, err
}

func parseDateTimeTimezoneOffset(formatString, input string, currentInput *string, tmp *tmpDateTime, colon bool) Value {
	if len(*currentInput) < 1 {
		return Ref(NewIncompatibleDateTimeFormatError(formatString, input))
	}

	signChar := (*currentInput)[0]
	*currentInput = (*currentInput)[1:]

	var sign int
	switch signChar {
	case '+':
		sign = 1
	case '-':
		sign = -1
	default:
		return Ref(Errorf(
			FormatErrorClass,
			"invalid datetime timezone sign: %c",
			signChar,
		))
	}

	var hours, minutes int

	hours, *currentInput = parseTemporalDigits(*currentInput, 2, false)
	if hours < 0 {
		return Ref(NewIncompatibleDateTimeFormatError(formatString, input))
	}
	if hours >= 24 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for datetime timezone offset hours out of range: %d",
			hours,
		))
	}

	if colon {
		err := parseDateTimeMatchText(formatString, input, currentInput, ":")
		if !err.IsUndefined() {
			return err
		}
	}

	minutes, *currentInput = parseTemporalDigits(*currentInput, 2, false)
	if minutes < 0 {
		return Ref(NewIncompatibleDateTimeFormatError(formatString, input))
	}
	if minutes >= 60 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for datetime timezone offset minutes out of range: %d",
			hours,
		))
	}

	offset := TimeSpan(sign) * (TimeSpan(hours)*Hour + TimeSpan(minutes)*Minute)
	tmp.zone = NewTimezoneFromOffset(offset)
	return Undefined
}

func parseDateTimeTimezoneIANAName(formatString, input string, currentInput *string, tmp *tmpDateTime) (err Value) {
	var buffer strings.Builder

inputLoop:
	for len(*currentInput) > 0 {
		char, size := utf8.DecodeRuneInString(*currentInput)
		switch char {
		case '/', '+', '-', '_':
		default:
			if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
				break inputLoop
			}
		}
		*currentInput = (*currentInput)[size:]
		buffer.WriteRune(char)
	}

	timezoneName := buffer.String()
	location, er := time.LoadLocation(timezoneName)
	if er != nil {
		return Ref(Errorf(
			FormatErrorClass,
			"invalid datetime timezone name: %s",
			timezoneName,
		))
	}

	tmp.zone = NewTimezone(location)
	return Undefined
}

func parseDateTimeTimezoneName(formatString, input string, currentInput *string, tmp *tmpDateTime) (err Value) {
	var buffer strings.Builder

	for len(*currentInput) > 0 {
		char, size := utf8.DecodeRuneInString(*currentInput)
		if !unicode.IsLetter(char) {
			break
		}
		*currentInput = (*currentInput)[size:]
		buffer.WriteRune(unicode.ToUpper(char))
	}

	timezoneName := buffer.String()
	offset, ok := tzAbbrevOffsets[timezoneName]
	if !ok {
		return Ref(Errorf(
			FormatErrorClass,
			"invalid datetime timezone abbreviation: %s",
			timezoneName,
		))
	}

	tmp.zone = NewTimezoneFromOffset(offset)
	return Undefined
}

func parseDate1(formatString, input string, currentInput *string, tmp *tmpDateTime) (err Value) {
	err = parseDateAbbreviatedDayOfWeekName(currentInput, &tmp.date)
	if !err.IsUndefined() {
		return err
	}
	err = parseDateTimeMatchText(formatString, input, currentInput, " ")
	if !err.IsUndefined() {
		return err
	}
	err = parseDateAbbreviatedMonthName(currentInput, &tmp.date)
	if !err.IsUndefined() {
		return err
	}
	err = parseDateTimeMatchText(formatString, input, currentInput, " ")
	if !err.IsUndefined() {
		return err
	}
	err = parseDateDayOfMonth(formatString, input, currentInput, &tmp.date, true)
	if !err.IsUndefined() {
		return err
	}
	err = parseDateTimeMatchText(formatString, input, currentInput, " ")
	if !err.IsUndefined() {
		return err
	}
	err = parseTime24Seconds(formatString, input, currentInput, &tmp.time)
	if !err.IsUndefined() {
		return err
	}
	err = parseDateTimeMatchText(formatString, input, currentInput, " ")
	if !err.IsUndefined() {
		return err
	}
	err = parseDateTimeTimezoneName(formatString, input, currentInput, tmp)
	if !err.IsUndefined() {
		return err
	}
	return Undefined
}

func parseDateAndTime(formatString, input string, currentInput *string, tmp *tmpDateTime) (err Value) {
	err = parseDateAbbreviatedDayOfWeekName(currentInput, &tmp.date)
	if !err.IsUndefined() {
		return err
	}
	err = parseDateTimeMatchText(formatString, input, currentInput, " ")
	if !err.IsUndefined() {
		return err
	}
	err = parseDateAbbreviatedMonthName(currentInput, &tmp.date)
	if !err.IsUndefined() {
		return err
	}
	err = parseDateTimeMatchText(formatString, input, currentInput, " ")
	if !err.IsUndefined() {
		return err
	}
	err = parseDateDayOfMonth(formatString, input, currentInput, &tmp.date, true)
	if !err.IsUndefined() {
		return err
	}
	err = parseDateTimeMatchText(formatString, input, currentInput, " ")
	if !err.IsUndefined() {
		return err
	}
	err = parseTime24Seconds(formatString, input, currentInput, &tmp.time)
	if !err.IsUndefined() {
		return err
	}
	err = parseDateYear(formatString, input, currentInput, &tmp.date, false)
	if !err.IsUndefined() {
		return err
	}

	return Undefined
}

// Create a string formatted according to the given format string.
func (t DateTime) Format(formatString string) (_ string, err Value) {
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
		case timescanner.FULL_ISO_YEAR:
			fmt.Fprintf(&buffer, "%d", t.ISOYear())
		case timescanner.FULL_ISO_YEAR_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%4d", t.ISOYear())
		case timescanner.FULL_ISO_YEAR_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%04d", t.ISOYear())
		case timescanner.FULL_YEAR:
			fmt.Fprintf(&buffer, "%d", t.Year())
		case timescanner.FULL_YEAR_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%4d", t.Year())
		case timescanner.FULL_YEAR_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%04d", t.Year())
		case timescanner.CENTURY:
			fmt.Fprintf(&buffer, "%d", t.Century())
		case timescanner.CENTURY_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.Century())
		case timescanner.CENTURY_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.Century())
		case timescanner.YEAR_LAST_TWO_WEEK_BASED:
			fmt.Fprintf(&buffer, "%d", t.ISOYearLastTwo())
		case timescanner.YEAR_LAST_TWO_WEEK_BASED_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.ISOYearLastTwo())
		case timescanner.YEAR_LAST_TWO_WEEK_BASED_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.ISOYearLastTwo())
		case timescanner.YEAR_LAST_TWO:
			fmt.Fprintf(&buffer, "%d", t.YearLastTwo())
		case timescanner.YEAR_LAST_TWO_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.YearLastTwo())
		case timescanner.YEAR_LAST_TWO_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.YearLastTwo())
		case timescanner.MONTH:
			fmt.Fprintf(&buffer, "%d", t.Month())
		case timescanner.MONTH_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.Month())
		case timescanner.MONTH_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.Month())
		case timescanner.MONTH_FULL_NAME:
			buffer.WriteString(t.MonthName())
		case timescanner.MONTH_FULL_NAME_UPPERCASE:
			buffer.WriteString(strings.ToUpper(t.MonthName()))
		case timescanner.MONTH_ABBREVIATED_NAME:
			buffer.WriteString(t.AbbreviatedMonthName())
		case timescanner.MONTH_ABBREVIATED_NAME_UPPERCASE:
			buffer.WriteString(strings.ToUpper(t.AbbreviatedMonthName()))
		case timescanner.DAY_OF_MONTH:
			fmt.Fprintf(&buffer, "%d", t.Day())
		case timescanner.DAY_OF_MONTH_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.Day())
		case timescanner.DAY_OF_MONTH_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.Day())
		case timescanner.DAY_OF_YEAR:
			fmt.Fprintf(&buffer, "%d", t.YearDay())
		case timescanner.DAY_OF_YEAR_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%3d", t.YearDay())
		case timescanner.DAY_OF_YEAR_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%03d", t.YearDay())
		case timescanner.HOUR_OF_DAY:
			fmt.Fprintf(&buffer, "%d", t.Hour())
		case timescanner.HOUR_OF_DAY_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.Hour())
		case timescanner.HOUR_OF_DAY_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.Hour())
		case timescanner.HOUR_OF_DAY12:
			fmt.Fprintf(&buffer, "%d", t.Hour12())
		case timescanner.HOUR_OF_DAY12_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.Hour12())
		case timescanner.HOUR_OF_DAY12_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.Hour12())
		case timescanner.MERIDIEM_INDICATOR_LOWERCASE:
			buffer.WriteString(t.MeridiemLowercase())
		case timescanner.MERIDIEM_INDICATOR_UPPERCASE:
			buffer.WriteString(t.Meridiem())
		case timescanner.MINUTE_OF_HOUR:
			fmt.Fprintf(&buffer, "%d", t.Minute())
		case timescanner.MINUTE_OF_HOUR_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.Minute())
		case timescanner.MINUTE_OF_HOUR_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.Minute())
		case timescanner.SECOND_OF_MINUTE:
			fmt.Fprintf(&buffer, "%d", t.Second())
		case timescanner.SECOND_OF_MINUTE_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.Second())
		case timescanner.SECOND_OF_MINUTE_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.Second())
		case timescanner.MILLISECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d", t.Millisecond())
		case timescanner.MILLISECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%3d", t.Millisecond())
		case timescanner.MILLISECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%03d", t.Millisecond())
		case timescanner.TIMEZONE_IANA_NAME:
			buffer.WriteString(t.ZoneName())
		case timescanner.TIMEZONE_NAME:
			buffer.WriteString(t.ZoneAbbreviatedName())
		case timescanner.TIMEZONE_OFFSET:
			offset := t.ZoneOffset()
			var sign string
			if offset >= 0 {
				sign = "+"
			} else {
				sign = "-"
				offset = -offset
			}

			hours := offset / Hour
			minutes := (offset % Hour) / Minute
			fmt.Fprintf(&buffer, "%s%02d%02d", sign, hours, minutes)
		case timescanner.TIMEZONE_OFFSET_COLON:
			offset := t.ZoneOffset()
			var sign string
			if offset >= 0 {
				sign = "+"
			} else {
				sign = "-"
				offset = -offset
			}

			hours := offset / Hour
			minutes := (offset % Hour) / Minute
			fmt.Fprintf(&buffer, "%s%02d:%02d", sign, hours, minutes)
		case timescanner.DAY_OF_WEEK_FULL_NAME:
			buffer.WriteString(t.WeekdayName())
		case timescanner.DAY_OF_WEEK_FULL_NAME_UPPERCASE:
			buffer.WriteString(strings.ToUpper(t.WeekdayName()))
		case timescanner.DAY_OF_WEEK_ABBREVIATED_NAME:
			buffer.WriteString(t.AbbreviatedWeekdayName())
		case timescanner.DAY_OF_WEEK_ABBREVIATED_NAME_UPPERCASE:
			buffer.WriteString(strings.ToUpper(t.AbbreviatedWeekdayName()))
		case timescanner.DAY_OF_WEEK_NUMBER:
			fmt.Fprintf(&buffer, "%d", t.WeekdayFromMonday())
		case timescanner.DAY_OF_WEEK_NUMBER_ALT:
			fmt.Fprintf(&buffer, "%d", t.WeekdayFromSunday())
		case timescanner.UNIX_SECONDS:
			fmt.Fprintf(&buffer, "%d", t.UnixSeconds())
		case timescanner.UNIX_MILLISECONDS:
			fmt.Fprintf(&buffer, "%d", t.UnixMilliseconds())
		case timescanner.UNIX_MICROSECONDS:
			fmt.Fprintf(&buffer, "%d%06d", t.UnixSeconds(), t.MicrosecondsInSecond())
		case timescanner.UNIX_NANOSECONDS:
			fmt.Fprintf(&buffer, "%d%09d", t.UnixSeconds(), t.NanosecondsInSecond())
		case timescanner.UNIX_PICOSECONDS:
			fmt.Fprintf(&buffer, "%d%012d", t.UnixSeconds(), t.PicosecondsInSecond())
		case timescanner.UNIX_FEMTOSECONDS:
			fmt.Fprintf(&buffer, "%d%015d", t.UnixSeconds(), t.FemtosecondsInSecond())
		case timescanner.UNIX_ATTOSECONDS:
			fmt.Fprintf(&buffer, "%d%018d", t.UnixSeconds(), t.AttosecondsInSecond())
		case timescanner.UNIX_ZEPTOSECONDS:
			fmt.Fprintf(&buffer, "%d%018d000", t.UnixSeconds(), t.AttosecondsInSecond())
		case timescanner.UNIX_YOCTOSECONDS:
			fmt.Fprintf(&buffer, "%d%018d000000", t.UnixSeconds(), t.AttosecondsInSecond())
		case timescanner.ISO_WEEK:
			fmt.Fprintf(&buffer, "%d", t.ISOWeek())
		case timescanner.ISO_WEEK_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.ISOWeek())
		case timescanner.ISO_WEEK_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.ISOWeek())
		case timescanner.WEEK_OF_YEAR:
			fmt.Fprintf(&buffer, "%d", t.WeekFromMonday())
		case timescanner.WEEK_OF_YEAR_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.WeekFromMonday())
		case timescanner.WEEK_OF_YEAR_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.WeekFromMonday())
		case timescanner.WEEK_OF_YEAR_ALT:
			fmt.Fprintf(&buffer, "%d", t.WeekFromSunday())
		case timescanner.WEEK_OF_YEAR_ALT_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.WeekFromSunday())
		case timescanner.WEEK_OF_YEAR_ALT_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.WeekFromSunday())
		case timescanner.MICROSECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d", t.MicrosecondsInSecond())
		case timescanner.MICROSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%6d", t.MicrosecondsInSecond())
		case timescanner.MICROSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%06d", t.MicrosecondsInSecond())
		case timescanner.NANOSECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d", t.NanosecondsInSecond())
		case timescanner.NANOSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%9d", t.NanosecondsInSecond())
		case timescanner.NANOSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%09d", t.NanosecondsInSecond())
		case timescanner.PICOSECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d", t.PicosecondsInSecond())
		case timescanner.PICOSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%12d", t.PicosecondsInSecond())
		case timescanner.PICOSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%012d", t.PicosecondsInSecond())
		case timescanner.FEMTOSECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d", t.FemtosecondsInSecond())
		case timescanner.FEMTOSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%15d", t.FemtosecondsInSecond())
		case timescanner.FEMTOSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%015d", t.FemtosecondsInSecond())
		case timescanner.ATTOSECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d", t.AttosecondsInSecond())
		case timescanner.ATTOSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%18d", t.AttosecondsInSecond())
		case timescanner.ATTOSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%018d", t.AttosecondsInSecond())
		case timescanner.ZEPTOSECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d000", t.AttosecondsInSecond())
		case timescanner.ZEPTOSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%18d000", t.AttosecondsInSecond())
		case timescanner.ZEPTOSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%018d000", t.AttosecondsInSecond())
		case timescanner.YOCTOSECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d000000", t.AttosecondsInSecond())
		case timescanner.YOCTOSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%18d000000", t.AttosecondsInSecond())
		case timescanner.YOCTOSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%018d000000", t.AttosecondsInSecond())
		case timescanner.DATE_AND_TIME:
			fmt.Fprintf(
				&buffer,
				"%s %s %2d %02d:%02d:%02d %04d",
				t.AbbreviatedWeekdayName(),
				t.AbbreviatedMonthName(),
				t.Day(),
				t.Hour(),
				t.Minute(),
				t.Second(),
				t.Year(),
			)
		case timescanner.DATE_AND_TIME_UPPERCASE:
			fmt.Fprintf(
				&buffer,
				"%s %s %2d %02d:%02d:%02d %04d",
				strings.ToUpper(t.AbbreviatedWeekdayName()),
				strings.ToUpper(t.AbbreviatedMonthName()),
				t.Day(),
				t.Hour(),
				t.Minute(),
				t.Second(),
				t.Year(),
			)
		case timescanner.DATE:
			fmt.Fprintf(&buffer, "%02d/%02d/%02d", t.Month(), t.Day(), t.YearLastTwo())
		case timescanner.ISO8601_DATE:
			fmt.Fprintf(&buffer, "%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
		case timescanner.TIME12:
			fmt.Fprintf(
				&buffer,
				"%02d:%02d:%02d %s",
				t.Hour12(),
				t.Minute(),
				t.Second(),
				t.Meridiem(),
			)
		case timescanner.TIME24:
			fmt.Fprintf(&buffer, "%02d:%02d", t.Hour(), t.Minute())
		case timescanner.TIME24_SECONDS:
			fmt.Fprintf(&buffer, "%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
		case timescanner.DATE1_FORMAT:
			fmt.Fprintf(
				&buffer,
				"%s %s %2d %02d:%02d:%02d %s %04d",
				t.AbbreviatedWeekdayName(),
				t.AbbreviatedMonthName(),
				t.Day(),
				t.Hour(),
				t.Minute(),
				t.Second(),
				t.ZoneAbbreviatedName(),
				t.Year(),
			)
		case timescanner.DATE1_FORMAT_UPPERCASE:
			fmt.Fprintf(
				&buffer,
				"%s %s %2d %02d:%02d:%02d %s %04d",
				strings.ToUpper(t.AbbreviatedWeekdayName()),
				strings.ToUpper(t.AbbreviatedMonthName()),
				t.Day(),
				t.Hour(),
				t.Minute(),
				t.Second(),
				t.ZoneAbbreviatedName(),
				t.Year(),
			)
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

// Check whether t is greater than other and return an error
// if something went wrong.
func (t *DateTime) GreaterThan(other Value) (result bool, err Value) {
	switch other.flag {
	case DATE_FLAG:
		return t.Cmp(other.AsDate().ToDateTime()) == 1, err
	case REFERENCE_FLAG:
	default:
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}

	if !other.IsReference() {
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case *DateTime:
		return t.Cmp(o) == 1, err
	default:
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
}

func (t *DateTime) GreaterThanVal(other Value) (Value, Value) {
	result, err := t.GreaterThan(other)
	return ToElkBool(result), err
}

// Check whether t is greater than or equal to other and return an error
// if something went wrong.
func (t *DateTime) GreaterThanEqual(other Value) (result bool, err Value) {
	switch other.flag {
	case DATE_FLAG:
		return t.Cmp(other.AsDate().ToDateTime()) >= 0, err
	case REFERENCE_FLAG:
	default:
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}

	if !other.IsReference() {
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case *DateTime:
		return t.Cmp(o) >= 0, err
	default:
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
}

func (t *DateTime) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := t.GreaterThanEqual(other)
	return ToElkBool(result), err
}

// Check whether t is less than other and return an error
// if something went wrong.
func (t *DateTime) LessThan(other Value) (result bool, err Value) {
	switch other.flag {
	case DATE_FLAG:
		return t.Cmp(other.AsDate().ToDateTime()) == -1, err
	case REFERENCE_FLAG:
	default:
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}

	if !other.IsReference() {
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case *DateTime:
		return t.Cmp(o) == -1, err
	default:
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
}

func (t *DateTime) LessThanVal(other Value) (Value, Value) {
	result, err := t.LessThan(other)
	return ToElkBool(result), err
}

// Check whether t is less than or equal to other and return an error
// if something went wrong.
func (t *DateTime) LessThanEqual(other Value) (result bool, err Value) {
	switch other.flag {
	case DATE_FLAG:
		return t.Cmp(other.AsDate().ToDateTime()) <= 0, err
	case REFERENCE_FLAG:
	default:
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}

	if !other.IsReference() {
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case *DateTime:
		return t.Cmp(o) <= 0, err
	default:
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
}

func (t *DateTime) LessThanEqualVal(other Value) (Value, Value) {
	result, err := t.LessThanEqual(other)
	return ToElkBool(result), err
}

func (t *DateTime) LaxEqual(other Value) bool {
	return t.Equal(other)
}

// Check whether t is equal to other and return an error
// if something went wrong.
func (t *DateTime) Equal(other Value) bool {
	if !other.IsReference() {
		return false
	}
	switch o := other.AsReference().(type) {
	case *DateTime:
		return t.Cmp(o) == 0
	default:
		return false
	}
}

func (t *DateTime) StrictEqual(other Value) bool {
	return t.Equal(other)
}

func initDateTime() {
	DateTimeClass = NewClass()
	StdModule.AddConstantString("DateTime", Ref(DateTimeClass))
	DateTimeClass.AddConstantString("DEFAULT_FORMAT", Ref(String(DefaultDateTimeFormat)))
}
