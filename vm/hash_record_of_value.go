package vm

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type HashRecordOfValue HashMapOfValue

var _ HashRecord = &HashRecordOfValue{}

func NewHashRecordOfValue(capacity int) *HashRecordOfValue {
	return &HashRecordOfValue{
		Table: make([]value.PairOfValue, capacity),
	}
}

func HashRecordOfValueConstructor(class *value.Class) value.Value {
	return value.Ref(&HashRecordOfValue{})
}

func (h *HashRecordOfValue) Keys() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		for _, pair := range h.Table {
			if pair.Key().IsUndefined() {
				continue
			}

			if !yield(pair.Key()) {
				return
			}
		}
	}
}

func (h *HashRecordOfValue) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		for _, pair := range h.Table {
			if pair.Key().IsUndefined() {
				continue
			}

			if !yield(pair) {
				return
			}
		}
	}
}

func (h *HashRecordOfValue) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		for _, pair := range h.Table {
			if pair.Key().IsUndefined() {
				continue
			}

			if !yield(pair.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (h *HashRecordOfValue) Clone() *HashRecordOfValue {
	return &HashRecordOfValue{
		Table:         slices.Clone(h.Table),
		OccupiedSlots: h.OccupiedSlots,
		Elements:      h.Elements,
	}
}

func (h *HashRecordOfValue) CloneHashRecord(thread *Thread, capacity int) (HashRecord, value.Value) {
	newRecord := NewHashRecordOfValue(capacity)
	err := HashRecordOfValueCopy(thread, newRecord, h)
	if !err.IsUndefined() {
		return nil, err
	}

	return newRecord, value.Undefined
}

func (h *HashRecordOfValue) NewHashRecord(capacity int) HashRecord {
	return NewHashRecordOfValue(capacity)
}

func (h *HashRecordOfValue) IterNative() *HashRecordOfValueIterator {
	return NewHashRecordOfValueIterator(h)
}

func (h *HashRecordOfValue) Iter() value.NativeIterator {
	return h.IterNative()
}

func (h *HashRecordOfValue) IterRecord() value.NativeResettableIterator {
	return h.IterNative()
}

func (h *HashRecordOfValue) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
	return HashRecordOfValueGet(thread, h, key)
}

func (h *HashRecordOfValue) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
	val, err := h.GetValUndefined(thread, key)
	if err.IsNotUndefined() {
		return value.Undefined, err
	}
	if val.IsUndefined() {
		return value.Nil, value.Undefined
	}
	return val, value.Undefined
}

func (h *HashRecordOfValue) SetVal(thread *Thread, key, val value.Value) value.Value {
	return HashRecordOfValueSet(thread, h, key, val)
}

func (h *HashRecordOfValue) Merge(thread *Thread, other value.Value) (*HashRecordOfValue, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashRecordOfValue:
		return HashRecordOfValueMerge(thread, h, o)
	case HashRecord:
		return HashRecordOfValueMergeInterface(thread, h, o)
	case Record:
		newMap := NewHashRecordOfValue(h.Length() + o.Length())
		err := HashRecordOfValueCopyInterface(thread, newMap, h)
		if err.IsNotUndefined() {
			return nil, err
		}

		err = HashRecordOfValueCopyInterface(thread, newMap, o)
		if err.IsNotUndefined() {
			return nil, err
		}

		return newMap, value.Undefined
	default:
		result := NewHashRecordOfValue(h.Length())

		err := HashRecordOfValueCopyInterface(thread, result, h)
		if err.IsNotUndefined() {
			return nil, err
		}

		for element, err := range Iterate(thread, other) {
			if err.IsNotUndefined() {
				return nil, err
			}

			entry := element.AsReference().(value.Pair)
			result.SetVal(thread, entry.Key(), entry.Value())
		}

		return result, value.Undefined
	}

}

func (h *HashRecordOfValue) MergeVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	return value.RefErr(h.Merge(thread, other))
}

func (h *HashRecordOfValue) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
	return HashRecordOfValueContains(thread, h, other)
}

func (h *HashRecordOfValue) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	return HashRecordOfValueContainsValue(thread, h, val)
}

func (h *HashRecordOfValue) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	return HashRecordOfValueContainsKey(thread, h, key)
}

func (h *HashRecordOfValue) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashRecordOfValue:
		return HashRecordOfValueEqual(thread, h, o)
	case HashRecord:
		return HashRecordOfValueEqualInterface(thread, h, o)
	}

	return false, value.Undefined
}

func (h *HashRecordOfValue) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashRecordOfValue:
		return HashRecordOfValueLaxEqual(thread, h, o)
	case HashRecord:
		return HashRecordOfValueLaxEqualInterface(thread, h, o)
	default:
		if !value.IsA(other, value.MapMixin) {
			return false, value.Undefined
		}

		for element, err := range Iterate(thread, other) {
			if err.IsNotUndefined() {
				return false, err
			}

			pair := element.AsReference().(value.Pair)
			val, err := h.GetValUndefined(thread, pair.Key())
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

func (*HashRecordOfValue) Class() *value.Class {
	return value.HashRecordClass
}

func (*HashRecordOfValue) DirectClass() *value.Class {
	return value.HashRecordClass
}

func (*HashRecordOfValue) SingletonClass() *value.Class {
	return nil
}

func (h *HashRecordOfValue) Copy() value.Reference {
	return h
}

func (h *HashRecordOfValue) ToValue() value.Value {
	return value.Ref(h)
}

func (h *HashRecordOfValue) Error() string {
	return h.Inspect()
}

const MAX_HASH_RECORD_ELEMENTS_IN_INSPECT = 300

func (h *HashRecordOfValue) Inspect() string {
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
	for _, entry := range h.Table {
		if entry.Key().IsUndefined() {
			continue
		}

		keyString := entry.Key().Inspect()
		keyStrings = append(keyStrings, keyString)

		valString := entry.Value().Inspect()
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

func (*HashRecordOfValue) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *HashRecordOfValue) Length() int {
	return h.Elements
}

type HashRecordOfValueIterator struct {
	HashRecord *HashRecordOfValue
	Index      int
	snapshot   []value.PairOfValue
}

var _ value.NativeResettableIterator = &HashRecordOfValueIterator{}

func NewHashRecordOfValueIterator(hrec *HashRecordOfValue) *HashRecordOfValueIterator {
	iterator := &HashRecordOfValueIterator{
		HashRecord: hrec,
	}
	iterator.captureSnapshot()
	return iterator
}

func NewHashRecordOfValueIteratorWithIndex(hrec *HashRecordOfValue, index int) *HashRecordOfValueIterator {
	iterator := &HashRecordOfValueIterator{
		HashRecord: hrec,
		Index:      index,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *HashRecordOfValueIterator) captureSnapshot() {
	snapshot := make([]value.PairOfValue, 0, h.HashRecord.Length())
	for pair := range h.HashRecord.All() {
		snapshot = append(snapshot, pair)
	}
	h.snapshot = snapshot
}

func (*HashRecordOfValueIterator) Class() *value.Class {
	return value.HashRecordIteratorClass
}

func (*HashRecordOfValueIterator) DirectClass() *value.Class {
	return value.HashRecordIteratorClass
}

func (*HashRecordOfValueIterator) SingletonClass() *value.Class {
	return nil
}

func (h *HashRecordOfValueIterator) Error() string {
	return h.Inspect()
}

func (h *HashRecordOfValueIterator) Copy() value.Reference {
	return &HashRecordOfValueIterator{
		HashRecord: h.HashRecord,
		Index:      h.Index,
		snapshot:   h.snapshot,
	}
}

func (i *HashRecordOfValueIterator) ToValue() value.Value {
	return value.Ref(i)
}

func (h *HashRecordOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::HashRecord::Iterator{hash_record: %s}", h.HashRecord.Inspect())
}

func (*HashRecordOfValueIterator) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *HashRecordOfValueIterator) NextValue() (value.Value, value.Value) {
	if h.Index >= len(h.snapshot) {
		return value.Undefined, value.S(symbol.L_stop_iteration).ToValue()
	}

	pair := &h.snapshot[h.Index]
	h.Index++
	return value.Ref(pair), value.Undefined
}

func (h *HashRecordOfValueIterator) Reset() {
	h.Index = 0
}

// Create a new hash record with the given entries.
func NewHashRecordOfValueWithElements(vm *Thread, elements ...value.PairOfValue) (*HashRecordOfValue, value.Value) {
	return NewHashRecordOfValueWithCapacityAndElements(vm, len(elements), elements...)
}

// Create a new hash record with the given entries.
func MustNewHashRecordOfValueWithElements(vm *Thread, elements ...value.PairOfValue) *HashRecordOfValue {
	hrec, err := NewHashRecordOfValueWithElements(vm, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return hrec
}

func NewHashRecordOfValueWithCapacityAndElements(vm *Thread, capacity int, elements ...value.PairOfValue) (*HashRecordOfValue, value.Value) {
	h := NewHashRecordOfValue(capacity)
	for _, element := range elements {
		err := HashRecordOfValueSet(vm, h, element.Key(), element.Value())
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return h, value.Undefined
}

func MustNewHashRecordWithCapacityAndElements(vm *Thread, capacity int, elements ...value.PairOfValue) *HashRecordOfValue {
	hrec, err := NewHashRecordOfValueWithCapacityAndElements(vm, capacity, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return hrec
}

// Delete the given key from the hashMap
func HashRecordOfValueDelete(vm *Thread, hashRecord *HashRecordOfValue, key value.Value) (bool, value.Value) {
	return HashMapOfValueDelete(vm, (*HashMapOfValue)(hashRecord), key)
}

// Get the element under the given key.
func HashRecordOfValueGet(vm *Thread, hashRecord *HashRecordOfValue, key value.Value) (value.Value, value.Value) {
	return HashMapOfValueGet(vm, (*HashMapOfValue)(hashRecord), key)
}

func HashRecordOfValueCopyTable(vm *Thread, target *HashRecordOfValue, source []value.PairOfValue) value.Value {
	return HashMapOfValueCopyTable(vm, (*HashMapOfValue)(target), source)
}

// Copy the pairs of one hash record to the other.
func HashRecordOfValueCopy(vm *Thread, target *HashRecordOfValue, source *HashRecordOfValue) value.Value {
	return HashMapOfValueCopy(vm, (*HashMapOfValue)(target), (*HashMapOfValue)(source))
}

func HashRecordOfValueCopyInterface(vm *Thread, target *HashRecordOfValue, source Record) value.Value {
	return HashMapOfValueCopyInterface(vm, (*HashMapOfValue)(target), source)
}

// Create a new map containing the pairs of both maps.
func HashRecordOfValueMerge(vm *Thread, x *HashRecordOfValue, y *HashRecordOfValue) (*HashRecordOfValue, value.Value) {
	m, err := HashMapOfValueMerge(vm, (*HashMapOfValue)(x), (*HashMapOfValue)(y))
	if err.IsNotUndefined() {
		return nil, err
	}
	return (*HashRecordOfValue)(m), value.Undefined
}

func HashRecordOfValueMergeInterface(vm *Thread, x *HashRecordOfValue, y HashRecord) (*HashRecordOfValue, value.Value) {
	m, err := HashMapOfValueMergeInterface(vm, (*HashMapOfValue)(x), y)
	if err.IsNotUndefined() {
		return nil, err
	}
	return (*HashRecordOfValue)(m), value.Undefined
}

// Check if the given pair is present in the record
func HashRecordOfValueContains(vm *Thread, hrec *HashRecordOfValue, pair value.Pair) (bool, value.Value) {
	return HashMapOfValueContains(vm, (*HashMapOfValue)(hrec), pair)
}

// Check if the given key is present in the record
func HashRecordOfValueContainsKey(vm *Thread, hrec *HashRecordOfValue, key value.Value) (bool, value.Value) {
	return HashMapOfValueContainsKey(vm, (*HashMapOfValue)(hrec), key)
}

// Check if the given value is present in the record
func HashRecordOfValueContainsValue(vm *Thread, hrec *HashRecordOfValue, val value.Value) (bool, value.Value) {
	return HashMapOfValueContainsValue(vm, (*HashMapOfValue)(hrec), val)
}

// Checks whether two hash records are equal (lax)
func HashRecordOfValueLaxEqual(vm *Thread, x *HashRecordOfValue, y *HashRecordOfValue) (bool, value.Value) {
	return HashMapOfValueLaxEqual(vm, (*HashMapOfValue)(x), (*HashMapOfValue)(y))
}

func HashRecordOfValueLaxEqualInterface(vm *Thread, x *HashRecordOfValue, y HashRecord) (bool, value.Value) {
	return HashMapOfValueLaxEqualInterface(vm, (*HashMapOfValue)(x), y)
}

// Checks whether two hash records are equal
func HashRecordOfValueEqual(vm *Thread, x *HashRecordOfValue, y *HashRecordOfValue) (bool, value.Value) {
	return HashMapOfValueEqual(vm, (*HashMapOfValue)(x), (*HashMapOfValue)(y))
}

func HashRecordOfValueEqualInterface(vm *Thread, x *HashRecordOfValue, y HashRecord) (bool, value.Value) {
	return HashMapOfValueEqualInterface(vm, (*HashMapOfValue)(x), y)
}

// Add additional n empty slots for new elements.
func HashRecordOfValueGrow(vm *Thread, hashRecord *HashRecordOfValue, newSlots int) value.Value {
	return HashMapOfValueGrow(vm, (*HashMapOfValue)(hashRecord), newSlots)
}

// Resize the given hash record to the desired capacity.
func HashRecordOfValueSetCapacity(vm *Thread, hashRecord *HashRecordOfValue, capacity int) value.Value {
	return HashMapOfValueSetCapacity(vm, (*HashMapOfValue)(hashRecord), capacity)
}

func HashRecordOfValueSetWithMaxLoad(vm *Thread, hashRecord *HashRecordOfValue, key, val value.Value, maxLoad float64) value.Value {
	return HashMapOfValueSetWithMaxLoad(vm, (*HashMapOfValue)(hashRecord), key, val, maxLoad)
}

// Set a value under the given key.
func HashRecordOfValueSet(vm *Thread, hashRecord *HashRecordOfValue, key, val value.Value) value.Value {
	return HashMapOfValueSet(vm, (*HashMapOfValue)(hashRecord), key, val)
}
