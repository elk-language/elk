package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// ::Std::EndlessOpenRange
func initEndlessOpenRange() {
	// Instance methods
	c := &value.EndlessOpenRangeClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.EndlessOpenRange)
			iterator := value.NewEndlessOpenRangeIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)
	Def(
		c,
		"==",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.EndlessOpenRange)
			other, ok := args[1].SafeAsReference().(*value.EndlessOpenRange)
			if !ok {
				return value.False.ToValue(), value.Undefined
			}
			equal, err := EndlessOpenRangeEqual(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.BoolVal(equal), value.Undefined
		},
		DefWithParameters(1),
	)
	// Special version of `contains` used in pattern matching.
	// Given value has to be an instance of the same class as `start`,
	// otherwise `false` will be returned
	Def(
		c,
		"#contains",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.EndlessOpenRange)
			other := args[1]
			if !value.IsA(other, self.Start.Class()) {
				return value.False.ToValue(), value.Undefined
			}
			contains, err := EndlessOpenRangeContains(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.BoolVal(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.EndlessOpenRange)
			other := args[1]
			contains, err := EndlessOpenRangeContains(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.BoolVal(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"is_left_closed",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return value.False.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"is_left_open",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return value.True.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"is_right_closed",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return value.False.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"is_right_open",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return value.True.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"start",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.EndlessOpenRange)
			return self.Start, value.Undefined
		},
	)
	Def(
		c,
		"end",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return value.Undefined, value.Undefined
		},
	)
}

// ::Std::EndlessOpenRange::Iterator
func initEndlessOpenRangeIterator() {
	// Instance methods
	c := &value.EndlessOpenRangeIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.EndlessOpenRangeIterator)(args[0].Pointer())
			return EndlessOpenRangeIteratorNext(vm, self)
		},
	)
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"reset",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.EndlessOpenRangeIterator)(args[0].Pointer())
			self.Reset()
			return args[0], value.Undefined
		},
	)

}

// Checks whether a value is contained in the endless open range
func EndlessOpenRangeContains(vm *Thread, r *value.EndlessOpenRange, val value.Value) (bool, value.Value) {
	eqVal, err := GreaterThan(vm, val, r.Start)
	if !err.IsUndefined() {
		return false, err
	}

	return value.Truthy(eqVal), value.Undefined
}

// Checks whether two Endless open ranges are equal
func EndlessOpenRangeEqual(vm *Thread, x, y *value.EndlessOpenRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.Start, y.Start)
	if !err.IsUndefined() {
		return false, err
	}

	return value.Truthy(eqVal), value.Undefined
}

// Get the next element of the range
func EndlessOpenRangeIteratorNext(vm *Thread, i *value.EndlessOpenRangeIterator) (value.Value, value.Value) {
	// i.CurrentElement++
	next, err := Increment(vm, i.CurrentElement)
	if !err.IsUndefined() {
		return value.Undefined, err
	}
	i.CurrentElement = next

	current := i.CurrentElement

	return current, value.Undefined
}

// Iterate over all elements of the iterator
func EndlessOpenRangeIteratorAll(vm *Thread, i *value.EndlessOpenRangeIterator) iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		for {
			element, err := EndlessOpenRangeIteratorNext(vm, i)
			if err.IsInlineSymbol() {
				if element.AsInlineSymbol() == symbol.L_stop_iteration {
					break
				}
			}

			if !yield(element, err) {
				return
			}
		}
	}
}
