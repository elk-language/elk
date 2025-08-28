package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::TimeSpan
func initTimeSpan() {
	// Class methods
	c := &value.TimeSpanClass.SingletonClass().MethodContainer
	Def(
		c,
		"parse",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			str := args[1].AsReference().(value.String)
			return value.ToValueErr(value.ParsTimeSpan(str))
		},
		DefWithParameters(1),
	)

	// Instance methods
	c = &value.TimeSpanClass.MethodContainer
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			dur := args[1].MustTimeSpan()
			return self.Add(dur).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			dur := args[1].MustTimeSpan()
			return self.Subtract(dur).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			other := args[1]
			return value.ToValueErr(self.Multiply(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			other := args[1]
			return value.ToValueErr(self.Divide(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return value.Ref(self.ToString()), value.Undefined
		},
	)
	Def(
		c,
		"total_nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.TotalNanoseconds(), value.Undefined
		},
	)
	Def(
		c,
		"nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return value.SmallInt(self.Nanoseconds()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.InNanoseconds().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.TotalMicroseconds(), value.Undefined
		},
	)
	Def(
		c,
		"microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return value.SmallInt(self.Microseconds()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.InMicroseconds().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.TotalMilliseconds(), value.Undefined
		},
	)
	Def(
		c,
		"milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return value.SmallInt(self.Milliseconds()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.InMilliseconds().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.TotalSeconds(), value.Undefined
		},
	)
	Def(
		c,
		"seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return value.SmallInt(self.Seconds()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.InSeconds().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.TotalMinutes(), value.Undefined
		},
	)
	Def(
		c,
		"minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return value.SmallInt(self.Minutes()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.InMinutes().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.TotalHours(), value.Undefined
		},
	)
	Def(
		c,
		"hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return value.SmallInt(self.Hours()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.InHours().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.TotalDays(), value.Undefined
		},
	)
	Def(
		c,
		"in_days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.InDays().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_weeks",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.TotalWeeks(), value.Undefined
		},
	)
	Def(
		c,
		"in_weeks",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.InWeeks().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_months",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.TotalMonths(), value.Undefined
		},
	)
	Def(
		c,
		"months",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return value.SmallInt(self.Months()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_months",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.InMonths().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.TotalYears(), value.Undefined
		},
	)
	Alias(c, "years", "total years")
	Def(
		c,
		"in_years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustTimeSpan()
			return self.InYears().ToValue(), value.Undefined
		},
	)
}
