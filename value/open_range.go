package value

import (
	"fmt"
	"strings"
)

// Represents an open range eg. `5<.<2`
type OpenRange struct {
	Start Value // start value
	End   Value // end value
}

// Create a new open range class.
func NewOpenRange(start, end Value) *OpenRange {
	return &OpenRange{
		Start: start,
		End:   end,
	}
}

func (r *OpenRange) Copy() Reference {
	return r
}

func (*OpenRange) Class() *Class {
	return OpenRangeClass
}

func (*OpenRange) DirectClass() *Class {
	return OpenRangeClass
}

func (*OpenRange) SingletonClass() *Class {
	return nil
}

func (r *OpenRange) Error() string {
	return r.Inspect()
}

func (r *OpenRange) Inspect() string {
	var buff strings.Builder
	buff.WriteString(r.Start.Inspect())
	buff.WriteString("<.<")
	buff.WriteString(r.End.Inspect())

	return buff.String()
}

func (r *OpenRange) InstanceVariables() SymbolMap {
	return nil
}

var OpenRangeClass *Class // ::Std::OpenRange

// ::Std::OpenRange::Iterator
//
// OpenRange iterator class.
var OpenRangeIteratorClass *Class

type OpenRangeIterator struct {
	Range          *OpenRange
	CurrentElement Value
}

func NewOpenRangeIterator(r *OpenRange) *OpenRangeIterator {
	return &OpenRangeIterator{
		Range:          r,
		CurrentElement: r.Start,
	}
}

func NewOpenRangeIteratorWithCurrentElement(r *OpenRange, currentElement Value) *OpenRangeIterator {
	return &OpenRangeIterator{
		Range:          r,
		CurrentElement: currentElement,
	}
}

func (r *OpenRangeIterator) Reset() {
	r.CurrentElement = r.Range.Start
}

func (*OpenRangeIterator) Class() *Class {
	return OpenRangeIteratorClass
}

func (*OpenRangeIterator) DirectClass() *Class {
	return OpenRangeIteratorClass
}

func (*OpenRangeIterator) SingletonClass() *Class {
	return nil
}

func (r *OpenRangeIterator) Copy() Reference {
	return &OpenRangeIterator{
		Range:          r.Range,
		CurrentElement: r.CurrentElement,
	}
}

func (r *OpenRangeIterator) Inspect() string {
	return fmt.Sprintf("Std::OpenRange::Iterator{range: %s, current_element: %s}", r.Range.Inspect(), r.CurrentElement.Inspect())
}

func (r *OpenRangeIterator) Error() string {
	return r.Inspect()
}

func (*OpenRangeIterator) InstanceVariables() SymbolMap {
	return nil
}

func initOpenRange() {
	OpenRangeClass = NewClass()
	OpenRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("OpenRange", Ref(OpenRangeClass))

	OpenRangeIteratorClass = NewClass()
	OpenRangeClass.AddConstantString("Iterator", Ref(OpenRangeIteratorClass))
}
