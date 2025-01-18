package value

import (
	"fmt"
	"strings"
)

// Represents a right open range eg. `5..<2`
type RightOpenRange struct {
	Start Value // start value
	End   Value // end value
}

// Create a new right open range class.
func NewRightOpenRange(start, end Value) *RightOpenRange {
	return &RightOpenRange{
		Start: start,
		End:   end,
	}
}

func (r *RightOpenRange) Copy() Reference {
	return r
}

func (*RightOpenRange) Class() *Class {
	return RightOpenRangeClass
}

func (*RightOpenRange) DirectClass() *Class {
	return RightOpenRangeClass
}

func (*RightOpenRange) SingletonClass() *Class {
	return nil
}

func (r *RightOpenRange) Inspect() string {
	var buff strings.Builder
	buff.WriteString(r.Start.Inspect())
	buff.WriteString("..<")
	buff.WriteString(r.End.Inspect())

	return buff.String()
}

func (r *RightOpenRange) Error() string {
	return r.Inspect()
}

func (r *RightOpenRange) InstanceVariables() SymbolMap {
	return nil
}

var RightOpenRangeClass *Class // ::Std::RightOpenRange

// ::Std::RightOpenRange::Iterator
//
// RightOpenRange iterator class.
var RightOpenRangeIteratorClass *Class

type RightOpenRangeIterator struct {
	Range          *RightOpenRange
	CurrentElement Value
}

func NewRightOpenRangeIterator(r *RightOpenRange) *RightOpenRangeIterator {
	return &RightOpenRangeIterator{
		Range:          r,
		CurrentElement: r.Start,
	}
}

func NewRightOpenRangeIteratorWithCurrentElement(r *RightOpenRange, currentElement Value) *RightOpenRangeIterator {
	return &RightOpenRangeIterator{
		Range:          r,
		CurrentElement: currentElement,
	}
}

func (r *RightOpenRangeIterator) Reset() {
	r.CurrentElement = r.Range.Start
}

func (*RightOpenRangeIterator) Class() *Class {
	return RightOpenRangeIteratorClass
}

func (*RightOpenRangeIterator) DirectClass() *Class {
	return RightOpenRangeIteratorClass
}

func (*RightOpenRangeIterator) SingletonClass() *Class {
	return nil
}

func (r *RightOpenRangeIterator) Copy() Reference {
	return &RightOpenRangeIterator{
		Range:          r.Range,
		CurrentElement: r.CurrentElement,
	}
}

func (r *RightOpenRangeIterator) Inspect() string {
	return fmt.Sprintf("Std::RightOpenRange::Iterator{range: %s, current_element: %s}", r.Range.Inspect(), r.CurrentElement.Inspect())
}

func (r *RightOpenRangeIterator) Error() string {
	return r.Inspect()
}

func (*RightOpenRangeIterator) InstanceVariables() SymbolMap {
	return nil
}

func initRightOpenRange() {
	RightOpenRangeClass = NewClass()
	RightOpenRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("RightOpenRange", Ref(RightOpenRangeClass))

	RightOpenRangeIteratorClass = NewClass()
	RightOpenRangeClass.AddConstantString("Iterator", Ref(RightOpenRangeIteratorClass))
}
