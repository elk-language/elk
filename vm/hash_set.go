package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

// ::Std::HashSet
func init() {
	// Instance methods
	c := &value.HashSetClass.MethodContainer
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashSet)
			iterator := value.NewHashSetIterator(self)
			return iterator, nil
		},
	)
	Def(
		c,
		"capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashSet)
			return value.SmallInt(self.Capacity()), nil
		},
	)
	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashSet)
			return value.SmallInt(self.Length()), nil
		},
	)
	Def(
		c,
		"left_capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashSet)
			return value.SmallInt(self.LeftCapacity()), nil
		},
	)
	Def(
		c,
		"grow",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashSet)
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
			HashSetGrow(vm, self, n)
			return self, nil
		},
		DefWithParameters("new_slots"),
		DefWithSealed(),
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashSet)
			other, ok := args[1].(*value.HashSet)
			if !ok {
				return value.False, nil
			}
			equal, err := HashSetEqual(vm, self, other)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(equal), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashSet)
			val := args[1]
			contains, err := HashSetContains(vm, self, val)
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(contains), nil
		},
		DefWithParameters("value"),
		DefWithSealed(),
	)
	Def(
		c,
		"<<",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashSet)
			val := args[1]
			err := HashSetAppend(vm, self, val)
			if err != nil {
				return nil, err
			}
			return self, nil
		},
		DefWithParameters("value"),
		DefWithSealed(),
	)
	Alias(c, "append", "<<")
}

// ::Std::HashSet::Iterator
func init() {
	// Instance methods
	c := &value.HashSetIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashSetIterator)
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

func NewHashSetComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *value.HashSet) bool {
		if x == y {
			return true
		}
		if x.Length() != y.Length() {
			return false
		}
		if x.Capacity() != y.Capacity() {
			return false
		}

		v := New()
		for _, xVal := range x.Table {
			if xVal == nil || xVal == value.Undefined {
				continue
			}

			contains, err := HashSetContains(v, y, xVal)
			if err != nil {
				return false
			}
			if !contains {
				return false
			}
		}

		return true
	})
}

// Create a new hash set with the given entries.
func NewHashSetWithElements(vm *VM, elements ...value.Value) (*value.HashSet, value.Value) {
	return NewHashSetWithCapacityAndElements(vm, len(elements), elements...)
}

// Create a new hash set with the given entries.
func MustNewHashSetWithElements(vm *VM, elements ...value.Value) *value.HashSet {
	set, err := NewHashSetWithElements(vm, elements...)
	if err != nil {
		panic(err)
	}

	return set
}

func NewHashSetWithCapacityAndElements(vm *VM, capacity int, elements ...value.Value) (*value.HashSet, value.Value) {
	s := value.NewHashSet(capacity)
	for _, element := range elements {
		err := HashSetAppend(vm, s, element)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func NewHashSetWithCapacityAndElementsMaxLoad(vm *VM, capacity int, maxLoad float64, elements ...value.Value) (*value.HashSet, value.Value) {
	s := value.NewHashSet(capacity)
	for _, element := range elements {
		err := HashSetAppendWithMaxLoad(vm, s, element, maxLoad)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func MustNewHashSetWithCapacityAndElements(vm *VM, capacity int, elements ...value.Value) *value.HashSet {
	set, err := NewHashSetWithCapacityAndElements(vm, capacity, elements...)
	if err != nil {
		panic(err)
	}

	return set
}

func MustNewHashSetWithCapacityAndElementsMaxLoad(vm *VM, capacity int, maxLoad float64, elements ...value.Value) *value.HashSet {
	set, err := NewHashSetWithCapacityAndElementsMaxLoad(vm, capacity, maxLoad, elements...)
	if err != nil {
		panic(err)
	}

	return set
}

// Checks whether two hash sets are equal
func HashSetEqual(vm *VM, x *value.HashSet, y *value.HashSet) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, nil
	}

	for _, xVal := range x.Table {
		if xVal == nil || xVal == value.Undefined {
			continue
		}

		contains, err := HashSetContains(vm, y, xVal)
		if err != nil {
			return false, err
		}
		if !contains {
			return false, nil
		}
	}

	return true, nil
}

// Delete the given value from the hash set
func HashSetDelete(vm *VM, hashSet *value.HashSet, val value.Value) (bool, value.Value) {
	if hashSet.Length() == 0 {
		return false, nil
	}

	index, err := HashSetIndex(vm, hashSet, val)
	if err != nil {
		return false, err
	}
	if index < 0 {
		return false, nil
	}
	existingVal := hashSet.Table[index]
	if existingVal == nil || existingVal == value.Undefined {
		return false, nil
	}

	// undefined means that the entry has been deleted
	hashSet.Table[index] = value.Undefined
	hashSet.Elements--

	return true, nil
}

// Check whether the given value is contained within the set.
func HashSetContains(vm *VM, set *value.HashSet, val value.Value) (bool, value.Value) {
	if set.Length() == 0 {
		return false, nil
	}

	index, err := HashSetIndex(vm, set, val)
	if err != nil {
		return false, err
	}
	if index == -1 {
		return false, nil
	}

	valInSlot := set.Table[index]
	if valInSlot == nil || valInSlot == value.Undefined {
		return false, nil
	}
	return true, nil
}

func HashSetCopyTable(vm *VM, target *value.HashSet, source []value.Value) value.Value {
	for _, entry := range source {
		if entry == nil || entry == value.Undefined {
			continue
		}

		err := HashSetAppendWithMaxLoad(vm, target, entry, 1)
		if err != nil {
			return err
		}
	}

	return nil
}

// Copy the pairs of one hashmap to the other.
func HashSetCopy(vm *VM, target *value.HashSet, source *value.HashSet) value.Value {
	requiredCapacity := target.Length() + source.Length()
	if target.Capacity() < requiredCapacity {
		HashSetSetCapacity(vm, target, requiredCapacity)
	}

	for _, entry := range source.Table {
		if entry == nil || entry == value.Undefined {
			continue
		}

		i, err := HashSetIndex(vm, target, entry)
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
func HashSetGrow(vm *VM, set *value.HashSet, newSlots int) value.Value {
	return HashSetSetCapacity(vm, set, set.Capacity()+newSlots)
}

// Resize the given set to the desired capacity.
func HashSetSetCapacity(vm *VM, set *value.HashSet, capacity int) value.Value {
	if set.Capacity() == capacity {
		return nil
	}

	oldTable := set.Table
	newTable := make([]value.Value, capacity)
	tmpHashSet := &value.HashSet{
		Table: newTable,
	}

	for _, entry := range oldTable {
		if entry == nil || entry == value.Undefined {
			continue
		}

		i, err := HashSetIndex(vm, tmpHashSet, entry)
		if err != nil {
			return err
		}
		if i == -1 {
			panic("no room in target hashset during resizing")
		}
		newTable[i] = entry
		tmpHashSet.OccupiedSlots++
		tmpHashSet.Elements++
	}

	set.OccupiedSlots = tmpHashSet.OccupiedSlots
	set.Elements = tmpHashSet.Elements
	set.Table = newTable
	return nil
}

func HashSetAppendWithMaxLoad(vm *VM, set *value.HashSet, val value.Value, maxLoad float64) value.Value {
	if set.Capacity() == 0 {
		HashSetSetCapacity(vm, set, 5)
	} else if float64(set.OccupiedSlots) >= float64(set.Capacity())*maxLoad {
		HashSetSetCapacity(vm, set, set.OccupiedSlots*2)
	}

	index, err := HashSetIndex(vm, set, val)
	if err != nil {
		return err
	}
	if index == -1 {
		panic(fmt.Sprintf("no room in target hashset when trying to add a new value: %s", set.Inspect()))
	}
	entry := set.Table[index]
	if entry == nil {
		set.OccupiedSlots++
		set.Elements++
	}

	set.Table[index] = val

	return nil
}

// Set a value under the given key.
func HashSetAppend(vm *VM, hashMap *value.HashSet, val value.Value) value.Value {
	return HashSetAppendWithMaxLoad(vm, hashMap, val, value.HashSetMaxLoad)
}

// Get the index that the value should be inserted into.
// Returns (nil, err) when an error has been encountered.
// Returns (-1, nil) when there's no room for new values.
func HashSetIndex(vm *VM, set *value.HashSet, val value.Value) (int, value.Value) {
	hash, err := Hash(vm, val)
	if err != nil {
		return 0, err
	}
	deletedIndex := -1

	capacity := set.Capacity()
	index := int(hash % value.UInt64(capacity))
	startIndex := index

	for {
		entry := set.Table[index]
		if entry == nil {
			// empty bucket
			if deletedIndex != -1 {
				return deletedIndex, nil
			}
			return index, nil
		} else if entry == value.Undefined {
			if deletedIndex == -1 {
				// deleted entry
				deletedIndex = index
			}
		} else {
			// present entry
			equal, err := Equal(vm, entry, val)
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
