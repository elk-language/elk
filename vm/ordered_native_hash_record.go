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

type OrderedNativeHashRecord[K value.ComparableValueInterface, V value.ValueInterface] struct {
	index   map[K]int
	entries []value.NativePair[K, V]
}

var _ OrderedHashRecord = &OrderedNativeHashRecord[value.String, value.String]{}

func NewOrderedNativeHashRecord[K value.ComparableValueInterface, V value.ValueInterface](capacity int) *OrderedNativeHashRecord[K, V] {
	return &OrderedNativeHashRecord[K, V]{
		index:   make(map[K]int, capacity),
		entries: make([]value.NativePair[K, V], 0, capacity),
	}
}

func NewOrderedNativeHashRecordWithElements[K value.ComparableValueInterface, V value.ValueInterface](elements ...value.NativePair[K, V]) *OrderedNativeHashRecord[K, V] {
	return NewOrderedNativeHashRecordWithElementsAndTotalCapacity(len(elements), elements...)
}

func NewOrderedNativeHashRecordWithElementsAndTotalCapacity[K value.ComparableValueInterface, V value.ValueInterface](capacity int, elements ...value.NativePair[K, V]) *OrderedNativeHashRecord[K, V] {
	m := NewOrderedNativeHashRecord[K, V](capacity)
	for _, pair := range elements {
		m.Set(pair.NativeKey(), pair.NativeValue())
	}
	return m
}

func (h *OrderedNativeHashRecord[K, V]) CloneOrderedHashRecord(thread *Thread, capacity int) (OrderedHashRecord, value.Value) {
	newIndex := make(map[K]int, capacity)
	maps.Copy(newIndex, h.index)

	entries := make([]value.NativePair[K, V], len(h.entries), capacity)
	copy(entries, h.entries)

	result := &OrderedNativeHashRecord[K, V]{
		index:   newIndex,
		entries: entries,
	}
	return result, value.Undefined
}

func (h *OrderedNativeHashRecord[K, V]) NewOrderedHashRecord(capacity int) OrderedHashRecord {
	return NewOrderedNativeHashRecord[K, V](capacity)
}

func (h *OrderedNativeHashRecord[K, V]) Keys() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		for _, pair := range h.entries {
			if !yield(pair.Key()) {
				return
			}
		}
	}
}

func (h *OrderedNativeHashRecord[K, V]) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		for _, pair := range h.entries {
			elkPair := value.MakePairOfValue(pair.Key(), pair.Value())
			if !yield(elkPair) {
				return
			}
		}
	}
}

func (h *OrderedNativeHashRecord[K, V]) AllNative() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, pair := range h.entries {
			if !yield(pair.NativeKey(), pair.NativeValue()) {
				return
			}
		}
	}
}

func (h *OrderedNativeHashRecord[K, V]) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		for _, pair := range h.entries {
			if !yield(pair.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (h *OrderedNativeHashRecord[K, V]) IterOrderedRecord() value.NativeResettableIterator {
	return NewOrderedNativeHashRecordIterator(h)
}

func (h *OrderedNativeHashRecord[K, V]) Iter() value.NativeIterator {
	return h.IterOrderedRecord()
}

func (h *OrderedNativeHashRecord[K, V]) GetIndex(key K) (int, bool) {
	index, ok := h.index[key]
	return index, ok
}

func (h *OrderedNativeHashRecord[K, V]) GetNativePair(key K) *value.NativePair[K, V] {
	i, ok := h.GetIndex(key)
	if !ok {
		return nil
	}
	return &h.entries[i]
}

func (h *OrderedNativeHashRecord[K, V]) GetPair(key K) value.Pair {
	return h.GetNativePair(key)
}

func (h *OrderedNativeHashRecord[K, V]) Get(key K) value.Value {
	i, ok := h.GetIndex(key)
	if !ok {
		return value.Undefined
	}
	return h.entries[i].Value()
}

func (h *OrderedNativeHashRecord[K, V]) GetOk(key K) (V, bool) {
	i, ok := h.GetIndex(key)
	if !ok {
		var zero V
		return zero, false
	}
	return h.entries[i].NativeValue(), true
}

func (h *OrderedNativeHashRecord[K, V]) At(index int) *value.NativePair[K, V] {
	return &h.entries[index]
}

func (h *OrderedNativeHashRecord[K, V]) AtInterface(index int) value.Pair {
	return h.At(index)
}

func (h *OrderedNativeHashRecord[K, V]) GetAt(index int) (result *value.NativePair[K, V], err value.Value) {
	index, err = value.NormalizeArrayIndex(index, h.Length())
	if err.IsNotUndefined() {
		return nil, err
	}

	return &h.entries[index], value.Undefined
}

func (h *OrderedNativeHashRecord[K, V]) GetAtInterface(index int) (result value.Pair, err value.Value) {
	result, err = h.GetAt(index)
	return result, err
}

func (h *OrderedNativeHashRecord[K, V]) TryAtVal(index value.Value) (value.Pair, value.Value) {
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

func (h *OrderedNativeHashRecord[K, V]) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return value.Undefined, value.Undefined
	}
	v, ok := h.GetOk(k)
	if !ok {
		return value.Undefined, value.Undefined
	}
	return v.ToValue(), value.Undefined
}

func (h *OrderedNativeHashRecord[K, V]) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return value.Nil, value.Undefined
	}
	v, ok := h.GetOk(k)
	if !ok {
		return value.Nil, value.Undefined
	}
	return v.ToValue(), value.Undefined
}

func (h *OrderedNativeHashRecord[K, V]) Set(key K, val V) {
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

func (h *OrderedNativeHashRecord[K, V]) SetVal(thread *Thread, key, val value.Value) value.Value {
	k, ok := value.Downcast[K](key)
	if !ok {
		return value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}
	v, ok := value.Downcast[V](val)
	if !ok {
		return value.NewInvalidValueInTypedMap(h, val.Class()).ToValue()
	}

	h.Set(k, v)
	return value.Undefined
}

func (h *OrderedNativeHashRecord[K, V]) MergeVal(thread *Thread, other value.Value) (OrderedHashRecord, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedNativeHashRecord[K, V]:
		result := NewOrderedNativeHashRecord[K, V](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for _, entry := range o.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
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

func (h *OrderedNativeHashRecord[K, V]) ContainsNativePair(thread *Thread, other *value.NativePair[K, V]) (bool, value.Value) {
	v, ok := h.GetOk(other.NativeKey())
	if !ok {
		return false, value.Undefined
	}

	eq, err := Equal(thread, v.ToValue(), other.Value())
	if err.IsNotUndefined() {
		return false, err
	}

	return value.Truthy(eq), value.Undefined
}

func (h *OrderedNativeHashRecord[K, V]) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
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

func (h *OrderedNativeHashRecord[K, V]) ContainsNativeValue(thread *Thread, val V) (bool, value.Value) {
	for _, entry := range h.entries {
		eq, err := Equal(thread, entry.Value(), val.ToValue())
		if err.IsNotUndefined() {
			return false, err
		}

		if value.Truthy(eq) {
			return true, value.Undefined
		}
	}

	return false, value.Undefined
}

func (h *OrderedNativeHashRecord[K, V]) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	v, ok := value.Downcast[V](val)
	if !ok {
		return false, value.NewInvalidValueInTypedMap(h, val.Class()).ToValue()
	}
	return h.ContainsNativeValue(thread, v)
}

func (h *OrderedNativeHashRecord[K, V]) ContainsNativeKey(thread *Thread, key K) bool {
	_, ok := h.GetOk(key)
	return ok
}

func (h *OrderedNativeHashRecord[K, V]) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return false, value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}
	return h.ContainsNativeKey(thread, k), value.Undefined
}

func (h *OrderedNativeHashRecord[K, V]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedNativeHashRecord[K, V]:
		return NativePairSliceEqual(thread, h.entries, o.entries)
	case OrderedHashRecord:
		otherEntries := slices.Collect(o.All())
		return PairSliceEqualNative(thread, h.entries, otherEntries)
	}

	return false, value.Undefined
}

func (h *OrderedNativeHashRecord[K, V]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case OrderedHashRecord:
		return OrderedHashRecordLaxEqual(thread, h, o)
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

func (*OrderedNativeHashRecord[K, V]) Class() *value.Class {
	return value.OrderedHashRecordClass
}

func (*OrderedNativeHashRecord[K, V]) DirectClass() *value.Class {
	return value.OrderedHashRecordClass
}

func (*OrderedNativeHashRecord[K, V]) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedNativeHashRecord[K, V]) Clone() *OrderedNativeHashRecord[K, V] {
	return &OrderedNativeHashRecord[K, V]{
		index:   maps.Clone(h.index),
		entries: slices.Clone(h.entries),
	}
}

func (h *OrderedNativeHashRecord[K, V]) Copy() value.Reference {
	return h.Clone()
}

func (h *OrderedNativeHashRecord[K, V]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *OrderedNativeHashRecord[K, V]) Error() string {
	return h.Inspect()
}

func (h *OrderedNativeHashRecord[K, V]) Inspect() string {
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

func (*OrderedNativeHashRecord[K, V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedNativeHashRecord[K, V]) Length() int {
	return len(h.entries)
}

type OrderedNativeHashRecordIterator[K value.ComparableValueInterface, V value.ValueInterface] struct {
	HashRecord *OrderedNativeHashRecord[K, V]
	index      int
	snapshot   []value.NativePair[K, V]
}

var _ value.NativeResettableIterator = &OrderedNativeHashRecordIterator[value.String, value.String]{}

func NewOrderedNativeHashRecordIterator[K value.ComparableValueInterface, V value.ValueInterface](hrec *OrderedNativeHashRecord[K, V]) *OrderedNativeHashRecordIterator[K, V] {
	iterator := &OrderedNativeHashRecordIterator[K, V]{
		HashRecord: hrec,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *OrderedNativeHashRecordIterator[K, V]) captureSnapshot() {
	snapshot := make([]value.NativePair[K, V], 0, h.HashRecord.Length())
	for k, v := range h.HashRecord.AllNative() {
		snapshot = append(snapshot, value.MakeNativePair(k, v))
	}
	h.snapshot = snapshot
}

func (*OrderedNativeHashRecordIterator[K, V]) Class() *value.Class {
	return value.OrderedHashRecordIteratorClass
}

func (*OrderedNativeHashRecordIterator[K, V]) DirectClass() *value.Class {
	return value.OrderedHashRecordIteratorClass
}

func (*OrderedNativeHashRecordIterator[K, V]) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedNativeHashRecordIterator[K, V]) Copy() value.Reference {
	return &OrderedNativeHashRecordIterator[K, V]{
		HashRecord: h.HashRecord,
		index:      h.index,
		snapshot:   h.snapshot,
	}
}

func (i *OrderedNativeHashRecordIterator[K, V]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *OrderedNativeHashRecordIterator[K, V]) Error() string {
	return h.Inspect()
}

func (h *OrderedNativeHashRecordIterator[K, V]) Inspect() string {
	return fmt.Sprintf("Std::OrderedHashRecord::Iterator{hash_record: %s}", h.HashRecord.Inspect())
}

func (*OrderedNativeHashRecordIterator[K, V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedNativeHashRecordIterator[K, V]) Next() (p value.NativePair[K, V], err value.Value) {
	if h.index >= len(h.snapshot) {
		return p, value.S(symbol.L_stop_iteration).ToValue()
	}

	v := h.snapshot[h.index]
	h.index++
	return v, value.Undefined
}

func (h *OrderedNativeHashRecordIterator[K, V]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *OrderedNativeHashRecordIterator[K, V]) Reset() {
	h.index = 0
}
