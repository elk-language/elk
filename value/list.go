package value

// ::Std::List
//
// Represents a dynamically sized array,
// that can shrink and grow.
var ListClass *Class

// Elk's List value
type List []Value

func (List) Class() *Class {
	return ListClass
}

func (List) DirectClass() *Class {
	return ListClass
}

func (List) SingletonClass() *Class {
	return nil
}

func (l List) Copy() Value {
	if l == nil {
		return l
	}

	newList := make(List, len(l))
	copy(newList, l)
	return newList
}

func (l List) Inspect() string {
	return InspectSlice(l)
}

func (List) InstanceVariables() SymbolMap {
	return nil
}

func initList() {
	ListClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("List", ListClass)
}
