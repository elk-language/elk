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

type SyncOrderedNativeHashSet[V value.ComparableValueInterface] struct {
	Set *OrderedNativeHashSet[V]
	mu  sync.Mutex
}

var _ SyncOrderedHashSet = &SyncOrderedNativeHashSet[value.String]{}

func NewSyncOrderedNativeHashSet[V value.ComparableValueInterface](capacity int) *SyncOrderedNativeHashSet[V] {
	return &SyncOrderedNativeHashSet[V]{
		Set: NewOrderedNativeHashSet[V](capacity),
	}
}

func (h *SyncOrderedNativeHashSet[V]) Lock() {
	h.mu.Lock()
}

func (h *SyncOrderedNativeHashSet[V]) Unlock() {
	h.mu.Unlock()
}

func (h *SyncOrderedNativeHashSet[V]) OrderedHashSet() OrderedHashSet {
	return h.Set
}

func (h *SyncOrderedNativeHashSet[V]) ToHashSet() (HashSet, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.ToHashSet()
}

func (h *SyncOrderedNativeHashSet[V]) ToSyncHashSet() (SyncHashSet, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.ToSyncHashSet()
}

func (h *SyncOrderedNativeHashSet[V]) ToOrderedHashSet() OrderedHashSet {
	h.Lock()
	defer h.Unlock()

	return h.Set.Clone()
}

func (h *SyncOrderedNativeHashSet[V]) CloneSyncOrderedHashSet(thread *Thread, capacity int) (SyncOrderedHashSet, value.Value) {
	h.Lock()
	defer h.Unlock()

	cloned, err := h.Set.CloneOrderedHashSet(thread, capacity)
	if err.IsNotUndefined() {
		return nil, err
	}
	return &SyncOrderedNativeHashSet[V]{Set: cloned.(*OrderedNativeHashSet[V])}, value.Undefined
}

func (h *SyncOrderedNativeHashSet[V]) NewSyncOrderedHashSet(capacity int) SyncOrderedHashSet {
	return NewSyncOrderedNativeHashSet[V](capacity)
}

func (h *SyncOrderedNativeHashSet[V]) All() iter.Seq[value.Value] {
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

func (h *SyncOrderedNativeHashSet[V]) AllNative() iter.Seq[V] {
	return func(yield func(V) bool) {
		h.Lock()
		defer h.Unlock()

		for element := range h.Set.AllNative() {
			if !yield(element) {
				return
			}
		}
	}
}

func (h *SyncOrderedNativeHashSet[V]) Iterate() iter.Seq2[value.Value, value.Value] {
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

func (*SyncOrderedNativeHashSet[V]) Class() *value.Class {
	return value.SyncOrderedHashSetClass
}

func (*SyncOrderedNativeHashSet[V]) DirectClass() *value.Class {
	return value.SyncOrderedHashSetClass
}

func (*SyncOrderedNativeHashSet[V]) SingletonClass() *value.Class {
	return nil
}

func (h *SyncOrderedNativeHashSet[V]) Iter() value.NativeIterator {
	return NewSyncOrderedNativeHashSetIterator(h)
}

func (h *SyncOrderedNativeHashSet[V]) IterSyncOrderedSet() value.NativeResettableIterator {
	return NewSyncOrderedNativeHashSetIterator(h)
}

func (h *SyncOrderedNativeHashSet[V]) Clone() *SyncOrderedNativeHashSet[V] {
	h.Lock()
	defer h.Unlock()

	return &SyncOrderedNativeHashSet[V]{
		Set: h.Set.Clone(),
	}
}

func (h *SyncOrderedNativeHashSet[V]) Copy() value.Reference {
	return h.Clone()
}

func (h *SyncOrderedNativeHashSet[V]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *SyncOrderedNativeHashSet[V]) Error() string {
	return h.Inspect()
}

func (h *SyncOrderedNativeHashSet[V]) Inspect() string {
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

func (*SyncOrderedNativeHashSet[V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncOrderedNativeHashSet[V]) Capacity() int {
	return h.Set.Capacity()
}

func (h *SyncOrderedNativeHashSet[V]) LeftCapacity() int {
	h.Lock()
	defer h.Unlock()

	return h.Set.LeftCapacity()
}

func (h *SyncOrderedNativeHashSet[V]) Length() int {
	return h.Set.Length()
}

func (h *SyncOrderedNativeHashSet[V]) GetIndex(val V) (int, bool) {
	h.Lock()
	defer h.Unlock()

	return h.Set.GetIndex(val)
}

func (h *SyncOrderedNativeHashSet[V]) At(index int) V {
	h.Lock()
	defer h.Unlock()

	return h.Set.At(index)
}

// Get the element with the given index
// Does not do bounds checking or index normalisation.
// Panics if the index is invalid.
func (h *SyncOrderedNativeHashSet[V]) AtVal(index int) value.Value {
	h.Lock()
	defer h.Unlock()

	return h.Set.AtVal(index)
}

func (h *SyncOrderedNativeHashSet[V]) GetAt(index int) (V, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.GetAt(index)
}

// Get the element with the given index
// Does not do bounds checking or index normalisation.
// Returns an error (undefined, err) if the index is out of range.
func (h *SyncOrderedNativeHashSet[V]) GetAtVal(index int) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.GetAtVal(index)
}

// Normalizes the index so that `-1` becomes `length-1`.
// Does bounds checking and returns (undefined, undefined) if the index is out of range.
func (h *SyncOrderedNativeHashSet[V]) TryAtVal(index value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.TryAtVal(index)
}

func (h *SyncOrderedNativeHashSet[V]) Remove(val V) bool {
	h.Lock()
	defer h.Unlock()

	return h.Set.Remove(val)
}

func (h *SyncOrderedNativeHashSet[V]) RemoveVal(thread *Thread, other value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.RemoveVal(thread, other)
}

func wrapSyncOrderedNativeHashSetResult[V value.ComparableValueInterface](result OrderedHashSet, err value.Value) (value.Value, value.Value) {
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	switch set := result.(type) {
	case *OrderedNativeHashSet[V]:
		return (&SyncOrderedNativeHashSet[V]{Set: set}).ToValue(), value.Undefined
	case *OrderedHashSetOfValue:
		return (&SyncOrderedHashSetOfValue{Set: set}).ToValue(), value.Undefined
	default:
		return result.ToValue(), value.Undefined
	}
}

func (h *SyncOrderedNativeHashSet[V]) UnionVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncOrderedNativeHashSet[V]:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return wrapSyncOrderedNativeHashSetResult[V](h.Set.Union(thread, o.Set.ToValue()))
	case *OrderedNativeHashSet[V]:
		h.Lock()
		defer h.Unlock()

		return wrapSyncOrderedNativeHashSetResult[V](h.Set.Union(thread, o.ToValue()))
	case SyncOrderedHashSet:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return wrapSyncOrderedNativeHashSetResult[V](h.Set.Union(thread, o.OrderedHashSet().ToValue()))
	default:
		h.Lock()
		defer h.Unlock()

		return wrapSyncOrderedNativeHashSetResult[V](h.Set.Union(thread, other))
	}
}

func (h *SyncOrderedNativeHashSet[V]) IntersectionVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncOrderedNativeHashSet[V]:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return wrapSyncOrderedNativeHashSetResult[V](h.Set.Intersection(thread, o.Set.ToValue()))
	case *OrderedNativeHashSet[V]:
		h.Lock()
		defer h.Unlock()

		return wrapSyncOrderedNativeHashSetResult[V](h.Set.Intersection(thread, o.ToValue()))
	case SyncOrderedHashSet:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return wrapSyncOrderedNativeHashSetResult[V](h.Set.Intersection(thread, o.OrderedHashSet().ToValue()))
	default:
		h.Lock()
		defer h.Unlock()

		return wrapSyncOrderedNativeHashSetResult[V](h.Set.Intersection(thread, other))
	}
}

func (h *SyncOrderedNativeHashSet[V]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncOrderedNativeHashSet[V]:
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

func (h *SyncOrderedNativeHashSet[V]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.LaxEqual(thread, other)
}

func (h *SyncOrderedNativeHashSet[V]) Append(val V) bool {
	h.Lock()
	defer h.Unlock()

	return h.Set.Append(val)
}

func (h *SyncOrderedNativeHashSet[V]) AppendVal(thread *Thread, val value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.AppendVal(thread, val)
}

func (h *SyncOrderedNativeHashSet[V]) ContainsNative(val V) bool {
	h.Lock()
	defer h.Unlock()

	return h.Set.ContainsNative(val)
}

func (h *SyncOrderedNativeHashSet[V]) Contains(thread *Thread, val value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.Contains(thread, val)
}

type SyncOrderedNativeHashSetIterator[V value.ComparableValueInterface] struct {
	HashSet  *SyncOrderedNativeHashSet[V]
	index    int
	snapshot []V
	version  int
}

var _ value.NativeResettableIterator = &SyncOrderedNativeHashSetIterator[value.String]{}

func NewSyncOrderedNativeHashSetIterator[V value.ComparableValueInterface](set *SyncOrderedNativeHashSet[V]) *SyncOrderedNativeHashSetIterator[V] {
	iterator := &SyncOrderedNativeHashSetIterator[V]{
		HashSet: set,
		version: set.Set.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *SyncOrderedNativeHashSetIterator[V]) captureSnapshot() {
	h.HashSet.Lock()
	defer h.HashSet.Unlock()

	snapshot := make([]V, 0, h.HashSet.Set.Length())
	for element := range h.HashSet.Set.AllNative() {
		snapshot = append(snapshot, element)
	}
	h.snapshot = snapshot
}

func (*SyncOrderedNativeHashSetIterator[V]) Class() *value.Class {
	return value.SyncOrderedHashSetIteratorClass
}

func (*SyncOrderedNativeHashSetIterator[V]) DirectClass() *value.Class {
	return value.SyncOrderedHashSetIteratorClass
}

func (*SyncOrderedNativeHashSetIterator[V]) SingletonClass() *value.Class {
	return nil
}

func (h *SyncOrderedNativeHashSetIterator[V]) Copy() value.Reference {
	return &SyncOrderedNativeHashSetIterator[V]{
		HashSet:  h.HashSet,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *SyncOrderedNativeHashSetIterator[V]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *SyncOrderedNativeHashSetIterator[V]) Error() string {
	return h.Inspect()
}

func (h *SyncOrderedNativeHashSetIterator[V]) Inspect() string {
	return fmt.Sprintf("Std::Sync::OrderedHashSet::Iterator{hash_set: %s}", h.HashSet.Inspect())
}

func (*SyncOrderedNativeHashSetIterator[V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncOrderedNativeHashSetIterator[V]) Next() (v V, err value.Value) {
	if h.version != h.HashSet.Set.version {
		return v, value.NewMutationDuringIterationError(h.HashSet.Class().Name()).ToValue()
	}
	if h.index >= len(h.snapshot) {
		return v, value.S(symbol.L_stop_iteration).ToValue()
	}

	v = h.snapshot[h.index]
	h.index++
	return v, value.Undefined
}

func (h *SyncOrderedNativeHashSetIterator[V]) NextValue() (value.Value, value.Value) {
	element, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return element.ToValue(), value.Undefined
}

func (h *SyncOrderedNativeHashSetIterator[V]) Reset() {
	h.index = 0
	h.version = h.HashSet.Set.version
}
