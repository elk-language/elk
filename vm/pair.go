package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::Pair
func init() {
	// Instance methods
	c := &value.PairClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Pair)
			self.Key = args[1]
			self.Value = args[2]
			return self, nil
		},
		DefWithParameters("key", "value"),
	)
	Def(
		c,
		"key",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Pair)
			return self.Key, nil
		},
	)
	Def(
		c,
		"value",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Pair)
			return self.Value, nil
		},
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Pair)
			other, ok := args[1].(*value.Pair)
			if !ok {
				return value.False, nil
			}
			equal, err := PairEqual(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(equal), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
}

// Checks whether two pairs are equal
func PairEqual(vm *VM, x *value.Pair, y *value.Pair) (bool, value.Value) {
	eqVal, err := Equal(vm, x.Key, y.Key)
	if err != nil {
		return false, err
	}

	if value.Falsy(eqVal) {
		return false, nil
	}

	eqVal, err = Equal(vm, x.Value, y.Value)
	if err != nil {
		return false, err
	}

	return value.Truthy(eqVal), nil
}
