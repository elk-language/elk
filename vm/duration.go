package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Duration
func initDuration() {
	// Class methods
	c := &value.DurationClass.SingletonClass().MethodContainer
	Def(
		c,
		"parse",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			str := args[1].(value.String)
			return value.ParseDuration(str)
		},
		DefWithParameters(1),
	)

	// Instance methods
	c = &value.DurationClass.MethodContainer
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			dur := args[1].(value.Duration)
			return self.Add(dur), nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			dur := args[1].(value.Duration)
			return self.Subtract(dur), nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			other := args[1]
			return self.Multiply(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			other := args[1]
			return self.Divide(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.ToString(), nil
		},
	)
	Def(
		c,
		"nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.Nanoseconds(), nil
		},
	)
	Def(
		c,
		"in_nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.InNanoseconds(), nil
		},
	)
	Def(
		c,
		"microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.Microseconds(), nil
		},
	)
	Def(
		c,
		"in_microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.InMicroseconds(), nil
		},
	)
	Def(
		c,
		"milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.Milliseconds(), nil
		},
	)
	Def(
		c,
		"in_milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.InMilliseconds(), nil
		},
	)
	Def(
		c,
		"seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.Seconds(), nil
		},
	)
	Def(
		c,
		"in_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.InSeconds(), nil
		},
	)
	Def(
		c,
		"minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.Minutes(), nil
		},
	)
	Def(
		c,
		"in_minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.InMinutes(), nil
		},
	)
	Def(
		c,
		"hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.Hours(), nil
		},
	)
	Def(
		c,
		"in_hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.InHours(), nil
		},
	)
	Def(
		c,
		"days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.Days(), nil
		},
	)
	Def(
		c,
		"in_days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.InDays(), nil
		},
	)
	Def(
		c,
		"weeks",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.Weeks(), nil
		},
	)
	Def(
		c,
		"in_weeks",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.InWeeks(), nil
		},
	)
	Def(
		c,
		"years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.Years(), nil
		},
	)
	Def(
		c,
		"in_years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.InYears(), nil
		},
	)
}
