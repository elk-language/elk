package vm

import (
	"fmt"
	"iter"
	"strings"
	"sync"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type SyncOrderedHashSetOfValue struct {
	Set *OrderedHashSetOfValue
	mu  sync.Mutex
}

var _ SyncOrderedHashSet = &SyncOrderedHashSetOfValue{}

func NewSyncOrderedHashSetOfValue(capacity int) *SyncOrderedHashSetOfValue {
	return &SyncOrderedHashSetOfValue{
		Set: NewOrderedHashSetOfValue(capacity),
	}
}

func (h *SyncOrderedHashSetOfValue) Lock() {
	h.mu.Lock()
}

func (h *SyncOrderedHashSetOfValue) Unlock() {
	h.mu.Unlock()
}

func (h *SyncOrderedHashSetOfValue) OrderedHashSet() OrderedHashSet {
	return h.Set
}

func (h *SyncOrderedHashSetOfValue) ToHashSet(thread *Thread) (HashSet, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.ToHashSet(thread)
}

func (h *SyncOrderedHashSetOfValue) ToSyncHashSet(thread *Thread) (SyncHashSet, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.ToSyncHashSet(thread)
}

func (h *SyncOrderedHashSetOfValue) ToOrderedHashSet() OrderedHashSet {
	h.Lock()
	defer h.Unlock()

	return h.Set.Clone()
}

func (h *SyncOrderedHashSetOfValue) CloneSyncOrderedHashSet(thread *Thread, capacity int) (SyncOrderedHashSet, value.Value) {
	h.Lock()
	defer h.Unlock()

	cloned, err := h.Set.CloneOrderedHashSet(thread, capacity)
	if err.IsNotUndefined() {
		return nil, err
	}
	return &SyncOrderedHashSetOfValue{Set: cloned.(*OrderedHashSetOfValue)}, value.Undefined
}

func (h *SyncOrderedHashSetOfValue) NewSyncOrderedHashSet(capacity int) SyncOrderedHashSet {
	return NewSyncOrderedHashSetOfValue(capacity)
}

func (h *SyncOrderedHashSetOfValue) All() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		h.Lock()
		defer h.Unlock()

		for element := range h.Set.All() {
			if !yield(element) {
				return
			}
		}
	}
}

func (h *SyncOrderedHashSetOfValue) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		h.Lock()
		defer h.Unlock()

		for element, err := range h.Set.Iterate() {
			if !yield(element, err) {
				return
			}
		}
	}
}

func (*SyncOrderedHashSetOfValue) Class() *value.Class {
	return value.SyncOrderedHashSetClass
}

func (*SyncOrderedHashSetOfValue) DirectClass() *value.Class {
	return value.SyncOrderedHashSetClass
}

func (*SyncOrderedHashSetOfValue) SingletonClass() *value.Class {
	return nil
}

func (h *SyncOrderedHashSetOfValue) Iter() value.NativeIterator {
	return NewSyncOrderedHashSetOfValueIterator(h)
}

func (h *SyncOrderedHashSetOfValue) IterSyncOrderedSet() value.NativeResettableIterator {
	return NewSyncOrderedHashSetOfValueIterator(h)
}

func (h *SyncOrderedHashSetOfValue) Clone() *SyncOrderedHashSetOfValue {
	h.Lock()
	defer h.Unlock()

	return &SyncOrderedHashSetOfValue{
		Set: h.Set.Clone(),
	}
}

func (h *SyncOrderedHashSetOfValue) Copy() value.Reference {
	return h.Clone()
}

func (h *SyncOrderedHashSetOfValue) ToValue() value.Value {
	return value.Ref(h)
}

func (h *SyncOrderedHashSetOfValue) Error() string {
	return h.Inspect()
}

func (h *SyncOrderedHashSetOfValue) Inspect() string {
	h.Lock()
	defer h.Unlock()

	var hasMultilineElements bool
	elementStrings := make(
		[]string,
		0,
		min(MAX_HASH_SET_ELEMENTS_IN_INSPECT, h.Set.Length()),
	)

	var i int
	for element := range h.Set.All() {
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
	buff.WriteString("Std::Sync::OrderedHashSet[")
	if hasMultilineElements || h.Set.Length() > 15 {
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

func (*SyncOrderedHashSetOfValue) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncOrderedHashSetOfValue) Capacity() int {
	return h.Set.Capacity()
}

func (h *SyncOrderedHashSetOfValue) LeftCapacity() int {
	h.Lock()
	defer h.Unlock()

	return h.Set.LeftCapacity()
}

func (h *SyncOrderedHashSetOfValue) Length() int {
	return h.Set.Length()
}

func (h *SyncOrderedHashSetOfValue) GetIndex(thread *Thread, key value.Value) (value.SmallInt, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.GetIndex(thread, key)
}

// Get the element with the given index
// Does not do bounds checking or index normalisation.
// Panics if the index is invalid.
func (h *SyncOrderedHashSetOfValue) AtVal(index int) value.Value {
	h.Lock()
	defer h.Unlock()

	return h.Set.AtVal(index)
}

// Get the element with the given index
// Does not do bounds checking or index normalisation.
// Returns an error (undefined, err) if the index is out of range.
func (h *SyncOrderedHashSetOfValue) GetAtVal(index int) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.GetAtVal(index)
}

// Normalizes the index so that `-1` becomes `length-1`.
// Does bounds checking and returns (nil, undefined) if the index is out of range.
func (h *SyncOrderedHashSetOfValue) TryAtVal(index value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.TryAtVal(index)
}

func (h *SyncOrderedHashSetOfValue) RemoveVal(thread *Thread, other value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.RemoveVal(thread, other)
}

func (h *SyncOrderedHashSetOfValue) UnionVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	result, err := h.Set.Union(thread, other)
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return (&SyncOrderedHashSetOfValue{Set: result}).ToValue(), value.Undefined
}

func (h *SyncOrderedHashSetOfValue) IntersectionVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	result, err := h.Set.Intersection(thread, other)
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return (&SyncOrderedHashSetOfValue{Set: result}).ToValue(), value.Undefined
}

func (h *SyncOrderedHashSetOfValue) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncOrderedHashSetOfValue:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return h.Set.Equal(thread, o.Set.ToValue())
	case SyncOrderedHashSet:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return h.Set.Equal(thread, o.OrderedHashSet().ToValue())
	default:
		return false, value.Undefined
	}
}

func (h *SyncOrderedHashSetOfValue) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.LaxEqual(thread, other)
}

func (h *SyncOrderedHashSetOfValue) AppendVal(thread *Thread, val value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.AppendVal(thread, val)
}

func (h *SyncOrderedHashSetOfValue) Contains(thread *Thread, val value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.Contains(thread, val)
}

type SyncOrderedHashSetOfValueIterator struct {
	HashSet  *SyncOrderedHashSetOfValue
	index    int
	snapshot []value.Value
	version  int
}

var _ value.NativeResettableIterator = &SyncOrderedHashSetOfValueIterator{}

func NewSyncOrderedHashSetOfValueIterator(set *SyncOrderedHashSetOfValue) *SyncOrderedHashSetOfValueIterator {
	iterator := &SyncOrderedHashSetOfValueIterator{
		HashSet: set,
		version: set.Set.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *SyncOrderedHashSetOfValueIterator) captureSnapshot() {
	h.HashSet.Lock()
	defer h.HashSet.Unlock()

	snapshot := make([]value.Value, 0, h.HashSet.Set.Length())
	for element := range h.HashSet.Set.All() {
		snapshot = append(snapshot, element)
	}
	h.snapshot = snapshot
}

func (*SyncOrderedHashSetOfValueIterator) Class() *value.Class {
	return value.SyncOrderedHashSetIteratorClass
}

func (*SyncOrderedHashSetOfValueIterator) DirectClass() *value.Class {
	return value.SyncOrderedHashSetIteratorClass
}

func (*SyncOrderedHashSetOfValueIterator) SingletonClass() *value.Class {
	return nil
}

func (h *SyncOrderedHashSetOfValueIterator) Copy() value.Reference {
	return &SyncOrderedHashSetOfValueIterator{
		HashSet:  h.HashSet,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *SyncOrderedHashSetOfValueIterator) ToValue() value.Value {
	return value.Ref(i)
}

func (h *SyncOrderedHashSetOfValueIterator) Error() string {
	return h.Inspect()
}

func (h *SyncOrderedHashSetOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::Sync::OrderedHashSet::Iterator{hash_set: %s}", h.HashSet.Inspect())
}

func (*SyncOrderedHashSetOfValueIterator) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncOrderedHashSetOfValueIterator) NextValue() (value.Value, value.Value) {
	if h.version != h.HashSet.Set.version {
		return value.Undefined, value.NewMutationDuringIterationError(h.HashSet.Class().Name()).ToValue()
	}
	if h.index >= len(h.snapshot) {
		return value.Undefined, value.S(symbol.L_stop_iteration).ToValue()
	}

	element := h.snapshot[h.index]
	h.index++
	return element, value.Undefined
}

func (h *SyncOrderedHashSetOfValueIterator) Reset() {
	h.index = 0
	h.version = h.HashSet.Set.version
}
