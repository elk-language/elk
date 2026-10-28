package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type SyncHashMap interface {
	Record
	HashMap() HashMap
	Lock()
	Unlock()
	IterSyncMap() value.NativeResettableIterator
	RemoveKey(thread *Thread, key value.Value) (bool, value.Value)
	SetVal(thread *Thread, key, val value.Value) value.Value
	CloneSyncHashMap(thread *Thread, capacity int) (SyncHashMap, value.Value)
	NewSyncHashMap(capacity int) SyncHashMap
	All() iter.Seq[value.PairOfValue]
	Length() int
	// If the key is not present should return (Undefined, Undefined)
	GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value)
	// If the key is not present should return (Nil, Undefined)
	GetValNil(thread *Thread, key value.Value) (value.Value, value.Value)
	MergeVal(thread *Thread, other value.Value) (SyncHashMap, value.Value)
	Contains(thread *Thread, pair value.Pair) (bool, value.Value)
	ContainsKey(thread *Thread, key value.Value) (bool, value.Value)
	ContainsValue(thread *Thread, val value.Value) (bool, value.Value)
	Equal(thread *Thread, other value.Value) (bool, value.Value)
	LaxEqual(thread *Thread, other value.Value) (bool, value.Value)
}

func syncHashMapConstructor(class *value.Class) value.Value {
	return NewSyncHashMapOfValue(5).ToValue()
}

// ::Std::Sync::HashMap
func initSyncHashMap() {
	value.SyncHashMapClass.ConstructorFunc = syncHashMapConstructor

	// Instance methods
	c := value.SyncHashMapClass
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncHashMap)
			iterator := self.IterSyncMap()
			return iterator.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncHashMap)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncHashMap)
			key := args[1]
			return self.GetValNil(vm, key)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"[]=",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncHashMap)
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
			self := args[0].AsReference().(SyncHashMap)
			other := args[1]
			return value.RefErr(self.MergeVal(vm, other))
		},
		DefWithParameters(1),
	)
	Alias(c, "merge", "+")

	Def(
		c,
		"remove_key",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncHashMap)
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
			self := args[0].AsReference().(SyncHashMap)
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
			self := args[0].AsReference().(SyncHashMap)
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
			self := args[0].AsReference().(SyncHashMap)
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
			self := args[0].AsReference().(SyncHashMap)
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
			self := args[0].MustReference().(SyncHashMap)
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
			self := args[0].AsReference().(SyncHashMap)
			callable := args[1]
			newMap := NewSyncHashMapOfValue(self.Length())

			self.Lock()
			defer self.Unlock()

			hashMap := self.HashMap()
			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for pair := range hashMap.All() {
					result, err := vm.CallClosure(function, pair.ToValue())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					r, ok := result.SafeAsReference().(value.Pair)
					if !ok {
						return value.Undefined, value.Ref(value.NewArgumentTypeError("pair", result.Class().Name(), value.PairClass.Name()))
					}
					err = HashMapOfValueSet(vm, newMap.Map, r.Key(), r.Value())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newMap), value.Undefined
			}

			// callable is another value
			for pair := range hashMap.All() {
				result, err := vm.CallMethodByName(value.S(symbol.L_call), callable, pair.ToValue())
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				r, ok := result.SafeAsReference().(value.Pair)
				if !ok {
					return value.Undefined, value.Ref(value.NewArgumentTypeError("pair", result.Class().Name(), value.PairClass.Name()))
				}
				err = HashMapOfValueSet(vm, newMap.Map, r.Key(), r.Value())
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
			self := args[0].AsReference().(SyncHashMap)
			callable := args[1]
			newMap := NewSyncHashMapOfValue(self.Length())

			self.Lock()
			defer self.Unlock()

			hashMap := self.HashMap()
			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for pair := range hashMap.All() {
					if pair.Key().IsUndefined() {
						continue
					}
					result, err := vm.CallClosure(function, pair.Value())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					err = HashMapOfValueSet(vm, newMap.Map, pair.Key(), result)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newMap), value.Undefined
			}

			// callable is another value
			for pair := range hashMap.All() {
				if pair.Key().IsUndefined() {
					continue
				}
				result, err := vm.CallMethodByName(value.S(symbol.L_call), callable, pair.Value())
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				err = HashMapOfValueSet(vm, newMap.Map, pair.Key(), result)
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
			self := args[0].AsReference().(SyncHashMap)
			callable := args[1]

			self.Lock()
			defer self.Unlock()

			hashMap := self.HashMap()
			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for pair := range hashMap.All() {
					if pair.Key().IsUndefined() {
						continue
					}
					result, err := vm.CallClosure(function, pair.Value())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					hashMap.SetVal(vm, pair.Key(), result)
				}
				return self.ToValue(), value.Undefined
			}

			// callable is another value
			for pair := range hashMap.All() {
				if pair.Key().IsUndefined() {
					continue
				}
				result, err := vm.CallMethodByName(value.S(symbol.L_call), callable, pair.Value())
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				hashMap.SetVal(vm, pair.Key(), result)
			}
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
}

// ::Std::Sync::HashMap::Iterator
func initSyncHashMapIterator() {
	// Instance methods
	c := value.SyncHashMapIteratorClass
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
