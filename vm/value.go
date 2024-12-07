package vm

import (
	"encoding/binary"
	"reflect"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Std::Value
func initValue() {
	// Instance methods
	c := &value.ValueClass.MethodContainer
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return value.Ref(value.String(self.Inspect())), value.Nil
		},
	)
	Def(
		c,
		"class",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return value.Ref(self.Class()), value.Nil
		},
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			return value.ToElkBool(self == other), value.Nil
		},
		DefWithParameters(1),
	)
	Alias(c, "===", "==")
	Def(
		c,
		"copy",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return self.Copy(), value.Nil
		},
	)
	Def(
		c,
		"hash",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			result, err := value.Hash(self)
			if err == value.Ref(value.NotBuiltinError) {
				return ObjectHash(self).ToValue(), value.Nil
			}
			if err.IsNil() {
				return value.Ref(result), value.Nil
			}
			return value.Nil, err
		},
	)

}

func ObjectHash(val value.Value) value.UInt64 {
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Ptr {
		if !v.CanAddr() {
			return value.UInt64(0)
		}

		v = v.Addr()
	}
	ptr := v.Pointer()
	b := make([]byte, 8)
	d := xxhash.New()
	binary.LittleEndian.PutUint64(b, uint64(ptr))
	d.Write(b)
	return value.UInt64(d.Sum64())
}

// Calculate the hash for the given value
func Hash(vm *VM, key value.Value) (value.UInt64, value.Value) {
	result, err := value.Hash(key)

	if err == value.Ref(value.NotBuiltinError) {
		if vm == nil {
			return 0, value.Nil
		}
		class := key.DirectClass()
		method := class.LookupMethod(symbol.L_hash)
		if method == nil {
			return ObjectHash(key), value.Nil
		}
		dynamicResult, dynamicErr := vm.CallMethod(method, key)
		if !dynamicErr.IsNil() {
			return 0, dynamicErr
		}
		if dynamicResult.IsUInt64() {
			return dynamicResult.AsUInt64(), value.Nil
		}
		return 0, value.Ref(value.NewCoerceError(
			value.UInt64Class,
			dynamicResult.Class(),
		))
	} else if !err.IsNil() {
		return 0, err
	}

	return result, value.Nil
}

// Check whether two values are equal
func Equal(vm *VM, left, right value.Value) (value.Value, value.Value) {
	result := value.Equal(left, right)

	if !result.IsNil() {
		return result, value.Nil
	}
	if vm == nil {
		return value.Nil, value.Nil
	}

	result, err := vm.CallMethodByName(symbol.OpEqual, left, right)
	if !err.IsNil() {
		return value.Nil, err
	}
	return result, value.Nil
}

// Check whether two values are equal (lax)
func LaxEqual(vm *VM, left, right value.Value) (value.Value, value.Value) {
	result := value.LaxEqual(left, right)

	if !result.IsNil() {
		return result, value.Nil
	}
	if vm == nil {
		return value.Nil, value.Nil
	}

	result, err := vm.CallMethodByName(symbol.OpLaxEqual, left, right)
	if !err.IsNil() {
		return value.Nil, err
	}
	return result, value.Nil
}

// Check whether the left value is greater than the right
func GreaterThan(vm *VM, left, right value.Value) (value.Value, value.Value) {
	result, err := value.GreaterThan(left, right)

	if !err.IsNil() {
		return value.Nil, err
	}
	if !result.IsNil() {
		return result, value.Nil
	}
	if vm == nil {
		return value.Nil, value.Nil
	}

	result, err2 := vm.CallMethodByName(symbol.OpGreaterThan, left, right)
	if !err2.IsNil() {
		return value.Nil, err2
	}
	return result, value.Nil
}

// Check whether the left value is greater than or equal to the right
func GreaterThanEqual(vm *VM, left, right value.Value) (value.Value, value.Value) {
	result, err := value.GreaterThanEqual(left, right)

	if !err.IsNil() {
		return value.Nil, err
	}
	if !result.IsNil() {
		return result, value.Nil
	}
	if vm == nil {
		return value.Nil, value.Nil
	}

	result, err2 := vm.CallMethodByName(symbol.OpGreaterThanEqual, left, right)
	if !err2.IsNil() {
		return value.Nil, err2
	}
	return result, value.Nil
}

// Check whether the left value is less than the right
func LessThan(vm *VM, left, right value.Value) (value.Value, value.Value) {
	result, err := value.LessThan(left, right)

	if !err.IsNil() {
		return value.Nil, err
	}
	if !result.IsNil() {
		return result, value.Nil
	}
	if vm == nil {
		return value.Nil, value.Nil
	}

	result, err2 := vm.CallMethodByName(symbol.OpLessThan, left, right)
	if !err2.IsNil() {
		return value.Nil, err2
	}
	return result, value.Nil
}

// Check whether the left value is less than or equal to the right
func LessThanEqual(vm *VM, left, right value.Value) (value.Value, value.Value) {
	result, err := value.LessThanEqual(left, right)

	if !err.IsNil() {
		return value.Nil, err
	}
	if !result.IsNil() {
		return result, value.Nil
	}
	if vm == nil {
		return value.Nil, value.Nil
	}

	result, err2 := vm.CallMethodByName(symbol.OpLessThanEqual, left, right)
	if !err2.IsNil() {
		return value.Nil, err2
	}
	return result, value.Nil
}

// Increment the given value
func Increment(vm *VM, val value.Value) (value.Value, value.Value) {
	result := value.Increment(val)

	if !result.IsNil() {
		return result, value.Nil
	}
	if vm == nil {
		return value.Nil, value.Nil
	}

	result, err := vm.CallMethodByName(symbol.OpIncrement, val)
	if !err.IsNil() {
		return value.Nil, err
	}
	return result, value.Nil
}

// Decrement the given value
func Decrement(vm *VM, val value.Value) (value.Value, value.Value) {
	result := value.Decrement(val)

	if !result.IsNil() {
		return result, value.Nil
	}
	if vm == nil {
		return value.Nil, value.Nil
	}

	result, err := vm.CallMethodByName(symbol.OpDecrement, val)
	if !err.IsNil() {
		return value.Nil, err
	}
	return result, value.Nil
}
