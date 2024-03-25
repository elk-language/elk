package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::OpenRange
func init() {
	// Instance methods
	c := &value.OpenRangeClass.MethodContainer
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.OpenRange)
			iterator := value.NewOpenRangeIterator(self)
			return iterator, nil
		},
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.OpenRange)
			other, ok := args[1].(*value.OpenRange)
			if !ok {
				return value.False, nil
			}
			equal, err := OpenRangeEqual(vm, self, other)
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
			self := args[0].(*value.OpenRange)
			other := args[1]
			contains, err := OpenRangeContains(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(contains), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
}

// ::Std::OpenRange::Iterator
func init() {
	// Instance methods
	c := &value.OpenRangeIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.OpenRangeIterator)
			return OpenRangeIteratorNext(vm, self)
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

// Checks whether a value is contained in the open range
func OpenRangeContains(vm *VM, r *value.OpenRange, val value.Value) (bool, value.Value) {
	eqVal, err := GreaterThan(vm, val, r.From)
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

// Checks whether two open ranges are equal
func OpenRangeEqual(vm *VM, x *value.OpenRange, y *value.OpenRange) (bool, value.Value) {
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
func OpenRangeIteratorNext(vm *VM, i *value.OpenRangeIterator) (value.Value, value.Value) {
	// i.CurrentElement++
	next, err := Increment(vm, i.CurrentElement)
	if err != nil {
		return nil, err
	}
	i.CurrentElement = next

	greater, err := GreaterThanEqual(vm, i.CurrentElement, i.Range.To)
	if err != nil {
		return nil, err
	}

	if value.Truthy(greater) {
		return nil, stopIterationSymbol
	}

	current := i.CurrentElement

	return current, nil
}
