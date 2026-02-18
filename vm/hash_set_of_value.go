package vm

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
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
				yield(value.Undefined, value.NewMutationDuringIterationError(h.Class().Name).ToValue())
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

func (h *HashSetOfValue) Copy() value.Reference {
	newTable := slices.Clone(h.table)
	return &HashSetOfValue{
		table:         newTable,
		occupiedSlots: h.occupiedSlots,
		elements:      h.elements,
	}
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

func (h *HashSetOfValue) UnionVal(thread *Thread, other value.Value) (result value.Value, err value.Value) {
	switch other := other.SafeAsReference().(type) {
	case *HashSetOfValue:
		return value.RefErr(HashSetOfValueUnion(thread, h, other))
	case HashSet:
		return value.RefErr(HashSetOfValueUnionInterface(thread, h, other))
	default:
		return value.Undefined, value.NewCoerceError(value.HashSetClass, other.Class()).ToValue()
	}
}

func (h *HashSetOfValue) IntersectionVal(thread *Thread, other value.Value) (result value.Value, err value.Value) {
	switch other := other.SafeAsReference().(type) {
	case *HashSetOfValue:
		return value.RefErr(HashSetOfValueIntersection(thread, h, other))
	case HashSet:
		return value.RefErr(HashSetOfValueIntersectionInterface(thread, h, other))
	default:
		return value.Undefined, value.NewCoerceError(value.HashSetClass, other.Class()).ToValue()
	}
}

func (h *HashSetOfValue) Equal(thread *Thread, other value.Value) (result bool, err value.Value) {
	switch other := other.SafeAsReference().(type) {
	case *HashSetOfValue:
		return HashSetOfValueEqual(thread, h, other)
	case HashSet:
		return HashSetOfValueEqualInterface(thread, h, other)
	default:
		return false, value.NewCoerceError(value.HashSetClass, other.Class()).ToValue()
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
		return value.Undefined, value.NewMutationDuringIterationError(h.HashSet.Class().Name).ToValue()
	}

	for {
		if h.index >= h.HashSet.Capacity() {
			return value.Undefined, stopIterationSymbol.ToValue()
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
