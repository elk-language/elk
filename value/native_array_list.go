package value

import (
	"fmt"
	"iter"
	"strings"

	"github.com/elk-language/elk/indent"
)

// Elk's native array list implementation
type nativeArrayList[T ValueInterface] []T

func newNativeArrayList[T ValueInterface](capacity int) *nativeArrayList[T] {
	l := make(nativeArrayList[T], 0, capacity)
	return &l
}

func newNativeArrayListWithLength[T ValueInterface](length int) *nativeArrayList[T] {
	l := make(nativeArrayList[T], length)
	return &l
}

func newNativeArrayListWithElements[T ValueInterface](capacity int, elements ...T) *nativeArrayList[T] {
	l := make(nativeArrayList[T], len(elements), len(elements)+capacity)
	copy(l, elements)
	return &l
}

func (l *nativeArrayList[T]) Elements() iter.Seq2[int, Value] {
	return func(yield func(int, Value) bool) {
		for i, element := range *l {
			if !yield(i, element.ToValue()) {
				return
			}
		}
	}
}

func (*nativeArrayList[T]) Class() *Class {
	return ArrayListClass
}

func (*nativeArrayList[T]) DirectClass() *Class {
	return ArrayListClass
}

func (l *nativeArrayList[T]) ToValue() Value {
	return Ref(l)
}

func (*nativeArrayList[T]) SingletonClass() *Class {
	return nil
}

func (l *nativeArrayList[T]) Copy() Reference {
	if l == nil {
		return l
	}

	newList := make(nativeArrayList[T], len(*l))
	copy(newList, *l)
	return &newList
}

func (l *nativeArrayList[T]) Error() string {
	return l.Inspect()
}

func (l *nativeArrayList[T]) Inspect() string {
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

func (*nativeArrayList[T]) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *nativeArrayList[T]) Capacity() int {
	return cap(*l)
}

func (l *nativeArrayList[T]) RemoveAtErr(index int) Value {
	index, err := NormalizeArrayIndex(index, l.Length())
	if !err.IsUndefined() {
		return err
	}

	l.RemoveAt(index)
	return Undefined
}

func (l *nativeArrayList[T]) RemoveAt(i int) {
	s := *l
	copy(s[i:], s[i+1:])
	*l = s[:len(s)-1]
}

func (l *nativeArrayList[T]) LeftCapacity() int {
	return l.Capacity() - l.Length()
}

func (l *nativeArrayList[T]) Length() int {
	return len(*l)
}

// Add new elements.
func (l *nativeArrayList[T]) Append(elements ...T) {
	*l = append(*l, elements...)
}

// Expand the array list to have
// empty slots for new elements.
func (l *nativeArrayList[T]) Grow(newSlots int) {
	newList := make(nativeArrayList[T], l.Length(), l.Capacity()+newSlots)
	copy(newList, *l)
	*l = newList
}

// Get an element under the given index.
func (l *nativeArrayList[T]) Get(index int) (T, Value) {
	return GetFromSlice((*[]T)(l), index)
}

// Get an element under the given index without bounds checking
func (l *nativeArrayList[T]) At(i int) T {
	return (*l)[i]
}

// Get an element under the given index.
func (l *nativeArrayList[T]) Subscript(key Value) (t T, err Value) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return t, Ref(NewIndexOutOfRangeError(key.Inspect(), len(*l)))
		}
		return t, Ref(NewCoerceError(IntClass, key.Class()))
	}

	return l.Get(i)
}

// Set an element under the given index.
func (l *nativeArrayList[T]) Set(index int, val T) Value {
	return SetInSlice((*[]T)(l), index, val)
}

// Set an element under the given index without bounds checking.
func (l *nativeArrayList[T]) SetAt(index int, val T) {
	(*l)[index] = val
}

// Set an element under the given index.
func (l *nativeArrayList[T]) SubscriptSet(key Value, val T) Value {
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
func (l *nativeArrayList[T]) Concat(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *nativeArrayList[T]:
			newList := make(nativeArrayList[T], len(*l), len(*l)+len(*o))
			copy(newList, *l)
			newList = append(newList, *o...)
			return Ref(&newList), Undefined
		case ArrayTuple:
			newList := make(ArrayListOfValue, len(*l), len(*l)+o.Length())
			for i, element := range *l {
				newList[i] = element.ToValue()
			}
			for i, element := range o.Elements() {
				newList[i+o.Length()] = element
			}
			return Ref(&newList), Undefined
		}
	}

	return Undefined, Ref(Errorf(TypeErrorClass, "cannot concat %s with list %s", other.Inspect(), l.Inspect()))
}

// Repeat the content of this list n times and return a new list containing the result.
// If the operation is illegal an error will be returned.
func (l *nativeArrayList[T]) Repeat(other Value) (*nativeArrayList[T], Value) {
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
		newList := make(nativeArrayList[T], 0, newLen)
		for range int(o) {
			newList = append(newList, *l...)
		}
		return &newList, Undefined
	default:
		return nil, Ref(Errorf(TypeErrorClass, "cannot repeat a list using %s", other.Inspect()))
	}
}

// Return an immutable box pointing to the slot with the given index.
func (l *nativeArrayList[T]) ImmutableBoxOfVal(index Value) (*ImmutableNativeBox[T], Value) {
	b, err := l.BoxOfVal(index)
	return b.ToImmutableBox(), err
}

// Return a box pointing to the slot with the given index.
func (l *nativeArrayList[T]) BoxOfVal(index Value) (*NativeBox[T], Value) {
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
func (l *nativeArrayList[T]) BoxOf(index int) (*NativeBox[T], Value) {
	index, err := NormalizeArrayIndex(index, l.Length())
	if !err.IsUndefined() {
		return nil, err
	}

	box := NewNativeBox(&(*l)[index])
	return box, Undefined
}

type nativeArrayListIterator[T ValueInterface] struct {
	ArrayList *nativeArrayList[T]
	Index     int
}

func newNativeArrayListIterator[T ValueInterface](list *nativeArrayList[T]) *nativeArrayListIterator[T] {
	return &nativeArrayListIterator[T]{
		ArrayList: list,
	}
}

func newNativeArrayListIteratorWithIndex[T ValueInterface](list *nativeArrayList[T], index int) *nativeArrayListIterator[T] {
	return &nativeArrayListIterator[T]{
		ArrayList: list,
		Index:     index,
	}
}

func (*nativeArrayListIterator[T]) SingletonClass() *Class {
	return nil
}

func (l *nativeArrayListIterator[T]) inspect() string {
	return fmt.Sprintf("{&: %p, list: %s, index: %d}", l, l.ArrayList.Inspect(), l.Index)
}

func (*nativeArrayListIterator[T]) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *nativeArrayListIterator[T]) Next() (t T, err Value) {
	if l.Index >= l.ArrayList.Length() {
		return t, stopIterationSymbol.ToValue()
	}

	next := (*l.ArrayList)[l.Index]
	l.Index++
	return next, Undefined
}

func (l *nativeArrayListIterator[T]) Reset() {
	l.Index = 0
}
