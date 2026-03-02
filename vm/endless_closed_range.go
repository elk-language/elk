package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// ::Std::EndlessClosedRange
func initEndlessClosedRange() {
	// Instance methods
	c := &value.EndlessClosedRangeClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.EndlessClosedRange)
			iterator := value.NewEndlessClosedRangeIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)
	Def(
		c,
		"==",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.EndlessClosedRange)
			other, ok := args[1].SafeAsReference().(*value.EndlessClosedRange)
			if !ok {
				return value.False.ToValue(), value.Undefined
			}
			equal, err := EndlessClosedRangeEqual(vm, self, other)
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
			self := args[0].MustReference().(*value.EndlessClosedRange)
			other := args[1]
			if !value.IsA(other, self.Start.Class()) {
				return value.False.ToValue(), value.Undefined
			}
			contains, err := EndlessClosedRangeContains(vm, self, other)
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
			self := args[0].MustReference().(*value.EndlessClosedRange)
			other := args[1]
			contains, err := EndlessClosedRangeContains(vm, self, other)
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
			return value.True.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"is_left_open",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return value.False.ToValue(), value.Undefined
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
			self := args[0].MustReference().(*value.EndlessClosedRange)
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

// ::Std::EndlessClosedRange::Iterator
func initEndlessClosedRangeIterator() {
	// Instance methods
	c := &value.EndlessClosedRangeIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.EndlessClosedRangeIterator)(args[0].Pointer())
			return EndlessClosedRangeIteratorNext(vm, self)
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
			self := (*value.EndlessClosedRangeIterator)(args[0].Pointer())
			self.Reset()
			return args[0], value.Undefined
		},
	)

}

// Checks whether a value is contained in the closed range
func EndlessClosedRangeContains(vm *Thread, r *value.EndlessClosedRange, val value.Value) (bool, value.Value) {
	eqVal, err := GreaterThanEqual(vm, val, r.Start)
	if !err.IsUndefined() {
		return false, err
	}

	return value.Truthy(eqVal), value.Undefined
}

// Checks whether two closed ranges are equal
func EndlessClosedRangeEqual(vm *Thread, x *value.EndlessClosedRange, y *value.EndlessClosedRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.Start, y.Start)
	if !err.IsUndefined() {
		return false, err
	}

	return value.Truthy(eqVal), value.Undefined
}

// Get the next element of the range
func EndlessClosedRangeIteratorNext(vm *Thread, i *value.EndlessClosedRangeIterator) (value.Value, value.Value) {
	current := i.CurrentElement

	// i.CurrentElement++
	next, err := Increment(vm, i.CurrentElement)
	if !err.IsUndefined() {
		return value.Undefined, err
	}
	i.CurrentElement = next

	return current, value.Undefined
}

// Iterate over all elements of the iterator
func EndlessClosedRangeIteratorAll(vm *Thread, i *value.EndlessClosedRangeIterator) iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		for {
			element, err := EndlessClosedRangeIteratorNext(vm, i)
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
