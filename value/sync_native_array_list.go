package value

import (
	"fmt"
	"iter"
	"strings"
	"sync"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value/symbol"
)

// Elk's native array list implementation
type SyncNativeArrayList[T ValueInterface] struct {
	List *NativeArrayList[T]
	mu   sync.Mutex
}

var _ SyncArrayList = &SyncNativeArrayList[String]{}

func NewSyncNativeArrayList[T ValueInterface](capacity int) *SyncNativeArrayList[T] {
	return &SyncNativeArrayList[T]{
		List: NewNativeArrayList[T](capacity),
	}
}

func (l *SyncNativeArrayList[T]) Lock() {
	l.mu.Lock()
}

func (l *SyncNativeArrayList[T]) Unlock() {
	l.mu.Unlock()
}

func (l *SyncNativeArrayList[T]) ArrayList() ArrayList {
	return l.List
}

func (l *SyncNativeArrayList[T]) NewSyncArrayList(capacity int) SyncArrayList {
	return NewSyncNativeArrayList[T](capacity)
}

func (l *SyncNativeArrayList[T]) CloneSyncArrayList(capacity int) SyncArrayList {
	l.Lock()
	defer l.Unlock()

	return &SyncNativeArrayList[T]{
		List: l.List.CloneArrayList(capacity).(*NativeArrayList[T]),
	}
}

func (l *SyncNativeArrayList[T]) SliceSyncArrayList(from, to int) SyncArrayList {
	l.Lock()
	defer l.Unlock()

	return &SyncNativeArrayList[T]{
		List: l.List.Slice(from, to),
	}
}

func (l *SyncNativeArrayList[T]) Elements() iter.Seq2[int, Value] {
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

func (l *SyncNativeArrayList[T]) Iterate() iter.Seq2[Value, Value] {
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

func (*SyncNativeArrayList[T]) Class() *Class {
	return SyncArrayListClass
}

func (*SyncNativeArrayList[T]) DirectClass() *Class {
	return SyncArrayListClass
}

func (l *SyncNativeArrayList[T]) ToValue() Value {
	return Ref(l)
}

func (*SyncNativeArrayList[T]) SingletonClass() *Class {
	return nil
}

func (l *SyncNativeArrayList[T]) Clone() *SyncNativeArrayList[T] {
	l.Lock()
	defer l.Unlock()

	return &SyncNativeArrayList[T]{
		List: l.List.Clone(),
	}
}

func (l *SyncNativeArrayList[T]) Copy() Reference {
	return l.Clone()
}

func (l *SyncNativeArrayList[T]) Error() string {
	return l.Inspect()
}

func (l *SyncNativeArrayList[T]) Inspect() string {
	l.Lock()
	defer l.Unlock()

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

func (*SyncNativeArrayList[T]) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *SyncNativeArrayList[T]) Capacity() int {
	return l.List.Capacity()
}

func (l *SyncNativeArrayList[T]) RemoveAtVal(indexVal Value) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.RemoveAtVal(indexVal)
}

func (l *SyncNativeArrayList[T]) RemoveAt(i int) {
	l.Lock()
	defer l.Unlock()

	l.List.RemoveAt(i)
}

func (l *SyncNativeArrayList[T]) RemoveAtRange(i, j int) {
	l.Lock()
	defer l.Unlock()

	l.List.RemoveAtRange(i, j)
}

func (l *SyncNativeArrayList[T]) RemoveAtRangeErr(i, j int) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.RemoveAtRangeErr(i, j)
}

func (l *SyncNativeArrayList[T]) LeftCapacity() int {
	l.Lock()
	defer l.Unlock()

	return l.List.LeftCapacity()
}

func (l *SyncNativeArrayList[T]) Length() int {
	return l.List.Length()
}

func (l *SyncNativeArrayList[T]) Append(elements ...T) {
	l.Lock()
	defer l.Unlock()

	l.List.Append(elements...)
}

func (l *SyncNativeArrayList[T]) AppendVal(elements ...Value) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.AppendVal(elements...)
}

func (l *SyncNativeArrayList[T]) Grow(newSlots int) {
	l.Lock()
	defer l.Unlock()

	l.List.Grow(newSlots)
}

func (l *SyncNativeArrayList[T]) Get(index int) (T, Value) {
	l.Lock()
	defer l.Unlock()

	return l.List.Get(index)
}

func (l *SyncNativeArrayList[T]) GetOk(index int) (T, bool) {
	l.Lock()
	defer l.Unlock()

	return l.List.GetOk(index)
}

func (l *SyncNativeArrayList[T]) At(i int) T {
	l.Lock()
	defer l.Unlock()

	return l.List.At(i)
}

func (l *SyncNativeArrayList[T]) AtVal(i int) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.AtVal(i)
}

func (l *SyncNativeArrayList[T]) Subscript(key Value) (Value, Value) {
	l.Lock()
	defer l.Unlock()

	return l.List.Subscript(key)
}

func (l *SyncNativeArrayList[T]) SubscriptInt(key int) (Value, Value) {
	l.Lock()
	defer l.Unlock()

	return l.List.SubscriptInt(key)
}

func (l *SyncNativeArrayList[T]) TryAtVal(key Value) (Value, Value) {
	l.Lock()
	defer l.Unlock()

	return l.List.TryAtVal(key)
}

func (l *SyncNativeArrayList[T]) Set(index int, val T) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.Set(index, val)
}

func (l *SyncNativeArrayList[T]) SetAt(index int, val T) {
	l.Lock()
	defer l.Unlock()

	l.List.SetAt(index, val)
}

func (l *SyncNativeArrayList[T]) SetAtVal(index int, val Value) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.SetAtVal(index, val)
}

func (l *SyncNativeArrayList[T]) SubscriptSet(key Value, val Value) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.SubscriptSet(key, val)
}

func (l *SyncNativeArrayList[T]) SubscriptSetInt(key int, val Value) Value {
	l.Lock()
	defer l.Unlock()

	return l.List.SubscriptSetInt(key, val)
}

func (l *SyncNativeArrayList[T]) ConcatValOk(other Value) (Value, bool) {
	switch o := other.SafeAsReference().(type) {
	case *SyncNativeArrayList[T]:
		l.Lock()
		defer l.Unlock()

		o.Lock()
		defer o.Unlock()

		newList := make(NativeArrayList[T], l.Length(), l.Length()+o.Length())
		copy(newList, *l.List)
		newList = append(newList, *o.List...)
		result := &SyncNativeArrayList[T]{
			List: &newList,
		}
		return result.ToValue(), true
	case SyncArrayList:
		l.Lock()
		defer l.Unlock()

		o.Lock()
		defer o.Unlock()

		lLen := l.Length()
		oLen := o.Length()
		newList := make(ArrayListOfValue, lLen, lLen+oLen)
		for i, element := range *l.List {
			newList[i] = element.ToValue()
		}
		for i, element := range o.ArrayList().Elements() {
			newList[i+oLen] = element
		}
		result := &SyncArrayListOfValue{
			List: &newList,
		}
		return result.ToValue(), true
	case Tuple:
		l.Lock()
		defer l.Unlock()

		newList := make(ArrayListOfValue, l.Length(), l.Length()+o.Length())
		for i, element := range *l.List {
			newList[i] = element.ToValue()
		}
		for i, element := range o.Elements() {
			newList[i+o.Length()] = element
		}
		result := &SyncArrayListOfValue{
			List: &newList,
		}
		return result.ToValue(), true
	default:
		return Undefined, false
	}
}

func (l *SyncNativeArrayList[T]) ConcatVal(other Value) (Value, Value) {
	result, ok := l.ConcatValOk(other)
	if !ok {
		return Undefined, Ref(Errorf(TypeErrorClass, "cannot concat %s with list %s", other.Inspect(), l.Inspect()))
	}

	return result, Undefined
}

func (l *SyncNativeArrayList[T]) Repeat(other Value) (*SyncNativeArrayList[T], Value) {
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
		newList := make(NativeArrayList[T], 0, newLen)
		for range int(o) {
			newList = append(newList, *l.List...)
		}
		result := &SyncNativeArrayList[T]{
			List: &newList,
		}
		return result, Undefined
	default:
		return nil, Ref(Errorf(TypeErrorClass, "cannot repeat a list using %s", other.Inspect()))
	}
}

func (l *SyncNativeArrayList[T]) RepeatVal(other Value) (Value, Value) {
	return RefErr(l.Repeat(other))
}

func (l *SyncNativeArrayList[T]) IterNative() *SyncNativeArrayListIterator[T] {
	return NewSyncNativeArrayListIterator(l)
}

func (l *SyncNativeArrayList[T]) Iter() NativeIterator {
	return l.IterNative()
}

func (l *SyncNativeArrayList[T]) IterSyncList() SyncArrayListIterator {
	return l.IterNative()
}

type SyncNativeArrayListIterator[T ValueInterface] struct {
	ArrayList *SyncNativeArrayList[T]
	Index     int
}

var _ SyncArrayListIterator = &SyncNativeArrayListIterator[String]{}

func NewSyncNativeArrayListIterator[T ValueInterface](list *SyncNativeArrayList[T]) *SyncNativeArrayListIterator[T] {
	return &SyncNativeArrayListIterator[T]{
		ArrayList: list,
	}
}

func NewSyncNativeArrayListIteratorWithIndex[T ValueInterface](list *SyncNativeArrayList[T], index int) *SyncNativeArrayListIterator[T] {
	return &SyncNativeArrayListIterator[T]{
		ArrayList: list,
		Index:     index,
	}
}

func (*SyncNativeArrayListIterator[T]) SingletonClass() *Class {
	return nil
}

func (*SyncNativeArrayListIterator[T]) Class() *Class {
	return SyncArrayListIteratorClass
}

func (*SyncNativeArrayListIterator[T]) DirectClass() *Class {
	return SyncArrayListIteratorClass
}

func (l *SyncNativeArrayListIterator[T]) Inspect() string {
	return fmt.Sprintf("Std::Sync::ArrayList::Iterator{&: %p, list: %s, index: %d}", l, l.ArrayList.Inspect(), l.Index)
}

func (l *SyncNativeArrayListIterator[T]) Error() string {
	return l.Inspect()
}

func (i *SyncNativeArrayListIterator[T]) ToValue() Value {
	return Ref(i)
}

func (l *SyncNativeArrayListIterator[T]) Copy() Reference {
	return &SyncNativeArrayListIterator[T]{
		ArrayList: l.ArrayList,
		Index:     l.Index,
	}
}

func (*SyncNativeArrayListIterator[T]) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *SyncNativeArrayListIterator[T]) Next() (t T, err Value) {
	l.ArrayList.Lock()
	defer l.ArrayList.Unlock()

	if l.Index >= l.ArrayList.List.Length() {
		return t, S(symbol.L_stop_iteration).ToValue()
	}

	next := (*l.ArrayList.List)[l.Index]
	l.Index++
	return next, Undefined
}

func (l *SyncNativeArrayListIterator[T]) NextValue() (Value, Value) {
	return ToValueErr(l.Next())
}

func (l *SyncNativeArrayListIterator[T]) Elements() iter.Seq[Value] {
	return func(yield func(Value) bool) {
		l.ArrayList.Lock()
		defer l.ArrayList.Unlock()

		for ; l.Index < l.ArrayList.List.Length(); l.Index++ {
			if !yield((*l.ArrayList.List)[l.Index].ToValue()) {
				return
			}
		}
	}
}

func (l *SyncNativeArrayListIterator[T]) Reset() {
	l.Index = 0
}
