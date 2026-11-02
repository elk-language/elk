package vm

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/google/go-cmp/cmp"
)

const HashSetMaxLoad = 0.75

type HashSetOfValue struct {
	table         []value.Value // underlying data container, it's `len` is always equal to it's `cap` and it also serves as the capacity of the HashSet
	occupiedSlots int           // number of slots taken by active values and those left by deleted values
	elements      int           // number of slots occupied by active values
	version       int           // version of the set, each mutation increments this counter, used for guarding against concurrent mutation during iteration
}

var _ HashSet = &HashSetOfValue{}

func NewHashSetOfValue(capacity int) *HashSetOfValue {
	return &HashSetOfValue{
		table: make([]value.Value, capacity),
	}
}

func (h *HashSetOfValue) CloneHashSet(thread *Thread, capacity int) (HashSet, value.Value) {
	newSet := NewHashSetOfValue(capacity)
	err := HashSetOfValueCopy(thread, newSet, h)
	if err.IsNotUndefined() {
		return nil, err
	}

	return newSet, value.Undefined
}

func (h *HashSetOfValue) NewHashSet(capacity int) HashSet {
	return NewHashSetOfValue(capacity)
}

func (h *HashSetOfValue) All() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		for _, element := range h.table {
			if element.IsUndefined() || element == DeletedHashSetValue {
				continue
			}

			if !yield(element) {
				return
			}
		}
	}
}

func (h *HashSetOfValue) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		originalVersion := h.version
		for _, element := range h.table {
			if element.IsUndefined() || element == DeletedHashSetValue {
				continue
			}

			if originalVersion != h.version {
				yield(value.Undefined, value.NewMutationDuringIterationError(h.Class().Name()).ToValue())
				return
			}

			if !yield(element, value.Undefined) {
				return
			}
		}
	}
}

func (*HashSetOfValue) Class() *value.Class {
	return value.HashSetClass
}

func (*HashSetOfValue) DirectClass() *value.Class {
	return value.HashSetClass
}

func (*HashSetOfValue) SingletonClass() *value.Class {
	return nil
}

func (h *HashSetOfValue) Iter() value.NativeIterator {
	return NewHashSetOfValueIterator(h)
}

func (h *HashSetOfValue) IterSet() value.NativeResettableIterator {
	return NewHashSetOfValueIterator(h)
}

func (h *HashSetOfValue) Clone() *HashSetOfValue {
	newTable := slices.Clone(h.table)
	return &HashSetOfValue{
		table:         newTable,
		occupiedSlots: h.occupiedSlots,
		elements:      h.elements,
	}
}

func (h *HashSetOfValue) ToSyncHashSet() SyncHashSet {
	return &SyncHashSetOfValue{Set: h.Clone()}
}

func (h *HashSetOfValue) Copy() value.Reference {
	return h.Clone()
}

func (h *HashSetOfValue) ToValue() value.Value {
	return value.Ref(h)
}

func (h *HashSetOfValue) Error() string {
	return h.Inspect()
}

const MAX_HASH_SET_ELEMENTS_IN_INSPECT = 300

func (h *HashSetOfValue) Inspect() string {
	var hasMultilineElements bool
	elementStrings := make(
		[]string,
		0,
		min(MAX_HASH_SET_ELEMENTS_IN_INSPECT, h.Length()),
	)

	var i int
	for element := range h.All() {
		elementString := element.Inspect()
		elementStrings = append(elementStrings, elementString)
		if strings.ContainsRune(elementString, '\n') {
			hasMultilineElements = true
		}

		if i >= MAX_HASH_SET_ELEMENTS_IN_INSPECT-1 {
			break
		}
		i++
	}

	var buff strings.Builder
	buff.WriteString("^[")
	if hasMultilineElements || h.Length() > 15 {
		buff.WriteRune('\n')
		for i, elementString := range elementStrings {
			if i != 0 {
				buff.WriteString(",\n")
			}

			indent.IndentString(&buff, elementString, 1)

			if i >= MAX_HASH_SET_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(",\n  ...")
				break
			}
		}
		buff.WriteRune('\n')
	} else {
		for i, elementString := range elementStrings {
			if i != 0 {
				buff.WriteString(", ")
			}

			buff.WriteString(elementString)

			if i >= MAX_HASH_SET_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(", ...")
				break
			}
		}
	}
	buff.WriteRune(']')

	return buff.String()
}

func (*HashSetOfValue) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *HashSetOfValue) Capacity() int {
	return len(h.table)
}

func (h *HashSetOfValue) LeftCapacity() int {
	return h.Capacity() - h.Length()
}

func (h *HashSetOfValue) Length() int {
	return h.elements
}

func (h *HashSetOfValue) RemoveVal(thread *Thread, other value.Value) (removed bool, err value.Value) {
	return HashSetOfValueDelete(thread, h, other)
}

func (h *HashSetOfValue) Union(thread *Thread, other value.Value) (result *HashSetOfValue, err value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashSetOfValue:
		return HashSetOfValueUnion(thread, h, o)
	case HashSet:
		return HashSetOfValueUnionInterface(thread, h, o)
	case SyncHashSet:
		o.Lock()
		defer o.Unlock()

		return HashSetOfValueUnionInterface(thread, h, o.HashSet())
	default:
		return HashSetOfValueUnionValue(thread, h, other)
	}
}

func (h *HashSetOfValue) UnionVal(thread *Thread, other value.Value) (result value.Value, err value.Value) {
	return value.RefErr(h.Union(thread, other))
}

func (h *HashSetOfValue) Intersection(thread *Thread, other value.Value) (result *HashSetOfValue, err value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashSetOfValue:
		return HashSetOfValueIntersection(thread, h, o)
	case HashSet:
		return HashSetOfValueIntersectionInterface(thread, h, o)
	case SyncHashSet:
		o.Lock()
		defer o.Unlock()

		return HashSetOfValueIntersectionInterface(thread, h, o.HashSet())
	default:
		return HashSetOfValueIntersectionValue(thread, h, other)
	}
}

func (h *HashSetOfValue) IntersectionVal(thread *Thread, other value.Value) (result value.Value, err value.Value) {
	return value.RefErr(h.Intersection(thread, other))
}

func (h *HashSetOfValue) Equal(thread *Thread, other value.Value) (result bool, err value.Value) {
	switch other := other.SafeAsReference().(type) {
	case *HashSetOfValue:
		return HashSetOfValueEqual(thread, h, other)
	case HashSet:
		return HashSetOfValueEqualInterface(thread, h, other)
	default:
		return false, value.Undefined
	}
}

func (h *HashSetOfValue) LaxEqual(thread *Thread, other value.Value) (result bool, err value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashSetOfValue:
		return HashSetOfValueEqual(thread, h, o)
	case HashSet:
		return HashSetOfValueEqualInterface(thread, h, o)
	case SyncHashSet:
		o.Lock()
		defer o.Unlock()

		return HashSetOfValueEqualInterface(thread, h, o.HashSet())
	default:
		return HashSetOfValueLaxEqualValue(thread, h, other)
	}
}

func (h *HashSetOfValue) AppendVal(thread *Thread, other value.Value) (result bool, err value.Value) {
	return HashSetOfValueAppend(thread, h, other)
}

func (h *HashSetOfValue) Contains(thread *Thread, other value.Value) (result bool, err value.Value) {
	return HashSetOfValueContains(thread, h, other)
}

type HashSetOfValueIterator struct {
	HashSet *HashSetOfValue
	index   int
	version int
}

var _ value.NativeResettableIterator = &HashSetOfValueIterator{}

func NewHashSetOfValueIterator(set *HashSetOfValue) *HashSetOfValueIterator {
	return &HashSetOfValueIterator{
		HashSet: set,
		version: set.version,
	}
}

func (*HashSetOfValueIterator) Class() *value.Class {
	return value.HashSetIteratorClass
}

func (*HashSetOfValueIterator) DirectClass() *value.Class {
	return value.HashSetIteratorClass
}

func (*HashSetOfValueIterator) SingletonClass() *value.Class {
	return nil
}

func (h *HashSetOfValueIterator) Copy() value.Reference {
	return &HashSetOfValueIterator{
		HashSet: h.HashSet,
		index:   h.index,
	}
}

func (i *HashSetOfValueIterator) ToValue() value.Value {
	return value.Ref(i)
}

func (h *HashSetOfValueIterator) Error() string {
	return h.Inspect()
}

func (h *HashSetOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::HashSet::Iterator{hash_set: %s}", h.HashSet.Inspect())
}

func (*HashSetOfValueIterator) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *HashSetOfValueIterator) NextValue() (value.Value, value.Value) {
	if h.version != h.HashSet.version {
		return value.Undefined, value.NewMutationDuringIterationError(h.HashSet.Class().Name()).ToValue()
	}

	for {
		if h.index >= h.HashSet.Capacity() {
			return value.Undefined, value.S(symbol.L_stop_iteration).ToValue()
		}

		element := h.HashSet.table[h.index]
		h.index++
		if !element.IsUndefined() {
			return element, value.Undefined
		}
	}
}

func (h *HashSetOfValueIterator) Reset() {
	h.index = 0
}

func NewHashSetComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y HashSet) bool {
		if x == y {
			return true
		}
		if x.Length() != y.Length() {
			return false
		}

		result := DefaultThreadPool.Call(func(vm *Thread) (result value.Value, err value.Value) {
			for xVal := range x.All() {
				contains, err := y.Contains(vm, xVal)
				if !err.IsUndefined() {
					return value.False.ToValue(), value.Undefined
				}
				if !contains {
					return value.False.ToValue(), value.Undefined
				}
			}

			return value.True.ToValue(), value.Undefined
		}).MustAwaitSync()

		return value.Truthy(result)
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

func HashSetOfValueLaxEqualValue(thread *Thread, h *HashSetOfValue, other value.Value) (bool, value.Value) {
	if !value.IsA(other, value.SetMixin) {
		return false, value.Undefined
	}

	oLengthVal, err := thread.CallMethodByName(value.S(symbol.L_length), other)
	if err.IsNotUndefined() {
		return false, err
	}

	oLength := oLengthVal.AsInt()
	if h.Length() != oLength {
		return false, value.Undefined
	}

	for _, element := range Iterate(thread, other) {
		contains, err := h.Contains(thread, element)
		if err.IsNotUndefined() {
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
	err := HashSetOfValueCopy(vm, newSet, longer)
	if err.IsNotUndefined() {
		return nil, err
	}
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
	err := HashSetOfValueCopy(vm, newSet, x)
	if err.IsNotUndefined() {
		return nil, err
	}

	for v := range y.All() {
		_, err := HashSetOfValueAppend(vm, newSet, v)
		if err.IsNotUndefined() {
			return nil, err
		}
	}

	return newSet, value.Undefined
}

// Create a new set that is the union of the given two sets
func HashSetOfValueUnionValue(thread *Thread, h *HashSetOfValue, other value.Value) (*HashSetOfValue, value.Value) {
	if !value.IsA(other, value.SetMixin) {
		return nil, value.NewCoerceError(value.HashSetClass, other.Class()).ToValue()
	}

	oLengthVal, err := thread.CallMethodByName(value.S(symbol.L_length), other)
	if err.IsNotUndefined() {
		return nil, err
	}

	oLength := oLengthVal.AsInt()
	newSet := NewHashSetOfValue(h.Length() + oLength)

	err = HashSetOfValueCopy(thread, newSet, h)
	if err.IsNotUndefined() {
		return nil, err
	}

	for element, err := range Iterate(thread, other) {
		if err.IsNotUndefined() {
			return nil, err
		}

		_, err := newSet.AppendVal(thread, element)
		if err.IsNotUndefined() {
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
		if err.IsNotUndefined() {
			return nil, err
		}
		if contains {
			_, err := HashSetOfValueAppend(vm, newSet, shorterVal)
			if err.IsNotUndefined() {
				return nil, err
			}
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
			_, err := HashSetOfValueAppend(vm, newSet, yVal)
			if err.IsNotUndefined() {
				return nil, err
			}
		}
	}

	return newSet, value.Undefined
}

// Create a new set that is the intersection of the given two sets
func HashSetOfValueIntersectionValue(thread *Thread, h *HashSetOfValue, other value.Value) (*HashSetOfValue, value.Value) {
	if !value.IsA(other, value.SetMixin) {
		return nil, value.NewCoerceError(value.HashSetClass, other.Class()).ToValue()
	}

	newSet := NewHashSetOfValue(5)

	for element, err := range Iterate(thread, other) {
		if err.IsNotUndefined() {
			return nil, err
		}

		contains, err := h.Contains(thread, element)
		if err.IsNotUndefined() {
			return nil, err
		}

		if contains {
			_, err := newSet.AppendVal(thread, element)
			if err.IsNotUndefined() {
				return nil, err
			}
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
	hashSet.version++

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

	target.version++
	return value.Undefined
}

// Copy the pairs of one hashmap to the other.
func HashSetOfValueCopy(vm *Thread, target *HashSetOfValue, source *HashSetOfValue) value.Value {
	requiredCapacity := target.Length() + source.Length()
	if target.Capacity() < requiredCapacity {
		err := HashSetOfValueSetCapacity(vm, target, requiredCapacity)
		if err.IsNotUndefined() {
			return err
		}
	}

	target.version++
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
	set.version++

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
		err := HashSetOfValueSetCapacity(vm, set, 5)
		if err.IsNotUndefined() {
			return false, err
		}
	} else if float64(set.occupiedSlots) >= float64(set.Capacity())*maxLoad {
		err := HashSetOfValueSetCapacity(vm, set, set.occupiedSlots*2)
		if err.IsNotUndefined() {
			return false, err
		}
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
	set.version++

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
