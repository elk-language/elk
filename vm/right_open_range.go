package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::RightOpenRange
func initRightOpenRange() {
	// Instance methods
	c := &value.RightOpenRangeClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.RightOpenRange)
			iterator := value.NewRightOpenRangeIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.RightOpenRange)
			other, ok := args[1].SafeAsReference().(*value.RightOpenRange)
			if !ok {
				return value.False, value.Undefined
			}
			equal, err := RightOpenRangeEqual(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.ToElkBool(equal), value.Undefined
		},
		DefWithParameters(1),
	)
	// Special version of `contains` used in pattern matching.
	// Given value has to be an instance of the same class as `start` or `end`,
	// otherwise `false` will be returned
	Def(
		c,
		"#contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.RightOpenRange)
			other := args[1]
			if !value.IsA(other, self.Start.Class()) && !value.IsA(other, self.End.Class()) {
				return value.False, value.Undefined
			}
			contains, err := RightOpenRangeContains(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.ToElkBool(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.RightOpenRange)
			other := args[1]
			contains, err := RightOpenRangeContains(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.ToElkBool(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"is_left_closed",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.True, value.Undefined
		},
	)
	Def(
		c,
		"is_left_open",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.False, value.Undefined
		},
	)
	Def(
		c,
		"is_right_closed",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.False, value.Undefined
		},
	)
	Def(
		c,
		"is_right_open",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.True, value.Undefined
		},
	)
	Def(
		c,
		"start",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.RightOpenRange)
			return self.Start, value.Undefined
		},
	)
	Def(
		c,
		"end",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.RightOpenRange)
			return self.End, value.Undefined
		},
	)
}

// ::Std::RightOpenRange::Iterator
func initRightOpenRangeIterator() {
	// Instance methods
	c := &value.RightOpenRangeIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.RightOpenRangeIterator)(args[0].Pointer())
			return RightOpenRangeIteratorNext(vm, self)
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"reset",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.RightOpenRangeIterator)(args[0].Pointer())
			self.Reset()
			return args[0], value.Undefined
		},
	)

}

// Checks whether a value is contained in the range
func RightOpenRangeContains(vm *VM, r *value.RightOpenRange, val value.Value) (bool, value.Value) {
	eqVal, err := GreaterThanEqual(vm, val, r.Start)
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

// Checks whether two right open ranges are equal
func RightOpenRangeEqual(vm *VM, x, y *value.RightOpenRange) (bool, value.Value) {
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
func RightOpenRangeIteratorNext(vm *VM, i *value.RightOpenRangeIterator) (value.Value, value.Value) {
	// i.CurrentElement >= i.Range.End
	greater, err := GreaterThanEqual(vm, i.CurrentElement, i.Range.End)
	if !err.IsUndefined() {
		return value.Undefined, err
	}

	if value.Truthy(greater) {
		return value.Undefined, stopIterationSymbol.ToValue()
	}

	current := i.CurrentElement

	// i.CurrentElement++
	next, err := Increment(vm, i.CurrentElement)
	if !err.IsUndefined() {
		return value.Undefined, err
	}
	i.CurrentElement = next

	return current, value.Undefined
}
