package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::EndlessClosedRange
func init() {
	// Instance methods
	c := &value.EndlessClosedRangeClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.EndlessClosedRange)
			iterator := value.NewEndlessClosedRangeIterator(self)
			return iterator, nil
		},
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.EndlessClosedRange)
			other, ok := args[1].(*value.EndlessClosedRange)
			if !ok {
				return value.False, nil
			}
			equal, err := EndlessClosedRangeEqual(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(equal), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	// Special version of `contains` used in pattern matching.
	// Given value has to be an instance of the same class as `start`,
	// otherwise `false` will be returned
	Def(
		c,
		"#contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.EndlessClosedRange)
			other := args[1]
			if !value.IsA(other, self.Start.Class()) {
				return value.False, nil
			}
			contains, err := EndlessClosedRangeContains(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(contains), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.EndlessClosedRange)
			other := args[1]
			contains, err := EndlessClosedRangeContains(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(contains), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"is_left_closed",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.True, nil
		},
		DefWithSealed(),
	)
	Def(
		c,
		"is_left_open",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.False, nil
		},
		DefWithSealed(),
	)
	Def(
		c,
		"is_right_closed",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.False, nil
		},
		DefWithSealed(),
	)
	Def(
		c,
		"is_right_open",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.True, nil
		},
		DefWithSealed(),
	)
	Def(
		c,
		"start",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.EndlessClosedRange)
			return self.Start, nil
		},
		DefWithSealed(),
	)
	Def(
		c,
		"end",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.Nil, nil
		},
		DefWithSealed(),
	)
}

// ::Std::EndlessClosedRange::Iterator
func init() {
	// Instance methods
	c := &value.EndlessClosedRangeIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.EndlessClosedRangeIterator)
			return EndlessClosedRangeIteratorNext(vm, self)
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], nil
		},
	)

}

// Checks whether a value is contained in the closed range
func EndlessClosedRangeContains(vm *VM, r *value.EndlessClosedRange, val value.Value) (bool, value.Value) {
	eqVal, err := GreaterThanEqual(vm, val, r.Start)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}

// Checks whether two closed ranges are equal
func EndlessClosedRangeEqual(vm *VM, x *value.EndlessClosedRange, y *value.EndlessClosedRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.Start, y.Start)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}

// Get the next element of the range
func EndlessClosedRangeIteratorNext(vm *VM, i *value.EndlessClosedRangeIterator) (value.Value, value.Value) {
	current := i.CurrentElement

	// i.CurrentElement++
	next, err := Increment(vm, i.CurrentElement)
	if err != nil {
		return nil, err
	}
	i.CurrentElement = next

	return current, nil
}
