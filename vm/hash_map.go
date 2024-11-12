package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

// ::Std::HashMap
func init() {
	// Instance methods
	c := &value.HashMapClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			iterator := value.NewHashMapIterator(self)
			return iterator, nil
		},
	)
	Def(
		c,
		"capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			return value.SmallInt(self.Capacity()), nil
		},
	)
	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			return value.SmallInt(self.Length()), nil
		},
	)
	Def(
		c,
		"left_capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			return value.SmallInt(self.LeftCapacity()), nil
		},
	)
	Def(
		c,
		"[]",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			key := args[1]
			result, err := HashMapGet(vm, self, key)
			if err != nil {
				return nil, err
			}
			if result == nil {
				return value.Nil, nil
			}
			return result, nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"[]=",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			key := args[1]
			val := args[2]
			err := HashMapSet(vm, self, key, val)
			if err != nil {
				return nil, err
			}
			return val, nil
		},
		DefWithParameters(2),
	)
	Def(
		c,
		"+",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			other := args[1]

			switch o := other.(type) {
			case *value.HashMap:
				result, err := HashMapConcat(vm, self, o)
				if err != nil {
					return nil, err
				}
				return result, nil
			case *value.HashRecord:
				result, err := HashMapConcat(vm, self, (*value.HashMap)(o))
				if err != nil {
					return nil, err
				}
				return result, nil
			default:
				return nil, value.NewCoerceError(value.HashMapClass, other.Class())
			}
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"grow",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			nValue := args[1]
			n, ok := value.IntToGoInt(nValue)
			if !ok && n == -1 {
				return nil, value.NewTooLargeCapacityError(nValue.Inspect())
			}
			if n < 0 {
				return nil, value.NewNegativeCapacityError(nValue.Inspect())
			}
			if !ok {
				return nil, value.NewCapacityTypeError(nValue.Inspect())
			}
			HashMapGrow(vm, self, n)
			return self, nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			otherVal := args[1]
			switch other := otherVal.(type) {
			case *value.Pair:
				contains, err := HashMapContains(vm, self, other)
				if err != nil {
					return nil, err
				}

				return value.ToElkBool(contains), nil
			default:
				return nil, value.NewCoerceError(value.PairClass, otherVal.Class())
			}
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains_key",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			contains, err := HashMapContainsKey(vm, self, args[1])
			if err != nil {
				return nil, err
			}

			return value.ToElkBool(contains), nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains_value",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			contains, err := HashMapContainsValue(vm, self, args[1])
			if err != nil {
				return nil, err
			}

			return value.ToElkBool(contains), nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			other, ok := args[1].(*value.HashMap)
			if !ok {
				return value.False, nil
			}
			equal, err := HashMapEqual(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(equal), nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			switch other := args[1].(type) {
			case *value.HashMap:
				equal, err := HashMapLaxEqual(vm, self, other)
				if err != nil {
					return nil, err
				}
				return value.ToElkBool(equal), nil
			case *value.HashRecord:
				equal, err := HashMapLaxEqual(vm, self, (*value.HashMap)(other))
				if err != nil {
					return nil, err
				}
				return value.ToElkBool(equal), nil
			default:
				return value.False, nil
			}
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			callable := args[1]
			newMap := value.NewHashMap(self.Length())

			// callable is a closure
			if function, ok := callable.(*Closure); ok {
				for i, pair := range self.Table {
					if pair.Key == nil {
						continue
					}
					result, err := vm.CallClosure(function, &self.Table[i])
					if err != nil {
						return nil, err
					}
					r, ok := result.(*value.Pair)
					if !ok {
						return nil, value.NewArgumentTypeError("pair", result.Class().Name, value.PairClass.Name)
					}
					err = HashMapSet(vm, newMap, r.Key, r.Value)
					if err != nil {
						return nil, err
					}
				}
				return newMap, nil
			}

			// callable is another value
			for i, pair := range self.Table {
				if pair.Key == nil {
					continue
				}
				result, err := vm.CallMethodByName(callSymbol, callable, &self.Table[i])
				if err != nil {
					return nil, err
				}
				r, ok := result.(*value.Pair)
				if !ok {
					return nil, value.NewArgumentTypeError("pair", result.Class().Name, value.PairClass.Name)
				}
				err = HashMapSet(vm, newMap, r.Key, r.Value)
				if err != nil {
					return nil, err
				}
			}
			return newMap, nil
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map_values",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			callable := args[1]
			newMap := value.NewHashMap(self.Length())

			// callable is a closure
			if function, ok := callable.(*Closure); ok {
				for _, pair := range self.Table {
					if pair.Key == nil {
						continue
					}
					result, err := vm.CallClosure(function, pair.Value)
					if err != nil {
						return nil, err
					}
					err = HashMapSet(vm, newMap, pair.Key, result)
					if err != nil {
						return nil, err
					}
				}
				return newMap, nil
			}

			// callable is another value
			for _, pair := range self.Table {
				if pair.Key == nil {
					continue
				}
				result, err := vm.CallMethodByName(callSymbol, callable, pair.Value)
				if err != nil {
					return nil, err
				}
				err = HashMapSet(vm, newMap, pair.Key, result)
				if err != nil {
					return nil, err
				}
			}
			return newMap, nil
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map_values_mut",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMap)
			callable := args[1]

			// callable is a closure
			if function, ok := callable.(*Closure); ok {
				for i := range len(self.Table) {
					pair := self.Table[i]
					if pair.Key == nil {
						continue
					}
					result, err := vm.CallClosure(function, pair.Value)
					if err != nil {
						return nil, err
					}
					self.Table[i].Value = result
				}
				return self, nil
			}

			// callable is another value
			for i := range len(self.Table) {
				pair := self.Table[i]
				if pair.Key == nil {
					continue
				}
				result, err := vm.CallMethodByName(callSymbol, callable, pair.Value)
				if err != nil {
					return nil, err
				}
				self.Table[i].Value = result
			}
			return self, nil
		},
		DefWithParameters(1),
	)
}

// ::Std::HashMap::Iterator
func init() {
	// Instance methods
	c := &value.HashMapIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashMapIterator)
			return self.Next()
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], nil
		},
	)

}

func NewHashMapComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *value.HashMap) bool {
		if x == y {
			return true
		}
		if x.Length() != y.Length() {
			return false
		}

		v := New()
		for _, xPair := range x.Table {
			if xPair.Key == nil {
				continue
			}

			yVal, err := HashMapGet(v, y, xPair.Key)
			if err != nil {
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
func NewHashMapWithElements(vm *VM, elements ...value.Pair) (*value.HashMap, value.Value) {
	return NewHashMapWithCapacityAndElements(vm, len(elements), elements...)
}

// Create a new hashmap with the given entries.
func MustNewHashMapWithElements(vm *VM, elements ...value.Pair) *value.HashMap {
	hmap, err := NewHashMapWithElements(vm, elements...)
	if err != nil {
		panic(err)
	}

	return hmap
}

func NewHashMapWithCapacityAndElements(vm *VM, capacity int, elements ...value.Pair) (*value.HashMap, value.Value) {
	h := value.NewHashMap(capacity)
	for _, element := range elements {
		err := HashMapSet(vm, h, element.Key, element.Value)
		if err != nil {
			return nil, err
		}
	}

	return h, nil
}

func MustNewHashMapWithCapacityAndElements(vm *VM, capacity int, elements ...value.Pair) *value.HashMap {
	hmap, err := NewHashMapWithCapacityAndElements(vm, capacity, elements...)
	if err != nil {
		panic(err)
	}

	return hmap
}

// Checks whether two hash maps are equal (lax)
func HashMapLaxEqual(vm *VM, x *value.HashMap, y *value.HashMap) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, nil
	}

	for _, xPair := range x.Table {
		if xPair.Key == nil {
			continue
		}

		yVal, err := HashMapGet(vm, y, xPair.Key)
		if err != nil {
			return false, err
		}
		if yVal == nil {
			return false, nil
		}
		eqVal, err := LaxEqual(vm, xPair.Value, yVal)
		if err != nil {
			return false, err
		}
		equal := value.Truthy(eqVal)
		if !equal {
			return false, nil
		}
	}

	return true, nil
}

// Create a new map containing the pairs of both maps.
func HashMapConcat(vm *VM, x *value.HashMap, y *value.HashMap) (*value.HashMap, value.Value) {
	result := x.Clone()
	err := HashMapCopy(vm, result, y)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Checks whether two hash maps are equal
func HashMapEqual(vm *VM, x *value.HashMap, y *value.HashMap) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, nil
	}

	for _, xPair := range x.Table {
		if xPair.Key == nil {
			continue
		}

		yVal, err := HashMapGet(vm, y, xPair.Key)
		if err != nil {
			return false, err
		}
		if yVal == nil {
			return false, nil
		}
		eqVal, err := Equal(vm, xPair.Value, yVal)
		if err != nil {
			return false, err
		}
		equal := value.Truthy(eqVal)
		if !equal {
			return false, nil
		}
	}

	return true, nil
}

// Delete the given key from the hashMap
func HashMapDelete(vm *VM, hashMap *value.HashMap, key value.Value) (bool, value.Value) {
	if hashMap.Length() == 0 {
		return false, nil
	}

	index, err := HashMapIndex(vm, hashMap, key)
	if err != nil {
		return false, err
	}
	if index < 0 {
		return false, nil
	}
	if hashMap.Table[index].Key == nil {
		return false, nil
	}

	hashMap.Table[index] = value.Pair{
		Key:   nil,
		Value: value.True,
	}
	hashMap.Elements--

	return true, nil
}

// Get the element under the given key.
// Returns (value, nil) when the value has been found.
// Returns (nil, nil) when the key is not present.
// Returns (nil, err) when there was an error.
func HashMapGet(vm *VM, hashMap *value.HashMap, key value.Value) (value.Value, value.Value) {
	if hashMap.Length() == 0 {
		return nil, nil
	}

	index, err := HashMapIndex(vm, hashMap, key)
	if err != nil {
		return nil, err
	}
	if index == -1 {
		return nil, nil
	}

	return hashMap.Table[index].Value, nil
}

// Check if the given pair is present in the map
func HashMapContains(vm *VM, hashMap *value.HashMap, pair *value.Pair) (bool, value.Value) {
	val, err := HashMapGet(vm, hashMap, pair.Key)
	if err != nil {
		return false, err
	}
	if val == nil {
		return false, nil
	}

	equal, err := Equal(vm, val, pair.Value)
	if err != nil {
		return false, err
	}

	return value.Truthy(equal), nil
}

// Check if the given key is present in the map
func HashMapContainsKey(vm *VM, hashMap *value.HashMap, key value.Value) (bool, value.Value) {
	if hashMap.Length() == 0 {
		return false, nil
	}

	index, err := HashMapIndex(vm, hashMap, key)
	if err != nil {
		return false, err
	}
	if index == -1 {
		return false, nil
	}

	pair := hashMap.Table[index]
	if pair.Key == nil {
		return false, nil
	}

	return true, nil
}

// Check if the given value is present in the map
func HashMapContainsValue(vm *VM, hashMap *value.HashMap, val value.Value) (bool, value.Value) {
	for _, pair := range hashMap.Table {
		// if the Key is nil the entry is empty or deleted
		// so we skip it
		if pair.Key == nil {
			continue
		}

		equal, err := Equal(vm, pair.Value, val)
		if err != nil {
			return false, err
		}

		if value.Truthy(equal) {
			return true, nil
		}
	}

	return false, nil
}

func HashMapCopyTable(vm *VM, target *value.HashMap, source []value.Pair) value.Value {
	for _, entry := range source {
		if entry.Key == nil {
			continue
		}

		err := HashMapSetWithMaxLoad(vm, target, entry.Key, entry.Value, 1)
		if err != nil {
			return err
		}
	}

	return nil
}

// Copy the pairs of one hashmap to the other.
func HashMapCopy(vm *VM, target *value.HashMap, source *value.HashMap) value.Value {
	requiredCapacity := target.Length() + source.Length()
	if target.Capacity() < requiredCapacity {
		HashMapSetCapacity(vm, target, requiredCapacity)
	}

	for _, entry := range source.Table {
		if entry.Key == nil {
			continue
		}

		i, err := HashMapIndex(vm, target, entry.Key)
		if err != nil {
			return err
		}
		if i == -1 {
			panic("no room in target hashmap during copy")
		}
		target.Table[i] = entry
		target.OccupiedSlots++
		target.Elements++
	}

	return nil
}

// Add additional n empty slots for new elements.
func HashMapGrow(vm *VM, hashMap *value.HashMap, newSlots int) value.Value {
	return HashMapSetCapacity(vm, hashMap, hashMap.Capacity()+newSlots)
}

// Resize the given hashmap to the desired capacity.
func HashMapSetCapacity(vm *VM, hashMap *value.HashMap, capacity int) value.Value {
	if hashMap.Capacity() == capacity {
		return nil
	}

	oldTable := hashMap.Table
	newTable := make([]value.Pair, capacity)
	tmpHashMap := &value.HashMap{
		Table: newTable,
	}

	for _, entry := range oldTable {
		if entry.Key == nil {
			continue
		}

		i, err := HashMapIndex(vm, tmpHashMap, entry.Key)
		if err != nil {
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
	return nil
}

func HashMapSetWithMaxLoad(vm *VM, hashMap *value.HashMap, key, val value.Value, maxLoad float64) value.Value {
	if hashMap.Capacity() == 0 {
		HashMapSetCapacity(vm, hashMap, 5)
	} else if float64(hashMap.OccupiedSlots) >= float64(hashMap.Capacity())*maxLoad {
		HashMapSetCapacity(vm, hashMap, hashMap.OccupiedSlots*2)
	}

	index, err := HashMapIndex(vm, hashMap, key)
	if err != nil {
		return err
	}
	if index == -1 {
		panic(fmt.Sprintf("no room in target hashmap when trying to add a new key: %s", hashMap.Inspect()))
	}
	entry := hashMap.Table[index]
	if entry.Key == nil && entry.Value == nil {
		hashMap.OccupiedSlots++
		hashMap.Elements++
	}

	hashMap.Table[index] = value.Pair{
		Key:   key,
		Value: val,
	}

	return nil
}

// Set a value under the given key.
func HashMapSet(vm *VM, hashMap *value.HashMap, key, val value.Value) value.Value {
	return HashMapSetWithMaxLoad(vm, hashMap, key, val, value.HashMapMaxLoad)
}

// Get the index that the key should be inserted into.
// Returns (nil, err) when an error has been encountered.
// Returns (-1, nil) when there's no room for new values.
func HashMapIndex(vm *VM, hashMap *value.HashMap, key value.Value) (int, value.Value) {
	hash, err := Hash(vm, key)
	if err != nil {
		return 0, err
	}
	deletedIndex := -1

	capacity := hashMap.Capacity()
	index := int(hash % value.UInt64(capacity))
	startIndex := index

	for {
		entry := hashMap.Table[index]
		if entry.Key == nil {
			// empty or deleted entry

			if entry.Value == nil {
				// empty bucket
				if deletedIndex != -1 {
					return deletedIndex, nil
				}
				return index, nil
			} else if deletedIndex == -1 {
				// deleted entry
				deletedIndex = index
			}
		} else {
			// present entry
			equal, err := Equal(vm, entry.Key, key)
			if err != nil {
				return 0, err
			}
			if value.Truthy(equal) {
				return index, nil
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
			return -1, nil
		}
	}
}
