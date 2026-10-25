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

type SyncNativeKeyHashMap[K value.ComparableValueInterface] struct {
	Map *NativeKeyHashMap[K]
	mu  sync.Mutex
}

var _ SyncHashMap = &SyncNativeKeyHashMap[value.String]{}

func NewSyncNativeKeyHashMapFromMap[K value.ComparableValueInterface](m map[K]value.Value) *SyncNativeKeyHashMap[K] {
	return &SyncNativeKeyHashMap[K]{
		Map: NewNativeKeyHashMapFromMap(m),
	}
}

func NewSyncNativeKeyHashMap[K value.ComparableValueInterface](capacity int) *SyncNativeKeyHashMap[K] {
	return &SyncNativeKeyHashMap[K]{
		Map: NewNativeKeyHashMap[K](capacity),
	}
}

func NewSyncNativeKeyHashMapWithElements[K value.ComparableValueInterface](elements ...value.NativePair[K, value.Value]) *SyncNativeKeyHashMap[K] {
	return NewSyncNativeKeyHashMapWithElementsAndTotalCapacity(len(elements), elements...)
}

func NewSyncNativeKeyHashMapWithElementsAndTotalCapacity[K value.ComparableValueInterface](capacity int, elements ...value.NativePair[K, value.Value]) *SyncNativeKeyHashMap[K] {
	m := NewSyncNativeKeyHashMap[K](capacity)
	for _, pair := range elements {
		m.Map.Set(pair.NativeKey(), pair.NativeValue())
	}
	return m
}

func (h *SyncNativeKeyHashMap[K]) HashMap() HashMap {
	return h.Map
}

func (h *SyncNativeKeyHashMap[K]) CloneSyncHashMap(thread *Thread, capacity int) (SyncHashMap, value.Value) {
	h.Lock()
	defer h.Unlock()

	newMap := NewSyncNativeKeyHashMap[K](capacity)
	maps.Copy(newMap.Map.m, h.Map.m)
	return newMap, value.Undefined
}

func (h *SyncNativeKeyHashMap[K]) Lock() {
	h.mu.Lock()
}

func (h *SyncNativeKeyHashMap[K]) Unlock() {
	h.mu.Unlock()
}

func (h *SyncNativeKeyHashMap[K]) NewSyncHashMap(capacity int) SyncHashMap {
	return NewSyncNativeKeyHashMap[K](capacity)
}

func (h *SyncNativeKeyHashMap[K]) All() iter.Seq[value.PairOfValue] {
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

func (h *SyncNativeKeyHashMap[K]) AllNative() iter.Seq2[K, value.Value] {
	return func(yield func(K, value.Value) bool) {
		h.Lock()
		defer h.Unlock()

		for k, v := range h.Map.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (h *SyncNativeKeyHashMap[K]) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		originalVersion := h.Map.version

		h.Lock()
		defer h.Unlock()

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

func (h *SyncNativeKeyHashMap[K]) IterSyncMap() value.NativeResettableIterator {
	return NewSyncNativeKeyHashMapIterator(h)
}

func (h *SyncNativeKeyHashMap[K]) Iter() value.NativeIterator {
	return h.IterSyncMap()
}

func (h *SyncNativeKeyHashMap[K]) Get(key K) value.Value {
	h.Lock()
	defer h.Unlock()

	return h.Map.Get(key)
}

func (h *SyncNativeKeyHashMap[K]) GetOk(key K) (value.Value, bool) {
	h.Lock()
	h.Unlock()

	return h.Map.GetOk(key)
}

func (h *SyncNativeKeyHashMap[K]) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
	h.Lock()
	h.Unlock()

	return h.Map.GetValUndefined(thread, key)
}

func (h *SyncNativeKeyHashMap[K]) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
	h.Lock()
	h.Unlock()

	return h.Map.GetValNil(thread, key)
}

func (h *SyncNativeKeyHashMap[K]) Set(key K, val value.Value) {
	h.Lock()
	defer h.Unlock()

	h.Map.Set(key, val)
}

func (h *SyncNativeKeyHashMap[K]) SetVal(thread *Thread, key, val value.Value) value.Value {
	h.Lock()
	defer h.Unlock()

	return h.Map.SetVal(thread, key, val)
}

func (h *SyncNativeKeyHashMap[K]) ConcatVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncNativeKeyHashMap[K]:
		newMap := NewSyncNativeKeyHashMap[K](h.Length() + o.Length())
		h.Lock()
		maps.Copy(newMap.Map.m, h.Map.m)
		h.Unlock()

		o.Lock()
		maps.Copy(newMap.Map.m, o.Map.m)
		o.Unlock()

		return newMap.ToValue(), value.Undefined
	case SyncHashMap:
		newMap := NewSyncHashMapOfValue(h.Length() + o.Length())

		h.Lock()
		err := HashMapOfValueCopyInterface(thread, newMap.Map, h.Map)
		h.Unlock()
		if err.IsNotUndefined() {
			return value.Undefined, err
		}

		o.Lock()
		err = HashMapOfValueCopyInterface(thread, newMap.Map, o)
		o.Unlock()
		if err.IsNotUndefined() {
			return value.Undefined, err
		}

		return newMap.ToValue(), value.Undefined
	case *NativeKeyHashMap[K]:
		h.Lock()
		newMap := NewSyncNativeKeyHashMap[K](h.Length() + o.Length())
		maps.Copy(newMap.Map.m, h.Map.m)
		h.Unlock()

		maps.Copy(newMap.Map.m, o.m)
		return newMap.ToValue(), value.Undefined
	case HashRecord:
		h.Lock()
		defer h.Unlock()

		newMap := NewSyncHashMapOfValue(h.Length() + o.Length())

		err := HashMapOfValueCopyInterface(thread, newMap.Map, h)
		if err.IsNotUndefined() {
			return value.Undefined, err
		}

		err = HashMapOfValueCopyInterface(thread, newMap.Map, o)
		if err.IsNotUndefined() {
			return value.Undefined, err
		}

		return newMap.ToValue(), value.Undefined
	}

	return value.Undefined, value.Ref(value.Errorf(value.TypeErrorClass, "cannot concat %s with map %s", other.Inspect(), h.Inspect()))
}

func (h *SyncNativeKeyHashMap[K]) ContainsNativePair(thread *Thread, other *value.NativePair[K, value.Value]) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.ContainsNativePair(thread, other)
}

func (h *SyncNativeKeyHashMap[K]) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Contains(thread, other)
}

func (h *SyncNativeKeyHashMap[K]) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsValue(thread, val)
}

func (h *SyncNativeKeyHashMap[K]) ContainsNativeKey(thread *Thread, key K) bool {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsNativeKey(thread, key)
}

func (h *SyncNativeKeyHashMap[K]) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsKey(thread, key)
}

func (h *SyncNativeKeyHashMap[K]) EqualNative(thread *Thread, other *SyncNativeKeyHashMap[K]) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	other.Lock()
	defer other.Unlock()

	return h.Map.EqualNative(thread, other.Map)
}

func (h *SyncNativeKeyHashMap[K]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncNativeKeyHashMap[K]:
		return h.EqualNative(thread, o)
	case SyncHashMap:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return HashRecordEqual(thread, h.Map, o.HashMap())
	}

	return false, value.Undefined
}

func (h *SyncNativeKeyHashMap[K]) LaxEqualNative(thread *Thread, other *SyncNativeKeyHashMap[K]) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	other.Lock()
	defer other.Unlock()

	return h.Map.LaxEqualNative(thread, other.Map)
}

func (h *SyncNativeKeyHashMap[K]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncNativeKeyHashMap[K]:
		return h.LaxEqualNative(thread, o)
	case SyncHashMap:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return HashRecordEqual(thread, h.Map, o.HashMap())
	case HashRecord:
		h.Lock()
		defer h.Unlock()

		return HashRecordLaxEqual(thread, h.Map, o)
	}

	return false, value.Undefined
}

func (*SyncNativeKeyHashMap[K]) Class() *value.Class {
	return value.SyncHashMapClass
}

func (*SyncNativeKeyHashMap[K]) DirectClass() *value.Class {
	return value.SyncHashMapClass
}

func (*SyncNativeKeyHashMap[K]) SingletonClass() *value.Class {
	return nil
}

func (h *SyncNativeKeyHashMap[K]) Clone() *SyncNativeKeyHashMap[K] {
	h.Lock()
	defer h.Unlock()

	newMap := NewSyncNativeKeyHashMap[K](h.Length())
	maps.Copy(newMap.Map.m, h.Map.m)
	return newMap
}

func (h *SyncNativeKeyHashMap[K]) Copy() value.Reference {
	return h.Clone()
}

func (h *SyncNativeKeyHashMap[K]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *SyncNativeKeyHashMap[K]) Error() string {
	return h.Inspect()
}

func (h *SyncNativeKeyHashMap[K]) Inspect() string {
	h.Lock()
	defer h.Unlock()

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

func (*SyncNativeKeyHashMap[K]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncNativeKeyHashMap[K]) Length() int {
	return h.Length()
}

type SyncNativeKeyHashMapIterator[K value.ComparableValueInterface] struct {
	HashMap  *SyncNativeKeyHashMap[K]
	index    int
	snapshot []value.NativePair[K, value.Value]
	version  int
}

var _ value.NativeResettableIterator = &SyncNativeKeyHashMapIterator[value.String]{}

func NewSyncNativeKeyHashMapIterator[K value.ComparableValueInterface](hmap *SyncNativeKeyHashMap[K]) *SyncNativeKeyHashMapIterator[K] {
	iterator := &SyncNativeKeyHashMapIterator[K]{
		HashMap: hmap,
		version: hmap.Map.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *SyncNativeKeyHashMapIterator[K]) captureSnapshot() {
	h.HashMap.Lock()
	defer h.HashMap.Unlock()

	snapshot := make([]value.NativePair[K, value.Value], 0, h.HashMap.Length())
	for k, v := range h.HashMap.Map.AllNative() {
		snapshot = append(snapshot, value.MakeNativePair(k, v))
	}
	h.snapshot = snapshot
}

func (*SyncNativeKeyHashMapIterator[K]) Class() *value.Class {
	return value.SyncHashMapIteratorClass
}

func (*SyncNativeKeyHashMapIterator[K]) DirectClass() *value.Class {
	return value.SyncHashMapIteratorClass
}

func (*SyncNativeKeyHashMapIterator[K]) SingletonClass() *value.Class {
	return nil
}

func (h *SyncNativeKeyHashMapIterator[K]) Copy() value.Reference {
	return &SyncNativeKeyHashMapIterator[K]{
		HashMap:  h.HashMap,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *SyncNativeKeyHashMapIterator[K]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *SyncNativeKeyHashMapIterator[K]) Error() string {
	return h.Inspect()
}

func (h *SyncNativeKeyHashMapIterator[K]) Inspect() string {
	return fmt.Sprintf("Std::Sync::HashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*SyncNativeKeyHashMapIterator[K]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncNativeKeyHashMapIterator[K]) Next() (p value.NativePair[K, value.Value], err value.Value) {
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

func (h *SyncNativeKeyHashMapIterator[K]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *SyncNativeKeyHashMapIterator[K]) Reset() {
	h.index = 0
}
