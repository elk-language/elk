package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/google/go-cmp/cmp"
)

type Record interface {
	value.Reference
	value.NativeIterable
	All() iter.Seq[value.PairOfValue]
	Length() int
	// If the key is not present should return (Undefined, Undefined)
	GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value)
	// If the key is not present should return (Nil, Undefined)
	GetValNil(thread *Thread, key value.Value) (value.Value, value.Value)
}

type HashRecord interface {
	Record
	IterRecord() value.NativeResettableIterator
	All() iter.Seq[value.PairOfValue]
	Keys() iter.Seq[value.Value]
	Length() int
	MergeVal(thread *Thread, other value.Value) (value.Value, value.Value)
	Contains(thread *Thread, pair value.Pair) (bool, value.Value)
	ContainsKey(thread *Thread, key value.Value) (bool, value.Value)
	ContainsValue(thread *Thread, val value.Value) (bool, value.Value)
	Equal(thread *Thread, other value.Value) (bool, value.Value)
	LaxEqual(thread *Thread, other value.Value) (bool, value.Value)
	CloneHashRecord(thread *Thread, capacity int) (HashRecord, value.Value)
	NewHashRecord(capacity int) HashRecord
	SetVal(thread *Thread, key, val value.Value) value.Value
}

func hashRecordConstructor(class *value.Class) value.Value {
	return NewHashRecordOfValue(5).ToValue()
}

func HashRecordEqual(vm *Thread, x HashRecord, y HashRecord) (bool, value.Value) {
	if x == y {
		return true, value.Undefined
	}
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for pair := range x.All() {
		yVal, err := y.GetValUndefined(vm, pair.Key())
		if !err.IsUndefined() {
			return false, err
		}
		if yVal.IsUndefined() {
			return false, value.Undefined
		}
		eqVal, err := Equal(vm, pair.Value(), yVal)
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

func HashRecordLaxEqual(vm *Thread, x HashRecord, y HashRecord) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for xPair := range x.All() {
		if xPair.Key().IsUndefined() {
			continue
		}

		yVal, err := y.GetValUndefined(vm, xPair.Key())
		if !err.IsUndefined() {
			return false, err
		}
		if yVal.IsUndefined() {
			return false, value.Undefined
		}
		eqVal, err := LaxEqual(vm, xPair.Value(), yVal)
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

// ::Std::HashRecord
func initHashRecord() {
	value.HashRecordClass.ConstructorFunc = hashRecordConstructor

	// Instance methods
	c := value.HashRecordClass
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashRecord)
			iterator := self.IterRecord()
			return iterator.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(HashRecord)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashRecord)
			key := args[1]
			return self.GetValNil(vm, key)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashRecord)
			other := args[1]
			return self.MergeVal(vm, other)
		},
		DefWithParameters(1),
	)
	Alias(c, "merge", "+")

	Def(
		c,
		"contains",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashRecord)
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
			self := args[0].AsReference().(HashRecord)
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
			self := args[0].AsReference().(HashRecord)
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
			self := args[0].AsReference().(HashRecord)
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
			self := args[0].MustReference().(HashRecord)
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
			self := args[0].MustReference().(HashRecord)
			callable := args[1]
			newRecord := NewHashRecordOfValue(self.Length())

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
					err = HashRecordOfValueSet(vm, newRecord, r.Key(), r.Value())
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
				err = HashRecordOfValueSet(vm, newRecord, r.Key(), r.Value())
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
			self := args[0].AsReference().(HashRecord)
			callable := args[1]
			newRecord := NewHashRecordOfValue(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for pair := range self.All() {
					result, err := vm.CallClosure(function, pair.Value())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					err = HashRecordOfValueSet(vm, newRecord, pair.Key(), result)
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
				err = HashRecordOfValueSet(vm, newRecord, pair.Key(), result)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
			}
			return value.Ref(newRecord), value.Undefined
		},
		DefWithParameters(1),
	)
}

// ::Std::HashRecord::Iterator
func initHashRecordIterator() {
	// Instance methods
	c := value.HashRecordIteratorClass
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

func NewHashRecordComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y HashRecord) bool {
		if x.Length() != y.Length() {
			return false
		}

		result := DefaultThreadPool.Call(func(vm *Thread) (result value.Value, err value.Value) {
			v := New()
			for xPair := range x.All() {
				yVal, err := y.GetValUndefined(v, xPair.Key())
				if !err.IsUndefined() {
					return value.False.ToValue(), value.Undefined
				}

				if !cmp.Equal(xPair.Value(), yVal, *opts...) {
					return value.False.ToValue(), value.Undefined
				}
			}

			return value.True.ToValue(), value.Undefined
		}).MustAwaitSync()

		return value.Truthy(result)
	})
}
