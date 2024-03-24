package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::ClosedRange
func init() {
	// Instance methods
	c := &value.ClosedRangeClass.MethodContainer
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ClosedRange)
			iterator := value.NewClosedRangeIterator(self)
			return iterator, nil
		},
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ClosedRange)
			other, ok := args[1].(*value.ClosedRange)
			if !ok {
				return value.False, nil
			}
			equal, err := ClosedRangeEqual(vm, self, other)
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
			self := args[0].(*value.ClosedRange)
			other := args[1]
			contains, err := ClosedRangeContains(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(contains), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
}

// ::Std::ClosedRange::Iterator
func init() {
	// Instance methods
	c := &value.ClosedRangeIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ClosedRangeIterator)
			return ClosedRangeIteratorNext(vm, self)
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

// Checks whether a value is contained in the closed range
func ClosedRangeContains(vm *VM, r *value.ClosedRange, val value.Value) (bool, value.Value) {
	eqVal, err := GreaterThanEqual(vm, val, r.From)
	if err != nil {
		return false, err
	}

	if value.Falsy(eqVal) {
		return false, nil
	}

	eqVal, err = LessThanEqual(vm, val, r.To)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}

// Checks whether two closed ranges are equal
func ClosedRangeEqual(vm *VM, x *value.ClosedRange, y *value.ClosedRange) (bool, value.Value) {
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
func ClosedRangeIteratorNext(vm *VM, i *value.ClosedRangeIterator) (value.Value, value.Value) {
	greater, err := GreaterThan(vm, i.CurrentElement, i.Range.To)
	if err != nil {
		return nil, err
	}

	if value.Truthy(greater) {
		return nil, stopIterationSymbol
	}

	current := i.CurrentElement

	// i.CurrentElement++
	next, err := vm.CallMethod(incrementSymbol, i.CurrentElement)
	if err != nil {
		return nil, err
	}
	i.CurrentElement = next

	return current, nil
}
