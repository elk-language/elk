package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// ::Std::Kernel
func initKernel() {
	// Instance methods
	c := &value.KernelModule.SingletonClass().MethodContainer
	Def(
		c,
		"inspect_stack",
		func(v *VM, args []value.Value) (value.Value, value.Value) {
			v.InspectStack()
			return value.Nil, nil
		},
	)
	Def(
		c,
		"print",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			values := args[1].(*value.ArrayTuple)
			for _, val := range *values {
				result, err := vm.CallMethodByName(toStringSymbol, val)
				if err != nil {
					return nil, err
				}
				fmt.Fprint(vm.Stdout, result)
			}

			return value.Nil, nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"println",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			values := args[1].(*value.ArrayTuple)
			for _, val := range *values {
				result, err := vm.CallMethodByName(toStringSymbol, val)
				if err != nil {
					return nil, err
				}
				fmt.Fprintln(vm.Stdout, result)
			}

			return value.Nil, nil
		},
		DefWithParameters(1),
	)
	Alias(c, "puts", "println")
}
