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

type OrderedNativeHashMap[K value.ComparableValueInterface, V value.ValueInterface] struct {
	index   map[K]int
	entries []value.NativePair[K, V]
	version int
}

var _ OrderedHashMap = &OrderedNativeHashMap[value.String, value.String]{}

func NewOrderedNativeHashMap[K value.ComparableValueInterface, V value.ValueInterface](capacity int) *OrderedNativeHashMap[K, V] {
	return &OrderedNativeHashMap[K, V]{
		index:   make(map[K]int, capacity),
		entries: make([]value.NativePair[K, V], 0, capacity),
	}
}

func NewOrderedNativeHashMapWithElements[K value.ComparableValueInterface, V value.ValueInterface](elements ...value.NativePair[K, V]) *OrderedNativeHashMap[K, V] {
	return NewOrderedNativeHashMapWithElementsAndTotalCapacity(len(elements), elements...)
}

func NewOrderedNativeHashMapWithElementsAndTotalCapacity[K value.ComparableValueInterface, V value.ValueInterface](capacity int, elements ...value.NativePair[K, V]) *OrderedNativeHashMap[K, V] {
	m := NewOrderedNativeHashMap[K, V](capacity)
	for _, pair := range elements {
		m.Set(pair.NativeKey(), pair.NativeValue())
	}
	return m
}

func (h *OrderedNativeHashMap[K, V]) CloneNative(thread *Thread, capacity int) (*OrderedNativeHashMap[K, V], value.Value) {
	newMap := make(map[K]int, capacity)
	maps.Copy(newMap, h.index)

	entries := make([]value.NativePair[K, V], len(h.entries), capacity)
	copy(entries, h.entries)

	result := &OrderedNativeHashMap[K, V]{
		index:   newMap,
		entries: entries,
	}
	return result, value.Undefined
}

func (h *OrderedNativeHashMap[K, V]) CloneOrderedHashMap(thread *Thread, capacity int) (OrderedHashMap, value.Value) {
	result, err := h.CloneNative(thread, capacity)
	if err.IsNotUndefined() {
		return nil, err
	}

	return result, value.Undefined
}

func (h *OrderedNativeHashMap[K, V]) NewOrderedHashMap(capacity int) OrderedHashMap {
	return NewOrderedNativeHashMap[K, V](capacity)
}

func (h *OrderedNativeHashMap[K, V]) NewOrderedHashRecord(capacity int) OrderedHashRecord {
	return h.NewOrderedHashMap(capacity)
}

func (h *OrderedNativeHashMap[K, V]) CloneOrderedHashRecord(thread *Thread, capacity int) (OrderedHashRecord, value.Value) {
	return h.CloneOrderedHashMap(thread, capacity)
}

func (h *OrderedNativeHashMap[K, V]) Keys() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		for _, pair := range h.entries {
			if !yield(pair.Key()) {
				return
			}
		}
	}
}

func (h *OrderedNativeHashMap[K, V]) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		for _, pair := range h.entries {
			elkPair := value.MakePairOfValue(pair.Key(), pair.Value())
			if !yield(elkPair) {
				return
			}
		}
	}
}

func (h *OrderedNativeHashMap[K, V]) AllNative() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, pair := range h.entries {
			if !yield(pair.NativeKey(), pair.NativeValue()) {
				return
			}
		}
	}
}

func (h *OrderedNativeHashMap[K, V]) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		originalVersion := h.version

		for _, pair := range h.entries {
			if originalVersion != h.version {
				yield(value.Undefined, value.NewMutationDuringIterationError(h.Class().Name()).ToValue())
				return
			}

			if !yield(pair.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (h *OrderedNativeHashMap[K, V]) IterOrderedMap() value.NativeResettableIterator {
	return NewOrderedNativeHashMapIterator(h)
}

func (h *OrderedNativeHashMap[K, V]) IterOrderedRecord() value.NativeResettableIterator {
	return h.IterOrderedMap()
}

func (h *OrderedNativeHashMap[K, V]) Iter() value.NativeIterator {
	return h.IterOrderedMap()
}

func (h *OrderedNativeHashMap[K, V]) GetIndex(key K) (int, bool) {
	index, ok := h.index[key]
	return index, ok
}

func (h *OrderedNativeHashMap[K, V]) GetNativePair(key K) *value.NativePair[K, V] {
	i, ok := h.GetIndex(key)
	if !ok {
		return nil
	}
	return &h.entries[i]
}

func (h *OrderedNativeHashMap[K, V]) GetPair(key K) value.Pair {
	return h.GetNativePair(key)
}

func (h *OrderedNativeHashMap[K, V]) Get(key K) value.Value {
	i, ok := h.GetIndex(key)
	if !ok {
		return value.Undefined
	}
	return h.entries[i].Value()
}

func (h *OrderedNativeHashMap[K, V]) GetOk(key K) (V, bool) {
	i, ok := h.GetIndex(key)
	if !ok {
		var zero V
		return zero, false
	}
	return h.entries[i].NativeValue(), true
}

func (h *OrderedNativeHashMap[K, V]) At(index int) *value.NativePair[K, V] {
	return &h.entries[index]
}

func (h *OrderedNativeHashMap[K, V]) AtInterface(index int) value.Pair {
	return h.At(index)
}

func (h *OrderedNativeHashMap[K, V]) GetAt(index int) (result *value.NativePair[K, V], err value.Value) {
	index, err = value.NormalizeArrayIndex(index, h.Length())
	if err.IsNotUndefined() {
		return nil, err
	}

	return &h.entries[index], value.Undefined
}

func (h *OrderedNativeHashMap[K, V]) GetAtInterface(index int) (result value.Pair, err value.Value) {
	result, err = h.GetAt(index)
	return result, err
}

// Normalizes the index so that `-1` becomes `length-1`.
// Does bounds checking and returns (undefined, undefined) if the index is out of range.
func (h *OrderedNativeHashMap[K, V]) TryAtVal(index value.Value) (value.Pair, value.Value) {
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

func (h *OrderedNativeHashMap[K, V]) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
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

func (h *OrderedNativeHashMap[K, V]) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
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

func (h *OrderedNativeHashMap[K, V]) Set(key K, val V) {
	entryIndex, ok := h.GetIndex(key)
	h.version++
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

func (h *OrderedNativeHashMap[K, V]) SetVal(thread *Thread, key, val value.Value) value.Value {
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

func (h *OrderedNativeHashMap[K, V]) MergeVal(thread *Thread, other value.Value) (OrderedHashRecord, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedNativeHashMap[K, V]:
		result := NewOrderedNativeHashMap[K, V](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for _, entry := range o.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}

		return result, value.Undefined
	case *NativeHashMap[K, V]:
		result := NewOrderedNativeHashMap[K, V](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case *NativeKeyHashMap[K]:
		result := NewOrderedNativeKeyHashMap[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue().ToValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case *OrderedNativeKeyHashMap[K]:
		result := NewOrderedNativeKeyHashMap[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue().ToValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case *SyncNativeHashMap[K, V]:
		o.Lock()
		defer o.Unlock()

		result := NewOrderedNativeHashMap[K, V](h.Length() + o.Map.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for k, v := range o.Map.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case *SyncNativeKeyHashMap[K]:
		result := NewOrderedNativeKeyHashMap[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue().ToValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case *SyncOrderedNativeHashMap[K, V]:
		result := NewOrderedNativeHashMap[K, V](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case *SyncOrderedNativeKeyHashMap[K]:
		result := NewOrderedNativeKeyHashMap[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue().ToValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case NativeHashRecord[K, V]:
		result := NewOrderedNativeHashMap[K, V](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case NativeKeyHashRecord[K]:
		result := NewOrderedNativeKeyHashMap[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue().ToValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case *OrderedNativeHashRecord[K, V]:
		result := NewOrderedNativeHashMap[K, V](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case *OrderedNativeKeyHashRecord[K]:
		result := NewOrderedNativeKeyHashMap[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue().ToValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case Record:
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

func (h *OrderedNativeHashMap[K, V]) ContainsNativePair(thread *Thread, other *value.NativePair[K, V]) (bool, value.Value) {
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

func (h *OrderedNativeHashMap[K, V]) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
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

func (h *OrderedNativeHashMap[K, V]) ContainsNativeValue(thread *Thread, val V) (bool, value.Value) {
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

func (h *OrderedNativeHashMap[K, V]) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	v, ok := value.Downcast[V](val)
	if !ok {
		return false, value.NewInvalidValueInTypedMap(h, val.Class()).ToValue()
	}
	return h.ContainsNativeValue(thread, v)
}

func (h *OrderedNativeHashMap[K, V]) ContainsNativeKey(thread *Thread, key K) bool {
	_, ok := h.GetOk(key)
	return ok
}

func (h *OrderedNativeHashMap[K, V]) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return false, value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}
	return h.ContainsNativeKey(thread, k), value.Undefined
}

func (h *OrderedNativeHashMap[K, V]) RemoveNativeKey(key K) bool {
	index, ok := h.GetIndex(key)
	if !ok {
		return false
	}

	h.version++
	h.entries = slices.Delete(h.entries, index, index+1)
	delete(h.index, key)
	for i := index; i < len(h.entries); i++ {
		h.index[h.entries[i].NativeKey()] = i
	}
	return true
}

func (h *OrderedNativeHashMap[K, V]) RemoveKey(thread *Thread, key value.Value) (bool, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return false, value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}
	return h.RemoveNativeKey(k), value.Undefined
}

func (h *OrderedNativeHashMap[K, V]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedNativeHashMap[K, V]:
		return NativePairSliceEqual(thread, h.entries, o.entries)
	case OrderedHashMap:
		otherEntries := slices.Collect(o.All())
		return PairSliceEqualNative(thread, h.entries, otherEntries)
	}

	return false, value.Undefined
}

func (h *OrderedNativeHashMap[K, V]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
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

func (*OrderedNativeHashMap[K, V]) Class() *value.Class {
	return value.OrderedHashMapClass
}

func (*OrderedNativeHashMap[K, V]) DirectClass() *value.Class {
	return value.OrderedHashMapClass
}

func (*OrderedNativeHashMap[K, V]) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedNativeHashMap[K, V]) Clone() *OrderedNativeHashMap[K, V] {
	return &OrderedNativeHashMap[K, V]{
		index:   maps.Clone(h.index),
		entries: slices.Clone(h.entries),
	}
}

func (h *OrderedNativeHashMap[K, V]) ToNativeHashMap() *NativeHashMap[K, V] {
	hashMap := NewNativeHashMap[K, V](len(h.entries))
	for _, pair := range h.entries {
		hashMap.Set(pair.NativeKey(), pair.NativeValue())
	}
	return hashMap
}

func (h *OrderedNativeHashMap[K, V]) ToHashMap(thread *Thread) (HashMap, value.Value) {
	return h.ToNativeHashMap(), value.Undefined
}

func (h *OrderedNativeHashMap[K, V]) ToNativeHashRecord() NativeHashRecord[K, V] {
	record := MakeNativeHashRecord[K, V](len(h.entries))
	for _, pair := range h.entries {
		record.Set(pair.NativeKey(), pair.NativeValue())
	}
	return record
}

func (h *OrderedNativeHashMap[K, V]) ToHashRecord(thread *Thread) (HashRecord, value.Value) {
	return h.ToNativeHashRecord(), value.Undefined
}

func (h *OrderedNativeHashMap[K, V]) ToSyncHashMap(thread *Thread) (SyncHashMap, value.Value) {
	return &SyncNativeHashMap[K, V]{Map: h.ToNativeHashMap()}, value.Undefined
}

func (h *OrderedNativeHashMap[K, V]) ToSyncOrderedHashMap() SyncOrderedHashMap {
	return &SyncOrderedNativeHashMap[K, V]{Map: h.Clone()}
}

func (h *OrderedNativeHashMap[K, V]) ToOrderedNativeHashRecord() *OrderedNativeHashRecord[K, V] {
	return &OrderedNativeHashRecord[K, V]{
		index:   maps.Clone(h.index),
		entries: slices.Clone(h.entries),
	}
}

func (h *OrderedNativeHashMap[K, V]) ToOrderedHashRecord() OrderedHashRecord {
	return h.ToOrderedNativeHashRecord()
}

func (h *OrderedNativeHashMap[K, V]) ToOrderedHashMap() OrderedHashMap {
	return h
}

func (h *OrderedNativeHashMap[K, V]) Copy() value.Reference {
	return h.Clone()
}

func (h *OrderedNativeHashMap[K, V]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *OrderedNativeHashMap[K, V]) Error() string {
	return h.Inspect()
}

func (h *OrderedNativeHashMap[K, V]) Inspect() string {
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
	for key, val := range h.AllNative() {
		keyString := key.Inspect()
		keyStrings = append(keyStrings, keyString)

		valString := val.Inspect()
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

func (*OrderedNativeHashMap[K, V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedNativeHashMap[K, V]) Length() int {
	return len(h.entries)
}

type OrderedNativeHashMapIterator[K value.ComparableValueInterface, V value.ValueInterface] struct {
	HashMap  *OrderedNativeHashMap[K, V]
	index    int
	snapshot []value.NativePair[K, V]
	version  int
}

var _ value.NativeResettableIterator = &OrderedNativeHashMapIterator[value.String, value.String]{}

func NewOrderedNativeHashMapIterator[K value.ComparableValueInterface, V value.ValueInterface](hmap *OrderedNativeHashMap[K, V]) *OrderedNativeHashMapIterator[K, V] {
	iterator := &OrderedNativeHashMapIterator[K, V]{
		HashMap: hmap,
		version: hmap.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *OrderedNativeHashMapIterator[K, V]) captureSnapshot() {
	snapshot := make([]value.NativePair[K, V], 0, h.HashMap.Length())
	for k, v := range h.HashMap.AllNative() {
		snapshot = append(snapshot, value.MakeNativePair(k, v))
	}
	h.snapshot = snapshot
}

func (*OrderedNativeHashMapIterator[K, V]) Class() *value.Class {
	return value.OrderedHashMapIteratorClass
}

func (*OrderedNativeHashMapIterator[K, V]) DirectClass() *value.Class {
	return value.OrderedHashMapIteratorClass
}

func (*OrderedNativeHashMapIterator[K, V]) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedNativeHashMapIterator[K, V]) Copy() value.Reference {
	return &OrderedNativeHashMapIterator[K, V]{
		HashMap:  h.HashMap,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *OrderedNativeHashMapIterator[K, V]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *OrderedNativeHashMapIterator[K, V]) Error() string {
	return h.Inspect()
}

func (h *OrderedNativeHashMapIterator[K, V]) Inspect() string {
	return fmt.Sprintf("Std::OrderedHashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*OrderedNativeHashMapIterator[K, V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedNativeHashMapIterator[K, V]) Next() (p value.NativePair[K, V], err value.Value) {
	if h.version != h.HashMap.version {
		return p, value.NewMutationDuringIterationError(h.HashMap.Class().Name()).ToValue()
	}
	if h.index >= len(h.snapshot) {
		return p, value.S(symbol.L_stop_iteration).ToValue()
	}

	v := h.snapshot[h.index]
	h.index++
	return v, value.Undefined
}

func (h *OrderedNativeHashMapIterator[K, V]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *OrderedNativeHashMapIterator[K, V]) Reset() {
	h.index = 0
}
