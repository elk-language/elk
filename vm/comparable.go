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
			if err != nil {
				return nil, err
			}

			switch r := result.(type) {
			case value.SmallInt:
				return value.ToElkBool(r > 0), nil
			default:
				return value.False, nil
			}
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
			if err != nil {
				return nil, err
			}

			switch r := result.(type) {
			case value.SmallInt:
				return value.ToElkBool(r >= 0), nil
			default:
				return value.False, nil
			}
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
			if err != nil {
				return nil, err
			}

			switch r := result.(type) {
			case value.SmallInt:
				return value.ToElkBool(r < 0), nil
			default:
				return value.False, nil
			}
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
			if err != nil {
				return nil, err
			}

			switch r := result.(type) {
			case value.SmallInt:
				return value.ToElkBool(r <= 0), nil
			default:
				return value.False, nil
			}
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
			if err != nil {
				return nil, err
			}

			switch r := result.(type) {
			case value.SmallInt:
				return value.ToElkBool(r == 0), nil
			default:
				return value.False, nil
			}
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
				return value.False, nil
			}
			result, err := vm.CallMethodByName(symbol.OpSpaceship, self, other)
			if err != nil {
				return nil, err
			}

			switch r := result.(type) {
			case value.SmallInt:
				return value.ToElkBool(r == 0), nil
			default:
				return value.False, nil
			}
		},
		DefWithParameters(1),
	)

}
