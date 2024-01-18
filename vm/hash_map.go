package vm

import (
	"github.com/elk-language/elk/value"
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

func HashMapDelete(vm *VM, hashMap *value.HashMap, key value.Value) (bool, value.Value) {
	if hashMap.Count == 0 {
		return false, nil
	}

	index, err := HashMapIndex(vm, hashMap, key)
	if err != nil {
		return false, err
	}
	if hashMap.Table[index].Key == nil {
		return false, nil
	}

	hashMap.Table[index] = value.Pair{
		Key:   nil,
		Value: value.True,
	}

	return true, nil
}

func HashMapGet(vm *VM, hashMap *value.HashMap, key, val value.Value) (value.Value, value.Value) {
	if hashMap.Count == 0 {
		return nil, nil
	}

	index, err := HashMapIndex(vm, hashMap, key)
	if err != nil {
		return nil, err
	}

	return hashMap.Table[index].Value, nil
}

func HashMapCopy(vm *VM, target *value.HashMap, source *value.HashMap) value.Value {
	requiredCapacity := target.Length() + source.Length()
	if target.Capacity() < requiredCapacity {
		HashMapGrow(vm, target, requiredCapacity)
	}

	for _, entry := range source.Table {
		if entry.Key == nil {
			continue
		}

		i, err := HashMapIndex(vm, target, entry.Key)
		if err != nil {
			return err
		}
		target.Table[i] = entry
		target.Count++
	}

	return nil
}

func HashMapGrow(vm *VM, hashMap *value.HashMap, capacity int) value.Value {
	newTable := make([]value.Pair, capacity)

	hashMap.Count = 0
	for _, entry := range hashMap.Table {
		if entry.Key == nil {
			continue
		}

		i, err := HashMapIndex(vm, hashMap, entry.Key)
		if err != nil {
			return err
		}
		newTable[i] = entry
		hashMap.Count++
	}

	hashMap.Table = newTable
	return nil
}

func HashMapSet(vm *VM, hashMap *value.HashMap, key, val value.Value) value.Value {
	if float64(hashMap.Count) >= float64(hashMap.Capacity())*value.HashMapMaxLoad {
		HashMapGrow(vm, hashMap, hashMap.Count*2)
	}

	index, err := HashMapIndex(vm, hashMap, key)
	if err != nil {
		return err
	}
	entry := hashMap.Table[index]
	if entry.Key == nil && entry.Value == nil {
		hashMap.Count++
	}

	hashMap.Table[index] = value.Pair{
		Key:   key,
		Value: val,
	}

	return nil
}

// Get the index that the key should be inserted into
func HashMapIndex(vm *VM, hashMap *value.HashMap, key value.Value) (int, value.Value) {
	hash, err := Hash(vm, key)
	if err != nil {
		return 0, err
	}
	deletedIndex := -1

	capacity := hashMap.Capacity()
	index := int(hash % value.UInt64(capacity))
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

		index = (index + 1) % capacity
	}
}
