package value

import (
	"fmt"
	"iter"
	"strings"

	"github.com/elk-language/elk/indent"
)

// Elk's native array tuple implementation
type NativeArrayTuple[T ValueInterface] []T

var _ ArrayTuple = &NativeArrayTuple[String]{}

func NewNativeArrayTuple[T ValueInterface](capacity int) *NativeArrayTuple[T] {
	l := make(NativeArrayTuple[T], 0, capacity)
	return &l
}

func NewNativeArrayTupleWithLength[T ValueInterface](length int) *NativeArrayTuple[T] {
	l := make(NativeArrayTuple[T], length)
	return &l
}

func NewNativeArrayTupleWithElements[T ValueInterface](capacity int, elements ...T) *NativeArrayTuple[T] {
	l := make(NativeArrayTuple[T], len(elements), len(elements)+capacity)
	copy(l, elements)
	return &l
}

func (t *NativeArrayTuple[T]) Iter() *NativeArrayTupleIterator[T] {
	return NewNativeArrayTupleIterator(t)
}

func (l *NativeArrayTuple[T]) IterTuple() ArrayTupleIterator {
	return l.Iter()
}

func (*NativeArrayTuple[T]) Class() *Class {
	return ArrayTupleClass
}

func (*NativeArrayTuple[T]) DirectClass() *Class {
	return ArrayTupleClass
}

func (*NativeArrayTuple[T]) SingletonClass() *Class {
	return nil
}

func (t *NativeArrayTuple[T]) Copy() Reference {
	return t
}

func (l *NativeArrayTuple[T]) Elements() iter.Seq2[int, Value] {
	return func(yield func(int, Value) bool) {
		for i, element := range *l {
			if !yield(i, element.ToValue()) {
				return
			}
		}
	}
}

func (l *NativeArrayTuple[T]) Iterate() iter.Seq2[Value, Value] {
	return func(yield func(Value, Value) bool) {
		for _, element := range *l {
			if !yield(element.ToValue(), Undefined) {
				return
			}
		}
	}
}

func (t *NativeArrayTuple[T]) ToValue() Value {
	return Ref(t)
}

func (t *NativeArrayTuple[T]) Error() string {
	return t.Inspect()
}

// Add a new element.
func (t *NativeArrayTuple[T]) Append(element T) {
	*t = append(*t, element)
}

func (t *NativeArrayTuple[T]) AppendVal(element Value) {
	*t = append(*t, element.ToInterface().(T))
}

func (t *NativeArrayTuple[T]) Inspect() string {
	var hasMultilineElements bool
	elementStrings := make(
		[]string,
		0,
		min(MAX_ARRAY_TUPLE_ELEMENTS_IN_INSPECT, t.Length()),
	)

	for i, element := range *t {
		elementString := element.Inspect()
		elementStrings = append(elementStrings, elementString)
		if strings.ContainsRune(elementString, '\n') {
			hasMultilineElements = true
		}

		if i >= MAX_ARRAY_TUPLE_ELEMENTS_IN_INSPECT-1 {
			break
		}
	}

	var buff strings.Builder

	buff.WriteString("%[")

	if hasMultilineElements || t.Length() > 15 {
		buff.WriteRune('\n')
		for i, elementString := range elementStrings {
			if i != 0 {
				buff.WriteString(",\n")
			}

			indent.IndentString(&buff, elementString, 1)

			if i >= MAX_ARRAY_TUPLE_ELEMENTS_IN_INSPECT-1 {
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

			if i >= MAX_ARRAY_TUPLE_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(", ...")
				break
			}
		}
	}

	buff.WriteRune(']')
	return buff.String()
}

func (*NativeArrayTuple[T]) InstanceVariables() *InstanceVariables {
	return nil
}

// Get an element under the given index.
func (t *NativeArrayTuple[T]) Get(index int) (T, Value) {
	return GetFromSlice((*[]T)(t), index)
}

// Get an element under the given index without bounds checking.
func (t *NativeArrayTuple[T]) At(i int) T {
	return (*t)[i]
}

func (t *NativeArrayTuple[T]) AtVal(i int) Value {
	return t.At(i).ToValue()
}

// Get an element under the given index.
func (t *NativeArrayTuple[T]) Subscript(key Value) (Value, Value) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Undefined, Ref(NewIndexOutOfRangeError(key.Inspect(), len(*t)))
		}
		return Undefined, Ref(NewCoerceError(IntClass, key.Class()))
	}

	return ToValueErr(t.Get(i))
}

// Set an element under the given index.
func (t *NativeArrayTuple[T]) Set(index int, val T) Value {
	return SetInSlice((*[]T)(t), index, val)
}

// Set an element under the given index without bounds checking.
func (t *NativeArrayTuple[T]) SetAt(index int, val T) {
	(*t)[index] = val
}

func (t *NativeArrayTuple[T]) SetAtVal(index int, val Value) {
	(*t)[index] = val.ToInterface().(T)
}

// Set an element under the given index.
func (t *NativeArrayTuple[T]) SubscriptSet(key, val Value) Value {
	length := len(*t)
	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Ref(NewIndexOutOfRangeError(key.Inspect(), length))
		}
		return Ref(NewCoerceError(IntClass, key.Class()))
	}

	return t.Set(i, val.ToInterface().(T))
}

// Concatenate another value with this arrayTuple, creating a new value, and return the result.
// If the operation is illegal an error will be returned.
func (t *NativeArrayTuple[T]) ConcatVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *NativeArrayTuple[T]:
			newList := make(NativeArrayTuple[T], len(*t), len(*t)+len(*o))
			copy(newList, *t)
			newList = append(newList, *o...)
			return Ref(&newList), Undefined
		case ArrayList:
			newList := make(ArrayListOfValue, len(*t), len(*t)+o.Length())
			for i, element := range *t {
				newList[i] = element.ToValue()
			}
			for i, element := range o.Elements() {
				newList[i+o.Length()] = element
			}
			return Ref(&newList), Undefined
		case ArrayTuple:
			newList := make(ArrayTupleOfValue, len(*t), len(*t)+o.Length())
			for i, element := range *t {
				newList[i] = element.ToValue()
			}
			for i, element := range o.Elements() {
				newList[i+o.Length()] = element
			}
			return Ref(&newList), Undefined
		}
	}

	return Undefined, Ref(Errorf(TypeErrorClass, "cannot concat %s with arrayTuple %s", other.Inspect(), t.Inspect()))
}

// Repeat the content of this arrayTuple n times and return a new arrayTuple containing the result.
// If the operation is illegal an error will be returned.
func (t *NativeArrayTuple[T]) Repeat(other Value) (*NativeArrayTuple[T], Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return nil, Ref(Errorf(
				OutOfRangeErrorClass,
				"arrayTuple repeat count is too large %s",
				o.Inspect(),
			))
		default:
			return nil, Ref(Errorf(TypeErrorClass, "cannot repeat a arrayTuple using %s", other.Inspect()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		if o < 0 {
			return nil, Ref(Errorf(
				OutOfRangeErrorClass,
				"arrayTuple repeat count cannot be negative: %s",
				o.Inspect(),
			))
		}
		newLen, ok := o.MultiplyOverflow(SmallInt(len(*t)))
		if !ok {
			return nil, Ref(Errorf(
				OutOfRangeErrorClass,
				"arrayTuple repeat count is too large %s",
				o.Inspect(),
			))
		}
		newArrayTuple := make(NativeArrayTuple[T], 0, newLen)
		for i := 0; i < int(o); i++ {
			newArrayTuple = append(newArrayTuple, *t...)
		}
		return &newArrayTuple, Undefined
	default:
		return nil, Ref(Errorf(TypeErrorClass, "cannot repeat a arrayTuple using %s", other.Inspect()))
	}
}

func (t *NativeArrayTuple[T]) RepeatVal(other Value) (Value, Value) {
	return RefErr(t.Repeat(other))
}

// Return a box pointing to the slot with the given index.
func (l *NativeArrayTuple[T]) ImmutableBoxOfVal(index Value) (Value, Value) {
	var i int

	i, ok := ToGoInt(index)
	if !ok {
		if i == -1 {
			return Undefined, Ref(NewIndexOutOfRangeError(index.Inspect(), len(*l)))
		}
		return Undefined, Ref(NewCoerceError(IntClass, index.Class()))
	}

	return RefErr(l.ImmutableBoxOf(i))
}

// Return a box pointing to the slot with the given index.
func (l *NativeArrayTuple[T]) ImmutableBoxOf(index int) (*ImmutableNativeBox[T], Value) {
	index, err := NormalizeArrayIndex(index, l.Length())
	if !err.IsUndefined() {
		return nil, err
	}

	box := NewImmutableNativeBox(&(*l)[index])
	return box, Undefined
}

func (t *NativeArrayTuple[T]) Length() int {
	return len(*t)
}

type NativeArrayTupleIterator[T ValueInterface] struct {
	ArrayTuple *NativeArrayTuple[T]
	Index      int
}

var _ ArrayTupleIterator = &NativeArrayListIterator[String]{}

func NewNativeArrayTupleIterator[T ValueInterface](arrayTuple *NativeArrayTuple[T]) *NativeArrayTupleIterator[T] {
	return &NativeArrayTupleIterator[T]{
		ArrayTuple: arrayTuple,
	}
}

func NewNativeArrayTupleIteratorWithIndex[T ValueInterface](arrayTuple *NativeArrayTuple[T], index int) *NativeArrayTupleIterator[T] {
	return &NativeArrayTupleIterator[T]{
		ArrayTuple: arrayTuple,
		Index:      index,
	}
}

func (*NativeArrayTupleIterator[T]) Class() *Class {
	return ArrayTupleIteratorClass
}

func (*NativeArrayTupleIterator[T]) DirectClass() *Class {
	return ArrayTupleIteratorClass
}

func (*NativeArrayTupleIterator[T]) SingletonClass() *Class {
	return nil
}

func (t *NativeArrayTupleIterator[T]) Copy() Reference {
	return &NativeArrayTupleIterator[T]{
		ArrayTuple: t.ArrayTuple,
		Index:      t.Index,
	}
}

func (i *NativeArrayTupleIterator[T]) ToValue() Value {
	return Ref(i)
}

func (t *NativeArrayTupleIterator[T]) Inspect() string {
	return fmt.Sprintf("Std::ArrayTuple::Iterator{&: %p, tuple: %s, index: %d}", t, t.ArrayTuple.Inspect(), t.Index)
}

func (t *NativeArrayTupleIterator[T]) Error() string {
	return t.Inspect()
}

func (*NativeArrayTupleIterator[T]) InstanceVariables() *InstanceVariables {
	return nil
}

func (t *NativeArrayTupleIterator[T]) Next() (r T, err Value) {
	if t.Index >= t.ArrayTuple.Length() {
		return r, stopIterationSymbol.ToValue()
	}

	next := (*t.ArrayTuple)[t.Index]
	t.Index++
	return next, Undefined
}

func (t *NativeArrayTupleIterator[T]) NextValue() (Value, Value) {
	return ToValueErr(t.Next())
}

func (t *NativeArrayTupleIterator[T]) Elements() iter.Seq[Value] {
	return func(yield func(Value) bool) {
		for ; t.Index >= t.ArrayTuple.Length(); t.Index++ {
			if !yield((*t.ArrayTuple)[t.Index].ToValue()) {
				return
			}
		}
	}
}

func (t *NativeArrayTupleIterator[T]) Reset() {
	t.Index = 0
}
