package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

// ::Std::HashMap
func initHashMap() {
	// Instance methods
	c := &value.HashMapClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashMap)
			iterator := value.NewHashMapIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)
	Def(
		c,
		"capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashMap)
			return value.SmallInt(self.Capacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashMap)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"left_capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashMap)
			return value.SmallInt(self.LeftCapacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashMap)
			key := args[1]
			result, err := HashMapGet(vm, self, key)
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
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashMap)
			key := args[1]
			val := args[2]
			err := HashMapSet(vm, self, key, val)
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
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashMap)
			other := args[1]

			switch o := other.SafeAsReference().(type) {
			case *value.HashMap:
				result, err := HashMapConcat(vm, self, o)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.Ref(result), value.Undefined
			case *value.HashRecord:
				result, err := HashMapConcat(vm, self, (*value.HashMap)(o))
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.Ref(result), value.Undefined
			default:
				return value.Undefined, value.Ref(value.NewCoerceError(value.HashMapClass, other.Class()))
			}
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"grow",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashMap)
			nValue := args[1]
			n, ok := value.IntToGoInt(nValue)
			if !ok && n == -1 {
				return value.Undefined, value.Ref(value.NewTooLargeCapacityError(nValue.Inspect()))
			}
			if n < 0 {
				return value.Undefined, value.Ref(value.NewNegativeCapacityError(nValue.Inspect()))
			}
			if !ok {
				return value.Undefined, value.Ref(value.NewCapacityTypeError(nValue.Inspect()))
			}
			HashMapGrow(vm, self, n)
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashMap)
			otherVal := args[1]
			switch other := otherVal.SafeAsReference().(type) {
			case *value.Pair:
				contains, err := HashMapContains(vm, self, other)
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
			self := args[0].MustReference().(*value.HashMap)
			contains, err := HashMapContainsKey(vm, self, args[1])
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
			self := args[0].MustReference().(*value.HashMap)
			contains, err := HashMapContainsValue(vm, self, args[1])
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
			self := args[0].MustReference().(*value.HashMap)
			other, ok := args[1].SafeAsReference().(*value.HashMap)
			if !ok {
				return value.False, value.Undefined
			}
			equal, err := HashMapEqual(vm, self, other)
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
			self := args[0].MustReference().(*value.HashMap)
			switch other := args[1].SafeAsReference().(type) {
			case *value.HashMap:
				equal, err := HashMapLaxEqual(vm, self, other)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.ToElkBool(equal), value.Undefined
			case *value.HashRecord:
				equal, err := HashMapLaxEqual(vm, self, (*value.HashMap)(other))
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
			self := args[0].MustReference().(*value.HashMap)
			callable := args[1]
			newMap := value.NewHashMap(self.Length())

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
					err = HashMapSet(vm, newMap, r.Key, r.Value)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newMap), value.Undefined
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
				err = HashMapSet(vm, newMap, r.Key, r.Value)
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
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashMap)
			callable := args[1]
			newMap := value.NewHashMap(self.Length())

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
					err = HashMapSet(vm, newMap, pair.Key, result)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newMap), value.Undefined
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
				err = HashMapSet(vm, newMap, pair.Key, result)
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
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashMap)
			callable := args[1]

			// callable is a closure
			if function, ok := callable.SafeAsReference().(*Closure); ok {
				for i := range len(self.Table) {
					pair := self.Table[i]
					if pair.Key.IsUndefined() {
						continue
					}
					result, err := vm.CallClosure(function, pair.Value)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					self.Table[i].Value = result
				}
				return value.Ref(self), value.Undefined
			}

			// callable is another value
			for i := range len(self.Table) {
				pair := self.Table[i]
				if pair.Key.IsUndefined() {
					continue
				}
				result, err := vm.CallMethodByName(callSymbol, callable, pair.Value)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				self.Table[i].Value = result
			}
			return value.Ref(self), value.Undefined
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
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashMapIterator)
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
			if xPair.Key.IsUndefined() {
				continue
			}

			yVal, err := HashMapGet(v, y, xPair.Key)
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
func NewHashMapWithElements(vm *VM, elements ...value.Pair) (*value.HashMap, value.Value) {
	return NewHashMapWithCapacityAndElements(vm, len(elements), elements...)
}

// Create a new hashmap with the given entries.
func MustNewHashMapWithElements(vm *VM, elements ...value.Pair) *value.HashMap {
	hmap, err := NewHashMapWithElements(vm, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return hmap
}

func NewHashMapWithCapacityAndElements(vm *VM, capacity int, elements ...value.Pair) (*value.HashMap, value.Value) {
	h := value.NewHashMap(capacity)
	for _, element := range elements {
		err := HashMapSet(vm, h, element.Key, element.Value)
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return h, value.Undefined
}

func MustNewHashMapWithCapacityAndElements(vm *VM, capacity int, elements ...value.Pair) *value.HashMap {
	hmap, err := NewHashMapWithCapacityAndElements(vm, capacity, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return hmap
}

// Checks whether two hash maps are equal (lax)
func HashMapLaxEqual(vm *VM, x *value.HashMap, y *value.HashMap) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for _, xPair := range x.Table {
		if xPair.Key.IsUndefined() {
			continue
		}

		yVal, err := HashMapGet(vm, y, xPair.Key)
		if !err.IsUndefined() {
			return false, err
		}
		if yVal.IsUndefined() {
			return false, value.Undefined
		}
		eqVal, err := LaxEqual(vm, xPair.Value, yVal)
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
func HashMapConcat(vm *VM, x *value.HashMap, y *value.HashMap) (*value.HashMap, value.Value) {
	result := x.Clone()
	err := HashMapCopy(vm, result, y)
	if !err.IsUndefined() {
		return nil, err
	}
	return result, value.Undefined
}

// Checks whether two hash maps are equal
func HashMapEqual(vm *VM, x *value.HashMap, y *value.HashMap) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for _, xPair := range x.Table {
		if xPair.Key.IsUndefined() {
			continue
		}

		yVal, err := HashMapGet(vm, y, xPair.Key)
		if !err.IsUndefined() {
			return false, err
		}
		if yVal.IsUndefined() {
			return false, value.Undefined
		}
		eqVal, err := Equal(vm, xPair.Value, yVal)
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
func HashMapDelete(vm *VM, hashMap *value.HashMap, key value.Value) (bool, value.Value) {
	if hashMap.Length() == 0 {
		return false, value.Undefined
	}

	index, err := HashMapIndex(vm, hashMap, key)
	if !err.IsUndefined() {
		return false, err
	}
	if index < 0 {
		return false, value.Undefined
	}
	if hashMap.Table[index].Key.IsUndefined() {
		return false, value.Undefined
	}

	hashMap.Table[index] = value.Pair{
		Key:   value.Undefined,
		Value: value.True,
	}
	hashMap.Elements--

	return true, value.Undefined
}

// Get the element under the given key.
// Returns (value, undefined) when the value has been found.
// Returns (undefined, undefine) when the key is not present.
// Returns (undefined, err) when there was an error.
func HashMapGet(vm *VM, hashMap *value.HashMap, key value.Value) (value.Value, value.Value) {
	if hashMap.Length() == 0 {
		return value.Undefined, value.Undefined
	}

	index, err := HashMapIndex(vm, hashMap, key)
	if !err.IsUndefined() {
		return value.Undefined, err
	}
	if index == -1 {
		return value.Undefined, value.Undefined
	}

	return hashMap.Table[index].Value, value.Undefined
}

// Check if the given pair is present in the map
func HashMapContains(vm *VM, hashMap *value.HashMap, pair *value.Pair) (bool, value.Value) {
	val, err := HashMapGet(vm, hashMap, pair.Key)
	if !err.IsUndefined() {
		return false, err
	}
	if val.IsUndefined() {
		return false, value.Undefined
	}

	equal, err := Equal(vm, val, pair.Value)
	if !err.IsUndefined() {
		return false, err
	}

	return value.Truthy(equal), value.Undefined
}

// Check if the given key is present in the map
func HashMapContainsKey(vm *VM, hashMap *value.HashMap, key value.Value) (bool, value.Value) {
	if hashMap.Length() == 0 {
		return false, value.Undefined
	}

	index, err := HashMapIndex(vm, hashMap, key)
	if !err.IsUndefined() {
		return false, err
	}
	if index == -1 {
		return false, value.Undefined
	}

	pair := hashMap.Table[index]
	if pair.Key.IsUndefined() {
		return false, value.Undefined
	}

	return true, value.Undefined
}

// Check if the given value is present in the map
func HashMapContainsValue(vm *VM, hashMap *value.HashMap, val value.Value) (bool, value.Value) {
	for _, pair := range hashMap.Table {
		// if the Key is undefined the entry is empty or deleted
		// so we skip it
		if pair.Key.IsUndefined() {
			continue
		}

		equal, err := Equal(vm, pair.Value, val)
		if !err.IsUndefined() {
			return false, err
		}

		if value.Truthy(equal) {
			return true, value.Undefined
		}
	}

	return false, value.Undefined
}

func HashMapCopyTable(vm *VM, target *value.HashMap, source []value.Pair) value.Value {
	for _, entry := range source {
		if entry.Key.IsUndefined() {
			continue
		}

		err := HashMapSetWithMaxLoad(vm, target, entry.Key, entry.Value, 1)
		if !err.IsUndefined() {
			return err
		}
	}

	return value.Undefined
}

// Copy the pairs of one hashmap to the other.
func HashMapCopy(vm *VM, target *value.HashMap, source *value.HashMap) value.Value {
	requiredCapacity := target.Length() + source.Length()
	if target.Capacity() < requiredCapacity {
		HashMapSetCapacity(vm, target, requiredCapacity)
	}

	for _, entry := range source.Table {
		if entry.Key.IsUndefined() {
			continue
		}

		i, err := HashMapIndex(vm, target, entry.Key)
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
func HashMapGrow(vm *VM, hashMap *value.HashMap, newSlots int) value.Value {
	return HashMapSetCapacity(vm, hashMap, hashMap.Capacity()+newSlots)
}

// Resize the given hashmap to the desired capacity.
func HashMapSetCapacity(vm *VM, hashMap *value.HashMap, capacity int) value.Value {
	if hashMap.Capacity() == capacity {
		return value.Undefined
	}

	oldTable := hashMap.Table
	newTable := make([]value.Pair, capacity)
	tmpHashMap := &value.HashMap{
		Table: newTable,
	}

	for _, entry := range oldTable {
		if entry.Key.IsUndefined() {
			continue
		}

		i, err := HashMapIndex(vm, tmpHashMap, entry.Key)
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

func HashMapSetWithMaxLoad(vm *VM, hashMap *value.HashMap, key, val value.Value, maxLoad float64) value.Value {
	if hashMap.Capacity() == 0 {
		HashMapSetCapacity(vm, hashMap, 5)
	} else if float64(hashMap.OccupiedSlots) >= float64(hashMap.Capacity())*maxLoad {
		HashMapSetCapacity(vm, hashMap, hashMap.OccupiedSlots*2)
	}

	index, err := HashMapIndex(vm, hashMap, key)
	if !err.IsUndefined() {
		return err
	}
	if index == -1 {
		panic(fmt.Sprintf("no room in target hashmap when trying to add a new key: %s", hashMap.Inspect()))
	}
	entry := hashMap.Table[index]
	if entry.Key.IsUndefined() && entry.Value.IsUndefined() {
		hashMap.OccupiedSlots++
		hashMap.Elements++
	}

	hashMap.Table[index] = value.Pair{
		Key:   key,
		Value: val,
	}

	return value.Undefined
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
	if !err.IsUndefined() {
		return 0, err
	}
	deletedIndex := -1

	capacity := hashMap.Capacity()
	index := int(hash % value.UInt64(capacity))
	startIndex := index

	for {
		entry := hashMap.Table[index]
		if entry.Key.IsUndefined() {
			// empty or deleted entry

			if entry.Value.IsUndefined() {
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
			equal, err := Equal(vm, entry.Key, key)
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
