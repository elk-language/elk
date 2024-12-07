package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Time
func initTime() {
	// Class methods
	c := &value.TimeClass.SingletonClass().MethodContainer
	Def(
		c,
		"now",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.Ref(value.TimeNow()), value.Undefined
		},
	)

	// Instance methods
	c = &value.TimeClass.MethodContainer
	Def(
		c,
		"format",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			arg := args[1]
			switch a := arg.SafeAsReference().(type) {
			case value.String:
				result, err := self.Format(string(a))
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.Ref(value.String(result)), value.Undefined
			default:
				return value.Undefined, value.Ref(value.Errorf(
					value.ArgumentErrorClass,
					"expected a format string, got: %s",
					arg.Inspect(),
				))
			}
		},
		DefWithParameters(1),
	)
	Alias(c, "strftime", "format")

	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			dur := args[1].MustDuration()
			return value.Ref(self.Subtract(dur)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			dur := args[1].MustDuration()
			return value.Ref(self.Add(dur)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"diff",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			other := args[1].MustReference().(*value.Time)
			return self.Diff(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"zone",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(self.Zone()), value.Undefined
		},
	)
	Alias(c, "timezone", "zone")

	Def(
		c,
		"zone_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.String(self.ZoneName())), value.Undefined
		},
	)
	Alias(c, "timezone_name", "zone_name")

	Def(
		c,
		"zone_offset_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.ZoneOffsetSeconds()).ToValue(), value.Undefined
		},
	)
	Alias(c, "timezone_offset_seconds", "zone_offset_seconds")

	Def(
		c,
		"zone_offset_hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.ZoneOffsetHours()).ToValue(), value.Undefined
		},
	)
	Alias(c, "timezone_offset_hours", "zone_offset_hours")

	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(self.ToString()), value.Undefined
		},
	)

	Def(
		c,
		"year",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Year()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"iso_year",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.ISOYear()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"month",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Month()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"week_from_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.WeekFromMonday()).ToValue(), value.Undefined
		},
	)
	Alias(c, "week", "week_from_monday")
	Def(
		c,
		"week_from_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.WeekFromSunday()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"iso_week",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.ISOWeek()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Day()).ToValue(), value.Undefined
		},
	)
	Alias(c, "month_day", "day")

	Def(
		c,
		"year_day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.YearDay()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"weekday_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.String(self.WeekdayName())), value.Undefined
		},
	)
	Def(
		c,
		"weekday_from_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.WeekdayFromMonday()).ToValue(), value.Undefined
		},
	)
	Alias(c, "weekday", "weekday_from_monday")

	Def(
		c,
		"weekday_from_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.WeekdayFromSunday()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"hour",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Hour()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"minute",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Minute()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Second()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"millisecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Millisecond()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"microsecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Microsecond()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"nanosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Nanosecond()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"picosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkInt(self.Picosecond()), value.Undefined
		},
	)
	Def(
		c,
		"femtosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkInt(self.Femtosecond()), value.Undefined
		},
	)
	Def(
		c,
		"attosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkInt(self.Attosecond()), value.Undefined
		},
	)
	Def(
		c,
		"zeptosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.Zeptosecond())), value.Undefined
		},
	)
	Def(
		c,
		"yoctosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.Yoctosecond())), value.Undefined
		},
	)
	Def(
		c,
		"unix_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkInt(self.UnixSeconds()), value.Undefined
		},
	)
	Alias(c, "unix", "unix_seconds")
	Def(
		c,
		"unix_milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkInt(self.UnixMilliseconds()), value.Undefined
		},
	)
	Def(
		c,
		"unix_microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixMicroseconds())), value.Undefined
		},
	)
	Def(
		c,
		"unix_nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixNanoseconds())), value.Undefined
		},
	)
	Def(
		c,
		"unix_picoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixPicoseconds())), value.Undefined
		},
	)
	Def(
		c,
		"unix_femtoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixFemtoseconds())), value.Undefined
		},
	)
	Def(
		c,
		"unix_attoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixAttoseconds())), value.Undefined
		},
	)
	Def(
		c,
		"unix_zeptoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixZeptoseconds())), value.Undefined
		},
	)
	Def(
		c,
		"unix_yoctoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixYoctoseconds())), value.Undefined
		},
	)

	Def(
		c,
		"to_utc",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(self.UTC()), value.Undefined
		},
	)
	Alias(c, "utc", "to_utc")

	Def(
		c,
		"to_local",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(self.Local()), value.Undefined
		},
	)
	Alias(c, "local", "to_local")
	Def(
		c,
		"is_utc",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsUTC()), value.Undefined
		},
	)
	Def(
		c,
		"is_local",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsLocal()), value.Undefined
		},
	)

	Def(
		c,
		"is_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsMonday()), value.Undefined
		},
	)
	Def(
		c,
		"is_tuesday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsTuesday()), value.Undefined
		},
	)
	Def(
		c,
		"is_wednesday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsWednesday()), value.Undefined
		},
	)
	Def(
		c,
		"is_thursday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsThursday()), value.Undefined
		},
	)
	Def(
		c,
		"is_friday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsFriday()), value.Undefined
		},
	)
	Def(
		c,
		"is_saturday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsSaturday()), value.Undefined
		},
	)
	Def(
		c,
		"is_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsSunday()), value.Undefined
		},
	)

	Def(
		c,
		"is_am",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsAM()), value.Undefined
		},
	)
	Def(
		c,
		"is_pm",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsPM()), value.Undefined
		},
	)
	Def(
		c,
		"meridiem",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.String(self.Meridiem())), value.Undefined
		},
	)

	Def(
		c,
		"hour12",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Hour12()).ToValue(), value.Undefined
		},
	)

	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			other := args[1]
			return self.GreaterThan(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			other := args[1]
			return self.GreaterThanEqual(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			other := args[1]
			return self.LessThan(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			other := args[1]
			return self.LessThanEqual(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			other := args[1]
			return self.Equal(other), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "===", "==")
}
