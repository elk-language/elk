package value

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
)

var HashRecordClass *Class         // ::Std::HashRecord
var HashRecordIteratorClass *Class // ::Std::HashRecord::Iterator

type HashRecord HashMap

func NewHashRecord(capacity int) *HashRecord {
	return &HashRecord{
		Table: make([]Pair, capacity),
	}
}

func HashRecordConstructor(class *Class) Value {
	return Ref(&HashRecord{})
}

func (*HashRecord) Class() *Class {
	return HashRecordClass
}

func (*HashRecord) DirectClass() *Class {
	return HashRecordClass
}

func (*HashRecord) SingletonClass() *Class {
	return nil
}

func (h *HashRecord) Copy() Reference {
	return h
}

func (h *HashRecord) Error() string {
	return h.Inspect()
}

const MAX_HASH_RECORD_ELEMENTS_IN_INSPECT = 300

func (h *HashRecord) Inspect() string {
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
		if entry.Key.IsUndefined() {
			continue
		}

		keyString := entry.Key.Inspect()
		keyStrings = append(keyStrings, keyString)

		valString := entry.Value.Inspect()
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

func (*HashRecord) InstanceVariables() SymbolMap {
	return nil
}

func (h *HashRecord) Length() int {
	return h.Elements
}

type HashRecordIterator struct {
	HashRecord *HashRecord
	Index      int
}

func NewHashRecordIterator(hrec *HashRecord) *HashRecordIterator {
	return &HashRecordIterator{
		HashRecord: hrec,
	}
}

func NewHashRecordIteratorWithIndex(hrec *HashRecord, index int) *HashRecordIterator {
	return &HashRecordIterator{
		HashRecord: hrec,
		Index:      index,
	}
}

func (*HashRecordIterator) Class() *Class {
	return HashRecordIteratorClass
}

func (*HashRecordIterator) DirectClass() *Class {
	return HashRecordIteratorClass
}

func (*HashRecordIterator) SingletonClass() *Class {
	return nil
}

func (h *HashRecordIterator) Error() string {
	return h.Inspect()
}

func (h *HashRecordIterator) Copy() Reference {
	return &HashRecordIterator{
		HashRecord: h.HashRecord,
		Index:      h.Index,
	}
}

func (h *HashRecordIterator) Inspect() string {
	return fmt.Sprintf("Std::HashRecord::Iterator{hash_record: %s}", h.HashRecord.Inspect())
}

func (*HashRecordIterator) InstanceVariables() SymbolMap {
	return nil
}

func (h *HashRecordIterator) Next() (Value, Value) {
	for {
		if h.Index >= len(h.HashRecord.Table) {
			return Undefined, stopIterationSymbol.ToValue()
		}

		pair := h.HashRecord.Table[h.Index]
		h.Index++
		if !pair.Key.IsUndefined() {
			return Ref(&h.HashRecord.Table[h.Index-1]), Undefined
		}
	}
}

func (h *HashRecordIterator) Reset() {
	h.Index = 0
}

func initHashRecord() {
	HashRecordClass = NewClassWithOptions(ClassWithConstructor(HashRecordConstructor))
	HashRecordClass.IncludeMixin(RecordMixin)
	StdModule.AddConstantString("HashRecord", Ref(HashRecordClass))

	HashRecordIteratorClass = NewClass()
	HashRecordClass.AddConstantString("Iterator", Ref(HashRecordIteratorClass))
}
