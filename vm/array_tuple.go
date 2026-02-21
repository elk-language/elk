package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type mutableArrayTuple interface {
	value.ArrayTuple
	SubscriptSet(key, val value.Value) value.Value
}

// ::Std::ArrayTuple
func initArrayTuple() {
	// Instance methods
	c := &value.ArrayTupleClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			iterator := self.IterTuple()
			return iterator.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			other := args[1]
			return self.Subscript(other)
		},
		DefWithParameters(1),
	)
	Alias(c, "at", "[]")
	Def(
		c,
		"+",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			other := args[1]
			return self.ConcatVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			other := args[1]
			return self.RepeatVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"immutable_box_of",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			other := args[1]
			return self.ImmutableBoxOfVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			switch other := args[1].SafeAsReference().(type) {
			case value.ArrayList:
				equal, err := ArrayTupleEqual(vm, self, other)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.BoolVal(equal), value.Undefined
			case value.ArrayTuple:
				equal, err := ArrayTupleEqual(vm, self, other)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.BoolVal(equal), value.Undefined
			default:
				return value.False.ToValue(), value.Undefined
			}
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			switch other := args[1].SafeAsReference().(type) {
			case value.ArrayTuple:
				equal, err := ArrayTupleEqual(vm, self, other)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.BoolVal(equal), value.Undefined
			default:
				return value.False.ToValue(), value.Undefined
			}
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			callable := args[1]
			newTuple := value.NewArrayTupleOfValueWithLength(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(*Closure); ok {
				for i := range self.Length() {
					element := self.AtVal(i)
					result, err := vm.CallClosure(function, element)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					newTuple.SetAt(i, result)
				}
				return value.Ref(newTuple), value.Undefined
			}

			// callable is another value
			for i := range self.Length() {
				element := self.AtVal(i)
				result, err := vm.CallMethodByName(callSymbol, callable, element)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				newTuple.SetAt(i, result)
			}
			return value.Ref(newTuple), value.Undefined
		},
		DefWithParameters(1),
	)

}

// ::Std::ArrayTupleIterator
func initArrayTupleIterator() {
	// Instance methods
	c := &value.ArrayTupleIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTupleIterator)
			return self.NextValue()
		},
	)
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"reset",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTupleIterator)
			self.Reset()
			return args[0], value.Undefined
		},
	)

}

func ArrayTupleEqual(vm *Thread, x, y value.ArrayTuple) (bool, value.Value) {
	xLen := x.Length()
	if xLen != y.Length() {
		return false, value.Undefined
	}

	for i := 0; i < xLen; i++ {
		equal, err := vm.CallMethodByName(symbol.OpEqual, x.AtVal(i), y.AtVal(i))
		if !err.IsUndefined() {
			return false, err
		}
		if value.Falsy(equal) {
			return false, value.Undefined
		}
	}
	return true, value.Undefined
}
