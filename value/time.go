package value

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/elk-language/elk/value/timescanner"
)

var TimeClass *Class // ::Std::Time

// Elk's Time value
type Time struct {
	Go time.Time
}

func (Time) Class() *Class {
	return TimeClass
}

func (Time) DirectClass() *Class {
	return TimeClass
}

func (Time) SingletonClass() *Class {
	return nil
}

func (t Time) Inspect() string {
	return fmt.Sprintf("Time('%s')", t.Go.Format(time.RFC3339Nano))
}

func (t Time) InstanceVariables() SymbolMap {
	return nil
}

func ToElkTime(time time.Time) Time {
	return Time{Go: time}
}

// Create a new Time object.
func NewTime(year, month, day, hour, min, sec, nsec int, zone *Timezone) Time {
	var location *time.Location
	if zone == nil {
		location = time.UTC
	} else {
		location = zone.ToGoLocation()
	}
	return Time{
		Go: time.Date(year, time.Month(month), day, hour, min, sec, nsec, location),
	}
}

func TimeNow() Time {
	return ToElkTime(time.Now())
}

func (t Time) ToGoTime() time.Time {
	return t.Go
}

func (t Time) Zone() *Timezone {
	return NewTimezone(t.Go.Location())
}

func (t Time) WeekBasedYear() int {
	year, _ := t.Go.ISOWeek()
	return year
}

func (t Time) WeekBasedYearLastTwo() int {
	return t.WeekBasedYear() % 100
}

func (t Time) YearLastTwo() int {
	return t.Go.Year() % 100
}

func (t Time) Year() int {
	return t.Go.Year()
}

func (t Time) Century() int {
	return t.Go.Year() / 100
}

func (t Time) Month() int {
	return int(t.Go.Month())
}

func (t Time) MonthName() string {
	return t.Go.Month().String()
}

func (t Time) AbbreviatedMonthName() string {
	return t.MonthName()[0:3]
}

func (t Time) Day() int {
	return t.Go.Day()
}

// Day of the year.
func (t Time) YearDay() int {
	return t.Go.YearDay()
}

// Hour in a 24 hour clock.
func (t Time) Hour() int {
	return t.Go.Hour()
}

// Whether the current hour is AM.
func (t Time) IsAM() bool {
	hour := t.Hour()

	return hour < 12
}

func (t Time) Meridiem() string {
	if t.IsAM() {
		return "AM"
	}

	return "PM"
}

func (t Time) MeridiemLowercase() string {
	if t.IsAM() {
		return "am"
	}

	return "pm"
}

// Hour in a twelve hour clock.
func (t Time) Hour12() int {
	hour := t.Hour()
	if hour == 0 {
		return 12
	}

	if hour <= 12 {
		return hour
	}

	return hour - 12
}

func (t Time) Minute() int {
	return t.Go.Minute()
}

func (t Time) Second() int {
	return t.Go.Second()
}

func (t Time) Millisecond() int {
	return t.Nanosecond() / 1000_000
}

func (t Time) Microsecond() int {
	return t.Nanosecond() / 1000
}

func (t Time) Nanosecond() int {
	return t.Go.Nanosecond()
}

func (t Time) Picosecond() int64 {
	return int64(t.Nanosecond()) * 1000
}

func (t Time) Femtosecond() int64 {
	return int64(t.Nanosecond()) * 1000_000
}

func (t Time) Attosecond() int64 {
	return int64(t.Nanosecond()) * 1000_000_000
}

func (t Time) Zeptosecond() *big.Int {
	i := big.NewInt(int64(t.Nanosecond()))
	i.Mul(i, big.NewInt(1000_000_000_000))
	return i
}

func (t Time) Yoctosecond() *big.Int {
	i := big.NewInt(int64(t.Nanosecond()))
	i.Mul(i, big.NewInt(1000_000_000_000_000))
	return i
}

func (t Time) ZoneName() string {
	return t.Go.Location().String()
}

func (t Time) ZoneAbbreviatedName() string {
	name, _ := t.Go.Zone()
	return name
}

func (t Time) ZoneOffsetSeconds() int {
	_, offset := t.Go.Zone()
	return offset
}

func (t Time) ZoneOffsetHours() int {
	_, offset := t.Go.Zone()
	return offset / 3600
}

func (t Time) ZoneOffsetHourMinutes() int {
	_, offset := t.Go.Zone()
	return (offset % 3600) / 60
}

func (t Time) WeekdayName() string {
	return t.Go.Weekday().String()
}

func (t Time) AbbreviatedWeekdayName() string {
	return t.WeekdayName()[0:3]
}

// Specifies the day of the week (Monday = 1, ...).
func (t Time) Weekday() int {
	weekday := int(t.Go.Weekday())
	if weekday == 0 {
		return 7
	}

	return weekday
}

// Specifies the day of the week (Sunday = 0, ...).
func (t Time) WeekdayAlt() int {
	return int(t.Go.Weekday())
}

func (t Time) UnixSeconds() int {
	return int(t.Go.Unix())
}

func (t Time) UnixMilliseconds() int {
	return int(t.Go.UnixMilli())
}

func (t Time) ISOWeek() int {
	_, week := t.Go.ISOWeek()
	return week
}

func (t Time) weekNumber(firstWeekday int) int {
	yday := t.YearDay()
	wday := t.WeekdayAlt()

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
func (t Time) WeekMonday() int {
	return t.weekNumber(1)
}

// The week number of the current year as a decimal number,
// range 00 to 53, starting with the first Sunday
// as the first day of week 01.
func (t Time) WeekSunday() int {
	return t.weekNumber(0)
}

// Create a string formatted according to the given format string.
func (t Time) Format(formatString string) (string, *Error) {
	scanner := timescanner.New(formatString)
	var buffer strings.Builder

tokenLoop:
	for {
		token, value := scanner.Next()
		switch token {
		case timescanner.END_OF_FILE:
			break tokenLoop
		case timescanner.INVALID_FORMAT_DIRECTIVE:
			return "", Errorf(
				FormatErrorClass,
				"invalid format directive: %s",
				value,
			)
		case timescanner.PERCENT:
			buffer.WriteByte('%')
		case timescanner.NEWLINE:
			buffer.WriteByte('\n')
		case timescanner.TAB:
			buffer.WriteByte('\t')
		case timescanner.TEXT:
			buffer.WriteString(value)
		case timescanner.FULL_YEAR_WEEK_BASED:
			fmt.Fprintf(&buffer, "%d", t.WeekBasedYear())
		case timescanner.FULL_YEAR_WEEK_BASED_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%4d", t.WeekBasedYear())
		case timescanner.FULL_YEAR_WEEK_BASED_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%04d", t.WeekBasedYear())
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
			fmt.Fprintf(&buffer, "%d", t.WeekBasedYearLastTwo())
		case timescanner.YEAR_LAST_TWO_WEEK_BASED_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.WeekBasedYearLastTwo())
		case timescanner.YEAR_LAST_TWO_WEEK_BASED_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.WeekBasedYearLastTwo())
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
			fmt.Fprintf(&buffer, "%d", t.Weekday())
		case timescanner.DAY_OF_WEEK_NUMBER_ALT:
			fmt.Fprintf(&buffer, "%d", t.WeekdayAlt())
		case timescanner.UNIX_SECONDS:
			fmt.Fprintf(&buffer, "%d", t.UnixSeconds())
		case timescanner.UNIX_MILLISECONDS:
			fmt.Fprintf(&buffer, "%d", t.UnixMilliseconds())
		case timescanner.WEEK_OF_WEEK_BASED_YEAR:
			fmt.Fprintf(&buffer, "%d", t.ISOWeek())
		case timescanner.WEEK_OF_WEEK_BASED_YEAR_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.ISOWeek())
		case timescanner.WEEK_OF_WEEK_BASED_YEAR_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.ISOWeek())
		case timescanner.WEEK_OF_YEAR:
			fmt.Fprintf(&buffer, "%d", t.WeekMonday())
		case timescanner.WEEK_OF_YEAR_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.WeekMonday())
		case timescanner.WEEK_OF_YEAR_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.WeekMonday())
		case timescanner.WEEK_OF_YEAR_ALT:
			fmt.Fprintf(&buffer, "%d", t.WeekSunday())
		case timescanner.WEEK_OF_YEAR_ALT_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%2d", t.WeekSunday())
		case timescanner.WEEK_OF_YEAR_ALT_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%02d", t.WeekSunday())
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
			return "", Errorf(
				FormatErrorClass,
				"unsupported format directive: %s",
				token.String(),
			)
		}

	}

	return buffer.String(), nil
}

func initTime() {
	TimeClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Time", TimeClass)
}
