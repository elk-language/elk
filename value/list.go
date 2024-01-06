package value

import "slices"

// ::Std::List
//
// Represents a dynamically sized array,
// that can shrink and grow.
var ListClass *Class

// Elk's List value
type List []Value

func (*List) Class() *Class {
	return ListClass
}

func (*List) DirectClass() *Class {
	return ListClass
}

func (*List) SingletonClass() *Class {
	return nil
}

func (l *List) Copy() Value {
	if l == nil {
		return l
	}

	newList := make(List, len(*l))
	copy(newList, *l)
	return &newList
}

// Add a new element.
func (l *List) Append(element Value) {
	*l = append(*l, element)
}

func (l *List) Inspect() string {
	return InspectSlice(*l)
}

func (*List) InstanceVariables() SymbolMap {
	return nil
}

// Expands the list by n nil elements.
func (l *List) Expand(newElements int) {
	if newElements < 1 {
		return
	}

	newCollection := slices.Grow(*l, newElements)
	for i := 0; i < newElements; i++ {
		newCollection = append(newCollection, Nil)
	}
	*l = newCollection
}

func initList() {
	ListClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("List", ListClass)
}
