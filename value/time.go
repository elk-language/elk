package value

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/elk-language/elk/value/timescanner"
)

const DefaultTimeFormat = "%H:%M:%S.%9N"

// Represents a time of day: hour, minute, second, nanosecond.
type Time struct {
	duration TimeSpan
}

var TimeClass *Class // ::Std::Time

func TimeNow() Time {
	t := time.Now()
	return MakeTime(t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}

// Create a new Time value.
func MakeTime(hour, min, sec, nsec int) Time {
	duration := TimeSpan(hour)*Hour +
		TimeSpan(min)*Minute +
		TimeSpan(sec)*Second +
		TimeSpan(nsec)*Nanosecond

	duration %= Day

	return Time{
		duration: duration,
	}
}

func (t Time) ToTimeSpan() TimeSpan {
	return t.duration
}

func (t Time) ToDateTime() *DateTime {
	d := t.ToDateTimeValue()
	return &d
}

func (t Time) ToDateTimeValue() DateTime {
	goTime := time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
	return ToElkDateTimeValue(goTime)
}

func (t Time) Copy() Reference {
	return t
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
	return fmt.Sprintf("Std::Time.parse('%s')", t.String())
}

func (Time) InstanceVariables() *InstanceVariables {
	return nil
}

func (t Time) Error() string {
	return t.Inspect()
}

func (t Time) ToString() String {
	return String(t.String())
}

func (t Time) String() string {
	return t.MustFormat(DefaultTimeFormat)
}

// Hour in a 24 hour clock.
func (t Time) Hour() int {
	return t.duration.Hours()
}

func (t Time) Minute() int {
	return t.duration.Minutes()
}

func (t Time) Second() int {
	return t.duration.Seconds()
}

func (t Time) MillisecondsInSecond() int {
	return int(t.duration % Second / Millisecond)
}

func (t Time) Millisecond() int {
	return t.duration.Milliseconds()
}

func (t Time) MicrosecondsInSecond() int {
	return int(t.duration % Second / Microsecond)
}

func (t Time) Microsecond() int {
	return t.duration.Microseconds()
}

func (t Time) NanosecondsInSecond() int {
	return int(t.duration % Second)
}

func (t Time) Nanosecond() int {
	return t.duration.Nanoseconds()
}

func (t Time) PicosecondsInSecond() int64 {
	return int64(t.NanosecondsInSecond()) * 1000
}

func (t Time) FemtosecondsInSecond() int64 {
	return int64(t.NanosecondsInSecond()) * 1000_000
}

func (t Time) AttosecondsInSecond() int64 {
	return int64(t.NanosecondsInSecond()) * 1000_000_000
}

func (t Time) ZeptosecondsInSecond() *BigInt {
	i := big.NewInt(int64(t.Nanosecond()))
	i.Mul(i, big.NewInt(1000_000_000_000))
	return ToElkBigInt(i)
}

func (t Time) YoctosecondsInSecond() *BigInt {
	i := big.NewInt(int64(t.Nanosecond()))
	i.Mul(i, big.NewInt(1000_000_000_000_000))
	return ToElkBigInt(i)
}

// Whether the current hour is AM.
func (t Time) IsAM() bool {
	hour := t.Hour()

	return hour < 12
}

// Whether the current hour is PM.
func (t Time) IsPM() bool {
	hour := t.Hour()

	return hour >= 12
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

func (t Time) MustFormat(formatString string) string {
	result, err := t.Format(formatString)
	if !err.IsUndefined() {
		panic(err)
	}

	return result
}

func parseDigits(s string, maxChars int, spacePadded bool) (int, string) {
	if len(s) == 0 {
		return -1, ""
	}

	n := 0
	i := 0

	if spacePadded {
		for ; i < len(s) && i < maxChars && s[i] == ' '; i++ {
		}
	}

	for {
		if i >= len(s) || i >= maxChars {
			break
		}
		ch := s[i]
		if ch < '0' || ch > '9' {
			break
		}
		n = n*10 + int(ch-'0')
		i++
	}

	if i == 0 {
		return -1, ""
	}

	return n, s[i:]
}

// Create a string formatted according to the given format string.
func ParseTime(formatString, input string) (result Time, err Value) {
	scanner := timescanner.New(formatString)
	var buffer strings.Builder
	currentInput := input

tokenLoop:
	for {
		token, value := scanner.Next()
		switch token {
		case timescanner.END_OF_FILE:
			if len(currentInput) > 0 {
				return Time{}, Ref(NewIncompatibleTimeFormatError(formatString, input))
			}
			break tokenLoop
		case timescanner.INVALID_FORMAT_DIRECTIVE:
			return Time{}, Ref(Errorf(
				FormatErrorClass,
				"invalid time format directive: %s",
				value,
			))
		case timescanner.PERCENT:
			err = parseTimeMatchText(formatString, input, &currentInput, "%")
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.NEWLINE:
			err = parseTimeMatchText(formatString, input, &currentInput, "\n")
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.TAB:
			err = parseTimeMatchText(formatString, input, &currentInput, "\t")
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.TEXT:
			err = parseTimeMatchText(formatString, input, &currentInput, value)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.HOUR_OF_DAY, timescanner.HOUR_OF_DAY_ZERO_PADDED:
			err = parseTimeHour(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.HOUR_OF_DAY_SPACE_PADDED:
			err = parseTimeHour(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.HOUR_OF_DAY12, timescanner.HOUR_OF_DAY12_ZERO_PADDED:
			err = parseTime12Hour(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.HOUR_OF_DAY12_SPACE_PADDED:
			err = parseTime12Hour(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.MERIDIEM_INDICATOR_LOWERCASE, timescanner.MERIDIEM_INDICATOR_UPPERCASE:
			err = parseTimeMeridiem(formatString, input, &currentInput, &result)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.MINUTE_OF_HOUR, timescanner.MINUTE_OF_HOUR_ZERO_PADDED:
			err = parseTimeMinute(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.MINUTE_OF_HOUR_SPACE_PADDED:
			err = parseTimeMinute(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.SECOND_OF_MINUTE, timescanner.SECOND_OF_MINUTE_ZERO_PADDED:
			err = parseTimeSecond(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.SECOND_OF_MINUTE_SPACE_PADDED:
			err = parseTimeSecond(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.MILLISECOND_OF_SECOND, timescanner.MILLISECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeMillisecond(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.MILLISECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeMillisecond(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.MICROSECOND_OF_SECOND, timescanner.MICROSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeMicrosecond(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.MICROSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeMicrosecond(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.NANOSECOND_OF_SECOND, timescanner.NANOSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeNanosecond(formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.NANOSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeNanosecond(formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.PICOSECOND_OF_SECOND, timescanner.PICOSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeSubNanosecond("picosecond", formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.PICOSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeSubNanosecond("picosecond", formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.FEMTOSECOND_OF_SECOND, timescanner.FEMTOSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeSubNanosecond("femtosecond", formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.FEMTOSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeSubNanosecond("femtosecond", formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.ATTOSECOND_OF_SECOND, timescanner.ATTOSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeSubNanosecond("attosecond", formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.ATTOSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeSubNanosecond("attosecond", formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.ZEPTOSECOND_OF_SECOND, timescanner.ZEPTOSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeSubNanosecond("zeptosecond", formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.ZEPTOSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeSubNanosecond("zeptosecond", formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.YOCTOSECOND_OF_SECOND, timescanner.YOCTOSECOND_OF_SECOND_ZERO_PADDED:
			err = parseTimeSubNanosecond("yoctosecond", formatString, input, &currentInput, &result, false)
			if !err.IsUndefined() {
				return Time{}, err
			}
		case timescanner.YOCTOSECOND_OF_SECOND_SPACE_PADDED:
			err = parseTimeSubNanosecond("yoctosecond", formatString, input, &currentInput, &result, true)
			if !err.IsUndefined() {
				return Time{}, err
			}
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

func parseTimeHour(formatString, input string, currentInput *string, result *Time, spacePadded bool) Value {
	var n int
	n, *currentInput = parseDigits(*currentInput, 2, spacePadded)
	if n < 0 {
		return Ref(NewIncompatibleTimeFormatError(formatString, input))
	}
	if n > 23 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for hour out of range: %d",
			n,
		))
	}

	result.duration += TimeSpan(n) * Hour
	return Undefined
}

func parseTimeMinute(formatString, input string, currentInput *string, result *Time, spacePadded bool) Value {
	var n int
	n, *currentInput = parseDigits(*currentInput, 2, spacePadded)
	if n < 0 {
		return Ref(NewIncompatibleTimeFormatError(formatString, input))
	}
	if n > 59 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for minute out of range: %d",
			n,
		))
	}

	result.duration += TimeSpan(n) * Minute
	return Undefined
}

func parseTimeSecond(formatString, input string, currentInput *string, result *Time, spacePadded bool) Value {
	var n int
	n, *currentInput = parseDigits(*currentInput, 2, spacePadded)
	if n < 0 {
		return Ref(NewIncompatibleTimeFormatError(formatString, input))
	}
	if n > 59 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for second out of range: %d",
			n,
		))
	}

	result.duration += TimeSpan(n) * Second
	return Undefined
}

func parseTimeMillisecond(formatString, input string, currentInput *string, result *Time, spacePadded bool) Value {
	var n int
	n, *currentInput = parseDigits(*currentInput, 3, spacePadded)
	if n < 0 {
		return Ref(NewIncompatibleTimeFormatError(formatString, input))
	}
	if n > 999 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for millisecond out of range: %d",
			n,
		))
	}

	result.duration += TimeSpan(n) * Millisecond
	return Undefined
}

func parseTimeMicrosecond(formatString, input string, currentInput *string, result *Time, spacePadded bool) Value {
	var n int
	n, *currentInput = parseDigits(*currentInput, 3, spacePadded)
	if n < 0 {
		return Ref(NewIncompatibleTimeFormatError(formatString, input))
	}
	if n > 999 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for microsecond out of range: %d",
			n,
		))
	}

	result.duration += TimeSpan(n) * Microsecond
	return Undefined
}

func parseTimeNanosecond(formatString, input string, currentInput *string, result *Time, spacePadded bool) Value {
	var n int
	n, *currentInput = parseDigits(*currentInput, 3, spacePadded)
	if n < 0 {
		return Ref(NewIncompatibleTimeFormatError(formatString, input))
	}
	if n > 999 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for nanosecond out of range: %d",
			n,
		))
	}

	result.duration += TimeSpan(n) * Nanosecond
	return Undefined
}

func parseTimeMatchText(formatString, input string, currentInput *string, text string) Value {
	if len(*currentInput) < len(text) {
		return Ref(NewIncompatibleTimeFormatError(formatString, input))
	}
	got := (*currentInput)[:len(text)]
	if got != text {
		return Ref(Errorf(
			FormatErrorClass,
			"cannot parse time string, expected `%s`, got `%s`",
			String(text).Inspect(),
			String(got).Inspect(),
		))
	}
	*currentInput = (*currentInput)[len(text):]

	return Undefined
}

func parseTimeSubNanosecond(name, formatString, input string, currentInput *string, result *Time, spacePadded bool) Value {
	var n int
	n, *currentInput = parseDigits(*currentInput, 3, spacePadded)
	if n < 0 {
		return Ref(NewIncompatibleTimeFormatError(formatString, input))
	}
	if n > 999 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for %s out of range: %d",
			name,
			n,
		))
	}

	return Undefined
}

func parseTime12Hour(formatString, input string, currentInput *string, result *Time, spacePadded bool) Value {
	var n int
	n, *currentInput = parseDigits(*currentInput, 2, false)
	if n == -1 {
		return Ref(NewIncompatibleTimeFormatError(formatString, input))
	}
	if n < 1 || n > 12 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for 12 hour out of range: %d",
			n,
		))
	}

	result.duration += TimeSpan(n) * Hour
	return Undefined
}

func parseTimeMeridiem(formatString, input string, currentInput *string, result *Time) Value {
	hours := result.duration / Hour % 24
	result.duration -= hours * Hour

	if hours < 1 || hours > 12 {
		return Ref(Errorf(
			FormatErrorClass,
			"value for 12 hour out of range: %d",
			hours,
		))
	}

	if len(*currentInput) < 2 {
		return Ref(NewIncompatibleTimeFormatError(formatString, input))
	}

	meridiem := (*currentInput)[:2]
	*currentInput = (*currentInput)[2:]

	switch meridiem {
	case "am", "AM":
		if hours == 12 {
			hours = 0
		}
	case "pm", "PM":
		if hours != 12 {
			hours += 12
		}
	default:
		return Ref(Errorf(
			FormatErrorClass,
			"invalid meridiem indicator: %s, expected `\"am\"` or `\"pm\"`",
			meridiem,
		))
	}

	result.duration += hours * Hour
	return Undefined
}

// Create a string formatted according to the given format string.
func (t Time) Format(formatString string) (_ string, err Value) {
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
				"invalid time format directive: %s",
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
			fmt.Fprintf(&buffer, "%d", t.MillisecondsInSecond())
		case timescanner.MILLISECOND_OF_SECOND_SPACE_PADDED:
			fmt.Fprintf(&buffer, "%3d", t.MillisecondsInSecond())
		case timescanner.MILLISECOND_OF_SECOND_ZERO_PADDED:
			fmt.Fprintf(&buffer, "%03d", t.MillisecondsInSecond())
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

func initTime() {
	TimeClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	StdModule.AddConstantString("Time", Ref(TimeClass))
	TimeClass.AddConstantString("DEFAULT_FORMAT", Ref(String(DefaultTimeFormat)))
}
