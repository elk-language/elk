package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// ::Std::OpenRange
func initOpenRange() {
	// Instance methods
	c := &value.OpenRangeClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.OpenRange)
			iterator := value.NewOpenRangeIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)
	Def(
		c,
		"==",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.OpenRange)
			other, ok := args[1].SafeAsReference().(*value.OpenRange)
			if !ok {
				return value.False.ToValue(), value.Undefined
			}
			equal, err := OpenRangeEqual(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.BoolVal(equal), value.Undefined
		},
		DefWithParameters(1),
	)
	// Special version of `contains` used in pattern matching.
	// Given value has to be an instance of the same class as `start` or `end`,
	// otherwise `false` will be returned
	Def(
		c,
		"#contains",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.OpenRange)
			other := args[1]
			if !value.IsA(other, self.Start.Class()) && !value.IsA(other, self.End.Class()) {
				return value.False.ToValue(), value.Undefined
			}
			contains, err := OpenRangeContains(vm, self, other)
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
			self := args[0].MustReference().(*value.OpenRange)
			other := args[1]
			contains, err := OpenRangeContains(vm, self, other)
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
			self := args[0].MustReference().(*value.OpenRange)
			return self.Start, value.Undefined
		},
	)
	Def(
		c,
		"end",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.OpenRange)
			return self.End, value.Undefined
		},
	)
}

// ::Std::OpenRange::Iterator
func initOpenRangeIterator() {
	// Instance methods
	c := &value.OpenRangeIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.OpenRangeIterator)(args[0].Pointer())
			return OpenRangeIteratorNext(vm, self)
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
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.OpenRangeIterator)(args[0].Pointer())
			self.Reset()
			return args[0], value.Undefined
		},
	)

}

// Checks whether a value is contained in the open range
func OpenRangeContains(vm *Thread, r *value.OpenRange, val value.Value) (bool, value.Value) {
	eqVal, err := GreaterThan(vm, val, r.Start)
	if !err.IsUndefined() {
		return false, err
	}

	if value.Falsy(eqVal) {
		return false, value.Undefined
	}

	eqVal, err = LessThan(vm, val, r.End)
	if !err.IsUndefined() {
		return false, err
	}

	return value.Truthy(eqVal), value.Undefined
}

// Checks whether two open ranges are equal
func OpenRangeEqual(vm *Thread, x *value.OpenRange, y *value.OpenRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.Start, y.Start)
	if !err.IsUndefined() {
		return false, err
	}

	if value.Falsy(eqVal) {
		return false, value.Undefined
	}

	eqVal, err = Equal(vm, x.End, y.End)
	if !err.IsUndefined() {
		return false, err
	}

	return value.Truthy(eqVal), value.Undefined
}

// Get the next element of the range
func OpenRangeIteratorNext(vm *Thread, i *value.OpenRangeIterator) (value.Value, value.Value) {
	// i.CurrentElement++
	next, err := Increment(vm, i.CurrentElement)
	if !err.IsUndefined() {
		return value.Undefined, err
	}
	i.CurrentElement = next

	greater, err := GreaterThanEqual(vm, i.CurrentElement, i.Range.End)
	if !err.IsUndefined() {
		return value.Undefined, err
	}

	if value.Truthy(greater) {
		return value.Undefined, stopIterationSymbol.ToValue()
	}

	current := i.CurrentElement

	return current, value.Undefined
}

// Iterate over all elements of the iterator
func OpenRangeIteratorAll(vm *Thread, i *value.OpenRangeIterator) iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		for {
			element, err := OpenRangeIteratorNext(vm, i)
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
