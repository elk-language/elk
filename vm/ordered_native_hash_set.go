package vm

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type OrderedNativeHashSet[V value.ComparableValueInterface] struct {
	index   map[V]int
	entries []V
	version int // version of the set, each mutation increments this counter, used for guarding against concurrent mutation during iteration
}

var _ OrderedHashSet = &OrderedNativeHashSet[value.String]{}

func NewOrderedNativeHashSet[V value.ComparableValueInterface](capacity int) *OrderedNativeHashSet[V] {
	return &OrderedNativeHashSet[V]{
		index:   make(map[V]int, capacity),
		entries: make([]V, 0, capacity),
	}
}

func (h *OrderedNativeHashSet[V]) CloneOrderedHashSet(thread *Thread, capacity int) (OrderedHashSet, value.Value) {
	newIndex := make(map[V]int, capacity)
	maps.Copy(newIndex, h.index)
	entries := make([]V, len(h.entries), capacity)
	copy(entries, h.entries)

	result := &OrderedNativeHashSet[V]{
		index:   newIndex,
		entries: entries,
	}
	return result, value.Undefined
}

func (h *OrderedNativeHashSet[V]) NewOrderedHashSet(capacity int) OrderedHashSet {
	return NewOrderedNativeHashSet[V](capacity)
}

func (h *OrderedNativeHashSet[V]) All() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		for _, element := range h.entries {
			if !yield(element.ToValue()) {
				return
			}
		}
	}
}

func (h *OrderedNativeHashSet[V]) AllNative() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, element := range h.entries {
			if !yield(element) {
				return
			}
		}
	}
}

func (h *OrderedNativeHashSet[V]) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		originalVersion := h.version
		for _, element := range h.entries {
			if originalVersion != h.version {
				yield(value.Undefined, value.NewMutationDuringIterationError(h.Class().Name()).ToValue())
				return
			}

			if !yield(element.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (*OrderedNativeHashSet[V]) Class() *value.Class {
	return value.OrderedHashSetClass
}

func (*OrderedNativeHashSet[V]) DirectClass() *value.Class {
	return value.OrderedHashSetClass
}

func (*OrderedNativeHashSet[V]) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedNativeHashSet[V]) Iter() value.NativeIterator {
	return NewOrderedNativeHashSetIterator(h)
}

func (h *OrderedNativeHashSet[V]) IterOrderedSet() value.NativeResettableIterator {
	return NewOrderedNativeHashSetIterator(h)
}

func (h *OrderedNativeHashSet[V]) Clone() *OrderedNativeHashSet[V] {
	return &OrderedNativeHashSet[V]{
		index:   maps.Clone(h.index),
		entries: slices.Clone(h.entries),
	}
}

func (h *OrderedNativeHashSet[V]) ToHashSet(thread *Thread) (HashSet, value.Value) {
	return h.ToNativeHashSet()
}

func (h *OrderedNativeHashSet[V]) ToNativeHashSet() (*NativeHashSet[V], value.Value) {
	set := NewNativeHashSet[V](len(h.entries))
	for _, element := range h.entries {
		set.Append(element)
	}
	return set, value.Undefined
}

func (h *OrderedNativeHashSet[V]) ToSyncHashSet(thread *Thread) (SyncHashSet, value.Value) {
	set, err := h.ToNativeHashSet()
	if err.IsNotUndefined() {
		return nil, err
	}
	return &SyncNativeHashSet[V]{Set: set}, value.Undefined
}

func (h *OrderedNativeHashSet[V]) ToSyncOrderedHashSet() SyncOrderedHashSet {
	return &SyncOrderedNativeHashSet[V]{Set: h.Clone()}
}

func (h *OrderedNativeHashSet[V]) Copy() value.Reference {
	return h.Clone()
}

func (h *OrderedNativeHashSet[V]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *OrderedNativeHashSet[V]) Error() string {
	return h.Inspect()
}

func (h *OrderedNativeHashSet[V]) Inspect() string {
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

func (*OrderedNativeHashSet[V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedNativeHashSet[V]) Capacity() int {
	return cap(h.entries)
}

func (h *OrderedNativeHashSet[V]) LeftCapacity() int {
	return h.Capacity() - h.Length()
}

func (h *OrderedNativeHashSet[V]) Length() int {
	return len(h.entries)
}

func (h *OrderedNativeHashSet[V]) GetIndex(val V) (int, bool) {
	index, ok := h.index[val]
	return index, ok
}

func (h *OrderedNativeHashSet[V]) At(index int) V {
	return h.entries[index]
}

// Get the element with the given index
// Does not do bounds checking or index normalisation.
// Panics if the index is invalid.
func (h *OrderedNativeHashSet[V]) AtVal(index int) value.Value {
	return h.At(index).ToValue()
}

func (h *OrderedNativeHashSet[V]) GetAt(index int) (v V, err value.Value) {
	index, err = value.NormalizeArrayIndex(index, h.Length())
	if err.IsNotUndefined() {
		return v, err
	}

	return h.entries[index], value.Undefined
}

// Get the element with the given index
// Does not do bounds checking or index normalisation.
// Returns an error (undefined, err) if the index is out of range.
func (h *OrderedNativeHashSet[V]) GetAtVal(index int) (value.Value, value.Value) {
	v, err := h.GetAt(index)
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return v.ToValue(), value.Undefined
}

// Normalizes the index so that `-1` becomes `length-1`.
// Does bounds checking and returns (undefined, undefined) if the index is out of range.
func (h *OrderedNativeHashSet[V]) TryAtVal(index value.Value) (value.Value, value.Value) {
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

	return h.entries[i].ToValue(), value.Undefined
}

func (h *OrderedNativeHashSet[V]) Remove(val V) bool {
	index, ok := h.GetIndex(val)
	if !ok {
		return false
	}

	h.version++
	h.entries = slices.Delete(h.entries, index, index+1)
	delete(h.index, val)
	for i := index; i < len(h.entries); i++ {
		h.index[h.entries[i]] = i
	}
	return true
}

func (h *OrderedNativeHashSet[V]) RemoveVal(thread *Thread, other value.Value) (bool, value.Value) {
	val, ok := value.Downcast[V](other)
	if !ok {
		return false, value.NewInvalidElementInTypedSet(h, other.Class()).ToValue()
	}
	return h.Remove(val), value.Undefined
}

func (h *OrderedNativeHashSet[V]) Union(thread *Thread, other value.Value) (OrderedHashSet, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedNativeHashSet[V]:
		result := NewOrderedNativeHashSet[V](h.Length() + o.Length())
		for _, element := range h.entries {
			result.Append(element)
		}
		for _, element := range o.entries {
			result.Append(element)
		}
		return result, value.Undefined
	case Set:
		return OrderedHashSetUnionSet(thread, h, o)
	default:
		return OrderedHashSetUnionValue(thread, h, other)
	}
}

func (h *OrderedNativeHashSet[V]) UnionVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	result, err := h.Union(thread, other)
	if err.IsNotUndefined() {
		return value.Undefined, err
	}
	return result.ToValue(), value.Undefined
}

func (h *OrderedNativeHashSet[V]) Intersection(thread *Thread, other value.Value) (OrderedHashSet, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedNativeHashSet[V]:
		result := NewOrderedNativeHashSet[V](h.Length())
		for _, element := range h.entries {
			if _, present := o.index[element]; present {
				result.Append(element)
			}
		}
		return result, value.Undefined
	case Set:
		return OrderedHashSetIntersectionSet(thread, h, o)
	default:
		return OrderedHashSetIntersectionValue(thread, h, other)
	}
}

func (h *OrderedNativeHashSet[V]) IntersectionVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	result, err := h.Intersection(thread, other)
	if err.IsNotUndefined() {
		return value.Undefined, err
	}
	return result.ToValue(), value.Undefined
}

func (h *OrderedNativeHashSet[V]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedNativeHashSet[V]:
		return NativeValueSliceEqual(thread, h.entries, o.entries)
	case OrderedHashSet:
		if len(h.entries) != o.Length() {
			return false, value.Undefined
		}

		i := 0
		for element := range o.All() {
			eq, err := Equal(thread, element, h.entries[i].ToValue())
			if err.IsNotUndefined() {
				return false, err
			}
			if !value.Truthy(eq) {
				return false, value.Undefined
			}
			i++
		}
		return true, value.Undefined
	default:
		return false, value.Undefined
	}
}

func (h *OrderedNativeHashSet[V]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case Set:
		return OrderedHashSetLaxEqualSet(thread, h, o)
	default:
		return OrderedHashSetLaxEqualValue(thread, h, other)
	}
}

func (h *OrderedNativeHashSet[V]) Append(val V) bool {
	if _, present := h.index[val]; present {
		return false
	}

	h.version++
	h.index[val] = len(h.entries)
	h.entries = append(h.entries, val)
	return true
}

func (h *OrderedNativeHashSet[V]) AppendVal(thread *Thread, val value.Value) (bool, value.Value) {
	element, ok := value.Downcast[V](val)
	if !ok {
		return false, value.NewInvalidElementInTypedSet(h, val.Class()).ToValue()
	}
	return h.Append(element), value.Undefined
}

func (h *OrderedNativeHashSet[V]) ContainsNative(val V) bool {
	_, present := h.index[val]
	return present
}

func (h *OrderedNativeHashSet[V]) Contains(thread *Thread, val value.Value) (bool, value.Value) {
	element, ok := value.Downcast[V](val)
	if !ok {
		return false, value.Undefined
	}
	return h.ContainsNative(element), value.Undefined
}

type OrderedNativeHashSetIterator[V value.ComparableValueInterface] struct {
	HashSet  *OrderedNativeHashSet[V]
	index    int
	snapshot []V
	version  int
}

var _ value.NativeResettableIterator = &OrderedNativeHashSetIterator[value.String]{}

func NewOrderedNativeHashSetIterator[V value.ComparableValueInterface](set *OrderedNativeHashSet[V]) *OrderedNativeHashSetIterator[V] {
	iterator := &OrderedNativeHashSetIterator[V]{
		HashSet: set,
		version: set.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *OrderedNativeHashSetIterator[V]) captureSnapshot() {
	snapshot := make([]V, 0, h.HashSet.Length())
	for element := range h.HashSet.AllNative() {
		snapshot = append(snapshot, element)
	}
	h.snapshot = snapshot
}

func (*OrderedNativeHashSetIterator[V]) Class() *value.Class {
	return value.OrderedHashSetIteratorClass
}

func (*OrderedNativeHashSetIterator[V]) DirectClass() *value.Class {
	return value.OrderedHashSetIteratorClass
}

func (*OrderedNativeHashSetIterator[V]) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedNativeHashSetIterator[V]) Copy() value.Reference {
	return &OrderedNativeHashSetIterator[V]{
		HashSet:  h.HashSet,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *OrderedNativeHashSetIterator[V]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *OrderedNativeHashSetIterator[V]) Error() string {
	return h.Inspect()
}

func (h *OrderedNativeHashSetIterator[V]) Inspect() string {
	return fmt.Sprintf("Std::OrderedHashSet::Iterator{hash_set: %s}", h.HashSet.Inspect())
}

func (*OrderedNativeHashSetIterator[V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedNativeHashSetIterator[V]) Next() (v V, err value.Value) {
	if h.version != h.HashSet.version {
		return v, value.NewMutationDuringIterationError(h.HashSet.Class().Name()).ToValue()
	}
	if h.index >= len(h.snapshot) {
		return v, value.S(symbol.L_stop_iteration).ToValue()
	}

	v = h.snapshot[h.index]
	h.index++
	return v, value.Undefined
}

func (h *OrderedNativeHashSetIterator[V]) NextValue() (value.Value, value.Value) {
	element, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return element.ToValue(), value.Undefined
}

func (h *OrderedNativeHashSetIterator[V]) Reset() {
	h.index = 0
	h.version = h.HashSet.version
}
