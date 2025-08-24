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
	duration Duration
}

var TimeClass *Class // ::Std::Time

func TimeNow() Time {
	t := time.Now()
	return MakeTime(t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}

// Create a new Time value.
func MakeTime(hour, min, sec, nsec int) Time {
	duration := Duration(hour)*Hour +
		Duration(min)*Minute +
		Duration(sec)*Second +
		Duration(nsec)*Nanosecond

	return Time{
		duration: duration,
	}
}

func (t Time) ToDuration() Duration {
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
	return t.duration.HoursMod()
}

func (t Time) Minute() int {
	return t.duration.MinutesMod()
}

func (t Time) Second() int {
	return t.duration.SecondsMod()
}

func (t Time) MillisecondsInSecond() int {
	return int(t.duration % Second / Millisecond)
}

func (t Time) Millisecond() int {
	return t.duration.MillisecondsMod()
}

func (t Time) MicrosecondsInSecond() int {
	return int(t.duration % Second / Microsecond)
}

func (t Time) Microsecond() int {
	return t.duration.MicrosecondsMod()
}

func (t Time) NanosecondsInSecond() int {
	return int(t.duration % Second)
}

func (t Time) Nanosecond() int {
	return t.duration.NanosecondsMod()
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
