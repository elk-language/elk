package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::LeftOpenRange
func init() {
	// Instance methods
	c := &value.LeftOpenRangeClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.LeftOpenRange)
			iterator := value.NewLeftOpenRangeIterator(self)
			return iterator, nil
		},
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.LeftOpenRange)
			other, ok := args[1].(*value.LeftOpenRange)
			if !ok {
				return value.False, nil
			}
			equal, err := LeftOpenRangeEqual(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(equal), nil
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
			self := args[0].(*value.LeftOpenRange)
			other := args[1]
			if !value.IsA(other, self.Start.Class()) && !value.IsA(other, self.End.Class()) {
				return value.False, nil
			}
			contains, err := LeftOpenRangeContains(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(contains), nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.LeftOpenRange)
			other := args[1]
			contains, err := LeftOpenRangeContains(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(contains), nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"is_left_closed",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.False, nil
		},
	)
	Def(
		c,
		"is_left_open",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.True, nil
		},
	)
	Def(
		c,
		"is_right_closed",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.True, nil
		},
	)
	Def(
		c,
		"is_right_open",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.False, nil
		},
	)
	Def(
		c,
		"start",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.LeftOpenRange)
			return self.Start, nil
		},
	)
	Def(
		c,
		"end",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.LeftOpenRange)
			return self.End, nil
		},
	)
}

// ::Std::LeftOpenRange::Iterator
func init() {
	// Instance methods
	c := &value.LeftOpenRangeIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.LeftOpenRangeIterator)
			return LeftOpenRangeIteratorNext(vm, self)
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

// Checks whether a value is contained in the left open range
func LeftOpenRangeContains(vm *VM, r *value.LeftOpenRange, val value.Value) (bool, value.Value) {
	eqVal, err := GreaterThan(vm, val, r.Start)
	if err != nil {
		return false, err
	}

	if value.Falsy(eqVal) {
		return false, nil
	}

	eqVal, err = LessThanEqual(vm, val, r.End)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}

// Checks whether two left open ranges are equal
func LeftOpenRangeEqual(vm *VM, x, y *value.LeftOpenRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.Start, y.Start)
	if err != nil {
		return false, err
	}

	if value.Falsy(eqVal) {
		return false, nil
	}

	eqVal, err = Equal(vm, x.End, y.End)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}

// Get the next element of the range
func LeftOpenRangeIteratorNext(vm *VM, i *value.LeftOpenRangeIterator) (value.Value, value.Value) {
	// i.CurrentElement++
	next, err := Increment(vm, i.CurrentElement)
	if err != nil {
		return nil, err
	}
	i.CurrentElement = next

	greater, err := GreaterThan(vm, i.CurrentElement, i.Range.End)
	if err != nil {
		return nil, err
	}

	if value.Truthy(greater) {
		return nil, stopIterationSymbol
	}

	current := i.CurrentElement

	return current, nil
}
