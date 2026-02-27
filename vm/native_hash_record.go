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

type NativeHashRecord[K value.ComparableValueInterface, V value.ValueInterface] map[K]V

var _ HashRecord = NativeHashRecord[value.String, value.String]{}

// UNSAFE! Cast a map with native go types to an Elk `NativeHashRecord` with corresponding Elk types.
// This is EXTREMELY unsafe, use it only if `IK`, OK` and `IV`, `OV` have the same
// underlying types eg. `CastNativeHashRecord[string, uint8, value.String, value.UInt8](m)`, this will convert `map[string]uint8` to `NativeHashRecord[value.String, value.UInt8]`
func CastNativeHashRecord[
	IK comparable,
	IV any,
	OK value.ComparableValueInterface,
	OV value.ValueInterface,
](m map[IK]IV) NativeHashRecord[OK, OV] {
	return NativeHashRecord[OK, OV](castNativeMap[IK, IV, OK, OV](m))
}

// Transform a map with native go types to a new Elk `NativeHashRecord` with corresponding Elk types
// using the given function.
// eg.
//
//	TransformIntoNativeHashRecord(m, func(k string, v uint8) (value.String, value.UInt8) {
//		return value.String(k), value.UInt8(v)
//	})
func TransformIntoNativeHashRecord[
	IK comparable,
	IV any,
	OK value.ComparableValueInterface,
	OV value.ValueInterface,
](
	m map[IK]IV,
	fn func(k IK, v IV) (OK, OV),
) NativeHashRecord[OK, OV] {
	newMap := MakeNativeHashRecord[OK, OV](len(m))
	for k, v := range m {
		ok, ov := fn(k, v)
		newMap[ok] = ov
	}
	return newMap
}

func MakeNativeHashRecordFromMap[K value.ComparableValueInterface, V value.ValueInterface](m map[K]V) NativeHashRecord[K, V] {
	return NativeHashRecord[K, V](m)
}

func MakeNativeHashRecord[K value.ComparableValueInterface, V value.ValueInterface](capacity int) NativeHashRecord[K, V] {
	return make(NativeHashRecord[K, V], capacity)
}

func (h NativeHashRecord[K, V]) CloneHashRecord(thread *Thread, capacity int) (HashRecord, value.Value) {
	newRecord := MakeNativeHashRecord[K, V](capacity)
	maps.Copy(newRecord, h)
	return newRecord, value.Undefined
}

func (h NativeHashRecord[K, V]) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		for k, v := range h {
			pair := value.MakePairOfValue(k.ToValue(), v.ToValue())
			if !yield(pair) {
				return
			}
		}
	}
}

func (h NativeHashRecord[K, V]) AllNative() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range h {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (h NativeHashRecord[K, V]) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		for k, v := range h {
			pair := value.NewNativePair(k, v)
			if !yield(pair.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (h NativeHashRecord[K, V]) IterRecord() value.NativeResettableIterator {
	return NewNativeHashRecordIterator(h)
}

func (h NativeHashRecord[K, V]) Iter() value.NativeIterator {
	return NewNativeHashRecordIterator(h)
}

func (h NativeHashRecord[K, V]) Get(key K) (V, bool) {
	v, ok := h[key]
	return v, ok
}

func (h NativeHashRecord[K, V]) GetVal(thread *Thread, key value.Value) (value.Value, value.Value) {
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

func (h NativeHashRecord[K, V]) SetVal(thread *Thread, key, val value.Value) value.Value {
	k, ok := value.Downcast[K](key)
	if !ok {
		return value.NewInvalidKeyInTypedMap(h, k.Class()).ToValue()
	}
	v, ok := value.Downcast[V](val)
	if !ok {
		return value.NewInvalidValueInTypedMap(h, v.Class()).ToValue()
	}

	h[k] = v
	return value.Undefined
}

func (h NativeHashRecord[K, V]) ConcatVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case NativeHashRecord[K, V]:
		newMap := MakeNativeHashRecord[K, V](h.Length() + o.Length())
		maps.Copy(newMap, h)
		maps.Copy(newMap, o)
		return newMap.ToValue(), value.Undefined
	case HashMap:
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
	case HashRecord:
		newMap := NewHashRecordOfValue(h.Length() + o.Length())

		err := HashRecordOfValueCopyInterface(thread, newMap, h)
		if err.IsNotUndefined() {
			return value.Undefined, err
		}

		err = HashRecordOfValueCopyInterface(thread, newMap, o)
		if err.IsNotUndefined() {
			return value.Undefined, err
		}

		return newMap.ToValue(), value.Undefined
	}

	return value.Undefined, value.Ref(value.Errorf(value.TypeErrorClass, "cannot concat %s with map %s", other.Inspect(), h.Inspect()))
}

func (h NativeHashRecord[K, V]) ContainsNativePair(thread *Thread, other *value.NativePair[K, V]) (bool, value.Value) {
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

func (h NativeHashRecord[K, V]) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
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

func (h NativeHashRecord[K, V]) ContainsNativeValue(thread *Thread, val V) (bool, value.Value) {
	for _, hval := range h {
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

func (h NativeHashRecord[K, V]) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	v, ok := value.Downcast[V](val)
	if !ok {
		return false, value.NewInvalidValueInTypedMap(h, val.Class()).ToValue()
	}
	return h.ContainsNativeValue(thread, v)
}

func (h NativeHashRecord[K, V]) ContainsNativeKey(thread *Thread, key K) bool {
	_, ok := h.Get(key)
	return ok
}

func (h NativeHashRecord[K, V]) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return false, value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}
	return h.ContainsNativeKey(thread, k), value.Undefined
}

func (h NativeHashRecord[K, V]) EqualNative(thread *Thread, other NativeHashRecord[K, V]) (bool, value.Value) {
	if h.Length() != other.Length() {
		return false, value.Undefined
	}

	for hkey, hval := range h {
		oval := other[hkey]
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

func (h NativeHashRecord[K, V]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case NativeHashRecord[K, V]:
		return h.EqualNative(thread, o)
	case HashRecord:
		return HashRecordEqual(thread, o, h)
	}

	return false, value.Undefined
}

func (h NativeHashRecord[K, V]) LaxEqualNative(thread *Thread, other NativeHashRecord[K, V]) (bool, value.Value) {
	if h.Length() != other.Length() {
		return false, value.Undefined
	}

	for hkey, hval := range h {
		oval := other[hkey]
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

func (h NativeHashRecord[K, V]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case NativeHashRecord[K, V]:
		return h.LaxEqualNative(thread, o)
	case HashRecord:
		return HashRecordLaxEqual(thread, h, o)
	}

	return false, value.Undefined
}

func (NativeHashRecord[K, V]) Class() *value.Class {
	return value.HashRecordClass
}

func (NativeHashRecord[K, V]) DirectClass() *value.Class {
	return value.HashRecordClass
}

func (NativeHashRecord[K, V]) SingletonClass() *value.Class {
	return nil
}

func (h NativeHashRecord[K, V]) Clone() NativeHashRecord[K, V] {
	newMap := MakeNativeHashRecord[K, V](h.Length())
	maps.Copy(newMap, h)
	return newMap
}

func (h NativeHashRecord[K, V]) Copy() value.Reference {
	return h.Clone()
}

func (h NativeHashRecord[K, V]) ToValue() value.Value {
	return value.Ref(h)
}

func (h NativeHashRecord[K, V]) Error() string {
	return h.Inspect()
}

func (h NativeHashRecord[K, V]) Inspect() string {
	var hasMultilineElements bool
	keyStrings := make(
		[]string,
		0,
		min(MAX_HASH_RECORD_ELEMENTS_IN_INSPECT, h.Length()),
	)
	valStrings := make(
		[]string,
		0,
		min(MAX_HASH_RECORD_ELEMENTS_IN_INSPECT, h.Length()),
	)

	i := 0
	for key, val := range h {
		keyString := key.Inspect()
		keyStrings = append(keyStrings, keyString)

		valString := val.Inspect()
		valStrings = append(valStrings, valString)

		if strings.ContainsRune(keyString, '\n') ||
			strings.ContainsRune(valString, '\n') {
			hasMultilineElements = true
		}

		if i >= MAX_HASH_RECORD_ELEMENTS_IN_INSPECT-1 {
			break
		}
		i++
	}

	var buff strings.Builder

	buff.WriteString("%{")
	if hasMultilineElements {
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

			if i >= MAX_HASH_RECORD_ELEMENTS_IN_INSPECT-1 {
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

			if i >= MAX_HASH_RECORD_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(", ...")
				break
			}
		}
	}
	buff.WriteRune('}')

	return buff.String()
}

func (NativeHashRecord[K, V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h NativeHashRecord[K, V]) Length() int {
	return len(h)
}

type NativeHashRecordIterator[K value.ComparableValueInterface, V value.ValueInterface] struct {
	HashRecord NativeHashRecord[K, V]
	index      int
	snapshot   []value.NativePair[K, V]
}

var _ value.NativeResettableIterator = &NativeHashRecordIterator[value.String, value.UInt8]{}

func NewNativeHashRecordIterator[K value.ComparableValueInterface, V value.ValueInterface](hrec NativeHashRecord[K, V]) *NativeHashRecordIterator[K, V] {
	iterator := &NativeHashRecordIterator[K, V]{
		HashRecord: hrec,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *NativeHashRecordIterator[K, V]) captureSnapshot() {
	snapshot := make([]value.NativePair[K, V], 0, h.HashRecord.Length())
	for k, v := range h.HashRecord.AllNative() {
		snapshot = append(snapshot, value.MakeNativePair(k, v))
	}
	h.snapshot = snapshot
}

func (*NativeHashRecordIterator[K, V]) Class() *value.Class {
	return value.HashRecordIteratorClass
}

func (*NativeHashRecordIterator[K, V]) DirectClass() *value.Class {
	return value.HashRecordIteratorClass
}

func (*NativeHashRecordIterator[K, V]) SingletonClass() *value.Class {
	return nil
}

func (h *NativeHashRecordIterator[K, V]) Error() string {
	return h.Inspect()
}

func (h *NativeHashRecordIterator[K, V]) Copy() value.Reference {
	return &NativeHashRecordIterator[K, V]{
		HashRecord: h.HashRecord,
		index:      h.index,
	}
}

func (i *NativeHashRecordIterator[K, V]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *NativeHashRecordIterator[K, V]) Inspect() string {
	return fmt.Sprintf("Std::HashRecord::Iterator{hash_record: %s}", h.HashRecord.Inspect())
}

func (*NativeHashRecordIterator[K, V]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *NativeHashRecordIterator[K, V]) Next() (p value.NativePair[K, V], err value.Value) {
	if h.index >= len(h.snapshot) {
		return p, symbol.L_stop_iteration.ToValue()
	}

	h.index++
	return h.snapshot[h.index], value.Undefined
}

func (h *NativeHashRecordIterator[K, V]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *NativeHashRecordIterator[K, V]) Reset() {
	h.index = 0
}
