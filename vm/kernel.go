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
			for val, err := range Iterate(vm, args[1]) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

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
			for val, err := range Iterate(vm, args[1]) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

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
			var duration value.TimeSpan
			if durationVal.IsReference() {
				duration = durationVal.AsReference().(value.TimeSpan)
			} else {
				duration = durationVal.AsInlineTimeSpan()
			}

			time.Sleep(duration.Go())

			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"timeout",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			durationVal := args[1]
			var duration value.TimeSpan
			if durationVal.IsReference() {
				duration = durationVal.AsReference().(value.TimeSpan)
			} else {
				duration = durationVal.AsInlineTimeSpan()
			}

			p := NewExternalPromise(vm.threadPool)

			go func(p *Promise, d time.Duration) {
				<-time.After(d)
				p.Resolve(value.Nil)
			}(p, duration.Go())

			return value.Ref(p), value.Undefined
		},
		DefWithParameters(1),
	)
}
