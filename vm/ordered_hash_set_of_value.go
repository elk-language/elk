package vm

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type OrderedHashSetOfValue struct {
	index   *HashMapOfValue
	entries []value.Value
	version int // version of the set, each mutation increments this counter, used for guarding against concurrent mutation during iteration
}

var _ OrderedHashSet = &OrderedHashSetOfValue{}

func NewOrderedHashSetOfValue(capacity int) *OrderedHashSetOfValue {
	return &OrderedHashSetOfValue{
		index:   NewHashMapOfValue(capacity),
		entries: make([]value.Value, 0, capacity),
	}
}

func (h *OrderedHashSetOfValue) CloneOrderedHashSet(thread *Thread, capacity int) (OrderedHashSet, value.Value) {
	newMap := NewHashMapOfValue(capacity)
	err := HashMapOfValueCopy(thread, newMap, h.index)
	if !err.IsUndefined() {
		return nil, err
	}
	entries := make([]value.Value, len(h.entries), capacity)
	copy(entries, h.entries)

	result := &OrderedHashSetOfValue{
		index:   newMap,
		entries: entries,
	}
	return result, value.Undefined
}

func (h *OrderedHashSetOfValue) NewOrderedHashSet(capacity int) OrderedHashSet {
	return NewOrderedHashSetOfValue(capacity)
}

func (h *OrderedHashSetOfValue) All() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		for _, val := range h.entries {
			if !yield(val) {
				return
			}
		}
	}
}

func (h *OrderedHashSetOfValue) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		originalVersion := h.version
		for _, val := range h.entries {
			if originalVersion != h.version {
				yield(value.Undefined, value.NewMutationDuringIterationError(h.Class().Name()).ToValue())
				return
			}

			if !yield(val, value.Undefined) {
				return
			}
		}
	}
}

func (*OrderedHashSetOfValue) Class() *value.Class {
	return value.OrderedHashSetClass
}

func (*OrderedHashSetOfValue) DirectClass() *value.Class {
	return value.OrderedHashSetClass
}

func (*OrderedHashSetOfValue) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedHashSetOfValue) Iter() value.NativeIterator {
	return NewOrderedHashSetOfValueIterator(h)
}

func (h *OrderedHashSetOfValue) IterOrderedSet() value.NativeResettableIterator {
	return NewOrderedHashSetOfValueIterator(h)
}

func (h *OrderedHashSetOfValue) Clone() *OrderedHashSetOfValue {
	return &OrderedHashSetOfValue{
		index:   h.index.Clone(),
		entries: slices.Clone(h.entries),
	}
}

func (h *OrderedHashSetOfValue) ToHashSetOfValue(thread *Thread) (*HashSetOfValue, value.Value) {
	set := NewHashSetOfValue(len(h.entries))
	for _, element := range h.entries {
		_, err := HashSetOfValueAppend(thread, set, element)
		if err.IsNotUndefined() {
			return nil, err
		}
	}
	return set, value.Undefined
}
func (h *OrderedHashSetOfValue) ToHashSet(thread *Thread) (HashSet, value.Value) {
	return h.ToHashSetOfValue(thread)
}

func (h *OrderedHashSetOfValue) ToSyncHashSet(thread *Thread) (SyncHashSet, value.Value) {
	set, err := h.ToHashSetOfValue(thread)
	if err.IsNotUndefined() {
		return nil, err
	}
	return &SyncHashSetOfValue{Set: set}, value.Undefined
}

func (h *OrderedHashSetOfValue) ToSyncOrderedHashSet() SyncOrderedHashSet {
	return &SyncOrderedHashSetOfValue{Set: h.Clone()}
}

func (h *OrderedHashSetOfValue) Copy() value.Reference {
	return h.Clone()
}

func (h *OrderedHashSetOfValue) ToValue() value.Value {
	return value.Ref(h)
}

func (h *OrderedHashSetOfValue) Error() string {
	return h.Inspect()
}

func (h *OrderedHashSetOfValue) Inspect() string {
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
	buff.WriteString("Std::OrderedHashSet[")
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

func (*OrderedHashSetOfValue) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedHashSetOfValue) Capacity() int {
	return cap(h.entries)
}

func (h *OrderedHashSetOfValue) LeftCapacity() int {
	return h.Capacity() - h.Length()
}

func (h *OrderedHashSetOfValue) Length() int {
	return len(h.entries)
}

func (h *OrderedHashSetOfValue) GetIndex(thread *Thread, key value.Value) (value.SmallInt, value.Value) {
	indexPair, err := h.index.GetPair(thread, key)
	if err.IsNotUndefined() {
		return 0, err
	}
	if indexPair.Key().IsUndefined() {
		return -1, value.Undefined
	}

	index := indexPair.Value().AsSmallInt()
	return index, value.Undefined
}

// Get the element with the given index
// Does not do bounds checking or index normalisation.
// Panics if the index is invalid.
func (h *OrderedHashSetOfValue) AtVal(index int) value.Value {
	return h.entries[index]
}

// Get the element with the given index
// Does not do bounds checking or index normalisation.
// Returns an error (undefined, err) if the index is out of range.
func (h *OrderedHashSetOfValue) GetAtVal(index int) (result value.Value, err value.Value) {
	index, err = value.NormalizeArrayIndex(index, h.Length())
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return h.entries[index], value.Undefined
}

// Normalizes the index so that `-1` becomes `length-1`.
// Does bounds checking and returns (undefined, undefined) if the index is out of range.
func (h *OrderedHashSetOfValue) TryAtVal(index value.Value) (value.Value, value.Value) {
	var i int

	i, ok := value.ToGoInt(index)
	if !ok {
		if i == -1 {
			return value.Undefined, value.Undefined
		}
		return value.Undefined, value.NewCoerceError(value.IntClass, index.Class()).ToValue()
	}

	i, ok = value.NormalizeArrayIndexOk(i, h.Length())
	if !ok {
		return value.Undefined, value.Undefined
	}

	return h.entries[i], value.Undefined
}

func (h *OrderedHashSetOfValue) RemoveVal(thread *Thread, other value.Value) (removed bool, err value.Value) {
	index, err := h.GetIndex(thread, other)
	if err.IsNotUndefined() {
		return false, err
	}
	if index == -1 {
		return false, value.Undefined
	}

	h.version++
	h.entries = slices.Delete(h.entries, int(index), int(index)+1)
	_, err = h.index.RemoveKey(thread, other)
	if err.IsNotUndefined() {
		return false, err
	}

	for i := int(index); i < len(h.entries); i++ {
		err = h.index.SetVal(thread, h.entries[i], value.SmallInt(i).ToValue())
		if err.IsNotUndefined() {
			return false, err
		}
	}

	return true, value.Undefined
}

func (h *OrderedHashSetOfValue) Union(thread *Thread, other value.Value) (result *OrderedHashSetOfValue, err value.Value) {
	switch o := other.SafeAsReference().(type) {
	case Set:
		return OrderedHashSetUnionSet(thread, h, o)
	default:
		return OrderedHashSetUnionValue(thread, h, other)
	}
}

func (h *OrderedHashSetOfValue) UnionVal(thread *Thread, other value.Value) (result value.Value, err value.Value) {
	return value.RefErr(h.Union(thread, other))
}

func (h *OrderedHashSetOfValue) Intersection(thread *Thread, other value.Value) (result *OrderedHashSetOfValue, err value.Value) {
	switch o := other.SafeAsReference().(type) {
	case Set:
		return OrderedHashSetIntersectionSet(thread, h, o)
	default:
		return OrderedHashSetIntersectionValue(thread, h, other)
	}
}

func (h *OrderedHashSetOfValue) IntersectionVal(thread *Thread, other value.Value) (result value.Value, err value.Value) {
	return value.RefErr(h.Intersection(thread, other))
}

func (h *OrderedHashSetOfValue) Equal(thread *Thread, other value.Value) (result bool, err value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedHashSetOfValue:
		return ValueSliceEqual(thread, h.entries, o.entries)
	case OrderedHashSet:
		otherEntries := slices.Collect(o.All())
		return ValueSliceEqual(thread, h.entries, otherEntries)
	default:
		return false, value.Undefined
	}
}

func (h *OrderedHashSetOfValue) LaxEqual(thread *Thread, other value.Value) (result bool, err value.Value) {
	switch o := other.SafeAsReference().(type) {
	case Set:
		return OrderedHashSetLaxEqualSet(thread, h, o)
	default:
		return OrderedHashSetLaxEqualValue(thread, h, other)
	}
}

func (h *OrderedHashSetOfValue) AppendVal(thread *Thread, val value.Value) (result bool, err value.Value) {
	existingIndex, err := h.index.GetValUndefined(thread, val)
	if err.IsNotUndefined() {
		return false, err
	}

	if existingIndex.IsNotUndefined() {
		return false, value.Undefined
	}

	// new key
	entryIndex := len(h.entries)
	h.entries = append(h.entries, val)
	h.index.SetVal(thread, val, value.SmallInt(entryIndex).ToValue())
	return true, value.Undefined
}

func (h *OrderedHashSetOfValue) Contains(thread *Thread, val value.Value) (result bool, err value.Value) {
	ok, err := h.index.ContainsKey(thread, val)
	if err.IsNotUndefined() {
		return false, err
	}
	return ok, value.Undefined
}

type OrderedHashSetOfValueIterator struct {
	HashSet  *OrderedHashSetOfValue
	index    int
	snapshot []value.Value
	version  int
}

var _ value.NativeResettableIterator = &OrderedHashSetOfValueIterator{}

func NewOrderedHashSetOfValueIterator(set *OrderedHashSetOfValue) *OrderedHashSetOfValueIterator {
	iterator := &OrderedHashSetOfValueIterator{
		HashSet: set,
		version: set.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *OrderedHashSetOfValueIterator) captureSnapshot() {
	snapshot := make([]value.Value, 0, h.HashSet.Length())
	for element := range h.HashSet.All() {
		snapshot = append(snapshot, element)
	}
	h.snapshot = snapshot
}

func (*OrderedHashSetOfValueIterator) Class() *value.Class {
	return value.OrderedHashSetIteratorClass
}

func (*OrderedHashSetOfValueIterator) DirectClass() *value.Class {
	return value.OrderedHashSetIteratorClass
}

func (*OrderedHashSetOfValueIterator) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedHashSetOfValueIterator) Copy() value.Reference {
	return &OrderedHashSetOfValueIterator{
		HashSet:  h.HashSet,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *OrderedHashSetOfValueIterator) ToValue() value.Value {
	return value.Ref(i)
}

func (h *OrderedHashSetOfValueIterator) Error() string {
	return h.Inspect()
}

func (h *OrderedHashSetOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::OrderedHashSet::Iterator{hash_set: %s}", h.HashSet.Inspect())
}

func (*OrderedHashSetOfValueIterator) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedHashSetOfValueIterator) NextValue() (value.Value, value.Value) {
	if h.version != h.HashSet.version {
		return value.Undefined, value.NewMutationDuringIterationError(h.HashSet.Class().Name()).ToValue()
	}
	if h.index >= len(h.snapshot) {
		return value.Undefined, value.S(symbol.L_stop_iteration).ToValue()
	}

	element := h.snapshot[h.index]
	h.index++
	return element, value.Undefined
}

func (h *OrderedHashSetOfValueIterator) Reset() {
	h.index = 0
	h.version = h.HashSet.version
}
