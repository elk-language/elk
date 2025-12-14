package value

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/elk-language/elk/indent"
)

const HashMapMaxLoad = 0.75

var HashMapClass *Class         // ::Std::HashMap
var HashMapIteratorClass *Class // ::Std::HashMap::Iterator

type HashMap struct {
	Table         []Pair
	OccupiedSlots int
	Elements      int
}

func HashMapConstructor(class *Class) Value {
	return Ref(&HashMap{})
}

func NewHashMap(capacity int) *HashMap {
	return &HashMap{
		Table: make([]Pair, capacity),
	}
}

func (h *HashMap) All() iter.Seq[Pair] {
	return func(yield func(Pair) bool) {
		for _, pair := range h.Table {
			if pair.Key.IsUndefined() {
				continue
			}

			if !yield(pair) {
				return
			}
		}
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

const MAX_HASH_MAP_ELEMENTS_IN_INSPECT = 300

func (h *HashMap) Inspect() string {
	var hasMultilineElements bool
	keyStrings := make(
		[]string,
		0,
		min(MAX_HASH_MAP_ELEMENTS_IN_INSPECT, h.Length()),
	)
	valStrings := make(
		[]string,
		0,
		min(MAX_HASH_MAP_ELEMENTS_IN_INSPECT, h.Length()),
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

		if i >= MAX_HASH_MAP_ELEMENTS_IN_INSPECT-1 {
			break
		}
		i++
	}

	var buff strings.Builder

	buff.WriteRune('{')
	if hasMultilineElements || h.Length() > 15 {
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

			if i >= MAX_HASH_MAP_ELEMENTS_IN_INSPECT-1 {
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

			if i >= MAX_HASH_MAP_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(", ...")
				break
			}
		}
	}
	buff.WriteRune('}')

	leftCapacity := h.LeftCapacity()
	if leftCapacity > 0 {
		fmt.Fprintf(&buff, ":%d", leftCapacity)
	}
	return buff.String()
}

func (*HashMap) InstanceVariables() *InstanceVariables {
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

func (*HashMapIterator) InstanceVariables() *InstanceVariables {
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

func (h *HashMapIterator) Reset() {
	h.Index = 0
}

func initHashMap() {
	HashMapClass = NewClassWithOptions(ClassWithConstructor(HashMapConstructor))
	HashMapClass.IncludeMixin(MapMixin)
	StdModule.AddConstantString("HashMap", Ref(HashMapClass))

	HashMapIteratorClass = NewClass()
	HashMapClass.AddConstantString("Iterator", Ref(HashMapIteratorClass))
}
