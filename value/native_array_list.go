package value

import (
	"fmt"
	"iter"
	"strings"
	"unsafe"

	"github.com/elk-language/elk/indent"
)

// Elk's native array list implementation
type NativeArrayList[T ValueInterface] []T

var _ ArrayList = &NativeArrayList[String]{}

// UNSAFE! Cast a slice with native go types to an Elk `NativeArrayList` with corresponding Elk types.
// This is EXTREMELY unsafe, use it only if `I` and `O` have the same
// underlying type eg. `UnsafeCastNativeArrayList[string, value.String](slice)`, this will convert `[]string` to `value.NativeArrayList[value.String]`
func UnsafeCastNativeArrayList[I any, O ValueInterface](slice []I) NativeArrayList[O] {
	return *(*NativeArrayList[O])(unsafe.Pointer(&slice))
}

func CastNativeArrayListPtr[E ValueInterface](slice *[]E) *NativeArrayList[E] {
	return (*NativeArrayList[E])(slice)
}

// Transform a slice with native go types to a new Elk `NativeArrayList` with corresponding Elk types
// using the given function.
// eg.
//
//	TransformSliceIntoNativeArrayList(m, func(v uint8) (value.UInt8) {
//		return value.UInt8(v)
//	})
func TransformSliceIntoNativeArrayList[
	I any,
	O ValueInterface,
](
	slice []I,
	fn func(v I) O,
) *NativeArrayList[O] {
	newList := NewNativeArrayListWithLength[O](len(slice))
	for i, v := range slice {
		ov := fn(v)
		newList.SetAt(i, ov)
	}
	return newList
}

// Transform any elk ArrayList to a new Elk `NativeArrayList` with Elk types
// using the given function.
// eg.
//
//	TransformSliceIntoNativeArrayList(m, func(v uint8) (value.UInt8) {
//		return value.UInt8(v)
//	})
func TransformArrayListIntoNativeArrayList[T ValueInterface](
	list ArrayList,
	fn func(v Value) T,
) *NativeArrayList[T] {
	switch list := list.(type) {
	case *NativeArrayList[T]:
		return list
	}

	newList := NewNativeArrayListWithLength[T](list.Length())
	for i, v := range list.Elements() {
		ov := fn(v)
		newList.SetAt(i, ov)
	}
	return newList
}

func NewNativeArrayList[T ValueInterface](capacity int) *NativeArrayList[T] {
	l := make(NativeArrayList[T], 0, capacity)
	return &l
}

func NewNativeArrayListWithLength[T ValueInterface](length int) *NativeArrayList[T] {
	l := make(NativeArrayList[T], length)
	return &l
}

func NewNativeArrayListWithElements[T ValueInterface](capacity int, elements ...T) *NativeArrayList[T] {
	l := make(NativeArrayList[T], len(elements), len(elements)+capacity)
	copy(l, elements)
	return &l
}

func NewNativeArrayListWithElementsAndTotalCapacity[T ValueInterface](capacity int, elements ...T) *NativeArrayList[T] {
	l := make(NativeArrayList[T], len(elements), capacity)
	copy(l, elements)
	return &l
}

func (l *NativeArrayList[T]) ToSlice() []T {
	return *l
}

func (l *NativeArrayList[T]) NewArrayList(capacity int) ArrayList {
	return NewNativeArrayList[T](capacity)
}

func (l *NativeArrayList[T]) CloneArrayList(capacity int) ArrayList {
	newList := NewNativeArrayList[T](capacity)
	newList.Append(*l...)
	return newList
}

func (l *NativeArrayList[T]) SliceArrayList(from, to int) ArrayList {
	n := (*l)[from:to]
	return &n
}

func (l *NativeArrayList[T]) SliceArrayTuple(from, to int) ArrayTuple {
	return l.SliceArrayList(from, to)
}

func (l *NativeArrayList[T]) NewArrayTuple(capacity int) ArrayTuple {
	return l.NewArrayList(capacity)
}

func (l *NativeArrayList[T]) CloneArrayTuple(capacity int) ArrayTuple {
	return l.CloneArrayList(capacity)
}

func (l *NativeArrayList[T]) Elements() iter.Seq2[int, Value] {
	return func(yield func(int, Value) bool) {
		for i, element := range *l {
			if !yield(i, element.ToValue()) {
				return
			}
		}
	}
}

func (l *NativeArrayList[T]) Iterate() iter.Seq2[Value, Value] {
	return func(yield func(Value, Value) bool) {
		for _, element := range *l {
			if !yield(element.ToValue(), Undefined) {
				return
			}
		}
	}
}

func (*NativeArrayList[T]) Class() *Class {
	return ArrayListClass
}

func (*NativeArrayList[T]) DirectClass() *Class {
	return ArrayListClass
}

func (l *NativeArrayList[T]) ToValue() Value {
	return Ref(l)
}

func (*NativeArrayList[T]) SingletonClass() *Class {
	return nil
}

func (l *NativeArrayList[T]) Copy() Reference {
	if l == nil {
		return l
	}

	newList := make(NativeArrayList[T], len(*l))
	copy(newList, *l)
	return &newList
}

func (l *NativeArrayList[T]) Error() string {
	return l.Inspect()
}

func (l *NativeArrayList[T]) Inspect() string {
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

func (*NativeArrayList[T]) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *NativeArrayList[T]) Capacity() int {
	return cap(*l)
}

func (l *NativeArrayList[T]) RemoveAtErr(index int) Value {
	index, err := NormalizeArrayIndex(index, l.Length())
	if !err.IsUndefined() {
		return err
	}

	l.RemoveAt(index)
	return Undefined
}

func (l *NativeArrayList[T]) RemoveAt(i int) {
	s := *l
	copy(s[i:], s[i+1:])
	*l = s[:len(s)-1]
}

func (l *NativeArrayList[T]) LeftCapacity() int {
	return l.Capacity() - l.Length()
}

func (l *NativeArrayList[T]) Length() int {
	return len(*l)
}

// Add new elements.
func (l *NativeArrayList[T]) Append(elements ...T) {
	*l = append(*l, elements...)
}

func (l *NativeArrayList[T]) AppendVal(elements ...Value) Value {
	for _, element := range elements {
		e, ok := Downcast[T](element)
		if !ok {
			return NewInvalidElementInTypedArray(l, element.Class()).ToValue()
		}
		*l = append(*l, e)
	}
	return Undefined
}

// Expand the array list to have
// empty slots for new elements.
func (l *NativeArrayList[T]) Grow(newSlots int) {
	newList := make(NativeArrayList[T], l.Length(), l.Capacity()+newSlots)
	copy(newList, *l)
	*l = newList
}

// Get an element under the given index.
func (l *NativeArrayList[T]) Get(index int) (T, Value) {
	return GetFromSlice((*[]T)(l), index)
}

// Get an element under the given index without bounds checking
func (l *NativeArrayList[T]) At(i int) T {
	return (*l)[i]
}

func (l *NativeArrayList[T]) AtVal(i int) Value {
	return l.At(i).ToValue()
}

// Get an element under the given index.
func (l *NativeArrayList[T]) Subscript(key Value) (t Value, err Value) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return t, Ref(NewIndexOutOfRangeError(key.Inspect(), len(*l)))
		}
		return t, Ref(NewCoerceError(IntClass, key.Class()))
	}

	return ToValueErr(l.Get(i))
}

// Set an element under the given index.
func (l *NativeArrayList[T]) Set(index int, val T) Value {
	return SetInSlice((*[]T)(l), index, val)
}

// Set an element under the given index without bounds checking.
func (l *NativeArrayList[T]) SetAt(index int, val T) {
	(*l)[index] = val
}

func (l *NativeArrayList[T]) SetAtVal(index int, val Value) Value {
	v, ok := Downcast[T](val)
	if !ok {
		return NewInvalidElementInTypedArray(l, val.Class()).ToValue()
	}

	(*l)[index] = v
	return Undefined
}

// Set an element under the given index.
func (l *NativeArrayList[T]) SubscriptSet(key Value, val Value) Value {
	length := len(*l)
	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Ref(NewIndexOutOfRangeError(key.Inspect(), length))
		}
		return Ref(NewCoerceError(IntClass, key.Class()))
	}

	v, ok := Downcast[T](val)
	if !ok {
		return NewInvalidElementInTypedArray(l, val.Class()).ToValue()
	}
	return l.Set(i, v)
}

// Concatenate another value with this list, creating a new list, and return the result.
// If the operation is illegal an error will be returned.
func (l *NativeArrayList[T]) ConcatVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *NativeArrayList[T]:
			newList := make(NativeArrayList[T], len(*l), len(*l)+len(*o))
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
func (l *NativeArrayList[T]) Repeat(other Value) (*NativeArrayList[T], Value) {
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
		newList := make(NativeArrayList[T], 0, newLen)
		for range int(o) {
			newList = append(newList, *l...)
		}
		return &newList, Undefined
	default:
		return nil, Ref(Errorf(TypeErrorClass, "cannot repeat a list using %s", other.Inspect()))
	}
}

func (l *NativeArrayList[T]) RepeatVal(other Value) (Value, Value) {
	return RefErr(l.Repeat(other))
}

// Return an immutable box pointing to the slot with the given index.
func (l *NativeArrayList[T]) ImmutableBoxOfVal(index Value) (Value, Value) {
	b, err := l.boxOf(index)
	return b.ToImmutableBox().ToValue(), err
}

func (l *NativeArrayList[T]) ImmutableBoxOf(index int) (*ImmutableNativeBox[T], Value) {
	b, err := l.BoxOf(index)
	return b.ToImmutableBox(), err
}

// Return a box pointing to the slot with the given index.
func (l *NativeArrayList[T]) BoxOfVal(index Value) (Value, Value) {
	return RefErr(l.boxOf(index))
}

func (l *NativeArrayList[T]) boxOf(index Value) (*NativeBox[T], Value) {
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
func (l *NativeArrayList[T]) BoxOf(index int) (*NativeBox[T], Value) {
	index, err := NormalizeArrayIndex(index, l.Length())
	if !err.IsUndefined() {
		return nil, err
	}

	box := NewNativeBox(&(*l)[index])
	return box, Undefined
}

func (l *NativeArrayList[T]) IterNative() *NativeArrayListIterator[T] {
	return NewNativeArrayListIterator(l)
}

func (l *NativeArrayList[T]) Iter() NativeIterator {
	return l.IterNative()
}

func (l *NativeArrayList[T]) IterTuple() ArrayTupleIterator {
	return l.IterNative()
}

func (l *NativeArrayList[T]) IterList() ArrayListIterator {
	return l.IterNative()
}

type NativeArrayListIterator[T ValueInterface] struct {
	ArrayList *NativeArrayList[T]
	Index     int
}

var _ ArrayListIterator = &NativeArrayListIterator[String]{}

func NewNativeArrayListIterator[T ValueInterface](list *NativeArrayList[T]) *NativeArrayListIterator[T] {
	return &NativeArrayListIterator[T]{
		ArrayList: list,
	}
}

func NewNativeArrayListIteratorWithIndex[T ValueInterface](list *NativeArrayList[T], index int) *NativeArrayListIterator[T] {
	return &NativeArrayListIterator[T]{
		ArrayList: list,
		Index:     index,
	}
}

func (*NativeArrayListIterator[T]) SingletonClass() *Class {
	return nil
}

func (*NativeArrayListIterator[T]) Class() *Class {
	return ArrayListIteratorClass
}

func (*NativeArrayListIterator[T]) DirectClass() *Class {
	return ArrayListIteratorClass
}

func (l *NativeArrayListIterator[T]) Inspect() string {
	return fmt.Sprintf("Std::ArrayList::Iterator{&: %p, list: %s, index: %d}", l, l.ArrayList.Inspect(), l.Index)
}

func (l *NativeArrayListIterator[T]) Error() string {
	return l.Inspect()
}

func (i *NativeArrayListIterator[T]) ToValue() Value {
	return Ref(i)
}

func (l *NativeArrayListIterator[T]) Copy() Reference {
	return &NativeArrayListIterator[T]{
		ArrayList: l.ArrayList,
		Index:     l.Index,
	}
}

func (*NativeArrayListIterator[T]) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *NativeArrayListIterator[T]) Next() (t T, err Value) {
	if l.Index >= l.ArrayList.Length() {
		return t, stopIterationSymbol.ToValue()
	}

	next := (*l.ArrayList)[l.Index]
	l.Index++
	return next, Undefined
}

func (l *NativeArrayListIterator[T]) NextValue() (t Value, err Value) {
	return ToValueErr(l.Next())
}

func (l *NativeArrayListIterator[T]) Elements() iter.Seq[Value] {
	return func(yield func(Value) bool) {
		for ; l.Index >= l.ArrayList.Length(); l.Index++ {
			if !yield((*l.ArrayList)[l.Index].ToValue()) {
				return
			}
		}
	}
}

func (l *NativeArrayListIterator[T]) Reset() {
	l.Index = 0
}
