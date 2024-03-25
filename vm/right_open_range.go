package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::rightOpenRange
func init() {
	// Instance methods
	c := &value.RightOpenRangeClass.MethodContainer
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.RightOpenRange)
			iterator := value.NewRightOpenRangeIterator(self)
			return iterator, nil
		},
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.RightOpenRange)
			other, ok := args[1].(*value.RightOpenRange)
			if !ok {
				return value.False, nil
			}
			equal, err := RightOpenRangeEqual(vm, self, other)
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
			self := args[0].(*value.RightOpenRange)
			other := args[1]
			contains, err := RightOpenRangeContains(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(contains), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
}

// ::Std::RightOpenRange::Iterator
func init() {
	// Instance methods
	c := &value.RightOpenRangeIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.RightOpenRangeIterator)
			return RightOpenRangeIteratorNext(vm, self)
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

// Checks whether a value is contained in the range
func RightOpenRangeContains(vm *VM, r *value.RightOpenRange, val value.Value) (bool, value.Value) {
	eqVal, err := GreaterThanEqual(vm, val, r.From)
	if err != nil {
		return false, err
	}

	if value.Falsy(eqVal) {
		return false, nil
	}

	eqVal, err = LessThan(vm, val, r.To)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}

// Checks whether two right open ranges are equal
func RightOpenRangeEqual(vm *VM, x, y *value.RightOpenRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.From, y.From)
	if err != nil {
		return false, err
	}

	if value.Falsy(eqVal) {
		return false, nil
	}

	eqVal, err = Equal(vm, x.To, y.To)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}

// Get the next element of the range
func RightOpenRangeIteratorNext(vm *VM, i *value.RightOpenRangeIterator) (value.Value, value.Value) {
	greater, err := GreaterThanEqual(vm, i.CurrentElement, i.Range.To)
	if err != nil {
		return nil, err
	}

	if value.Truthy(greater) {
		return nil, stopIterationSymbol
	}

	current := i.CurrentElement

	// i.CurrentElement++
	next, err := Increment(vm, i.CurrentElement)
	if err != nil {
		return nil, err
	}
	i.CurrentElement = next

	return current, nil
}
