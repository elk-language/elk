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
				result, err := vm.CallMethod(toStringSymbol, val)
				if err != nil {
					return nil, err
				}
				fmt.Fprint(vm.Stdout, result)
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
				result, err := vm.CallMethod(toStringSymbol, val)
				if err != nil {
					return nil, err
				}
				fmt.Fprintln(vm.Stdout, result)
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
	Def(
		c,
		"copy",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return self.Copy(), nil
		},
	)

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

// Check whether two values are equal
func Equal(vm *VM, left, right value.Value) (value.Value, value.Value) {
	result := value.Equal(left, right)

	if result != nil {
		return result, nil
	}
	if vm == nil {
		return nil, value.Nil
	}

	result, err := vm.CallMethod(equalSymbol, left, right)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Check whether the left value is greater than the right
func GreaterThan(vm *VM, left, right value.Value) (value.Value, value.Value) {
	result, err := value.GreaterThan(left, right)

	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}
	if vm == nil {
		return nil, value.Nil
	}

	result, err2 := vm.CallMethod(greaterThanSymbol, left, right)
	if err2 != nil {
		return nil, err2
	}
	return result, nil
}

// Check whether the left value is greater than or equal to the right
func GreaterThanEqual(vm *VM, left, right value.Value) (value.Value, value.Value) {
	result, err := value.GreaterThanEqual(left, right)

	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}
	if vm == nil {
		return nil, value.Nil
	}

	result, err2 := vm.CallMethod(greaterThanEqualSymbol, left, right)
	if err2 != nil {
		return nil, err2
	}
	return result, nil
}

// Check whether the left value is less than the right
func LessThan(vm *VM, left, right value.Value) (value.Value, value.Value) {
	result, err := value.LessThan(left, right)

	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}
	if vm == nil {
		return nil, value.Nil
	}

	result, err2 := vm.CallMethod(lessThanSymbol, left, right)
	if err2 != nil {
		return nil, err2
	}
	return result, nil
}

// Check whether the left value is less than or equal to the right
func LessThanEqual(vm *VM, left, right value.Value) (value.Value, value.Value) {
	result, err := value.LessThanEqual(left, right)

	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}
	if vm == nil {
		return nil, value.Nil
	}

	result, err2 := vm.CallMethod(lessThanEqualSymbol, left, right)
	if err2 != nil {
		return nil, err2
	}
	return result, nil
}

// Increment the given value
func Increment(vm *VM, val value.Value) (value.Value, value.Value) {
	result := value.Increment(val)

	if result != nil {
		return result, nil
	}
	if vm == nil {
		return nil, value.Nil
	}

	result, err := vm.CallMethod(incrementSymbol, val)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Decrement the given value
func Decrement(vm *VM, val value.Value) (value.Value, value.Value) {
	result := value.Decrement(val)

	if result != nil {
		return result, nil
	}
	if vm == nil {
		return nil, value.Nil
	}

	result, err := vm.CallMethod(decrementSymbol, val)
	if err != nil {
		return nil, err
	}
	return result, nil
}
