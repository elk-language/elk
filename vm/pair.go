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
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Pair)
			err := self.SetKey(args[1])
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			err = self.SetValue(args[2])
			if err.IsNotUndefined() {
				return value.Undefined, err
			}
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(2),
	)
	Def(
		c,
		"key",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Pair)
			return self.Key(), value.Undefined
		},
	)
	Def(
		c,
		"value",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Pair)
			return self.Value(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, _ []value.Value) (value.Value, value.Value) {
			return value.SmallInt(2).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Pair)
			other := args[1]
			return self.Subscript(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"[]=",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Pair)
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
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Pair)
			other, ok := args[1].AsReference().(value.Pair)
			if !ok {
				return value.False.ToValue(), value.Undefined
			}
			equal, err := PairEqual(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.BoolVal(equal), value.Undefined
		},
		DefWithParameters(1),
	)
}

// Checks whether two pairs are equal
func PairEqual(vm *Thread, x value.Pair, y value.Pair) (bool, value.Value) {
	eqVal, err := Equal(vm, x.Key(), y.Key())
	if !err.IsUndefined() {
		return false, err
	}

	if value.Falsy(eqVal) {
		return false, value.Undefined
	}

	eqVal, err = Equal(vm, x.Value(), y.Value())
	if !err.IsUndefined() {
		return false, err
	}

	return value.Truthy(eqVal), value.Undefined
}
