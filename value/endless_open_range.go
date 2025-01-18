package value

import (
	"fmt"
	"strings"
)

// Represents an endless open range eg. `5<..`
type EndlessOpenRange struct {
	Start Value // start value
}

// Create a new endless open range class.
func NewEndlessOpenRange(start Value) *EndlessOpenRange {
	return &EndlessOpenRange{
		Start: start,
	}
}

func (r *EndlessOpenRange) Copy() Reference {
	return r
}

func (*EndlessOpenRange) Class() *Class {
	return EndlessOpenRangeClass
}

func (*EndlessOpenRange) DirectClass() *Class {
	return EndlessOpenRangeClass
}

func (*EndlessOpenRange) SingletonClass() *Class {
	return nil
}

func (r *EndlessOpenRange) Error() string {
	return r.Inspect()
}

func (r *EndlessOpenRange) Inspect() string {
	var buff strings.Builder
	buff.WriteString(r.Start.Inspect())
	buff.WriteString("<..")

	return buff.String()
}

func (r *EndlessOpenRange) InstanceVariables() SymbolMap {
	return nil
}

var EndlessOpenRangeClass *Class // ::Std::EndlessOpenRange

// ::Std::EndlessOpenRange::Iterator
//
// EndlessOpenRange iterator class.
var EndlessOpenRangeIteratorClass *Class

type EndlessOpenRangeIterator struct {
	Range          *EndlessOpenRange
	CurrentElement Value
}

func NewEndlessOpenRangeIterator(r *EndlessOpenRange) *EndlessOpenRangeIterator {
	return &EndlessOpenRangeIterator{
		Range:          r,
		CurrentElement: r.Start,
	}
}

func NewEndlessOpenRangeIteratorWithCurrentElement(r *EndlessOpenRange, currentElement Value) *EndlessOpenRangeIterator {
	return &EndlessOpenRangeIterator{
		Range:          r,
		CurrentElement: currentElement,
	}
}

func (r *EndlessOpenRangeIterator) Reset() {
	r.CurrentElement = r.Range.Start
}

func (*EndlessOpenRangeIterator) Class() *Class {
	return EndlessOpenRangeIteratorClass
}

func (*EndlessOpenRangeIterator) DirectClass() *Class {
	return EndlessOpenRangeIteratorClass
}

func (*EndlessOpenRangeIterator) SingletonClass() *Class {
	return nil
}

func (r *EndlessOpenRangeIterator) Copy() Reference {
	return &EndlessOpenRangeIterator{
		Range:          r.Range,
		CurrentElement: r.CurrentElement,
	}
}

func (r *EndlessOpenRangeIterator) Error() string {
	return r.Inspect()
}

func (r *EndlessOpenRangeIterator) Inspect() string {
	return fmt.Sprintf("Std::EndlessOpenRange::Iterator{range: %s, current_element: %s}", r.Range.Inspect(), r.CurrentElement.Inspect())
}

func (*EndlessOpenRangeIterator) InstanceVariables() SymbolMap {
	return nil
}

func initEndlessOpenRange() {
	EndlessOpenRangeClass = NewClass()
	EndlessOpenRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("EndlessOpenRange", Ref(EndlessOpenRangeClass))

	EndlessOpenRangeIteratorClass = NewClass()
	EndlessOpenRangeClass.AddConstantString("Iterator", Ref(EndlessOpenRangeIteratorClass))
}
