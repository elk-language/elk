package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::OpenRange
func initOpenRange() {
	// Instance methods
	c := &value.OpenRangeClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.OpenRange)
			iterator := value.NewOpenRangeIterator(self)
			return value.Ref(iterator), value.Nil
		},
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.OpenRange)
			other, ok := args[1].SafeAsReference().(*value.OpenRange)
			if !ok {
				return value.False, value.Nil
			}
			equal, err := OpenRangeEqual(vm, self, other)
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
			self := args[0].MustReference().(*value.OpenRange)
			other := args[1]
			if !value.IsA(other, self.Start.Class()) && !value.IsA(other, self.End.Class()) {
				return value.False, value.Nil
			}
			contains, err := OpenRangeContains(vm, self, other)
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
			self := args[0].MustReference().(*value.OpenRange)
			other := args[1]
			contains, err := OpenRangeContains(vm, self, other)
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
			self := args[0].MustReference().(*value.OpenRange)
			return self.Start, value.Nil
		},
	)
	Def(
		c,
		"end",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.OpenRange)
			return self.End, value.Nil
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
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.OpenRangeIterator)
			return OpenRangeIteratorNext(vm, self)
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

// Checks whether a value is contained in the open range
func OpenRangeContains(vm *VM, r *value.OpenRange, val value.Value) (bool, value.Value) {
	eqVal, err := GreaterThan(vm, val, r.Start)
	if !err.IsNil() {
		return false, err
	}

	if value.Falsy(eqVal) {
		return false, value.Nil
	}

	eqVal, err = LessThan(vm, val, r.End)
	if !err.IsNil() {
		return false, err
	}

	return value.Truthy(eqVal), value.Nil
}

// Checks whether two open ranges are equal
func OpenRangeEqual(vm *VM, x *value.OpenRange, y *value.OpenRange) (bool, value.Value) {
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
func OpenRangeIteratorNext(vm *VM, i *value.OpenRangeIterator) (value.Value, value.Value) {
	// i.CurrentElement++
	next, err := Increment(vm, i.CurrentElement)
	if !err.IsNil() {
		return value.Nil, err
	}
	i.CurrentElement = next

	greater, err := GreaterThanEqual(vm, i.CurrentElement, i.Range.End)
	if !err.IsNil() {
		return value.Nil, err
	}

	if value.Truthy(greater) {
		return value.Nil, stopIterationSymbol.ToValue()
	}

	current := i.CurrentElement

	return current, value.Nil
}
