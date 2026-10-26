package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type HashMap interface {
	HashRecord
	IterMap() value.NativeResettableIterator
	SetVal(thread *Thread, key, val value.Value) value.Value
	RemoveKey(thread *Thread, key value.Value) (bool, value.Value)
	CloneHashMap(thread *Thread, capacity int) (HashMap, value.Value)
	NewHashMap(capacity int) HashMap
}

func hashMapConstructor(class *value.Class) value.Value {
	return NewHashMapOfValue(5).ToValue()
}

// ::Std::HashMap
func initHashMap() {
	value.HashMapClass.ConstructorFunc = hashMapConstructor

	// Instance methods
	c := value.HashMapClass
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashMap)
			iterator := self.IterMap()
			return iterator.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashMap)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashMap)
			key := args[1]
			return self.GetValNil(vm, key)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"[]=",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashMap)
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
			self := args[0].AsReference().(HashMap)
			other := args[1]
			return self.MergeVal(vm, other)
		},
		DefWithParameters(1),
	)
	Alias(c, "merge", "+")

	Def(
		c,
		"remove_key",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashMap)
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
		"contains",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashMap)
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
			self := args[0].AsReference().(HashMap)
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
			self := args[0].AsReference().(HashMap)
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
			self := args[0].AsReference().(HashMap)
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
			self := args[0].MustReference().(HashMap)
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
			self := args[0].AsReference().(HashMap)
			callable := args[1]
			newMap := NewHashMapOfValue(self.Length())

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
					err = HashMapOfValueSet(vm, newMap, r.Key(), r.Value())
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
				err = HashMapOfValueSet(vm, newMap, r.Key(), r.Value())
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
			self := args[0].AsReference().(HashMap)
			callable := args[1]
			newMap := NewHashMapOfValue(self.Length())

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
					err = HashMapOfValueSet(vm, newMap, pair.Key(), result)
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
				err = HashMapOfValueSet(vm, newMap, pair.Key(), result)
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
			self := args[0].AsReference().(HashMap)
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

// ::Std::HashMap::Iterator
func initHashMapIterator() {
	// Instance methods
	c := value.HashMapIteratorClass
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
