package value

import (
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

func initClosedRange() {
	ClosedRangeClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	ClosedRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("ClosedRange", ClosedRangeClass)
}
