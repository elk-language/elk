package value

import (
	"fmt"
	"strings"
)

// Represents a left open range eg. `5<..2`
type LeftOpenRange struct {
	Start Value // start value
	End   Value // end value
}

// Create a new left open range class.
func NewLeftOpenRange(start, end Value) *LeftOpenRange {
	return &LeftOpenRange{
		Start: start,
		End:   end,
	}
}

func (r *LeftOpenRange) Copy() Reference {
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

func (r *LeftOpenRange) Error() string {
	return r.Inspect()
}

func (r *LeftOpenRange) Inspect() string {
	var buff strings.Builder
	buff.WriteString(r.Start.Inspect())
	buff.WriteString("<..")
	buff.WriteString(r.End.Inspect())

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
		CurrentElement: r.Start,
	}
}

func NewLeftOpenRangeIteratorWithCurrentElement(r *LeftOpenRange, currentElement Value) *LeftOpenRangeIterator {
	return &LeftOpenRangeIterator{
		Range:          r,
		CurrentElement: currentElement,
	}
}

func (r *LeftOpenRangeIterator) Reset() {
	r.CurrentElement = r.Range.Start
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

func (r *LeftOpenRangeIterator) Copy() Reference {
	return &LeftOpenRangeIterator{
		Range:          r.Range,
		CurrentElement: r.CurrentElement,
	}
}

func (r *LeftOpenRangeIterator) Inspect() string {
	return fmt.Sprintf("Std::LeftOpenRange::Iterator{range: %s, current_element: %s}", r.Range.Inspect(), r.CurrentElement.Inspect())
}

func (r *LeftOpenRangeIterator) Error() string {
	return r.Inspect()
}

func (*LeftOpenRangeIterator) InstanceVariables() SymbolMap {
	return nil
}

func initLeftOpenRange() {
	LeftOpenRangeClass = NewClass()
	LeftOpenRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("LeftOpenRange", Ref(LeftOpenRangeClass))

	LeftOpenRangeIteratorClass = NewClass()
	LeftOpenRangeClass.AddConstantString("Iterator", Ref(LeftOpenRangeIteratorClass))
}
