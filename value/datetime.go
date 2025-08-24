package value

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/elk-language/elk/value/timescanner"
)

var DateTimeClass *Class // ::Std::DateTime

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

const DefaultTimeFormat = "%Y-%m-%d %H:%M:%S.%9N %:z"

func (t DateTime) Inspect() string {
	return fmt.Sprintf("Std::DateTime('%s')", t.ToString().String())
}

func (t DateTime) Error() string {
	return t.Inspect()
}

func (t DateTime) InstanceVariables() *InstanceVariables {
	return nil
}

func ToElkTime(time time.Time) *DateTime {
	return &DateTime{Go: time}
}

func ToElkTimeValue(time time.Time) DateTime {
	return DateTime{Go: time}
}

func (t DateTime) ToString() String {
	return String(t.String())
}

func (t DateTime) String() string {
	return t.MustFormat(DefaultTimeFormat)
}

// Create a new Time object.
func NewTime(year, month, day, hour, min, sec, nsec int, zone *Timezone) *DateTime {
	t := MakeTime(year, month, day, hour, min, sec, nsec, zone)
	return &t
}

// Create a new Time value.
func MakeTime(year, month, day, hour, min, sec, nsec int, zone *Timezone) DateTime {
	var location *time.Location
	if zone == nil {
		location = time.UTC
	} else {
		location = zone.ToGoLocation()
	}
	return DateTime{
		Go: time.Date(year, time.Month(month), day, hour, min, sec, nsec, location),
	}
}

func TimeNow() *DateTime {
	return ToElkTime(time.Now())
}

// Adds the given duration to the time.
// Returns a new time structure.
func (t *DateTime) Add(val Duration) *DateTime {
	return ToElkTime(t.Go.Add(val.Go()))
}

// Subtracts the given duration from the time.
// Returns a new time structure.
func (t *DateTime) Subtract(val Duration) *DateTime {
	return ToElkTime(t.Go.Add(-val.Go()))
}

// Calculates the difference between two time objects.
// Returns a duration.
func (t *DateTime) Diff(val *DateTime) Duration {
	return Duration(t.Go.Sub(val.Go))
}

func (t DateTime) ToGoTime() time.Time {
	return t.Go
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
	return t.Nanosecond() / 1000_000
}

func (t DateTime) Microsecond() int {
	return t.Nanosecond() / 1000
}

func (t DateTime) Nanosecond() int {
	return t.Go.Nanosecond()
}

func (t DateTime) Picosecond() int64 {
	return int64(t.Nanosecond()) * 1000
}

func (t DateTime) Femtosecond() int64 {
	return int64(t.Nanosecond()) * 1000_000
}

func (t DateTime) Attosecond() int64 {
	return int64(t.Nanosecond()) * 1000_000_000
}

func (t DateTime) Zeptosecond() *big.Int {
	i := big.NewInt(int64(t.Nanosecond()))
	i.Mul(i, big.NewInt(1000_000_000_000))
	return i
}

func (t DateTime) Yoctosecond() *big.Int {
	i := big.NewInt(int64(t.Nanosecond()))
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
	return i.Add(i, big.NewInt(int64(t.Microsecond())))
}

func (t DateTime) UnixNanoseconds() *big.Int {
	i := big.NewInt(t.UnixSeconds())
	i = i.Mul(i, big.NewInt(1000_000_000))
	return i.Add(i, big.NewInt(int64(t.Nanosecond())))
}

func (t DateTime) UnixPicoseconds() *big.Int {
	i := big.NewInt(t.UnixSeconds())
	i = i.Mul(i, big.NewInt(1000_000_000_000))
	return i.Add(i, big.NewInt(t.Picosecond()))
}

func (t DateTime) UnixFemtoseconds() *big.Int {
	i := big.NewInt(t.UnixSeconds())
	i = i.Mul(i, big.NewInt(1000_000_000_000_000))
	return i.Add(i, big.NewInt(t.Femtosecond()))
}

func (t DateTime) UnixAttoseconds() *big.Int {
	i := big.NewInt(t.UnixSeconds())
	i = i.Mul(i, big.NewInt(1000_000_000_000_000_000))
	return i.Add(i, big.NewInt(t.Attosecond()))
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
	return ToElkTime(t.Go.UTC())
}

// Convert the time to the local timezone.
func (t *DateTime) Local() *DateTime {
	return ToElkTime(t.Go.Local())
}

// Cmp compares x and y and returns:
//
//	  -1 if x <  y
//		 0 if x == y
//	  +1 if x >  y
func (x *DateTime) Cmp(y *DateTime) int {
	return x.Go.Compare(y.Go)
}

func (t DateTime) MustFormat(formatString string) string {
	result, err := t.Format(formatString)
	if !err.IsUndefined() {
		panic(err)
	}

	return result
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
		case timescanner.FULL_YEAR_WEEK_BASED:
			fmt.Fprintf(&buffer, "%d", t.ISOYear())
		case timescanner.FULL_YEAR_WEEK_BASED_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%4d", t.ISOYear())
		case timescanner.FULL_YEAR_WEEK_BASED_ZERO_PADDED:
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
		case timescanner.TIMEZONE_NAME:
			buffer.WriteString(t.ZoneAbbreviatedName())
		case timescanner.TIMEZONE_OFFSET:
			hours := t.ZoneOffsetHours()
			minutes := t.ZoneOffsetHourMinutes()
			var sign string
			if hours >= 0 {
				sign = "+"
			} else {
				sign = "-"
			}
			fmt.Fprintf(&buffer, "%s%02d%02d", sign, hours, minutes)
		case timescanner.TIMEZONE_OFFSET_COLON:
			hours := t.ZoneOffsetHours()
			minutes := t.ZoneOffsetHourMinutes()
			var sign string
			if hours >= 0 {
				sign = "+"
			} else {
				sign = "-"
			}
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
			fmt.Fprintf(&buffer, "%d%06d", t.UnixSeconds(), t.Microsecond())
		case timescanner.UNIX_NANOSECONDS:
			fmt.Fprintf(&buffer, "%d%09d", t.UnixSeconds(), t.Nanosecond())
		case timescanner.UNIX_PICOSECONDS:
			fmt.Fprintf(&buffer, "%d%012d", t.UnixSeconds(), t.Picosecond())
		case timescanner.UNIX_FEMTOSECONDS:
			fmt.Fprintf(&buffer, "%d%015d", t.UnixSeconds(), t.Femtosecond())
		case timescanner.UNIX_ATTOSECONDS:
			fmt.Fprintf(&buffer, "%d%018d", t.UnixSeconds(), t.Attosecond())
		case timescanner.UNIX_ZEPTOSECONDS:
			fmt.Fprintf(&buffer, "%d%018d000", t.UnixSeconds(), t.Attosecond())
		case timescanner.UNIX_YOCTOSECONDS:
			fmt.Fprintf(&buffer, "%d%018d000000", t.UnixSeconds(), t.Attosecond())
		case timescanner.WEEK_OF_WEEK_BASED_YEAR:
			fmt.Fprintf(&buffer, "%d", t.ISOWeek())
		case timescanner.WEEK_OF_WEEK_BASED_YEAR_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.ISOWeek())
		case timescanner.WEEK_OF_WEEK_BASED_YEAR_ZERO_PADDED:
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
			fmt.Fprintf(&buffer, "%d", t.Microsecond())
		case timescanner.MICROSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%6d", t.Microsecond())
		case timescanner.MICROSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%06d", t.Microsecond())
		case timescanner.NANOSECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d", t.Nanosecond())
		case timescanner.NANOSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%9d", t.Nanosecond())
		case timescanner.NANOSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%09d", t.Nanosecond())
		case timescanner.PICOSECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d", t.Picosecond())
		case timescanner.PICOSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%12d", t.Picosecond())
		case timescanner.PICOSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%012d", t.Picosecond())
		case timescanner.FEMTOSECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d", t.Femtosecond())
		case timescanner.FEMTOSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%15d", t.Femtosecond())
		case timescanner.FEMTOSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%015d", t.Femtosecond())
		case timescanner.ATTOSECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d", t.Attosecond())
		case timescanner.ATTOSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%18d", t.Attosecond())
		case timescanner.ATTOSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%018d", t.Attosecond())
		case timescanner.ZEPTOSECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d000", t.Attosecond())
		case timescanner.ZEPTOSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%18d000", t.Attosecond())
		case timescanner.ZEPTOSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%018d000", t.Attosecond())
		case timescanner.YOCTOSECOND_OF_SECOND:
			fmt.Fprintf(&buffer, "%d000000", t.Attosecond())
		case timescanner.YOCTOSECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%18d000000", t.Attosecond())
		case timescanner.YOCTOSECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%018d000000", t.Attosecond())
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
func (t *DateTime) GreaterThan(other Value) (result Value, err Value) {
	if !other.IsReference() {
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case *DateTime:
		return ToElkBool(t.Cmp(o) == 1), err
	default:
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
}

// Check whether t is greater than or equal to other and return an error
// if something went wrong.
func (t *DateTime) GreaterThanEqual(other Value) (result Value, err Value) {
	if !other.IsReference() {
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case *DateTime:
		return ToElkBool(t.Cmp(o) >= 0), err
	default:
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
}

// Check whether t is less than other and return an error
// if something went wrong.
func (t *DateTime) LessThan(other Value) (result Value, err Value) {
	if !other.IsReference() {
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case *DateTime:
		return ToElkBool(t.Cmp(o) == -1), err
	default:
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
}

// Check whether t is less than or equal to other and return an error
// if something went wrong.
func (t *DateTime) LessThanEqual(other Value) (result Value, err Value) {
	if !other.IsReference() {
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case *DateTime:
		return ToElkBool(t.Cmp(o) <= 0), err
	default:
		return result, Ref(NewCoerceError(t.Class(), other.Class()))
	}
}

func (t *DateTime) LaxEqual(other Value) Value {
	return t.Equal(other)
}

// Check whether t is equal to other and return an error
// if something went wrong.
func (t *DateTime) Equal(other Value) Value {
	if !other.IsReference() {
		return False
	}
	switch o := other.AsReference().(type) {
	case *DateTime:
		return ToElkBool(t.Cmp(o) == 0)
	default:
		return False
	}
}

func (t *DateTime) StrictEqual(other Value) Value {
	return t.Equal(other)
}

func initDateTime() {
	DateTimeClass = NewClass()
	StdModule.AddConstantString("DateTime", Ref(DateTimeClass))
	DateTimeClass.AddConstantString("DEFAULT_FORMAT", Ref(String(DefaultTimeFormat)))
}
