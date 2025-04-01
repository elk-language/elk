package value

import (
	"fmt"
	"slices"
	"strings"

	"github.com/elk-language/elk/position/diagnostic"
)

var DiagnosticListIteratorClass *Class // ::Std::DiagnosticList::Iterator
var DiagnosticListClass *Class         // ::Std::DiagnosticList

type DiagnosticList diagnostic.DiagnosticList

// Creates a new DiagnosticList.
func DiagnosticListConstructor(class *Class) Value {
	d := make(DiagnosticList, 5)
	return Ref(&d)
}

func (*DiagnosticList) Class() *Class {
	return DiagnosticListClass
}

func (*DiagnosticList) DirectClass() *Class {
	return DiagnosticListClass
}

func (*DiagnosticList) SingletonClass() *Class {
	return nil
}

func (d *DiagnosticList) Copy() Reference {
	newList := make(DiagnosticList, len(*d))
	copy(newList, *d)
	return &newList
}

func (*DiagnosticList) InstanceVariables() SymbolMap {
	return nil
}

func (d *DiagnosticList) Inspect() string {
	var buff strings.Builder

	buff.WriteString("Std::DiagnosticList[")
	for i, diag := range *d {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString((*Diagnostic)(diag).Inspect())
	}
	buff.WriteRune(']')

	return buff.String()
}

func (d *DiagnosticList) Error() string {
	return d.Inspect()
}

// Get an element under the given index.
func (dl *DiagnosticList) Subscript(key Value) (Value, Value) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Undefined, Ref(NewIndexOutOfRangeError(key.Inspect(), len(*dl)))
		}
		return Undefined, Ref(NewCoerceError(IntClass, key.Class()))
	}

	return dl.Get(i)
}

// Get an element under the given index.
func (dl *DiagnosticList) Get(index int) (Value, Value) {
	diag, err := GetFromSlice((*[]*diagnostic.Diagnostic)(dl), index)
	if !err.IsUndefined() {
		return Undefined, err
	}

	return Ref((*Diagnostic)(diag)), Undefined
}

func (dl *DiagnosticList) Length() int {
	return len(*dl)
}

func (l *DiagnosticList) LeftCapacity() int {
	return l.Capacity() - l.Length()
}

// Expand the array list to have
// empty slots for new elements.
func (dl *DiagnosticList) Grow(newSlots int) {
	newList := make(DiagnosticList, dl.Length(), dl.Capacity()+newSlots)
	copy(newList, *dl)
	*dl = newList
}

// Add new elements.
func (dl *DiagnosticList) Append(elements ...*Diagnostic) {
	d := (*[]*diagnostic.Diagnostic)(dl)
	*d = slices.Grow(*d, len(elements))
	for _, element := range elements {
		*d = append(*d, (*diagnostic.Diagnostic)(element))
	}
}

func (dl *DiagnosticList) Capacity() int {
	return cap(*dl)
}

// Get an element under the given index without bounds checking
func (dl *DiagnosticList) At(i int) *Diagnostic {
	return (*Diagnostic)((*dl)[i])
}

// Set an element under the given index.
func (l *DiagnosticList) Set(index int, val *Diagnostic) Value {
	return SetInSlice((*[]*diagnostic.Diagnostic)(l), index, (*diagnostic.Diagnostic)(val))
}

// Set an element under the given index without bounds checking.
func (l *DiagnosticList) SetAt(index int, val *Diagnostic) {
	(*l)[index] = (*diagnostic.Diagnostic)(val)
}

// Set an element under the given index.
func (dl *DiagnosticList) SubscriptSet(key Value, val *Diagnostic) Value {
	length := len(*dl)
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
func (l *DiagnosticList) Concat(other Value) (*ArrayList, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *ArrayList:
			newList := make(ArrayList, 0, len(*l)+len(*o))
			for _, element := range *l {
				newList = append(newList, Ref((*Diagnostic)(element)))
			}
			newList = append(newList, *o...)
			return &newList, Undefined
		case *ArrayTuple:
			newList := make(ArrayList, 0, len(*l)+len(*o))
			for _, element := range *l {
				newList = append(newList, Ref((*Diagnostic)(element)))
			}
			newList = append(newList, *o...)
			return &newList, Undefined
		case *DiagnosticList:
			newList := make(ArrayList, 0, len(*l)+len(*o))
			for _, element := range *l {
				newList = append(newList, Ref((*Diagnostic)(element)))
			}
			for _, element := range *o {
				newList = append(newList, Ref((*Diagnostic)(element)))
			}
			return &newList, Undefined
		}
	}

	return nil, Ref(Errorf(TypeErrorClass, "cannot concat %s with list %s", other.Inspect(), l.Inspect()))
}

// Repeat the content of this list n times and return a new list containing the result.
// If the operation is illegal an error will be returned.
func (l *DiagnosticList) Repeat(other Value) (*DiagnosticList, Value) {
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
		newList := make(DiagnosticList, 0, newLen)
		for range int(o) {
			newList = append(newList, *l...)
		}
		return &newList, Undefined
	default:
		return nil, Ref(Errorf(TypeErrorClass, "cannot repeat a list using %s", other.Inspect()))
	}
}

type DiagnosticListIterator struct {
	DiagnosticList *DiagnosticList
	Index          int
}

func NewDiagnosticListIterator(list *DiagnosticList) *DiagnosticListIterator {
	return &DiagnosticListIterator{
		DiagnosticList: list,
	}
}

func NewDiagnosticListIteratorWithIndex(list *DiagnosticList, index int) *DiagnosticListIterator {
	return &DiagnosticListIterator{
		DiagnosticList: list,
		Index:          index,
	}
}

func (*DiagnosticListIterator) Class() *Class {
	return DiagnosticListIteratorClass
}

func (*DiagnosticListIterator) DirectClass() *Class {
	return DiagnosticListIteratorClass
}

func (*DiagnosticListIterator) SingletonClass() *Class {
	return nil
}

func (l *DiagnosticListIterator) Copy() Reference {
	return &DiagnosticListIterator{
		DiagnosticList: l.DiagnosticList,
		Index:          l.Index,
	}
}

func (l *DiagnosticListIterator) Inspect() string {
	return fmt.Sprintf("Std::DiagnosticList::Iterator{&: %p, list: %s, index: %d}", l, l.DiagnosticList.Inspect(), l.Index)
}

func (l *DiagnosticListIterator) Error() string {
	return l.Inspect()
}

func (*DiagnosticListIterator) InstanceVariables() SymbolMap {
	return nil
}

func (l *DiagnosticListIterator) Next() (Value, Value) {
	if l.Index >= l.DiagnosticList.Length() {
		return Undefined, stopIterationSymbol.ToValue()
	}

	next := (*l.DiagnosticList)[l.Index]
	l.Index++
	return Ref((*Diagnostic)(next)), Undefined
}

func (l *DiagnosticListIterator) Reset() {
	l.Index = 0
}

func initDiagnosticList() {
	DiagnosticListClass = NewClassWithOptions(ClassWithConstructor(DiagnosticListConstructor))
	DiagnosticListClass.IncludeMixin(ListMixin)
	StdModule.AddConstantString("DiagnosticList", Ref(DiagnosticListClass))

	DiagnosticListIteratorClass = NewClass()
	DiagnosticListClass.AddConstantString("Iterator", Ref(DiagnosticListIteratorClass))
}
