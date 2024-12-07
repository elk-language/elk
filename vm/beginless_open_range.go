package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::BeginlessOpenRange
func initBeginlessOpenRange() {
	// Instance methods
	c := &value.BeginlessOpenRangeClass.MethodContainer
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.BeginlessOpenRange)
			other, ok := args[1].SafeAsReference().(*value.BeginlessOpenRange)
			if !ok {
				return value.False, value.Nil
			}
			equal, err := BeginlessOpenRangeEqual(vm, self, other)
			if !err.IsNil() {
				return value.Nil, err
			}
			return value.ToElkBool(equal), value.Nil
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
			self := args[0].MustReference().(*value.BeginlessOpenRange)
			other := args[1]
			if !value.IsA(other, self.End.Class()) {
				return value.False, value.Nil
			}
			contains, err := BeginlessOpenRangeContains(vm, self, other)
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
			self := args[0].MustReference().(*value.BeginlessOpenRange)
			other := args[1]
			contains, err := BeginlessOpenRangeContains(vm, self, other)
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
			return value.Nil, value.Nil
		},
	)
	Def(
		c,
		"end",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.BeginlessOpenRange)
			return self.End, value.Nil
		},
	)
}

// Checks whether a value is contained in the open range
func BeginlessOpenRangeContains(vm *VM, r *value.BeginlessOpenRange, val value.Value) (bool, value.Value) {
	eqVal, err := LessThan(vm, val, r.End)
	if !err.IsNil() {
		return false, err
	}

	return value.Truthy(eqVal), value.Nil
}

// Checks whether two open ranges are equal
func BeginlessOpenRangeEqual(vm *VM, x, y *value.BeginlessOpenRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.End, y.End)
	if !err.IsNil() {
		return false, err
	}

	return value.Truthy(eqVal), value.Nil
}
