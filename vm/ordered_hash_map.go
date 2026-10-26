package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type OrderedHashMap interface {
	OrderedHashRecord
	IterOrderedMap() value.NativeResettableIterator
	RemoveKey(thread *Thread, key value.Value) (bool, value.Value)
	SetVal(thread *Thread, key, val value.Value) value.Value
	CloneOrderedHashMap(thread *Thread, capacity int) (OrderedHashMap, value.Value)
	NewOrderedHashMap(capacity int) OrderedHashMap
}

func orderedHashMapConstructor(class *value.Class) value.Value {
	return NewOrderedHashMapOfValue(5).ToValue()
}

// ::Std::OrderedHashMap
func initOrderedHashMap() {
	value.OrderedHashMapClass.ConstructorFunc = orderedHashMapConstructor

	// Instance methods
	c := value.OrderedHashMapClass
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			iterator := self.IterOrderedMap()
			return iterator.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			key := args[1]
			return self.GetValNil(vm, key)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"[]=",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			key := args[1]
			val := args[2]
			err := self.SetVal(vm, key, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return val, value.Undefined
		},
		DefWithParameters(2),
	)
	Def(
		c,
		"+",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			other := args[1]
			return self.ConcatVal(vm, other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"remove_key",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			key := args[1]
			ok, err := self.RemoveKey(vm, key)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			return value.Bool(ok).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"try_first",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			length := self.Length()
			if length == 0 {
				return value.Nil, value.Undefined
			}

			return self.AtInterface(0).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"first",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			length := self.Length()
			if length == 0 {
				return value.Undefined, value.Ref(
					value.Errorf(
						value.IterableNotFoundErrorClass,
						"cannot get first element of `%s`",
						self.Inspect(),
					),
				)
			}

			return self.AtInterface(0).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"try_last",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			length := self.Length()
			if length == 0 {
				return value.Nil, value.Undefined
			}

			return self.AtInterface(length - 1).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"last",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			length := self.Length()
			if length == 0 {
				return value.Undefined, value.Ref(
					value.Errorf(
						value.IterableNotFoundErrorClass,
						"cannot get last element of `%s`",
						self.Inspect(),
					),
				)
			}

			return self.AtInterface(length - 1).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			otherVal := args[1]
			switch other := otherVal.SafeAsReference().(type) {
			case value.Pair:
				contains, err := self.Contains(vm, other)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				return value.BoolVal(contains), value.Undefined
			default:
				return value.Undefined, value.Ref(value.NewCoerceError(value.PairClass, otherVal.Class()))
			}
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains_key",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			contains, err := self.ContainsKey(vm, args[1])
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			return value.BoolVal(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains_value",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			contains, err := self.ContainsValue(vm, args[1])
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			return value.BoolVal(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			other := args[1]
			equal, err := self.Equal(vm, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.BoolVal(equal), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(OrderedHashMap)
			equal, err := self.LaxEqual(vm, args[1])
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.BoolVal(equal), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			callable := args[1]
			newMap := NewOrderedHashMapOfValue(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for pair := range self.All() {
					result, err := vm.CallClosure(function, pair.ToValue())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					r, ok := result.SafeAsReference().(value.Pair)
					if !ok {
						return value.Undefined, value.Ref(value.NewArgumentTypeError("pair", result.Class().Name(), value.PairClass.Name()))
					}
					err = newMap.SetVal(vm, r.Key(), r.Value())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newMap), value.Undefined
			}

			// callable is another value
			for pair := range self.All() {
				result, err := vm.CallMethodByName(value.S(symbol.L_call), callable, pair.ToValue())
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				r, ok := result.SafeAsReference().(value.Pair)
				if !ok {
					return value.Undefined, value.Ref(value.NewArgumentTypeError("pair", result.Class().Name(), value.PairClass.Name()))
				}
				err = newMap.SetVal(vm, r.Key(), r.Value())
				if !err.IsUndefined() {
					return value.Undefined, err
				}
			}
			return value.Ref(newMap), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map_values",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			callable := args[1]
			newMap := NewOrderedHashMapOfValue(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for pair := range self.All() {
					if pair.Key().IsUndefined() {
						continue
					}
					result, err := vm.CallClosure(function, pair.Value())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					newMap.SetVal(vm, pair.Key(), result)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newMap), value.Undefined
			}

			// callable is another value
			for pair := range self.All() {
				if pair.Key().IsUndefined() {
					continue
				}
				result, err := vm.CallMethodByName(value.S(symbol.L_call), callable, pair.Value())
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				newMap.SetVal(vm, pair.Key(), result)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
			}
			return value.Ref(newMap), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map_values_mut",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashMap)
			callable := args[1]

			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for pair := range self.All() {
					if pair.Key().IsUndefined() {
						continue
					}
					result, err := vm.CallClosure(function, pair.Value())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					self.SetVal(vm, pair.Key(), result)
				}
				return self.ToValue(), value.Undefined
			}

			// callable is another value
			for pair := range self.All() {
				if pair.Key().IsUndefined() {
					continue
				}
				result, err := vm.CallMethodByName(value.S(symbol.L_call), callable, pair.Value())
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				self.SetVal(vm, pair.Key(), result)
			}
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
}

// ::Std::OrderedHashMap::Iterator
func initOrderedHashMapIterator() {
	// Instance methods
	c := value.OrderedHashMapIteratorClass
	Def(
		c,
		"next",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.NativeResettableIterator)
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
			self := args[0].AsReference().(value.NativeResettableIterator)
			self.Reset()
			return args[0], value.Undefined
		},
	)

}
