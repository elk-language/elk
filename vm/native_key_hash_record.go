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

type NativeKeyHashRecord[K value.ComparableValueInterface] map[K]value.Value

var _ HashRecord = NativeKeyHashRecord[value.String]{}

// Transform a map with native go types to a new Elk `NativeKeyHashRecord` with corresponding Elk types
// using the given function.
// eg.
//
//	TransformIntoNativeKeyHashRecord(m, func(k string, v uint8) (value.String, value.Value) {
//		return value.String(k), value.UInt8(v).ToValue()
//	})
func TransformIntoNativeKeyHashRecord[
	IK comparable,
	IV any,
	OK value.ComparableValueInterface,
](
	m map[IK]IV,
	fn func(k IK, v IV) (OK, value.Value),
) NativeHashRecord[OK, value.Value] {
	newMap := MakeNativeHashRecord[OK, value.Value](len(m))
	for k, v := range m {
		ok, ov := fn(k, v)
		newMap[ok] = ov
	}
	return newMap
}

func MakeNativeKeyHashRecordFromMap[K value.ComparableValueInterface](m map[K]value.Value) NativeKeyHashRecord[K] {
	return NativeKeyHashRecord[K](m)
}

func MakeNativeKeyHashRecord[K value.ComparableValueInterface](capacity int) NativeKeyHashRecord[K] {
	return make(NativeKeyHashRecord[K], capacity)
}

func (h NativeKeyHashRecord[K]) CloneHashRecord(thread *Thread, capacity int) (HashRecord, value.Value) {
	newMap := MakeNativeKeyHashRecord[K](capacity)
	maps.Copy(newMap, h)
	return newMap, value.Undefined
}

func (h NativeKeyHashRecord[K]) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		for k, v := range h {
			pair := value.MakePairOfValue(k.ToValue(), v.ToValue())
			if !yield(pair) {
				return
			}
		}
	}
}

func (h NativeKeyHashRecord[K]) AllNative() iter.Seq2[K, value.Value] {
	return func(yield func(K, value.Value) bool) {
		for k, v := range h {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (h NativeKeyHashRecord[K]) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		for k, v := range h {
			pair := value.NewNativePair(k, v)
			if !yield(pair.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (h NativeKeyHashRecord[K]) IterNative() *NativeKeyHashRecordIterator[K] {
	return NewNativeKeyHashRecordIterator(h)
}

func (h NativeKeyHashRecord[K]) IterRecord() value.NativeResettableIterator {
	return h.IterNative()
}

func (h NativeKeyHashRecord[K]) Iter() value.NativeIterator {
	return h.IterNative()
}

func (h NativeKeyHashRecord[K]) Get(key K) (value.Value, bool) {
	v, ok := h[key]
	return v, ok
}

func (h NativeKeyHashRecord[K]) GetVal(thread *Thread, key value.Value) (value.Value, value.Value) {
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

func (h NativeKeyHashRecord[K]) SetVal(thread *Thread, key, val value.Value) value.Value {
	k, ok := value.Downcast[K](key)
	if !ok {
		return value.NewInvalidKeyInTypedMap(h, k.Class()).ToValue()
	}
	h[k] = val
	return value.Undefined
}

func (h NativeKeyHashRecord[K]) ConcatVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case NativeKeyHashRecord[K]:
		newMap := MakeNativeKeyHashRecord[K](h.Length() + o.Length())
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

func (h NativeKeyHashRecord[K]) ContainsNativePair(thread *Thread, other *value.NativePair[K, value.Value]) (bool, value.Value) {
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

func (h NativeKeyHashRecord[K]) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
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

func (h NativeKeyHashRecord[K]) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
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

func (h NativeKeyHashRecord[K]) ContainsNativeKey(thread *Thread, key K) bool {
	_, ok := h.Get(key)
	return ok
}

func (h NativeKeyHashRecord[K]) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	k, ok := value.Downcast[K](key)
	if !ok {
		return false, value.NewInvalidKeyInTypedMap(h, key.Class()).ToValue()
	}
	return h.ContainsNativeKey(thread, k), value.Undefined
}

func (h NativeKeyHashRecord[K]) EqualNative(thread *Thread, other NativeKeyHashRecord[K]) (bool, value.Value) {
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

func (h NativeKeyHashRecord[K]) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case NativeKeyHashRecord[K]:
		return h.EqualNative(thread, o)
	case HashRecord:
		return HashRecordEqual(thread, o, h)
	}

	return false, value.Undefined
}

func (h NativeKeyHashRecord[K]) LaxEqualNative(thread *Thread, other NativeKeyHashRecord[K]) (bool, value.Value) {
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

func (h NativeKeyHashRecord[K]) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case NativeKeyHashRecord[K]:
		return h.LaxEqualNative(thread, o)
	case HashRecord:
		return HashRecordLaxEqual(thread, h, o)
	}

	return false, value.Undefined
}

func (NativeKeyHashRecord[K]) Class() *value.Class {
	return value.HashRecordClass
}

func (NativeKeyHashRecord[K]) DirectClass() *value.Class {
	return value.HashRecordClass
}

func (NativeKeyHashRecord[K]) SingletonClass() *value.Class {
	return nil
}

func (h NativeKeyHashRecord[K]) Clone() NativeKeyHashRecord[K] {
	newMap := MakeNativeKeyHashRecord[K](h.Length())
	maps.Copy(newMap, h)
	return newMap
}

func (h NativeKeyHashRecord[K]) Copy() value.Reference {
	return h.Clone()
}

func (h NativeKeyHashRecord[K]) ToValue() value.Value {
	return value.Ref(h)
}

func (h NativeKeyHashRecord[K]) Error() string {
	return h.Inspect()
}

func (h NativeKeyHashRecord[K]) Inspect() string {
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

func (NativeKeyHashRecord[K]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h NativeKeyHashRecord[K]) Length() int {
	return len(h)
}

type NativeKeyHashRecordIterator[K value.ComparableValueInterface] struct {
	HashRecord NativeKeyHashRecord[K]
	index      int
	snapshot   []value.NativePair[K, value.Value]
}

var _ value.NativeResettableIterator = &NativeKeyHashRecordIterator[value.String]{}

func NewNativeKeyHashRecordIterator[K value.ComparableValueInterface](hrec NativeKeyHashRecord[K]) *NativeKeyHashRecordIterator[K] {
	iterator := &NativeKeyHashRecordIterator[K]{
		HashRecord: hrec,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *NativeKeyHashRecordIterator[K]) captureSnapshot() {
	snapshot := make([]value.NativePair[K, value.Value], 0, h.HashRecord.Length())
	for k, v := range h.HashRecord.AllNative() {
		snapshot = append(snapshot, value.MakeNativePair(k, v))
	}
	h.snapshot = snapshot
}

func (*NativeKeyHashRecordIterator[K]) Class() *value.Class {
	return value.HashRecordIteratorClass
}

func (*NativeKeyHashRecordIterator[K]) DirectClass() *value.Class {
	return value.HashRecordIteratorClass
}

func (*NativeKeyHashRecordIterator[K]) SingletonClass() *value.Class {
	return nil
}

func (h *NativeKeyHashRecordIterator[K]) Error() string {
	return h.Inspect()
}

func (h *NativeKeyHashRecordIterator[K]) Copy() value.Reference {
	return &NativeKeyHashRecordIterator[K]{
		HashRecord: h.HashRecord,
		index:      h.index,
	}
}

func (i *NativeKeyHashRecordIterator[K]) ToValue() value.Value {
	return value.Ref(i)
}

func (h *NativeKeyHashRecordIterator[K]) Inspect() string {
	return fmt.Sprintf("Std::HashRecord::Iterator{hash_record: %s}", h.HashRecord.Inspect())
}

func (*NativeKeyHashRecordIterator[K]) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *NativeKeyHashRecordIterator[K]) Next() (p value.NativePair[K, value.Value], err value.Value) {
	if h.index >= len(h.snapshot) {
		return p, symbol.L_stop_iteration.ToValue()
	}

	h.index++
	return h.snapshot[h.index], value.Undefined
}

func (h *NativeKeyHashRecordIterator[K]) NextValue() (value.Value, value.Value) {
	p, err := h.Next()
	if err.IsNotUndefined() {
		return value.Undefined, err
	}

	return p.ToValue(), value.Undefined
}

func (h *NativeKeyHashRecordIterator[K]) Reset() {
	h.index = 0
}
