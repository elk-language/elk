package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

func init() {
	DefineMethodWithOptions(
		value.ObjectClass.Methods,
		"print",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			values := args[1].(value.List)
			for _, val := range values {
				fmt.Fprint(vm.Stdout, val)
			}

			return value.Nil, nil
		},
		NativeMethodWithStringParameters("values"),
		NativeMethodWithPositionalRestParameter(),
	)
	DefineMethodWithOptions(
		value.ObjectClass.Methods,
		"println",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			values := args[1].(value.List)
			for _, val := range values {
				fmt.Fprintln(vm.Stdout, val)
			}

			return value.Nil, nil
		},
		NativeMethodWithStringParameters("values"),
		NativeMethodWithPositionalRestParameter(),
	)
	DefineMethodWithOptions(
		value.ObjectClass.Methods,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return value.String(self.Inspect()), nil
		},
	)
	DefineMethodWithOptions(
		value.ObjectClass.Methods,
		"class",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return self.Class(), nil
		},
	)

	DefineMethodWithOptions(
		value.ObjectClass.Methods,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			return value.ToElkBool(self == other), nil
		},
		NativeMethodWithStringParameters("other"),
	)
	value.ObjectClass.DefineAliasString("===", "==")

}
