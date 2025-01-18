package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

// ::Std::HashSet
func initHashSet() {
	// Instance methods
	c := &value.HashSetClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashSet)
			iterator := value.NewHashSetIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)
	Def(
		c,
		"capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashSet)
			return value.SmallInt(self.Capacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashSet)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"left_capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashSet)
			return value.SmallInt(self.LeftCapacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"grow",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashSet)
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
			HashSetGrow(vm, self, n)
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashSet)
			otherVal := args[1]
			other, ok := otherVal.SafeAsReference().(*value.HashSet)
			if !ok {
				return value.Undefined, value.Ref(value.NewCoerceError(value.HashSetClass, otherVal.Class()))
			}
			result, err := HashSetUnion(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.Ref(result), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "union", "+")
	Alias(c, "|", "+")

	Def(
		c,
		"&",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashSet)
			otherVal := args[1]
			other, ok := otherVal.SafeAsReference().(*value.HashSet)
			if !ok {
				return value.Undefined, value.Ref(value.NewCoerceError(value.HashSetClass, otherVal.Class()))
			}
			result, err := HashSetIntersection(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.Ref(result), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "intersection", "&")

	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashSet)
			other, ok := args[1].SafeAsReference().(*value.HashSet)
			if !ok {
				return value.False, value.Undefined
			}
			equal, err := HashSetEqual(vm, self, other)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.ToElkBool(equal), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "=~", "==")

	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashSet)
			val := args[1]
			contains, err := HashSetContains(vm, self, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.ToElkBool(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"append",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashSet)
			val := args[1].SafeAsReference().(*value.ArrayList)
			for _, element := range *val {
				err := HashSetAppend(vm, self, element)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
			}
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<<",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashSet)
			val := args[1]
			err := HashSetAppend(vm, self, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.HashSet)
			callable := args[1]
			newSet := value.NewHashSet(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(*Closure); ok {
				for _, val := range self.Table {
					if val == DeletedHashSetValue || val.IsUndefined() {
						continue
					}
					result, err := vm.CallClosure(function, val)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					err = HashSetAppend(vm, newSet, result)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newSet), value.Undefined
			}

			// callable is another value
			for _, val := range self.Table {
				if val == DeletedHashSetValue || val.IsUndefined() {
					continue
				}
				result, err := vm.CallMethodByName(callSymbol, callable, val)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				err = HashSetAppend(vm, newSet, result)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
			}
			return value.Ref(newSet), value.Undefined
		},
		DefWithParameters(1),
	)
}

// ::Std::HashSet::Iterator
func initHashSetIterator() {
	// Instance methods
	c := &value.HashSetIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.HashSetIterator)(args[0].Pointer())
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
			self := (*value.HashSetIterator)(args[0].Pointer())
			self.Reset()
			return args[0], value.Undefined
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
			if xVal == DeletedHashSetValue || xVal.IsUndefined() {
				continue
			}

			contains, err := HashSetContains(v, y, xVal)
			if !err.IsUndefined() {
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
	if !err.IsUndefined() {
		panic(err)
	}

	return set
}

func NewHashSetWithCapacityAndElements(vm *VM, capacity int, elements ...value.Value) (*value.HashSet, value.Value) {
	s := value.NewHashSet(capacity)
	for _, element := range elements {
		err := HashSetAppend(vm, s, element)
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return s, value.Undefined
}

func NewHashSetWithCapacityAndElementsMaxLoad(vm *VM, capacity int, maxLoad float64, elements ...value.Value) (*value.HashSet, value.Value) {
	s := value.NewHashSet(capacity)
	for _, element := range elements {
		err := HashSetAppendWithMaxLoad(vm, s, element, maxLoad)
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return s, value.Undefined
}

func MustNewHashSetWithCapacityAndElements(vm *VM, capacity int, elements ...value.Value) *value.HashSet {
	set, err := NewHashSetWithCapacityAndElements(vm, capacity, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return set
}

func MustNewHashSetWithCapacityAndElementsMaxLoad(vm *VM, capacity int, maxLoad float64, elements ...value.Value) *value.HashSet {
	set, err := NewHashSetWithCapacityAndElementsMaxLoad(vm, capacity, maxLoad, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return set
}

// Checks whether two hash sets are equal
func HashSetEqual(vm *VM, x *value.HashSet, y *value.HashSet) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for _, xVal := range x.Table {
		if xVal == DeletedHashSetValue || xVal.IsUndefined() {
			continue
		}

		contains, err := HashSetContains(vm, y, xVal)
		if !err.IsUndefined() {
			return false, err
		}
		if !contains {
			return false, value.Undefined
		}
	}

	return true, value.Undefined
}

// Create a new set that is the union of the given two sets
func HashSetUnion(vm *VM, x *value.HashSet, y *value.HashSet) (*value.HashSet, value.Value) {
	var longer *value.HashSet
	var shorter *value.HashSet
	if x.Length() > y.Length() {
		longer = x
		shorter = y
	} else {
		longer = y
		shorter = x
	}

	newSet := value.NewHashSet(longer.Elements)
	HashSetCopy(vm, newSet, longer)
	for _, shorterVal := range shorter.Table {
		if shorterVal == DeletedHashSetValue || shorterVal.IsUndefined() {
			continue
		}

		err := HashSetAppend(vm, newSet, shorterVal)
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return newSet, value.Undefined
}

// Create a new set that is the intersection of the given two sets
func HashSetIntersection(vm *VM, x *value.HashSet, y *value.HashSet) (*value.HashSet, value.Value) {
	var longer *value.HashSet
	var shorter *value.HashSet
	if x.Length() > y.Length() {
		longer = x
		shorter = y
	} else {
		longer = y
		shorter = x
	}

	newSet := value.NewHashSet(5)
	for _, shorterVal := range shorter.Table {
		if shorterVal == DeletedHashSetValue || shorterVal.IsUndefined() {
			continue
		}

		contains, err := HashSetContains(vm, longer, shorterVal)
		if !err.IsUndefined() {
			return nil, err
		}
		if contains {
			HashSetAppend(vm, newSet, shorterVal)
		}
	}

	return newSet, value.Undefined
}

// Delete the given value from the hash set
func HashSetDelete(vm *VM, hashSet *value.HashSet, val value.Value) (bool, value.Value) {
	if hashSet.Length() == 0 {
		return false, value.Undefined
	}

	index, err := HashSetIndex(vm, hashSet, val)
	if !err.IsUndefined() {
		return false, err
	}
	if index < 0 {
		return false, value.Undefined
	}
	existingVal := hashSet.Table[index]
	if existingVal == DeletedHashSetValue || existingVal.IsUndefined() {
		return false, value.Undefined
	}

	// `DeletedHashSetValue` means that the entry has been deleted
	hashSet.Table[index] = DeletedHashSetValue
	hashSet.Elements--

	return true, value.Undefined
}

// Check whether the given value is contained within the set.
func HashSetContains(vm *VM, set *value.HashSet, val value.Value) (bool, value.Value) {
	if set.Length() == 0 {
		return false, value.Undefined
	}

	index, err := HashSetIndex(vm, set, val)
	if !err.IsUndefined() {
		return false, err
	}
	if index == -1 {
		return false, value.Undefined
	}

	valInSlot := set.Table[index]
	if valInSlot == DeletedHashSetValue || valInSlot.IsUndefined() {
		return false, value.Undefined
	}
	return true, value.Undefined
}

func HashSetCopyTable(vm *VM, target *value.HashSet, source []value.Value) value.Value {
	for _, entry := range source {
		if entry == DeletedHashSetValue || entry.IsUndefined() {
			continue
		}

		err := HashSetAppendWithMaxLoad(vm, target, entry, 1)
		if !err.IsUndefined() {
			return err
		}
	}

	return value.Undefined
}

// Copy the pairs of one hashmap to the other.
func HashSetCopy(vm *VM, target *value.HashSet, source *value.HashSet) value.Value {
	requiredCapacity := target.Length() + source.Length()
	if target.Capacity() < requiredCapacity {
		HashSetSetCapacity(vm, target, requiredCapacity)
	}

	for _, entry := range source.Table {
		if entry == DeletedHashSetValue || entry.IsUndefined() {
			continue
		}

		i, err := HashSetIndex(vm, target, entry)
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
func HashSetGrow(vm *VM, set *value.HashSet, newSlots int) value.Value {
	return HashSetSetCapacity(vm, set, set.Capacity()+newSlots)
}

// Resize the given set to the desired capacity.
func HashSetSetCapacity(vm *VM, set *value.HashSet, capacity int) value.Value {
	if set.Capacity() == capacity {
		return value.Undefined
	}

	oldTable := set.Table
	newTable := make([]value.Value, capacity)
	tmpHashSet := &value.HashSet{
		Table: newTable,
	}

	for _, entry := range oldTable {
		if entry == DeletedHashSetValue || entry.IsUndefined() {
			continue
		}

		i, err := HashSetIndex(vm, tmpHashSet, entry)
		if !err.IsUndefined() {
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
	return value.Undefined
}

func HashSetAppendWithMaxLoad(vm *VM, set *value.HashSet, val value.Value, maxLoad float64) value.Value {
	if set.Capacity() == 0 {
		HashSetSetCapacity(vm, set, 5)
	} else if float64(set.OccupiedSlots) >= float64(set.Capacity())*maxLoad {
		HashSetSetCapacity(vm, set, set.OccupiedSlots*2)
	}

	index, err := HashSetIndex(vm, set, val)
	if !err.IsUndefined() {
		return err
	}
	if index == -1 {
		panic(fmt.Sprintf("no room in target hashset when trying to add a new value: %s", set.Inspect()))
	}
	entry := set.Table[index]
	if entry.IsUndefined() {
		set.OccupiedSlots++
		set.Elements++
	}

	set.Table[index] = val

	return value.Undefined
}

// Set a value under the given key.
func HashSetAppend(vm *VM, set *value.HashSet, val value.Value) value.Value {
	return HashSetAppendWithMaxLoad(vm, set, val, value.HashSetMaxLoad)
}

var DeletedHashSetValue value.Value = value.Ref(DeletedHashSetValueType{})

type DeletedHashSetValueType struct{}

func (DeletedHashSetValueType) Class() *value.Class {
	return nil
}
func (DeletedHashSetValueType) DirectClass() *value.Class {
	return nil
}
func (DeletedHashSetValueType) SingletonClass() *value.Class {
	return nil
}
func (DeletedHashSetValueType) Inspect() string {
	return "<empty_hash_set_slot>"
}
func (DeletedHashSetValueType) InstanceVariables() value.SymbolMap {
	return nil
}
func (e DeletedHashSetValueType) Copy() value.Reference {
	return e
}
func (e DeletedHashSetValueType) Error() string {
	return e.Inspect()
}

// Get the index that the value should be inserted into.
// Returns (0, err) when an error has been encountered.
// Returns (-1, undefined) when there's no room for new values.
func HashSetIndex(vm *VM, set *value.HashSet, val value.Value) (int, value.Value) {
	hash, err := Hash(vm, val)
	if !err.IsUndefined() {
		return 0, err
	}
	deletedIndex := -1

	capacity := set.Capacity()
	index := int(hash % value.UInt64(capacity))
	startIndex := index

	for {
		entry := set.Table[index]
		if entry.IsUndefined() {
			// empty bucket
			if deletedIndex != -1 {
				return deletedIndex, value.Undefined
			}
			return index, value.Undefined
		} else if entry == DeletedHashSetValue {
			if deletedIndex == -1 {
				// deleted entry
				deletedIndex = index
			}
		} else {
			// present entry
			equal, err := Equal(vm, entry, val)
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
