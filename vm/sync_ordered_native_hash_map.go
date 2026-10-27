package vm

import (
	"fmt"
	"iter"
	"strings"
	"sync"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type SyncOrderedNativeHashMap[K value.ComparableValueInterface, V value.ValueInterface] struct {
	Map *OrderedNativeHashMap[K, V]
	mu  sync.Mutex
}

var _ SyncOrderedHashMap = &SyncOrderedNativeHashMap[value.String, value.String]{}

func NewSyncOrderedNativeHashMap[K value.ComparableValueInterface, V value.ValueInterface](capacity int) *SyncOrderedNativeHashMap[K, V] {
	return &SyncOrderedNativeHashMap[K, V]{
		Map: NewOrderedNativeHashMap[K, V](capacity),
	}
}

func NewSyncOrderedNativeHashMapWithElements[K value.ComparableValueInterface, V value.ValueInterface](elements ...value.NativePair[K, V]) *SyncOrderedNativeHashMap[K, V] {
	return NewSyncOrderedNativeHashMapWithElementsAndTotalCapacity(len(elements), elements...)
}

func NewSyncOrderedNativeHashMapWithElementsAndTotalCapacity[K value.ComparableValueInterface, V value.ValueInterface](capacity int, elements ...value.NativePair[K, V]) *SyncOrderedNativeHashMap[K, V] {
	m := NewSyncOrderedNativeHashMap[K, V](capacity)
	for _, pair := range elements {
		m.Map.Set(pair.NativeKey(), pair.NativeValue())
	}
	return m
}

func (h *SyncOrderedNativeHashMap[K, V]) OrderedHashMap() OrderedHashMap {
	return h.Map
}

func (h *SyncOrderedNativeHashMap[K, V]) Lock() {
	h.mu.Lock()
}

func (h *SyncOrderedNativeHashMap[K, V]) Unlock() {
	h.mu.Unlock()
}

func (h *SyncOrderedNativeHashMap[K, V]) Keys() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		h.Lock()
		defer h.Unlock()

		for key := range h.Map.Keys() {
			if !yield(key) {
				return
			}
		}
	}
}

func (h *SyncOrderedNativeHashMap[K, V]) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		h.Lock()
		defer h.Unlock()

		for pair := range h.Map.All() {
			if !yield(pair) {
				return
			}
		}
	}
}

func (h *SyncOrderedNativeHashMap[K, V]) AllNative() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		h.Lock()
		defer h.Unlock()

		for key, val := range h.Map.AllNative() {
			if !yield(key, val) {
				return
			}
		}
	}
}

func (h *SyncOrderedNativeHashMap[K, V]) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		h.Lock()
		defer h.Unlock()

		for key, val := range h.Map.Iterate() {
			if !yield(key, val) {
				return
			}
		}
	}
}

func (h *SyncOrderedNativeHashMap[K, V]) CloneSyncOrderedHashMap(thread *Thread, capacity int) (SyncOrderedHashMap, value.Value) {
	h.Lock()
	defer h.Unlock()

	cloned, err := h.Map.CloneNative(thread, capacity)
	if err.IsNotUndefined() {
		return nil, err
	}
	return &SyncOrderedNativeHashMap[K, V]{Map: cloned}, value.Undefined
}

func (h *SyncOrderedNativeHashMap[K, V]) NewSyncOrderedHashMap(capacity int) SyncOrderedHashMap {
	return NewSyncOrderedNativeHashMap[K, V](capacity)
}

func (h *SyncOrderedNativeHashMap[K, V]) IterNative() *SyncOrderedNativeHashMapIterator[K, V] {
	return NewSyncOrderedNativeHashMapIterator(h)
}

func (h *SyncOrderedNativeHashMap[K, V]) Iter() value.NativeIterator {
	return h.IterNative()
}

func (h *SyncOrderedNativeHashMap[K, V]) IterSyncOrderedMap() value.NativeResettableIterator {
	return h.IterNative()
}

func (h *SyncOrderedNativeHashMap[K, V]) GetIndex(key K) (int, bool) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetIndex(key)
}

func (h *SyncOrderedNativeHashMap[K, V]) GetNativePair(key K) *value.NativePair[K, V] {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetNativePair(key)
}

func (h *SyncOrderedNativeHashMap[K, V]) GetPair(key K) value.Pair {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetPair(key)
}

func (h *SyncOrderedNativeHashMap[K, V]) Get(key K) value.Value {
	h.Lock()
	defer h.Unlock()

	return h.Map.Get(key)
}

func (h *SyncOrderedNativeHashMap[K, V]) GetOk(key K) (V, bool) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetOk(key)
}

func (h *SyncOrderedNativeHashMap[K, V]) At(index int) *value.NativePair[K, V] {
	h.Lock()
	defer h.Unlock()

	return h.Map.At(index)
}

func (h *SyncOrderedNativeHashMap[K, V]) AtInterface(index int) value.Pair {
	h.Lock()
	defer h.Unlock()

	return h.Map.AtInterface(index)
}

func (h *SyncOrderedNativeHashMap[K, V]) GetAt(index int) (*value.NativePair[K, V], value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetAt(index)
}

func (h *SyncOrderedNativeHashMap[K, V]) GetAtInterface(index int) (value.Pair, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetAtInterface(index)
}

func (h *SyncOrderedNativeHashMap[K, V]) TryAtVal(index value.Value) (value.Pair, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.TryAtVal(index)
}

func (h *SyncOrderedNativeHashMap[K, V]) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetValUndefined(thread, key)
}

func (h *SyncOrderedNativeHashMap[K, V]) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetValNil(thread, key)
}

func (h *SyncOrderedNativeHashMap[K, V]) Set(key K, val V) {
	h.Lock()
	defer h.Unlock()

	h.Map.Set(key, val)
}

func (h *SyncOrderedNativeHashMap[K, V]) SetVal(thread *Thread, key, val value.Value) value.Value {
	h.Lock()
	defer h.Unlock()

	return h.Map.SetVal(thread, key, val)
}

func (h *SyncOrderedNativeHashMap[K, V]) syncMergeResult(result OrderedHashRecord, err value.Value) (SyncOrderedHashMap, value.Value) {
	if err.IsNotUndefined() {
		return nil, err
	}
	switch m := result.(type) {
	case *OrderedNativeHashMap[K, V]:
		return &SyncOrderedNativeHashMap[K, V]{Map: m}, value.Undefined
	case *OrderedHashMapOfValue:
		return &SyncOrderedHashMapOfValue{Map: m}, value.Undefined
	default:
		return nil, value.Ref(value.Errorf(value.TypeErrorClass, "unexpected merge result"))
	}
}

func (h *SyncOrderedNativeHashMap[K, V]) MergeVal(thread *Thread, other value.Value) (SyncOrderedHashMap, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncOrderedNativeHashMap[K, V]:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return h.syncMergeResult(h.Map.MergeVal(thread, o.Map.ToValue()))
	case SyncOrderedHashMap:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return h.syncMergeResult(h.Map.MergeVal(thread, o.OrderedHashMap().ToValue()))
	default:
		h.Lock()
		defer h.Unlock()

		return h.syncMergeResult(h.Map.MergeVal(thread, other))
	}
}

func (h *SyncOrderedNativeHashMap[K, V]) ContainsNativePair(thread *Thread, other *value.NativePair[K, V]) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsNativePair(thread, other)
}

func (h *SyncOrderedNativeHashMap[K, V]) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.Contains(thread, other)
}

func (h *SyncOrderedNativeHashMap[K, V]) ContainsNativeValue(thread *Thread, val V) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsNativeValue(thread, val)
}

func (h *SyncOrderedNativeHashMap[K, V]) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsValue(thread, val)
}

func (h *SyncOrderedNativeHashMap[K, V]) ContainsNativeKey(thread *Thread, key K) bool {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsNativeKey(thread, key)
}

func (h *SyncOrderedNativeHashMap[K, V]) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsKey(thread, key)
}

func (h *SyncOrderedNativeHashMap[K, V]) RemoveNativeKey(key K) bool {
	h.Lock()
	defer h.Unlock()

	return h.Map.RemoveNativeKey(key)
}

func (h *SyncOrderedNativeHashMap[K, V]) RemoveKey(thread *Thread, key value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.RemoveKey(thread, key)
}

func (h *SyncOrderedNativeHashMap[K, V]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncOrderedNativeHashMap[K, V]:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return h.Map.Equal(thread, o.Map.ToValue())
	case SyncOrderedHashMap:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return h.Map.Equal(thread, o.OrderedHashMap().ToValue())
	default:
		return false, value.Undefined
	}
}

func (h *SyncOrderedNativeHashMap[K, V]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncOrderedNativeHashMap[K, V]:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return h.Map.LaxEqual(thread, o.Map.ToValue())
	case SyncOrderedHashMap:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return h.Map.LaxEqual(thread, o.OrderedHashMap().ToValue())
	default:
		if !value.IsA(other, value.MapMixin) {
			return false, value.Undefined
		}

		h.Lock()
		defer h.Unlock()

		for element, err := range Iterate(thread, other) {
			if err.IsNotUndefined() {
				return false, err
			}

			pair := element.AsReference().(value.Pair)
			val, err := h.Map.GetValUndefined(thread, pair.Key())
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

func (*SyncOrderedNativeHashMap[K, V]) Class() *value.Class {
	return value.SyncOrderedHashMapClass
}

func (*SyncOrderedNativeHashMap[K, V]) DirectClass() *value.Class {
	return value.SyncOrderedHashMapClass
}

func (*SyncOrderedNativeHashMap[K, V]) SingletonClass() *value.Class {
	return nil
}

func (h *SyncOrderedNativeHashMap[K, V]) Clone() *SyncOrderedNativeHashMap[K, V] {
	h.Lock()
	defer h.Unlock()

	return &SyncOrderedNativeHashMap[K, V]{
		Map: h.Map.Clone(),
	}
}

func (h *SyncOrderedNativeHashMap[K, V]) Copy() value.Reference {
	return h.Clone()
}

func (h *SyncOrderedNativeHashMap[K, V]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *SyncOrderedNativeHashMap[K, V]) Error() string {
	return h.Inspect()
}

func (h *SyncOrderedNativeHashMap[K, V]) Inspect() string {
	h.Lock()
	defer h.Unlock()

	var hasMultilineElements bool
	keyStrings := make(
		[]string,
		0,
		min(MAX_HASH_MAP_ELEMENTS_IN_INSPECT, h.Map.Length()),
	)
	valStrings := make(
		[]string,
		0,
		min(MAX_HASH_MAP_ELEMENTS_IN_INSPECT, h.Map.Length()),
	)

	i := 0
	for key, val := range h.Map.AllNative() {
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

	buff.WriteString("Std::Sync::OrderedHashMap{")
	if hasMultilineElements || h.Map.Length() > 15 {
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

func (*SyncOrderedNativeHashMap[K, V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncOrderedNativeHashMap[K, V]) Length() int {
	h.Lock()
	defer h.Unlock()

	return h.Map.Length()
}

type SyncOrderedNativeHashMapIterator[K value.ComparableValueInterface, V value.ValueInterface] struct {
	HashMap  *SyncOrderedNativeHashMap[K, V]
	index    int
	snapshot []value.NativePair[K, V]
	version  int
}

var _ value.NativeResettableIterator = &SyncOrderedNativeHashMapIterator[value.String, value.String]{}

func NewSyncOrderedNativeHashMapIterator[K value.ComparableValueInterface, V value.ValueInterface](hmap *SyncOrderedNativeHashMap[K, V]) *SyncOrderedNativeHashMapIterator[K, V] {
	iterator := &SyncOrderedNativeHashMapIterator[K, V]{
		HashMap: hmap,
		version: hmap.Map.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *SyncOrderedNativeHashMapIterator[K, V]) captureSnapshot() {
	h.HashMap.Lock()
	defer h.HashMap.Unlock()

	snapshot := make([]value.NativePair[K, V], 0, h.HashMap.Map.Length())
	for k, v := range h.HashMap.Map.AllNative() {
		snapshot = append(snapshot, value.MakeNativePair(k, v))
	}
	h.snapshot = snapshot
}

func (*SyncOrderedNativeHashMapIterator[K, V]) Class() *value.Class {
	return value.SyncOrderedHashMapIteratorClass
}

func (*SyncOrderedNativeHashMapIterator[K, V]) DirectClass() *value.Class {
	return value.SyncOrderedHashMapIteratorClass
}

func (*SyncOrderedNativeHashMapIterator[K, V]) SingletonClass() *value.Class {
	return nil
}

func (h *SyncOrderedNativeHashMapIterator[K, V]) Copy() value.Reference {
	return &SyncOrderedNativeHashMapIterator[K, V]{
		HashMap:  h.HashMap,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *SyncOrderedNativeHashMapIterator[K, V]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *SyncOrderedNativeHashMapIterator[K, V]) Error() string {
	return h.Inspect()
}

func (h *SyncOrderedNativeHashMapIterator[K, V]) Inspect() string {
	return fmt.Sprintf("Std::Sync::OrderedHashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*SyncOrderedNativeHashMapIterator[K, V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncOrderedNativeHashMapIterator[K, V]) Next() (p value.NativePair[K, V], err value.Value) {
	if h.version != h.HashMap.Map.version {
		return p, value.NewMutationDuringIterationError(h.HashMap.Class().Name()).ToValue()
	}
	if h.index >= len(h.snapshot) {
		return p, value.S(symbol.L_stop_iteration).ToValue()
	}

	v := h.snapshot[h.index]
	h.index++
	return v, value.Undefined
}

func (h *SyncOrderedNativeHashMapIterator[K, V]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *SyncOrderedNativeHashMapIterator[K, V]) Reset() {
	h.index = 0
	h.version = h.HashMap.Map.version
}
