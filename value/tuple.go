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

// ::Std::Tuple::Iterator
//
// Tuple iterator class.
var TupleIteratorClass *Class

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
	return GetFromSlice((*[]Value)(t), index)
}

// Get an element under the given index.
func (t *Tuple) Subscript(key Value) (Value, *Error) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return nil, NewIndexOutOfRangeError(key.Inspect(), len(*t))
		}
		return nil, NewCoerceError(IntClass, key.Class())
	}

	return t.Get(i)
}

// Concatenate another value with this tuple, creating a new value, and return the result.
// If the operation is illegal an error will be returned.
func (t *Tuple) Concat(other Value) (Value, *Error) {
	switch o := other.(type) {
	case *List:
		newList := make(List, len(*t), len(*t)+len(*o))
		copy(newList, *t)
		newList = append(newList, *o...)
		return &newList, nil
	case *Tuple:
		newTuple := make(Tuple, len(*t), len(*t)+len(*o))
		copy(newTuple, *t)
		newTuple = append(newTuple, *o...)
		return &newTuple, nil
	default:
		return nil, Errorf(TypeErrorClass, "cannot concat %s with tuple %s", other.Inspect(), t.Inspect())
	}
}

// Repeat the content of this tuple n times and return a new tuple containing the result.
// If the operation is illegal an error will be returned.
func (t *Tuple) Repeat(other Value) (*Tuple, *Error) {
	switch o := other.(type) {
	case SmallInt:
		if o < 0 {
			return nil, Errorf(
				OutOfRangeErrorClass,
				"tuple repeat count cannot be negative: %s",
				o.Inspect(),
			)
		}
		newLen, ok := o.MultiplyOverflow(SmallInt(len(*t)))
		if !ok {
			return nil, Errorf(
				OutOfRangeErrorClass,
				"tuple repeat count is too large %s",
				o.Inspect(),
			)
		}
		newTuple := make(Tuple, 0, newLen)
		for i := 0; i < int(o); i++ {
			newTuple = append(newTuple, *t...)
		}
		return &newTuple, nil
	case *BigInt:
		return nil, Errorf(
			OutOfRangeErrorClass,
			"tuple repeat count is too large %s",
			o.Inspect(),
		)
	default:
		return nil, Errorf(TypeErrorClass, "cannot repeat a tuple using %s", other.Inspect())
	}
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

func (t *Tuple) Length() int {
	return len(*t)
}

type TupleIterator struct {
	Tuple *Tuple
	Index int
}

func NewTupleIterator(tuple *Tuple) *TupleIterator {
	return &TupleIterator{
		Tuple: tuple,
	}
}

func NewTupleIteratorWithIndex(tuple *Tuple, index int) *TupleIterator {
	return &TupleIterator{
		Tuple: tuple,
		Index: index,
	}
}

func (*TupleIterator) Class() *Class {
	return TupleIteratorClass
}

func (*TupleIterator) DirectClass() *Class {
	return TupleIteratorClass
}

func (*TupleIterator) SingletonClass() *Class {
	return nil
}

func (t *TupleIterator) Copy() Value {
	return &TupleIterator{
		Tuple: t.Tuple,
		Index: t.Index,
	}
}

func (t *TupleIterator) Inspect() string {
	return fmt.Sprintf("Std::Tuple::Iterator{tuple: %s, index: %d}", t.Tuple.Inspect(), t.Index)
}

func (*TupleIterator) InstanceVariables() SymbolMap {
	return nil
}

func (t *TupleIterator) Next() (Value, Value) {
	if t.Index >= t.Tuple.Length() {
		return nil, stopIterationSymbol
	}

	next := (*t.Tuple)[t.Index]
	t.Index++
	return next, nil
}

func initTuple() {
	TupleClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Tuple", TupleClass)
}
