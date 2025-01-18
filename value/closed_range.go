package value

import (
	"fmt"
	"strings"
)

// Represents a closed range eg. `5...2`
type ClosedRange struct {
	Start Value // start value
	End   Value // end value
}

// Create a new closed range class.
func NewClosedRange(start, end Value) *ClosedRange {
	return &ClosedRange{
		Start: start,
		End:   end,
	}
}

func (r *ClosedRange) Copy() Reference {
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

func (r *ClosedRange) Error() string {
	return r.Inspect()
}

func (r *ClosedRange) Inspect() string {
	var buff strings.Builder
	buff.WriteString(r.Start.Inspect())
	buff.WriteString("...")
	buff.WriteString(r.End.Inspect())

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
		CurrentElement: r.Start,
	}
}

func NewClosedRangeIteratorWithCurrentElement(r *ClosedRange, currentElement Value) *ClosedRangeIterator {
	return &ClosedRangeIterator{
		Range:          r,
		CurrentElement: currentElement,
	}
}

func (r *ClosedRangeIterator) Reset() {
	r.CurrentElement = r.Range.Start
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

func (r *ClosedRangeIterator) Copy() Reference {
	return &ClosedRangeIterator{
		Range:          r.Range,
		CurrentElement: r.CurrentElement,
	}
}

func (r *ClosedRangeIterator) Error() string {
	return r.Inspect()
}

func (r *ClosedRangeIterator) Inspect() string {
	return fmt.Sprintf("Std::ClosedRange::Iterator{range: %s, current_element: %s}", r.Range.Inspect(), r.CurrentElement.Inspect())
}

func (*ClosedRangeIterator) InstanceVariables() SymbolMap {
	return nil
}

func initClosedRange() {
	ClosedRangeClass = NewClass()
	ClosedRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("ClosedRange", Ref(ClosedRangeClass))

	ClosedRangeIteratorClass = NewClass()
	ClosedRangeClass.AddConstantString("Iterator", Ref(ClosedRangeIteratorClass))
}
