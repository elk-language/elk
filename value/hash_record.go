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

func (*HashRecord) Class() *Class {
	return HashRecordClass
}

func (*HashRecord) DirectClass() *Class {
	return HashRecordClass
}

func (*HashRecord) SingletonClass() *Class {
	return nil
}

func (h *HashRecord) Copy() Value {
	return h
}

func (h *HashRecord) Error() string {
	return h.Inspect()
}

func (h *HashRecord) Inspect() string {
	var buffer strings.Builder
	buffer.WriteString("%{")

	first := true
	for _, entry := range h.Table {
		if entry.Key == nil {
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

func (h *HashRecordIterator) Copy() Value {
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
			return nil, stopIterationSymbol
		}

		pair := h.HashRecord.Table[h.Index]
		h.Index++
		if pair.Key != nil {
			return &h.HashRecord.Table[h.Index-1], nil
		}
	}
}

func initHashRecord() {
	HashRecordClass = NewClass()
	HashRecordClass.IncludeMixin(RecordMixin)
	StdModule.AddConstantString("HashRecord", HashRecordClass)

	HashRecordIteratorClass = NewClass()
	HashRecordClass.AddConstantString("Iterator", HashRecordIteratorClass)
}
