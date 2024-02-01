package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

func init() {
	// Instance methods
	c := &value.ValueClass.MethodContainer
	Def(
		c,
		"print",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			values := args[1].(*value.ArrayList)
			for _, val := range *values {
				fmt.Fprint(vm.Stdout, val)
			}

			return value.Nil, nil
		},
		DefWithParameters("values"),
		DefWithPositionalRestParameter(),
	)
	Def(
		c,
		"println",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			values := args[1].(*value.ArrayList)
			for _, val := range *values {
				fmt.Fprintln(vm.Stdout, val)
			}

			return value.Nil, nil
		},
		DefWithParameters("values"),
		DefWithPositionalRestParameter(),
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return value.String(self.Inspect()), nil
		},
	)
	Def(
		c,
		"inspect_stack",
		func(v *VM, args []value.Value) (value.Value, value.Value) {
			v.InspectStack()
			return value.Nil, nil
		},
	)
	Def(
		c,
		"class",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return self.Class(), nil
		},
	)

	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			return value.ToElkBool(self == other), nil
		},
		DefWithParameters("other"),
	)
	Alias(c, "===", "==")

}

var hashSymbol = value.ToSymbol("hash")

// Calculate the hash for the given value
func Hash(vm *VM, key value.Value) (value.UInt64, value.Value) {
	result, err := value.Hash(key)

	if err == value.NotBuiltinError {
		if vm == nil {
			return 0, value.Nil
		}
		dynamicResult, dynamicErr := vm.CallMethod(hashSymbol, key)
		if dynamicErr != nil {
			return 0, dynamicErr
		}
		uintResult, ok := dynamicResult.(value.UInt64)
		if !ok {
			return 0, value.NewCoerceError(
				value.UInt64Class,
				dynamicResult.Class(),
			)
		}
		return uintResult, nil
	} else if err != nil {
		return 0, err
	}

	return result, nil
}

// Calculate the hash for the given value
func StrictEqual(vm *VM, left, right value.Value) (value.Value, value.Value) {
	result := value.StrictEqual(left, right)

	if result != nil {
		return result, nil
	}
	if vm == nil {
		return nil, value.Nil
	}

	result, err := vm.CallMethod(strictEqualSymbol, left, right)
	if err != nil {
		return nil, err
	}
	return result, nil
}
