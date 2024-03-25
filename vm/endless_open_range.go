package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::EndlessOpenRange
func init() {
	// Instance methods
	c := &value.EndlessOpenRangeClass.MethodContainer
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.EndlessOpenRange)
			iterator := value.NewEndlessOpenRangeIterator(self)
			return iterator, nil
		},
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.EndlessOpenRange)
			other, ok := args[1].(*value.EndlessOpenRange)
			if !ok {
				return value.False, nil
			}
			equal, err := EndlessOpenRangeEqual(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(equal), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.EndlessOpenRange)
			other := args[1]
			contains, err := EndlessOpenRangeContains(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(contains), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
}

// ::Std::EndlessOpenRange::Iterator
func init() {
	// Instance methods
	c := &value.EndlessOpenRangeIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.EndlessOpenRangeIterator)
			return EndlessOpenRangeIteratorNext(vm, self)
		},
	)
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], nil
		},
	)

}

// Checks whether a value is contained in the endless open range
func EndlessOpenRangeContains(vm *VM, r *value.EndlessOpenRange, val value.Value) (bool, value.Value) {
	eqVal, err := GreaterThan(vm, val, r.From)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}

// Checks whether two Endless open ranges are equal
func EndlessOpenRangeEqual(vm *VM, x, y *value.EndlessOpenRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.From, y.From)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}

// Get the next element of the range
func EndlessOpenRangeIteratorNext(vm *VM, i *value.EndlessOpenRangeIterator) (value.Value, value.Value) {
	// i.CurrentElement++
	next, err := Increment(vm, i.CurrentElement)
	if err != nil {
		return nil, err
	}
	i.CurrentElement = next

	current := i.CurrentElement

	return current, nil
}
