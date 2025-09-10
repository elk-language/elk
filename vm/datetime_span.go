package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Datetime::Span
func initDateTimeSpan() {
	// Instance methods
	c := &value.DateTimeSpanClass.MethodContainer

	Def(
		c,
		"-@",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.Ref(self.Negate()), value.Undefined
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
			self := (*value.DateTimeSpan)(args[0].Pointer())
			other := args[1]
			return self.Add(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+@1",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			other := (*value.DateTimeSpan)(args[1].Pointer())
			return value.Ref(self.AddDateTimeSpan(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+@2",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			other := args[1].AsTimeSpan()
			return value.Ref(self.AddTimeSpan(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+@3",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			other := args[1].AsDateSpan()
			return value.Ref(self.AddDateSpan(other)), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			other := args[1]
			return self.Subtract(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@1",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			other := (*value.DateTimeSpan)(args[1].Pointer())
			return value.Ref(self.SubtractDateTimeSpan(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@2",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			other := args[1].AsTimeSpan()
			return value.Ref(self.SubtractTimeSpan(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@3",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			other := args[1].AsDateSpan()
			return value.Ref(self.SubtractDateSpan(other)), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"date_span",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return self.DateSpan.ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"time_span",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return self.TimeSpan.ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.Years()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.InYears()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return self.TotalYears(), value.Undefined
		},
	)
	Def(
		c,
		"months",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.Months()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_months",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.InMonths()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_months",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return self.TotalMonths(), value.Undefined
		},
	)
	Def(
		c,
		"days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.Days()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.InDays()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return self.TotalDays(), value.Undefined
		},
	)
	Def(
		c,
		"hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.Hours()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.InHours()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return self.TotalHours(), value.Undefined
		},
	)
	Def(
		c,
		"minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.Minutes()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.InMinutes()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return self.TotalMinutes(), value.Undefined
		},
	)
	Def(
		c,
		"seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.Seconds()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.InSeconds()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return self.TotalSeconds(), value.Undefined
		},
	)
	Def(
		c,
		"milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.Milliseconds()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.InMilliseconds()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return self.TotalMilliseconds(), value.Undefined
		},
	)
	Def(
		c,
		"microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.Microseconds()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.InMicroseconds()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return self.TotalMicroseconds(), value.Undefined
		},
	)
	Def(
		c,
		"nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.Nanoseconds()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"in_nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.SmallInt(self.InNanoseconds()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"total_nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return self.TotalNanoseconds(), value.Undefined
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DateTimeSpan)(args[0].Pointer())
			return value.Ref(self.ToString()), value.Undefined
		},
	)
}
