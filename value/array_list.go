package value

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
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

// Elk's ArrayList value
type ArrayList []Value

func ArrayListConstructor(class *Class) Value {
	return Ref(&ArrayList{})
}

func NewArrayList(capacity int) *ArrayList {
	l := make(ArrayList, 0, capacity)
	return &l
}

func NewArrayListWithLength(length int) *ArrayList {
	l := make(ArrayList, length)
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

func (l *ArrayList) Copy() Reference {
	if l == nil {
		return l
	}

	newList := make(ArrayList, len(*l))
	copy(newList, *l)
	return &newList
}

func (l *ArrayList) Error() string {
	return l.Inspect()
}

const MAX_ARRAY_LIST_ELEMENTS_IN_INSPECT = 300

func (l *ArrayList) Inspect() string {
	var hasMultilineElements bool
	elementStrings := make(
		[]string,
		0,
		min(MAX_ARRAY_LIST_ELEMENTS_IN_INSPECT, l.Length()),
	)

	for i, element := range *l {
		elementString := element.Inspect()
		elementStrings = append(elementStrings, elementString)
		if strings.ContainsRune(elementString, '\n') {
			hasMultilineElements = true
		}

		if i >= MAX_ARRAY_LIST_ELEMENTS_IN_INSPECT-1 {
			break
		}
	}

	var buff strings.Builder

	buff.WriteRune('[')
	if hasMultilineElements || l.Length() > 15 {
		buff.WriteRune('\n')
		for i, elementString := range elementStrings {
			if i != 0 {
				buff.WriteString(",\n")
			}

			indent.IndentString(&buff, elementString, 1)

			if i >= MAX_ARRAY_LIST_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(",\n  ...")
				break
			}
		}
		buff.WriteRune('\n')
	} else {
		for i, elementString := range elementStrings {
			if i != 0 {
				buff.WriteString(", ")
			}

			buff.WriteString(elementString)

			if i >= MAX_ARRAY_LIST_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(", ...")
				break
			}
		}
	}

	buff.WriteRune(']')
	leftCap := l.LeftCapacity()
	if leftCap > 0 {
		buff.WriteByte(':')
		fmt.Fprintf(&buff, "%d", leftCap)
	}
	return buff.String()
}

func (*ArrayList) InstanceVariables() *InstanceVariables {
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
func GetFromSlice[V any](collection *[]V, index int) (ret V, err Value) {
	l := len(*collection)
	if index >= l || index < -l {
		return ret, Ref(NewIndexOutOfRangeError(fmt.Sprint(index), len(*collection)))
	}

	if index < 0 {
		index = l + index
	}

	return (*collection)[index], Undefined
}

// Set an element under the given index.
func SetInSlice[V any](collection *[]V, index int, val V) (err Value) {
	l := len(*collection)
	if index >= l || index < -l {
		return Ref(NewIndexOutOfRangeError(fmt.Sprint(index), len(*collection)))
	}

	if index < 0 {
		index = l + index
	}

	(*collection)[index] = val
	return Undefined
}

// Get an element under the given index.
func (l *ArrayList) Get(index int) (Value, Value) {
	return GetFromSlice((*[]Value)(l), index)
}

// Get an element under the given index without bounds checking
func (l *ArrayList) At(i int) Value {
	return (*l)[i]
}

// Get an element under the given index.
func (l *ArrayList) Subscript(key Value) (Value, Value) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Undefined, Ref(NewIndexOutOfRangeError(key.Inspect(), len(*l)))
		}
		return Undefined, Ref(NewCoerceError(IntClass, key.Class()))
	}

	return l.Get(i)
}

// Set an element under the given index.
func (l *ArrayList) Set(index int, val Value) Value {
	return SetInSlice((*[]Value)(l), index, val)
}

// Set an element under the given index without bounds checking.
func (l *ArrayList) SetAt(index int, val Value) {
	(*l)[index] = val
}

// Set an element under the given index.
func (l *ArrayList) SubscriptSet(key, val Value) Value {
	length := len(*l)
	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Ref(NewIndexOutOfRangeError(key.Inspect(), length))
		}
		return Ref(NewCoerceError(IntClass, key.Class()))
	}

	return l.Set(i, val)
}

// Concatenate another value with this list, creating a new list, and return the result.
// If the operation is illegal an error will be returned.
func (l *ArrayList) Concat(other Value) (*ArrayList, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *ArrayList:
			newList := make(ArrayList, len(*l), len(*l)+len(*o))
			copy(newList, *l)
			newList = append(newList, *o...)
			return &newList, Undefined
		case *ArrayTuple:
			newList := make(ArrayList, len(*l), len(*l)+len(*o))
			copy(newList, *l)
			newList = append(newList, *o...)
			return &newList, Undefined
		}
	}

	return nil, Ref(Errorf(TypeErrorClass, "cannot concat %s with list %s", other.Inspect(), l.Inspect()))
}

// Repeat the content of this list n times and return a new list containing the result.
// If the operation is illegal an error will be returned.
func (l *ArrayList) Repeat(other Value) (*ArrayList, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return nil, Ref(Errorf(
				OutOfRangeErrorClass,
				"list repeat count is too large %s",
				o.Inspect(),
			))
		default:
			return nil, Ref(Errorf(TypeErrorClass, "cannot repeat a list using %s", other.Inspect()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		if o < 0 {
			return nil, Ref(Errorf(
				OutOfRangeErrorClass,
				"list repeat count cannot be negative: %s",
				o.Inspect(),
			))
		}
		newLen, ok := o.MultiplyOverflow(SmallInt(len(*l)))
		if !ok {
			return nil, Ref(Errorf(
				OutOfRangeErrorClass,
				"list repeat count is too large %s",
				o.Inspect(),
			))
		}
		newList := make(ArrayList, 0, newLen)
		for range int(o) {
			newList = append(newList, *l...)
		}
		return &newList, Undefined
	default:
		return nil, Ref(Errorf(TypeErrorClass, "cannot repeat a list using %s", other.Inspect()))
	}
}

// Return an immutable box pointing to the slot with the given index.
func (l *ArrayList) ImmutableBoxOfVal(index Value) (*ImmutableBox, Value) {
	b, err := l.BoxOfVal(index)
	return b.ToImmutableBox(), err
}

// Return a box pointing to the slot with the given index.
func (l *ArrayList) BoxOfVal(index Value) (*Box, Value) {
	var i int

	i, ok := ToGoInt(index)
	if !ok {
		if i == -1 {
			return nil, Ref(NewIndexOutOfRangeError(index.Inspect(), len(*l)))
		}
		return nil, Ref(NewCoerceError(IntClass, index.Class()))
	}

	return l.BoxOf(i)
}

// Return a box pointing to the slot with the given index.
func (l *ArrayList) BoxOf(index int) (*Box, Value) {
	len := l.Length()
	if index >= len || index < -len {
		return nil, Ref(NewIndexOutOfRangeError(fmt.Sprint(index), len))
	}

	if index < 0 {
		index = len + index
	}

	box := (*Box)(&(*l)[index])
	return box, Undefined
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

func (l *ArrayListIterator) Copy() Reference {
	return &ArrayListIterator{
		ArrayList: l.ArrayList,
		Index:     l.Index,
	}
}

func (l *ArrayListIterator) Inspect() string {
	return fmt.Sprintf("Std::ArrayList::Iterator{&: %p, list: %s, index: %d}", l, l.ArrayList.Inspect(), l.Index)
}

func (l *ArrayListIterator) Error() string {
	return l.Inspect()
}

func (*ArrayListIterator) InstanceVariables() *InstanceVariables {
	return nil
}

var stopIterationSymbol = ToSymbol("stop_iteration")

func (l *ArrayListIterator) Next() (Value, Value) {
	if l.Index >= l.ArrayList.Length() {
		return Undefined, stopIterationSymbol.ToValue()
	}

	next := (*l.ArrayList)[l.Index]
	l.Index++
	return next, Undefined
}

func (l *ArrayListIterator) Reset() {
	l.Index = 0
}

func initArrayList() {
	ArrayListClass = NewClassWithOptions(ClassWithConstructor(ArrayListConstructor))
	ArrayListClass.IncludeMixin(ListMixin)
	StdModule.AddConstantString("ArrayList", Ref(ArrayListClass))

	ArrayListIteratorClass = NewClass()
	ArrayListClass.AddConstantString("Iterator", Ref(ArrayListIteratorClass))
}
