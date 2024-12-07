package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Symbol
func initSymbol() {
	// Instance methods
	c := &value.SymbolClass.MethodContainer

	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustSymbol()
			other := args[1]
			return self.Equal(other), value.Nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"to_symbol",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Nil
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustSymbol()
			return value.Ref(self.ToString()), value.Nil
		},
	)
	Alias(c, "name", "to_string")

	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustSymbol()
			return value.Ref(value.String(self.Inspect())), value.Nil
		},
	)

}
