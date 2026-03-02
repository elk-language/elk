package value

import (
	"fmt"
	"iter"
	"strings"

	"github.com/elk-language/elk/indent"
)

// Elk's ArrayListOfValue value
type ArrayListOfValue []Value

var _ ArrayList = &ArrayListOfValue{}

func NormalizeArrayIndex(index, length int) (int, Value) {
	if index >= length || index < -length {
		return 0, Ref(NewIndexOutOfRangeError(fmt.Sprint(index), length))
	}

	if index < 0 {
		index = length + index
	}

	return index, Undefined
}

func ArrayListOfValueConstructor(class *Class) Value {
	return Ref(&ArrayListOfValue{})
}

func NewArrayListOfValue(capacity int) *ArrayListOfValue {
	l := make(ArrayListOfValue, 0, capacity)
	return &l
}

func NewArrayListOfValueWithLength(length int) *ArrayListOfValue {
	l := make(ArrayListOfValue, length)
	return &l
}

func NewArrayListOfValueWithElements(capacity int, elements ...Value) *ArrayListOfValue {
	l := make(ArrayListOfValue, len(elements), len(elements)+capacity)
	copy(l, elements)
	return &l
}

func (l *ArrayListOfValue) NewArrayList(capacity int) ArrayList {
	return NewArrayListOfValue(capacity)
}

func (l *ArrayListOfValue) CloneArrayList(capacity int) ArrayList {
	newList := NewArrayListOfValue(capacity)
	newList.Append(*l...)
	return newList
}

func (l *ArrayListOfValue) SliceArrayList(from, to int) ArrayList {
	n := (*l)[from:to]
	return &n
}

func (l *ArrayListOfValue) SliceArrayTuple(from, to int) ArrayTuple {
	return l.SliceArrayList(from, to)
}

func (l *ArrayListOfValue) NewArrayTuple(capacity int) ArrayTuple {
	return l.NewArrayList(capacity)
}

func (l *ArrayListOfValue) CloneArrayTuple(capacity int) ArrayTuple {
	return l.CloneArrayList(capacity)
}

func (*ArrayListOfValue) Class() *Class {
	return ArrayListClass
}

func (*ArrayListOfValue) DirectClass() *Class {
	return ArrayListClass
}

func (*ArrayListOfValue) SingletonClass() *Class {
	return nil
}

func (l *ArrayListOfValue) Copy() Reference {
	if l == nil {
		return l
	}

	newList := make(ArrayListOfValue, len(*l))
	copy(newList, *l)
	return &newList
}

func (l *ArrayListOfValue) Elements() iter.Seq2[int, Value] {
	return func(yield func(int, Value) bool) {
		for i, element := range *l {
			if !yield(i, element) {
				return
			}
		}
	}
}

func (l *ArrayListOfValue) Iterate() iter.Seq2[Value, Value] {
	return func(yield func(Value, Value) bool) {
		for _, element := range *l {
			if !yield(element, Undefined) {
				return
			}
		}
	}
}

func (l *ArrayListOfValue) ToValue() Value {
	return Ref(l)
}

func (l *ArrayListOfValue) Error() string {
	return l.Inspect()
}

func (l *ArrayListOfValue) Iter() NativeIterator {
	return l.IterNative()
}

func (l *ArrayListOfValue) IterNative() *ArrayListOfValueIterator {
	return NewArrayListOfValueIterator(l)
}

func (l *ArrayListOfValue) IterTuple() ArrayTupleIterator {
	return l.IterNative()
}

func (l *ArrayListOfValue) IterList() ArrayListIterator {
	return l.IterNative()
}

const MAX_ARRAY_LIST_ELEMENTS_IN_INSPECT = 300

func (l *ArrayListOfValue) Inspect() string {
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

func (*ArrayListOfValue) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *ArrayListOfValue) Capacity() int {
	return cap(*l)
}

func (l *ArrayListOfValue) RemoveAtErr(index int) Value {
	index, err := NormalizeArrayIndex(index, l.Length())
	if !err.IsUndefined() {
		return err
	}

	l.RemoveAt(index)
	return Undefined
}

func (l *ArrayListOfValue) RemoveAt(i int) {
	s := *l
	copy(s[i:], s[i+1:])
	*l = s[:len(s)-1]
}

func (l *ArrayListOfValue) LeftCapacity() int {
	return l.Capacity() - l.Length()
}

func (l *ArrayListOfValue) Length() int {
	return len(*l)
}

// Add new elements.
func (l *ArrayListOfValue) Append(elements ...Value) {
	*l = append(*l, elements...)
}

func (l *ArrayListOfValue) AppendVal(elements ...Value) Value {
	l.Append(elements...)
	return Undefined
}

// Expand the array list to have
// empty slots for new elements.
func (l *ArrayListOfValue) Grow(newSlots int) {
	newList := make(ArrayListOfValue, l.Length(), l.Capacity()+newSlots)
	copy(newList, *l)
	*l = newList
}

// Get an element under the given index.
func GetFromSlice[V any](collection *[]V, index int) (ret V, err Value) {
	l := len(*collection)
	index, err = NormalizeArrayIndex(index, l)
	if err.IsNotUndefined() {
		return ret, err
	}

	return (*collection)[index], Undefined
}

// Set an element under the given index.
func SetInSlice[V any](collection *[]V, index int, val V) (err Value) {
	index, err = NormalizeArrayIndex(index, len(*collection))
	if !err.IsUndefined() {
		return err
	}

	(*collection)[index] = val
	return Undefined
}

// Get an element under the given index.
func (l *ArrayListOfValue) Get(index int) (Value, Value) {
	return GetFromSlice((*[]Value)(l), index)
}

// Get an element under the given index without bounds checking
func (l *ArrayListOfValue) At(i int) Value {
	return (*l)[i]
}

func (l *ArrayListOfValue) AtVal(i int) Value {
	return l.At(i)
}

// Get an element under the given index.
func (l *ArrayListOfValue) Subscript(key Value) (Value, Value) {
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
func (l *ArrayListOfValue) Set(index int, val Value) Value {
	return SetInSlice((*[]Value)(l), index, val)
}

// Set an element under the given index without bounds checking.
func (l *ArrayListOfValue) SetAt(index int, val Value) {
	(*l)[index] = val
}

func (l *ArrayListOfValue) SetAtVal(index int, val Value) Value {
	l.SetAtVal(index, val)
	return Undefined
}

// Set an element under the given index.
func (l *ArrayListOfValue) SubscriptSet(key, val Value) Value {
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
func (l *ArrayListOfValue) Concat(other Value) (*ArrayListOfValue, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *ArrayListOfValue:
			newList := make(ArrayListOfValue, len(*l), len(*l)+len(*o))
			copy(newList, *l)
			newList = append(newList, *o...)
			return &newList, Undefined
		case *ArrayTupleOfValue:
			newList := make(ArrayListOfValue, len(*l), len(*l)+len(*o))
			copy(newList, *l)
			newList = append(newList, *o...)
			return &newList, Undefined
		case ArrayTuple:
			newList := make(ArrayListOfValue, len(*l), len(*l)+o.Length())
			copy(newList, *l)

			for i, element := range o.Elements() {
				newList[len(*l)+i] = element
			}

			return &newList, Undefined
		}
	}

	return nil, Ref(Errorf(TypeErrorClass, "cannot concat %s with list %s", other.Inspect(), l.Inspect()))
}

func (l *ArrayListOfValue) ConcatVal(other Value) (Value, Value) {
	return RefErr(l.Concat(other))
}

// Repeat the content of this list n times and return a new list containing the result.
// If the operation is illegal an error will be returned.
func (l *ArrayListOfValue) Repeat(other Value) (*ArrayListOfValue, Value) {
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
		newList := make(ArrayListOfValue, 0, newLen)
		for range int(o) {
			newList = append(newList, *l...)
		}
		return &newList, Undefined
	default:
		return nil, Ref(Errorf(TypeErrorClass, "cannot repeat a list using %s", other.Inspect()))
	}
}

func (l *ArrayListOfValue) RepeatVal(other Value) (Value, Value) {
	return RefErr(l.Repeat(other))
}

// Return an immutable box pointing to the slot with the given index.
func (l *ArrayListOfValue) ImmutableBoxOfVal(index Value) (Value, Value) {
	b, err := l.boxOf(index)
	return Ref(b.ToImmutableBox()), err
}

// Return a box pointing to the slot with the given index.
func (l *ArrayListOfValue) BoxOfVal(index Value) (Value, Value) {
	return RefErr(l.boxOf(index))
}

func (l *ArrayListOfValue) boxOf(index Value) (*BoxOfValue, Value) {
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
func (l *ArrayListOfValue) BoxOf(index int) (*BoxOfValue, Value) {
	index, err := NormalizeArrayIndex(index, l.Length())
	if !err.IsUndefined() {
		return nil, err
	}

	box := (*BoxOfValue)(&(*l)[index])
	return box, Undefined
}

// Expands the list by n nil elements.
func (l *ArrayListOfValue) Expand(newElements int) {
	if newElements < 1 {
		return
	}

	newCollection := make(ArrayListOfValue, len(*l), cap(*l)+newElements)
	copy(newCollection, *l)
	for i := 0; i < newElements; i++ {
		newCollection = append(newCollection, Nil)
	}
	*l = newCollection
}

func (l *ArrayListOfValue) AppendAt(key Value, val Value) Value {
	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Ref(NewIndexOutOfRangeError(key.Inspect(), l.Length()))
		}
		return Ref(NewCoerceError(IntClass, key.Class()))
	}

	return l.AppendAtInt(i, val)
}

func (l *ArrayListOfValue) AppendAtInt(index int, val Value) Value {
	length := l.Length()

	if index < 0 {
		return Ref(NewNegativeIndicesInCollectionLiteralsError(fmt.Sprint(index)))
	}

	if index >= length {
		newElementsCount := (index + 1) - length
		l.Expand(newElementsCount)
	}

	(*l)[index] = val
	return Undefined
}

type ArrayListOfValueIterator struct {
	ArrayList *ArrayListOfValue
	Index     int
}

var _ ArrayListIterator = &ArrayListOfValueIterator{}

func NewArrayListOfValueIterator(list *ArrayListOfValue) *ArrayListOfValueIterator {
	return &ArrayListOfValueIterator{
		ArrayList: list,
	}
}

func NewArrayListIteratorWithIndex(list *ArrayListOfValue, index int) *ArrayListOfValueIterator {
	return &ArrayListOfValueIterator{
		ArrayList: list,
		Index:     index,
	}
}

func (*ArrayListOfValueIterator) Class() *Class {
	return ArrayListIteratorClass
}

func (*ArrayListOfValueIterator) DirectClass() *Class {
	return ArrayListIteratorClass
}

func (*ArrayListOfValueIterator) SingletonClass() *Class {
	return nil
}

func (l *ArrayListOfValueIterator) Copy() Reference {
	return &ArrayListOfValueIterator{
		ArrayList: l.ArrayList,
		Index:     l.Index,
	}
}

func (i *ArrayListOfValueIterator) ToValue() Value {
	return Ref(i)
}

func (l *ArrayListOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::ArrayList::Iterator{&: %p, list: %s, index: %d}", l, l.ArrayList.Inspect(), l.Index)
}

func (l *ArrayListOfValueIterator) Error() string {
	return l.Inspect()
}

func (*ArrayListOfValueIterator) InstanceVariables() *InstanceVariables {
	return nil
}

var stopIterationSymbol = ToSymbol("stop_iteration")

func (l *ArrayListOfValueIterator) NextValue() (Value, Value) {
	if l.Index >= l.ArrayList.Length() {
		return Undefined, stopIterationSymbol.ToValue()
	}

	next := (*l.ArrayList)[l.Index]
	l.Index++
	return next, Undefined
}

func (l *ArrayListOfValueIterator) Elements() iter.Seq[Value] {
	return func(yield func(Value) bool) {
		for ; l.Index >= l.ArrayList.Length(); l.Index++ {
			if !yield((*l.ArrayList)[l.Index]) {
				return
			}
		}
	}
}

func (l *ArrayListOfValueIterator) Reset() {
	l.Index = 0
}
