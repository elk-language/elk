package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Date
func initDate() {
	// Singleton methods
	c := &value.DateClass.SingletonClass().MethodContainer
	Def(
		c,
		"now",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.DateNow().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"parse",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			formatString := args[1].AsString().String()
			input := args[2].AsString().String()
			result, err := value.ParseDate(formatString, input)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return result.ToValue(), value.Undefined
		},
		DefWithParameters(2),
	)

	// Instance methods
	c = &value.DateClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			var argYear int
			if !args[1].IsUndefined() {
				if args[1].IsSmallInt() {
					argYear = int(args[1].AsSmallInt())
				} else {
					argYear = value.MaxSmallInt
				}
			}

			var argMonth int
			if !args[2].IsUndefined() {
				if args[2].IsSmallInt() {
					argMonth = int(args[2].AsSmallInt())
				} else {
					argMonth = value.MaxSmallInt
				}
			} else {
				argMonth = 1
			}

			var argDay int
			if !args[3].IsUndefined() {
				if args[3].IsSmallInt() {
					argDay = int(args[3].AsSmallInt())
				} else {
					argDay = value.MaxSmallInt
				}
			} else {
				argDay = 1
			}

			self, err := value.MakeValidatedDate(
				argYear,
				argMonth,
				argDay,
			)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(3),
	)
	Def(
		c,
		"format",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
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
		"year",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.SmallInt(self.Year()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"iso_year",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.SmallInt(self.ISOYear()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"month",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDate()
			return value.SmallInt(self.Month()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDate()
			return value.Ref(self.ToString()), value.Undefined
		},
	)
	Def(
		c,
		"week_from_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.SmallInt(self.WeekFromMonday()).ToValue(), value.Undefined
		},
	)
	Alias(c, "week", "week_from_monday")
	Def(
		c,
		"week_from_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.SmallInt(self.WeekFromSunday()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"iso_week",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.SmallInt(self.ISOWeek()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.SmallInt(self.Day()).ToValue(), value.Undefined
		},
	)
	Alias(c, "month_day", "day")

	Def(
		c,
		"year_day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.SmallInt(self.YearDay()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"weekday_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.Ref(value.String(self.WeekdayName())), value.Undefined
		},
	)
	Def(
		c,
		"weekday_from_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.SmallInt(self.WeekdayFromMonday()).ToValue(), value.Undefined
		},
	)
	Alias(c, "weekday", "weekday_from_monday")

	Def(
		c,
		"weekday_from_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.SmallInt(self.WeekdayFromSunday()).ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"is_monday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.ToElkBool(self.IsMonday()), value.Undefined
		},
	)
	Def(
		c,
		"is_tuesday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.ToElkBool(self.IsTuesday()), value.Undefined
		},
	)
	Def(
		c,
		"is_wednesday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.ToElkBool(self.IsWednesday()), value.Undefined
		},
	)
	Def(
		c,
		"is_thursday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.ToElkBool(self.IsThursday()), value.Undefined
		},
	)
	Def(
		c,
		"is_friday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.ToElkBool(self.IsFriday()), value.Undefined
		},
	)
	Def(
		c,
		"is_saturday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.ToElkBool(self.IsSaturday()), value.Undefined
		},
	)
	Def(
		c,
		"is_sunday",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.ToElkBool(self.IsSunday()), value.Undefined
		},
	)
}
