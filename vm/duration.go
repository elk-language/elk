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
		"microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.Microseconds(), nil
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
		"seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.Seconds(), nil
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
		"hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Duration)
			return self.Hours(), nil
		},
	)
}
