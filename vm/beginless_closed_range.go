package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::BeginlessClosedRange
func init() {
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
		DefWithParameters("other"),
		DefWithSealed(),
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
		DefWithParameters("other"),
		DefWithSealed(),
	)
}

// Checks whether a value is contained in the beginless closed range
func BeginlessClosedRangeContains(vm *VM, r *value.BeginlessClosedRange, val value.Value) (bool, value.Value) {
	eqVal, err := LessThanEqual(vm, val, r.To)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}

// Checks whether two beginless closed ranges are equal
func BeginlessClosedRangeEqual(vm *VM, x, y *value.BeginlessClosedRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.To, y.To)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}
