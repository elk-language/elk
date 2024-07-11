package vm

import (
	"github.com/elk-language/elk/value"
)

func init() {
	// Instance methods
	c := &value.SymbolClass.MethodContainer

	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Symbol)
			other := args[1]
			return self.Equal(other), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"to_symbol",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], nil
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Symbol)
			return self.ToString(), nil
		},
	)
	Alias(c, "name", "to_string")

	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Symbol)
			return value.String(self.Inspect()), nil
		},
	)

}
