package vm

import (
	"fmt"
	"iter"
	"maps"
	"strings"
	"unsafe"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type NativeHashMap[K value.ComparableValueInterface, V value.ValueInterface] struct {
	m       map[K]V
	version int
}

var _ HashMap = &NativeHashMap[value.String, value.String]{}

// UNSAFE! Cast a map with native go types to a map with corresponding Elk types.
// This is EXTREMELY unsafe, use it only if `IK`, OK` and `IV`, `OV` have the same
// underlying types eg. `unsafeCastNativeMap[string, uint8, value.String, value.UInt8](m)`, this will convert `map[string]uint8` to `map[value.String]value.UInt8`
func unsafeCastNativeMap[
	IK comparable,
	IV any,
	OK value.ComparableValueInterface,
	OV value.ValueInterface,
](m map[IK]IV) map[OK]OV {
	return *(*map[OK]OV)(unsafe.Pointer(&m))
}

// UNSAFE! Cast a map with native go types to a new Elk `NativeHashMap` with corresponding Elk types.
// This is EXTREMELY unsafe, use it only if `IK`, OK` and `IV`, `OV` have the same
// underlying types eg. `NewUnsafeCastNativeHashMap[string, uint8, value.String, value.UInt8](m)`, this will convert `map[string]uint8` to `*NativeHashMap[value.String, value.UInt8]`
func NewUnsafeCastNativeHashMap[
	IK comparable,
	IV any,
	OK value.ComparableValueInterface,
	OV value.ValueInterface,
](m map[IK]IV) *NativeHashMap[OK, OV] {
	return &NativeHashMap[OK, OV]{
		m: unsafeCastNativeMap[IK, IV, OK, OV](m),
	}
}

// Transform a map with native go types to a new Elk `NativeHashMap` with corresponding Elk types
// using the given function.
// eg.
//
//	TransformMapIntoNativeHashMap(m, func(k string, v uint8) (value.String, value.UInt8) {
//		return value.String(k), value.UInt8(v)
//	})
func TransformMapIntoNativeHashMap[
	IK comparable,
	IV any,
	OK value.ComparableValueInterface,
	OV value.ValueInterface,
](
	m map[IK]IV,
	fn func(k IK, v IV) (OK, OV),
) *NativeHashMap[OK, OV] {
	newMap := NewNativeHashMap[OK, OV](len(m))
	for k, v := range m {
		ok, ov := fn(k, v)
		newMap.m[ok] = ov
	}
	return newMap
}

func NewNativeHashMapFromMap[K value.ComparableValueInterface, V value.ValueInterface](m map[K]V) *NativeHashMap[K, V] {
	return &NativeHashMap[K, V]{
		m: m,
	}
}

func NewNativeHashMap[K value.ComparableValueInterface, V value.ValueInterface](capacity int) *NativeHashMap[K, V] {
	return &NativeHashMap[K, V]{
		m: make(map[K]V, capacity),
	}
}

func NewNativeHashMapWithElements[K value.ComparableValueInterface, V value.ValueInterface](elements ...value.NativePair[K, V]) *NativeHashMap[K, V] {
	return NewNativeHashMapWithElementsAndTotalCapacity(len(elements), elements...)
}

func NewNativeHashMapWithElementsAndTotalCapacity[K value.ComparableValueInterface, V value.ValueInterface](capacity int, elements ...value.NativePair[K, V]) *NativeHashMap[K, V] {
	m := NewNativeHashMap[K, V](capacity)
	for _, pair := range elements {
		m.Set(pair.NativeKey(), pair.NativeValue())
	}
	return m
}

func (h *NativeHashMap[K, V]) CloneHashMap(thread *Thread, capacity int) (HashMap, value.Value) {
	newMap := NewNativeHashMap[K, V](capacity)
	maps.Copy(newMap.m, h.m)
	return newMap, value.Undefined
}

func (h *NativeHashMap[K, V]) NewHashMap(capacity int) HashMap {
	return NewNativeHashMap[K, V](capacity)
}

func (h *NativeHashMap[K, V]) NewHashRecord(capacity int) HashRecord {
	return h.NewHashMap(capacity)
}

func (h *NativeHashMap[K, V]) CloneHashRecord(thread *Thread, capacity int) (HashRecord, value.Value) {
	return h.CloneHashMap(thread, capacity)
}

func (h *NativeHashMap[K, V]) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		for k, v := range h.m {
			pair := value.MakePairOfValue(k.ToValue(), v.ToValue())
			if !yield(pair) {
				return
			}
		}
	}
}

func (h *NativeHashMap[K, V]) AllNative() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range h.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (h *NativeHashMap[K, V]) Iterate() iter.Seq2[value.Value, value.Value] {
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

func (h *NativeHashMap[K, V]) IterMap() value.NativeResettableIterator {
	return NewNativeHashMapIterator(h)
}

func (h *NativeHashMap[K, V]) IterRecord() value.NativeResettableIterator {
	return h.IterMap()
}

func (h *NativeHashMap[K, V]) Iter() value.NativeIterator {
	return h.IterMap()
}

func (h *NativeHashMap[K, V]) Get(key K) (V, bool) {
	v, ok := h.m[key]
	return v, ok
}

func (h *NativeHashMap[K, V]) GetVal(thread *Thread, key value.Value) (value.Value, value.Value) {
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

func (h *NativeHashMap[K, V]) Set(key K, val V) {
	h.m[key] = val
}

func (h *NativeHashMap[K, V]) SetVal(thread *Thread, key, val value.Value) value.Value {
	k, ok := value.Downcast[K](key)
	if !ok {
		return value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}
	v, ok := value.Downcast[V](val)
	if !ok {
		return value.NewInvalidValueInTypedMap(h, val.Class()).ToValue()
	}

	h.m[k] = v
	return value.Undefined
}

func (h *NativeHashMap[K, V]) ConcatVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *NativeHashMap[K, V]:
		newMap := NewNativeHashMap[K, V](h.Length() + o.Length())
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

func (h *NativeHashMap[K, V]) ContainsNativePair(thread *Thread, other *value.NativePair[K, V]) (bool, value.Value) {
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

func (h *NativeHashMap[K, V]) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
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

func (h *NativeHashMap[K, V]) ContainsNativeValue(thread *Thread, val V) (bool, value.Value) {
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

func (h *NativeHashMap[K, V]) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	v, ok := value.Downcast[V](val)
	if !ok {
		return false, value.NewInvalidValueInTypedMap(h, val.Class()).ToValue()
	}
	return h.ContainsNativeValue(thread, v)
}

func (h *NativeHashMap[K, V]) ContainsNativeKey(thread *Thread, key K) bool {
	_, ok := h.Get(key)
	return ok
}

func (h *NativeHashMap[K, V]) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return false, value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}
	return h.ContainsNativeKey(thread, k), value.Undefined
}

func (h *NativeHashMap[K, V]) EqualNative(thread *Thread, other *NativeHashMap[K, V]) (bool, value.Value) {
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

func (h *NativeHashMap[K, V]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *NativeHashMap[K, V]:
		return h.EqualNative(thread, o)
	case HashMap:
		return HashRecordEqual(thread, o, h)
	}

	return false, value.Undefined
}

func (h *NativeHashMap[K, V]) LaxEqualNative(thread *Thread, other *NativeHashMap[K, V]) (bool, value.Value) {
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

func (h *NativeHashMap[K, V]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *NativeHashMap[K, V]:
		return h.LaxEqualNative(thread, o)
	case HashRecord:
		return HashRecordLaxEqual(thread, h, o)
	}

	return false, value.Undefined
}

func (*NativeHashMap[K, V]) Class() *value.Class {
	return value.HashMapClass
}

func (*NativeHashMap[K, V]) DirectClass() *value.Class {
	return value.HashMapClass
}

func (*NativeHashMap[K, V]) SingletonClass() *value.Class {
	return nil
}

func (h *NativeHashMap[K, V]) Clone() *NativeHashMap[K, V] {
	newMap := NewNativeHashMap[K, V](h.Length())
	maps.Copy(newMap.m, h.m)
	return newMap
}

func (h *NativeHashMap[K, V]) Copy() value.Reference {
	return h.Clone()
}

func (h *NativeHashMap[K, V]) ToValue() value.Value {
	return value.Ref(h)
}

func (h *NativeHashMap[K, V]) Error() string {
	return h.Inspect()
}

func (h *NativeHashMap[K, V]) Inspect() string {
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

func (*NativeHashMap[K, V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *NativeHashMap[K, V]) Length() int {
	return len(h.m)
}

type NativeHashMapIterator[K value.ComparableValueInterface, V value.ValueInterface] struct {
	HashMap  *NativeHashMap[K, V]
	index    int
	snapshot []value.NativePair[K, V]
	version  int
}

var _ value.NativeResettableIterator = &NativeHashMapIterator[value.String, value.String]{}

func NewNativeHashMapIterator[K value.ComparableValueInterface, V value.ValueInterface](hmap *NativeHashMap[K, V]) *NativeHashMapIterator[K, V] {
	iterator := &NativeHashMapIterator[K, V]{
		HashMap: hmap,
		version: hmap.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *NativeHashMapIterator[K, V]) captureSnapshot() {
	snapshot := make([]value.NativePair[K, V], 0, h.HashMap.Length())
	for k, v := range h.HashMap.AllNative() {
		snapshot = append(snapshot, value.MakeNativePair(k, v))
	}
	h.snapshot = snapshot
}

func (*NativeHashMapIterator[K, V]) Class() *value.Class {
	return value.HashMapIteratorClass
}

func (*NativeHashMapIterator[K, V]) DirectClass() *value.Class {
	return value.HashMapIteratorClass
}

func (*NativeHashMapIterator[K, V]) SingletonClass() *value.Class {
	return nil
}

func (h *NativeHashMapIterator[K, V]) Copy() value.Reference {
	return &NativeHashMapIterator[K, V]{
		HashMap:  h.HashMap,
		index:    h.index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *NativeHashMapIterator[K, V]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *NativeHashMapIterator[K, V]) Error() string {
	return h.Inspect()
}

func (h *NativeHashMapIterator[K, V]) Inspect() string {
	return fmt.Sprintf("Std::HashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*NativeHashMapIterator[K, V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *NativeHashMapIterator[K, V]) Next() (p value.NativePair[K, V], err value.Value) {
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

func (h *NativeHashMapIterator[K, V]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *NativeHashMapIterator[K, V]) Reset() {
	h.index = 0
}
