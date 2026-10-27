package vm

import (
	"fmt"
	"iter"
	"maps"
	"strings"
	"sync"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type SyncNativeHashMap[K value.ComparableValueInterface, V value.ValueInterface] struct {
	Map *NativeHashMap[K, V]
	mu  sync.Mutex
}

var _ SyncHashMap = &SyncNativeHashMap[value.String, value.String]{}

func NewSyncNativeHashMapFromMap[K value.ComparableValueInterface, V value.ValueInterface](m *NativeHashMap[K, V]) *SyncNativeHashMap[K, V] {
	return &SyncNativeHashMap[K, V]{
		Map: m,
	}
}

func NewSyncNativeHashMap[K value.ComparableValueInterface, V value.ValueInterface](capacity int) *SyncNativeHashMap[K, V] {
	return &SyncNativeHashMap[K, V]{
		Map: NewNativeHashMap[K, V](capacity),
	}
}

func (h *SyncNativeHashMap[K, V]) HashMap() HashMap {
	return h.Map
}

func (h *SyncNativeHashMap[K, V]) Lock() {
	h.mu.Lock()
}

func (h *SyncNativeHashMap[K, V]) Unlock() {
	h.mu.Unlock()
}

func (h *SyncNativeHashMap[K, V]) CloneSyncHashMap(thread *Thread, capacity int) (SyncHashMap, value.Value) {
	h.Lock()
	defer h.Unlock()

	newMap := NewSyncNativeHashMap[K, V](capacity)
	maps.Copy(newMap.Map.m, h.Map.m)
	return newMap, value.Undefined
}

func (h *SyncNativeHashMap[K, V]) NewSyncHashMap(capacity int) SyncHashMap {
	return NewSyncNativeHashMap[K, V](capacity)
}

func (h *SyncNativeHashMap[K, V]) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		h.Lock()
		defer h.Unlock()

		for k, v := range h.Map.m {
			pair := value.MakePairOfValue(k.ToValue(), v.ToValue())
			if !yield(pair) {
				return
			}
		}
	}
}

func (h *SyncNativeHashMap[K, V]) AllNative() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		h.Lock()
		defer h.Unlock()

		for k, v := range h.Map.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (h *SyncNativeHashMap[K, V]) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		h.Lock()
		defer h.Unlock()

		originalVersion := h.Map.version

		for k, v := range h.Map.m {
			if originalVersion != h.Map.version {
				yield(value.Undefined, value.NewMutationDuringIterationError(h.Class().Name()).ToValue())
				return
			}

			pair := value.NewNativePair(k, v)
			if !yield(pair.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (h *SyncNativeHashMap[K, V]) IterSyncMap() value.NativeResettableIterator {
	return NewSyncNativeHashMapIterator(h)
}

func (h *SyncNativeHashMap[K, V]) Iter() value.NativeIterator {
	return h.IterSyncMap()
}

func (h *SyncNativeHashMap[K, V]) GetOk(key K) (V, bool) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetOk(key)
}

func (h *SyncNativeHashMap[K, V]) Get(key K) value.Value {
	h.Lock()
	defer h.Unlock()

	return h.Map.Get(key)
}

func (h *SyncNativeHashMap[K, V]) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.GetValUndefined(thread, key)
}

func (h *SyncNativeHashMap[K, V]) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.GetValNil(thread, key)
}

func (h *SyncNativeHashMap[K, V]) Set(key K, val V) {
	h.Lock()
	defer h.Unlock()

	h.Map.Set(key, val)
}

func (h *SyncNativeHashMap[K, V]) SetVal(thread *Thread, key, val value.Value) value.Value {
	h.Lock()
	defer h.Unlock()

	return h.Map.SetVal(thread, key, val)
}

func (h *SyncNativeHashMap[K, V]) MergeVal(thread *Thread, other value.Value) (SyncHashMap, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncNativeHashMap[K, V]:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		newMap := NewSyncNativeHashMap[K, V](h.Length() + o.Length())
		maps.Copy(newMap.Map.m, h.Map.m)
		maps.Copy(newMap.Map.m, o.Map.m)
		return newMap, value.Undefined
	case SyncHashMap:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		newMap := NewSyncHashMapOfValue(h.Length() + o.Length())

		err := HashMapOfValueCopyInterface(thread, newMap.Map, h.Map)
		if err.IsNotUndefined() {
			return nil, err
		}

		err = HashMapOfValueCopyInterface(thread, newMap.Map, o.HashMap())
		if err.IsNotUndefined() {
			return nil, err
		}

		return newMap, value.Undefined
	case *NativeHashMap[K, V]:
		h.Lock()
		defer h.Unlock()

		newMap := NewSyncNativeHashMap[K, V](h.Length() + o.Length())
		maps.Copy(newMap.Map.m, h.Map.m)
		maps.Copy(newMap.Map.m, o.m)
		return newMap, value.Undefined
	case Record:
		h.Lock()
		defer h.Unlock()

		newMap := NewSyncHashMapOfValue(h.Length() + o.Length())

		err := HashMapOfValueCopyInterface(thread, newMap.Map, h.Map)
		if err.IsNotUndefined() {
			return nil, err
		}

		err = HashMapOfValueCopyInterface(thread, newMap.Map, o)
		if err.IsNotUndefined() {
			return nil, err
		}

		return newMap, value.Undefined
	}

	return nil, value.Ref(value.Errorf(value.TypeErrorClass, "cannot concat %s with map %s", other.Inspect(), h.Inspect()))
}

func (h *SyncNativeHashMap[K, V]) ContainsNativePair(thread *Thread, other *value.NativePair[K, V]) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsNativePair(thread, other)
}

func (h *SyncNativeHashMap[K, V]) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.Contains(thread, other)
}

func (h *SyncNativeHashMap[K, V]) ContainsNativeValue(thread *Thread, val V) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsNativeValue(thread, val)
}

func (h *SyncNativeHashMap[K, V]) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsValue(thread, val)
}

func (h *SyncNativeHashMap[K, V]) ContainsNativeKey(thread *Thread, key K) bool {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsNativeKey(thread, key)
}

func (h *SyncNativeHashMap[K, V]) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsKey(thread, key)
}

func (h *SyncNativeHashMap[K, V]) RemoveNativeKey(thread *Thread, key K) bool {
	h.Lock()
	defer h.Unlock()

	return h.Map.RemoveNativeKey(thread, key)
}

func (h *SyncNativeHashMap[K, V]) RemoveKey(thread *Thread, key value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.RemoveKey(thread, key)
}

func (h *SyncNativeHashMap[K, V]) EqualNative(thread *Thread, other *SyncNativeHashMap[K, V]) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	other.Lock()
	defer other.Unlock()

	return h.Map.EqualNative(thread, other.Map)
}

func (h *SyncNativeHashMap[K, V]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncNativeHashMap[K, V]:
		return h.EqualNative(thread, o)
	case SyncHashMap:
		h.Lock()
		defer h.Lock()

		o.Lock()
		defer o.Lock()

		return HashRecordEqual(thread, h.Map, o.HashMap())
	}

	return false, value.Undefined
}

func (h *SyncNativeHashMap[K, V]) LaxEqualNative(thread *Thread, other *SyncNativeHashMap[K, V]) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	other.Lock()
	defer other.Unlock()

	return h.Map.LaxEqualNative(thread, other.Map)
}

func (h *SyncNativeHashMap[K, V]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncNativeHashMap[K, V]:
		return h.LaxEqualNative(thread, o)
	case SyncHashMap:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return HashRecordLaxEqual(thread, h.Map, o.HashMap())
	case HashRecord:
		h.Lock()
		defer h.Unlock()

		return HashRecordLaxEqual(thread, h.Map, o)
	}

	return false, value.Undefined
}

func (*SyncNativeHashMap[K, V]) Class() *value.Class {
	return value.SyncHashMapClass
}

func (*SyncNativeHashMap[K, V]) DirectClass() *value.Class {
	return value.SyncHashMapClass
}

func (*SyncNativeHashMap[K, V]) SingletonClass() *value.Class {
	return nil
}

func (h *SyncNativeHashMap[K, V]) Clone() *SyncNativeHashMap[K, V] {
	return NewSyncNativeHashMapFromMap(h.Map.Clone())
}

func (h *SyncNativeHashMap[K, V]) Copy() value.Reference {
	return h.Clone()
}

func (h *SyncNativeHashMap[K, V]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *SyncNativeHashMap[K, V]) Error() string {
	return h.Inspect()
}

func (h *SyncNativeHashMap[K, V]) Inspect() string {
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
	for key, val := range h.Map.m {
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

	buff.WriteString("Std::Sync::HashMap{")
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

func (*SyncNativeHashMap[K, V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncNativeHashMap[K, V]) Length() int {
	return h.Map.Length()
}

type SyncNativeHashMapIterator[K value.ComparableValueInterface, V value.ValueInterface] struct {
	HashMap  *SyncNativeHashMap[K, V]
	index    int
	snapshot []value.NativePair[K, V]
	version  int
}

var _ value.NativeResettableIterator = &SyncNativeHashMapIterator[value.String, value.String]{}

func NewSyncNativeHashMapIterator[K value.ComparableValueInterface, V value.ValueInterface](hmap *SyncNativeHashMap[K, V]) *SyncNativeHashMapIterator[K, V] {
	iterator := &SyncNativeHashMapIterator[K, V]{
		HashMap: hmap,
		version: hmap.Map.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *SyncNativeHashMapIterator[K, V]) captureSnapshot() {
	h.HashMap.Lock()
	defer h.HashMap.Unlock()

	snapshot := make([]value.NativePair[K, V], 0, h.HashMap.Length())
	for k, v := range h.HashMap.Map.AllNative() {
		snapshot = append(snapshot, value.MakeNativePair(k, v))
	}
	h.snapshot = snapshot
}

func (*SyncNativeHashMapIterator[K, V]) Class() *value.Class {
	return value.SyncHashMapIteratorClass
}

func (*SyncNativeHashMapIterator[K, V]) DirectClass() *value.Class {
	return value.SyncHashMapIteratorClass
}

func (*SyncNativeHashMapIterator[K, V]) SingletonClass() *value.Class {
	return nil
}

func (h *SyncNativeHashMapIterator[K, V]) Copy() value.Reference {
	return &SyncNativeHashMapIterator[K, V]{
		HashMap:  h.HashMap,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *SyncNativeHashMapIterator[K, V]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *SyncNativeHashMapIterator[K, V]) Error() string {
	return h.Inspect()
}

func (h *SyncNativeHashMapIterator[K, V]) Inspect() string {
	return fmt.Sprintf("Std::Sync::HashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*SyncNativeHashMapIterator[K, V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncNativeHashMapIterator[K, V]) Next() (p value.NativePair[K, V], err value.Value) {
	if h.version != h.HashMap.Map.version {
		return p, value.NewMutationDuringIterationError(h.Class().Name()).ToValue()
	}
	if h.index >= len(h.snapshot) {
		return p, value.S(symbol.L_stop_iteration).ToValue()
	}

	v := h.snapshot[h.index]
	h.index++
	return v, value.Undefined
}

func (h *SyncNativeHashMapIterator[K, V]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *SyncNativeHashMapIterator[K, V]) Reset() {
	h.index = 0
}
