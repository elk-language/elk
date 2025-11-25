package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Date::Span
func initDateSpan() {
	// Singleton methods
	c := &value.DateSpanClass.SingletonClass().MethodContainer
	Def(
		c,
		"parse",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			str := string(args[1].AsString())
			return value.RefErr(value.ParseDateSpan(str))
		},
		DefWithParameters(1),
	)

	// Instance methods
	c = &value.DateSpanClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			var argYear int
			if !args[1].IsUndefined() {
				if args[1].IsSmallInt() {
					argYear = int(args[1].AsSmallInt())
				} else {
					argYear = int(args[1].AsBigInt().ToSmallInt())
				}
			}

			var argMonth int
			if !args[2].IsUndefined() {
				if args[2].IsSmallInt() {
					argMonth = int(args[2].AsSmallInt())
				} else {
					argMonth = int(args[2].AsBigInt().ToSmallInt())
				}
			}

			var argDay int
			if !args[3].IsUndefined() {
				if args[3].IsSmallInt() {
					argDay = int(args[3].AsSmallInt())
				} else {
					argDay = int(args[3].AsBigInt().ToSmallInt())
				}
			}

			self := value.MakeDateSpan(
				argYear,
				argMonth,
				argDay,
			)
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(3),
	)

	Def(
		c,
		"-@",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.Negate().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"+@",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)

	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := args[1]
			return self.Add(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+@1",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := args[1].AsDateSpan()
			return self.AddDateSpan(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+@2",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := args[1].AsTimeSpan()
			return value.Ref(self.AddTimeSpan(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+@3",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := (*value.DateTimeSpan)(args[1].Pointer())
			return value.Ref(self.AddDateTimeSpan(other)), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := args[1]
			return self.Subtract(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@1",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := args[1].AsDateSpan()
			return self.SubtractDateSpan(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@2",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := args[1].AsTimeSpan()
			return value.Ref(self.SubtractTimeSpan(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@3",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := (*value.DateTimeSpan)(args[1].Pointer())
			return value.Ref(self.SubtractDateTimeSpan(other)), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := args[1]
			return self.Multiply(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*@1",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := args[1]
			return self.MultiplyInt(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*@2",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := args[1].AsFloat()
			return value.Ref(self.MultiplyFloat(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*@3",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := (*value.BigFloat)(args[1].Pointer())
			return value.Ref(self.MultiplyBigFloat(other)), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"/",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := args[1]
			return self.Divide(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/@1",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := args[1]
			return value.Ref(self.DivideInt(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/@2",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := args[1].AsFloat()
			return value.Ref(self.DivideFloat(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/@3",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := (*value.BigFloat)(args[1].Pointer())
			return value.Ref(self.DivideBigFloat(other)), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return value.SmallInt(self.Years()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.TotalYears(), value.Undefined
		},
	)
	Def(
		c,
		"in_years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.InYears().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"months",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return value.SmallInt(self.Months()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_months",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.TotalMonths(), value.Undefined
		},
	)
	Def(
		c,
		"in_months",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.InMonths().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_weeks",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.TotalWeeks(), value.Undefined
		},
	)
	Def(
		c,
		"in_weeks",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.InWeeks().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return value.SmallInt(self.Days()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.TotalDays(), value.Undefined
		},
	)
	Def(
		c,
		"in_days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.InDays().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.SmallInt(0).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.TotalHours(), value.Undefined
		},
	)
	Def(
		c,
		"in_hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.InHours().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.SmallInt(0).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.TotalMinutes(), value.Undefined
		},
	)
	Def(
		c,
		"in_minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.InMinutes().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.SmallInt(0).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.TotalSeconds(), value.Undefined
		},
	)
	Def(
		c,
		"in_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.InSeconds().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.SmallInt(0).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.TotalMilliseconds(), value.Undefined
		},
	)
	Def(
		c,
		"in_milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.InMilliseconds().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.SmallInt(0).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.TotalMicroseconds(), value.Undefined
		},
	)
	Def(
		c,
		"in_microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.InMicroseconds().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.SmallInt(0).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.TotalNanoseconds(), value.Undefined
		},
	)
	Def(
		c,
		"in_nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.InNanoseconds().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return value.Ref(self.ToString()), value.Undefined
		},
	)
	Def(
		c,
		"to_datetime_span",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return value.Ref(self.ToDateTimeSpan()), value.Undefined
		},
	)
	Def(
		c,
		"to_date",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.ToDate().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_date_span",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return self.CompareVal(args[1])
		},
		DefWithParameters(1),
	)

	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			ok, err := self.GreaterThanEqual(args[1])
			return value.ToElkBool(ok), err
		},
		DefWithParameters(1),
	)

	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			ok, err := self.GreaterThan(args[1])
			return value.ToElkBool(ok), err
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			ok, err := self.LessThanEqual(args[1])
			return value.ToElkBool(ok), err
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			ok, err := self.LessThan(args[1])
			return value.ToElkBool(ok), err
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "===", "==")
}
