package vm

import (
	"github.com/elk-language/elk/value"
)

func init() {
	DefineMethodWithOptions(
		value.TimeClass.SingletonClass().Methods,
		"now",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.TimeNow(), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"format",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
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
		NativeMethodWithStringParameters("format_string"),
	)
	value.TimeClass.DefineAliasString("strftime", "format")

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"zone",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return self.Zone(), nil
		},
	)
	value.TimeClass.DefineAliasString("timezone", "zone")

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"zone_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.String(self.ZoneName()), nil
		},
	)
	value.TimeClass.DefineAliasString("timezone_name", "zone_name")

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"zone_offset_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.ZoneOffsetSeconds()), nil
		},
	)
	value.TimeClass.DefineAliasString("timezone_offset_seconds", "zone_offset_seconds")

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"zone_offset_hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.ZoneOffsetHours()), nil
		},
	)
	value.TimeClass.DefineAliasString("timezone_offset_hours", "zone_offset_hours")

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return self.ToString(), nil
		},
	)

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"year",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.Year()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"iso_year",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.ISOYear()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"month",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.Month()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"week_from_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.WeekFromMonday()), nil
		},
	)
	value.TimeClass.DefineAliasString("week", "week_from_monday")
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"week_from_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.WeekFromSunday()), nil
		},
	)

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"iso_week",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.ISOWeek()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.Day()), nil
		},
	)
	value.TimeClass.DefineAliasString("month_day", "day")

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"year_day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.YearDay()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"weekday_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.String(self.WeekdayName()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"weekday_from_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.WeekdayFromMonday()), nil
		},
	)
	value.TimeClass.DefineAliasString("weekday", "weekday_from_monday")

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"weekday_from_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.WeekdayFromSunday()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"hour",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.Hour()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"minute",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.Minute()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.Second()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"millisecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.Millisecond()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"microsecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.Microsecond()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"nanosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.Nanosecond()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"picosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkInt(self.Picosecond()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"femtosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkInt(self.Femtosecond()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"attosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkInt(self.Attosecond()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"zeptosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBigInt(self.Zeptosecond()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"yoctosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBigInt(self.Yoctosecond()), nil
		},
	)

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"unix_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkInt(self.UnixSeconds()), nil
		},
	)
	value.TimeClass.DefineAliasString("unix", "unix_seconds")
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"unix_milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkInt(self.UnixMilliseconds()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"unix_microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBigInt(self.UnixMicroseconds()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"unix_nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBigInt(self.UnixNanoseconds()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"unix_picoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBigInt(self.UnixPicoseconds()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"unix_femtoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBigInt(self.UnixFemtoseconds()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"unix_attoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBigInt(self.UnixAttoseconds()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"unix_zeptoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBigInt(self.UnixZeptoseconds()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"unix_yoctoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBigInt(self.UnixYoctoseconds()), nil
		},
	)

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"to_utc",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return self.UTC(), nil
		},
	)
	value.TimeClass.DefineAliasString("utc", "to_utc")

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"to_local",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return self.Local(), nil
		},
	)
	value.TimeClass.DefineAliasString("local", "to_local")
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"is_utc",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBool(self.IsUTC()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"is_local",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBool(self.IsLocal()), nil
		},
	)

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"is_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBool(self.IsMonday()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"is_tuesday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBool(self.IsTuesday()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"is_wednesday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBool(self.IsWednesday()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"is_thursday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBool(self.IsThursday()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"is_friday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBool(self.IsFriday()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"is_saturday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBool(self.IsSaturday()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"is_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBool(self.IsSunday()), nil
		},
	)

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"is_am",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBool(self.IsAM()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"is_pm",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.ToElkBool(self.IsPM()), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"meridiem",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.String(self.Meridiem()), nil
		},
	)

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"hour12",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return value.SmallInt(self.Hour12()), nil
		},
	)

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			other := args[1]
			return value.ToValueErr(self.GreaterThan(other))
		},
		NativeMethodWithStringParameters("other"),
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			other := args[1]
			return value.ToValueErr(self.GreaterThanEqual(other))
		},
		NativeMethodWithStringParameters("other"),
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			other := args[1]
			return value.ToValueErr(self.LessThan(other))
		},
		NativeMethodWithStringParameters("other"),
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			other := args[1]
			return value.ToValueErr(self.LessThanEqual(other))
		},
		NativeMethodWithStringParameters("other"),
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			other := args[1]
			return self.Equal(other), nil
		},
		NativeMethodWithStringParameters("other"),
	)
	value.TimeClass.DefineAliasString("===", "==")
}
