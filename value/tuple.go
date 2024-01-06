package value

import (
	"fmt"
	"slices"
	"strings"
)

// ::Std::Tuple
//
// Represents an immutable array.
var TupleClass *Class

// Elk's Tuple value
type Tuple []Value

func (*Tuple) Class() *Class {
	return TupleClass
}

func (*Tuple) DirectClass() *Class {
	return TupleClass
}

func (*Tuple) SingletonClass() *Class {
	return nil
}

func (t *Tuple) Copy() Value {
	return t
}

// Add a new element.
func (t *Tuple) Append(element Value) {
	*t = append(*t, element)
}

func (t *Tuple) Inspect() string {
	var builder strings.Builder

	builder.WriteString("%[")

	for i, element := range *t {
		if i != 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(element.Inspect())
	}

	builder.WriteString("]")
	return builder.String()
}

func (*Tuple) InstanceVariables() SymbolMap {
	return nil
}

// Get an element under the given index.
func (t *Tuple) Get(index int) (Value, *Error) {
	l := len(*t)
	if index >= l || index < -l {
		return nil, NewIndexOutOfRangeError(fmt.Sprint(index), fmt.Sprint(len(*t)))
	}

	if index < 0 {
		index = l - index
	}

	return (*t)[index], nil
}

// Get an element under the given index.
func (t *Tuple) Subscript(key Value) (Value, *Error) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return nil, NewIndexOutOfRangeError(key.Inspect(), fmt.Sprint(len(*t)))
		}
		return nil, NewCoerceError(IntClass, key.Class())
	}

	return t.Get(i)
}

// Set an element under the given index.
func (t *Tuple) Set(index int, val Value) *Error {
	l := len(*t)
	if index >= l || index < -l {
		return NewIndexOutOfRangeError(fmt.Sprint(index), fmt.Sprint(len(*t)))
	}

	if index < 0 {
		index = l - index
	}

	(*t)[index] = val
	return nil
}

// Set an element under the given index.
func (t *Tuple) SubscriptSet(key, val Value) *Error {
	l := len(*t)
	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return NewIndexOutOfRangeError(key.Inspect(), fmt.Sprint(l))
		}
		return NewCoerceError(IntClass, key.Class())
	}

	return t.Set(i, val)
}

// Expands the tuple by n nil elements.
func (t *Tuple) Expand(newElements int) {
	if newElements < 1 {
		return
	}

	newCollection := slices.Grow(*t, newElements)
	for i := 0; i < newElements; i++ {
		newCollection = append(newCollection, Nil)
	}
	*t = newCollection
}

func initTuple() {
	TupleClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Tuple", TupleClass)
}
