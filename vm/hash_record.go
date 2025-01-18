package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

// ::Std::HashRecord
func initHashRecord() {
	// Instance methods
	c := &value.HashRecordClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashRecord)
			iterator := value.NewHashRecordIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashRecord)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashRecord)
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
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashRecord)
			other := args[1]

			switch o := other.SafeAsReference().(type) {
			case *value.HashMap:
				result, err := HashRecordConcat(vm, self, (*value.HashRecord)(o))
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.Ref(result), value.Undefined
			case *value.HashRecord:
				result, err := HashRecordConcat(vm, self, o)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.Ref((*value.HashRecord)(result)), value.Undefined
			default:
				return value.Undefined, value.Ref(value.NewCoerceError(value.HashRecordClass, other.Class()))
			}
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashRecord)
			otherVal := args[1]
			switch other := otherVal.SafeAsReference().(type) {
			case *value.Pair:
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
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashRecord)
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
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashRecord)
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
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashRecord)
			other, ok := args[1].SafeAsReference().(*value.HashRecord)
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
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashRecord)
			switch other := args[1].SafeAsReference().(type) {
			case *value.HashRecord:
				equal, err := HashRecordLaxEqual(vm, self, other)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.ToElkBool(equal), value.Undefined
			case *value.HashMap:
				equal, err := HashRecordLaxEqual(vm, self, (*value.HashRecord)(other))
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
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashRecord)
			callable := args[1]
			newRecord := value.NewHashRecord(self.Length())

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
					r, ok := result.SafeAsReference().(*value.Pair)
					if !ok {
						return value.Undefined, value.Ref(value.NewArgumentTypeError("pair", result.Class().Name, value.PairClass.Name))
					}
					err = HashRecordSet(vm, newRecord, r.Key, r.Value)
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
				r, ok := result.SafeAsReference().(*value.Pair)
				if !ok {
					return value.Undefined, value.Ref(value.NewArgumentTypeError("pair", result.Class().Name, value.PairClass.Name))
				}
				err = HashRecordSet(vm, newRecord, r.Key, r.Value)
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
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashRecord)
			callable := args[1]
			newRecord := value.NewHashRecord(self.Length())

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
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.HashRecordIterator)(args[0].Pointer())
			return self.Next()
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"reset",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.HashRecordIterator)(args[0].Pointer())
			self.Reset()
			return args[0], value.Undefined
		},
	)

}

// Create a new hash record with the given entries.
func NewHashRecordWithElements(vm *VM, elements ...value.Pair) (*value.HashRecord, value.Value) {
	return NewHashRecordWithCapacityAndElements(vm, len(elements), elements...)
}

// Create a new hash record with the given entries.
func MustNewHashRecordWithElements(vm *VM, elements ...value.Pair) *value.HashRecord {
	hrec, err := NewHashRecordWithElements(vm, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return hrec
}

func NewHashRecordWithCapacityAndElements(vm *VM, capacity int, elements ...value.Pair) (*value.HashRecord, value.Value) {
	h := value.NewHashRecord(capacity)
	for _, element := range elements {
		err := HashRecordSet(vm, h, element.Key, element.Value)
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return h, value.Undefined
}

func MustNewHashRecordWithCapacityAndElements(vm *VM, capacity int, elements ...value.Pair) *value.HashRecord {
	hrec, err := NewHashRecordWithCapacityAndElements(vm, capacity, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return hrec
}

// Delete the given key from the hashMap
func HashRecordDelete(vm *VM, hashRecord *value.HashRecord, key value.Value) (bool, value.Value) {
	return HashMapDelete(vm, (*value.HashMap)(hashRecord), key)
}

// Get the element under the given key.
func HashRecordGet(vm *VM, hashRecord *value.HashRecord, key value.Value) (value.Value, value.Value) {
	return HashMapGet(vm, (*value.HashMap)(hashRecord), key)
}

func HashRecordCopyTable(vm *VM, target *value.HashRecord, source []value.Pair) value.Value {
	return HashMapCopyTable(vm, (*value.HashMap)(target), source)
}

// Copy the pairs of one hash record to the other.
func HashRecordCopy(vm *VM, target *value.HashRecord, source *value.HashRecord) value.Value {
	return HashMapCopy(vm, (*value.HashMap)(target), (*value.HashMap)(source))
}

// Create a new map containing the pairs of both maps.
func HashRecordConcat(vm *VM, x *value.HashRecord, y *value.HashRecord) (*value.HashMap, value.Value) {
	return HashMapConcat(vm, (*value.HashMap)(x), (*value.HashMap)(y))
}

// Check if the given pair is present in the record
func HashRecordContains(vm *VM, hrec *value.HashRecord, pair *value.Pair) (bool, value.Value) {
	return HashMapContains(vm, (*value.HashMap)(hrec), pair)
}

// Check if the given key is present in the record
func HashRecordContainsKey(vm *VM, hrec *value.HashRecord, key value.Value) (bool, value.Value) {
	return HashMapContainsKey(vm, (*value.HashMap)(hrec), key)
}

// Check if the given value is present in the record
func HashRecordContainsValue(vm *VM, hrec *value.HashRecord, val value.Value) (bool, value.Value) {
	return HashMapContainsValue(vm, (*value.HashMap)(hrec), val)
}

// Checks whether two hash records are equal (lax)
func HashRecordLaxEqual(vm *VM, x *value.HashRecord, y *value.HashRecord) (bool, value.Value) {
	return HashMapLaxEqual(vm, (*value.HashMap)(x), (*value.HashMap)(y))
}

// Checks whether two hash records are equal
func HashRecordEqual(vm *VM, x *value.HashRecord, y *value.HashRecord) (bool, value.Value) {
	return HashMapEqual(vm, (*value.HashMap)(x), (*value.HashMap)(y))
}

// Add additional n empty slots for new elements.
func HashRecordGrow(vm *VM, hashRecord *value.HashRecord, newSlots int) value.Value {
	return HashMapGrow(vm, (*value.HashMap)(hashRecord), newSlots)
}

// Resize the given hash record to the desired capacity.
func HashRecordSetCapacity(vm *VM, hashRecord *value.HashRecord, capacity int) value.Value {
	return HashMapSetCapacity(vm, (*value.HashMap)(hashRecord), capacity)
}

func HashRecordSetWithMaxLoad(vm *VM, hashRecord *value.HashRecord, key, val value.Value, maxLoad float64) value.Value {
	return HashMapSetWithMaxLoad(vm, (*value.HashMap)(hashRecord), key, val, maxLoad)
}

// Set a value under the given key.
func HashRecordSet(vm *VM, hashRecord *value.HashRecord, key, val value.Value) value.Value {
	return HashMapSet(vm, (*value.HashMap)(hashRecord), key, val)
}

func NewHashRecordComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *value.HashRecord) bool {
		if x == y {
			return true
		}
		if x.Length() != y.Length() {
			return false
		}

		v := New()
		for _, xPair := range x.Table {
			if xPair.Key.IsUndefined() {
				continue
			}

			yVal, err := HashRecordGet(v, y, xPair.Key)
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
