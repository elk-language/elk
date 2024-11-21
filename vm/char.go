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
			self := args[0].(value.Char)
			return self + 1, nil
		},
	)
	Def(
		c,
		"--",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			return self - 1, nil
		},
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return self.Concat(other)
		},
		DefWithParameters(1),
	)
	Alias(c, "concat", "+")

	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return self.Repeat(other)
		},
		DefWithParameters(1),
	)
	Alias(c, "repeat", "*")

	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return self.Compare(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return self.LessThanEqual(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return self.LessThan(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return self.GreaterThan(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return self.GreaterThanEqual(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return self.Equal(other), nil
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"uppercase",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			return self.Uppercase(), nil
		},
	)
	Def(
		c,
		"lowercase",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			return self.Lowercase(), nil
		},
	)

	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			return value.SmallInt(self.CharCount()), nil
		},
	)
	Alias(c, "char_count", "length")

	Def(
		c,
		"byte_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			return value.SmallInt(self.ByteCount()), nil
		},
	)
	Def(
		c,
		"grapheme_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			return value.SmallInt(self.GraphemeCount()), nil
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			return value.String(string(self)), nil
		},
	)
	Def(
		c,
		"to_symbol",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			return value.ToSymbol(string(self)), nil
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			return value.String(self.Inspect()), nil
		},
	)
	Def(
		c,
		"is_empty",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.False, nil
		},
	)

}
