package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type SyncOrderedHashMap interface {
	Record
	OrderedHashMap() OrderedHashMap
	IterSyncOrderedMap() value.NativeResettableIterator
	RemoveKey(thread *Thread, key value.Value) (bool, value.Value)
	SetVal(thread *Thread, key, val value.Value) value.Value
	Lock()
	Unlock()

	All() iter.Seq[value.PairOfValue]
	Length() int
	// Normalizes the index so that `-1` becomes `length-1`.
	// Does bounds checking and returns (undefined, undefined) if the index is out of range.
	TryAtVal(value.Value) (value.Pair, value.Value)
	AtInterface(int) value.Pair
	GetAtInterface(int) (value.Pair, value.Value)
	// If the key is not present should return (Undefined, Undefined)
	GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value)
	// If the key is not present should return (Nil, Undefined)
	GetValNil(thread *Thread, key value.Value) (value.Value, value.Value)
	MergeVal(thread *Thread, other value.Value) (SyncOrderedHashMap, value.Value)
	Contains(thread *Thread, pair value.Pair) (bool, value.Value)
	ContainsKey(thread *Thread, key value.Value) (bool, value.Value)
	ContainsValue(thread *Thread, val value.Value) (bool, value.Value)
	Equal(thread *Thread, other value.Value) (bool, value.Value)
	LaxEqual(thread *Thread, other value.Value) (bool, value.Value)
	CloneSyncOrderedHashMap(thread *Thread, capacity int) (SyncOrderedHashMap, value.Value)
	NewSyncOrderedHashMap(capacity int) SyncOrderedHashMap
	ToHashRecord(thread *Thread) (HashRecord, value.Value)
	ToHashMap(thread *Thread) (HashMap, value.Value)
	ToSyncHashMap(thread *Thread) (SyncHashMap, value.Value)
	ToOrderedHashMap() OrderedHashMap
	ToOrderedHashRecord() OrderedHashRecord
}

func syncOrderedHashMapConstructor(class *value.Class) value.Value {
	return NewSyncOrderedHashMapOfValue(5).ToValue()
}

// ::Std::Sync::OrderedHashMap
func initSyncOrderedHashMap() {
	value.SyncOrderedHashMapClass.ConstructorFunc = syncOrderedHashMapConstructor

	// Instance methods
	c := value.SyncOrderedHashMapClass
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
			iterator := self.IterSyncOrderedMap()
			return iterator.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_hash_map",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
			result, err := self.ToHashMap(vm)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}
			return result.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_sync_hash_map",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
			result, err := self.ToSyncHashMap(vm)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}
			return result.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_hash_record",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
			result, err := self.ToHashRecord(vm)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}
			return result.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_ordered_hash_map",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
			return self.ToOrderedHashMap().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_sync_ordered_hash_map",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"to_ordered_hash_record",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
			return self.ToOrderedHashRecord().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
			key := args[1]
			return self.GetValNil(vm, key)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"[]=",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
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
			self := args[0].AsReference().(SyncOrderedHashMap)
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
			self := args[0].AsReference().(SyncOrderedHashMap)
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
			self := args[0].AsReference().(SyncOrderedHashMap)
			length := self.Length()
			if length == 0 {
				return value.Nil, value.Undefined
			}

			return self.AtInterface(0).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"first",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
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
	)
	Def(
		c,
		"try_last",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
			length := self.Length()
			if length == 0 {
				return value.Nil, value.Undefined
			}

			return self.AtInterface(length - 1).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"last",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
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
	)

	Def(
		c,
		"get_at",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
			index := args[1].AsAnyInt()
			return value.RefErr(self.GetAtInterface(index))
		},
		DefWithParameters(1),
	)
	Alias(c, "at", "get_at")

	Def(
		c,
		"try_at",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
			other := args[1]
			pair, err := self.TryAtVal(other)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			if pair == nil {
				return value.Nil, value.Undefined
			}
			return pair.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"contains",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashMap)
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
			self := args[0].AsReference().(SyncOrderedHashMap)
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
			self := args[0].AsReference().(SyncOrderedHashMap)
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
			self := args[0].AsReference().(SyncOrderedHashMap)
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
			self := args[0].MustReference().(SyncOrderedHashMap)
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
			self := args[0].AsReference().(SyncOrderedHashMap)
			callable := args[1]
			newMap := NewSyncOrderedHashMapOfValue(self.Length())

			self.Lock()
			defer self.Unlock()

			hashMap := self.OrderedHashMap()
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
					err = newMap.Map.SetVal(vm, r.Key(), r.Value())
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
				err = newMap.Map.SetVal(vm, r.Key(), r.Value())
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
			self := args[0].AsReference().(SyncOrderedHashMap)
			callable := args[1]
			newMap := NewSyncOrderedHashMapOfValue(self.Length())

			self.Lock()
			defer self.Unlock()

			hashMap := self.OrderedHashMap()
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
					newMap.Map.SetVal(vm, pair.Key(), result)
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
				newMap.Map.SetVal(vm, pair.Key(), result)
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
			self := args[0].AsReference().(SyncOrderedHashMap)
			callable := args[1]

			self.Lock()
			defer self.Unlock()

			hashMap := self.OrderedHashMap()
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

// ::Std::Sync::OrderedHashMap::Iterator
func initSyncOrderedHashMapIterator() {
	// Instance methods
	c := value.SyncOrderedHashMapIteratorClass
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
