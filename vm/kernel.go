package vm

import (
	"fmt"
	"os"
	"time"

	"github.com/elk-language/elk/value"
)

// ::Std::Kernel
func initKernel() {
	c := &value.KernelModule.SingletonClass().MethodContainer
	Def(
		c,
		"exit",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			var code int
			if !args[1].IsUndefined() {
				code = args[1].AsInt()
			}
			os.Exit(code)
			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"print",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			for val, err := range Iterate(vm, args[1]) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				result, err := ToString(vm, val)
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
		"print@1",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			val := args[1]
			result, err := ToString(vm, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			r := result.MustReference().(value.String).String()
			fmt.Fprint(vm.Stdout, r)

			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"println",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			var iterated bool
			for val, err := range Iterate(vm, args[1]) {
				iterated = true

				if !err.IsUndefined() {
					return value.Undefined, err
				}

				result, err := ToString(vm, val)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				r := result.MustReference().(value.String).String()
				fmt.Fprintln(vm.Stdout, r)
			}

			if !iterated {
				fmt.Fprintln(vm.Stdout)
			}

			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"println@1",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			val := args[1]

			result, err := ToString(vm, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			r := result.MustReference().(value.String).String()
			fmt.Fprintln(vm.Stdout, r)

			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "puts", "println")
	Alias(c, "puts@1", "println@1")

	Def(
		c,
		"sleep",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			durationVal := args[1]
			var duration value.TimeSpan
			if durationVal.IsReference() {
				duration = durationVal.AsReference().(value.TimeSpan)
			} else {
				duration = durationVal.AsInlineTimeSpan()
			}

			time.Sleep(duration.Native())

			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"timeout",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
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
			}(p, duration.Native())

			return value.Ref(p), value.Undefined
		},
		DefWithParameters(2),
	)
}
