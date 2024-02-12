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
		"iterator",
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
		DefWithParameters("key"),
		DefWithSealed(),
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
		DefWithParameters("key", "value"),
		DefWithSealed(),
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
		DefWithParameters("new_slots"),
		DefWithSealed(),
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
		"iterator",
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
	var reachedEnd bool

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
			equal, err := StrictEqual(vm, entry.Key, key)
			if err != nil {
				return 0, err
			}
			if value.Truthy(equal) {
				return index, nil
			}
		}

		if index == capacity-1 {
			if reachedEnd {
				return -1, nil
			}
			index = 0
			reachedEnd = true
		} else {
			index++
		}
	}
}
