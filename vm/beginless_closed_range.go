package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::BeginlessClosedRange
func initBeginlessClosedRange() {
	// Instance methods
	c := &value.BeginlessClosedRangeClass.MethodContainer
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BeginlessClosedRange)
			other, ok := args[1].(*value.BeginlessClosedRange)
			if !ok {
				return value.False, nil
			}
			equal, err := BeginlessClosedRangeEqual(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(equal), nil
		},
		DefWithParameters(1),
	)
	// Special version of `contains` used in pattern matching.
	// Given value has to be an instance of the same class as `end`,
	// otherwise `false` will be returned
	Def(
		c,
		"#contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BeginlessClosedRange)
			other := args[1]
			if !value.IsA(other, self.End.Class()) {
				return value.False, nil
			}
			contains, err := BeginlessClosedRangeContains(vm, self, other)
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
			self := args[0].(*value.BeginlessClosedRange)
			other := args[1]
			contains, err := BeginlessClosedRangeContains(vm, self, other)
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
			return value.Nil, nil
		},
	)
	Def(
		c,
		"end",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BeginlessClosedRange)
			return self.End, nil
		},
	)
}

// Checks whether a value is contained in the beginless closed range
func BeginlessClosedRangeContains(vm *VM, r *value.BeginlessClosedRange, val value.Value) (bool, value.Value) {
	eqVal, err := LessThanEqual(vm, val, r.End)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}

// Checks whether two beginless closed ranges are equal
func BeginlessClosedRangeEqual(vm *VM, x, y *value.BeginlessClosedRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.End, y.End)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}
