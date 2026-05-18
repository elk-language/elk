package vm

import (
	"fmt"
	"iter"
	"maps"
	"strings"
	"unsafe"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type NativeHashSet[V value.ComparableValueInterface] struct {
	m       map[V]struct{}
	version int // version of the set, each mutation increments this counter, used for guarding against concurrent mutation during iteration
}

var _ HashSet = &NativeHashSet[value.String]{}

// UNSAFE! Cast a set with native go types to a map with corresponding Elk types.
// This is EXTREMELY unsafe, use it only if `I` and `O` have the same
// underlying types eg. `unsafeCastNativeSet[string, value.String](s)`, this will convert `map[string]struct{}` to `map[value.String]struct{}`
func unsafeCastNativeSet[
	I comparable,
	O value.ComparableValueInterface,
](m map[I]struct{}) map[O]struct{} {
	return *(*map[O]struct{})(unsafe.Pointer(&m))
}

// UNSAFE! Cast a map with native go types to a new Elk `NativeHashSet` with corresponding Elk types.
// This is EXTREMELY unsafe, use it only if `I` and `O` have the same
// underlying types eg. `NewUnsafeCastNativeHashSet[string, value.String](s)`, this will convert `map[string]struct{}` to `*NativeHashSet[value.String]`
func NewUnsafeCastNativeHashSet[
	I comparable,
	O value.ComparableValueInterface,
](m map[I]struct{}) *NativeHashSet[O] {
	return &NativeHashSet[O]{
		m: unsafeCastNativeSet[I, O](m),
	}
}

// Transform a map with native go types to a new Elk `NativeHashMap` with corresponding Elk types
// using the given function.
// eg.
//
//	TransformMapIntoNativeHashSet(m, func(v string (value.String) {
//		return value.String(k)
//	})
func TransformMapIntoNativeHashSet[
	I comparable,
	O value.ComparableValueInterface,
](
	m map[I]struct{},
	fn func(v I) O,
) *NativeHashSet[O] {
	newSet := NewNativeHashSet[O](len(m))
	for v := range m {
		ov := fn(v)
		newSet.m[ov] = struct{}{}
	}
	return newSet
}

func NewNativeHashSet[V value.ComparableValueInterface](capacity int) *NativeHashSet[V] {
	return &NativeHashSet[V]{
		m: make(map[V]struct{}, capacity),
	}
}

func NewNativeHashSetWithElements[V value.ComparableValueInterface](elements ...V) *NativeHashSet[V] {
	m := make(map[V]struct{}, len(elements))
	for _, element := range elements {
		m[element] = struct{}{}
	}

	return &NativeHashSet[V]{
		m: m,
	}
}

func NewNativeHashSetWithElementsAndTotalCapacity[V value.ComparableValueInterface](capacity int, elements ...V) *NativeHashSet[V] {
	m := make(map[V]struct{}, capacity)
	for _, element := range elements {
		m[element] = struct{}{}
	}

	return &NativeHashSet[V]{
		m: m,
	}
}

func (h *NativeHashSet[V]) CloneHashSet(thread *Thread, capacity int) (HashSet, value.Value) {
	newSet := NewNativeHashSet[V](capacity)
	maps.Copy(newSet.m, h.m)
	return newSet, value.Undefined
}

func (h *NativeHashSet[V]) NewHashSet(capacity int) HashSet {
	return NewNativeHashSet[V](capacity)
}

func (h *NativeHashSet[V]) All() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		for element := range h.m {
			if !yield(element.ToValue()) {
				return
			}
		}
	}
}

func (h *NativeHashSet[V]) AllNative() iter.Seq[V] {
	return func(yield func(V) bool) {
		for element := range h.m {
			if !yield(element) {
				return
			}
		}
	}
}

func (h *NativeHashSet[V]) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		originalVersion := h.version
		for element := range h.m {
			if originalVersion != h.version {
				yield(value.Undefined, value.NewMutationDuringIterationError(h.Class().Name).ToValue())
				return
			}

			if !yield(element.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (*NativeHashSet[V]) Class() *value.Class {
	return value.HashSetClass
}

func (*NativeHashSet[V]) DirectClass() *value.Class {
	return value.HashSetClass
}

func (*NativeHashSet[V]) SingletonClass() *value.Class {
	return nil
}

func (h *NativeHashSet[V]) Iter() value.NativeIterator {
	return NewNativeHashSetIterator(h)
}

func (h *NativeHashSet[V]) IterSet() value.NativeResettableIterator {
	return NewNativeHashSetIterator(h)
}

func (h *NativeHashSet[V]) Copy() value.Reference {
	return &NativeHashSet[V]{
		m: maps.Clone(h.m),
	}
}

func (h *NativeHashSet[V]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *NativeHashSet[V]) Error() string {
	return h.Inspect()
}

func (h *NativeHashSet[V]) Inspect() string {
	var hasMultilineElements bool
	elementStrings := make(
		[]string,
		0,
		min(MAX_HASH_SET_ELEMENTS_IN_INSPECT, h.Length()),
	)

	var i int
	for element := range h.m {
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

func (*NativeHashSet[V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *NativeHashSet[V]) Length() int {
	return len(h.m)
}

func (h *NativeHashSet[V]) Remove(other V) (removed bool) {
	_, present := h.m[other]
	if !present {
		return false
	}

	delete(h.m, other)
	return true
}

func (h *NativeHashSet[V]) RemoveVal(thread *Thread, other value.Value) (removed bool, err value.Value) {
	o, ok := value.Downcast[V](other)
	if !ok {
		return false, value.NewInvalidElementInTypedSet(h, other.Class()).ToValue()
	}
	return h.Remove(o), value.Undefined
}

func (h *NativeHashSet[V]) Union(other *NativeHashSet[V]) *NativeHashSet[V] {
	newSet := make(map[V]struct{}, h.Length()+other.Length())

	maps.Copy(newSet, h.m)
	maps.Copy(newSet, other.m)

	return &NativeHashSet[V]{
		m: newSet,
	}
}

func (h *NativeHashSet[V]) UnionInterface(thread *Thread, other HashSet) (result *HashSetOfValue, err value.Value) {
	newSet := NewHashSetOfValue(h.Length() + other.Length())

	for v := range h.All() {
		_, err := newSet.AppendVal(thread, v)
		if !err.IsUndefined() {
			return nil, err
		}
	}

	for v := range other.All() {
		_, err := newSet.AppendVal(thread, v)
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return newSet, value.Undefined
}

func (h *NativeHashSet[V]) UnionVal(thread *Thread, other value.Value) (result value.Value, err value.Value) {
	switch other := other.SafeAsReference().(type) {
	case *NativeHashSet[V]:
		return h.Union(other).ToValue(), value.Undefined
	case HashSet:
		return value.RefErr(h.UnionInterface(thread, other))
	default:
		return value.Undefined, value.NewCoerceError(value.HashSetClass, other.Class()).ToValue()
	}
}

func (h *NativeHashSet[V]) Intersection(other *NativeHashSet[V]) *NativeHashSet[V] {
	newSet := make(map[V]struct{}, 5)

	for hVal := range h.AllNative() {
		_, present := other.m[hVal]
		if present {
			newSet[hVal] = struct{}{}
		}
	}

	return &NativeHashSet[V]{
		m: newSet,
	}
}

func (h *NativeHashSet[V]) IntersectionInterface(thread *Thread, other HashSet) (result *HashSetOfValue, err value.Value) {
	newSet := NewHashSetOfValue(5)

	for v := range h.All() {
		contains, err := other.Contains(thread, v)
		if !err.IsUndefined() {
			return nil, err
		}
		if !contains {
			continue
		}
		_, err = newSet.AppendVal(thread, v)
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return newSet, value.Undefined
}

func (h *NativeHashSet[V]) IntersectionVal(thread *Thread, other value.Value) (result value.Value, err value.Value) {
	switch other := other.SafeAsReference().(type) {
	case *NativeHashSet[V]:
		return h.Intersection(other).ToValue(), value.Undefined
	case HashSet:
		return value.RefErr(h.IntersectionInterface(thread, other))
	default:
		return value.Undefined, value.NewCoerceError(value.HashSetClass, other.Class()).ToValue()
	}
}

func (h *NativeHashSet[V]) EqualNative(other *NativeHashSet[V]) bool {
	return maps.Equal(h.m, other.m)
}

func (h *NativeHashSet[V]) EqualInterface(thread *Thread, other HashSet) (bool, value.Value) {
	if h.Length() != other.Length() {
		return false, value.Undefined
	}

	for hVal := range h.All() {
		contains, err := other.Contains(thread, hVal)
		if !err.IsUndefined() {
			return false, err
		}
		if !contains {
			return false, value.Undefined
		}
	}

	return true, value.Undefined
}

func (h *NativeHashSet[V]) Equal(thread *Thread, other value.Value) (result bool, err value.Value) {
	switch other := other.SafeAsReference().(type) {
	case *NativeHashSet[V]:
		return h.EqualNative(other), value.Undefined
	case HashSet:
		return h.EqualInterface(thread, other)
	default:
		return false, value.NewCoerceError(value.HashSetClass, other.Class()).ToValue()
	}
}

func (h *NativeHashSet[V]) Append(val V) (result bool) {
	_, present := h.m[val]
	if present {
		return false
	}

	h.m[val] = struct{}{}
	return true
}

func (h *NativeHashSet[V]) AppendVal(thread *Thread, other value.Value) (result bool, err value.Value) {
	o, ok := value.Downcast[V](other)
	if !ok {
		return false, value.NewInvalidElementInTypedSet(h, other.Class()).ToValue()
	}
	return h.Append(o), value.Undefined
}

func (h *NativeHashSet[V]) ContainsNative(val V) (result bool) {
	_, present := h.m[val]
	return present
}

func (h *NativeHashSet[V]) Contains(thread *Thread, other value.Value) (result bool, err value.Value) {
	o, ok := value.Downcast[V](other)
	if !ok {
		return false, value.Undefined
	}
	return h.ContainsNative(o), value.Undefined
}

type NativeHashSetIterator[V value.ComparableValueInterface] struct {
	HashSet  *NativeHashSet[V]
	index    int
	snapshot []V
	version  int
}

var _ value.NativeResettableIterator = &HashSetOfValueIterator{}

func NewNativeHashSetIterator[V value.ComparableValueInterface](set *NativeHashSet[V]) *NativeHashSetIterator[V] {
	iterator := &NativeHashSetIterator[V]{
		HashSet: set,
		version: set.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *NativeHashSetIterator[V]) captureSnapshot() {
	snapshot := make([]V, 0, h.HashSet.Length())
	for v := range h.HashSet.AllNative() {
		snapshot = append(snapshot, v)
	}
	h.snapshot = snapshot
}

func (*NativeHashSetIterator[V]) Class() *value.Class {
	return value.HashSetIteratorClass
}

func (*NativeHashSetIterator[V]) DirectClass() *value.Class {
	return value.HashSetIteratorClass
}

func (*NativeHashSetIterator[V]) SingletonClass() *value.Class {
	return nil
}

func (h *NativeHashSetIterator[V]) Copy() value.Reference {
	return &NativeHashSetIterator[V]{
		HashSet: h.HashSet,
		index:   h.index,
	}
}

func (i *NativeHashSetIterator[V]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *NativeHashSetIterator[V]) Error() string {
	return h.Inspect()
}

func (h *NativeHashSetIterator[V]) Inspect() string {
	return fmt.Sprintf("Std::HashSet::Iterator{hash_set: %s}", h.HashSet.Inspect())
}

func (*NativeHashSetIterator[V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *NativeHashSetIterator[V]) Next() (v V, err value.Value) {
	if h.version != h.HashSet.version {
		return v, value.NewMutationDuringIterationError(h.Class().Name).ToValue()
	}
	if h.index >= len(h.snapshot) {
		return v, symbol.L_stop_iteration.ToValue()
	}

	v = h.snapshot[h.index]
	h.index++
	return v, value.Undefined
}

func (h *NativeHashSetIterator[V]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *NativeHashSetIterator[V]) Reset() {
	h.index = 0
}
