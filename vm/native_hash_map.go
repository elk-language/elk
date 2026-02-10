package vm

import (
	"fmt"
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type NativeHashMap[K value.ComparableValueInterface, V value.ValueInterface] struct {
	m       map[K]V
	version int
}

var _ HashMap = &NativeHashMap[value.String, value.String]{}

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

func (h *NativeHashMap[K, V]) All() iter.Seq[value.NativePair[K, V]] {
	return func(yield func(value.NativePair[K, V]) bool) {
		for k, v := range h.m {
			pair := value.MakeNativePair(k, v)
			if !yield(pair) {
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

func (h *NativeHashMap[K, V]) IterMap() HashMapIterator {
	return NewNativeHashMapIterator(h)
}

func (h *NativeHashMap[K, V]) IterRecord() HashRecordIterator {
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
		return value.NewInvalidKeyInTypedMap(h, k.Class()).ToValue()
	}
	v, ok := value.Downcast[V](val)
	if !ok {
		return value.NewInvalidValueInTypedMap(h, v.Class()).ToValue()
	}

	h.m[k] = v
	return value.Undefined
}

// func (h *HashMapOfValue) ConcatVal(thread *Thread, other value.Value) (value.Value, value.Value) {
// 	switch o := other.SafeAsReference().(type) {
// 	case *HashMapOfValue:
// 		return value.RefErr(HashMapOfValueConcat(thread, h, o))
// 	case HashRecord:
// 		return value.RefErr(HashMapOfValueConcatInterface(thread, h, o))
// 	}

// 	return value.Undefined, value.Ref(value.Errorf(value.TypeErrorClass, "cannot concat %s with map %s", other.Inspect(), h.Inspect()))
// }

// func (h *HashMapOfValue) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
// 	return HashMapOfValueContains(thread, h, other)
// }

// func (h *HashMapOfValue) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
// 	return HashMapOfValueContainsValue(thread, h, val)
// }

// func (h *HashMapOfValue) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
// 	return HashMapOfValueContainsKey(thread, h, key)
// }

// func (h *HashMapOfValue) Equal(thread *Thread, other value.Value) (bool, value.Value) {
// 	switch o := other.SafeAsReference().(type) {
// 	case *HashMapOfValue:
// 		return HashMapOfValueEqual(thread, h, o)
// 	case HashMap:
// 		return HashMapOfValueEqualInterface(thread, h, o)
// 	}

// 	return false, value.Undefined
// }

// func (h *HashMapOfValue) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
// 	switch o := other.SafeAsReference().(type) {
// 	case *HashMapOfValue:
// 		return HashMapOfValueLaxEqual(thread, h, o)
// 	case HashMap:
// 		return HashMapOfValueLaxEqualInterface(thread, h, o)
// 	}

// 	return false, value.Undefined
// }

// func (h *HashMapOfValue) Grow(thread *Thread, newSlots int) value.Value {
// 	return HashMapOfValueGrow(thread, h, newSlots)
// }

// func (*HashMapOfValue) Class() *value.Class {
// 	return value.HashMapClass
// }

// func (*HashMapOfValue) DirectClass() *value.Class {
// 	return value.HashMapClass
// }

// func (*HashMapOfValue) SingletonClass() *value.Class {
// 	return nil
// }

// func (h *HashMapOfValue) Clone() *HashMapOfValue {
// 	newTable := slices.Clone(h.Table)
// 	return &HashMapOfValue{
// 		Table:         newTable,
// 		OccupiedSlots: h.OccupiedSlots,
// 		Elements:      h.Elements,
// 	}
// }

// func (h *HashMapOfValue) Copy() value.Reference {
// 	return h.Clone()
// }

// func (h *HashMapOfValue) ToValue() value.Value {
// 	return value.Ref(h)
// }

// func (h *HashMapOfValue) Error() string {
// 	return h.Inspect()
// }

// const MAX_HASH_MAP_ELEMENTS_IN_INSPECT = 300

// func (h *HashMapOfValue) Inspect() string {
// 	var hasMultilineElements bool
// 	keyStrings := make(
// 		[]string,
// 		0,
// 		min(MAX_HASH_MAP_ELEMENTS_IN_INSPECT, h.Length()),
// 	)
// 	valStrings := make(
// 		[]string,
// 		0,
// 		min(MAX_HASH_MAP_ELEMENTS_IN_INSPECT, h.Length()),
// 	)

// 	i := 0
// 	for _, entry := range h.Table {
// 		if entry.Key().IsUndefined() {
// 			continue
// 		}

// 		keyString := entry.Key().Inspect()
// 		keyStrings = append(keyStrings, keyString)

// 		valString := entry.Value().Inspect()
// 		valStrings = append(valStrings, valString)

// 		if strings.ContainsRune(keyString, '\n') ||
// 			strings.ContainsRune(valString, '\n') {
// 			hasMultilineElements = true
// 		}

// 		if i >= MAX_HASH_MAP_ELEMENTS_IN_INSPECT-1 {
// 			break
// 		}
// 		i++
// 	}

// 	var buff strings.Builder

// 	buff.WriteRune('{')
// 	if hasMultilineElements || h.Length() > 15 {
// 		buff.WriteRune('\n')
// 		for i := range len(keyStrings) {
// 			keyString := keyStrings[i]
// 			valString := valStrings[i]

// 			if i != 0 {
// 				buff.WriteString(",\n")
// 			}
// 			indent.IndentString(&buff, keyString, 1)
// 			buff.WriteString(" => ")
// 			indent.IndentStringFromSecondLine(&buff, valString, 1)

// 			if i >= MAX_HASH_MAP_ELEMENTS_IN_INSPECT-1 {
// 				buff.WriteString(",\n  ...")
// 				break
// 			}
// 		}
// 		buff.WriteRune('\n')
// 	} else {
// 		for i := range len(keyStrings) {
// 			keyString := keyStrings[i]
// 			valString := valStrings[i]

// 			if i != 0 {
// 				buff.WriteString(", ")
// 			}
// 			buff.WriteString(keyString)
// 			buff.WriteString(" => ")
// 			buff.WriteString(valString)

// 			if i >= MAX_HASH_MAP_ELEMENTS_IN_INSPECT-1 {
// 				buff.WriteString(", ...")
// 				break
// 			}
// 		}
// 	}
// 	buff.WriteRune('}')

// 	leftCapacity := h.LeftCapacity()
// 	if leftCapacity > 0 {
// 		fmt.Fprintf(&buff, ":%d", leftCapacity)
// 	}
// 	return buff.String()
// }

// func (*HashMapOfValue) InstanceVariables() *value.InstanceVariables {
// 	return nil
// }

// func (h *HashMapOfValue) Capacity() int {
// 	return len(h.Table)
// }

// func (h *HashMapOfValue) LeftCapacity() int {
// 	return h.Capacity() - h.Length()
// }

// func (h *HashMapOfValue) Length() int {
// 	return h.Elements
// }

type NativeHashMapIterator[K value.ComparableValueInterface, V value.ValueInterface] struct {
	HashMap *NativeHashMap[K, V]
	Index   int
}

var _ HashMapIterator = &NativeHashMapIterator[value.String, value.String]{}

func NewNativeHashMapIterator[K value.ComparableValueInterface, V value.ValueInterface](hmap *NativeHashMap[K, V]) *NativeHashMapIterator[K, V] {
	return &NativeHashMapIterator[K, V]{
		HashMap: hmap,
	}
}

func NewNativeHashMapIteratorWithIndex[K value.ComparableValueInterface, V value.ValueInterface](hmap *NativeHashMap[K, V], index int) *NativeHashMapIterator[K, V] {
	return &NativeHashMapIterator[K, V]{
		HashMap: hmap,
		Index:   index,
	}
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
		HashMap: h.HashMap,
		Index:   h.Index,
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
	for {
		if h.Index >= h.HashMap.Capacity() {
			return p, symbol.L_stop_iteration.ToValue()
		}

		pair := h.HashMap.Table[h.Index]
		h.Index++
		if !pair.Key().IsUndefined() {
			return value.Ref(&h.HashMap.Table[h.Index-1]), value.Undefined
		}
	}
}

func (h *NativeHashMapIterator[K, V]) NextValue() (value.Value, value.Value) {
	for {
		if h.Index >= h.HashMap.Capacity() {
			return value.Undefined, symbol.L_stop_iteration.ToValue()
		}

		pair := h.HashMap.Table[h.Index]
		h.Index++
		if !pair.Key().IsUndefined() {
			return value.Ref(&h.HashMap.Table[h.Index-1]), value.Undefined
		}
	}
}

func (h *NativeHashMapIterator[K, V]) Reset() {
	h.Index = 0
}
