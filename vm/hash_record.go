package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

type HashRecord interface {
	value.ValueInterface
	value.NativeIterable
	IterRecord() value.NativeResettableIterator
	All() iter.Seq[value.PairOfValue]
	Length() int
	GetVal(thread *Thread, key value.Value) (value.Value, value.Value)
	ConcatVal(thread *Thread, other value.Value) (value.Value, value.Value)
	Contains(thread *Thread, pair value.Pair) (bool, value.Value)
	ContainsKey(thread *Thread, key value.Value) (bool, value.Value)
	ContainsValue(thread *Thread, val value.Value) (bool, value.Value)
	Equal(thread *Thread, other value.Value) (bool, value.Value)
	LaxEqual(thread *Thread, other value.Value) (bool, value.Value)
}

type mutableHashRecord interface {
	HashRecord
	SetVal(thread *Thread, key, val value.Value) value.Value
}

// ::Std::HashRecord
func initHashRecord() {
	// Instance methods
	c := &value.HashRecordClass.MethodContainer
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
			result, err := self.GetVal(vm, key)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			if result.IsUndefined() {
				return value.Nil, value.Undefined
			}
			return result, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashRecord)
			other := args[1]
			return self.ConcatVal(vm, other)
		},
		DefWithParameters(1),
	)
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
			if function, ok := callable.SafeAsReference().(*Closure); ok {
				for pair := range self.All() {
					result, err := vm.CallClosure(function, pair.ToValue())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					r, ok := result.SafeAsReference().(value.Pair)
					if !ok {
						return value.Undefined, value.Ref(value.NewArgumentTypeError("pair", result.Class().Name, value.PairClass.Name))
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
				result, err := vm.CallMethodByName(callSymbol, callable, pair.ToValue())
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				r, ok := result.SafeAsReference().(value.Pair)
				if !ok {
					return value.Undefined, value.Ref(value.NewArgumentTypeError("pair", result.Class().Name, value.PairClass.Name))
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
			if function, ok := callable.SafeAsReference().(*Closure); ok {
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
				result, err := vm.CallMethodByName(callSymbol, callable, pair.Value())
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
	c := &value.HashRecordIteratorClass.MethodContainer
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

// Create a new hash record with the given entries.
func NewHashRecordWithElements(vm *Thread, elements ...value.PairOfValue) (*HashRecordOfValue, value.Value) {
	return NewHashRecordWithCapacityAndElements(vm, len(elements), elements...)
}

// Create a new hash record with the given entries.
func MustNewHashRecordWithElements(vm *Thread, elements ...value.PairOfValue) *HashRecordOfValue {
	hrec, err := NewHashRecordWithElements(vm, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return hrec
}

func NewHashRecordWithCapacityAndElements(vm *Thread, capacity int, elements ...value.PairOfValue) (*HashRecordOfValue, value.Value) {
	h := NewHashRecordOfValue(capacity)
	for _, element := range elements {
		err := HashRecordOfValueSet(vm, h, element.Key(), element.Value())
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return h, value.Undefined
}

func MustNewHashRecordWithCapacityAndElements(vm *Thread, capacity int, elements ...value.PairOfValue) *HashRecordOfValue {
	hrec, err := NewHashRecordWithCapacityAndElements(vm, capacity, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return hrec
}

// Delete the given key from the hashMap
func HashRecordOfValueDelete(vm *Thread, hashRecord *HashRecordOfValue, key value.Value) (bool, value.Value) {
	return HashMapOfValueDelete(vm, (*HashMapOfValue)(hashRecord), key)
}

// Get the element under the given key.
func HashRecordOfValueGet(vm *Thread, hashRecord *HashRecordOfValue, key value.Value) (value.Value, value.Value) {
	return HashMapOfValueGet(vm, (*HashMapOfValue)(hashRecord), key)
}

func HashRecordOfValueCopyTable(vm *Thread, target *HashRecordOfValue, source []value.PairOfValue) value.Value {
	return HashMapOfValueCopyTable(vm, (*HashMapOfValue)(target), source)
}

// Copy the pairs of one hash record to the other.
func HashRecordOfValueCopy(vm *Thread, target *HashRecordOfValue, source *HashRecordOfValue) value.Value {
	return HashMapOfValueCopy(vm, (*HashMapOfValue)(target), (*HashMapOfValue)(source))
}

func HashRecordOfValueCopyInterface(vm *Thread, target *HashRecordOfValue, source HashRecord) value.Value {
	return HashMapOfValueCopyInterface(vm, (*HashMapOfValue)(target), source)
}

// Create a new map containing the pairs of both maps.
func HashRecordOfValueConcat(vm *Thread, x *HashRecordOfValue, y *HashRecordOfValue) (*HashRecordOfValue, value.Value) {
	m, err := HashMapOfValueConcat(vm, (*HashMapOfValue)(x), (*HashMapOfValue)(y))
	if err.IsNotUndefined() {
		return nil, err
	}
	return (*HashRecordOfValue)(m), value.Undefined
}

func HashRecordOfValueConcatInterface(vm *Thread, x *HashRecordOfValue, y HashRecord) (*HashRecordOfValue, value.Value) {
	m, err := HashMapOfValueConcatInterface(vm, (*HashMapOfValue)(x), y)
	if err.IsNotUndefined() {
		return nil, err
	}
	return (*HashRecordOfValue)(m), value.Undefined
}

// Check if the given pair is present in the record
func HashRecordOfValueContains(vm *Thread, hrec *HashRecordOfValue, pair value.Pair) (bool, value.Value) {
	return HashMapOfValueContains(vm, (*HashMapOfValue)(hrec), pair)
}

// Check if the given key is present in the record
func HashRecordOfValueContainsKey(vm *Thread, hrec *HashRecordOfValue, key value.Value) (bool, value.Value) {
	return HashMapOfValueContainsKey(vm, (*HashMapOfValue)(hrec), key)
}

// Check if the given value is present in the record
func HashRecordOfValueContainsValue(vm *Thread, hrec *HashRecordOfValue, val value.Value) (bool, value.Value) {
	return HashMapOfValueContainsValue(vm, (*HashMapOfValue)(hrec), val)
}

// Checks whether two hash records are equal (lax)
func HashRecordOfValueLaxEqual(vm *Thread, x *HashRecordOfValue, y *HashRecordOfValue) (bool, value.Value) {
	return HashMapOfValueLaxEqual(vm, (*HashMapOfValue)(x), (*HashMapOfValue)(y))
}

func HashRecordOfValueLaxEqualInterface(vm *Thread, x *HashRecordOfValue, y HashRecord) (bool, value.Value) {
	return HashMapOfValueLaxEqualInterface(vm, (*HashMapOfValue)(x), y)
}

// Checks whether two hash records are equal
func HashRecordOfValueEqual(vm *Thread, x *HashRecordOfValue, y *HashRecordOfValue) (bool, value.Value) {
	return HashMapOfValueEqual(vm, (*HashMapOfValue)(x), (*HashMapOfValue)(y))
}

func HashRecordOfValueEqualInterface(vm *Thread, x *HashRecordOfValue, y HashRecord) (bool, value.Value) {
	return HashMapOfValueEqualInterface(vm, (*HashMapOfValue)(x), y)
}

// Add additional n empty slots for new elements.
func HashRecordOfValueGrow(vm *Thread, hashRecord *HashRecordOfValue, newSlots int) value.Value {
	return HashMapOfValueGrow(vm, (*HashMapOfValue)(hashRecord), newSlots)
}

// Resize the given hash record to the desired capacity.
func HashRecordOfValueSetCapacity(vm *Thread, hashRecord *HashRecordOfValue, capacity int) value.Value {
	return HashMapOfValueSetCapacity(vm, (*HashMapOfValue)(hashRecord), capacity)
}

func HashRecordOfValueSetWithMaxLoad(vm *Thread, hashRecord *HashRecordOfValue, key, val value.Value, maxLoad float64) value.Value {
	return HashMapOfValueSetWithMaxLoad(vm, (*HashMapOfValue)(hashRecord), key, val, maxLoad)
}

// Set a value under the given key.
func HashRecordOfValueSet(vm *Thread, hashRecord *HashRecordOfValue, key, val value.Value) value.Value {
	return HashMapOfValueSet(vm, (*HashMapOfValue)(hashRecord), key, val)
}

func NewHashRecordOfValueComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *HashRecordOfValue) bool {
		if x == y {
			return true
		}
		if x.Length() != y.Length() {
			return false
		}

		v := New()
		for _, xPair := range x.Table {
			if xPair.Key().IsUndefined() {
				continue
			}

			yVal, err := HashRecordOfValueGet(v, y, xPair.Key())
			if !err.IsUndefined() {
				return false
			}

			if !cmp.Equal(xPair.Value(), yVal, *opts...) {
				return false
			}

		}

		return true
	})
}
