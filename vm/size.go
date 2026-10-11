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
		"-",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return (self - other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	// ?: Idk if restricting `other` argument to Size type is good idea, what does 3.bytes * 3.bytes mean?
	// ?: Under the hood Size is a int64, but allowing integers in general should be smarter approach
	// ?: for both `*` and `/` operators.
	Def(
		c,
		"*",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return (self * other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return (self / other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=>",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return value.SmallInt(self.Cmp(other)).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return value.BoolVal(self.GreaterThan(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return value.BoolVal(self.GreaterThanEqual(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return value.BoolVal(self.LessThan(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return value.BoolVal(self.LessThanEqual(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "=~", "==")

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

	Def(
		c,
		"megabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.Megabytes(), value.Undefined
		},
	)
	Def(
		c,
		"in_megabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.InMegabytes().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"gigabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.Gigabytes(), value.Undefined
		},
	)
	Def(
		c,
		"in_gigabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.InGigabytes().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"terabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.Terabytes(), value.Undefined
		},
	)
	Def(
		c,
		"in_terabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.InTerabytes().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"petabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.Petabytes(), value.Undefined
		},
	)
	Def(
		c,
		"in_petabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.InPetabytes().ToValue(), value.Undefined
		},
	)

}
