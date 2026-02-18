package vm

import (
	"fmt"
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

type HashSet interface {
	value.ValueInterface
	value.NativeIterable
	All() iter.Seq[value.Value]
	Length() int
	UnionVal(*Thread, value.Value) (result value.Value, err value.Value)
	IntersectionVal(*Thread, value.Value) (result value.Value, err value.Value)
	Equal(*Thread, value.Value) (result bool, err value.Value)
	Contains(*Thread, value.Value) (result bool, err value.Value)
	AppendVal(*Thread, value.Value) (added bool, err value.Value)
	RemoveVal(thread *Thread, other value.Value) (removed bool, err value.Value)
	IterSet() value.NativeResettableIterator
}

// ::Std::HashSet
func initHashSet() {
	// Instance methods
	c := &value.HashSetClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			iterator := self.IterSet()
			return iterator.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"remove",
		func(thread *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			other := args[1]
			removed, err := self.RemoveVal(thread, other)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			return value.ToElkBool(removed), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			otherVal := args[1]
			return self.UnionVal(vm, otherVal)
		},
		DefWithParameters(1),
	)
	Alias(c, "union", "+")
	Alias(c, "|", "+")

	Def(
		c,
		"&",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			otherVal := args[1]
			return self.IntersectionVal(vm, otherVal)
		},
		DefWithParameters(1),
	)
	Alias(c, "intersection", "&")

	Def(
		c,
		"==",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			otherVal := args[1]
			result, err := self.Equal(vm, otherVal)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			return value.ToElkBool(result), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "=~", "==")

	Def(
		c,
		"contains",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			val := args[1]
			contains, err := self.Contains(vm, val)
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
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			for element, err := range Iterate(vm, args[1]) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				_, err := self.AppendVal(vm, element)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
			}
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<<",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			val := args[1]
			_, err := self.AppendVal(vm, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"push",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			val := args[1]
			added, err := self.AppendVal(vm, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.ToElkBool(added), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			callable := args[1]
			newSet := NewHashSetOfValue(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(*Closure); ok {
				for val := range self.All() {
					result, err := vm.CallClosure(function, val)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					_, err = newSet.AppendVal(vm, result)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newSet), value.Undefined
			}

			// callable is another value
			for val := range self.All() {
				result, err := vm.CallMethodByName(callSymbol, callable, val)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				_, err = newSet.AppendVal(vm, result)
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

func NewHashSetOfValueComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *HashSetOfValue) bool {
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
		for _, xVal := range x.table {
			if xVal == DeletedHashSetValue || xVal.IsUndefined() {
				continue
			}

			contains, err := HashSetOfValueContains(v, y, xVal)
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
func NewHashSetOfValueWithElements(vm *Thread, elements ...value.Value) (*HashSetOfValue, value.Value) {
	return NewHashSetOfValueWithCapacityAndElements(vm, len(elements), elements...)
}

// Create a new hash set with the given entries.
func MustNewHashSetOfValueWithElements(vm *Thread, elements ...value.Value) *HashSetOfValue {
	set, err := NewHashSetOfValueWithElements(vm, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return set
}

func NewHashSetOfValueWithCapacityAndElements(vm *Thread, capacity int, elements ...value.Value) (*HashSetOfValue, value.Value) {
	s := NewHashSetOfValue(capacity)
	for _, element := range elements {
		_, err := HashSetOfValueAppend(vm, s, element)
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return s, value.Undefined
}

func NewHashSetOfValueWithCapacityAndElementsMaxLoad(vm *Thread, capacity int, maxLoad float64, elements ...value.Value) (*HashSetOfValue, value.Value) {
	s := NewHashSetOfValue(capacity)
	for _, element := range elements {
		_, err := HashSetOfValueAppendWithMaxLoad(vm, s, element, maxLoad)
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return s, value.Undefined
}

func MustNewHashSetOfValueWithCapacityAndElements(vm *Thread, capacity int, elements ...value.Value) *HashSetOfValue {
	set, err := NewHashSetOfValueWithCapacityAndElements(vm, capacity, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return set
}

func MustNewHashSetOfValueWithCapacityAndElementsMaxLoad(vm *Thread, capacity int, maxLoad float64, elements ...value.Value) *HashSetOfValue {
	set, err := NewHashSetOfValueWithCapacityAndElementsMaxLoad(vm, capacity, maxLoad, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return set
}

// Checks whether two hash sets are equal
func HashSetOfValueEqual(vm *Thread, x *HashSetOfValue, y *HashSetOfValue) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for _, xVal := range x.table {
		if xVal == DeletedHashSetValue || xVal.IsUndefined() {
			continue
		}

		contains, err := HashSetOfValueContains(vm, y, xVal)
		if !err.IsUndefined() {
			return false, err
		}
		if !contains {
			return false, value.Undefined
		}
	}

	return true, value.Undefined
}

func HashSetOfValueEqualInterface(vm *Thread, x *HashSetOfValue, y HashSet) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for _, xVal := range x.table {
		if xVal == DeletedHashSetValue || xVal.IsUndefined() {
			continue
		}

		contains, err := y.Contains(vm, xVal)
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
func HashSetOfValueUnion(vm *Thread, x *HashSetOfValue, y *HashSetOfValue) (*HashSetOfValue, value.Value) {
	var longer *HashSetOfValue
	var shorter *HashSetOfValue
	if x.Length() > y.Length() {
		longer = x
		shorter = y
	} else {
		longer = y
		shorter = x
	}

	newSet := NewHashSetOfValue(shorter.Length() + longer.Length())
	HashSetOfValueCopy(vm, newSet, longer)
	for _, shorterVal := range shorter.table {
		if shorterVal == DeletedHashSetValue || shorterVal.IsUndefined() {
			continue
		}

		_, err := HashSetOfValueAppend(vm, newSet, shorterVal)
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return newSet, value.Undefined
}

// Create a new set that is the union of the given two sets
func HashSetOfValueUnionInterface(vm *Thread, x *HashSetOfValue, y HashSet) (*HashSetOfValue, value.Value) {
	newSet := NewHashSetOfValue(x.Length() + y.Length())
	HashSetOfValueCopy(vm, newSet, x)
	for v := range y.All() {
		_, err := HashSetOfValueAppend(vm, newSet, v)
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return newSet, value.Undefined
}

// Create a new set that is the intersection of the given two sets
func HashSetOfValueIntersection(vm *Thread, x *HashSetOfValue, y *HashSetOfValue) (*HashSetOfValue, value.Value) {
	var longer *HashSetOfValue
	var shorter *HashSetOfValue
	if x.Length() > y.Length() {
		longer = x
		shorter = y
	} else {
		longer = y
		shorter = x
	}

	newSet := NewHashSetOfValue(5)
	for _, shorterVal := range shorter.table {
		if shorterVal == DeletedHashSetValue || shorterVal.IsUndefined() {
			continue
		}

		contains, err := HashSetOfValueContains(vm, longer, shorterVal)
		if !err.IsUndefined() {
			return nil, err
		}
		if contains {
			HashSetOfValueAppend(vm, newSet, shorterVal)
		}
	}

	return newSet, value.Undefined
}

// Create a new set that is the intersection of the given two sets
func HashSetOfValueIntersectionInterface(vm *Thread, x *HashSetOfValue, y HashSet) (*HashSetOfValue, value.Value) {
	newSet := NewHashSetOfValue(5)
	for yVal := range y.All() {
		if yVal == DeletedHashSetValue || yVal.IsUndefined() {
			continue
		}

		contains, err := HashSetOfValueContains(vm, x, yVal)
		if !err.IsUndefined() {
			return nil, err
		}
		if contains {
			HashSetOfValueAppend(vm, newSet, yVal)
		}
	}

	return newSet, value.Undefined
}

// Delete the given value from the hash set
func HashSetOfValueDelete(vm *Thread, hashSet *HashSetOfValue, val value.Value) (bool, value.Value) {
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
	existingVal := hashSet.table[index]
	if existingVal == DeletedHashSetValue || existingVal.IsUndefined() {
		return false, value.Undefined
	}

	// `DeletedHashSetValue` means that the entry has been deleted
	hashSet.table[index] = DeletedHashSetValue
	hashSet.elements--

	return true, value.Undefined
}

// Check whether the given value is contained within the set.
func HashSetOfValueContains(vm *Thread, set *HashSetOfValue, val value.Value) (bool, value.Value) {
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

	valInSlot := set.table[index]
	if valInSlot == DeletedHashSetValue || valInSlot.IsUndefined() {
		return false, value.Undefined
	}
	return true, value.Undefined
}

func HashSetOfValueCopyTable(vm *Thread, target *HashSetOfValue, source []value.Value) value.Value {
	for _, entry := range source {
		if entry == DeletedHashSetValue || entry.IsUndefined() {
			continue
		}

		_, err := HashSetOfValueAppendWithMaxLoad(vm, target, entry, 1)
		if !err.IsUndefined() {
			return err
		}
	}

	return value.Undefined
}

// Copy the pairs of one hashmap to the other.
func HashSetOfValueCopy(vm *Thread, target *HashSetOfValue, source *HashSetOfValue) value.Value {
	requiredCapacity := target.Length() + source.Length()
	if target.Capacity() < requiredCapacity {
		HashSetOfValueSetCapacity(vm, target, requiredCapacity)
	}

	for _, entry := range source.table {
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
		target.table[i] = entry
		target.occupiedSlots++
		target.elements++
	}

	return value.Undefined
}

// Add additional n empty slots for new elements.
func HashSetOfValueGrow(vm *Thread, set *HashSetOfValue, newSlots int) value.Value {
	return HashSetOfValueSetCapacity(vm, set, set.Capacity()+newSlots)
}

// Resize the given set to the desired capacity.
func HashSetOfValueSetCapacity(vm *Thread, set *HashSetOfValue, capacity int) value.Value {
	if set.Capacity() == capacity {
		return value.Undefined
	}

	oldTable := set.table
	newTable := make([]value.Value, capacity)
	tmpHashSet := &HashSetOfValue{
		table: newTable,
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
		tmpHashSet.occupiedSlots++
		tmpHashSet.elements++
	}

	set.occupiedSlots = tmpHashSet.occupiedSlots
	set.elements = tmpHashSet.elements
	set.table = newTable
	return value.Undefined
}

func HashSetOfValueAppendWithMaxLoad(vm *Thread, set *HashSetOfValue, val value.Value, maxLoad float64) (bool, value.Value) {
	if set.Capacity() == 0 {
		HashSetOfValueSetCapacity(vm, set, 5)
	} else if float64(set.occupiedSlots) >= float64(set.Capacity())*maxLoad {
		HashSetOfValueSetCapacity(vm, set, set.occupiedSlots*2)
	}

	index, err := HashSetIndex(vm, set, val)
	if !err.IsUndefined() {
		return false, err
	}
	if index == -1 {
		panic(fmt.Sprintf("no room in target hashset when trying to add a new value: %s", set.Inspect()))
	}
	entry := set.table[index]

	var newValue bool

	if entry.IsUndefined() {
		// the slot is empty
		set.occupiedSlots++
		set.elements++
		newValue = true
	} else if entry == DeletedHashSetValue {
		// this is a zombie slot, just overwrite it's content
		set.elements++
		newValue = true
	}

	set.table[index] = val

	return newValue, value.Undefined
}

// Set a value under the given key.
func HashSetOfValueAppend(vm *Thread, set *HashSetOfValue, val value.Value) (bool, value.Value) {
	return HashSetOfValueAppendWithMaxLoad(vm, set, val, HashSetMaxLoad)
}

var DeletedHashSetValue value.Value = value.Ref(deletedHashSetValueType{})

type deletedHashSetValueType struct{}

var _ value.Reference = deletedHashSetValueType{}

func (deletedHashSetValueType) Class() *value.Class {
	return nil
}
func (deletedHashSetValueType) DirectClass() *value.Class {
	return nil
}
func (deletedHashSetValueType) SingletonClass() *value.Class {
	return nil
}
func (deletedHashSetValueType) Inspect() string {
	return "<empty_hash_set_slot>"
}
func (deletedHashSetValueType) InstanceVariables() *value.InstanceVariables {
	return nil
}
func (e deletedHashSetValueType) Copy() value.Reference {
	return e
}
func (e deletedHashSetValueType) ToValue() value.Value {
	return value.Ref(e)
}
func (e deletedHashSetValueType) Error() string {
	return e.Inspect()
}

// Get the index that the value should be inserted into.
// Returns (0, err) when an error has been encountered.
// Returns (-1, undefined) when there's no room for new values.
func HashSetIndex(vm *Thread, set *HashSetOfValue, val value.Value) (int, value.Value) {
	hash, err := Hash(vm, val)
	if !err.IsUndefined() {
		return 0, err
	}
	deletedIndex := -1

	capacity := set.Capacity()
	index := int(hash % value.UInt64(capacity))
	startIndex := index

	for {
		entry := set.table[index]
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
