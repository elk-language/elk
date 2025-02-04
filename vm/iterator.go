package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Iterate over an iterable value
func Iterate(vm *VM, collectionValue value.Value) iter.Seq2[value.Value, value.Value] {
	var collection []value.Value

	switch c := collectionValue.AsReference().(type) {
	case *value.ArrayList:
		collection = *(c)
	case *value.ArrayTuple:
		collection = *(c)
	default:
		return func(yield func(value.Value, value.Value) bool) {
			iterator, err := vm.CallMethodByName(symbol.L_iter, collectionValue)
			if !err.IsUndefined() {
				yield(value.Undefined, err)
				return
			}

			for {
				element, err := vm.CallMethodByName(symbol.L_next, iterator)
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

	return func(yield func(value.Value, value.Value) bool) {
		for _, element := range collection {
			if !yield(element, value.Undefined) {
				return
			}
		}
	}
}
