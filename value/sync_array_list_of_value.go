package value

import (
	"fmt"
	"iter"
	"strings"
	"sync"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value/symbol"
)

// Elk's SyncArrayListOfValue value
type SyncArrayListOfValue struct {
	List *ArrayListOfValue
	mu   sync.Mutex
}

var _ SyncArrayList = &SyncArrayListOfValue{}

func SyncArrayListOfValueConstructor(class *Class) Value {
	return Ref(&SyncArrayListOfValue{List: &ArrayListOfValue{}})
}

func NewSyncArrayListOfValue(capacity int) *SyncArrayListOfValue {
	return &SyncArrayListOfValue{
		List: NewArrayListOfValue(capacity),
	}
}

func NewSyncArrayListOfValueWithLength(length int) *SyncArrayListOfValue {
	return &SyncArrayListOfValue{
		List: NewArrayListOfValueWithLength(length),
	}
}

func (l *SyncArrayListOfValue) Lock() {
	l.mu.Lock()
}

func (l *SyncArrayListOfValue) Unlock() {
	l.mu.Unlock()
}

func (l *SyncArrayListOfValue) ArrayList() ArrayList {
	return l.List
}

func (l *SyncArrayListOfValue) NewSyncArrayList(capacity int) SyncArrayList {
	return NewSyncArrayListOfValue(capacity)
}

func (l *SyncArrayListOfValue) CloneSyncArrayList(capacity int) SyncArrayList {
	l.Lock()
	defer l.Unlock()

	return &SyncArrayListOfValue{
		List: l.List.CloneArrayList(capacity).(*ArrayListOfValue),
	}
}

func (l *SyncArrayListOfValue) SliceSyncArrayList(from, to int) SyncArrayList {
	l.Lock()
	defer l.Unlock()

	return &SyncArrayListOfValue{
		List: l.List.SliceArrayList(from, to).(*ArrayListOfValue),
	}
}

func (l *SyncArrayListOfValue) ToArrayList() ArrayList {
	l.Lock()
	defer l.Unlock()

	return l.List.Clone()
}

func (l *SyncArrayListOfValue) ToArrayTuple() ArrayTuple {
	l.Lock()
	defer l.Unlock()

	return l.List.ToArrayTuple()
}

func (*SyncArrayListOfValue) Class() *Class {
	return SyncArrayListClass
}

func (*SyncArrayListOfValue) DirectClass() *Class {
	return SyncArrayListClass
}

func (*SyncArrayListOfValue) SingletonClass() *Class {
	return nil
}

func (l *SyncArrayListOfValue) Copy() Reference {
	return &SyncArrayListOfValue{
		List: l.List.Clone(),
	}
}

func (l *SyncArrayListOfValue) Elements() iter.Seq2[int, Value] {
	return func(yield func(int, Value) bool) {
		l.Lock()
		defer l.Unlock()

		for i, element := range l.List.Elements() {
			if !yield(i, element) {
				return
			}
		}
	}
}

func (l *SyncArrayListOfValue) Iterate() iter.Seq2[Value, Value] {
	return func(yield func(Value, Value) bool) {
		l.Lock()
		defer l.Unlock()

		for element, err := range l.List.Iterate() {
			if !yield(element, err) {
				return
			}
		}
	}
}

func (l *SyncArrayListOfValue) ToValue() Value {
	return Ref(l)
}

func (l *SyncArrayListOfValue) Error() string {
	return l.Inspect()
}

func (l *SyncArrayListOfValue) Iter() NativeIterator {
	return l.IterNative()
}

func (l *SyncArrayListOfValue) IterNative() *SyncArrayListOfValueIterator {
	return NewSyncArrayListOfValueIterator(l)
}

func (l *SyncArrayListOfValue) IterSyncList() SyncArrayListIterator {
	return l.IterNative()
}

func (l *SyncArrayListOfValue) Inspect() string {
	var hasMultilineElements bool
	elementStrings := make(
		[]string,
		0,
		min(MAX_ARRAY_LIST_ELEMENTS_IN_INSPECT, l.Length()),
	)

	for i, element := range *l.List {
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

	buff.WriteString("Std::Sync::ArrayList[")
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

func (*SyncArrayListOfValue) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *SyncArrayListOfValue) Capacity() int {
	return l.List.Capacity()
}

func (l *SyncArrayListOfValue) RemoveAtRange(i, j int) {
	l.Lock()
	defer l.Unlock()

	l.List.RemoveAtRange(i, j)
}

func (l *SyncArrayListOfValue) RemoveAtRangeErr(i, j int) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.RemoveAtRangeErr(i, j)
}

func (l *SyncArrayListOfValue) RemoveAt(i int) {
	l.Lock()
	defer l.Unlock()

	l.List.RemoveAt(i)
}

func (l *SyncArrayListOfValue) RemoveAtVal(indexVal Value) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.RemoveAtVal(indexVal)
}

func (l *SyncArrayListOfValue) LeftCapacity() int {
	l.Lock()
	defer l.Unlock()

	return l.List.LeftCapacity()
}

func (l *SyncArrayListOfValue) Length() int {
	return l.List.Length()
}

func (l *SyncArrayListOfValue) Append(elements ...Value) {
	l.Lock()
	defer l.Unlock()

	l.List.Append(elements...)
}

func (l *SyncArrayListOfValue) AppendVal(elements ...Value) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.AppendVal(elements...)
}

func (l *SyncArrayListOfValue) Grow(newSlots int) {
	l.Lock()
	defer l.Unlock()

	l.List.Grow(newSlots)
}

func (l *SyncArrayListOfValue) Get(index int) (Value, Value) {
	l.Lock()
	defer l.Unlock()

	return l.List.Get(index)
}

func (l *SyncArrayListOfValue) GetOk(index int) (Value, bool) {
	l.Lock()
	defer l.Unlock()

	return l.List.GetOk(index)
}

func (l *SyncArrayListOfValue) At(i int) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.At(i)
}

func (l *SyncArrayListOfValue) AtVal(i int) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.AtVal(i)
}

func (l *SyncArrayListOfValue) Subscript(key Value) (Value, Value) {
	l.Lock()
	defer l.Unlock()

	return l.List.Subscript(key)
}

func (l *SyncArrayListOfValue) SubscriptInt(key int) (Value, Value) {
	l.Lock()
	defer l.Unlock()

	return l.List.SubscriptInt(key)
}

func (l *SyncArrayListOfValue) TryAtVal(key Value) (Value, Value) {
	l.Lock()
	defer l.Unlock()

	return l.List.TryAtVal(key)
}

func (l *SyncArrayListOfValue) Set(index int, val Value) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.Set(index, val)
}

func (l *SyncArrayListOfValue) SetAt(index int, val Value) {
	l.Lock()
	defer l.Unlock()

	l.List.SetAt(index, val)
}

func (l *SyncArrayListOfValue) SetAtVal(index int, val Value) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.SetAtVal(index, val)
}

func (l *SyncArrayListOfValue) SubscriptSet(key, val Value) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.SubscriptSet(key, val)
}

func (l *SyncArrayListOfValue) SubscriptSetInt(key int, val Value) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.SubscriptSetInt(key, val)
}

func (l *SyncArrayListOfValue) ConcatOk(other Value) (*SyncArrayListOfValue, bool) {
	switch o := other.SafeAsReference().(type) {
	case *SyncArrayListOfValue:
		l.Lock()
		defer l.Unlock()

		o.Lock()
		defer o.Unlock()

		llen := l.List.Length()
		newList := make(ArrayListOfValue, llen, llen+o.Length())
		copy(newList, *l.List)
		newList = append(newList, *o.List...)
		result := &SyncArrayListOfValue{
			List: &newList,
		}
		return result, true
	case *ArrayListOfValue:
		l.Lock()
		defer l.Unlock()

		llen := l.List.Length()
		newList := make(ArrayListOfValue, llen, llen+len(*o))
		copy(newList, *l.List)
		newList = append(newList, *o...)
		result := &SyncArrayListOfValue{
			List: &newList,
		}
		return result, true
	case *ArrayTupleOfValue:
		l.Lock()
		defer l.Unlock()

		llen := l.List.Length()
		newList := make(ArrayListOfValue, llen, llen+len(*o))
		copy(newList, *l.List)
		newList = append(newList, *o...)
		result := &SyncArrayListOfValue{
			List: &newList,
		}
		return result, true
	case Tuple:
		llen := l.List.Length()
		newList := make(ArrayListOfValue, llen, llen+o.Length())
		copy(newList, *l.List)

		for i, element := range o.Elements() {
			newList[llen+i] = element
		}

		result := &SyncArrayListOfValue{
			List: &newList,
		}
		return result, true
	default:
		return nil, false
	}
}

func (l *SyncArrayListOfValue) Concat(other Value) (*SyncArrayListOfValue, Value) {
	result, ok := l.ConcatOk(other)
	if !ok {
		return nil, Ref(Errorf(TypeErrorClass, "cannot concat %s with list %s", other.Inspect(), l.Inspect()))
	}

	return result, Undefined
}

func (l *SyncArrayListOfValue) ConcatVal(other Value) (Value, Value) {
	return RefErr(l.Concat(other))
}

func (l *SyncArrayListOfValue) ConcatValOk(other Value) (Value, bool) {
	result, ok := l.ConcatOk(other)
	if !ok {
		return Undefined, false
	}

	return result.ToValue(), true
}

func (l *SyncArrayListOfValue) Repeat(other Value) (*SyncArrayListOfValue, Value) {
	l.Lock()
	defer l.Unlock()

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
		newLen, ok := o.MultiplyOverflow(SmallInt(l.Length()))
		if !ok {
			return nil, Ref(Errorf(
				OutOfRangeErrorClass,
				"list repeat count is too large %s",
				o.Inspect(),
			))
		}
		newList := make(ArrayListOfValue, 0, newLen)
		for range int(o) {
			newList = append(newList, *l.List...)
		}
		result := &SyncArrayListOfValue{
			List: &newList,
		}
		return result, Undefined
	default:
		return nil, Ref(Errorf(TypeErrorClass, "cannot repeat a list using %s", other.Inspect()))
	}
}

func (l *SyncArrayListOfValue) RepeatVal(other Value) (Value, Value) {
	return l.RepeatVal(other)
}

func (l *SyncArrayListOfValue) Expand(newElements int) {
	l.Lock()
	defer l.Unlock()

	l.List.Expand(newElements)
}

func (l *SyncArrayListOfValue) AppendAt(key Value, val Value) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.AppendAt(key, val)
}

func (l *SyncArrayListOfValue) AppendAtInt(index int, val Value) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.AppendAtInt(index, val)
}

type SyncArrayListOfValueIterator struct {
	ArrayList *SyncArrayListOfValue
	Index     int
}

var _ SyncArrayListIterator = &SyncArrayListOfValueIterator{}

func NewSyncArrayListOfValueIterator(list *SyncArrayListOfValue) *SyncArrayListOfValueIterator {
	return &SyncArrayListOfValueIterator{
		ArrayList: list,
	}
}

func NewSyncArrayListIteratorWithIndex(list *SyncArrayListOfValue, index int) *SyncArrayListOfValueIterator {
	return &SyncArrayListOfValueIterator{
		ArrayList: list,
		Index:     index,
	}
}

func (*SyncArrayListOfValueIterator) Class() *Class {
	return SyncArrayListIteratorClass
}

func (*SyncArrayListOfValueIterator) DirectClass() *Class {
	return SyncArrayListIteratorClass
}

func (*SyncArrayListOfValueIterator) SingletonClass() *Class {
	return nil
}

func (l *SyncArrayListOfValueIterator) Copy() Reference {
	return &SyncArrayListOfValueIterator{
		ArrayList: l.ArrayList,
		Index:     l.Index,
	}
}

func (i *SyncArrayListOfValueIterator) ToValue() Value {
	return Ref(i)
}

func (l *SyncArrayListOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::Sync::ArrayList::Iterator{&: %p, list: %s, index: %d}", l, l.ArrayList.Inspect(), l.Index)
}

func (l *SyncArrayListOfValueIterator) Error() string {
	return l.Inspect()
}

func (*SyncArrayListOfValueIterator) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *SyncArrayListOfValueIterator) NextValue() (Value, Value) {
	l.ArrayList.Lock()
	defer l.ArrayList.Unlock()

	if l.Index >= l.ArrayList.List.Length() {
		return Undefined, S(symbol.L_stop_iteration).ToValue()
	}

	next := (*l.ArrayList.List)[l.Index]
	l.Index++
	return next, Undefined
}

func (l *SyncArrayListOfValueIterator) Elements() iter.Seq[Value] {
	return func(yield func(Value) bool) {
		l.ArrayList.Lock()
		defer l.ArrayList.Unlock()

		for ; l.Index < l.ArrayList.List.Length(); l.Index++ {
			if !yield((*l.ArrayList.List)[l.Index]) {
				return
			}
		}
	}
}

func (l *SyncArrayListOfValueIterator) Reset() {
	l.Index = 0
}
