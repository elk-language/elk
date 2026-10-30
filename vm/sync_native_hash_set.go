package vm

import (
	"fmt"
	"iter"
	"maps"
	"strings"
	"sync"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type SyncNativeHashSet[V value.ComparableValueInterface] struct {
	Set *NativeHashSet[V]
	mu  sync.Mutex
}

var _ SyncHashSet = &SyncNativeHashSet[value.String]{}

func NewSyncNativeHashSet[V value.ComparableValueInterface](capacity int) *SyncNativeHashSet[V] {
	return &SyncNativeHashSet[V]{
		Set: NewNativeHashSet[V](capacity),
	}
}

func (h *SyncNativeHashSet[V]) Lock() {
	h.mu.Lock()
}

func (h *SyncNativeHashSet[V]) Unlock() {
	h.mu.Unlock()
}

func (h *SyncNativeHashSet[V]) HashSet() HashSet {
	return h.Set
}

func (h *SyncNativeHashSet[V]) CloneSyncHashSet(thread *Thread, capacity int) (SyncHashSet, value.Value) {
	h.Lock()
	defer h.Unlock()

	newSet := NewSyncNativeHashSet[V](capacity)
	maps.Copy(newSet.Set.m, h.Set.m)
	return newSet, value.Undefined
}

func (h *SyncNativeHashSet[V]) NewSyncHashSet(capacity int) SyncHashSet {
	return NewSyncNativeHashSet[V](capacity)
}

func (h *SyncNativeHashSet[V]) All() iter.Seq[value.Value] {
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

func (h *SyncNativeHashSet[V]) AllNative() iter.Seq[V] {
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

func (h *SyncNativeHashSet[V]) Iterate() iter.Seq2[value.Value, value.Value] {
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

func (*SyncNativeHashSet[V]) Class() *value.Class {
	return value.SyncHashSetClass
}

func (*SyncNativeHashSet[V]) DirectClass() *value.Class {
	return value.SyncHashSetClass
}

func (*SyncNativeHashSet[V]) SingletonClass() *value.Class {
	return nil
}

func (h *SyncNativeHashSet[V]) Iter() value.NativeIterator {
	return NewSyncNativeHashSetIterator(h)
}

func (h *SyncNativeHashSet[V]) IterSyncSet() value.NativeResettableIterator {
	return NewSyncNativeHashSetIterator(h)
}

func (h *SyncNativeHashSet[V]) Copy() value.Reference {
	h.Lock()
	defer h.Unlock()

	return &SyncNativeHashSet[V]{
		Set: h.Set.Copy().(*NativeHashSet[V]),
	}
}

func (h *SyncNativeHashSet[V]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *SyncNativeHashSet[V]) Error() string {
	return h.Inspect()
}

func (h *SyncNativeHashSet[V]) Inspect() string {
	h.Lock()
	defer h.Unlock()

	var hasMultilineElements bool
	elementStrings := make(
		[]string,
		0,
		min(MAX_HASH_SET_ELEMENTS_IN_INSPECT, h.Set.Length()),
	)

	var i int
	for element := range h.Set.m {
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

func (*SyncNativeHashSet[V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncNativeHashSet[V]) Length() int {
	return h.Set.Length()
}

func (h *SyncNativeHashSet[V]) Remove(other V) bool {
	h.Lock()
	defer h.Unlock()

	return h.Set.Remove(other)
}

func (h *SyncNativeHashSet[V]) RemoveVal(thread *Thread, other value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.RemoveVal(thread, other)
}

func (h *SyncNativeHashSet[V]) Union(other *SyncNativeHashSet[V]) *SyncNativeHashSet[V] {
	h.Lock()
	defer h.Unlock()

	other.Lock()
	defer other.Unlock()

	return &SyncNativeHashSet[V]{Set: h.Set.Union(other.Set)}
}

func (h *SyncNativeHashSet[V]) UnionInterface(thread *Thread, other SyncHashSet) (*SyncHashSetOfValue, value.Value) {
	h.Lock()
	defer h.Unlock()

	other.Lock()
	defer other.Unlock()

	result, err := h.Set.UnionInterface(thread, other.HashSet())
	if err.IsNotUndefined() {
		return nil, err
	}

	return &SyncHashSetOfValue{Set: result}, value.Undefined
}

func (h *SyncNativeHashSet[V]) unionValue(thread *Thread, other value.Value) (*SyncHashSetOfValue, value.Value) {
	h.Lock()
	defer h.Unlock()

	result, err := h.Set.unionValue(thread, other)
	if err.IsNotUndefined() {
		return nil, err
	}

	return &SyncHashSetOfValue{Set: result}, value.Undefined
}

func (h *SyncNativeHashSet[V]) UnionVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncNativeHashSet[V]:
		return h.Union(o).ToValue(), value.Undefined
	case SyncHashSet:
		return value.RefErr(h.UnionInterface(thread, o))
	case *NativeHashSet[V]:
		h.Lock()
		defer h.Unlock()

		return (&SyncNativeHashSet[V]{Set: h.Set.Union(o)}).ToValue(), value.Undefined
	case HashSet:
		h.Lock()
		defer h.Unlock()

		result, err := h.Set.UnionInterface(thread, o)
		if err.IsNotUndefined() {
			return value.Undefined, err
		}
		return (&SyncHashSetOfValue{Set: result}).ToValue(), value.Undefined
	default:
		return value.RefErr(h.unionValue(thread, other))
	}
}

func (h *SyncNativeHashSet[V]) Intersection(other *SyncNativeHashSet[V]) *SyncNativeHashSet[V] {
	h.Lock()
	defer h.Unlock()

	other.Lock()
	defer other.Unlock()

	return (&SyncNativeHashSet[V]{Set: h.Set.Intersection(other.Set)})
}

func (h *SyncNativeHashSet[V]) IntersectionInterface(thread *Thread, other SyncHashSet) (*SyncHashSetOfValue, value.Value) {
	h.Lock()
	defer h.Unlock()

	other.Lock()
	defer other.Unlock()

	result, err := h.Set.IntersectionInterface(thread, other.HashSet())
	if err.IsNotUndefined() {
		return nil, err
	}

	return &SyncHashSetOfValue{Set: result}, value.Undefined
}

func (h *SyncNativeHashSet[V]) intersectionValue(thread *Thread, other value.Value) (*HashSetOfValue, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.intersectionValue(thread, other)
}

func (h *SyncNativeHashSet[V]) IntersectionVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncNativeHashSet[V]:
		return h.Intersection(o).ToValue(), value.Undefined
	case *NativeHashSet[V]:
		return (&SyncNativeHashSet[V]{Set: h.Set.Intersection(o)}).ToValue(), value.Undefined
	case SyncHashSet:
		return value.RefErr(h.IntersectionInterface(thread, o))
	case HashSet:
		h.Lock()
		defer h.Unlock()

		result, err := h.Set.IntersectionInterface(thread, o)
		if err.IsNotUndefined() {
			return value.Undefined, err
		}
		return (&SyncHashSetOfValue{Set: result}).ToValue(), value.Undefined
	default:
		return value.RefErr(h.intersectionValue(thread, other))
	}
}

func (h *SyncNativeHashSet[V]) EqualNative(other *SyncNativeHashSet[V]) bool {
	h.Lock()
	defer h.Unlock()

	other.Lock()
	defer other.Unlock()

	return h.Set.EqualNative(other.Set)
}

func (h *SyncNativeHashSet[V]) EqualInterface(thread *Thread, other SyncHashSet) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	other.Lock()
	defer other.Unlock()

	return h.Set.EqualInterface(thread, other.HashSet())
}

func (h *SyncNativeHashSet[V]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncNativeHashSet[V]:
		return h.EqualNative(o), value.Undefined
	case SyncHashSet:
		return h.EqualInterface(thread, o)
	default:
		return false, value.Undefined
	}
}

func (h *SyncNativeHashSet[V]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.LaxEqual(thread, other)
}

func (h *SyncNativeHashSet[V]) Append(val V) bool {
	h.Lock()
	defer h.Unlock()

	return h.Set.Append(val)
}

func (h *SyncNativeHashSet[V]) AppendVal(thread *Thread, other value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.AppendVal(thread, other)
}

func (h *SyncNativeHashSet[V]) ContainsNative(val V) bool {
	h.Lock()
	defer h.Unlock()

	return h.Set.ContainsNative(val)
}

func (h *SyncNativeHashSet[V]) Contains(thread *Thread, other value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Set.Contains(thread, other)
}

type SyncNativeHashSetIterator[V value.ComparableValueInterface] struct {
	HashSet  *SyncNativeHashSet[V]
	index    int
	snapshot []V
	version  int
}

var _ value.NativeResettableIterator = &SyncNativeHashSetIterator[value.String]{}

func NewSyncNativeHashSetIterator[V value.ComparableValueInterface](set *SyncNativeHashSet[V]) *SyncNativeHashSetIterator[V] {
	iterator := &SyncNativeHashSetIterator[V]{
		HashSet: set,
		version: set.Set.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *SyncNativeHashSetIterator[V]) captureSnapshot() {
	h.HashSet.Lock()
	defer h.HashSet.Unlock()

	snapshot := make([]V, 0, h.HashSet.Set.Length())
	for element := range h.HashSet.Set.AllNative() {
		snapshot = append(snapshot, element)
	}
	h.snapshot = snapshot
}

func (*SyncNativeHashSetIterator[V]) Class() *value.Class {
	return value.SyncHashSetIteratorClass
}

func (*SyncNativeHashSetIterator[V]) DirectClass() *value.Class {
	return value.SyncHashSetIteratorClass
}

func (*SyncNativeHashSetIterator[V]) SingletonClass() *value.Class {
	return nil
}

func (h *SyncNativeHashSetIterator[V]) Copy() value.Reference {
	return &SyncNativeHashSetIterator[V]{
		HashSet:  h.HashSet,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *SyncNativeHashSetIterator[V]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *SyncNativeHashSetIterator[V]) Error() string {
	return h.Inspect()
}

func (h *SyncNativeHashSetIterator[V]) Inspect() string {
	return fmt.Sprintf("Std::Sync::HashSet::Iterator{hash_set: %s}", h.HashSet.Inspect())
}

func (*SyncNativeHashSetIterator[V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncNativeHashSetIterator[V]) Next() (v V, err value.Value) {
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

func (h *SyncNativeHashSetIterator[V]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *SyncNativeHashSetIterator[V]) Reset() {
	h.index = 0
}
