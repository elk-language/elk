package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Date::Span
func initSize() {
	// Instance methods
	c := &value.SizeClass.MethodContainer

	Def(
		c,
		"-@",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return (-self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"+@",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)

	Def(
		c,
		"+",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return (self + other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"bytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.Bytes(), value.Undefined
		},
	)
	Def(
		c,
		"in_bytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.InBytes().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"kilobytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.Kilobytes(), value.Undefined
		},
	)
	Def(
		c,
		"in_kilobytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.InKilobytes().ToValue(), value.Undefined
		},
	)
}
