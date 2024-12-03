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
			return value.TimeNow(), nil
		},
	)

	// Instance methods
	c = &value.TimeClass.MethodContainer
	Def(
		c,
		"format",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			arg := args[1]
			switch a := arg.(type) {
			case value.String:
				result, err := self.Format(string(a))
				if err != nil {
					return nil, err
				}
				return value.String(result), nil
			default:
				return nil, value.Errorf(
					value.ArgumentErrorClass,
					"expected a format string, got: %s",
					arg.Inspect(),
				)
			}
		},
		DefWithParameters(1),
	)
	Alias(c, "strftime", "format")

	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			dur := args[1].(value.Duration)
			return self.Subtract(dur), nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			dur := args[1].(value.Duration)
			return self.Add(dur), nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"diff",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			other := args[1].(*value.Time)
			return self.Diff(other), nil
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"zone",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return self.Zone(), nil
		},
	)
	Alias(c, "timezone", "zone")

	Def(
		c,
		"zone_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.String(self.ZoneName()), nil
		},
	)
	Alias(c, "timezone_name", "zone_name")

	Def(
		c,
		"zone_offset_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.ZoneOffsetSeconds()), nil
		},
	)
	Alias(c, "timezone_offset_seconds", "zone_offset_seconds")

	Def(
		c,
		"zone_offset_hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.ZoneOffsetHours()), nil
		},
	)
	Alias(c, "timezone_offset_hours", "zone_offset_hours")

	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return self.ToString(), nil
		},
	)

	Def(
		c,
		"year",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.Year()), nil
		},
	)
	Def(
		c,
		"iso_year",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.ISOYear()), nil
		},
	)
	Def(
		c,
		"month",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.Month()), nil
		},
	)
	Def(
		c,
		"week_from_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.WeekFromMonday()), nil
		},
	)
	Alias(c, "week", "week_from_monday")
	Def(
		c,
		"week_from_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.WeekFromSunday()), nil
		},
	)

	Def(
		c,
		"iso_week",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.ISOWeek()), nil
		},
	)
	Def(
		c,
		"day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.Day()), nil
		},
	)
	Alias(c, "month_day", "day")

	Def(
		c,
		"year_day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.YearDay()), nil
		},
	)
	Def(
		c,
		"weekday_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.String(self.WeekdayName()), nil
		},
	)
	Def(
		c,
		"weekday_from_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.WeekdayFromMonday()), nil
		},
	)
	Alias(c, "weekday", "weekday_from_monday")

	Def(
		c,
		"weekday_from_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.WeekdayFromSunday()), nil
		},
	)
	Def(
		c,
		"hour",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.Hour()), nil
		},
	)
	Def(
		c,
		"minute",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.Minute()), nil
		},
	)
	Def(
		c,
		"second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.Second()), nil
		},
	)
	Def(
		c,
		"millisecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.Millisecond()), nil
		},
	)
	Def(
		c,
		"microsecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.Microsecond()), nil
		},
	)
	Def(
		c,
		"nanosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.Nanosecond()), nil
		},
	)
	Def(
		c,
		"picosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkInt(self.Picosecond()), nil
		},
	)
	Def(
		c,
		"femtosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkInt(self.Femtosecond()), nil
		},
	)
	Def(
		c,
		"attosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkInt(self.Attosecond()), nil
		},
	)
	Def(
		c,
		"zeptosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBigInt(self.Zeptosecond()), nil
		},
	)
	Def(
		c,
		"yoctosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBigInt(self.Yoctosecond()), nil
		},
	)

	Def(
		c,
		"unix_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkInt(self.UnixSeconds()), nil
		},
	)
	Alias(c, "unix", "unix_seconds")
	Def(
		c,
		"unix_milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkInt(self.UnixMilliseconds()), nil
		},
	)
	Def(
		c,
		"unix_microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBigInt(self.UnixMicroseconds()), nil
		},
	)
	Def(
		c,
		"unix_nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBigInt(self.UnixNanoseconds()), nil
		},
	)
	Def(
		c,
		"unix_picoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBigInt(self.UnixPicoseconds()), nil
		},
	)
	Def(
		c,
		"unix_femtoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBigInt(self.UnixFemtoseconds()), nil
		},
	)
	Def(
		c,
		"unix_attoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBigInt(self.UnixAttoseconds()), nil
		},
	)
	Def(
		c,
		"unix_zeptoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBigInt(self.UnixZeptoseconds()), nil
		},
	)
	Def(
		c,
		"unix_yoctoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBigInt(self.UnixYoctoseconds()), nil
		},
	)

	Def(
		c,
		"to_utc",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return self.UTC(), nil
		},
	)
	Alias(c, "utc", "to_utc")

	Def(
		c,
		"to_local",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return self.Local(), nil
		},
	)
	Alias(c, "local", "to_local")
	Def(
		c,
		"is_utc",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBool(self.IsUTC()), nil
		},
	)
	Def(
		c,
		"is_local",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBool(self.IsLocal()), nil
		},
	)

	Def(
		c,
		"is_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBool(self.IsMonday()), nil
		},
	)
	Def(
		c,
		"is_tuesday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBool(self.IsTuesday()), nil
		},
	)
	Def(
		c,
		"is_wednesday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBool(self.IsWednesday()), nil
		},
	)
	Def(
		c,
		"is_thursday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBool(self.IsThursday()), nil
		},
	)
	Def(
		c,
		"is_friday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBool(self.IsFriday()), nil
		},
	)
	Def(
		c,
		"is_saturday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBool(self.IsSaturday()), nil
		},
	)
	Def(
		c,
		"is_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBool(self.IsSunday()), nil
		},
	)

	Def(
		c,
		"is_am",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBool(self.IsAM()), nil
		},
	)
	Def(
		c,
		"is_pm",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.ToElkBool(self.IsPM()), nil
		},
	)
	Def(
		c,
		"meridiem",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.String(self.Meridiem()), nil
		},
	)

	Def(
		c,
		"hour12",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			return value.SmallInt(self.Hour12()), nil
		},
	)

	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			other := args[1]
			return self.GreaterThan(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			other := args[1]
			return self.GreaterThanEqual(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			other := args[1]
			return self.LessThan(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			other := args[1]
			return self.LessThanEqual(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Time)
			other := args[1]
			return self.Equal(other), nil
		},
		DefWithParameters(1),
	)
	Alias(c, "===", "==")
}
