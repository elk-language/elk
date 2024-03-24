package value

import (
	"fmt"
	"strings"
)

// Represents an endless closed range eg. `5...`
type EndlessClosedRange struct {
	From Value // start value
}

// Create a new endless closed range class.
func NewEndlessClosedRange(from Value) *EndlessClosedRange {
	return &EndlessClosedRange{
		From: from,
	}
}

func (r *EndlessClosedRange) Copy() Value {
	return r
}

func (*EndlessClosedRange) Class() *Class {
	return EndlessClosedRangeClass
}

func (*EndlessClosedRange) DirectClass() *Class {
	return EndlessClosedRangeClass
}

func (*EndlessClosedRange) SingletonClass() *Class {
	return nil
}

func (r *EndlessClosedRange) Inspect() string {
	var buff strings.Builder
	buff.WriteString(r.From.Inspect())
	buff.WriteString("...")

	return buff.String()
}

func (r *EndlessClosedRange) InstanceVariables() SymbolMap {
	return nil
}

var EndlessClosedRangeClass *Class // ::Std::EndlessClosedRange

// ::Std::EndlessClosedRange::Iterator
//
// EndlessClosedRange iterator class.
var EndlessClosedRangeIteratorClass *Class

type EndlessClosedRangeIterator struct {
	Range          *EndlessClosedRange
	CurrentElement Value
}

func NewEndlessClosedRangeIterator(r *EndlessClosedRange) *EndlessClosedRangeIterator {
	return &EndlessClosedRangeIterator{
		Range:          r,
		CurrentElement: r.From,
	}
}

func NewEndlessClosedRangeIteratorWithCurrentElement(r *EndlessClosedRange, currentElement Value) *EndlessClosedRangeIterator {
	return &EndlessClosedRangeIterator{
		Range:          r,
		CurrentElement: currentElement,
	}
}

func (*EndlessClosedRangeIterator) Class() *Class {
	return EndlessClosedRangeIteratorClass
}

func (*EndlessClosedRangeIterator) DirectClass() *Class {
	return EndlessClosedRangeIteratorClass
}

func (*EndlessClosedRangeIterator) SingletonClass() *Class {
	return nil
}

func (r *EndlessClosedRangeIterator) Copy() Value {
	return &EndlessClosedRangeIterator{
		Range:          r.Range,
		CurrentElement: r.CurrentElement,
	}
}

func (r *EndlessClosedRangeIterator) Inspect() string {
	return fmt.Sprintf("Std::EndlessClosedRange::Iterator{range: %s, current_element: %s}", r.Range.Inspect(), r.CurrentElement.Inspect())
}

func (*EndlessClosedRangeIterator) InstanceVariables() SymbolMap {
	return nil
}

func initEndlessClosedRange() {
	EndlessClosedRangeClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	EndlessClosedRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("EndlessClosedRange", EndlessClosedRangeClass)

	EndlessClosedRangeIteratorClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	EndlessClosedRangeClass.AddConstantString("Iterator", EndlessClosedRangeIteratorClass)
}
