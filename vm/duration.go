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
			str := args[1].AsReference().(value.String)
			return value.ToValueErr(value.ParseDuration(str))
		},
		DefWithParameters(1),
	)

	// Instance methods
	c = &value.DurationClass.MethodContainer
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			dur := args[1].MustDuration()
			return self.Add(dur).ToValue(), value.Nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			dur := args[1].MustDuration()
			return self.Subtract(dur).ToValue(), value.Nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			other := args[1]
			return value.ToValueErr(self.Multiply(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			other := args[1]
			return value.ToValueErr(self.Divide(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return value.Ref(self.ToString()), value.Nil
		},
	)
	Def(
		c,
		"nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.Nanoseconds(), value.Nil
		},
	)
	Def(
		c,
		"in_nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.InNanoseconds().ToValue(), value.Nil
		},
	)
	Def(
		c,
		"microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.Microseconds(), value.Nil
		},
	)
	Def(
		c,
		"in_microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.InMicroseconds().ToValue(), value.Nil
		},
	)
	Def(
		c,
		"milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.Milliseconds(), value.Nil
		},
	)
	Def(
		c,
		"in_milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.InMilliseconds().ToValue(), value.Nil
		},
	)
	Def(
		c,
		"seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.Seconds(), value.Nil
		},
	)
	Def(
		c,
		"in_seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.InSeconds().ToValue(), value.Nil
		},
	)
	Def(
		c,
		"minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.Minutes(), value.Nil
		},
	)
	Def(
		c,
		"in_minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.InMinutes().ToValue(), value.Nil
		},
	)
	Def(
		c,
		"hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.Hours(), value.Nil
		},
	)
	Def(
		c,
		"in_hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.InHours().ToValue(), value.Nil
		},
	)
	Def(
		c,
		"days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.Days(), value.Nil
		},
	)
	Def(
		c,
		"in_days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.InDays().ToValue(), value.Nil
		},
	)
	Def(
		c,
		"weeks",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.Weeks(), value.Nil
		},
	)
	Def(
		c,
		"in_weeks",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.InWeeks().ToValue(), value.Nil
		},
	)
	Def(
		c,
		"years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.Years(), value.Nil
		},
	)
	Def(
		c,
		"in_years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDuration()
			return self.InYears().ToValue(), value.Nil
		},
	)
}
