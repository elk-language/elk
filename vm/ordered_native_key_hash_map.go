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

type OrderedNativeKeyHashMap[K value.ComparableValueInterface] struct {
	index   map[K]int
	entries []value.NativePair[K, value.Value] // slice with key value pairs
	version int
}

var _ OrderedHashMap = &OrderedNativeKeyHashMap[value.String]{}

func NewOrderedNativeKeyHashMap[K value.ComparableValueInterface](capacity int) *OrderedNativeKeyHashMap[K] {
	return &OrderedNativeKeyHashMap[K]{
		index:   make(map[K]int, capacity),
		entries: make([]value.NativePair[K, value.Value], capacity),
	}
}

func NewOrderedNativeKeyHashMapWithElements[K value.ComparableValueInterface](elements ...value.NativePair[K, value.Value]) *OrderedNativeKeyHashMap[K] {
	return NewOrderedNativeKeyHashMapWithElementsAndTotalCapacity(len(elements), elements...)
}

func NewOrderedNativeKeyHashMapWithElementsAndTotalCapacity[K value.ComparableValueInterface](capacity int, elements ...value.NativePair[K, value.Value]) *OrderedNativeKeyHashMap[K] {
	m := NewOrderedNativeKeyHashMap[K](capacity)
	for _, pair := range elements {
		m.Set(pair.NativeKey(), pair.NativeValue())
	}
	return m
}

func (h *OrderedNativeKeyHashMap[K]) CloneOrderedHashMap(thread *Thread, capacity int) (OrderedHashMap, value.Value) {
	newMap := make(map[K]int, capacity)
	maps.Copy(newMap, h.index)

	entries := make([]value.NativePair[K, value.Value], len(h.entries), capacity)
	copy(entries, h.entries)

	result := &OrderedNativeKeyHashMap[K]{
		index:   newMap,
		entries: entries,
	}
	return result, value.Undefined
}

func (h *OrderedNativeKeyHashMap[K]) NewOrderedHashMap(capacity int) OrderedHashMap {
	return NewOrderedNativeKeyHashMap[K](capacity)
}

func (h *OrderedNativeKeyHashMap[K]) NewOrderedHashRecord(capacity int) OrderedHashRecord {
	return h.NewOrderedHashMap(capacity)
}

func (h *OrderedNativeKeyHashMap[K]) CloneOrderedHashRecord(thread *Thread, capacity int) (OrderedHashRecord, value.Value) {
	return h.CloneOrderedHashMap(thread, capacity)
}

func (h *OrderedNativeKeyHashMap[K]) Keys() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		for _, pair := range h.entries {
			if !yield(pair.Key()) {
				return
			}
		}
	}
}

func (h *OrderedNativeKeyHashMap[K]) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		for _, pair := range h.entries {
			pair := value.MakePairOfValue(pair.Key(), pair.NativeValue())
			if !yield(pair) {
				return
			}
		}
	}
}

func (h *OrderedNativeKeyHashMap[K]) AllNative() iter.Seq2[K, value.Value] {
	return func(yield func(K, value.Value) bool) {
		for _, pair := range h.entries {
			if !yield(pair.NativeKey(), pair.NativeValue()) {
				return
			}
		}
	}
}

func (h *OrderedNativeKeyHashMap[K]) Iterate() iter.Seq2[value.Value, value.Value] {
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

func (h *OrderedNativeKeyHashMap[K]) IterOrderedMap() value.NativeResettableIterator {
	return NewOrderedNativeKeyHashMapIterator(h)
}

func (h *OrderedNativeKeyHashMap[K]) IterOrderedRecord() value.NativeResettableIterator {
	return h.IterOrderedMap()
}

func (h *OrderedNativeKeyHashMap[K]) Iter() value.NativeIterator {
	return h.IterOrderedMap()
}

func (h *OrderedNativeKeyHashMap[K]) GetIndex(key K) (int, bool) {
	index, ok := h.index[key]
	return index, ok
}

func (h *OrderedNativeKeyHashMap[K]) GetNativePair(key K) *value.NativePair[K, value.Value] {
	i, ok := h.GetIndex(key)
	if !ok {
		return nil
	}
	return &h.entries[i]
}

func (h *OrderedNativeKeyHashMap[K]) GetPair(key K) value.Pair {
	return h.GetNativePair(key)
}

func (h *OrderedNativeKeyHashMap[K]) Get(key K) value.Value {
	i, ok := h.GetIndex(key)
	if !ok {
		return value.Undefined
	}
	return h.entries[i].NativeValue()
}

func (h *OrderedNativeKeyHashMap[K]) GetOk(key K) (value.Value, bool) {
	i, ok := h.GetIndex(key)
	if !ok {
		return value.Undefined, false
	}
	return h.entries[i].NativeValue(), true
}

func (h *OrderedNativeKeyHashMap[K]) At(index int) *value.NativePair[K, value.Value] {
	return &h.entries[index]
}

func (h *OrderedNativeKeyHashMap[K]) AtInterface(index int) value.Pair {
	return h.At(index)
}

func (h *OrderedNativeKeyHashMap[K]) GetAt(index int) (result *value.NativePair[K, value.Value], err value.Value) {
	index, err = value.NormalizeArrayIndex(index, h.Length())
	if err.IsNotUndefined() {
		return nil, err
	}

	return &h.entries[index], value.Undefined
}

func (h *OrderedNativeKeyHashMap[K]) GetAtInterface(index int) (result value.Pair, err value.Value) {
	result, err = h.GetAt(index)
	return result, err
}

func (h *OrderedNativeKeyHashMap[K]) TryAtVal(index value.Value) (value.Pair, value.Value) {
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

func (h *OrderedNativeKeyHashMap[K]) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
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

func (h *OrderedNativeKeyHashMap[K]) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
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

func (h *OrderedNativeKeyHashMap[K]) Set(key K, val value.Value) {
	entryIndex, ok := h.GetIndex(key)
	h.version++
	if !ok {
		// key is already present
		h.entries[entryIndex] = value.MakeNativePair(key, val)
	} else {
		// new key
		entryIndex := len(h.entries)
		h.entries = append(h.entries, value.MakeNativePair(key, val))
		h.index[key] = entryIndex
	}
}

func (h *OrderedNativeKeyHashMap[K]) SetVal(thread *Thread, key, val value.Value) value.Value {
	k, ok := value.Downcast[K](key)
	if !ok {
		return value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}

	h.Set(k, val)
	return value.Undefined
}

func (h *OrderedNativeKeyHashMap[K]) MergeVal(thread *Thread, other value.Value) (OrderedHashRecord, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedNativeKeyHashMap[K]:
		result := NewOrderedNativeKeyHashMap[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.SetVal(thread, entry.Key(), entry.NativeValue())
		}
		for _, entry := range o.entries {
			result.SetVal(thread, entry.Key(), entry.NativeValue())
		}

		return result, value.Undefined
	case *NativeKeyHashMap[K]:
		result := NewOrderedNativeKeyHashMap[K](h.Length() + o.Length())
		for _, entry := range h.entries {
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case *SyncNativeKeyHashMap[K]:
		result := NewOrderedNativeKeyHashMap[K](h.Length() + o.Length())
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
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case NativeKeyHashRecord[K]:
		result := NewOrderedNativeKeyHashMap[K](h.Length() + o.Length())
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
			result.Set(entry.NativeKey(), entry.NativeValue())
		}
		for k, v := range o.AllNative() {
			result.Set(k, v)
		}
		return result, value.Undefined
	case Record:
		result := NewOrderedHashMapOfValue(h.Length() + o.Length())
		for _, entry := range h.entries {
			result.SetVal(thread, entry.Key(), entry.NativeValue())
		}
		for entry := range o.All() {
			result.SetVal(thread, entry.Key(), entry.Value())
		}

		return result, value.Undefined
	default:
		result := NewOrderedHashMapOfValue(h.Length())

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

func (h *OrderedNativeKeyHashMap[K]) ContainsNativePair(thread *Thread, other *value.NativePair[K, value.Value]) (bool, value.Value) {
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

func (h *OrderedNativeKeyHashMap[K]) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
	v, err := h.GetValUndefined(thread, other.Key())
	if err.IsNotUndefined() {
		return false, err
	}

	eq, err := Equal(thread, v, other.Value())
	if err.IsNotUndefined() {
		return false, err
	}

	return value.Truthy(eq), value.Undefined
}

func (h *OrderedNativeKeyHashMap[K]) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	for _, hval := range h.entries {
		eq, err := Equal(thread, hval.ToValue(), val)
		if err.IsNotUndefined() {
			return false, err
		}

		if value.Truthy(eq) {
			return true, value.Undefined
		}
	}

	return false, value.Undefined
}

func (h *OrderedNativeKeyHashMap[K]) ContainsNativeKey(thread *Thread, key K) bool {
	_, ok := h.GetOk(key)
	return ok
}

func (h *OrderedNativeKeyHashMap[K]) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return false, value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}
	return h.ContainsNativeKey(thread, k), value.Undefined
}

func (h *OrderedNativeKeyHashMap[K]) RemoveNativeKey(key K) bool {
	index, ok := h.GetIndex(key)
	if !ok {
		return false
	}

	h.version++
	h.entries = slices.Delete(h.entries, int(index), int(index)+1)
	delete(h.index, key)
	return true
}

func (h *OrderedNativeKeyHashMap[K]) RemoveKey(thread *Thread, key value.Value) (bool, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return false, value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}
	return h.RemoveNativeKey(k), value.Undefined
}

func NativePairSliceEqual[K1, K2 value.ComparableValueInterface, V1, V2 value.ValueInterface](thread *Thread, a []value.NativePair[K1, V1], b []value.NativePair[K2, V2]) (bool, value.Value) {
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

func PairSliceEqualNative[K value.ComparableValueInterface, V value.ValueInterface](thread *Thread, a []value.NativePair[K, V], b []value.PairOfValue) (bool, value.Value) {
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

func (h *OrderedNativeKeyHashMap[K]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *OrderedNativeKeyHashMap[K]:
		return NativePairSliceEqual(thread, h.entries, o.entries)
	case OrderedHashMap:
		otherEntries := slices.Collect(o.All())
		return PairSliceEqualNative(thread, h.entries, otherEntries)
	}

	return false, value.Undefined
}

func (h *OrderedNativeKeyHashMap[K]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
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

func (*OrderedNativeKeyHashMap[K]) Class() *value.Class {
	return value.OrderedHashMapClass
}

func (*OrderedNativeKeyHashMap[K]) DirectClass() *value.Class {
	return value.OrderedHashMapClass
}

func (*OrderedNativeKeyHashMap[K]) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedNativeKeyHashMap[K]) Clone() *OrderedNativeKeyHashMap[K] {
	return &OrderedNativeKeyHashMap[K]{
		index:   maps.Clone(h.index),
		entries: slices.Clone(h.entries),
	}
}

func (h *OrderedNativeKeyHashMap[K]) Copy() value.Reference {
	return h.Clone()
}

func (h *OrderedNativeKeyHashMap[K]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *OrderedNativeKeyHashMap[K]) Error() string {
	return h.Inspect()
}

func (h *OrderedNativeKeyHashMap[K]) Inspect() string {
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

func (*OrderedNativeKeyHashMap[K]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedNativeKeyHashMap[K]) Length() int {
	return len(h.entries)
}

type OrderedNativeKeyHashMapIterator[K value.ComparableValueInterface] struct {
	HashMap  *OrderedNativeKeyHashMap[K]
	index    int
	snapshot []value.NativePair[K, value.Value]
	version  int
}

var _ value.NativeResettableIterator = &OrderedNativeKeyHashMapIterator[value.String]{}

func NewOrderedNativeKeyHashMapIterator[K value.ComparableValueInterface](hmap *OrderedNativeKeyHashMap[K]) *OrderedNativeKeyHashMapIterator[K] {
	iterator := &OrderedNativeKeyHashMapIterator[K]{
		HashMap: hmap,
		version: hmap.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *OrderedNativeKeyHashMapIterator[K]) captureSnapshot() {
	snapshot := make([]value.NativePair[K, value.Value], 0, h.HashMap.Length())
	for k, v := range h.HashMap.AllNative() {
		snapshot = append(snapshot, value.MakeNativePair(k, v))
	}
	h.snapshot = snapshot
}

func (*OrderedNativeKeyHashMapIterator[K]) Class() *value.Class {
	return value.OrderedHashMapIteratorClass
}

func (*OrderedNativeKeyHashMapIterator[K]) DirectClass() *value.Class {
	return value.OrderedHashMapIteratorClass
}

func (*OrderedNativeKeyHashMapIterator[K]) SingletonClass() *value.Class {
	return nil
}

func (h *OrderedNativeKeyHashMapIterator[K]) Copy() value.Reference {
	return &OrderedNativeKeyHashMapIterator[K]{
		HashMap:  h.HashMap,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *OrderedNativeKeyHashMapIterator[K]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *OrderedNativeKeyHashMapIterator[K]) Error() string {
	return h.Inspect()
}

func (h *OrderedNativeKeyHashMapIterator[K]) Inspect() string {
	return fmt.Sprintf("Std::OrderedHashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*OrderedNativeKeyHashMapIterator[K]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *OrderedNativeKeyHashMapIterator[K]) Next() (p value.NativePair[K, value.Value], err value.Value) {
	if h.version != h.HashMap.version {
		return p, value.NewMutationDuringIterationError(h.Class().Name()).ToValue()
	}
	if h.index >= len(h.snapshot) {
		return p, value.S(symbol.L_stop_iteration).ToValue()
	}

	v := h.snapshot[h.index]
	h.index++
	return v, value.Undefined
}

func (h *OrderedNativeKeyHashMapIterator[K]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *OrderedNativeKeyHashMapIterator[K]) Reset() {
	h.index = 0
}
