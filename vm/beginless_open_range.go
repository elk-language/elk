package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::BeginlessOpenRange
func init() {
	// Instance methods
	c := &value.BeginlessOpenRangeClass.MethodContainer
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BeginlessOpenRange)
			other, ok := args[1].(*value.BeginlessOpenRange)
			if !ok {
				return value.False, nil
			}
			equal, err := BeginlessOpenRangeEqual(vm, self, other)
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
			self := args[0].(*value.BeginlessOpenRange)
			other := args[1]
			contains, err := BeginlessOpenRangeContains(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(contains), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
}

// Checks whether a value is contained in the open range
func BeginlessOpenRangeContains(vm *VM, r *value.BeginlessOpenRange, val value.Value) (bool, value.Value) {
	eqVal, err := LessThan(vm, val, r.To)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}

// Checks whether two open ranges are equal
func BeginlessOpenRangeEqual(vm *VM, x, y *value.BeginlessOpenRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.To, y.To)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}
