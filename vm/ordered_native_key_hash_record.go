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

type OrderedNativeKeyHashRecord[K value.ComparableValueInterface] struct {
	index   map[K]int
	entries []value.NativePair[K, value.Value] // slice with key value pairs
}

var _ OrderedHashRecord = &OrderedNativeKeyHashRecord[value.String]{}

func NewOrderedNativeKeyHashRecord[K value.ComparableValueInterface](capacity int) *OrderedNativeKeyHashRecord[K] {
	return &OrderedNativeKeyHashRecord[K]{
		index:   make(map[K]int, capacity),
		entries: make([]value.NativePair[K, value.Value], 0, capacity),
	}
}

func NewOrderedNativeKeyHashRecordWithElements[K value.ComparableValueInterface](elements ...value.NativePair[K, value.Value]) *OrderedNativeKeyHashRecord[K] {
	return NewOrderedNativeKeyHashRecordWithElementsAndTotalCapacity(len(elements), elements...)
}

func NewOrderedNativeKeyHashRecordWithElementsAndTotalCapacity[K value.ComparableValueInterface](capacity int, elements ...value.NativePair[K, value.Value]) *OrderedNativeKeyHashRecord[K] {
	m := NewOrderedNativeKeyHashRecord[K](capacity)
	for _, pair := range elements {
		m.Set(pair.NativeKey(), pair.NativeValue())
	}
	return m
}

func (h *OrderedNativeKeyHashRecord[K]) CloneOrderedHashRecord(thread *Thread, capacity int) (OrderedHashRecord, value.Value) {
	newIndex := make(map[K]int, capacity)
	maps.Copy(newIndex, h.index)

	entries := make([]value.NativePair[K, value.Value], len(h.entries), capacity)
	copy(entries, h.entries)

	result := &OrderedNativeKeyHashRecord[K]{
		index:   newIndex,
		entries: entries,
	}
	return result, value.Undefined
}

func (h *OrderedNativeKeyHashRecord[K]) NewOrderedHashRecord(capacity int) OrderedHashRecord {
	return NewOrderedNativeKeyHashRecord[K](capacity)
}

func (h *OrderedNativeKeyHashRecord[K]) Keys() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		for _, pair := range h.entries {
			if !yield(pair.Key()) {
				return
			}
		}
	}
}

func (h *OrderedNativeKeyHashRecord[K]) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		for _, pair := range h.entries {
			elkPair := value.MakePairOfValue(pair.Key(), pair.NativeValue())
			if !yield(elkPair) {
				return
			}
		}
	}
}

func (h *OrderedNativeKeyHashRecord[K]) AllNative() iter.Seq2[K, value.Value] {
	return func(yield func(K, value.Value) bool) {
		for _, pair := range h.entries {
			if !yield(pair.NativeKey(), pair.NativeValue()) {
				return
			}
		}
	}
}

func (h *OrderedNativeKeyHashRecord[K]) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		for _, pair := range h.entries {
			if !yield(pair.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (h *OrderedNativeKeyHashRecord[K]) IterOrderedRecord() value.NativeResettableIterator {
	return NewOrderedNativeKeyHashRecordIterator(h)
}

func (h *OrderedNativeKeyHashRecord[K]) Iter() value.NativeIterator {
	return h.IterOrderedRecord()
}

func (h *OrderedNativeKeyHashRecord[K]) GetIndex(key K) (int, bool) {
	index, ok := h.index[key]
	return index, ok
}

func (h *OrderedNativeKeyHashRecord[K]) GetNativePair(key K) *value.NativePair[K, value.Value] {
	i, ok := h.GetIndex(key)
	if !ok {
		return nil
	}
	return &h.entries[i]
}

func (h *OrderedNativeKeyHashRecord[K]) GetPair(key K) value.Pair {
	return h.GetNativePair(key)
}

func (h *OrderedNativeKeyHashRecord[K]) Get(key K) value.Value {
	i, ok := h.GetIndex(key)
	if !ok {
		return value.Undefined
	}
	return h.entries[i].NativeValue()
}

func (h *OrderedNativeKeyHashRecord[K]) GetOk(key K) (value.Value, bool) {
	i, ok := h.GetIndex(key)
	if !ok {
		return value.Undefined, false
	}
	return h.entries[i].NativeValue(), true
}

func (h *OrderedNativeKeyHashRecord[K]) At(index int) *value.NativePair[K, value.Value] {
	return &h.entries[index]
}

func (h *OrderedNativeKeyHashRecord[K]) AtInterface(index int) value.Pair {
	return h.At(index)
}

func (h *OrderedNativeKeyHashRecord[K]) GetAt(index int) (result *value.NativePair[K, value.Value], err value.Value) {
	index, err = value.NormalizeArrayIndex(index, h.Length())
	if err.IsNotUndefined() {
		return nil, err
	}

	return &h.entries[index], value.Undefined
}

func (h *OrderedNativeKeyHashRecord[K]) GetAtInterface(index int) (result value.Pair, err value.Value) {
	result, err = h.GetAt(index)
	return result, err
}

func (h *OrderedNativeKeyHashRecord[K]) TryAtVal(index value.Value) (value.Pair, value.Value) {
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

func (h *OrderedNativeKeyHashRecord[K]) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return value.Undefined, value.Undefined
	}
	v, ok := h.GetOk(k)
	if !ok {
		return value.Undefined, value.Undefined
	}
	return v, value.Undefined
}

func (h *OrderedNativeKeyHashRecord[K]) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return value.Nil, value.Undefined
	}
	v, ok := h.GetOk(k)
	if !ok {
		return value.Nil, value.Undefined
	}
	return v, value.Undefined
}

func (h *OrderedNativeKeyHashRecord[K]) Set(key K, val value.Value) {
	entryIndex, ok := h.GetIndex(key)
	if ok {
		// key is already present
		h.entries[entryIndex] = value.MakeNativePair(key, val)
	} else {
		// new key
		entryIndex = len(h.entries)
		h.entries = append(h.entries, value.MakeNativePair(key, val))
		h.index[key] = entryIndex
	}
}

func (h *OrderedNativeKeyHashRecord[K]) SetVal(thread *Thread, key, val value.Value) value.Value {
	k, ok := value.Downcast[K](key)
	if !ok {
		return value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}

	h.Set(k, val)
	return value.Undefined
}

func (h *OrderedNativeKeyHashRecord[K]) MergeVal(thread *Thread, other value.Value) (OrderedHashRecord, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedNativeKeyHashRecord[K]:
		result := NewOrderedNativeKeyHashRecord[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.SetVal(thread, entry.Key(), entry.NativeValue())
		}
		for _, entry := range o.entries {
			result.SetVal(thread, entry.Key(), entry.NativeValue())
		}

		return result, value.Undefined
	case *NativeKeyHashMap[K]:
		result := NewOrderedNativeKeyHashRecord[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case *OrderedNativeKeyHashMap[K]:
		result := NewOrderedNativeKeyHashRecord[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case *SyncNativeKeyHashMap[K]:
		result := NewOrderedNativeKeyHashRecord[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case *SyncOrderedNativeKeyHashMap[K]:
		result := NewOrderedNativeKeyHashRecord[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case NativeKeyHashRecord[K]:
		result := NewOrderedNativeKeyHashRecord[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case Record:
		result := NewOrderedHashRecordOfValue(h.Length() + o.Length())
		for _, entry := range h.entries {
			result.SetVal(thread, entry.Key(), entry.NativeValue())
		}
		for entry := range o.All() {
			result.SetVal(thread, entry.Key(), entry.Value())
		}

		return result, value.Undefined
	default:
		result := NewOrderedHashRecordOfValue(h.Length())

		for _, entry := range h.entries {
			result.SetVal(thread, entry.Key(), entry.NativeValue())
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

func (h *OrderedNativeKeyHashRecord[K]) ContainsNativePair(thread *Thread, other *value.NativePair[K, value.Value]) (bool, value.Value) {
	v, ok := h.GetOk(other.NativeKey())
	if !ok {
		return false, value.Undefined
	}

	eq, err := Equal(thread, v, other.NativeValue())
	if err.IsNotUndefined() {
		return false, err
	}

	return value.Truthy(eq), value.Undefined
}

func (h *OrderedNativeKeyHashRecord[K]) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
	v, err := h.GetValUndefined(thread, other.Key())
	if err.IsNotUndefined() {
		return false, err
	}
	if v.IsUndefined() {
		return false, value.Undefined
	}

	eq, err := Equal(thread, v, other.Value())
	if err.IsNotUndefined() {
		return false, err
	}

	return value.Truthy(eq), value.Undefined
}

func (h *OrderedNativeKeyHashRecord[K]) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	for _, entry := range h.entries {
		eq, err := Equal(thread, entry.NativeValue(), val)
		if err.IsNotUndefined() {
			return false, err
		}

		if value.Truthy(eq) {
			return true, value.Undefined
		}
	}

	return false, value.Undefined
}

func (h *OrderedNativeKeyHashRecord[K]) ContainsNativeKey(thread *Thread, key K) bool {
	_, ok := h.GetOk(key)
	return ok
}

func (h *OrderedNativeKeyHashRecord[K]) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return false, value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}
	return h.ContainsNativeKey(thread, k), value.Undefined
}

func (h *OrderedNativeKeyHashRecord[K]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedNativeKeyHashRecord[K]:
		return NativePairSliceEqual(thread, h.entries, o.entries)
	case OrderedHashRecord:
		otherEntries := slices.Collect(o.All())
		return PairSliceEqualNative(thread, h.entries, otherEntries)
	}

	return false, value.Undefined
}

func (h *OrderedNativeKeyHashRecord[K]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
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

func (*OrderedNativeKeyHashRecord[K]) Class() *value.Class {
	return value.OrderedHashRecordClass
}

func (*OrderedNativeKeyHashRecord[K]) DirectClass() *value.Class {
	return value.OrderedHashRecordClass
}

func (*OrderedNativeKeyHashRecord[K]) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedNativeKeyHashRecord[K]) Clone() *OrderedNativeKeyHashRecord[K] {
	return &OrderedNativeKeyHashRecord[K]{
		index:   maps.Clone(h.index),
		entries: slices.Clone(h.entries),
	}
}

func (h *OrderedNativeKeyHashRecord[K]) Copy() value.Reference {
	return h.Clone()
}

func (h *OrderedNativeKeyHashRecord[K]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *OrderedNativeKeyHashRecord[K]) Error() string {
	return h.Inspect()
}

func (h *OrderedNativeKeyHashRecord[K]) Inspect() string {
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
	for key, val := range h.AllNative() {
		keyString := key.Inspect()
		keyStrings = append(keyStrings, keyString)

		valString := val.Inspect()
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

	buff.WriteString("Std::OrderedHashRecord{")
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
	buff.WriteRune('}')

	return buff.String()
}

func (*OrderedNativeKeyHashRecord[K]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedNativeKeyHashRecord[K]) Length() int {
	return len(h.entries)
}

type OrderedNativeKeyHashRecordIterator[K value.ComparableValueInterface] struct {
	HashRecord *OrderedNativeKeyHashRecord[K]
	index      int
	snapshot   []value.NativePair[K, value.Value]
}

var _ value.NativeResettableIterator = &OrderedNativeKeyHashRecordIterator[value.String]{}

func NewOrderedNativeKeyHashRecordIterator[K value.ComparableValueInterface](hrec *OrderedNativeKeyHashRecord[K]) *OrderedNativeKeyHashRecordIterator[K] {
	iterator := &OrderedNativeKeyHashRecordIterator[K]{
		HashRecord: hrec,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *OrderedNativeKeyHashRecordIterator[K]) captureSnapshot() {
	snapshot := make([]value.NativePair[K, value.Value], 0, h.HashRecord.Length())
	for k, v := range h.HashRecord.AllNative() {
		snapshot = append(snapshot, value.MakeNativePair(k, v))
	}
	h.snapshot = snapshot
}

func (*OrderedNativeKeyHashRecordIterator[K]) Class() *value.Class {
	return value.OrderedHashRecordIteratorClass
}

func (*OrderedNativeKeyHashRecordIterator[K]) DirectClass() *value.Class {
	return value.OrderedHashRecordIteratorClass
}

func (*OrderedNativeKeyHashRecordIterator[K]) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedNativeKeyHashRecordIterator[K]) Copy() value.Reference {
	return &OrderedNativeKeyHashRecordIterator[K]{
		HashRecord: h.HashRecord,
		index:      h.index,
		snapshot:   h.snapshot,
	}
}

func (i *OrderedNativeKeyHashRecordIterator[K]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *OrderedNativeKeyHashRecordIterator[K]) Error() string {
	return h.Inspect()
}

func (h *OrderedNativeKeyHashRecordIterator[K]) Inspect() string {
	return fmt.Sprintf("Std::OrderedHashRecord::Iterator{hash_record: %s}", h.HashRecord.Inspect())
}

func (*OrderedNativeKeyHashRecordIterator[K]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedNativeKeyHashRecordIterator[K]) Next() (p value.NativePair[K, value.Value], err value.Value) {
	if h.index >= len(h.snapshot) {
		return p, value.S(symbol.L_stop_iteration).ToValue()
	}

	v := h.snapshot[h.index]
	h.index++
	return v, value.Undefined
}

func (h *OrderedNativeKeyHashRecordIterator[K]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *OrderedNativeKeyHashRecordIterator[K]) Reset() {
	h.index = 0
}
