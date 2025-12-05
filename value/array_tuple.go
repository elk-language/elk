package value

import (
	"fmt"
	"slices"
	"strings"

	"github.com/elk-language/elk/indent"
)

// ::Std::ArrayTuple
//
// Represents an immutable array.
var ArrayTupleClass *Class

// ::Std::ArrayTuple::Iterator
//
// ArrayTuple iterator class.
var ArrayTupleIteratorClass *Class

// Elk's ArrayTuple value
type ArrayTuple []Value

func ArrayTupleConstructor(class *Class) Value {
	return Ref(&ArrayTuple{})
}

func NewArrayTuple(capacity int) *ArrayTuple {
	l := make(ArrayTuple, 0, capacity)
	return &l
}

func NewArrayTupleWithLength(length int) *ArrayTuple {
	l := make(ArrayTuple, length)
	return &l
}

func NewArrayTupleWithElements(capacity int, elements ...Value) *ArrayTuple {
	l := make(ArrayTuple, len(elements), len(elements)+capacity)
	copy(l, elements)
	return &l
}

func (*ArrayTuple) Class() *Class {
	return ArrayTupleClass
}

func (*ArrayTuple) DirectClass() *Class {
	return ArrayTupleClass
}

func (*ArrayTuple) SingletonClass() *Class {
	return nil
}

func (t *ArrayTuple) Copy() Reference {
	return t
}

func (t *ArrayTuple) Error() string {
	return t.Inspect()
}

// Add a new element.
func (t *ArrayTuple) Append(element Value) {
	*t = append(*t, element)
}

const MAX_ARRAY_TUPLE_ELEMENTS_IN_INSPECT = 300

func (t *ArrayTuple) Inspect() string {
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

func (*ArrayTuple) InstanceVariables() *InstanceVariables {
	return nil
}

// Get an element under the given index.
func (t *ArrayTuple) Get(index int) (Value, Value) {
	return GetFromSlice((*[]Value)(t), index)
}

// Get an element under the given index without bounds checking.
func (t *ArrayTuple) At(i int) Value {
	return (*t)[i]
}

// Get an element under the given index.
func (t *ArrayTuple) Subscript(key Value) (Value, Value) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Undefined, Ref(NewIndexOutOfRangeError(key.Inspect(), len(*t)))
		}
		return Undefined, Ref(NewCoerceError(IntClass, key.Class()))
	}

	return t.Get(i)
}

// Set an element under the given index.
func (t *ArrayTuple) Set(index int, val Value) Value {
	return SetInSlice((*[]Value)(t), index, val)
}

// Set an element under the given index without bounds checking.
func (t *ArrayTuple) SetAt(index int, val Value) {
	(*t)[index] = val
}

// Set an element under the given index.
func (t *ArrayTuple) SubscriptSet(key, val Value) Value {
	length := len(*t)
	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Ref(NewIndexOutOfRangeError(key.Inspect(), length))
		}
		return Ref(NewCoerceError(IntClass, key.Class()))
	}

	return t.Set(i, val)
}

// Concatenate another value with this arrayTuple, creating a new value, and return the result.
// If the operation is illegal an error will be returned.
func (t *ArrayTuple) ConcatVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *ArrayList:
			newList := make(ArrayList, len(*t), len(*t)+len(*o))
			copy(newList, *t)
			newList = append(newList, *o...)
			return Ref(&newList), Undefined
		case *ArrayTuple:
			newArrayTuple := make(ArrayTuple, len(*t), len(*t)+len(*o))
			copy(newArrayTuple, *t)
			newArrayTuple = append(newArrayTuple, *o...)
			return Ref(&newArrayTuple), Undefined
		}
	}

	return Undefined, Ref(Errorf(TypeErrorClass, "cannot concat %s with arrayTuple %s", other.Inspect(), t.Inspect()))
}

// Repeat the content of this arrayTuple n times and return a new arrayTuple containing the result.
// If the operation is illegal an error will be returned.
func (t *ArrayTuple) Repeat(other Value) (*ArrayTuple, Value) {
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
		newArrayTuple := make(ArrayTuple, 0, newLen)
		for i := 0; i < int(o); i++ {
			newArrayTuple = append(newArrayTuple, *t...)
		}
		return &newArrayTuple, Undefined
	default:
		return nil, Ref(Errorf(TypeErrorClass, "cannot repeat a arrayTuple using %s", other.Inspect()))
	}
}

// Return a box pointing to the slot with the given index.
func (l *ArrayTuple) ImmutableBoxOfVal(index Value) (*ImmutableBox, Value) {
	var i int

	i, ok := ToGoInt(index)
	if !ok {
		if i == -1 {
			return nil, Ref(NewIndexOutOfRangeError(index.Inspect(), len(*l)))
		}
		return nil, Ref(NewCoerceError(IntClass, index.Class()))
	}

	return l.ImmutableBoxOf(i)
}

// Return a box pointing to the slot with the given index.
func (l *ArrayTuple) ImmutableBoxOf(index int) (*ImmutableBox, Value) {
	index, err := NormalizeArrayIndex(index, l.Length())
	if !err.IsUndefined() {
		return nil, err
	}

	box := (*ImmutableBox)(&(*l)[index])
	return box, Undefined
}

// Expands the arrayTuple by n nil elements.
func (t *ArrayTuple) Expand(newElements int) {
	if newElements < 1 {
		return
	}

	newCollection := slices.Grow(*t, newElements)
	for i := 0; i < newElements; i++ {
		newCollection = append(newCollection, Nil)
	}
	*t = newCollection
}

func (t *ArrayTuple) Length() int {
	return len(*t)
}

type ArrayTupleIterator struct {
	ArrayTuple *ArrayTuple
	Index      int
}

func NewArrayTupleIterator(arrayTuple *ArrayTuple) *ArrayTupleIterator {
	return &ArrayTupleIterator{
		ArrayTuple: arrayTuple,
	}
}

func NewArrayTupleIteratorWithIndex(arrayTuple *ArrayTuple, index int) *ArrayTupleIterator {
	return &ArrayTupleIterator{
		ArrayTuple: arrayTuple,
		Index:      index,
	}
}

func (*ArrayTupleIterator) Class() *Class {
	return ArrayTupleIteratorClass
}

func (*ArrayTupleIterator) DirectClass() *Class {
	return ArrayTupleIteratorClass
}

func (*ArrayTupleIterator) SingletonClass() *Class {
	return nil
}

func (t *ArrayTupleIterator) Copy() Reference {
	return &ArrayTupleIterator{
		ArrayTuple: t.ArrayTuple,
		Index:      t.Index,
	}
}

func (t *ArrayTupleIterator) Inspect() string {
	return fmt.Sprintf("Std::ArrayTuple::Iterator{&: %p, tuple: %s, index: %d}", t, t.ArrayTuple.Inspect(), t.Index)
}

func (t *ArrayTupleIterator) Error() string {
	return t.Inspect()
}

func (*ArrayTupleIterator) InstanceVariables() *InstanceVariables {
	return nil
}

func (t *ArrayTupleIterator) Next() (Value, Value) {
	if t.Index >= t.ArrayTuple.Length() {
		return Undefined, stopIterationSymbol.ToValue()
	}

	next := (*t.ArrayTuple)[t.Index]
	t.Index++
	return next, Undefined
}

func (t *ArrayTupleIterator) Reset() {
	t.Index = 0
}

func initArrayTuple() {
	ArrayTupleClass = NewClassWithOptions(ClassWithConstructor(ArrayTupleConstructor))
	ArrayTupleClass.IncludeMixin(TupleMixin)
	StdModule.AddConstantString("ArrayTuple", Ref(ArrayTupleClass))

	ArrayTupleIteratorClass = NewClass()
	ArrayTupleClass.AddConstantString("Iterator", Ref(ArrayTupleIteratorClass))
}
