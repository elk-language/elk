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
			if !err.IsNil() {
				return value.Nil, err
			}

			if result.IsSmallInt() {
				r := result.AsSmallInt()
				return value.ToElkBool(r > 0), value.Nil
			}

			return value.False, value.Nil
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
			if !err.IsNil() {
				return value.Nil, err
			}

			if result.IsSmallInt() {
				r := result.AsSmallInt()
				return value.ToElkBool(r >= 0), value.Nil
			}

			return value.False, value.Nil
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
			if !err.IsNil() {
				return value.Nil, err
			}

			if result.IsSmallInt() {
				r := result.AsSmallInt()
				return value.ToElkBool(r < 0), value.Nil
			}

			return value.False, value.Nil
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
			if !err.IsNil() {
				return value.Nil, err
			}

			if result.IsSmallInt() {
				r := result.AsSmallInt()
				return value.ToElkBool(r <= 0), value.Nil
			}

			return value.False, value.Nil
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
			if !err.IsNil() {
				return value.Nil, err
			}

			if result.IsSmallInt() {
				r := result.AsSmallInt()
				return value.ToElkBool(r == 0), value.Nil
			}

			return value.False, value.Nil
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
				return value.False, value.Nil
			}
			result, err := vm.CallMethodByName(symbol.OpSpaceship, self, other)
			if !err.IsNil() {
				return value.Nil, err
			}

			if result.IsSmallInt() {
				r := result.AsSmallInt()
				return value.ToElkBool(r == 0), value.Nil
			}

			return value.False, value.Nil
		},
		DefWithParameters(1),
	)

}
