package vm

import (
	"fmt"
	"iter"
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

func (h *HashRecordOfValue) IterRecord() HashRecordIterator {
	return NewHashRecordOfValueIterator(h)
}

func (h *HashRecordOfValue) GetVal(thread *Thread, key value.Value) (value.Value, value.Value) {
	return HashRecordOfValueGet(thread, h, key)
}

func (h *HashRecordOfValue) SetVal(thread *Thread, key, val value.Value) value.Value {
	return HashRecordOfValueSet(thread, h, key, val)
}

func (h *HashRecordOfValue) ConcatVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashRecordOfValue:
		return value.RefErr(HashRecordOfValueConcat(thread, h, o))
	case HashRecord:
		return value.RefErr(HashRecordOfValueConcatInterface(thread, h, o))
	}

	return value.Undefined, value.Ref(value.Errorf(value.TypeErrorClass, "cannot concat %s with record %s", other.Inspect(), h.Inspect()))
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
	}

	return false, value.Undefined
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
}

var _ HashRecordIterator = &HashRecordOfValueIterator{}

func NewHashRecordOfValueIterator(hrec *HashRecordOfValue) *HashRecordOfValueIterator {
	return &HashRecordOfValueIterator{
		HashRecord: hrec,
	}
}

func NewHashRecordOfValueIteratorWithIndex(hrec *HashRecordOfValue, index int) *HashRecordOfValueIterator {
	return &HashRecordOfValueIterator{
		HashRecord: hrec,
		Index:      index,
	}
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
	for {
		if h.Index >= len(h.HashRecord.Table) {
			return value.Undefined, symbol.L_stop_iteration.ToValue()
		}

		pair := h.HashRecord.Table[h.Index]
		h.Index++
		if !pair.Key().IsUndefined() {
			return value.Ref(&h.HashRecord.Table[h.Index-1]), value.Undefined
		}
	}
}

func (h *HashRecordOfValueIterator) Reset() {
	h.Index = 0
}
