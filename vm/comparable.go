package vm

import (
	"github.com/elk-language/elk/value"
)

var compareMethodSymbol = value.ToSymbol("<=>")

func init() {
	DefineMethodWithOptions(
		value.ComparableMixin.Methods,
		">",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			result, err := vm.CallMethod(self, compareMethodSymbol, []value.Value{other})
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
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.ComparableMixin.Methods,
		">=",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			result, err := vm.CallMethod(self, compareMethodSymbol, []value.Value{other})
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
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.ComparableMixin.Methods,
		"<",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			result, err := vm.CallMethod(self, compareMethodSymbol, []value.Value{other})
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
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.ComparableMixin.Methods,
		"<=",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			result, err := vm.CallMethod(self, compareMethodSymbol, []value.Value{other})
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
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.ComparableMixin.Methods,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			result, err := vm.CallMethod(self, compareMethodSymbol, []value.Value{other})
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
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.ComparableMixin.Methods,
		"===",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.Class() != other.Class() {
				return value.False, nil
			}
			result, err := vm.CallMethod(self, compareMethodSymbol, []value.Value{other})
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
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)

}
