package value

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/elk-language/elk/indent"
)

// Elk's ArrayTupleOfValue value
type ArrayTupleOfValue []Value

var _ ArrayTuple = &ArrayTupleOfValue{}

func ArrayTupleOfValueConstructor(class *Class) Value {
	return Ref(&ArrayTupleOfValue{})
}

func NewArrayTupleOfValue(capacity int) *ArrayTupleOfValue {
	l := make(ArrayTupleOfValue, 0, capacity)
	return &l
}

func NewArrayTupleOfValueWithLength(length int) *ArrayTupleOfValue {
	l := make(ArrayTupleOfValue, length)
	return &l
}

func NewArrayTupleOfValueWithElements(capacity int, elements ...Value) *ArrayTupleOfValue {
	l := make(ArrayTupleOfValue, len(elements), len(elements)+capacity)
	copy(l, elements)
	return &l
}

func (t *ArrayTupleOfValue) IterNative() *ArrayTupleOfValueIterator {
	return NewArrayTupleOfValueIterator(t)
}

func (l *ArrayTupleOfValue) Iter() NativeIterator {
	return l.IterNative()
}

func (l *ArrayTupleOfValue) IterTuple() ArrayTupleIterator {
	return l.IterNative()
}

func (*ArrayTupleOfValue) Class() *Class {
	return ArrayTupleClass
}

func (*ArrayTupleOfValue) DirectClass() *Class {
	return ArrayTupleClass
}

func (*ArrayTupleOfValue) SingletonClass() *Class {
	return nil
}

func (t *ArrayTupleOfValue) Copy() Reference {
	return t
}

func (l *ArrayTupleOfValue) Elements() iter.Seq2[int, Value] {
	return func(yield func(int, Value) bool) {
		for i, element := range *l {
			if !yield(i, element) {
				return
			}
		}
	}
}

func (l *ArrayTupleOfValue) Iterate() iter.Seq2[Value, Value] {
	return func(yield func(Value, Value) bool) {
		for _, element := range *l {
			if !yield(element, Undefined) {
				return
			}
		}
	}
}

func (t *ArrayTupleOfValue) ToValue() Value {
	return Ref(t)
}

func (t *ArrayTupleOfValue) Error() string {
	return t.Inspect()
}

func (t *ArrayTupleOfValue) NewArrayTuple(capacity int) ArrayTuple {
	return NewArrayTupleOfValue(capacity)
}

func (t *ArrayTupleOfValue) CloneArrayTuple(capacity int) ArrayTuple {
	newTuple := NewArrayTupleOfValue(capacity)
	newTuple.Append(*t...)
	return newTuple
}

func (t *ArrayTupleOfValue) SliceArrayTuple(from, to int) ArrayTuple {
	n := (*t)[from:to]
	return &n
}

// Add a new element.
func (t *ArrayTupleOfValue) AppendVal(elements ...Value) Value {
	*t = append(*t, elements...)
	return Undefined
}

func (t *ArrayTupleOfValue) Append(elements ...Value) {
	*t = append(*t, elements...)
}

const MAX_ARRAY_TUPLE_ELEMENTS_IN_INSPECT = 300

func (t *ArrayTupleOfValue) Inspect() string {
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

func (*ArrayTupleOfValue) InstanceVariables() *InstanceVariables {
	return nil
}

// Get an element under the given index.
func (t *ArrayTupleOfValue) Get(index int) (Value, Value) {
	return GetFromSlice((*[]Value)(t), index)
}

// Get an element under the given index without bounds checking.
func (t *ArrayTupleOfValue) At(i int) Value {
	return (*t)[i]
}

func (t *ArrayTupleOfValue) AtVal(i int) Value {
	return t.At(i)
}

// Get an element under the given index.
func (t *ArrayTupleOfValue) Subscript(key Value) (Value, Value) {
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
func (t *ArrayTupleOfValue) Set(index int, val Value) Value {
	return SetInSlice((*[]Value)(t), index, val)
}

// Set an element under the given index without bounds checking.
func (t *ArrayTupleOfValue) SetAt(index int, val Value) {
	(*t)[index] = val
}

// Set an element under the given index.
func (t *ArrayTupleOfValue) SubscriptSet(key, val Value) Value {
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
func (t *ArrayTupleOfValue) ConcatVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *ArrayListOfValue:
			newList := make(ArrayListOfValue, len(*t), len(*t)+len(*o))
			copy(newList, *t)
			newList = append(newList, *o...)
			return Ref(&newList), Undefined
		case *ArrayTupleOfValue:
			newArrayTuple := make(ArrayTupleOfValue, len(*t), len(*t)+len(*o))
			copy(newArrayTuple, *t)
			newArrayTuple = append(newArrayTuple, *o...)
			return Ref(&newArrayTuple), Undefined
		case ArrayList:
			newArrayTuple := make(ArrayListOfValue, len(*t), len(*t)+o.Length())
			copy(newArrayTuple, *t)

			for i, element := range o.Elements() {
				newArrayTuple[len(*t)+i] = element
			}

			return Ref(&newArrayTuple), Undefined
		}
	}

	return Undefined, Ref(Errorf(TypeErrorClass, "cannot concat %s with arrayTuple %s", other.Inspect(), t.Inspect()))
}

// Repeat the content of this arrayTuple n times and return a new arrayTuple containing the result.
// If the operation is illegal an error will be returned.
func (t *ArrayTupleOfValue) Repeat(other Value) (*ArrayTupleOfValue, Value) {
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
		newArrayTuple := make(ArrayTupleOfValue, 0, newLen)
		for i := 0; i < int(o); i++ {
			newArrayTuple = append(newArrayTuple, *t...)
		}
		return &newArrayTuple, Undefined
	default:
		return nil, Ref(Errorf(TypeErrorClass, "cannot repeat a arrayTuple using %s", other.Inspect()))
	}
}

func (t *ArrayTupleOfValue) RepeatVal(other Value) (Value, Value) {
	return RefErr(t.Repeat(other))
}

// Return a box pointing to the slot with the given index.
func (l *ArrayTupleOfValue) ImmutableBoxOfVal(index Value) (Value, Value) {
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
func (l *ArrayTupleOfValue) ImmutableBoxOf(index int) (*ImmutableBoxOfValue, Value) {
	index, err := NormalizeArrayIndex(index, l.Length())
	if !err.IsUndefined() {
		return nil, err
	}

	box := (*ImmutableBoxOfValue)(&(*l)[index])
	return box, Undefined
}

// Expands the arrayTuple by n nil elements.
func (t *ArrayTupleOfValue) Expand(newElements int) {
	if newElements < 1 {
		return
	}

	newCollection := slices.Grow(*t, newElements)
	for range newElements {
		newCollection = append(newCollection, Nil)
	}
	*t = newCollection
}

func (t *ArrayTupleOfValue) AppendAt(key Value, val Value) Value {
	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Ref(NewIndexOutOfRangeError(key.Inspect(), t.Length()))
		}
		return Ref(NewCoerceError(IntClass, key.Class()))
	}

	return t.AppendAtInt(i, val)
}

func (t *ArrayTupleOfValue) AppendAtInt(index int, val Value) Value {
	l := len(*t)

	if index < 0 {
		return Ref(NewNegativeIndicesInCollectionLiteralsError(fmt.Sprint(index)))
	}

	if index >= l {
		newElementsCount := (index + 1) - l
		t.Expand(newElementsCount)
	}

	(*t)[index] = val
	return Undefined
}

func (t *ArrayTupleOfValue) Length() int {
	return len(*t)
}

type ArrayTupleOfValueIterator struct {
	ArrayTuple *ArrayTupleOfValue
	Index      int
}

var _ ArrayTupleIterator = &ArrayTupleOfValueIterator{}

func NewArrayTupleOfValueIterator(arrayTuple *ArrayTupleOfValue) *ArrayTupleOfValueIterator {
	return &ArrayTupleOfValueIterator{
		ArrayTuple: arrayTuple,
	}
}

func NewArrayTupleOfValueIteratorWithIndex(arrayTuple *ArrayTupleOfValue, index int) *ArrayTupleOfValueIterator {
	return &ArrayTupleOfValueIterator{
		ArrayTuple: arrayTuple,
		Index:      index,
	}
}

func (*ArrayTupleOfValueIterator) Class() *Class {
	return ArrayTupleIteratorClass
}

func (*ArrayTupleOfValueIterator) DirectClass() *Class {
	return ArrayTupleIteratorClass
}

func (*ArrayTupleOfValueIterator) SingletonClass() *Class {
	return nil
}

func (t *ArrayTupleOfValueIterator) Copy() Reference {
	return &ArrayTupleOfValueIterator{
		ArrayTuple: t.ArrayTuple,
		Index:      t.Index,
	}
}

func (i *ArrayTupleOfValueIterator) ToValue() Value {
	return Ref(i)
}

func (t *ArrayTupleOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::ArrayTuple::Iterator{&: %p, tuple: %s, index: %d}", t, t.ArrayTuple.Inspect(), t.Index)
}

func (t *ArrayTupleOfValueIterator) Error() string {
	return t.Inspect()
}

func (*ArrayTupleOfValueIterator) InstanceVariables() *InstanceVariables {
	return nil
}

func (t *ArrayTupleOfValueIterator) NextValue() (Value, Value) {
	if t.Index >= t.ArrayTuple.Length() {
		return Undefined, stopIterationSymbol.ToValue()
	}

	next := (*t.ArrayTuple)[t.Index]
	t.Index++
	return next, Undefined
}

func (t *ArrayTupleOfValueIterator) Elements() iter.Seq[Value] {
	return func(yield func(Value) bool) {
		for ; t.Index >= t.ArrayTuple.Length(); t.Index++ {
			if !yield((*t.ArrayTuple)[t.Index]) {
				return
			}
		}
	}
}

func (t *ArrayTupleOfValueIterator) Reset() {
	t.Index = 0
}
