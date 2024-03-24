package value

import (
	"fmt"
	"strings"
)

// Represents a closed range eg. `5...2`
type ClosedRange struct {
	From Value // start value
	To   Value // end value
}

// Create a new closed range class.
func NewClosedRange(from, to Value) *ClosedRange {
	return &ClosedRange{
		From: from,
		To:   to,
	}
}

func (r *ClosedRange) Copy() Value {
	return r
}

func (*ClosedRange) Class() *Class {
	return ClosedRangeClass
}

func (*ClosedRange) DirectClass() *Class {
	return ClosedRangeClass
}

func (*ClosedRange) SingletonClass() *Class {
	return nil
}

func (r *ClosedRange) Inspect() string {
	var buff strings.Builder
	buff.WriteString(r.From.Inspect())
	buff.WriteString("...")
	buff.WriteString(r.To.Inspect())

	return buff.String()
}

func (r *ClosedRange) InstanceVariables() SymbolMap {
	return nil
}

var ClosedRangeClass *Class // ::Std::ClosedRange

// ::Std::ClosedRange::Iterator
//
// ClosedRange iterator class.
var ClosedRangeIteratorClass *Class

type ClosedRangeIterator struct {
	Range          *ClosedRange
	CurrentElement Value
}

func NewClosedRangeIterator(r *ClosedRange) *ClosedRangeIterator {
	return &ClosedRangeIterator{
		Range:          r,
		CurrentElement: r.From,
	}
}

func NewClosedRangeIteratorWithCurrentElement(r *ClosedRange, currentElement Value) *ClosedRangeIterator {
	return &ClosedRangeIterator{
		Range:          r,
		CurrentElement: currentElement,
	}
}

func (*ClosedRangeIterator) Class() *Class {
	return ClosedRangeIteratorClass
}

func (*ClosedRangeIterator) DirectClass() *Class {
	return ClosedRangeIteratorClass
}

func (*ClosedRangeIterator) SingletonClass() *Class {
	return nil
}

func (r *ClosedRangeIterator) Copy() Value {
	return &ClosedRangeIterator{
		Range:          r.Range,
		CurrentElement: r.CurrentElement,
	}
}

func (r *ClosedRangeIterator) Inspect() string {
	return fmt.Sprintf("Std::ClosedRange::Iterator{range: %s, current_element: %s}", r.Range.Inspect(), r.CurrentElement.Inspect())
}

func (*ClosedRangeIterator) InstanceVariables() SymbolMap {
	return nil
}

func initClosedRange() {
	ClosedRangeClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	ClosedRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("ClosedRange", ClosedRangeClass)

	ClosedRangeIteratorClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	ClosedRangeClass.AddConstantString("Iterator", ClosedRangeIteratorClass)
}
