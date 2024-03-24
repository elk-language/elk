package value

import (
	"fmt"
	"strings"
)

// Represents an open range eg. `5<.<2`
type OpenRange struct {
	From Value // start value
	To   Value // end value
}

// Create a new open range class.
func NewOpenRange(from, to Value) *OpenRange {
	return &OpenRange{
		From: from,
		To:   to,
	}
}

func (r *OpenRange) Copy() Value {
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

func (r *OpenRange) Inspect() string {
	var buff strings.Builder
	buff.WriteString(r.From.Inspect())
	buff.WriteString("<.<")
	buff.WriteString(r.To.Inspect())

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
		CurrentElement: r.From,
	}
}

func NewOpenRangeIteratorWithCurrentElement(r *OpenRange, currentElement Value) *OpenRangeIterator {
	return &OpenRangeIterator{
		Range:          r,
		CurrentElement: currentElement,
	}
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

func (r *OpenRangeIterator) Copy() Value {
	return &OpenRangeIterator{
		Range:          r.Range,
		CurrentElement: r.CurrentElement,
	}
}

func (r *OpenRangeIterator) Inspect() string {
	return fmt.Sprintf("Std::OpenRange::Iterator{range: %s, current_element: %s}", r.Range.Inspect(), r.CurrentElement.Inspect())
}

func (*OpenRangeIterator) InstanceVariables() SymbolMap {
	return nil
}

func initOpenRange() {
	OpenRangeClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	OpenRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("OpenRange", OpenRangeClass)

	OpenRangeIteratorClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	OpenRangeClass.AddConstantString("Iterator", OpenRangeIteratorClass)
}
