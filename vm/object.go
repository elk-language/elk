package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

func init() {
	// Instance methods
	c := &value.ObjectClass.MethodContainer
	Def(
		c,
		"print",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			values := args[1].(value.List)
			for _, val := range values {
				fmt.Fprint(vm.Stdout, val)
			}

			return value.Nil, nil
		},
		DefWithParameters("values"),
		DefWithPositionalRestParameter(),
	)
	Def(
		c,
		"println",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			values := args[1].(value.List)
			for _, val := range values {
				fmt.Fprintln(vm.Stdout, val)
			}

			return value.Nil, nil
		},
		DefWithParameters("values"),
		DefWithPositionalRestParameter(),
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return value.String(self.Inspect()), nil
		},
	)
	Def(
		c,
		"class",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return self.Class(), nil
		},
	)

	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			return value.ToElkBool(self == other), nil
		},
		DefWithParameters("other"),
	)
	Alias(c, "===", "==")

}
