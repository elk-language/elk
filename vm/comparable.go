package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Std::Comparable
func initComparable() {
	// Instance methods
	c := &value.ComparableMixin.MethodContainer
	Def(
		c,
		">",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			result, err := vm.CallMethodByName(symbol.OpSpaceship, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			if result.IsSmallInt() {
				r := result.AsSmallInt()
				return value.ToElkBool(r > 0), value.Undefined
			}

			return value.False, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			result, err := vm.CallMethodByName(symbol.OpSpaceship, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			if result.IsSmallInt() {
				r := result.AsSmallInt()
				return value.ToElkBool(r >= 0), value.Undefined
			}

			return value.False, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			result, err := vm.CallMethodByName(symbol.OpSpaceship, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			if result.IsSmallInt() {
				r := result.AsSmallInt()
				return value.ToElkBool(r < 0), value.Undefined
			}

			return value.False, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			result, err := vm.CallMethodByName(symbol.OpSpaceship, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			if result.IsSmallInt() {
				r := result.AsSmallInt()
				return value.ToElkBool(r <= 0), value.Undefined
			}

			return value.False, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			result, err := vm.CallMethodByName(symbol.OpSpaceship, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			if result.IsSmallInt() {
				r := result.AsSmallInt()
				return value.ToElkBool(r == 0), value.Undefined
			}

			return value.False, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"===",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.Class() != other.Class() {
				return value.False, value.Undefined
			}
			result, err := vm.CallMethodByName(symbol.OpSpaceship, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			if result.IsSmallInt() {
				r := result.AsSmallInt()
				return value.ToElkBool(r == 0), value.Undefined
			}

			return value.False, value.Undefined
		},
		DefWithParameters(1),
	)

}
