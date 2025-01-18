package value

import (
	"fmt"
	"slices"
	"strings"
)

const HashSetMaxLoad = 0.75

var HashSetClass *Class         // ::Std::HashSet
var HashSetIteratorClass *Class // ::Std::HashSet::Iterator

type HashSet struct {
	Table         []Value
	OccupiedSlots int
	Elements      int
}

func NewHashSet(capacity int) *HashSet {
	return &HashSet{
		Table: make([]Value, capacity),
	}
}

func (*HashSet) Class() *Class {
	return HashSetClass
}

func (*HashSet) DirectClass() *Class {
	return HashSetClass
}

func (*HashSet) SingletonClass() *Class {
	return nil
}

func (h *HashSet) Copy() Reference {
	newTable := slices.Clone(h.Table)
	return &HashSet{
		Table:         newTable,
		OccupiedSlots: h.OccupiedSlots,
		Elements:      h.Elements,
	}
}

func (h *HashSet) Error() string {
	return h.Inspect()
}

const MAX_HASH_SET_ELEMENTS_IN_INSPECT = 50

func (h *HashSet) Inspect() string {
	var buffer strings.Builder
	buffer.WriteString("^[")

	first := true
	i := 0
	for _, entry := range h.Table {
		if entry.IsUndefined() {
			continue
		}
		if first {
			first = false
		} else {
			buffer.WriteString(", ")
		}
		buffer.WriteString(entry.Inspect())

		if i >= MAX_HASH_SET_ELEMENTS_IN_INSPECT-1 {
			buffer.WriteString(", ...")
			break
		}
		i++
	}
	buffer.WriteRune(']')

	leftCapacity := h.LeftCapacity()
	if leftCapacity > 0 {
		fmt.Fprintf(&buffer, ":%d", leftCapacity)
	}
	return buffer.String()
}

func (*HashSet) InstanceVariables() SymbolMap {
	return nil
}

func (h *HashSet) Capacity() int {
	return len(h.Table)
}

func (h *HashSet) LeftCapacity() int {
	return h.Capacity() - h.Length()
}

func (h *HashSet) Length() int {
	return h.Elements
}

type HashSetIterator struct {
	HashSet *HashSet
	Index   int
}

func NewHashSetIterator(set *HashSet) *HashSetIterator {
	return &HashSetIterator{
		HashSet: set,
	}
}

func NewHashSetIteratorWithIndex(set *HashSet, index int) *HashSetIterator {
	return &HashSetIterator{
		HashSet: set,
		Index:   index,
	}
}

func (*HashSetIterator) Class() *Class {
	return HashSetIteratorClass
}

func (*HashSetIterator) DirectClass() *Class {
	return HashSetIteratorClass
}

func (*HashSetIterator) SingletonClass() *Class {
	return nil
}

func (h *HashSetIterator) Copy() Reference {
	return &HashSetIterator{
		HashSet: h.HashSet,
		Index:   h.Index,
	}
}

func (h *HashSetIterator) Error() string {
	return h.Inspect()
}

func (h *HashSetIterator) Inspect() string {
	return fmt.Sprintf("Std::HashSet::Iterator{hash_set: %s}", h.HashSet.Inspect())
}

func (*HashSetIterator) InstanceVariables() SymbolMap {
	return nil
}

func (h *HashSetIterator) Next() (Value, Value) {
	for {
		if h.Index >= h.HashSet.Capacity() {
			return Undefined, stopIterationSymbol.ToValue()
		}

		element := h.HashSet.Table[h.Index]
		h.Index++
		if !element.IsUndefined() {
			return element, Undefined
		}
	}
}

func (h *HashSetIterator) Reset() {
	h.Index = 0
}

func initHashSet() {
	HashSetClass = NewClass()
	HashSetClass.IncludeMixin(SetMixin)
	StdModule.AddConstantString("HashSet", Ref(HashSetClass))

	HashSetIteratorClass = NewClass()
	HashSetClass.AddConstantString("Iterator", Ref(HashSetIteratorClass))
}
