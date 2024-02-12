package value

import (
	"fmt"
	"strings"
)

// ::Std::ArrayList
//
// Represents a dynamically sized array,
// that can shrink and grow.
var ArrayListClass *Class

// ::Std::ArrayList::Iterator
//
// ArrayList iterator class.
var ArrayListIteratorClass *Class

var NilArrayList ArrayList = nil

// Elk's ArrayList value
type ArrayList []Value

func NewArrayList(capacity int) *ArrayList {
	l := make(ArrayList, 0, capacity)
	return &l
}

func NewArrayListWithElements(capacity int, elements ...Value) *ArrayList {
	l := make(ArrayList, len(elements), len(elements)+capacity)
	copy(l, elements)
	return &l
}

func (*ArrayList) Class() *Class {
	return ArrayListClass
}

func (*ArrayList) DirectClass() *Class {
	return ArrayListClass
}

func (*ArrayList) SingletonClass() *Class {
	return nil
}

func (l *ArrayList) Copy() Value {
	if l == nil {
		return l
	}

	newList := make(ArrayList, len(*l))
	copy(newList, *l)
	return &newList
}

func (l *ArrayList) Inspect() string {
	var builder strings.Builder

	builder.WriteString("[")

	for i, element := range *l {
		if i != 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(element.Inspect())
	}

	builder.WriteString("]")
	leftCap := l.LeftCapacity()
	if leftCap > 0 {
		builder.WriteByte(':')
		fmt.Fprintf(&builder, "%d", leftCap)
	}
	return builder.String()
}

func (*ArrayList) InstanceVariables() SymbolMap {
	return nil
}

func (l *ArrayList) Capacity() int {
	return cap(*l)
}

func (l *ArrayList) LeftCapacity() int {
	return l.Capacity() - l.Length()
}

func (l *ArrayList) Length() int {
	return len(*l)
}

// Add new elements.
func (l *ArrayList) Append(elements ...Value) {
	*l = append(*l, elements...)
}

// Expand the array list to have
// empty slots for new elements.
func (l *ArrayList) Grow(newSlots int) {
	newList := make(ArrayList, l.Length(), l.Capacity()+newSlots)
	copy(newList, *l)
	*l = newList
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
func (l *ArrayList) Get(index int) (Value, *Error) {
	return GetFromSlice((*[]Value)(l), index)
}

// Get an element under the given index.
func (l *ArrayList) Subscript(key Value) (Value, *Error) {
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
func (l *ArrayList) Set(index int, val Value) *Error {
	return SetInSlice((*[]Value)(l), index, val)
}

// Set an element under the given index.
func (l *ArrayList) SubscriptSet(key, val Value) *Error {
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
func (l *ArrayList) Concat(other Value) (*ArrayList, *Error) {
	switch o := other.(type) {
	case *ArrayList:
		newList := make(ArrayList, len(*l), len(*l)+len(*o))
		copy(newList, *l)
		newList = append(newList, *o...)
		return &newList, nil
	case *ArrayTuple:
		newList := make(ArrayList, len(*l), len(*l)+len(*o))
		copy(newList, *l)
		newList = append(newList, *o...)
		return &newList, nil
	default:
		return nil, Errorf(TypeErrorClass, "cannot concat %s with list %s", other.Inspect(), l.Inspect())
	}
}

// Repeat the content of this list n times and return a new list containing the result.
// If the operation is illegal an error will be returned.
func (l *ArrayList) Repeat(other Value) (*ArrayList, *Error) {
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
		newList := make(ArrayList, 0, newLen)
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
func (l *ArrayList) Expand(newElements int) {
	if newElements < 1 {
		return
	}

	newCollection := make(ArrayList, len(*l), cap(*l)+newElements)
	copy(newCollection, *l)
	for i := 0; i < newElements; i++ {
		newCollection = append(newCollection, Nil)
	}
	*l = newCollection
}

type ArrayListIterator struct {
	ArrayList *ArrayList
	Index     int
}

func NewArrayListIterator(list *ArrayList) *ArrayListIterator {
	return &ArrayListIterator{
		ArrayList: list,
	}
}

func NewArrayListIteratorWithIndex(list *ArrayList, index int) *ArrayListIterator {
	return &ArrayListIterator{
		ArrayList: list,
		Index:     index,
	}
}

func (*ArrayListIterator) Class() *Class {
	return ArrayListIteratorClass
}

func (*ArrayListIterator) DirectClass() *Class {
	return ArrayListIteratorClass
}

func (*ArrayListIterator) SingletonClass() *Class {
	return nil
}

func (l *ArrayListIterator) Copy() Value {
	return &ArrayListIterator{
		ArrayList: l.ArrayList,
		Index:     l.Index,
	}
}

func (l *ArrayListIterator) Inspect() string {
	return fmt.Sprintf("Std::ArrayList::Iterator{list: %s, index: %d}", l.ArrayList.Inspect(), l.Index)
}

func (*ArrayListIterator) InstanceVariables() SymbolMap {
	return nil
}

var stopIterationSymbol = ToSymbol("stop_iteration")

func (l *ArrayListIterator) Next() (Value, Value) {
	if l.Index >= l.ArrayList.Length() {
		return nil, stopIterationSymbol
	}

	next := (*l.ArrayList)[l.Index]
	l.Index++
	return next, nil
}

func initArrayList() {
	ArrayListClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("ArrayList", ArrayListClass)

	ArrayListIteratorClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	ArrayListClass.AddConstantString("Iterator", ArrayListIteratorClass)
}
