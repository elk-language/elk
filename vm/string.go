package vm

import (
	"github.com/elk-language/elk/value"
)

func init() {
	// Instance methods
	c := &value.StringClass.MethodContainer
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.Concat(other))
		},
		DefWithParameters("other"),
		DefWithFrozen(),
	)
	value.StringClass.DefineAliasString("concat", "+")

	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.RemoveSuffix(other))
		},
		DefWithParameters("other"),
		DefWithFrozen(),
	)
	value.StringClass.DefineAliasString("remove_suffix", "-")

	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.Repeat(other))
		},
		DefWithParameters("other"),
		DefWithFrozen(),
	)
	value.StringClass.DefineAliasString("repeat", "*")

	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.Compare(other))
		},
		DefWithParameters("other"),
		DefWithFrozen(),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.LessThanEqual(other))
		},
		DefWithParameters("other"),
		DefWithFrozen(),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.LessThan(other))
		},
		DefWithParameters("other"),
		DefWithFrozen(),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.GreaterThan(other))
		},
		DefWithParameters("other"),
		DefWithFrozen(),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.GreaterThanEqual(other))
		},
		DefWithParameters("other"),
		DefWithFrozen(),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return self.Equal(other), nil
		},
		DefWithParameters("other"),
		DefWithFrozen(),
	)
	Def(
		c,
		"===",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return self.StrictEqual(other), nil
		},
		DefWithParameters("other"),
		DefWithFrozen(),
	)

	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.SmallInt(self.CharCount()), nil
		},
	)
	value.StringClass.DefineAliasString("char_count", "length")

	Def(
		c,
		"byte_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.SmallInt(self.ByteCount()), nil
		},
	)
	Def(
		c,
		"grapheme_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.SmallInt(self.GraphemeCount()), nil
		},
	)
	Def(
		c,
		"to_symbol",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.ToSymbol(self), nil
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.String(self.Inspect()), nil
		},
	)
	Def(
		c,
		"is_empty",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.ToElkBool(self.IsEmpty()), nil
		},
	)

}
