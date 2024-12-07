package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::Pair
func initPair() {
	// Instance methods
	c := &value.PairClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Pair)
			self.Key = args[1]
			self.Value = args[2]
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(2),
	)
	Def(
		c,
		"key",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Pair)
			return self.Key, value.Undefined
		},
	)
	Def(
		c,
		"value",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Pair)
			return self.Value, value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *VM, _ []value.Value) (value.Value, value.Value) {
			return value.SmallInt(2).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Pair)
			other := args[1]
			return self.Subscript(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"[]=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Pair)
			key := args[1]
			val := args[2]
			err := self.SubscriptSet(key, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return val, value.Undefined
		},
		DefWithParameters(2),
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Pair)
			other, ok := args[1].SafeAsReference().(*value.Pair)
			if !ok {
				return value.False, value.Undefined
			}
			equal, err := PairEqual(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.ToElkBool(equal), value.Undefined
		},
		DefWithParameters(1),
	)
}

// Checks whether two pairs are equal
func PairEqual(vm *VM, x *value.Pair, y *value.Pair) (bool, value.Value) {
	eqVal, err := Equal(vm, x.Key, y.Key)
	if !err.IsUndefined() {
		return false, err
	}

	if value.Falsy(eqVal) {
		return false, value.Undefined
	}

	eqVal, err = Equal(vm, x.Value, y.Value)
	if !err.IsUndefined() {
		return false, err
	}

	return value.Truthy(eqVal), value.Undefined
}
