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
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.BeginlessClosedRange)
			other, ok := args[1].SafeAsReference().(*value.BeginlessClosedRange)
			if !ok {
				return value.False.ToValue(), value.Undefined
			}
			equal, err := BeginlessClosedRangeEqual(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.BoolVal(equal), value.Undefined
		},
		DefWithParameters(1),
	)
	// Special version of `contains` used in pattern matching.
	// Given value has to be an instance of the same class as `end`,
	// otherwise `false` will be returned
	Def(
		c,
		"#contains",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.BeginlessClosedRange)
			other := args[1]
			if !value.IsA(other, self.End.Class()) {
				return value.False.ToValue(), value.Undefined
			}
			contains, err := BeginlessClosedRangeContains(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.BoolVal(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.BeginlessClosedRange)
			other := args[1]
			contains, err := BeginlessClosedRangeContains(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.BoolVal(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"is_left_closed",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return value.False.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"is_left_open",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return value.True.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"is_right_closed",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return value.True.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"is_right_open",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return value.False.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"start",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return value.Undefined, value.Undefined
		},
	)
	Def(
		c,
		"end",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.BeginlessClosedRange)
			return self.End, value.Undefined
		},
	)
}

// Checks whether a value is contained in the beginless closed range
func BeginlessClosedRangeContains(vm *Thread, r *value.BeginlessClosedRange, val value.Value) (bool, value.Value) {
	eqVal, err := LessThanEqual(vm, val, r.End)
	if !err.IsUndefined() {
		return false, err
	}

	return value.Truthy(eqVal), value.Undefined
}

// Checks whether two beginless closed ranges are equal
func BeginlessClosedRangeEqual(vm *Thread, x, y *value.BeginlessClosedRange) (bool, value.Value) {
	eqVal, err := Equal(vm, x.End, y.End)
	if !err.IsUndefined() {
		return false, err
	}

	return value.Truthy(eqVal), value.Undefined
}
