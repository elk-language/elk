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

type SyncOrderedHashMapOfValue struct {
	Map *OrderedHashMapOfValue
	mu  sync.Mutex
}

var _ SyncOrderedHashMap = &SyncOrderedHashMapOfValue{}

func SyncOrderedHashMapOfValueConstructor(class *value.Class) value.Value {
	return value.Ref(NewSyncOrderedHashMapOfValue(5))
}

func NewSyncOrderedHashMapOfValue(capacity int) *SyncOrderedHashMapOfValue {
	return &SyncOrderedHashMapOfValue{
		Map: NewOrderedHashMapOfValue(capacity),
	}
}

func (h *SyncOrderedHashMapOfValue) OrderedHashMap() OrderedHashMap {
	return h.Map
}

func (h *SyncOrderedHashMapOfValue) Lock() {
	h.mu.Lock()
}

func (h *SyncOrderedHashMapOfValue) Unlock() {
	h.mu.Unlock()
}

func (h *SyncOrderedHashMapOfValue) Keys() iter.Seq[value.Value] {
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

func (h *SyncOrderedHashMapOfValue) All() iter.Seq[value.PairOfValue] {
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

func (h *SyncOrderedHashMapOfValue) Iterate() iter.Seq2[value.Value, value.Value] {
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

func (h *SyncOrderedHashMapOfValue) CloneSyncOrderedHashMap(thread *Thread, capacity int) (SyncOrderedHashMap, value.Value) {
	h.Lock()
	defer h.Unlock()

	cloned, err := h.Map.CloneOrderedHashMap(thread, capacity)
	if err.IsNotUndefined() {
		return nil, err
	}
	return &SyncOrderedHashMapOfValue{Map: cloned.(*OrderedHashMapOfValue)}, value.Undefined
}

func (h *SyncOrderedHashMapOfValue) NewSyncOrderedHashMap(capacity int) SyncOrderedHashMap {
	return NewSyncOrderedHashMapOfValue(capacity)
}

func (h *SyncOrderedHashMapOfValue) IterNative() *SyncOrderedHashMapOfValueIterator {
	return NewSyncOrderedHashMapOfValueIterator(h)
}

func (h *SyncOrderedHashMapOfValue) Iter() value.NativeIterator {
	return h.IterNative()
}

func (h *SyncOrderedHashMapOfValue) IterSyncOrderedMap() value.NativeResettableIterator {
	return h.IterNative()
}

func (h *SyncOrderedHashMapOfValue) TryAtVal(index value.Value) (value.Pair, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.TryAtVal(index)
}

func (h *SyncOrderedHashMapOfValue) AtInterface(index int) value.Pair {
	h.Lock()
	defer h.Unlock()

	return h.Map.AtInterface(index)
}

func (h *SyncOrderedHashMapOfValue) GetAtInterface(index int) (value.Pair, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetAtInterface(index)
}

func (h *SyncOrderedHashMapOfValue) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetValUndefined(thread, key)
}

func (h *SyncOrderedHashMapOfValue) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.GetValNil(thread, key)
}

func (h *SyncOrderedHashMapOfValue) SetVal(thread *Thread, key, val value.Value) value.Value {
	h.Lock()
	defer h.Unlock()

	return h.Map.SetVal(thread, key, val)
}

func (h *SyncOrderedHashMapOfValue) MergeVal(thread *Thread, other value.Value) (SyncOrderedHashMap, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncOrderedHashMapOfValue:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		result, err := h.Map.MergeVal(thread, o.Map.ToValue())
		if err.IsNotUndefined() {
			return nil, err
		}
		return &SyncOrderedHashMapOfValue{Map: result.(*OrderedHashMapOfValue)}, value.Undefined
	case SyncOrderedHashMap:
		h.Lock()
		defer h.Unlock()

		o.Lock()
		defer o.Unlock()

		result, err := h.Map.MergeVal(thread, o.OrderedHashMap().ToValue())
		if err.IsNotUndefined() {
			return nil, err
		}
		return &SyncOrderedHashMapOfValue{Map: result.(*OrderedHashMapOfValue)}, value.Undefined
	default:
		h.Lock()
		defer h.Unlock()

		result, err := h.Map.MergeVal(thread, other)
		if err.IsNotUndefined() {
			return nil, err
		}
		return &SyncOrderedHashMapOfValue{Map: result.(*OrderedHashMapOfValue)}, value.Undefined
	}
}

func (h *SyncOrderedHashMapOfValue) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.Contains(thread, other)
}

func (h *SyncOrderedHashMapOfValue) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsValue(thread, val)
}

func (h *SyncOrderedHashMapOfValue) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.ContainsKey(thread, key)
}

func (h *SyncOrderedHashMapOfValue) RemoveKey(thread *Thread, key value.Value) (bool, value.Value) {
	h.Lock()
	defer h.Unlock()

	return h.Map.RemoveKey(thread, key)
}

func (h *SyncOrderedHashMapOfValue) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncOrderedHashMapOfValue:
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

func (h *SyncOrderedHashMapOfValue) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *SyncOrderedHashMapOfValue:
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

func (*SyncOrderedHashMapOfValue) Class() *value.Class {
	return value.SyncOrderedHashMapClass
}

func (*SyncOrderedHashMapOfValue) DirectClass() *value.Class {
	return value.SyncOrderedHashMapClass
}

func (*SyncOrderedHashMapOfValue) SingletonClass() *value.Class {
	return nil
}

func (h *SyncOrderedHashMapOfValue) Clone() *SyncOrderedHashMapOfValue {
	h.Lock()
	defer h.Unlock()

	return &SyncOrderedHashMapOfValue{
		Map: h.Map.Clone(),
	}
}

func (h *SyncOrderedHashMapOfValue) Copy() value.Reference {
	return h.Clone()
}

func (h *SyncOrderedHashMapOfValue) ToValue() value.Value {
	return value.Ref(h)
}

func (h *SyncOrderedHashMapOfValue) Error() string {
	return h.Inspect()
}

func (h *SyncOrderedHashMapOfValue) Inspect() string {
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

func (*SyncOrderedHashMapOfValue) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncOrderedHashMapOfValue) Capacity() int {
	h.Lock()
	defer h.Unlock()

	return h.Map.Capacity()
}

func (h *SyncOrderedHashMapOfValue) LeftCapacity() int {
	h.Lock()
	defer h.Unlock()

	return h.Map.LeftCapacity()
}

func (h *SyncOrderedHashMapOfValue) Length() int {
	h.Lock()
	defer h.Unlock()

	return h.Map.Length()
}

type SyncOrderedHashMapOfValueIterator struct {
	HashMap  *SyncOrderedHashMapOfValue
	Index    int
	snapshot []value.PairOfValue
	version  int
}

var _ value.NativeResettableIterator = &SyncOrderedHashMapOfValueIterator{}

func NewSyncOrderedHashMapOfValueIterator(hmap *SyncOrderedHashMapOfValue) *SyncOrderedHashMapOfValueIterator {
	iterator := &SyncOrderedHashMapOfValueIterator{
		HashMap: hmap,
		version: hmap.Map.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *SyncOrderedHashMapOfValueIterator) captureSnapshot() {
	h.HashMap.Lock()
	defer h.HashMap.Unlock()

	snapshot := make([]value.PairOfValue, 0, h.HashMap.Map.Length())
	for pair := range h.HashMap.Map.All() {
		snapshot = append(snapshot, pair)
	}
	h.snapshot = snapshot
}

func (*SyncOrderedHashMapOfValueIterator) Class() *value.Class {
	return value.SyncOrderedHashMapIteratorClass
}

func (*SyncOrderedHashMapOfValueIterator) DirectClass() *value.Class {
	return value.SyncOrderedHashMapIteratorClass
}

func (*SyncOrderedHashMapOfValueIterator) SingletonClass() *value.Class {
	return nil
}

func (h *SyncOrderedHashMapOfValueIterator) Copy() value.Reference {
	return &SyncOrderedHashMapOfValueIterator{
		HashMap:  h.HashMap,
		Index:    h.Index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *SyncOrderedHashMapOfValueIterator) ToValue() value.Value {
	return value.Ref(i)
}

func (h *SyncOrderedHashMapOfValueIterator) Error() string {
	return h.Inspect()
}

func (h *SyncOrderedHashMapOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::Sync::OrderedHashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*SyncOrderedHashMapOfValueIterator) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *SyncOrderedHashMapOfValueIterator) NextValue() (value.Value, value.Value) {
	if h.version != h.HashMap.Map.version {
		return value.Undefined, value.NewMutationDuringIterationError(h.HashMap.Class().Name()).ToValue()
	}

	if h.Index >= len(h.snapshot) {
		return value.Undefined, value.S(symbol.L_stop_iteration).ToValue()
	}

	pair := &h.snapshot[h.Index]
	h.Index++
	return value.Ref(pair), value.Undefined
}

func (h *SyncOrderedHashMapOfValueIterator) Reset() {
	h.Index = 0
	h.version = h.HashMap.Map.version
}
