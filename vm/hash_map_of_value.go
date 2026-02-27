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

const HashMapOfValueMaxLoad = 0.75

type HashMapOfValue struct {
	Table         []value.PairOfValue // underlying data container, it's `len` is always equal to it's `cap` and it also serves as the capacity of the HashMap
	OccupiedSlots int                 // number of slots taken by active pairs and those left by deleted pairs
	Elements      int                 // number of slots occupied by active pairs
	version       int                 // version of the map, each mutation increments this counter, used for guarding against concurrent mutation during iteration
}

var _ HashMap = &HashMapOfValue{}

func HashMapOfValueConstructor(class *value.Class) value.Value {
	return value.Ref(&HashMapOfValue{})
}

func NewHashMapOfValue(capacity int) *HashMapOfValue {
	return &HashMapOfValue{
		Table: make([]value.PairOfValue, capacity),
	}
}

func (h *HashMapOfValue) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		for _, pair := range h.Table {
			if pair.Key().IsUndefined() {
				continue
			}

			if !yield(pair) {
				return
			}
		}
	}
}

func (h *HashMapOfValue) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		originalVersion := h.version
		for _, pair := range h.Table {
			if pair.Key().IsUndefined() {
				continue
			}

			if originalVersion != h.version {
				yield(value.Undefined, value.NewMutationDuringIterationError(h.Class().Name).ToValue())
				return
			}

			if !yield(pair.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (h *HashMapOfValue) CloneHashMap(thread *Thread, capacity int) (HashMap, value.Value) {
	newMap := NewHashMapOfValue(capacity)
	err := HashMapOfValueCopy(thread, newMap, h)
	if !err.IsUndefined() {
		return nil, err
	}

	return newMap, value.Undefined
}

func (h *HashMapOfValue) CloneHashRecord(thread *Thread, capacity int) (HashRecord, value.Value) {
	return h.CloneHashMap(thread, capacity)
}

func (h *HashMapOfValue) IterNative() *HashMapOfValueIterator {
	return NewHashMapOfValueIterator(h)
}

func (h *HashMapOfValue) Iter() value.NativeIterator {
	return h.IterNative()
}

func (h *HashMapOfValue) IterMap() value.NativeResettableIterator {
	return h.IterNative()
}

func (h *HashMapOfValue) IterRecord() value.NativeResettableIterator {
	return h.IterNative()
}

func (h *HashMapOfValue) GetVal(thread *Thread, key value.Value) (value.Value, value.Value) {
	return HashMapOfValueGet(thread, h, key)
}

func (h *HashMapOfValue) SetVal(thread *Thread, key, val value.Value) value.Value {
	return HashMapOfValueSet(thread, h, key, val)
}

func (h *HashMapOfValue) ConcatVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashMapOfValue:
		return value.RefErr(HashMapOfValueConcat(thread, h, o))
	case HashRecord:
		return value.RefErr(HashMapOfValueConcatInterface(thread, h, o))
	}

	return value.Undefined, value.Ref(value.Errorf(value.TypeErrorClass, "cannot concat %s with map %s", other.Inspect(), h.Inspect()))
}

func (h *HashMapOfValue) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
	return HashMapOfValueContains(thread, h, other)
}

func (h *HashMapOfValue) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	return HashMapOfValueContainsValue(thread, h, val)
}

func (h *HashMapOfValue) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	return HashMapOfValueContainsKey(thread, h, key)
}

func (h *HashMapOfValue) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashMapOfValue:
		return HashMapOfValueEqual(thread, h, o)
	case HashMap:
		return HashMapOfValueEqualInterface(thread, h, o)
	}

	return false, value.Undefined
}

func (h *HashMapOfValue) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashMapOfValue:
		return HashMapOfValueLaxEqual(thread, h, o)
	case HashMap:
		return HashMapOfValueLaxEqualInterface(thread, h, o)
	}

	return false, value.Undefined
}

func (h *HashMapOfValue) Grow(thread *Thread, newSlots int) value.Value {
	return HashMapOfValueGrow(thread, h, newSlots)
}

func (*HashMapOfValue) Class() *value.Class {
	return value.HashMapClass
}

func (*HashMapOfValue) DirectClass() *value.Class {
	return value.HashMapClass
}

func (*HashMapOfValue) SingletonClass() *value.Class {
	return nil
}

func (h *HashMapOfValue) Clone() *HashMapOfValue {
	newTable := slices.Clone(h.Table)
	return &HashMapOfValue{
		Table:         newTable,
		OccupiedSlots: h.OccupiedSlots,
		Elements:      h.Elements,
	}
}

func (h *HashMapOfValue) Copy() value.Reference {
	return h.Clone()
}

func (h *HashMapOfValue) ToValue() value.Value {
	return value.Ref(h)
}

func (h *HashMapOfValue) Error() string {
	return h.Inspect()
}

const MAX_HASH_MAP_ELEMENTS_IN_INSPECT = 300

func (h *HashMapOfValue) Inspect() string {
	var hasMultilineElements bool
	keyStrings := make(
		[]string,
		0,
		min(MAX_HASH_MAP_ELEMENTS_IN_INSPECT, h.Length()),
	)
	valStrings := make(
		[]string,
		0,
		min(MAX_HASH_MAP_ELEMENTS_IN_INSPECT, h.Length()),
	)

	i := 0
	for _, entry := range h.Table {
		if entry.Key().IsUndefined() {
			continue
		}

		keyString := entry.Key().Inspect()
		keyStrings = append(keyStrings, keyString)

		valString := entry.Value().Inspect()
		valStrings = append(valStrings, valString)

		if strings.ContainsRune(keyString, '\n') ||
			strings.ContainsRune(valString, '\n') {
			hasMultilineElements = true
		}

		if i >= MAX_HASH_MAP_ELEMENTS_IN_INSPECT-1 {
			break
		}
		i++
	}

	var buff strings.Builder

	buff.WriteRune('{')
	if hasMultilineElements || h.Length() > 15 {
		buff.WriteRune('\n')
		for i := range len(keyStrings) {
			keyString := keyStrings[i]
			valString := valStrings[i]

			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, keyString, 1)
			buff.WriteString(" => ")
			indent.IndentStringFromSecondLine(&buff, valString, 1)

			if i >= MAX_HASH_MAP_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(",\n  ...")
				break
			}
		}
		buff.WriteRune('\n')
	} else {
		for i := range len(keyStrings) {
			keyString := keyStrings[i]
			valString := valStrings[i]

			if i != 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(keyString)
			buff.WriteString(" => ")
			buff.WriteString(valString)

			if i >= MAX_HASH_MAP_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(", ...")
				break
			}
		}
	}
	buff.WriteRune('}')

	return buff.String()
}

func (*HashMapOfValue) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *HashMapOfValue) Capacity() int {
	return len(h.Table)
}

func (h *HashMapOfValue) LeftCapacity() int {
	return h.Capacity() - h.Length()
}

func (h *HashMapOfValue) Length() int {
	return h.Elements
}

type HashMapOfValueIterator struct {
	HashMap *HashMapOfValue
	Index   int
	version int
}

var _ value.NativeResettableIterator = &HashMapOfValueIterator{}

func NewHashMapOfValueIterator(hmap *HashMapOfValue) *HashMapOfValueIterator {
	return &HashMapOfValueIterator{
		HashMap: hmap,
		version: hmap.version,
	}
}

func (*HashMapOfValueIterator) Class() *value.Class {
	return value.HashMapIteratorClass
}

func (*HashMapOfValueIterator) DirectClass() *value.Class {
	return value.HashMapIteratorClass
}

func (*HashMapOfValueIterator) SingletonClass() *value.Class {
	return nil
}

func (h *HashMapOfValueIterator) Copy() value.Reference {
	return &HashMapOfValueIterator{
		HashMap: h.HashMap,
		Index:   h.Index,
		version: h.version,
	}
}

func (i *HashMapOfValueIterator) ToValue() value.Value {
	return value.Ref(i)
}

func (h *HashMapOfValueIterator) Error() string {
	return h.Inspect()
}

func (h *HashMapOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::HashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*HashMapOfValueIterator) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *HashMapOfValueIterator) NextValue() (value.Value, value.Value) {
	if h.version != h.HashMap.version {
		return value.Undefined, value.NewMutationDuringIterationError(h.HashMap.Class().Name).ToValue()
	}

	for {
		if h.Index >= h.HashMap.Capacity() {
			return value.Undefined, symbol.L_stop_iteration.ToValue()
		}

		pair := h.HashMap.Table[h.Index]
		h.Index++
		if !pair.Key().IsUndefined() {
			return value.Ref(&h.HashMap.Table[h.Index-1]), value.Undefined
		}
	}
}

func (h *HashMapOfValueIterator) Reset() {
	h.Index = 0
	h.version = h.HashMap.version
}
