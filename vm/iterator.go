package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Iterate over an iterable value
func Iterate(vm *Thread, iterableVal value.Value) iter.Seq2[value.Value, value.Value] {
	switch c := iterableVal.SafeAsReference().(type) {
	case value.NativeIterator:
		return value.IterateNativeIterator(c)
	case value.NativeIterable:
		return c.Iterate()
	default:
		return func(yield func(value.Value, value.Value) bool) {
			iterator, err := vm.CallMethodByName(symbol.L_iter, iterableVal)
			if !err.IsUndefined() {
				yield(value.Undefined, err)
				return
			}

			for element, err := range Iterate(vm, iterator) {
				if !yield(element, value.Undefined) {
					return
				}

				if err != symbol.L_stop_iteration.ToValue() {
					yield(value.Undefined, err)
				}
				return
			}
		}
	}
}

func initIterator() {
	// ::Std::Iterator::Base
	// Instance methods
	c := &value.IteratorBaseMixin.MethodContainer
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
}
