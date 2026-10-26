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

type OrderedHashMapOfValue struct {
	index   *HashMapOfValue     // map of key to small int index of key-value pair in `entries`
	entries []value.PairOfValue // slice with key value pairs
}

var _ OrderedHashMap = &OrderedHashMapOfValue{}

func OrderedHashMapOfValueConstructor(class *value.Class) value.Value {
	return value.Ref(&OrderedHashMapOfValue{index: &HashMapOfValue{}})
}

func NewOrderedHashMapOfValue(capacity int) *OrderedHashMapOfValue {
	return &OrderedHashMapOfValue{
		index:   NewHashMapOfValue(capacity),
		entries: make([]value.PairOfValue, 0, capacity),
	}
}

func NewOrderedHashMapOfValueFromMap(m *HashMapOfValue) *OrderedHashMapOfValue {
	entries := make([]value.PairOfValue, m.Length())
	i := 0
	for pair := range m.All() {
		entries[i] = pair
		i++
	}

	return &OrderedHashMapOfValue{
		index:   m,
		entries: entries,
	}
}

func (h *OrderedHashMapOfValue) Keys() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		for _, pair := range h.entries {
			if !yield(pair.Key()) {
				return
			}
		}
	}
}

func (h *OrderedHashMapOfValue) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		for _, pair := range h.entries {
			if !yield(pair) {
				return
			}
		}
	}
}

func (h *OrderedHashMapOfValue) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		originalVersion := h.index.version
		for _, pair := range h.entries {
			if originalVersion != h.index.version {
				yield(value.Undefined, value.NewMutationDuringIterationError(h.Class().Name()).ToValue())
				return
			}

			if !yield(pair.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (h *OrderedHashMapOfValue) Clone() *OrderedHashMapOfValue {
	return &OrderedHashMapOfValue{
		index:   h.index.Clone(),
		entries: slices.Clone(h.entries),
	}
}

func (h *OrderedHashMapOfValue) CloneOrderedHashMap(thread *Thread, capacity int) (OrderedHashMap, value.Value) {
	newMap := NewHashMapOfValue(capacity)
	err := HashMapOfValueCopy(thread, newMap, h.index)
	if !err.IsUndefined() {
		return nil, err
	}
	entries := make([]value.PairOfValue, len(h.entries), capacity)
	copy(entries, h.entries)

	result := &OrderedHashMapOfValue{
		index:   newMap,
		entries: entries,
	}
	return result, value.Undefined
}

func (h *OrderedHashMapOfValue) NewOrderedHashMap(capacity int) OrderedHashMap {
	return NewOrderedHashMapOfValue(capacity)
}

func (h *OrderedHashMapOfValue) NewOrderedHashRecord(capacity int) OrderedHashRecord {
	return h.NewOrderedHashMap(capacity)
}

func (h *OrderedHashMapOfValue) CloneOrderedHashRecord(thread *Thread, capacity int) (OrderedHashRecord, value.Value) {
	return h.CloneOrderedHashMap(thread, capacity)
}

func (h *OrderedHashMapOfValue) IterNative() *OrderedHashMapOfValueIterator {
	return NewOrderedHashMapOfValueIterator(h)
}

func (h *OrderedHashMapOfValue) Iter() value.NativeIterator {
	return h.IterNative()
}

func (h *OrderedHashMapOfValue) IterOrderedMap() value.NativeResettableIterator {
	return h.IterNative()
}

func (h *OrderedHashMapOfValue) IterOrderedRecord() value.NativeResettableIterator {
	return h.IterNative()
}

func (h *OrderedHashMapOfValue) GetIndex(thread *Thread, key value.Value) (value.SmallInt, value.Value) {
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

func (h *OrderedHashMapOfValue) At(index int) *value.PairOfValue {
	return &h.entries[index]
}

func (h *OrderedHashMapOfValue) AtInterface(index int) value.Pair {
	return h.At(index)
}

func (h *OrderedHashMapOfValue) GetAt(index int) (result *value.PairOfValue, err value.Value) {
	index, err = value.NormalizeArrayIndex(index, h.Length())
	if err.IsNotUndefined() {
		return nil, err
	}

	return &h.entries[index], value.Undefined
}

func (h *OrderedHashMapOfValue) GetAtInterface(index int) (result value.Pair, err value.Value) {
	result, err = h.GetAt(index)
	return result, err
}

func (h *OrderedHashMapOfValue) TryAtVal(index value.Value) (value.Pair, value.Value) {
	var i int

	i, ok := value.ToGoInt(index)
	if !ok {
		if i == -1 {
			return nil, value.Undefined
		}
		return nil, value.NewCoerceError(value.IntClass, index.Class()).ToValue()
	}

	i, ok = value.NormalizeArrayIndexOk(i, h.Length())
	if !ok {
		return nil, value.Undefined
	}

	return &h.entries[i], value.Undefined
}

func (h *OrderedHashMapOfValue) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
	index, err := h.GetIndex(thread, key)
	if err.IsNotUndefined() {
		return value.Undefined, err
	}
	if index == -1 {
		return value.Undefined, value.Undefined
	}
	return h.entries[index].Value(), value.Undefined
}

func (h *OrderedHashMapOfValue) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
	index, err := h.GetIndex(thread, key)
	if err.IsNotUndefined() {
		return value.Undefined, err
	}
	if index == -1 {
		return value.Nil, value.Undefined
	}
	return h.entries[index].Value(), value.Undefined
}

func (h *OrderedHashMapOfValue) SetVal(thread *Thread, key, val value.Value) value.Value {
	existingIndex, err := h.index.GetValUndefined(thread, key)
	if err.IsNotUndefined() {
		return err
	}

	if existingIndex.IsUndefined() {
		// new key
		entryIndex := len(h.entries)
		h.entries = append(h.entries, value.MakePairOfValue(key, val))
		h.index.SetVal(thread, key, value.SmallInt(entryIndex).ToValue())
	} else {
		// key is already present
		entryIndex := existingIndex.AsSmallInt()
		h.entries[entryIndex] = value.MakePairOfValue(key, val)
	}
	return value.Undefined
}

func (h *OrderedHashMapOfValue) MergeVal(thread *Thread, other value.Value) (OrderedHashRecord, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedHashMapOfValue:
		result := NewOrderedHashMapOfValue(h.Length() + o.Length())

		for _, entry := range h.entries {
			result.SetVal(thread, entry.Key(), entry.Value())
		}
		for _, entry := range o.entries {
			result.SetVal(thread, entry.Key(), entry.Value())
		}

		return result, value.Undefined
	case OrderedHashRecord:
		result := NewOrderedHashMapOfValue(h.Length() + o.Length())

		for _, entry := range h.entries {
			result.SetVal(thread, entry.Key(), entry.Value())
		}
		for entry := range o.All() {
			result.SetVal(thread, entry.Key(), entry.Value())
		}

		return result, value.Undefined
	case HashRecord:
		result := NewOrderedHashMapOfValue(h.Length() + o.Length())

		for _, entry := range h.entries {
			result.SetVal(thread, entry.Key(), entry.Value())
		}
		for entry := range o.All() {
			result.SetVal(thread, entry.Key(), entry.Value())
		}

		return result, value.Undefined
	default:
		result := NewOrderedHashMapOfValue(h.Length())

		for _, entry := range h.entries {
			result.SetVal(thread, entry.Key(), entry.Value())
		}

		for element, err := range Iterate(thread, other) {
			if err.IsNotUndefined() {
				return nil, err
			}

			entry := element.AsReference().(value.Pair)
			result.SetVal(thread, entry.Key(), entry.Value())
		}

		return result, value.Undefined
	}
}

func (h *OrderedHashMapOfValue) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
	val, err := h.GetValUndefined(thread, other.Key())
	if err.IsNotUndefined() {
		return false, err
	}
	if val.IsUndefined() {
		return false, value.Undefined
	}

	eq, err := Equal(thread, val, other.Value())
	if err.IsNotUndefined() {
		return false, err
	}
	return value.Truthy(eq), value.Undefined
}

func (h *OrderedHashMapOfValue) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	for _, entry := range h.entries {
		eq, err := Equal(thread, entry.Value(), val)
		if err.IsNotUndefined() {
			return false, err
		}

		if value.Truthy(eq) {
			return true, value.Undefined
		}
	}

	return false, value.Undefined
}

func (h *OrderedHashMapOfValue) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	val, err := h.GetValUndefined(thread, key)
	if err.IsNotUndefined() {
		return false, err
	}
	if val.IsUndefined() {
		return false, value.Undefined
	}
	return true, value.Undefined
}

func (h *OrderedHashMapOfValue) RemoveKey(thread *Thread, key value.Value) (bool, value.Value) {
	index, err := h.GetIndex(thread, key)
	if err.IsNotUndefined() {
		return false, err
	}

	h.entries = slices.Delete(h.entries, int(index), int(index)+1)
	_, err = h.index.RemoveKey(thread, key)
	if err.IsNotUndefined() {
		return false, err
	}

	return true, value.Undefined
}

func ValueSliceEqual(thread *Thread, a []value.Value, b []value.Value) (bool, value.Value) {
	if len(a) != len(b) {
		return false, value.Undefined
	}

	for i := range len(a) {
		aEntry := a[i]
		bEntry := b[i]

		eq, err := Equal(thread, aEntry, bEntry)
		if err.IsNotUndefined() {
			return false, err
		}

		if value.Falsy(eq) {
			return false, value.Undefined
		}
	}

	return true, value.Undefined
}

func PairSliceEqual(thread *Thread, a []value.PairOfValue, b []value.PairOfValue) (bool, value.Value) {
	if len(a) != len(b) {
		return false, value.Undefined
	}

	for i := range len(a) {
		aEntry := a[i]
		bEntry := b[i]

		keyEq, err := Equal(thread, aEntry.Key(), bEntry.Key())
		if err.IsNotUndefined() {
			return false, err
		}

		if value.Falsy(keyEq) {
			return false, value.Undefined
		}

		valEq, err := Equal(thread, aEntry.Value(), bEntry.Value())
		if err.IsNotUndefined() {
			return false, err
		}

		if value.Falsy(valEq) {
			return false, value.Undefined
		}
	}

	return true, value.Undefined
}

func (h *OrderedHashMapOfValue) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedHashMapOfValue:
		return PairSliceEqual(thread, h.entries, o.entries)
	case OrderedHashMap:
		otherEntries := slices.Collect(o.All())
		return PairSliceEqual(thread, h.entries, otherEntries)
	}

	return false, value.Undefined
}

func (h *OrderedHashMapOfValue) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case OrderedHashRecord:
		return OrderedHashRecordLaxEqual(thread, h, o)
	case HashRecord:
		return OrderedHashRecordLaxEqualHashRecord(thread, h, o)
	}

	return false, value.Undefined
}

func (*OrderedHashMapOfValue) Class() *value.Class {
	return value.OrderedHashMapClass
}

func (*OrderedHashMapOfValue) DirectClass() *value.Class {
	return value.OrderedHashMapClass
}

func (*OrderedHashMapOfValue) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedHashMapOfValue) Copy() value.Reference {
	return h.Clone()
}

func (h *OrderedHashMapOfValue) ToValue() value.Value {
	return value.Ref(h)
}

func (h *OrderedHashMapOfValue) Error() string {
	return h.Inspect()
}

func (h *OrderedHashMapOfValue) Inspect() string {
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
	for entry := range h.All() {
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

	buff.WriteString("Std::OrderedHashMap{")
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

func (*OrderedHashMapOfValue) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedHashMapOfValue) Capacity() int {
	return cap(h.entries)
}

func (h *OrderedHashMapOfValue) LeftCapacity() int {
	return h.Capacity() - h.Length()
}

func (h *OrderedHashMapOfValue) Length() int {
	return len(h.entries)
}

type OrderedHashMapOfValueIterator struct {
	HashMap *OrderedHashMapOfValue
	Index   int
	version int
}

var _ value.NativeResettableIterator = &OrderedHashMapOfValueIterator{}

func NewOrderedHashMapOfValueIterator(hmap *OrderedHashMapOfValue) *OrderedHashMapOfValueIterator {
	return &OrderedHashMapOfValueIterator{
		HashMap: hmap,
		version: hmap.index.version,
	}
}

func (*OrderedHashMapOfValueIterator) Class() *value.Class {
	return value.OrderedHashMapIteratorClass
}

func (*OrderedHashMapOfValueIterator) DirectClass() *value.Class {
	return value.OrderedHashMapIteratorClass
}

func (*OrderedHashMapOfValueIterator) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedHashMapOfValueIterator) Copy() value.Reference {
	return &OrderedHashMapOfValueIterator{
		HashMap: h.HashMap,
		Index:   h.Index,
		version: h.version,
	}
}

func (i *OrderedHashMapOfValueIterator) ToValue() value.Value {
	return value.Ref(i)
}

func (h *OrderedHashMapOfValueIterator) Error() string {
	return h.Inspect()
}

func (h *OrderedHashMapOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::OrderedHashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*OrderedHashMapOfValueIterator) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedHashMapOfValueIterator) NextValue() (value.Value, value.Value) {
	if h.version != h.HashMap.index.version {
		return value.Undefined, value.NewMutationDuringIterationError(h.HashMap.Class().Name()).ToValue()
	}

	for {
		if h.Index >= h.HashMap.Capacity() {
			return value.Undefined, value.S(symbol.L_stop_iteration).ToValue()
		}

		pair := &h.HashMap.entries[h.Index]
		h.Index++
		if !pair.Key().IsUndefined() {
			return value.Ref(pair), value.Undefined
		}
	}
}

func (h *OrderedHashMapOfValueIterator) Reset() {
	h.Index = 0
	h.version = h.HashMap.index.version
}
