package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

type HashMap interface {
	HashRecord
	IterMap() HashMapIterator
	SetVal(thread *Thread, key, val value.Value) value.Value
}

type HashMapIterator interface {
	HashRecordIterator
}

// ::Std::HashMap
func initHashMap() {
	// Instance methods
	c := &value.HashMapClass.MethodContainer
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
			return self.ConcatVal(vm, other)
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
			self := args[0].AsReference().(HashMap)
			contains, err := self.ContainsKey(vm, args[1])
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
			self := args[0].AsReference().(HashMap)
			contains, err := self.ContainsValue(vm, args[1])
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
			self := args[0].AsReference().(HashMap)
			other := args[1]
			equal, err := self.Equal(vm, other)
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
			self := args[0].MustReference().(HashMap)
			equal, err := self.LaxEqual(vm, args[1])
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.ToElkBool(equal), value.Undefined
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
					err = HashMapOfValueSet(vm, newMap, r.Key(), r.Value())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newMap), value.Undefined
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
			if function, ok := callable.SafeAsReference().(*Closure); ok {
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
				result, err := vm.CallMethodByName(callSymbol, callable, pair.Value())
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
			if function, ok := callable.SafeAsReference().(*Closure); ok {
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
				result, err := vm.CallMethodByName(callSymbol, callable, pair.Value())
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
	c := &value.HashMapIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashMapIterator)
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
			self := args[0].AsReference().(HashMapIterator)
			self.Reset()
			return args[0], value.Undefined
		},
	)

}

func NewHashMapOfValueComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *HashMapOfValue) bool {
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

			yVal, err := HashMapOfValueGet(v, y, xPair.Key())
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

// Create a new hashmap with the given entries.
func NewHashMapOfValueWithElements(vm *Thread, elements ...value.PairOfValue) (*HashMapOfValue, value.Value) {
	return NewHashMapOfValueWithCapacityAndElements(vm, len(elements), elements...)
}

// Create a new hashmap with the given entries.
func MustNewHashMapOfValueWithElements(vm *Thread, elements ...value.PairOfValue) *HashMapOfValue {
	hmap, err := NewHashMapOfValueWithElements(vm, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return hmap
}

func NewHashMapOfValueWithCapacityAndElements(vm *Thread, capacity int, elements ...value.PairOfValue) (*HashMapOfValue, value.Value) {
	h := NewHashMapOfValue(capacity)
	for _, element := range elements {
		err := HashMapOfValueSet(vm, h, element.Key(), element.Value())
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return h, value.Undefined
}

func MustNewHashMapOfValueWithCapacityAndElements(vm *Thread, capacity int, elements ...value.PairOfValue) *HashMapOfValue {
	hmap, err := NewHashMapOfValueWithCapacityAndElements(vm, capacity, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return hmap
}

// Checks whether two hash maps are equal (lax)
func HashMapOfValueLaxEqual(vm *Thread, x *HashMapOfValue, y *HashMapOfValue) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for _, xPair := range x.Table {
		if xPair.Key().IsUndefined() {
			continue
		}

		yVal, err := HashMapOfValueGet(vm, y, xPair.Key())
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

func HashMapOfValueLaxEqualInterface(vm *Thread, x *HashMapOfValue, y HashRecord) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for _, xPair := range x.Table {
		if xPair.Key().IsUndefined() {
			continue
		}

		yVal, err := y.GetVal(vm, xPair.Key())
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

func HashRecordLaxEqual(vm *Thread, x HashRecord, y HashRecord) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for xPair := range x.All() {
		if xPair.Key().IsUndefined() {
			continue
		}

		yVal, err := y.GetVal(vm, xPair.Key())
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

// Create a new map containing the pairs of both maps.
func HashMapOfValueConcat(vm *Thread, x, y *HashMapOfValue) (*HashMapOfValue, value.Value) {
	result := x.Clone()
	err := HashMapOfValueCopy(vm, result, y)
	if !err.IsUndefined() {
		return nil, err
	}
	return result, value.Undefined
}

func HashMapOfValueConcatInterface(vm *Thread, x *HashMapOfValue, y HashRecord) (*HashMapOfValue, value.Value) {
	result := x.Clone()
	err := HashMapOfValueCopyInterface(vm, result, y)
	if !err.IsUndefined() {
		return nil, err
	}
	return result, value.Undefined
}

// Checks whether two hash maps are equal
func HashMapOfValueEqual(vm *Thread, x *HashMapOfValue, y *HashMapOfValue) (bool, value.Value) {
	if x == y {
		return true, value.Undefined
	}
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for _, xPair := range x.Table {
		if xPair.Key().IsUndefined() {
			continue
		}

		yVal, err := HashMapOfValueGet(vm, y, xPair.Key())
		if !err.IsUndefined() {
			return false, err
		}
		if yVal.IsUndefined() {
			return false, value.Undefined
		}
		eqVal, err := Equal(vm, xPair.Value(), yVal)
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

func HashMapOfValueEqualInterface(vm *Thread, x *HashMapOfValue, y HashRecord) (bool, value.Value) {
	if x == y {
		return true, value.Undefined
	}
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for _, xPair := range x.Table {
		if xPair.Key().IsUndefined() {
			continue
		}

		yVal, err := y.GetVal(vm, xPair.Key())
		if !err.IsUndefined() {
			return false, err
		}
		if yVal.IsUndefined() {
			return false, value.Undefined
		}
		eqVal, err := Equal(vm, xPair.Value(), yVal)
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

func HashRecordEqual(vm *Thread, x HashRecord, y HashRecord) (bool, value.Value) {
	if x == y {
		return true, value.Undefined
	}
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for pair := range x.All() {
		yVal, err := y.GetVal(vm, pair.Key())
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

// Delete the given key from the hashMap
func HashMapOfValueDelete(vm *Thread, hashMap *HashMapOfValue, key value.Value) (bool, value.Value) {
	if hashMap.Length() == 0 {
		return false, value.Undefined
	}

	index, err := HashMapOfValueIndex(vm, hashMap, key)
	if !err.IsUndefined() {
		return false, err
	}
	if index < 0 {
		return false, value.Undefined
	}
	if hashMap.Table[index].Key().IsUndefined() {
		return false, value.Undefined
	}

	hashMap.Table[index] = value.MakePairOfValue(
		value.Undefined,
		value.True,
	)
	hashMap.Elements--
	hashMap.version++

	return true, value.Undefined
}

// Get the element under the given key.
// Returns (value, undefined) when the value has been found.
// Returns (undefined, undefined) when the key is not present.
// Returns (undefined, err) when there was an error.
func HashMapOfValueGet(vm *Thread, hashMap *HashMapOfValue, key value.Value) (value.Value, value.Value) {
	if hashMap.Length() == 0 {
		return value.Undefined, value.Undefined
	}

	index, err := HashMapOfValueIndex(vm, hashMap, key)
	if !err.IsUndefined() {
		return value.Undefined, err
	}
	if index == -1 {
		return value.Undefined, value.Undefined
	}

	return hashMap.Table[index].Value(), value.Undefined
}

// Check if the given pair is present in the map
func HashMapOfValueContains(vm *Thread, hashMap *HashMapOfValue, pair value.Pair) (bool, value.Value) {
	val, err := HashMapOfValueGet(vm, hashMap, pair.Key())
	if !err.IsUndefined() {
		return false, err
	}
	if val.IsUndefined() {
		return false, value.Undefined
	}

	equal, err := Equal(vm, val, pair.Value())
	if !err.IsUndefined() {
		return false, err
	}

	return value.Truthy(equal), value.Undefined
}

// Check if the given key is present in the map
func HashMapOfValueContainsKey(vm *Thread, hashMap *HashMapOfValue, key value.Value) (bool, value.Value) {
	if hashMap.Length() == 0 {
		return false, value.Undefined
	}

	index, err := HashMapOfValueIndex(vm, hashMap, key)
	if !err.IsUndefined() {
		return false, err
	}
	if index == -1 {
		return false, value.Undefined
	}

	pair := hashMap.Table[index]
	if pair.Key().IsUndefined() {
		return false, value.Undefined
	}

	return true, value.Undefined
}

// Check if the given value is present in the map
func HashMapOfValueContainsValue(vm *Thread, hashMap *HashMapOfValue, val value.Value) (bool, value.Value) {
	for _, pair := range hashMap.Table {
		// if the Key is undefined the entry is empty or deleted
		// so we skip it
		if pair.Key().IsUndefined() {
			continue
		}

		equal, err := Equal(vm, pair.Value(), val)
		if !err.IsUndefined() {
			return false, err
		}

		if value.Truthy(equal) {
			return true, value.Undefined
		}
	}

	return false, value.Undefined
}

func HashMapOfValueCopyTable(vm *Thread, target *HashMapOfValue, source []value.PairOfValue) value.Value {
	for _, entry := range source {
		if entry.Key().IsUndefined() {
			continue
		}

		err := HashMapOfValueSetWithMaxLoad(vm, target, entry.Key(), entry.Value(), 1)
		if !err.IsUndefined() {
			return err
		}
	}

	target.version++
	return value.Undefined
}

// Copy the pairs of one hashmap to the other.
func HashMapOfValueCopy(vm *Thread, target *HashMapOfValue, source *HashMapOfValue) value.Value {
	requiredCapacity := target.Length() + source.Length()
	if target.Capacity() < requiredCapacity {
		HashMapOfValueSetCapacity(vm, target, requiredCapacity)
	}

	target.version++
	for _, entry := range source.Table {
		if entry.Key().IsUndefined() {
			continue
		}

		i, err := HashMapOfValueIndex(vm, target, entry.Key())
		if !err.IsUndefined() {
			return err
		}
		if i == -1 {
			panic("no room in target hashmap during copy")
		}
		target.Table[i] = entry
		target.OccupiedSlots++
		target.Elements++
	}

	return value.Undefined
}

func HashMapOfValueCopyInterface(vm *Thread, target *HashMapOfValue, source HashRecord) value.Value {
	requiredCapacity := target.Length() + source.Length()
	if target.Capacity() < requiredCapacity {
		HashMapOfValueSetCapacity(vm, target, requiredCapacity)
	}

	target.version++
	for entry := range source.All() {
		i, err := HashMapOfValueIndex(vm, target, entry.Key())
		if !err.IsUndefined() {
			return err
		}
		if i == -1 {
			panic("no room in target hashmap during copy")
		}
		target.Table[i] = entry
		target.OccupiedSlots++
		target.Elements++
	}

	return value.Undefined
}

// Add additional n empty slots for new elements.
func HashMapOfValueGrow(vm *Thread, hashMap *HashMapOfValue, newSlots int) value.Value {
	return HashMapOfValueSetCapacity(vm, hashMap, hashMap.Capacity()+newSlots)
}

// Resize the given hashmap to the desired capacity.
func HashMapOfValueSetCapacity(vm *Thread, hashMap *HashMapOfValue, capacity int) value.Value {
	if hashMap.Capacity() == capacity {
		return value.Undefined
	}

	oldTable := hashMap.Table
	newTable := make([]value.PairOfValue, capacity)
	tmpHashMap := &HashMapOfValue{
		Table: newTable,
	}
	hashMap.version++

	for _, entry := range oldTable {
		if entry.Key().IsUndefined() {
			continue
		}

		i, err := HashMapOfValueIndex(vm, tmpHashMap, entry.Key())
		if !err.IsUndefined() {
			return err
		}
		if i == -1 {
			panic("no room in target hashmap during resizing")
		}
		newTable[i] = entry
		tmpHashMap.OccupiedSlots++
		tmpHashMap.Elements++
	}

	hashMap.OccupiedSlots = tmpHashMap.OccupiedSlots
	hashMap.Elements = tmpHashMap.Elements
	hashMap.Table = newTable
	return value.Undefined
}

func HashMapOfValueSetWithMaxLoad(vm *Thread, hashMap *HashMapOfValue, key, val value.Value, maxLoad float64) value.Value {
	if hashMap.Capacity() == 0 {
		HashMapOfValueSetCapacity(vm, hashMap, 5)
	} else if float64(hashMap.OccupiedSlots) >= float64(hashMap.Capacity())*maxLoad {
		HashMapOfValueSetCapacity(vm, hashMap, hashMap.OccupiedSlots*2)
	}

	index, err := HashMapOfValueIndex(vm, hashMap, key)
	if !err.IsUndefined() {
		return err
	}
	if index == -1 {
		panic(fmt.Sprintf("no room in target hashmap when trying to add a new key: %s", hashMap.Inspect()))
	}
	entry := hashMap.Table[index]
	if entry.Key().IsUndefined() && entry.Value().IsUndefined() {
		hashMap.OccupiedSlots++
		hashMap.Elements++
	}

	hashMap.Table[index] = value.MakePairOfValue(
		key,
		val,
	)
	hashMap.version++

	return value.Undefined
}

// Set a value under the given key.
func HashMapOfValueSet(vm *Thread, hashMap *HashMapOfValue, key, val value.Value) value.Value {
	return HashMapOfValueSetWithMaxLoad(vm, hashMap, key, val, HashMapOfValueMaxLoad)
}

// Get the index that the key should be inserted into.
// Returns (nil, err) when an error has been encountered.
// Returns (-1, nil) when there's no room for new values.
func HashMapOfValueIndex(vm *Thread, hashMap *HashMapOfValue, key value.Value) (int, value.Value) {
	hash, err := Hash(vm, key)
	if !err.IsUndefined() {
		return 0, err
	}
	deletedIndex := -1

	capacity := hashMap.Capacity()
	index := int(hash % value.UInt64(capacity))
	startIndex := index

	for {
		entry := hashMap.Table[index]
		if entry.Key().IsUndefined() {
			// empty or deleted entry

			if entry.Value().IsUndefined() {
				// empty bucket
				if deletedIndex != -1 {
					return deletedIndex, value.Undefined
				}
				return index, value.Undefined
			} else if deletedIndex == -1 {
				// deleted entry
				deletedIndex = index
			}
		} else {
			// present entry
			equal, err := Equal(vm, entry.Key(), key)
			if !err.IsUndefined() {
				return 0, err
			}
			if value.Truthy(equal) {
				return index, value.Undefined
			}
		}

		if index == capacity-1 {
			index = 0
		} else {
			index++
		}

		// when we reach the start index
		// all slots are checked
		if index == startIndex {
			return -1, value.Undefined
		}
	}
}
