package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::EndlessOpenRange
func initEndlessOpenRange() {
	// Instance methods
	c := &value.EndlessOpenRangeClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.EndlessOpenRange)
			iterator := value.NewEndlessOpenRangeIterator(self)
			return value.Ref(iterator), value.Nil
		},
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.EndlessOpenRange)
			other, ok := args[1].SafeAsReference().(*value.EndlessOpenRange)
			if !ok {
				return value.False, value.Nil
			}
			equal, err := EndlessOpenRangeEqual(vm, self, other)
			if !err.IsNil() {
				return value.Nil, err
			}
			return value.ToElkBool(equal), value.Nil
		},
		DefWithParameters(1),
	)
	// Special version of `contains` used in pattern matching.
	// Given value has to be an instance of the same class as `start`,
	// otherwise `false` will be returned
	Def(
		c,
		"#contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.EndlessOpenRange)
			other := args[1]
			if !value.IsA(other, self.Start.Class()) {
				return value.False, value.Nil
			}
			contains, err := EndlessOpenRangeContains(vm, self, other)
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
			self := args[0].MustReference().(*value.EndlessOpenRange)
			other := args[1]
			contains, err := EndlessOpenRangeContains(vm, self, other)
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
			return value.False, value.Nil
		},
	)
	Def(
		c,
		"is_left_open",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.True, value.Nil
		},
	)
	Def(
		c,
		"is_right_closed",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.False, value.Nil
		},
	)
	Def(
		c,
		"is_right_open",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.True, value.Nil
		},
	)
	Def(
		c,
		"start",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.EndlessOpenRange)
			return self.Start, value.Nil
		},
	)
	Def(
		c,
		"end",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.Nil, value.Nil
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
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.EndlessOpenRangeIterator)
			return EndlessOpenRangeIteratorNext(vm, self)
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

// Checks whether a value is contained in the endless open range
func EndlessOpenRangeContains(vm *VM, r *value.EndlessOpenRange, val value.Value) (bool, value.Value) {
	eqVal, err := GreaterThan(vm, val, r.Start)
	if !err.IsNil() {
		return false, err
	}

	return value.Truthy(eqVal), value.Nil
}

// Checks whether two Endless open ranges are equal
func EndlessOpenRangeEqual(vm *VM, x, y *value.EndlessOpenRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.Start, y.Start)
	if !err.IsNil() {
		return false, err
	}

	return value.Truthy(eqVal), value.Nil
}

// Get the next element of the range
func EndlessOpenRangeIteratorNext(vm *VM, i *value.EndlessOpenRangeIterator) (value.Value, value.Value) {
	// i.CurrentElement++
	next, err := Increment(vm, i.CurrentElement)
	if !err.IsNil() {
		return value.Nil, err
	}
	i.CurrentElement = next

	current := i.CurrentElement

	return current, value.Nil
}
