package vm

import (
	"fmt"
	"iter"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type SyncHashSetOfValue struct {
	Set *HashSetOfValue
	mu  value.Mutex
}

var _ SyncHashSet = &SyncHashSetOfValue{}

func NewSyncHashSetOfValue(capacity int) *SyncHashSetOfValue {
	return &SyncHashSetOfValue{
		Set: NewHashSetOfValue(capacity),
	}
}

func (h *SyncHashSetOfValue) Lock() {
	h.mu.Lock()
}

func (h *SyncHashSetOfValue) Unlock() {
	h.mu.Unlock()
}

func (h *SyncHashSetOfValue) HashSet() HashSet {
	return h.Set
}

func (h *SyncHashSetOfValue) CloneSyncHashSet(thread *Thread, capacity int) (SyncHashSet, value.Value) {
	h.Lock()
	defer h.Unlock()

	newSet := NewSyncHashSetOfValue(capacity)
	err := HashSetOfValueCopy(thread, newSet.Set, h.Set)
	if err.IsNotUndefined() {
		return nil, err
	}

	return newSet, value.Undefined
}

func (h *SyncHashSetOfValue) NewSyncHashSet(capacity int) SyncHashSet {
	return NewSyncHashSetOfValue(capacity)
}

func (h *SyncHashSetOfValue) All() iter.Seq[value.Value] {
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

func (h *SyncHashSetOfValue) Iterate() iter.Seq2[value.Value, value.Value] {
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

func (*SyncHashSetOfValue) Class() *value.Class {
	return value.SyncHashSetClass
}

func (*SyncHashSetOfValue) DirectClass() *value.Class {
	return value.SyncHashSetClass
}

func (*SyncHashSetOfValue) SingletonClass() *value.Class {
	return nil
}

func (h *SyncHashSetOfValue) Iter() value.NativeIterator {
	return NewSyncHashSetOfValueIterator(h)
}

func (h *SyncHashSetOfValue) IterSet() value.NativeResettableIterator {
	return NewSyncHashSetOfValueIterator(h)
}

func (h *SyncHashSetOfValue) IterSyncSet() value.NativeResettableIterator {
	return NewSyncHashSetOfValueIterator(h)
}

func (h *SyncHashSetOfValue) Copy() value.Reference {
	h.Lock()
	defer h.Unlock()

	return &SyncHashSetOfValue{
		Set: h.Set.Copy().(*HashSetOfValue),
	}
}

func (h *SyncHashSetOfValue) ToValue() value.Value {
	return value.Ref(h)
}

func (h *SyncHashSetOfValue) Error() string {
	return h.Inspect()
}

func (h *SyncHashSetOfValue) Inspect() string {
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
	buff.WriteString("Std::Sync::HashSet[")
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

func (*SyncHashSetOfValue) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncHashSetOfValue) Capacity() int {
	return h.Set.Capacity()
}

func (h *SyncHashSetOfValue) LeftCapacity() int {
	h.Lock()
	defer h.Unlock()

	return h.Set.LeftCapacity()
}

func (h *SyncHashSetOfValue) Length() int {
	return h.Set.Length()
}

func (h *SyncHashSetOfValue) RemoveVal(thread *Thread, other value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.RemoveVal(thread, other)
}

func (h *SyncHashSetOfValue) UnionVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	result, err := h.Set.Union(thread, other)
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return (&SyncHashSetOfValue{Set: result}).ToValue(), value.Undefined
}

func (h *SyncHashSetOfValue) IntersectionVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	result, err := h.Set.Intersection(thread, other)
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return (&SyncHashSetOfValue{Set: result}).ToValue(), value.Undefined
}

func (h *SyncHashSetOfValue) Equal(thread *Thread, other value.Value) (result bool, err value.Value) {
	switch other := other.SafeAsReference().(type) {
	case *SyncHashSetOfValue:
		h.Lock()
		defer h.Unlock()

		other.Lock()
		defer other.Unlock()

		return HashSetOfValueEqual(thread, h.Set, other.Set)
	case SyncHashSet:
		h.Lock()
		defer h.Unlock()

		other.Lock()
		defer other.Unlock()

		return HashSetOfValueEqualInterface(thread, h.Set, other.HashSet())
	default:
		return false, value.Undefined
	}
}

func (h *SyncHashSetOfValue) LaxEqual(thread *Thread, other value.Value) (result bool, err value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.LaxEqual(thread, other)
}

func (h *SyncHashSetOfValue) AppendVal(thread *Thread, other value.Value) (result bool, err value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.AppendVal(thread, other)
}

func (h *SyncHashSetOfValue) Contains(thread *Thread, other value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.Contains(thread, other)
}

type SyncHashSetOfValueIterator struct {
	HashSet  *SyncHashSetOfValue
	index    int
	snapshot []value.Value
	version  int
}

var _ value.NativeResettableIterator = &SyncHashSetOfValueIterator{}

func NewSyncHashSetOfValueIterator(set *SyncHashSetOfValue) *SyncHashSetOfValueIterator {
	iterator := &SyncHashSetOfValueIterator{
		HashSet: set,
		version: set.Set.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *SyncHashSetOfValueIterator) captureSnapshot() {
	h.HashSet.Lock()
	defer h.HashSet.Unlock()

	snapshot := make([]value.Value, 0, h.HashSet.Set.Length())
	for element := range h.HashSet.Set.All() {
		snapshot = append(snapshot, element)
	}
	h.snapshot = snapshot
}

func (*SyncHashSetOfValueIterator) Class() *value.Class {
	return value.SyncHashSetIteratorClass
}

func (*SyncHashSetOfValueIterator) DirectClass() *value.Class {
	return value.SyncHashSetIteratorClass
}

func (*SyncHashSetOfValueIterator) SingletonClass() *value.Class {
	return nil
}

func (h *SyncHashSetOfValueIterator) Copy() value.Reference {
	return &SyncHashSetOfValueIterator{
		HashSet:  h.HashSet,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *SyncHashSetOfValueIterator) ToValue() value.Value {
	return value.Ref(i)
}

func (h *SyncHashSetOfValueIterator) Error() string {
	return h.Inspect()
}

func (h *SyncHashSetOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::Sync::HashSet::Iterator{hash_set: %s}", h.HashSet.Inspect())
}

func (*SyncHashSetOfValueIterator) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncHashSetOfValueIterator) NextValue() (value.Value, value.Value) {
	if h.version != h.HashSet.Set.version {
		return value.Undefined, value.NewMutationDuringIterationError(h.Class().Name()).ToValue()
	}
	if h.index >= len(h.snapshot) {
		return value.Undefined, value.S(symbol.L_stop_iteration).ToValue()
	}

	element := h.snapshot[h.index]
	h.index++
	return element, value.Undefined
}

func (h *SyncHashSetOfValueIterator) Reset() {
	h.index = 0
}
