package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Iterate over an iterable value
func Iterate(vm *Thread, collectionValue value.Value) iter.Seq2[value.Value, value.Value] {
	switch c := collectionValue.AsReference().(type) {
	case value.NativeIterator:
		return value.IterateNativeIterator(c)
	case value.NativeIterable:
		return c.Iterate()
	case *value.HashSet:
		return func(yield func(value.Value, value.Value) bool) {
			for _, element := range c.Table {
				if element.IsUndefined() {
					continue
				}

				if !yield(element, value.Undefined) {
					return
				}
			}
		}
	case *value.HashMap:
		return func(yield func(value.Value, value.Value) bool) {
			for index, _ := range c.Table {
				pair := &c.Table[index]
				if pair.Key.IsUndefined() {
					continue
				}

				if !yield(value.Ref(pair), value.Undefined) {
					return
				}
			}
		}
	case *value.HashRecord:
		return func(yield func(value.Value, value.Value) bool) {
			for index, _ := range c.Table {
				pair := &c.Table[index]
				if pair.Key.IsUndefined() {
					continue
				}

				if !yield(value.Ref(pair), value.Undefined) {
					return
				}
			}
		}
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
