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

type OrderedHashRecordOfValue struct {
	index   *HashMapOfValue     // record of key to small int index of key-value pair in `entries`
	entries []value.PairOfValue // slice with key value pairs
}

var _ OrderedHashRecord = &OrderedHashRecordOfValue{}

func OrderedHashRecordOfValueConstructor(class *value.Class) value.Value {
	return value.Ref(&OrderedHashRecordOfValue{index: &HashMapOfValue{}})
}

func NewOrderedHashRecordOfValue(capacity int) *OrderedHashRecordOfValue {
	return &OrderedHashRecordOfValue{
		index:   NewHashMapOfValue(capacity),
		entries: make([]value.PairOfValue, 0, capacity),
	}
}

func (h *OrderedHashRecordOfValue) Keys() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		for _, pair := range h.entries {
			if !yield(pair.Key()) {
				return
			}
		}
	}
}

func (h *OrderedHashRecordOfValue) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		for _, pair := range h.entries {
			if !yield(pair) {
				return
			}
		}
	}
}

func (h *OrderedHashRecordOfValue) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		for _, pair := range h.entries {
			if !yield(pair.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (h *OrderedHashRecordOfValue) Clone() *OrderedHashRecordOfValue {
	return &OrderedHashRecordOfValue{
		index:   h.index.Clone(),
		entries: slices.Clone(h.entries),
	}
}

func (h *OrderedHashRecordOfValue) ToHashMapOfValue(thread *Thread) (*HashMapOfValue, value.Value) {
	hashMap := NewHashMapOfValue(len(h.entries))
	for _, pair := range h.entries {
		err := HashMapOfValueSet(thread, hashMap, pair.Key(), pair.Value())
		if err.IsNotUndefined() {
			return nil, err
		}
	}
	return hashMap, value.Undefined
}

func (h *OrderedHashRecordOfValue) ToHashMap(thread *Thread) (HashMap, value.Value) {
	return h.ToHashMapOfValue(thread)
}

func (h *OrderedHashRecordOfValue) ToHashRecordOfValue(thread *Thread) (*HashRecordOfValue, value.Value) {
	hashRecord := NewHashRecordOfValue(len(h.entries))
	for _, pair := range h.entries {
		err := HashRecordOfValueSet(thread, hashRecord, pair.Key(), pair.Value())
		if err.IsNotUndefined() {
			return nil, err
		}
	}
	return hashRecord, value.Undefined
}

func (h *OrderedHashRecordOfValue) ToHashRecord(thread *Thread) (HashRecord, value.Value) {
	return h.ToHashRecordOfValue(thread)
}

func (h *OrderedHashRecordOfValue) ToSyncHashMap(thread *Thread) (SyncHashMap, value.Value) {
	hashMap, err := h.ToHashMapOfValue(thread)
	if err.IsNotUndefined() {
		return nil, err
	}
	return &SyncHashMapOfValue{Map: hashMap}, value.Undefined
}

func (h *OrderedHashRecordOfValue) ToOrderedHashMapOfValue() *OrderedHashMapOfValue {
	return &OrderedHashMapOfValue{
		index:   h.index.Clone(),
		entries: slices.Clone(h.entries),
	}
}

func (h *OrderedHashRecordOfValue) ToOrderedHashMap() OrderedHashMap {
	return h.ToOrderedHashMapOfValue()
}

func (h *OrderedHashRecordOfValue) ToSyncOrderedHashMap() SyncOrderedHashMap {
	return &SyncOrderedHashMapOfValue{Map: h.ToOrderedHashMapOfValue()}
}

func (h *OrderedHashRecordOfValue) CloneOrderedHashRecord(thread *Thread, capacity int) (OrderedHashRecord, value.Value) {
	newIndex := NewHashMapOfValue(capacity)
	err := HashMapOfValueCopy(thread, newIndex, h.index)
	if !err.IsUndefined() {
		return nil, err
	}
	entries := make([]value.PairOfValue, len(h.entries), capacity)
	copy(entries, h.entries)

	result := &OrderedHashRecordOfValue{
		index:   newIndex,
		entries: entries,
	}
	return result, value.Undefined
}

func (h *OrderedHashRecordOfValue) NewOrderedHashRecord(capacity int) OrderedHashRecord {
	return NewOrderedHashRecordOfValue(capacity)
}

func (h *OrderedHashRecordOfValue) IterNative() *OrderedHashRecordOfValueIterator {
	return NewOrderedHashRecordOfValueIterator(h)
}

func (h *OrderedHashRecordOfValue) Iter() value.NativeIterator {
	return h.IterNative()
}

func (h *OrderedHashRecordOfValue) IterOrderedRecord() value.NativeResettableIterator {
	return h.IterNative()
}

func (h *OrderedHashRecordOfValue) GetIndex(thread *Thread, key value.Value) (value.SmallInt, value.Value) {
	indexVal, err := h.index.GetValUndefined(thread, key)
	if err.IsNotUndefined() {
		return 0, err
	}
	if indexVal.IsUndefined() {
		return -1, value.Undefined
	}

	return indexVal.AsSmallInt(), value.Undefined
}

func (h *OrderedHashRecordOfValue) At(index int) *value.PairOfValue {
	return &h.entries[index]
}

func (h *OrderedHashRecordOfValue) AtInterface(index int) value.Pair {
	return h.At(index)
}

func (h *OrderedHashRecordOfValue) GetAt(index int) (result *value.PairOfValue, err value.Value) {
	index, err = value.NormalizeArrayIndex(index, h.Length())
	if err.IsNotUndefined() {
		return nil, err
	}

	return &h.entries[index], value.Undefined
}

// Normalizes the index so that `-1` becomes `length-1`.
// Does bounds checking and returns (nil, undefined) if the index is out of range.
func (h *OrderedHashRecordOfValue) TryAtVal(index value.Value) (value.Pair, value.Value) {
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

func (h *OrderedHashRecordOfValue) GetAtInterface(index int) (result value.Pair, err value.Value) {
	result, err = h.GetAt(index)
	return result, err
}

func (h *OrderedHashRecordOfValue) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
	index, err := h.GetIndex(thread, key)
	if err.IsNotUndefined() {
		return value.Undefined, err
	}
	if index == -1 {
		return value.Undefined, value.Undefined
	}
	return h.entries[index].Value(), value.Undefined
}

func (h *OrderedHashRecordOfValue) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
	index, err := h.GetIndex(thread, key)
	if err.IsNotUndefined() {
		return value.Undefined, err
	}
	if index == -1 {
		return value.Nil, value.Undefined
	}
	return h.entries[index].Value(), value.Undefined
}

func (h *OrderedHashRecordOfValue) SetVal(thread *Thread, key, val value.Value) value.Value {
	existingIndex, err := h.index.GetValUndefined(thread, key)
	if err.IsNotUndefined() {
		return err
	}

	if existingIndex.IsUndefined() {
		// new key
		entryIndex := len(h.entries)
		h.entries = append(h.entries, value.MakePairOfValue(key, val))
		err := h.index.SetVal(thread, key, value.SmallInt(entryIndex).ToValue())
		if err.IsNotUndefined() {
			return err
		}
	} else {
		// key is already present
		entryIndex := existingIndex.AsSmallInt()
		h.entries[entryIndex] = value.MakePairOfValue(key, val)
	}
	return value.Undefined
}

func (h *OrderedHashRecordOfValue) MergeVal(thread *Thread, other value.Value) (OrderedHashRecord, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedHashRecordOfValue:
		result := NewOrderedHashRecordOfValue(h.Length() + o.Length())

		for _, entry := range h.entries {
			result.SetVal(thread, entry.Key(), entry.Value())
		}
		for _, entry := range o.entries {
			result.SetVal(thread, entry.Key(), entry.Value())
		}

		return result, value.Undefined
	case Record:
		result := NewOrderedHashRecordOfValue(h.Length() + o.Length())

		for _, entry := range h.entries {
			result.SetVal(thread, entry.Key(), entry.Value())
		}
		for entry := range o.All() {
			result.SetVal(thread, entry.Key(), entry.Value())
		}

		return result, value.Undefined
	default:
		result := NewOrderedHashRecordOfValue(h.Length())

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

func (h *OrderedHashRecordOfValue) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
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

func (h *OrderedHashRecordOfValue) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
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

func (h *OrderedHashRecordOfValue) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	val, err := h.GetValUndefined(thread, key)
	if err.IsNotUndefined() {
		return false, err
	}
	if val.IsUndefined() {
		return false, value.Undefined
	}
	return true, value.Undefined
}

func (h *OrderedHashRecordOfValue) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedHashRecordOfValue:
		return PairSliceEqual(thread, h.entries, o.entries)
	case OrderedHashRecord:
		otherEntries := slices.Collect(o.All())
		return PairSliceEqual(thread, h.entries, otherEntries)
	}

	return false, value.Undefined
}

func (h *OrderedHashRecordOfValue) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case Record:
		return OrderedHashRecordLaxEqualRecord(thread, h, o)
	default:
		if !value.IsA(other, value.MapMixin) {
			return false, value.Undefined
		}

		for element, err := range Iterate(thread, other) {
			if err.IsNotUndefined() {
				return false, err
			}

			pair := element.AsReference().(value.Pair)
			val, err := h.GetValUndefined(thread, pair.Key())
			if err.IsNotUndefined() {
				return false, err
			}
			if val.IsUndefined() {
				return false, value.Undefined
			}

			eq, err := Equal(thread, val, pair.Value())
			if err.IsNotUndefined() {
				return false, err
			}

			if value.Falsy(eq) {
				return false, value.Undefined
			}
		}

		return true, value.Undefined
	}
}

func (*OrderedHashRecordOfValue) Class() *value.Class {
	return value.OrderedHashRecordClass
}

func (*OrderedHashRecordOfValue) DirectClass() *value.Class {
	return value.OrderedHashRecordClass
}

func (*OrderedHashRecordOfValue) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedHashRecordOfValue) Copy() value.Reference {
	return h.Clone()
}

func (h *OrderedHashRecordOfValue) ToValue() value.Value {
	return value.Ref(h)
}

func (h *OrderedHashRecordOfValue) Error() string {
	return h.Inspect()
}

func (h *OrderedHashRecordOfValue) Inspect() string {
	var hasMultilineElements bool
	keyStrings := make(
		[]string,
		0,
		min(MAX_HASH_RECORD_ELEMENTS_IN_INSPECT, h.Length()),
	)
	valStrings := make(
		[]string,
		0,
		min(MAX_HASH_RECORD_ELEMENTS_IN_INSPECT, h.Length()),
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

		if i >= MAX_HASH_RECORD_ELEMENTS_IN_INSPECT-1 {
			break
		}
		i++
	}

	var buff strings.Builder

	buff.WriteString("Std::OrderedHashRecord(")
	if hasMultilineElements {
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

			if i >= MAX_HASH_RECORD_ELEMENTS_IN_INSPECT-1 {
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

			if i >= MAX_HASH_RECORD_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(", ...")
				break
			}
		}
	}
	buff.WriteRune(')')

	return buff.String()
}

func (*OrderedHashRecordOfValue) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedHashRecordOfValue) Capacity() int {
	return cap(h.entries)
}

func (h *OrderedHashRecordOfValue) LeftCapacity() int {
	return h.Capacity() - h.Length()
}

func (h *OrderedHashRecordOfValue) Length() int {
	return len(h.entries)
}

type OrderedHashRecordOfValueIterator struct {
	HashRecord *OrderedHashRecordOfValue
	Index      int
	snapshot   []value.PairOfValue
}

var _ value.NativeResettableIterator = &OrderedHashRecordOfValueIterator{}

func NewOrderedHashRecordOfValueIterator(hrec *OrderedHashRecordOfValue) *OrderedHashRecordOfValueIterator {
	iterator := &OrderedHashRecordOfValueIterator{
		HashRecord: hrec,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *OrderedHashRecordOfValueIterator) captureSnapshot() {
	snapshot := make([]value.PairOfValue, 0, h.HashRecord.Length())
	for pair := range h.HashRecord.All() {
		snapshot = append(snapshot, pair)
	}
	h.snapshot = snapshot
}

func (*OrderedHashRecordOfValueIterator) Class() *value.Class {
	return value.OrderedHashRecordIteratorClass
}

func (*OrderedHashRecordOfValueIterator) DirectClass() *value.Class {
	return value.OrderedHashRecordIteratorClass
}

func (*OrderedHashRecordOfValueIterator) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedHashRecordOfValueIterator) Copy() value.Reference {
	return &OrderedHashRecordOfValueIterator{
		HashRecord: h.HashRecord,
		Index:      h.Index,
		snapshot:   h.snapshot,
	}
}

func (i *OrderedHashRecordOfValueIterator) ToValue() value.Value {
	return value.Ref(i)
}

func (h *OrderedHashRecordOfValueIterator) Error() string {
	return h.Inspect()
}

func (h *OrderedHashRecordOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::OrderedHashRecord::Iterator{hash_record: %s}", h.HashRecord.Inspect())
}

func (*OrderedHashRecordOfValueIterator) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedHashRecordOfValueIterator) NextValue() (value.Value, value.Value) {
	if h.Index >= len(h.snapshot) {
		return value.Undefined, value.S(symbol.L_stop_iteration).ToValue()
	}

	pair := &h.snapshot[h.Index]
	h.Index++
	return value.Ref(pair), value.Undefined
}

func (h *OrderedHashRecordOfValueIterator) Reset() {
	h.Index = 0
}
