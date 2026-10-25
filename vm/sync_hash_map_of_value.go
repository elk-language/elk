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

type SyncHashMapOfValue struct {
	Map *HashMapOfValue
	mu  sync.Mutex
}

var _ SyncHashMap = &SyncHashMapOfValue{}

func SyncHashMapOfValueConstructor(class *value.Class) value.Value {
	return value.Ref(&SyncHashMapOfValue{})
}

func NewSyncHashMapOfValue(capacity int) *SyncHashMapOfValue {
	return &SyncHashMapOfValue{
		Map: NewHashMapOfValue(capacity),
	}
}

func NewSyncHashMapOfValueWithMap(m *HashMapOfValue) *SyncHashMapOfValue {
	return &SyncHashMapOfValue{
		Map: m,
	}
}

func (h *SyncHashMapOfValue) HashMap() HashMap {
	return h.Map
}

func (h *SyncHashMapOfValue) Lock() {
	h.Lock()
}

func (h *SyncHashMapOfValue) Unlock() {
	h.Unlock()
}

func (h *SyncHashMapOfValue) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		h.Lock()
		defer h.Unlock()

		for _, pair := range h.Map.Table {
			if pair.Key().IsUndefined() {
				continue
			}

			if !yield(pair) {
				return
			}
		}
	}
}

func (h *SyncHashMapOfValue) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		h.Lock()
		defer h.Unlock()

		originalVersion := h.Map.version
		for _, pair := range h.Map.Table {
			if pair.Key().IsUndefined() {
				continue
			}

			if originalVersion != h.Map.version {
				yield(value.Undefined, value.NewMutationDuringIterationError(h.Class().Name()).ToValue())
				return
			}

			if !yield(pair.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (h *SyncHashMapOfValue) CloneSyncHashMap(thread *Thread, capacity int) (SyncHashMap, value.Value) {
	h.Lock()
	defer h.Unlock()

	newMap := NewSyncHashMapOfValue(capacity)
	err := HashMapOfValueCopy(thread, newMap.Map, h.Map)
	if !err.IsUndefined() {
		return nil, err
	}

	return newMap, value.Undefined
}

func (h *SyncHashMapOfValue) NewSyncHashMap(capacity int) SyncHashMap {
	return NewSyncHashMapOfValue(capacity)
}

func (h *SyncHashMapOfValue) IterNative() *SyncHashMapOfValueIterator {
	return NewSyncHashMapOfValueIterator(h)
}

func (h *SyncHashMapOfValue) Iter() value.NativeIterator {
	return h.IterNative()
}

func (h *SyncHashMapOfValue) IterSyncMap() value.NativeResettableIterator {
	return h.IterNative()
}

func (h *SyncHashMapOfValue) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.Map.GetValUndefined(thread, key)
}

func (h *SyncHashMapOfValue) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetValNil(thread, key)
}

func (h *SyncHashMapOfValue) SetVal(thread *Thread, key, val value.Value) value.Value {
	h.Lock()
	defer h.Unlock()

	return h.Map.SetVal(thread, key, val)
}

func (h *SyncHashMapOfValue) ConcatVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashMapOfValue:
		h.Lock()
		defer h.Unlock()

		m, err := HashMapOfValueConcat(thread, h.Map, o)
		if err.IsNotUndefined() {
			return value.Undefined, err
		}
		return NewSyncHashMapOfValueWithMap(m).ToValue(), value.Undefined
	case *SyncHashMapOfValue:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		m, err := HashMapOfValueConcat(thread, h.Map, o.Map)
		if err.IsNotUndefined() {
			return value.Undefined, err
		}
		return NewSyncHashMapOfValueWithMap(m).ToValue(), value.Undefined
	case HashRecord:
		h.Lock()
		defer h.Unlock()

		m, err := HashMapOfValueConcatInterface(thread, h.Map, o)
		if err.IsNotUndefined() {
			return value.Undefined, err
		}
		return NewSyncHashMapOfValueWithMap(m).ToValue(), value.Undefined
	case SyncHashMap:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		m, err := HashMapOfValueConcatInterface(thread, h.Map, o.HashMap())
		if err.IsNotUndefined() {
			return value.Undefined, err
		}
		return NewSyncHashMapOfValueWithMap(m).ToValue(), value.Undefined
	}

	return value.Undefined, value.Ref(value.Errorf(value.TypeErrorClass, "cannot concat %s with map %s", other.Inspect(), h.Inspect()))
}

func (h *SyncHashMapOfValue) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.Contains(thread, other)
}

func (h *SyncHashMapOfValue) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsValue(thread, val)
}

func (h *SyncHashMapOfValue) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsKey(thread, key)
}

func (h *SyncHashMapOfValue) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncHashMapOfValue:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return HashMapOfValueEqual(thread, h.Map, o.Map)
	case SyncHashMap:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return HashMapOfValueEqualInterface(thread, h.Map, o.HashMap())
	}

	return false, value.Undefined
}

func (h *SyncHashMapOfValue) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncHashMapOfValue:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return HashMapOfValueLaxEqual(thread, h.Map, o.Map)
	case SyncHashMap:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		return HashMapOfValueLaxEqualInterface(thread, h.Map, o.HashMap())
	case *HashMapOfValue:
		h.Lock()
		defer h.Unlock()

		return HashMapOfValueLaxEqual(thread, h.Map, o)
	case HashMap:
		h.Lock()
		defer h.Unlock()

		return HashMapOfValueLaxEqualInterface(thread, h.Map, o)
	}

	return false, value.Undefined
}

func (h *SyncHashMapOfValue) Grow(thread *Thread, newSlots int) value.Value {
	h.Lock()
	defer h.Unlock()

	return h.Map.Grow(thread, newSlots)
}

func (*SyncHashMapOfValue) Class() *value.Class {
	return value.SyncHashMapClass
}

func (*SyncHashMapOfValue) DirectClass() *value.Class {
	return value.SyncHashMapClass
}

func (*SyncHashMapOfValue) SingletonClass() *value.Class {
	return nil
}

func (h *SyncHashMapOfValue) Clone() *SyncHashMapOfValue {
	return &SyncHashMapOfValue{
		Map: h.Map.Clone(),
	}
}

func (h *SyncHashMapOfValue) Copy() value.Reference {
	return h.Clone()
}

func (h *SyncHashMapOfValue) ToValue() value.Value {
	return value.Ref(h)
}

func (h *SyncHashMapOfValue) Error() string {
	return h.Inspect()
}

func (h *SyncHashMapOfValue) Inspect() string {
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
	for entry := range h.Map.All() {
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

func (*SyncHashMapOfValue) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncHashMapOfValue) Capacity() int {
	return h.Map.Capacity()
}

func (h *SyncHashMapOfValue) LeftCapacity() int {
	h.Lock()
	defer h.Unlock()

	return h.Map.LeftCapacity()
}

func (h *SyncHashMapOfValue) Length() int {
	return h.Map.Length()
}

type SyncHashMapOfValueIterator struct {
	HashMap  *SyncHashMapOfValue
	Index    int
	snapshot []value.PairOfValue
	version  int
}

var _ value.NativeResettableIterator = &SyncHashMapOfValueIterator{}

func NewSyncHashMapOfValueIterator(hmap *SyncHashMapOfValue) *SyncHashMapOfValueIterator {
	iterator := &SyncHashMapOfValueIterator{
		HashMap: hmap,
		version: hmap.Map.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *SyncHashMapOfValueIterator) captureSnapshot() {
	h.HashMap.Lock()
	defer h.HashMap.Unlock()

	snapshot := make([]value.PairOfValue, 0, h.HashMap.Length())
	for pair := range h.HashMap.Map.All() {
		snapshot = append(snapshot, pair)
	}
	h.snapshot = snapshot
}

func (*SyncHashMapOfValueIterator) Class() *value.Class {
	return value.SyncHashMapIteratorClass
}

func (*SyncHashMapOfValueIterator) DirectClass() *value.Class {
	return value.SyncHashMapIteratorClass
}

func (*SyncHashMapOfValueIterator) SingletonClass() *value.Class {
	return nil
}

func (h *SyncHashMapOfValueIterator) Copy() value.Reference {
	return &SyncHashMapOfValueIterator{
		HashMap:  h.HashMap,
		Index:    h.Index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *SyncHashMapOfValueIterator) ToValue() value.Value {
	return value.Ref(i)
}

func (h *SyncHashMapOfValueIterator) Error() string {
	return h.Inspect()
}

func (h *SyncHashMapOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::Sync::HashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*SyncHashMapOfValueIterator) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncHashMapOfValueIterator) NextValue() (p value.Value, err value.Value) {
	if h.version != h.HashMap.Map.version {
		return p, value.NewMutationDuringIterationError(h.Class().Name()).ToValue()
	}
	if h.Index >= len(h.snapshot) {
		return p, value.S(symbol.L_stop_iteration).ToValue()
	}

	v := h.snapshot[h.Index]
	h.Index++
	return v.Key(), value.Undefined
}

func (h *SyncHashMapOfValueIterator) Reset() {
	h.Index = 0
}
