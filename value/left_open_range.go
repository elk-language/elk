package value

import (
	"fmt"
	"strings"
)

// Represents a left open range eg. `5<..2`
type LeftOpenRange struct {
	From Value // start value
	To   Value // end value
}

// Create a new left open range class.
func NewLeftOpenRange(from, to Value) *LeftOpenRange {
	return &LeftOpenRange{
		From: from,
		To:   to,
	}
}

func (r *LeftOpenRange) Copy() Value {
	return r
}

func (*LeftOpenRange) Class() *Class {
	return LeftOpenRangeClass
}

func (*LeftOpenRange) DirectClass() *Class {
	return LeftOpenRangeClass
}

func (*LeftOpenRange) SingletonClass() *Class {
	return nil
}

func (r *LeftOpenRange) Inspect() string {
	var buff strings.Builder
	buff.WriteString(r.From.Inspect())
	buff.WriteString("<..")
	buff.WriteString(r.To.Inspect())

	return buff.String()
}

func (r *LeftOpenRange) InstanceVariables() SymbolMap {
	return nil
}

var LeftOpenRangeClass *Class // ::Std::LeftOpenRange

// ::Std::LeftOpenRange::Iterator
//
// LeftOpenRange iterator class.
var LeftOpenRangeIteratorClass *Class

type LeftOpenRangeIterator struct {
	Range          *LeftOpenRange
	CurrentElement Value
}

func NewLeftOpenRangeIterator(r *LeftOpenRange) *LeftOpenRangeIterator {
	return &LeftOpenRangeIterator{
		Range:          r,
		CurrentElement: r.From,
	}
}

func NewLeftOpenRangeIteratorWithCurrentElement(r *LeftOpenRange, currentElement Value) *LeftOpenRangeIterator {
	return &LeftOpenRangeIterator{
		Range:          r,
		CurrentElement: currentElement,
	}
}

func (*LeftOpenRangeIterator) Class() *Class {
	return LeftOpenRangeIteratorClass
}

func (*LeftOpenRangeIterator) DirectClass() *Class {
	return LeftOpenRangeIteratorClass
}

func (*LeftOpenRangeIterator) SingletonClass() *Class {
	return nil
}

func (r *LeftOpenRangeIterator) Copy() Value {
	return &LeftOpenRangeIterator{
		Range:          r.Range,
		CurrentElement: r.CurrentElement,
	}
}

func (r *LeftOpenRangeIterator) Inspect() string {
	return fmt.Sprintf("Std::LeftOpenRange::Iterator{range: %s, current_element: %s}", r.Range.Inspect(), r.CurrentElement.Inspect())
}

func (*LeftOpenRangeIterator) InstanceVariables() SymbolMap {
	return nil
}

func initLeftOpenRange() {
	LeftOpenRangeClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	LeftOpenRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("LeftOpenRange", LeftOpenRangeClass)

	LeftOpenRangeIteratorClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	LeftOpenRangeClass.AddConstantString("Iterator", LeftOpenRangeIteratorClass)
}
