package value

import (
	"fmt"
	"strings"
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

const MAX_HASH_RECORD_ELEMENTS_IN_INSPECT = 50

func (h *HashRecord) Inspect() string {
	var buffer strings.Builder
	buffer.WriteString("%{")

	first := true
	i := 0
	for _, entry := range h.Table {
		if entry.Key.IsUndefined() {
			continue
		}
		if first {
			first = false
		} else {
			buffer.WriteString(", ")
		}
		buffer.WriteString(entry.Key.Inspect())
		buffer.WriteString("=>")
		buffer.WriteString(entry.Value.Inspect())

		if i >= MAX_HASH_RECORD_ELEMENTS_IN_INSPECT-1 {
			buffer.WriteString(", ...")
			break
		}
		i++
	}
	buffer.WriteRune('}')

	return buffer.String()
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

func initHashRecord() {
	HashRecordClass = NewClassWithOptions(ClassWithConstructor(HashRecordConstructor))
	HashRecordClass.IncludeMixin(RecordMixin)
	StdModule.AddConstantString("HashRecord", Ref(HashRecordClass))

	HashRecordIteratorClass = NewClass()
	HashRecordClass.AddConstantString("Iterator", Ref(HashRecordIteratorClass))
}
