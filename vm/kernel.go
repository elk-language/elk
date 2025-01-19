package vm

import (
	"fmt"
	"time"

	"github.com/elk-language/elk/value"
)

// ::Std::Kernel
func initKernel() {
	c := &value.KernelModule.SingletonClass().MethodContainer
	Def(
		c,
		"print",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			values := args[1].MustReference().(*value.ArrayTuple)
			for _, val := range *values {
				result, err := vm.CallMethodByName(toStringSymbol, val)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				r := result.MustReference().(value.String).String()
				fmt.Fprint(vm.Stdout, r)
			}

			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"println",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			values := args[1].MustReference().(*value.ArrayTuple)
			for _, val := range *values {
				result, err := vm.CallMethodByName(toStringSymbol, val)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				r := result.MustReference().(value.String).String()
				fmt.Fprintln(vm.Stdout, r)
			}

			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "puts", "println")

	Def(
		c,
		"sleep",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			durationVal := args[1]
			var duration value.Duration
			if durationVal.IsReference() {
				duration = durationVal.AsReference().(value.Duration)
			} else {
				duration = durationVal.AsDuration()
			}

			time.Sleep(duration.Go())

			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)
}
