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
			return value.Ref(value.TimeNow()), value.Nil
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
				if !err.IsNil() {
					return value.Nil, err
				}
				return value.Ref(value.String(result)), value.Nil
			default:
				return value.Nil, value.Ref(value.Errorf(
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
			return value.Ref(self.Subtract(dur)), value.Nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			dur := args[1].MustDuration()
			return value.Ref(self.Add(dur)), value.Nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"diff",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			other := args[1].MustReference().(*value.Time)
			return self.Diff(other).ToValue(), value.Nil
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"zone",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(self.Zone()), value.Nil
		},
	)
	Alias(c, "timezone", "zone")

	Def(
		c,
		"zone_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.String(self.ZoneName())), value.Nil
		},
	)
	Alias(c, "timezone_name", "zone_name")

	Def(
		c,
		"zone_offset_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.ZoneOffsetSeconds()).ToValue(), value.Nil
		},
	)
	Alias(c, "timezone_offset_seconds", "zone_offset_seconds")

	Def(
		c,
		"zone_offset_hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.ZoneOffsetHours()).ToValue(), value.Nil
		},
	)
	Alias(c, "timezone_offset_hours", "zone_offset_hours")

	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(self.ToString()), value.Nil
		},
	)

	Def(
		c,
		"year",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Year()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"iso_year",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.ISOYear()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"month",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Month()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"week_from_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.WeekFromMonday()).ToValue(), value.Nil
		},
	)
	Alias(c, "week", "week_from_monday")
	Def(
		c,
		"week_from_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.WeekFromSunday()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"iso_week",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.ISOWeek()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Day()).ToValue(), value.Nil
		},
	)
	Alias(c, "month_day", "day")

	Def(
		c,
		"year_day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.YearDay()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"weekday_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.String(self.WeekdayName())), value.Nil
		},
	)
	Def(
		c,
		"weekday_from_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.WeekdayFromMonday()).ToValue(), value.Nil
		},
	)
	Alias(c, "weekday", "weekday_from_monday")

	Def(
		c,
		"weekday_from_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.WeekdayFromSunday()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"hour",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Hour()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"minute",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Minute()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Second()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"millisecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Millisecond()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"microsecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Microsecond()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"nanosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Nanosecond()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"picosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkInt(self.Picosecond()), value.Nil
		},
	)
	Def(
		c,
		"femtosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkInt(self.Femtosecond()), value.Nil
		},
	)
	Def(
		c,
		"attosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkInt(self.Attosecond()), value.Nil
		},
	)
	Def(
		c,
		"zeptosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.Zeptosecond())), value.Nil
		},
	)
	Def(
		c,
		"yoctosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.Yoctosecond())), value.Nil
		},
	)
	Def(
		c,
		"unix_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkInt(self.UnixSeconds()), value.Nil
		},
	)
	Alias(c, "unix", "unix_seconds")
	Def(
		c,
		"unix_milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkInt(self.UnixMilliseconds()), value.Nil
		},
	)
	Def(
		c,
		"unix_microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixMicroseconds())), value.Nil
		},
	)
	Def(
		c,
		"unix_nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixNanoseconds())), value.Nil
		},
	)
	Def(
		c,
		"unix_picoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixPicoseconds())), value.Nil
		},
	)
	Def(
		c,
		"unix_femtoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixFemtoseconds())), value.Nil
		},
	)
	Def(
		c,
		"unix_attoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixAttoseconds())), value.Nil
		},
	)
	Def(
		c,
		"unix_zeptoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixZeptoseconds())), value.Nil
		},
	)
	Def(
		c,
		"unix_yoctoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.ToElkBigInt(self.UnixYoctoseconds())), value.Nil
		},
	)

	Def(
		c,
		"to_utc",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(self.UTC()), value.Nil
		},
	)
	Alias(c, "utc", "to_utc")

	Def(
		c,
		"to_local",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(self.Local()), value.Nil
		},
	)
	Alias(c, "local", "to_local")
	Def(
		c,
		"is_utc",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsUTC()), value.Nil
		},
	)
	Def(
		c,
		"is_local",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsLocal()), value.Nil
		},
	)

	Def(
		c,
		"is_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsMonday()), value.Nil
		},
	)
	Def(
		c,
		"is_tuesday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsTuesday()), value.Nil
		},
	)
	Def(
		c,
		"is_wednesday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsWednesday()), value.Nil
		},
	)
	Def(
		c,
		"is_thursday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsThursday()), value.Nil
		},
	)
	Def(
		c,
		"is_friday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsFriday()), value.Nil
		},
	)
	Def(
		c,
		"is_saturday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsSaturday()), value.Nil
		},
	)
	Def(
		c,
		"is_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsSunday()), value.Nil
		},
	)

	Def(
		c,
		"is_am",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsAM()), value.Nil
		},
	)
	Def(
		c,
		"is_pm",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.ToElkBool(self.IsPM()), value.Nil
		},
	)
	Def(
		c,
		"meridiem",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.Ref(value.String(self.Meridiem())), value.Nil
		},
	)

	Def(
		c,
		"hour12",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Time)
			return value.SmallInt(self.Hour12()).ToValue(), value.Nil
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
			return self.Equal(other), value.Nil
		},
		DefWithParameters(1),
	)
	Alias(c, "===", "==")
}
