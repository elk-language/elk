package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

type HashRecord interface {
	value.ValueInterface
	value.NativeIterable
	IterRecord() HashRecordIterator
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

type HashRecordIterator interface {
	value.NativeIterator
	value.ValueInterface
	Reset()
}

// ::Std::HashRecord
func initHashRecord() {
	// Instance methods
	c := &value.HashRecordClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*HashRecordOfValue)
			iterator := value.NewHashRecordIteratorOfValue(self)
			return value.Ref(iterator), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*HashRecordOfValue)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*HashRecordOfValue)
			key := args[1]
			result, err := HashRecordGet(vm, self, key)
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
			self := args[0].MustReference().(*HashRecordOfValue)
			other := args[1]

			switch o := other.SafeAsReference().(type) {
			case *HashMapOfValue:
				result, err := HashRecordConcat(vm, self, (*HashRecordOfValue)(o))
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.Ref(result), value.Undefined
			case *HashRecordOfValue:
				result, err := HashRecordConcat(vm, self, o)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.Ref((*HashRecordOfValue)(result)), value.Undefined
			default:
				return value.Undefined, value.Ref(value.NewCoerceError(value.HashRecordClass, other.Class()))
			}
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*HashRecordOfValue)
			otherVal := args[1]
			switch other := otherVal.SafeAsReference().(type) {
			case *value.PairOfValue:
				contains, err := HashRecordContains(vm, self, other)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				return value.ToElkBool(contains), value.Undefined
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
			self := args[0].MustReference().(*HashRecordOfValue)
			contains, err := HashRecordContainsKey(vm, self, args[1])
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			return value.ToElkBool(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains_value",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*HashRecordOfValue)
			contains, err := HashRecordContainsValue(vm, self, args[1])
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			return value.ToElkBool(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*HashRecordOfValue)
			other, ok := args[1].SafeAsReference().(*HashRecordOfValue)
			if !ok {
				return value.False, value.Undefined
			}
			equal, err := HashRecordEqual(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.ToElkBool(equal), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*HashRecordOfValue)
			switch other := args[1].SafeAsReference().(type) {
			case *HashRecordOfValue:
				equal, err := HashRecordLaxEqual(vm, self, other)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.ToElkBool(equal), value.Undefined
			case *HashMapOfValue:
				equal, err := HashRecordLaxEqual(vm, self, (*HashRecordOfValue)(other))
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.ToElkBool(equal), value.Undefined
			default:
				return value.False, value.Undefined
			}
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*HashRecordOfValue)
			callable := args[1]
			newRecord := value.NewHashRecordOfValue(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(*Closure); ok {
				for i, pair := range self.Table {
					if pair.Key.IsUndefined() {
						continue
					}
					result, err := vm.CallClosure(function, value.Ref(&self.Table[i]))
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					r, ok := result.SafeAsReference().(*value.PairOfValue)
					if !ok {
						return value.Undefined, value.Ref(value.NewArgumentTypeError("pair", result.Class().Name, value.PairClass.Name))
					}
					err = HashRecordSet(vm, newRecord, r.Key(), r.Value())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newRecord), value.Undefined
			}

			// callable is another value
			for i, pair := range self.Table {
				if pair.Key.IsUndefined() {
					continue
				}
				result, err := vm.CallMethodByName(callSymbol, callable, value.Ref(&self.Table[i]))
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				r, ok := result.SafeAsReference().(*value.PairOfValue)
				if !ok {
					return value.Undefined, value.Ref(value.NewArgumentTypeError("pair", result.Class().Name, value.PairClass.Name))
				}
				err = HashRecordSet(vm, newRecord, r.Key(), r.Value())
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
			self := args[0].MustReference().(*HashRecordOfValue)
			callable := args[1]
			newRecord := value.NewHashRecordOfValue(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(*Closure); ok {
				for _, pair := range self.Table {
					if pair.Key.IsUndefined() {
						continue
					}
					result, err := vm.CallClosure(function, pair.Value)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					err = HashRecordSet(vm, newRecord, pair.Key, result)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newRecord), value.Undefined
			}

			// callable is another value
			for _, pair := range self.Table {
				if pair.Key.IsUndefined() {
					continue
				}
				result, err := vm.CallMethodByName(callSymbol, callable, pair.Value)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				err = HashRecordSet(vm, newRecord, pair.Key, result)
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
			self := args[0].AsReference().(HashRecordIterator)
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
			self := args[0].AsReference().(HashRecordIterator)
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
		err := HashRecordSet(vm, h, element.Key(), element.Value())
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
func HashRecordDelete(vm *Thread, hashRecord *HashRecordOfValue, key value.Value) (bool, value.Value) {
	return HashMapOfValueDelete(vm, (*HashMapOfValue)(hashRecord), key)
}

// Get the element under the given key.
func HashRecordGet(vm *Thread, hashRecord *HashRecordOfValue, key value.Value) (value.Value, value.Value) {
	return HashMapOfValueGet(vm, (*HashMapOfValue)(hashRecord), key)
}

func HashRecordCopyTable(vm *Thread, target *HashRecordOfValue, source []value.PairOfValue) value.Value {
	return HashMapOfValueCopyTable(vm, (*HashMapOfValue)(target), source)
}

// Copy the pairs of one hash record to the other.
func HashRecordCopy(vm *Thread, target *HashRecordOfValue, source *HashRecordOfValue) value.Value {
	return HashMapOfValueCopy(vm, (*HashMapOfValue)(target), (*HashMapOfValue)(source))
}

// Create a new map containing the pairs of both maps.
func HashRecordConcat(vm *Thread, x *HashRecordOfValue, y *HashRecordOfValue) (*HashMapOfValue, value.Value) {
	return HashMapOfValueConcat(vm, (*HashMapOfValue)(x), (*HashMapOfValue)(y))
}

// Check if the given pair is present in the record
func HashRecordContains(vm *Thread, hrec *HashRecordOfValue, pair *value.PairOfValue) (bool, value.Value) {
	return HashMapOfValueContains(vm, (*HashMapOfValue)(hrec), pair)
}

// Check if the given key is present in the record
func HashRecordContainsKey(vm *Thread, hrec *HashRecordOfValue, key value.Value) (bool, value.Value) {
	return HashMapOfValueContainsKey(vm, (*HashMapOfValue)(hrec), key)
}

// Check if the given value is present in the record
func HashRecordContainsValue(vm *Thread, hrec *HashRecordOfValue, val value.Value) (bool, value.Value) {
	return HashMapOfValueContainsValue(vm, (*HashMapOfValue)(hrec), val)
}

// Checks whether two hash records are equal (lax)
func HashRecordLaxEqual(vm *Thread, x *HashRecordOfValue, y *HashRecordOfValue) (bool, value.Value) {
	return HashMapOfValueLaxEqual(vm, (*HashMapOfValue)(x), (*HashMapOfValue)(y))
}

// Checks whether two hash records are equal
func HashRecordEqual(vm *Thread, x *HashRecordOfValue, y *HashRecordOfValue) (bool, value.Value) {
	return HashMapOfValueEqual(vm, (*HashMapOfValue)(x), (*HashMapOfValue)(y))
}

// Add additional n empty slots for new elements.
func HashRecordGrow(vm *Thread, hashRecord *HashRecordOfValue, newSlots int) value.Value {
	return HashMapOfValueGrow(vm, (*HashMapOfValue)(hashRecord), newSlots)
}

// Resize the given hash record to the desired capacity.
func HashRecordSetCapacity(vm *Thread, hashRecord *HashRecordOfValue, capacity int) value.Value {
	return HashMapOfValueSetCapacity(vm, (*HashMapOfValue)(hashRecord), capacity)
}

func HashRecordSetWithMaxLoad(vm *Thread, hashRecord *HashRecordOfValue, key, val value.Value, maxLoad float64) value.Value {
	return HashMapOfValueSetWithMaxLoad(vm, (*HashMapOfValue)(hashRecord), key, val, maxLoad)
}

// Set a value under the given key.
func HashRecordSet(vm *Thread, hashRecord *HashRecordOfValue, key, val value.Value) value.Value {
	return HashMapOfValueSet(vm, (*HashMapOfValue)(hashRecord), key, val)
}

func NewHashRecordComparer(opts *cmp.Options) cmp.Option {
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

			yVal, err := HashRecordGet(v, y, xPair.Key())
			if !err.IsUndefined() {
				return false
			}

			if !cmp.Equal(xPair.Value, yVal, *opts...) {
				return false
			}

		}

		return true
	})
}
