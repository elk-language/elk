package vm

import (
	"github.com/elk-language/elk/value"
)

func init() {
	// Instance methods
	c := &value.CharClass.MethodContainer
	Def(
		c,
		"++",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			return self + 1, nil
		},
		DefWithSealed(),
	)
	Def(
		c,
		"--",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			return self - 1, nil
		},
		DefWithSealed(),
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return value.ToValueErr(self.Concat(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Alias(c, "concat", "+")

	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return value.ToValueErr(self.Repeat(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Alias(c, "repeat", "*")

	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return value.ToValueErr(self.Compare(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return value.ToValueErr(self.LessThanEqual(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return value.ToValueErr(self.LessThan(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return value.ToValueErr(self.GreaterThan(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return value.ToValueErr(self.GreaterThanEqual(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return self.Equal(other), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"===",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Char)
			other := args[1]
			return self.StrictEqual(other), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
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
