package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::DateTime
func initDateTime() {
	// Class methods
	c := &value.DateTimeClass.SingletonClass().MethodContainer
	Def(
		c,
		"now",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.Ref(value.DateTimeNow()), value.Undefined
		},
	)
	Def(
		c,
		"parse",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			formatString := args[1].AsString().String()
			input := args[2].AsString().String()
			result, err := value.ParseDateTime(formatString, input)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.Ref(result), value.Undefined
		},
		DefWithParameters(2),
	)

	// Instance methods
	c = &value.DateTimeClass.MethodContainer
	Def(
		c,
		"#init",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			var year int
			if !args[1].IsUndefined() {
				year = args[1].AsInt()
			}

			var month int
			if args[2].IsUndefined() {
				month = 1
			} else {
				month = args[2].AsInt()
			}

			var day int
			if args[3].IsUndefined() {
				day = 1
			} else {
				day = args[3].AsInt()
			}

			var hour int
			if !args[4].IsUndefined() {
				hour = args[4].AsInt()
			}

			var minute int
			if !args[5].IsUndefined() {
				minute = args[5].AsInt()
			}

			var second int
			if !args[6].IsUndefined() {
				second = args[6].AsInt()
			}

			var millisecond int
			if !args[7].IsUndefined() {
				millisecond = args[7].AsInt()
			}

			var microsecond int
			if !args[8].IsUndefined() {
				microsecond = args[8].AsInt()
			}

			var nanosecond int
			if !args[9].IsUndefined() {
				nanosecond = args[9].AsInt()
			}

			var zone *value.Timezone
			if args[10].IsUndefined() {
				zone = value.LocalTimezone
			} else {
				zone = (*value.Timezone)(args[10].Pointer())
			}

			return value.Ref(value.NewDateTime(year, month, day, hour, minute, second, millisecond, microsecond, nanosecond, zone)), value.Undefined
		},
		DefWithParameters(10),
	)

	Def(
		c,
		"format",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
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
			self := args[0].MustReference().(*value.DateTime)
			return self.Subtract(args[1])
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@1",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			span := (*value.DateTimeSpan)(args[1].Pointer())
			return value.Ref(self.SubtractDateTimeSpan(span)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@2",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(self.SubtractTimeSpan(args[1].AsTimeSpan())), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@3",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(self.SubtractDateSpan(args[1].AsDateSpan())), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@4",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			other := (*value.DateTime)(args[1].Pointer())
			return value.Ref(self.DiffDateTime(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@5",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			other := args[1].AsDate()
			return value.Ref(self.DiffDate(other)), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return self.Add(args[1])
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+@1",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			span := (*value.DateTimeSpan)(args[1].Pointer())
			return value.Ref(self.AddDateTimeSpan(span)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+@2",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(self.AddTimeSpan(args[1].AsTimeSpan())), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+@3",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(self.AddDateSpan(args[1].AsDateSpan())), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"diff",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return self.Diff(args[1])
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"diff@1",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			other := args[1].MustReference().(*value.DateTime)
			return value.Ref(self.DiffDateTime(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"diff@2",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			other := args[1].AsDate()
			return value.Ref(self.DiffDate(other)), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"date",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(*value.DateTime)
			return self.Date().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"time",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(*value.DateTime)
			return self.Time().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"zone",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(self.Zone()), value.Undefined
		},
	)
	Alias(c, "timezone", "zone")

	Def(
		c,
		"zone_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(value.String(self.ZoneName())), value.Undefined
		},
	)
	Alias(c, "timezone_name", "zone_name")

	Def(
		c,
		"zone_offset",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.ZoneOffset()).ToValue(), value.Undefined
		},
	)
	Alias(c, "timezone_offset", "zone_offset")

	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(self.ToString()), value.Undefined
		},
	)

	Def(
		c,
		"year",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.Year()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"iso_year",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.ISOYear()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"month",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.Month()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"week_from_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.WeekFromMonday()).ToValue(), value.Undefined
		},
	)
	Alias(c, "week", "week_from_monday")
	Def(
		c,
		"week_from_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.WeekFromSunday()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"iso_week",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.ISOWeek()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.Day()).ToValue(), value.Undefined
		},
	)
	Alias(c, "month_day", "day")

	Def(
		c,
		"year_day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.YearDay()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"iso_year_day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.ISOYearDay()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"weekday_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(value.String(self.WeekdayName())), value.Undefined
		},
	)
	Def(
		c,
		"abbreviated_weekday_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(value.String(self.AbbreviatedWeekdayName())), value.Undefined
		},
	)
	Def(
		c,
		"weekday_from_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.WeekdayFromMonday()).ToValue(), value.Undefined
		},
	)
	Alias(c, "weekday", "weekday_from_monday")

	Def(
		c,
		"weekday_from_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.WeekdayFromSunday()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"hour",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.Hour()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"minute",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.Minute()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.Second()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"millisecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.Millisecond()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"microsecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.Microsecond()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"microseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.MicrosecondsInSecond()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"nanosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.Nanosecond()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"nanoseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.NanosecondsInSecond()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"picoseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkInt(self.PicosecondsInSecond()), value.Undefined
		},
	)
	Def(
		c,
		"femtoseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkInt(self.FemtosecondsInSecond()), value.Undefined
		},
	)
	Def(
		c,
		"attoseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkInt(self.AttosecondsInSecond()), value.Undefined
		},
	)
	Def(
		c,
		"zeptoseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(value.ToElkBigInt(self.ZeptosecondsInSecond())), value.Undefined
		},
	)
	Def(
		c,
		"yoctoseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(value.ToElkBigInt(self.YoctosecondsInSecond())), value.Undefined
		},
	)
	Def(
		c,
		"unix_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkInt(self.UnixSeconds()), value.Undefined
		},
	)
	Alias(c, "unix", "unix_seconds")
	Def(
		c,
		"unix_milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkInt(self.UnixMilliseconds()), value.Undefined
		},
	)
	Def(
		c,
		"unix_microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(value.ToElkBigInt(self.UnixMicroseconds())), value.Undefined
		},
	)
	Def(
		c,
		"unix_nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(value.ToElkBigInt(self.UnixNanoseconds())), value.Undefined
		},
	)
	Def(
		c,
		"unix_picoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(value.ToElkBigInt(self.UnixPicoseconds())), value.Undefined
		},
	)
	Def(
		c,
		"unix_femtoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(value.ToElkBigInt(self.UnixFemtoseconds())), value.Undefined
		},
	)
	Def(
		c,
		"unix_attoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(value.ToElkBigInt(self.UnixAttoseconds())), value.Undefined
		},
	)
	Def(
		c,
		"unix_zeptoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(value.ToElkBigInt(self.UnixZeptoseconds())), value.Undefined
		},
	)
	Def(
		c,
		"unix_yoctoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(value.ToElkBigInt(self.UnixYoctoseconds())), value.Undefined
		},
	)

	Def(
		c,
		"to_utc",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(self.UTC()), value.Undefined
		},
	)
	Alias(c, "utc", "to_utc")

	Def(
		c,
		"to_local",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(self.Local()), value.Undefined
		},
	)
	Alias(c, "local", "to_local")
	Def(
		c,
		"is_utc",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkBool(self.IsUTC()), value.Undefined
		},
	)
	Def(
		c,
		"is_local",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkBool(self.IsLocal()), value.Undefined
		},
	)

	Def(
		c,
		"is_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkBool(self.IsMonday()), value.Undefined
		},
	)
	Def(
		c,
		"is_tuesday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkBool(self.IsTuesday()), value.Undefined
		},
	)
	Def(
		c,
		"is_wednesday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkBool(self.IsWednesday()), value.Undefined
		},
	)
	Def(
		c,
		"is_thursday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkBool(self.IsThursday()), value.Undefined
		},
	)
	Def(
		c,
		"is_friday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkBool(self.IsFriday()), value.Undefined
		},
	)
	Def(
		c,
		"is_saturday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkBool(self.IsSaturday()), value.Undefined
		},
	)
	Def(
		c,
		"is_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkBool(self.IsSunday()), value.Undefined
		},
	)

	Def(
		c,
		"is_am",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkBool(self.IsAM()), value.Undefined
		},
	)
	Def(
		c,
		"is_pm",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.ToElkBool(self.IsPM()), value.Undefined
		},
	)
	Def(
		c,
		"meridiem",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.Ref(value.String(self.Meridiem())), value.Undefined
		},
	)

	Def(
		c,
		"hour12",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return value.SmallInt(self.Hour12()).ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"in_zone",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTime)(args[0].Pointer())
			zone := (*value.Timezone)(args[1].Pointer())
			return value.Ref(self.InZone(zone)), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"with_zone",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTime)(args[0].Pointer())
			zone := (*value.Timezone)(args[1].Pointer())
			return value.Ref(self.WithZone(zone)), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			return self.CompareVal(args[1])
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			other := args[1]
			return self.GreaterThanVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			other := args[1]
			return self.GreaterThanEqualVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			other := args[1]
			return self.LessThanVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			other := args[1]
			return self.LessThanEqualVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DateTime)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "===", "==")
}
