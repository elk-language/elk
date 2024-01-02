package value

import "strings"

// ::Std::Tuple
//
// Represents an immutable array.
var TupleClass *Class

// Elk's Tuple value
type Tuple []Value

func (Tuple) Class() *Class {
	return TupleClass
}

func (Tuple) DirectClass() *Class {
	return TupleClass
}

func (Tuple) SingletonClass() *Class {
	return nil
}

func (t Tuple) Copy() Value {
	return t
}

func (t Tuple) Inspect() string {
	var builder strings.Builder

	builder.WriteString("%[")

	for i, element := range t {
		if i != 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(element.Inspect())
	}

	builder.WriteString("]")
	return builder.String()
}

func (Tuple) InstanceVariables() SymbolMap {
	return nil
}

func initTuple() {
	TupleClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Tuple", TupleClass)
}
