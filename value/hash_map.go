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

func (h *HashMap) Copy() Value {
	newTable := slices.Clone(h.Table)
	return &HashMap{
		Table:         newTable,
		OccupiedSlots: h.OccupiedSlots,
		Elements:      h.Elements,
	}
}

func (h *HashMap) Inspect() string {
	var buffer strings.Builder
	buffer.WriteRune('{')

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

func (h *HashMapIterator) Copy() Value {
	return &HashMapIterator{
		HashMap: h.HashMap,
		Index:   h.Index,
	}
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
			return nil, stopIterationSymbol
		}

		pair := h.HashMap.Table[h.Index]
		h.Index++
		if pair.Key != nil {
			return &pair, nil
		}
	}
}

func initHashMap() {
	HashMapClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	StdModule.AddConstantString("HashMap", HashMapClass)

	HashMapIteratorClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	HashMapClass.AddConstantString("Iterator", HashMapIteratorClass)
}
