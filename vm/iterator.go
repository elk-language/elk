package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Iterate over an iterable value
func Iterate(vm *Thread, iterableVal value.Value) iter.Seq2[value.Value, value.Value] {
	switch c := iterableVal.ToInterface().(type) {
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

			for element, err := range IterateIterator(vm, iterator) {
				if !err.IsUndefined() {
					yield(value.Undefined, err)
					return
				}

				if !yield(element, value.Undefined) {
					return
				}
			}
		}
	}
}

// Iterate over an iterator
func IterateIterator(vm *Thread, iteratorVal value.Value) iter.Seq2[value.Value, value.Value] {
	switch c := iteratorVal.ToInterface().(type) {
	case value.NativeIterator:
		return value.IterateNativeIterator(c)
	default:
		return func(yield func(value.Value, value.Value) bool) {
			for {
				element, err := vm.CallMethodByName(symbol.L_next, iteratorVal)
				if err.IsUndefined() {
					if !yield(element, value.Undefined) {
						return
					}
					continue
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
