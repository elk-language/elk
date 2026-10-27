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

type SyncOrderedNativeKeyHashMap[K value.ComparableValueInterface] struct {
	Map *OrderedNativeKeyHashMap[K]
	mu  sync.Mutex
}

var _ SyncOrderedHashMap = &SyncOrderedNativeKeyHashMap[value.String]{}

func NewSyncOrderedNativeKeyHashMap[K value.ComparableValueInterface](capacity int) *SyncOrderedNativeKeyHashMap[K] {
	return &SyncOrderedNativeKeyHashMap[K]{
		Map: NewOrderedNativeKeyHashMap[K](capacity),
	}
}

func NewSyncOrderedNativeKeyHashMapWithElements[K value.ComparableValueInterface](elements ...value.NativePair[K, value.Value]) *SyncOrderedNativeKeyHashMap[K] {
	return NewSyncOrderedNativeKeyHashMapWithElementsAndTotalCapacity(len(elements), elements...)
}

func NewSyncOrderedNativeKeyHashMapWithElementsAndTotalCapacity[K value.ComparableValueInterface](capacity int, elements ...value.NativePair[K, value.Value]) *SyncOrderedNativeKeyHashMap[K] {
	m := NewSyncOrderedNativeKeyHashMap[K](capacity)
	for _, pair := range elements {
		m.Map.Set(pair.NativeKey(), pair.NativeValue())
	}
	return m
}

func (h *SyncOrderedNativeKeyHashMap[K]) OrderedHashMap() OrderedHashMap {
	return h.Map
}

func (h *SyncOrderedNativeKeyHashMap[K]) Lock() {
	h.mu.Lock()
}

func (h *SyncOrderedNativeKeyHashMap[K]) Unlock() {
	h.mu.Unlock()
}

func (h *SyncOrderedNativeKeyHashMap[K]) Keys() iter.Seq[value.Value] {
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

func (h *SyncOrderedNativeKeyHashMap[K]) All() iter.Seq[value.PairOfValue] {
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

func (h *SyncOrderedNativeKeyHashMap[K]) AllNative() iter.Seq2[K, value.Value] {
	return func(yield func(K, value.Value) bool) {
		h.Lock()
		defer h.Unlock()

		for key, val := range h.Map.AllNative() {
			if !yield(key, val) {
				return
			}
		}
	}
}

func (h *SyncOrderedNativeKeyHashMap[K]) Iterate() iter.Seq2[value.Value, value.Value] {
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

func (h *SyncOrderedNativeKeyHashMap[K]) CloneSyncOrderedHashMap(thread *Thread, capacity int) (SyncOrderedHashMap, value.Value) {
	h.Lock()
	defer h.Unlock()

	cloned, err := h.Map.CloneOrderedHashMap(thread, capacity)
	if err.IsNotUndefined() {
		return nil, err
	}
	return &SyncOrderedNativeKeyHashMap[K]{Map: cloned.(*OrderedNativeKeyHashMap[K])}, value.Undefined
}

func (h *SyncOrderedNativeKeyHashMap[K]) NewSyncOrderedHashMap(capacity int) SyncOrderedHashMap {
	return NewSyncOrderedNativeKeyHashMap[K](capacity)
}

func (h *SyncOrderedNativeKeyHashMap[K]) IterNative() *SyncOrderedNativeKeyHashMapIterator[K] {
	return NewSyncOrderedNativeKeyHashMapIterator(h)
}

func (h *SyncOrderedNativeKeyHashMap[K]) Iter() value.NativeIterator {
	return h.IterNative()
}

func (h *SyncOrderedNativeKeyHashMap[K]) IterSyncOrderedMap() value.NativeResettableIterator {
	return h.IterNative()
}

func (h *SyncOrderedNativeKeyHashMap[K]) GetIndex(key K) (int, bool) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetIndex(key)
}

func (h *SyncOrderedNativeKeyHashMap[K]) GetNativePair(key K) *value.NativePair[K, value.Value] {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetNativePair(key)
}

func (h *SyncOrderedNativeKeyHashMap[K]) GetPair(key K) value.Pair {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetPair(key)
}

func (h *SyncOrderedNativeKeyHashMap[K]) Get(key K) value.Value {
	h.Lock()
	defer h.Unlock()

	return h.Map.Get(key)
}

func (h *SyncOrderedNativeKeyHashMap[K]) GetOk(key K) (value.Value, bool) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetOk(key)
}

func (h *SyncOrderedNativeKeyHashMap[K]) At(index int) *value.NativePair[K, value.Value] {
	h.Lock()
	defer h.Unlock()

	return h.Map.At(index)
}

func (h *SyncOrderedNativeKeyHashMap[K]) AtInterface(index int) value.Pair {
	h.Lock()
	defer h.Unlock()

	return h.Map.AtInterface(index)
}

func (h *SyncOrderedNativeKeyHashMap[K]) GetAt(index int) (*value.NativePair[K, value.Value], value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetAt(index)
}

func (h *SyncOrderedNativeKeyHashMap[K]) GetAtInterface(index int) (value.Pair, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetAtInterface(index)
}

func (h *SyncOrderedNativeKeyHashMap[K]) TryAtVal(index value.Value) (value.Pair, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.TryAtVal(index)
}

func (h *SyncOrderedNativeKeyHashMap[K]) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetValUndefined(thread, key)
}

func (h *SyncOrderedNativeKeyHashMap[K]) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetValNil(thread, key)
}

func (h *SyncOrderedNativeKeyHashMap[K]) Set(key K, val value.Value) {
	h.Lock()
	defer h.Unlock()

	h.Map.Set(key, val)
}

func (h *SyncOrderedNativeKeyHashMap[K]) SetVal(thread *Thread, key, val value.Value) value.Value {
	h.Lock()
	defer h.Unlock()

	return h.Map.SetVal(thread, key, val)
}

func (h *SyncOrderedNativeKeyHashMap[K]) syncMergeResult(result OrderedHashRecord, err value.Value) (SyncOrderedHashMap, value.Value) {
	if err.IsNotUndefined() {
		return nil, err
	}
	switch m := result.(type) {
	case *OrderedNativeKeyHashMap[K]:
		return &SyncOrderedNativeKeyHashMap[K]{Map: m}, value.Undefined
	case *OrderedHashMapOfValue:
		return &SyncOrderedHashMapOfValue{Map: m}, value.Undefined
	default:
		return nil, value.Ref(value.Errorf(value.TypeErrorClass, "unexpected merge result"))
	}
}

func (h *SyncOrderedNativeKeyHashMap[K]) MergeVal(thread *Thread, other value.Value) (SyncOrderedHashMap, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncOrderedNativeKeyHashMap[K]:
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

func (h *SyncOrderedNativeKeyHashMap[K]) ContainsNativePair(thread *Thread, other *value.NativePair[K, value.Value]) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsNativePair(thread, other)
}

func (h *SyncOrderedNativeKeyHashMap[K]) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.Contains(thread, other)
}

func (h *SyncOrderedNativeKeyHashMap[K]) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsValue(thread, val)
}

func (h *SyncOrderedNativeKeyHashMap[K]) ContainsNativeKey(thread *Thread, key K) bool {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsNativeKey(thread, key)
}

func (h *SyncOrderedNativeKeyHashMap[K]) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsKey(thread, key)
}

func (h *SyncOrderedNativeKeyHashMap[K]) RemoveNativeKey(key K) bool {
	h.Lock()
	defer h.Unlock()

	return h.Map.RemoveNativeKey(key)
}

func (h *SyncOrderedNativeKeyHashMap[K]) RemoveKey(thread *Thread, key value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.RemoveKey(thread, key)
}

func (h *SyncOrderedNativeKeyHashMap[K]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncOrderedNativeKeyHashMap[K]:
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

func (h *SyncOrderedNativeKeyHashMap[K]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncOrderedNativeKeyHashMap[K]:
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

func (*SyncOrderedNativeKeyHashMap[K]) Class() *value.Class {
	return value.SyncOrderedHashMapClass
}

func (*SyncOrderedNativeKeyHashMap[K]) DirectClass() *value.Class {
	return value.SyncOrderedHashMapClass
}

func (*SyncOrderedNativeKeyHashMap[K]) SingletonClass() *value.Class {
	return nil
}

func (h *SyncOrderedNativeKeyHashMap[K]) Clone() *SyncOrderedNativeKeyHashMap[K] {
	h.Lock()
	defer h.Unlock()

	return &SyncOrderedNativeKeyHashMap[K]{
		Map: h.Map.Clone(),
	}
}

func (h *SyncOrderedNativeKeyHashMap[K]) Copy() value.Reference {
	return h.Clone()
}

func (h *SyncOrderedNativeKeyHashMap[K]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *SyncOrderedNativeKeyHashMap[K]) Error() string {
	return h.Inspect()
}

func (h *SyncOrderedNativeKeyHashMap[K]) Inspect() string {
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

func (*SyncOrderedNativeKeyHashMap[K]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncOrderedNativeKeyHashMap[K]) Length() int {
	h.Lock()
	defer h.Unlock()

	return h.Map.Length()
}

type SyncOrderedNativeKeyHashMapIterator[K value.ComparableValueInterface] struct {
	HashMap  *SyncOrderedNativeKeyHashMap[K]
	index    int
	snapshot []value.NativePair[K, value.Value]
	version  int
}

var _ value.NativeResettableIterator = &SyncOrderedNativeKeyHashMapIterator[value.String]{}

func NewSyncOrderedNativeKeyHashMapIterator[K value.ComparableValueInterface](hmap *SyncOrderedNativeKeyHashMap[K]) *SyncOrderedNativeKeyHashMapIterator[K] {
	iterator := &SyncOrderedNativeKeyHashMapIterator[K]{
		HashMap: hmap,
		version: hmap.Map.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *SyncOrderedNativeKeyHashMapIterator[K]) captureSnapshot() {
	h.HashMap.Lock()
	defer h.HashMap.Unlock()

	snapshot := make([]value.NativePair[K, value.Value], 0, h.HashMap.Map.Length())
	for k, v := range h.HashMap.Map.AllNative() {
		snapshot = append(snapshot, value.MakeNativePair(k, v))
	}
	h.snapshot = snapshot
}

func (*SyncOrderedNativeKeyHashMapIterator[K]) Class() *value.Class {
	return value.SyncOrderedHashMapIteratorClass
}

func (*SyncOrderedNativeKeyHashMapIterator[K]) DirectClass() *value.Class {
	return value.SyncOrderedHashMapIteratorClass
}

func (*SyncOrderedNativeKeyHashMapIterator[K]) SingletonClass() *value.Class {
	return nil
}

func (h *SyncOrderedNativeKeyHashMapIterator[K]) Copy() value.Reference {
	return &SyncOrderedNativeKeyHashMapIterator[K]{
		HashMap:  h.HashMap,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *SyncOrderedNativeKeyHashMapIterator[K]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *SyncOrderedNativeKeyHashMapIterator[K]) Error() string {
	return h.Inspect()
}

func (h *SyncOrderedNativeKeyHashMapIterator[K]) Inspect() string {
	return fmt.Sprintf("Std::Sync::OrderedHashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*SyncOrderedNativeKeyHashMapIterator[K]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncOrderedNativeKeyHashMapIterator[K]) Next() (p value.NativePair[K, value.Value], err value.Value) {
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

func (h *SyncOrderedNativeKeyHashMapIterator[K]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *SyncOrderedNativeKeyHashMapIterator[K]) Reset() {
	h.index = 0
	h.version = h.HashMap.Map.version
}
