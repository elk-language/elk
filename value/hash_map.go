package value

import (
	"fmt"
	"slices"
	"strings"
)

const HashMapMaxLoad = 0.75

var HashMapClass *Class         // ::Std::HashMap
var HashMapIteratorClass *Class // ::Std::HashMap::Iterator

type HashMap struct {
	Table         []Pair
	OccupiedSlots int
	Elements      int
}

func NewHashMap(capacity int) *HashMap {
	return &HashMap{
		Table: make([]Pair, capacity),
	}
}

func (*HashMap) Class() *Class {
	return HashMapClass
}

func (*HashMap) DirectClass() *Class {
	return HashMapClass
}

func (*HashMap) SingletonClass() *Class {
	return nil
}

func (h *HashMap) Clone() *HashMap {
	newTable := slices.Clone(h.Table)
	return &HashMap{
		Table:         newTable,
		OccupiedSlots: h.OccupiedSlots,
		Elements:      h.Elements,
	}
}

func (h *HashMap) Copy() Reference {
	return h.Clone()
}

func (h *HashMap) Error() string {
	return h.Inspect()
}

func (h *HashMap) Inspect() string {
	var buffer strings.Builder
	buffer.WriteRune('{')

	first := true
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
	}
	buffer.WriteRune('}')

	leftCapacity := h.LeftCapacity()
	if leftCapacity > 0 {
		fmt.Fprintf(&buffer, ":%d", leftCapacity)
	}
	return buffer.String()
}

func (*HashMap) InstanceVariables() SymbolMap {
	return nil
}

func (h *HashMap) Capacity() int {
	return len(h.Table)
}

func (h *HashMap) LeftCapacity() int {
	return h.Capacity() - h.Length()
}

func (h *HashMap) Length() int {
	return h.Elements
}

type HashMapIterator struct {
	HashMap *HashMap
	Index   int
}

func NewHashMapIterator(hmap *HashMap) *HashMapIterator {
	return &HashMapIterator{
		HashMap: hmap,
	}
}

func NewHashMapIteratorWithIndex(hmap *HashMap, index int) *HashMapIterator {
	return &HashMapIterator{
		HashMap: hmap,
		Index:   index,
	}
}

func (*HashMapIterator) Class() *Class {
	return HashMapIteratorClass
}

func (*HashMapIterator) DirectClass() *Class {
	return HashMapIteratorClass
}

func (*HashMapIterator) SingletonClass() *Class {
	return nil
}

func (h *HashMapIterator) Copy() Reference {
	return &HashMapIterator{
		HashMap: h.HashMap,
		Index:   h.Index,
	}
}

func (h *HashMapIterator) Error() string {
	return h.Inspect()
}

func (h *HashMapIterator) Inspect() string {
	return fmt.Sprintf("Std::HashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*HashMapIterator) InstanceVariables() SymbolMap {
	return nil
}

func (h *HashMapIterator) Next() (Value, Value) {
	for {
		if h.Index >= h.HashMap.Capacity() {
			return Undefined, stopIterationSymbol.ToValue()
		}

		pair := h.HashMap.Table[h.Index]
		h.Index++
		if !pair.Key.IsUndefined() {
			return Ref(&h.HashMap.Table[h.Index-1]), Undefined
		}
	}
}

func initHashMap() {
	HashMapClass = NewClass()
	HashMapClass.IncludeMixin(MapMixin)
	StdModule.AddConstantString("HashMap", Ref(HashMapClass))

	HashMapIteratorClass = NewClass()
	HashMapClass.AddConstantString("Iterator", Ref(HashMapIteratorClass))
}
