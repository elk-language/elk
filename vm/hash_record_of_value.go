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

type HashRecordIteratorOfValue struct {
	HashRecord *HashRecordOfValue
	Index      int
}

func NewHashRecordIteratorOfValue(hrec *HashRecordOfValue) *HashRecordIteratorOfValue {
	return &HashRecordIteratorOfValue{
		HashRecord: hrec,
	}
}

func NewHashRecordIteratorOfValueWithIndex(hrec *HashRecordOfValue, index int) *HashRecordIteratorOfValue {
	return &HashRecordIteratorOfValue{
		HashRecord: hrec,
		Index:      index,
	}
}

func (*HashRecordIteratorOfValue) Class() *value.Class {
	return value.HashRecordIteratorClass
}

func (*HashRecordIteratorOfValue) DirectClass() *value.Class {
	return value.HashRecordIteratorClass
}

func (*HashRecordIteratorOfValue) SingletonClass() *value.Class {
	return nil
}

func (h *HashRecordIteratorOfValue) Error() string {
	return h.Inspect()
}

func (h *HashRecordIteratorOfValue) Copy() value.Reference {
	return &HashRecordIteratorOfValue{
		HashRecord: h.HashRecord,
		Index:      h.Index,
	}
}

func (i *HashRecordIteratorOfValue) ToValue() value.Value {
	return value.Ref(i)
}

func (h *HashRecordIteratorOfValue) Inspect() string {
	return fmt.Sprintf("Std::HashRecord::Iterator{hash_record: %s}", h.HashRecord.Inspect())
}

func (*HashRecordIteratorOfValue) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *HashRecordIteratorOfValue) NextValue() (value.Value, value.Value) {
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

func (h *HashRecordIteratorOfValue) Reset() {
	h.Index = 0
}
