package vm

import (
	"github.com/elk-language/elk/value"
)

func init() {
	DefineMethodWithOptions(
		value.StringClass.Methods,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.Concat(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	value.StringClass.DefineAliasString("concat", "+")

	DefineMethodWithOptions(
		value.StringClass.Methods,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.RemoveSuffix(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	value.StringClass.DefineAliasString("remove_suffix", "-")

	DefineMethodWithOptions(
		value.StringClass.Methods,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.Repeat(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	value.StringClass.DefineAliasString("repeat", "*")

	DefineMethodWithOptions(
		value.StringClass.Methods,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.Compare(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.StringClass.Methods,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.LessThanEqual(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.StringClass.Methods,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.LessThan(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.StringClass.Methods,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.GreaterThan(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.StringClass.Methods,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.GreaterThanEqual(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.StringClass.Methods,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return self.Equal(other), nil
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.StringClass.Methods,
		"===",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return self.StrictEqual(other), nil
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)

	DefineMethodWithOptions(
		value.StringClass.Methods,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.SmallInt(self.CharCount()), nil
		},
	)
	value.StringClass.DefineAliasString("char_count", "length")

	DefineMethodWithOptions(
		value.StringClass.Methods,
		"byte_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.SmallInt(self.ByteCount()), nil
		},
	)
	DefineMethodWithOptions(
		value.StringClass.Methods,
		"grapheme_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.SmallInt(self.GraphemeCount()), nil
		},
	)
	DefineMethodWithOptions(
		value.StringClass.Methods,
		"to_symbol",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.ToSymbol(self), nil
		},
	)
	DefineMethodWithOptions(
		value.StringClass.Methods,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.String(self.Inspect()), nil
		},
	)
	DefineMethodWithOptions(
		value.StringClass.Methods,
		"is_empty",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.ToElkBool(self.IsEmpty()), nil
		},
	)

}
