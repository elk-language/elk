package value

import (
	"fmt"
	"slices"
)

// ::Std::List
//
// Represents a dynamically sized array,
// that can shrink and grow.
var ListClass *Class

// ::Std::List::Iterator
//
// List iterator class.
var ListIteratorClass *Class

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

func (l *List) Inspect() string {
	return InspectSlice(*l)
}

func (*List) InstanceVariables() SymbolMap {
	return nil
}

func (l *List) Length() int {
	return len(*l)
}

// Add a new element.
func (l *List) Append(element Value) {
	*l = append(*l, element)
}

// Get an element under the given index.
func GetFromSlice(collection *[]Value, index int) (Value, *Error) {
	l := len(*collection)
	if index >= l || index < -l {
		return nil, NewIndexOutOfRangeError(fmt.Sprint(index), len(*collection))
	}

	if index < 0 {
		index = l + index
	}

	return (*collection)[index], nil
}

// Set an element under the given index.
func SetInSlice(collection *[]Value, index int, val Value) *Error {
	l := len(*collection)
	if index >= l || index < -l {
		return NewIndexOutOfRangeError(fmt.Sprint(index), len(*collection))
	}

	if index < 0 {
		index = l + index
	}

	(*collection)[index] = val
	return nil
}

// Get an element under the given index.
func (l *List) Get(index int) (Value, *Error) {
	return GetFromSlice((*[]Value)(l), index)
}

// Get an element under the given index.
func (l *List) Subscript(key Value) (Value, *Error) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return nil, NewIndexOutOfRangeError(key.Inspect(), len(*l))
		}
		return nil, NewCoerceError(IntClass, key.Class())
	}

	return l.Get(i)
}

// Set an element under the given index.
func (l *List) Set(index int, val Value) *Error {
	return SetInSlice((*[]Value)(l), index, val)
}

// Set an element under the given index.
func (l *List) SubscriptSet(key, val Value) *Error {
	length := len(*l)
	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return NewIndexOutOfRangeError(key.Inspect(), length)
		}
		return NewCoerceError(IntClass, key.Class())
	}

	return l.Set(i, val)
}

// Concatenate another value with this list, creating a new list, and return the result.
// If the operation is illegal an error will be returned.
func (l *List) Concat(other Value) (*List, *Error) {
	switch o := other.(type) {
	case *List:
		newList := make(List, len(*l), len(*l)+len(*o))
		copy(newList, *l)
		newList = append(newList, *o...)
		return &newList, nil
	case *ArrayTuple:
		newList := make(List, len(*l), len(*l)+len(*o))
		copy(newList, *l)
		newList = append(newList, *o...)
		return &newList, nil
	default:
		return nil, Errorf(TypeErrorClass, "cannot concat %s with list %s", other.Inspect(), l.Inspect())
	}
}

// Repeat the content of this list n times and return a new list containing the result.
// If the operation is illegal an error will be returned.
func (l *List) Repeat(other Value) (*List, *Error) {
	switch o := other.(type) {
	case SmallInt:
		if o < 0 {
			return nil, Errorf(
				OutOfRangeErrorClass,
				"list repeat count cannot be negative: %s",
				o.Inspect(),
			)
		}
		newLen, ok := o.MultiplyOverflow(SmallInt(len(*l)))
		if !ok {
			return nil, Errorf(
				OutOfRangeErrorClass,
				"list repeat count is too large %s",
				o.Inspect(),
			)
		}
		newList := make(List, 0, newLen)
		for i := 0; i < int(o); i++ {
			newList = append(newList, *l...)
		}
		return &newList, nil
	case *BigInt:
		return nil, Errorf(
			OutOfRangeErrorClass,
			"list repeat count is too large %s",
			o.Inspect(),
		)
	default:
		return nil, Errorf(TypeErrorClass, "cannot repeat a list using %s", other.Inspect())
	}
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

type ListIterator struct {
	List  *List
	Index int
}

func NewListIterator(list *List) *ListIterator {
	return &ListIterator{
		List: list,
	}
}

func NewListIteratorWithIndex(list *List, index int) *ListIterator {
	return &ListIterator{
		List:  list,
		Index: index,
	}
}

func (*ListIterator) Class() *Class {
	return ListIteratorClass
}

func (*ListIterator) DirectClass() *Class {
	return ListIteratorClass
}

func (*ListIterator) SingletonClass() *Class {
	return nil
}

func (l *ListIterator) Copy() Value {
	return &ListIterator{
		List:  l.List,
		Index: l.Index,
	}
}

func (l *ListIterator) Inspect() string {
	return fmt.Sprintf("Std::List::Iterator{list: %s, index: %d}", l.List.Inspect(), l.Index)
}

func (*ListIterator) InstanceVariables() SymbolMap {
	return nil
}

var stopIterationSymbol = ToSymbol("stop_iteration")

func (l *ListIterator) Next() (Value, Value) {
	if l.Index >= l.List.Length() {
		return nil, stopIterationSymbol
	}

	next := (*l.List)[l.Index]
	l.Index++
	return next, nil
}

func initList() {
	ListClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("List", ListClass)

	ListIteratorClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	ListClass.AddConstantString("Iterator", ListIteratorClass)
}
