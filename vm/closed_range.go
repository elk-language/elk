package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::ClosedRange
func initClosedRange() {
	// Instance methods
	c := &value.ClosedRangeClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ClosedRange)
			iterator := value.NewClosedRangeIterator(self)
			return value.Ref(iterator), value.Nil
		},
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ClosedRange)
			other, ok := args[1].MustReference().(*value.ClosedRange)
			if !ok {
				return value.False, value.Nil
			}
			equal, err := ClosedRangeEqual(vm, self, other)
			if !err.IsNil() {
				return value.Nil, err
			}
			return value.ToElkBool(equal), value.Nil
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
			self := args[0].MustReference().(*value.ClosedRange)
			other := args[1]
			if !value.IsA(other, self.Start.Class()) && !value.IsA(other, self.End.Class()) {
				return value.False, value.Nil
			}
			contains, err := ClosedRangeContains(vm, self, other)
			if !err.IsNil() {
				return value.Nil, err
			}
			return value.ToElkBool(contains), value.Nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ClosedRange)
			other := args[1]
			contains, err := ClosedRangeContains(vm, self, other)
			if !err.IsNil() {
				return value.Nil, err
			}
			return value.ToElkBool(contains), value.Nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"is_left_closed",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.True, value.Nil
		},
	)
	Def(
		c,
		"is_left_open",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.False, value.Nil
		},
	)
	Def(
		c,
		"is_right_closed",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.True, value.Nil
		},
	)
	Def(
		c,
		"is_right_open",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.False, value.Nil
		},
	)
	Def(
		c,
		"start",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ClosedRange)
			return self.Start, value.Nil
		},
	)
	Def(
		c,
		"end",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ClosedRange)
			return self.End, value.Nil
		},
	)
}

// ::Std::ClosedRange::Iterator
func initClosedRangeIterator() {
	// Instance methods
	c := &value.ClosedRangeIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ClosedRangeIterator)
			return ClosedRangeIteratorNext(vm, self)
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Nil
		},
	)

}

// Checks whether a value is contained in the closed range
func ClosedRangeContains(vm *VM, r *value.ClosedRange, val value.Value) (bool, value.Value) {
	eqVal, err := GreaterThanEqual(vm, val, r.Start)
	if !err.IsNil() {
		return false, err
	}

	if value.Falsy(eqVal) {
		return false, value.Nil
	}

	eqVal, err = LessThanEqual(vm, val, r.End)
	if !err.IsNil() {
		return false, err
	}

	return value.Truthy(eqVal), value.Nil
}

// Checks whether two closed ranges are equal
func ClosedRangeEqual(vm *VM, x *value.ClosedRange, y *value.ClosedRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.Start, y.Start)
	if !err.IsNil() {
		return false, err
	}

	if value.Falsy(eqVal) {
		return false, value.Nil
	}

	eqVal, err = Equal(vm, x.End, y.End)
	if !err.IsNil() {
		return false, err
	}

	return value.Truthy(eqVal), value.Nil
}

// Get the next element of the range
func ClosedRangeIteratorNext(vm *VM, i *value.ClosedRangeIterator) (value.Value, value.Value) {
	greater, err := GreaterThan(vm, i.CurrentElement, i.Range.End)
	if !err.IsNil() {
		return value.Nil, err
	}

	if value.Truthy(greater) {
		return value.Nil, value.ToSymbol("stop_iteration").ToValue()
	}

	current := i.CurrentElement

	// i.CurrentElement++
	next, err := Increment(vm, i.CurrentElement)
	if !err.IsNil() {
		return value.Nil, err
	}
	i.CurrentElement = next

	return current, value.Nil
}
