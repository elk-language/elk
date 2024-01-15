package value

import "strings"

const HashMapMaxLoad = 0.75

var HashMapClass *Class // ::Std::HashMap

type HashMap struct {
	Table []Pair
	Count int
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

func (h *HashMap) Capacity() int {
	return len(h.Table)
}

// TODO
func (h *HashMap) Copy() Value {
	return h
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
	return buffer.String()
}

func (*HashMap) InstanceVariables() SymbolMap {
	return nil
}

func initHashMap() {
	HashMapClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	StdModule.AddConstantString("HashMap", HashMapClass)
}
