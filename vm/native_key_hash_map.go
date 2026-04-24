package vm

import (
	"fmt"
	"iter"
	"maps"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type NativeKeyHashMap[K value.ComparableValueInterface] struct {
	m       map[K]value.Value
	version int
}

var _ HashMap = &NativeKeyHashMap[value.String]{}

// Transform a map with native go types to a new Elk `NativeKeyHashMap` with corresponding Elk types
// using the given function.
// eg.
//
//	TransformMapIntoNativeKeyHashMap(m, func(k string, v uint8) (value.String, value.Value) {
//		return value.String(k), value.UInt8(v).ToValue()
//	})
func TransformMapIntoNativeKeyHashMap[
	IK comparable,
	IV any,
	OK value.ComparableValueInterface,
](
	m map[IK]IV,
	fn func(k IK, v IV) (OK, value.Value),
) *NativeKeyHashMap[OK] {
	newMap := NewNativeKeyHashMap[OK](len(m))
	for k, v := range m {
		ok, ov := fn(k, v)
		newMap.m[ok] = ov
	}
	return newMap
}

func NewNativeKeyHashMapFromMap[K value.ComparableValueInterface](m map[K]value.Value) *NativeKeyHashMap[K] {
	return &NativeKeyHashMap[K]{
		m: m,
	}
}

func NewNativeKeyHashMap[K value.ComparableValueInterface](capacity int) *NativeKeyHashMap[K] {
	return &NativeKeyHashMap[K]{
		m: make(map[K]value.Value, capacity),
	}
}

func NewNativeKeyHashMapWithElements[K value.ComparableValueInterface](elements ...value.NativePair[K, value.Value]) *NativeKeyHashMap[K] {
	return NewNativeKeyHashMapWithElementsAndTotalCapacity(len(elements), elements...)
}

func NewNativeKeyHashMapWithElementsAndTotalCapacity[K value.ComparableValueInterface](capacity int, elements ...value.NativePair[K, value.Value]) *NativeKeyHashMap[K] {
	m := NewNativeKeyHashMap[K](capacity)
	for _, pair := range elements {
		m.Set(pair.NativeKey(), pair.NativeValue())
	}
	return m
}

func (h *NativeKeyHashMap[K]) CloneHashMap(thread *Thread, capacity int) (HashMap, value.Value) {
	newMap := NewNativeKeyHashMap[K](capacity)
	maps.Copy(newMap.m, h.m)
	return newMap, value.Undefined
}

func (h *NativeKeyHashMap[K]) NewHashMap(capacity int) HashMap {
	return NewNativeKeyHashMap[K](capacity)
}

func (h *NativeKeyHashMap[K]) NewHashRecord(capacity int) HashRecord {
	return h.NewHashMap(capacity)
}

func (h *NativeKeyHashMap[K]) CloneHashRecord(thread *Thread, capacity int) (HashRecord, value.Value) {
	return h.CloneHashMap(thread, capacity)
}

func (h *NativeKeyHashMap[K]) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		for k, v := range h.m {
			pair := value.MakePairOfValue(k.ToValue(), v.ToValue())
			if !yield(pair) {
				return
			}
		}
	}
}

func (h *NativeKeyHashMap[K]) AllNative() iter.Seq2[K, value.Value] {
	return func(yield func(K, value.Value) bool) {
		for k, v := range h.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (h *NativeKeyHashMap[K]) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		originalVersion := h.version

		for k, v := range h.m {
			if originalVersion != h.version {
				yield(value.Undefined, value.NewMutationDuringIterationError(h.Class().Name).ToValue())
				return
			}

			pair := value.NewNativePair(k, v)
			if !yield(pair.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (h *NativeKeyHashMap[K]) IterMap() value.NativeResettableIterator {
	return NewNativeKeyHashMapIterator(h)
}

func (h *NativeKeyHashMap[K]) IterRecord() value.NativeResettableIterator {
	return h.IterMap()
}

func (h *NativeKeyHashMap[K]) Iter() value.NativeIterator {
	return h.IterMap()
}

func (h *NativeKeyHashMap[K]) Get(key K) (value.Value, bool) {
	v, ok := h.m[key]
	return v, ok
}

func (h *NativeKeyHashMap[K]) GetVal(thread *Thread, key value.Value) (value.Value, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return value.Undefined, value.Undefined
	}
	v, ok := h.Get(k)
	if !ok {
		return value.Undefined, value.Undefined
	}
	return v.ToValue(), value.Undefined
}

func (h *NativeKeyHashMap[K]) Set(key K, val value.Value) {
	h.m[key] = val
}

func (h *NativeKeyHashMap[K]) SetVal(thread *Thread, key, val value.Value) value.Value {
	k, ok := value.Downcast[K](key)
	if !ok {
		return value.NewInvalidKeyInTypedMap(h, k.Class()).ToValue()
	}

	h.m[k] = val
	return value.Undefined
}

func (h *NativeKeyHashMap[K]) ConcatVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *NativeKeyHashMap[K]:
		newMap := NewNativeKeyHashMap[K](h.Length() + o.Length())
		maps.Copy(newMap.m, h.m)
		maps.Copy(newMap.m, o.m)
		return newMap.ToValue(), value.Undefined
	case HashRecord:
		newMap := NewHashMapOfValue(h.Length() + o.Length())

		err := HashMapOfValueCopyInterface(thread, newMap, h)
		if err.IsNotUndefined() {
			return value.Undefined, err
		}

		err = HashMapOfValueCopyInterface(thread, newMap, o)
		if err.IsNotUndefined() {
			return value.Undefined, err
		}

		return newMap.ToValue(), value.Undefined
	}

	return value.Undefined, value.Ref(value.Errorf(value.TypeErrorClass, "cannot concat %s with map %s", other.Inspect(), h.Inspect()))
}

func (h *NativeKeyHashMap[K]) ContainsNativePair(thread *Thread, other *value.NativePair[K, value.Value]) (bool, value.Value) {
	v, ok := h.Get(other.NativeKey())
	if !ok {
		return false, value.Undefined
	}

	eq, err := Equal(thread, v.ToValue(), other.Value())
	if err.IsNotUndefined() {
		return false, err
	}

	return value.Truthy(eq), value.Undefined
}

func (h *NativeKeyHashMap[K]) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
	v, err := h.GetVal(thread, other.Key())
	if err.IsNotUndefined() {
		return false, err
	}

	eq, err := Equal(thread, v, other.Value())
	if err.IsNotUndefined() {
		return false, err
	}

	return value.Truthy(eq), value.Undefined
}

func (h *NativeKeyHashMap[K]) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	for _, hval := range h.m {
		eq, err := Equal(thread, hval.ToValue(), val.ToValue())
		if err.IsNotUndefined() {
			return false, err
		}

		if value.Truthy(eq) {
			return true, value.Undefined
		}
	}

	return false, value.Undefined
}

func (h *NativeKeyHashMap[K]) ContainsNativeKey(thread *Thread, key K) bool {
	_, ok := h.Get(key)
	return ok
}

func (h *NativeKeyHashMap[K]) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return false, value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}
	return h.ContainsNativeKey(thread, k), value.Undefined
}

func (h *NativeKeyHashMap[K]) EqualNative(thread *Thread, other *NativeKeyHashMap[K]) (bool, value.Value) {
	if h == other {
		return true, value.Undefined
	}
	if h.Length() != other.Length() {
		return false, value.Undefined
	}

	for hkey, hval := range h.m {
		oval := other.m[hkey]
		eqVal, err := Equal(thread, hval.ToValue(), oval.ToValue())
		if !err.IsUndefined() {
			return false, err
		}
		equal := value.Truthy(eqVal)
		if !equal {
			return false, value.Undefined
		}
	}

	return true, value.Undefined
}

func (h *NativeKeyHashMap[K]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *NativeKeyHashMap[K]:
		return h.EqualNative(thread, o)
	case HashMap:
		return HashRecordEqual(thread, o, h)
	}

	return false, value.Undefined
}

func (h *NativeKeyHashMap[K]) LaxEqualNative(thread *Thread, other *NativeKeyHashMap[K]) (bool, value.Value) {
	if h == other {
		return true, value.Undefined
	}
	if h.Length() != other.Length() {
		return false, value.Undefined
	}

	for hkey, hval := range h.m {
		oval := other.m[hkey]
		eqVal, err := LaxEqual(thread, hval.ToValue(), oval.ToValue())
		if !err.IsUndefined() {
			return false, err
		}
		equal := value.Truthy(eqVal)
		if !equal {
			return false, value.Undefined
		}
	}

	return true, value.Undefined
}

func (h *NativeKeyHashMap[K]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *NativeKeyHashMap[K]:
		return h.LaxEqualNative(thread, o)
	case HashRecord:
		return HashRecordLaxEqual(thread, h, o)
	}

	return false, value.Undefined
}

func (*NativeKeyHashMap[K]) Class() *value.Class {
	return value.HashMapClass
}

func (*NativeKeyHashMap[K]) DirectClass() *value.Class {
	return value.HashMapClass
}

func (*NativeKeyHashMap[K]) SingletonClass() *value.Class {
	return nil
}

func (h *NativeKeyHashMap[K]) Clone() *NativeKeyHashMap[K] {
	newMap := NewNativeKeyHashMap[K](h.Length())
	maps.Copy(newMap.m, h.m)
	return newMap
}

func (h *NativeKeyHashMap[K]) Copy() value.Reference {
	return h.Clone()
}

func (h *NativeKeyHashMap[K]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *NativeKeyHashMap[K]) Error() string {
	return h.Inspect()
}

func (h *NativeKeyHashMap[K]) Inspect() string {
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
	for key, val := range h.m {
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

	buff.WriteRune('{')
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

func (*NativeKeyHashMap[K]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *NativeKeyHashMap[K]) Length() int {
	return len(h.m)
}

type NativeKeyHashMapIterator[K value.ComparableValueInterface] struct {
	HashMap  *NativeKeyHashMap[K]
	index    int
	snapshot []value.NativePair[K, value.Value]
	version  int
}

var _ value.NativeResettableIterator = &NativeHashMapIterator[value.String, value.String]{}

func NewNativeKeyHashMapIterator[K value.ComparableValueInterface](hmap *NativeKeyHashMap[K]) *NativeKeyHashMapIterator[K] {
	iterator := &NativeKeyHashMapIterator[K]{
		HashMap: hmap,
		version: hmap.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *NativeKeyHashMapIterator[K]) captureSnapshot() {
	snapshot := make([]value.NativePair[K, value.Value], 0, h.HashMap.Length())
	for k, v := range h.HashMap.AllNative() {
		snapshot = append(snapshot, value.MakeNativePair(k, v))
	}
	h.snapshot = snapshot
}

func (*NativeKeyHashMapIterator[K]) Class() *value.Class {
	return value.HashMapIteratorClass
}

func (*NativeKeyHashMapIterator[K]) DirectClass() *value.Class {
	return value.HashMapIteratorClass
}

func (*NativeKeyHashMapIterator[K]) SingletonClass() *value.Class {
	return nil
}

func (h *NativeKeyHashMapIterator[K]) Copy() value.Reference {
	return &NativeKeyHashMapIterator[K]{
		HashMap:  h.HashMap,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *NativeKeyHashMapIterator[K]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *NativeKeyHashMapIterator[K]) Error() string {
	return h.Inspect()
}

func (h *NativeKeyHashMapIterator[K]) Inspect() string {
	return fmt.Sprintf("Std::HashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*NativeKeyHashMapIterator[K]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *NativeKeyHashMapIterator[K]) Next() (p value.NativePair[K, value.Value], err value.Value) {
	if h.version != h.HashMap.version {
		return p, value.NewMutationDuringIterationError(h.Class().Name).ToValue()
	}
	if h.index >= len(h.snapshot) {
		return p, symbol.L_stop_iteration.ToValue()
	}

	v := h.snapshot[h.index]
	h.index++
	return v, value.Undefined
}

func (h *NativeKeyHashMapIterator[K]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *NativeKeyHashMapIterator[K]) Reset() {
	h.index = 0
}
