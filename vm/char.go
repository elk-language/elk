package vm

import (
	"github.com/elk-language/elk/value"
)

func initChar() {
	// Instance methods
	c := &value.CharClass.MethodContainer
	Def(
		c,
		"++",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			return (self + 1).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"--",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			return (self - 1).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			other := args[1]
			return value.RefErr(self.Concat(other))
		},
		DefWithParameters(1),
	)
	Alias(c, "concat", "+")

	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			other := args[1]
			return value.RefErr(self.Repeat(other))
		},
		DefWithParameters(1),
	)
	Alias(c, "repeat", "*")

	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			other := args[1]
			return self.CompareVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			other := args[1]
			return self.LessThanEqualVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			other := args[1]
			return self.LessThanVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			other := args[1]
			return self.GreaterThanVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			other := args[1]
			return self.GreaterThanEqualVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			other := args[1]
			return self.EqualVal(other), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"uppercase",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			return self.Uppercase().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"lowercase",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			return self.Lowercase().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			return value.SmallInt(self.CharCount()).ToValue(), value.Undefined
		},
	)
	Alias(c, "char_count", "length")

	Def(
		c,
		"byte_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			return value.SmallInt(self.ByteCount()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"grapheme_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			return value.SmallInt(self.GraphemeCount()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			return value.Ref(value.String(string(self))), value.Undefined
		},
	)
	Def(
		c,
		"to_symbol",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			return value.ToSymbol(string(self)).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustChar()
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)
	Def(
		c,
		"is_empty",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.False, value.Undefined
		},
	)

}
