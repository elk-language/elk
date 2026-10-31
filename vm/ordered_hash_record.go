package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type OrderedHashRecord interface {
	Record
	IterOrderedRecord() value.NativeResettableIterator
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
	MergeVal(thread *Thread, other value.Value) (OrderedHashRecord, value.Value)
	Contains(thread *Thread, pair value.Pair) (bool, value.Value)
	ContainsKey(thread *Thread, key value.Value) (bool, value.Value)
	ContainsValue(thread *Thread, val value.Value) (bool, value.Value)
	Equal(thread *Thread, other value.Value) (bool, value.Value)
	LaxEqual(thread *Thread, other value.Value) (bool, value.Value)
	CloneOrderedHashRecord(thread *Thread, capacity int) (OrderedHashRecord, value.Value)
	NewOrderedHashRecord(capacity int) OrderedHashRecord
	SetVal(thread *Thread, key, val value.Value) value.Value
	ToHashRecord(thread *Thread) (HashRecord, value.Value)
	ToHashMap(thread *Thread) (HashMap, value.Value)
	ToSyncHashMap(thread *Thread) (SyncHashMap, value.Value)
	ToSyncOrderedHashMap() SyncOrderedHashMap
	ToOrderedHashMap() OrderedHashMap
}

func orderedHashRecordConstructor(class *value.Class) value.Value {
	return NewOrderedHashRecordOfValue(5).ToValue()
}

// ::Std::OrderedHashRecord
func initOrderedHashRecord() {
	value.OrderedHashRecordClass.ConstructorFunc = orderedHashRecordConstructor

	// Instance methods
	c := value.OrderedHashRecordClass
	Def(
		c,
		"#init",
		func(thread *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(*OrderedHashRecordOfValue)
			entries := args[1]

			for entry, err := range Iterate(thread, entries) {
				if err.IsNotUndefined() {
					return value.Undefined, err
				}

				pair := entry.AsReference().(value.Pair)
				self.SetVal(thread, pair.Key(), pair.Value())
			}

			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashRecord)
			iterator := self.IterOrderedRecord()
			return iterator.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_hash_map",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashRecord)
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
			self := args[0].AsReference().(OrderedHashRecord)
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
			self := args[0].AsReference().(OrderedHashRecord)
			result, err := self.ToHashRecord(vm)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}
			return result.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_sync_ordered_hash_map",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashRecord)
			return self.ToSyncOrderedHashMap().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_ordered_hash_map",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashRecord)
			return self.ToOrderedHashMap().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_ordered_hash_record",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(OrderedHashRecord)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashRecord)
			key := args[1]
			return self.GetValNil(vm, key)
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"+",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashRecord)
			other := args[1]
			return value.RefErr(self.MergeVal(vm, other))
		},
		DefWithParameters(1),
	)
	Alias(c, "merge", "+")

	Def(
		c,
		"try_first",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashRecord)
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
			self := args[0].AsReference().(OrderedHashRecord)
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
			self := args[0].AsReference().(OrderedHashRecord)
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
			self := args[0].AsReference().(OrderedHashRecord)
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
			self := args[0].AsReference().(OrderedHashRecord)
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
			self := args[0].AsReference().(OrderedHashRecord)
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
			self := args[0].AsReference().(OrderedHashRecord)
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
			self := args[0].AsReference().(OrderedHashRecord)
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
			self := args[0].AsReference().(OrderedHashRecord)
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
			self := args[0].AsReference().(OrderedHashRecord)
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
			self := args[0].MustReference().(OrderedHashRecord)
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
			self := args[0].MustReference().(OrderedHashRecord)
			callable := args[1]
			newRecord := NewOrderedHashRecordOfValue(self.Length())

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
					err = newRecord.SetVal(vm, r.Key(), r.Value())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newRecord), value.Undefined
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
				err = newRecord.SetVal(vm, r.Key(), r.Value())
				if !err.IsUndefined() {
					return value.Undefined, err
				}
			}
			return value.Ref(newRecord), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map_values",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(OrderedHashRecord)
			callable := args[1]
			newRecord := NewOrderedHashRecordOfValue(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for pair := range self.All() {
					result, err := vm.CallClosure(function, pair.Value())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					err = newRecord.SetVal(vm, pair.Key(), result)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newRecord), value.Undefined
			}

			// callable is another value
			for pair := range self.All() {
				result, err := vm.CallMethodByName(value.S(symbol.L_call), callable, pair.Value())
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				err = newRecord.SetVal(vm, pair.Key(), result)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
			}
			return value.Ref(newRecord), value.Undefined
		},
		DefWithParameters(1),
	)
}

// ::Std::OrderedHashRecord::Iterator
func initOrderedHashRecordIterator() {
	// Instance methods
	c := value.OrderedHashRecordIteratorClass
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

func OrderedHashRecordLaxEqualRecord(thread *Thread, a OrderedHashRecord, b Record) (bool, value.Value) {
	if a.Length() != b.Length() {
		return false, value.Undefined
	}

	for aEntry := range a.All() {
		bVal, err := b.GetValUndefined(thread, aEntry.Key())
		if err.IsNotUndefined() {
			return false, err
		}
		if bVal.IsUndefined() {
			return false, value.Undefined
		}

		eqVal, err := LaxEqual(thread, aEntry.Value(), bVal)
		if !err.IsUndefined() {
			return false, err
		}
		equal := value.Truthy(eqVal)
		if !equal {
			return false, value.Undefined
		}
	}

	return true, value.Undefined
}
