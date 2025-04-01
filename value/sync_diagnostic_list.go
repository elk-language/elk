package value

import (
	"fmt"
	"slices"
	"strings"

	"github.com/elk-language/elk/position/diagnostic"
)

var SyncDiagnosticListIteratorClass *Class // ::Std::Sync::DiagnosticList::Iterator
var SyncDiagnosticListClass *Class         // ::Std::Sync::DiagnosticList

type SyncDiagnosticList diagnostic.SyncDiagnosticList

// Creates a new SyncDiagnosticList.
func SyncDiagnosticListConstructor(class *Class) Value {
	d := (*SyncDiagnosticList)(diagnostic.NewSyncDiagnosticList())
	return Ref(d)
}

func (*SyncDiagnosticList) Class() *Class {
	return SyncDiagnosticListClass
}

func (*SyncDiagnosticList) DirectClass() *Class {
	return SyncDiagnosticListClass
}

func (*SyncDiagnosticList) SingletonClass() *Class {
	return nil
}

func (d *SyncDiagnosticList) Copy() Reference {
	return &SyncDiagnosticList{
		DiagnosticList: slices.Clone(d.DiagnosticList),
	}
}

func (*SyncDiagnosticList) InstanceVariables() SymbolMap {
	return nil
}

func (d *SyncDiagnosticList) Inspect() string {
	d.Mutex.Lock()

	var buff strings.Builder
	buff.WriteString("Std::SyncDiagnosticList[")
	for i, diag := range d.DiagnosticList {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString((*Diagnostic)(diag).Inspect())
	}
	buff.WriteRune(']')

	d.Mutex.Unlock()

	return buff.String()
}

func (d *SyncDiagnosticList) Error() string {
	return d.Inspect()
}

// Get an element under the given index.
func (dl *SyncDiagnosticList) Subscript(key Value) (Value, Value) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Undefined, Ref(NewIndexOutOfRangeError(key.Inspect(), dl.Length()))
		}
		return Undefined, Ref(NewCoerceError(IntClass, key.Class()))
	}

	return dl.Get(i)
}

// Get an element under the given index.
func (dl *SyncDiagnosticList) Get(index int) (Value, Value) {
	dl.Mutex.Lock()
	diag, err := GetFromSlice((*[]*diagnostic.Diagnostic)(&dl.DiagnosticList), index)
	dl.Mutex.Unlock()

	if !err.IsUndefined() {
		return Undefined, err
	}

	return Ref((*Diagnostic)(diag)), Undefined
}

func (dl *SyncDiagnosticList) Length() int {
	return len(dl.DiagnosticList)
}

func (l *SyncDiagnosticList) LeftCapacity() int {
	l.Mutex.Lock()
	c := l.Capacity() - l.Length()
	l.Mutex.Unlock()

	return c
}

// Expand the array list to have
// empty slots for new elements.
func (dl *SyncDiagnosticList) Grow(newSlots int) {
	dl.Mutex.Lock()

	newList := make(diagnostic.DiagnosticList, dl.Length(), dl.Capacity()+newSlots)
	copy(newList, dl.DiagnosticList)
	dl.DiagnosticList = newList

	dl.Mutex.Unlock()
}

// Add new elements.
func (dl *SyncDiagnosticList) Append(elements ...*Diagnostic) {
	dl.Mutex.Lock()
	d := (*[]*diagnostic.Diagnostic)(&dl.DiagnosticList)
	*d = slices.Grow(*d, len(elements))
	for _, element := range elements {
		*d = append(*d, (*diagnostic.Diagnostic)(element))
	}
	dl.Mutex.Unlock()
}

func (dl *SyncDiagnosticList) Capacity() int {
	return cap(dl.DiagnosticList)
}

// Get an element under the given index without bounds checking
func (dl *SyncDiagnosticList) At(i int) *Diagnostic {
	dl.Mutex.Lock()
	v := (*Diagnostic)((dl.DiagnosticList)[i])
	dl.Mutex.Unlock()

	return v
}

// Set an element under the given index.
func (l *SyncDiagnosticList) Set(index int, val *Diagnostic) Value {
	l.Mutex.Lock()
	err := SetInSlice((*[]*diagnostic.Diagnostic)(&l.DiagnosticList), index, (*diagnostic.Diagnostic)(val))
	l.Mutex.Unlock()

	return err
}

// Set an element under the given index without bounds checking.
func (l *SyncDiagnosticList) SetAt(index int, val *Diagnostic) {
	l.Mutex.Lock()
	(l.DiagnosticList)[index] = (*diagnostic.Diagnostic)(val)
	l.Mutex.Unlock()
}

// Set an element under the given index.
func (dl *SyncDiagnosticList) SubscriptSet(key Value, val *Diagnostic) Value {
	length := len(dl.DiagnosticList)
	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Ref(NewIndexOutOfRangeError(key.Inspect(), length))
		}
		return Ref(NewCoerceError(IntClass, key.Class()))
	}

	return dl.Set(i, val)
}

// Concatenate another value with this list, creating a new list, and return the result.
// If the operation is illegal an error will be returned.
func (l *SyncDiagnosticList) Concat(other Value) (*ArrayList, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *ArrayList:
			l.Mutex.Lock()

			newList := make(ArrayList, 0, l.Length()+len(*o))
			for _, element := range l.DiagnosticList {
				newList = append(newList, Ref((*Diagnostic)(element)))
			}

			l.Mutex.Unlock()

			newList = append(newList, *o...)
			return &newList, Undefined
		case *ArrayTuple:
			l.Mutex.Lock()

			newList := make(ArrayList, 0, l.Length()+len(*o))
			for _, element := range l.DiagnosticList {
				newList = append(newList, Ref((*Diagnostic)(element)))
			}

			l.Mutex.Unlock()

			for _, element := range *o {
				newList = append(newList, element)
			}
			return &newList, Undefined
		case *DiagnosticList:
			l.Mutex.Lock()

			newList := make(ArrayList, 0, l.Length()+len(*o))
			for _, element := range l.DiagnosticList {
				newList = append(newList, Ref((*Diagnostic)(element)))
			}

			l.Mutex.Unlock()

			for _, element := range *o {
				newList = append(newList, Ref((*Diagnostic)(element)))
			}
			return &newList, Undefined
		case *SyncDiagnosticList:
			l.Mutex.Lock()
			newList := make(ArrayList, 0, l.Length()+o.Length())
			for _, element := range l.DiagnosticList {
				newList = append(newList, Ref((*Diagnostic)(element)))
			}
			l.Mutex.Unlock()

			o.Mutex.Lock()
			for _, element := range o.DiagnosticList {
				newList = append(newList, Ref((*Diagnostic)(element)))
			}
			o.Mutex.Unlock()

			return &newList, Undefined
		}
	}

	return nil, Ref(Errorf(TypeErrorClass, "cannot concat %s with list %s", other.Inspect(), l.Inspect()))
}

// Repeat the content of this list n times and return a new list containing the result.
// If the operation is illegal an error will be returned.
func (l *SyncDiagnosticList) Repeat(other Value) (*SyncDiagnosticList, Value) {
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

		l.Mutex.Lock()

		newLen, ok := o.MultiplyOverflow(SmallInt(l.Length()))
		if !ok {
			return nil, Ref(Errorf(
				OutOfRangeErrorClass,
				"list repeat count is too large %s",
				o.Inspect(),
			))
		}
		newList := make(diagnostic.DiagnosticList, 0, newLen)
		for range int(o) {
			newList = append(newList, l.DiagnosticList...)
		}

		l.Mutex.Unlock()

		result := &SyncDiagnosticList{
			DiagnosticList: newList,
		}
		return result, Undefined
	default:
		return nil, Ref(Errorf(TypeErrorClass, "cannot repeat a list using %s", other.Inspect()))
	}
}

type SyncDiagnosticListIterator struct {
	SyncDiagnosticList *SyncDiagnosticList
	Index              int
}

func NewSyncDiagnosticListIterator(list *SyncDiagnosticList) *SyncDiagnosticListIterator {
	return &SyncDiagnosticListIterator{
		SyncDiagnosticList: list,
	}
}

func NewSyncDiagnosticListIteratorWithIndex(list *SyncDiagnosticList, index int) *SyncDiagnosticListIterator {
	return &SyncDiagnosticListIterator{
		SyncDiagnosticList: list,
		Index:              index,
	}
}

func (*SyncDiagnosticListIterator) Class() *Class {
	return SyncDiagnosticListIteratorClass
}

func (*SyncDiagnosticListIterator) DirectClass() *Class {
	return SyncDiagnosticListIteratorClass
}

func (*SyncDiagnosticListIterator) SingletonClass() *Class {
	return nil
}

func (l *SyncDiagnosticListIterator) Copy() Reference {
	return &SyncDiagnosticListIterator{
		SyncDiagnosticList: l.SyncDiagnosticList,
		Index:              l.Index,
	}
}

func (l *SyncDiagnosticListIterator) Inspect() string {
	return fmt.Sprintf("Std::SyncDiagnosticList::Iterator{&: %p, list: %s, index: %d}", l, l.SyncDiagnosticList.Inspect(), l.Index)
}

func (l *SyncDiagnosticListIterator) Error() string {
	return l.Inspect()
}

func (*SyncDiagnosticListIterator) InstanceVariables() SymbolMap {
	return nil
}

func (l *SyncDiagnosticListIterator) Next() (Value, Value) {
	if l.Index >= l.SyncDiagnosticList.Length() {
		return Undefined, stopIterationSymbol.ToValue()
	}

	next := l.SyncDiagnosticList.At(l.Index)
	l.Index++
	return Ref((*Diagnostic)(next)), Undefined
}

func (l *SyncDiagnosticListIterator) Reset() {
	l.Index = 0
}

func initSyncDiagnosticList() {
	SyncDiagnosticListClass = NewClassWithOptions(ClassWithConstructor(SyncDiagnosticListConstructor))
	SyncDiagnosticListClass.IncludeMixin(ListMixin)
	SyncModule.AddConstantString("DiagnosticList", Ref(SyncDiagnosticListClass))

	SyncDiagnosticListIteratorClass = NewClass()
	SyncDiagnosticListClass.AddConstantString("Iterator", Ref(SyncDiagnosticListIteratorClass))
}
