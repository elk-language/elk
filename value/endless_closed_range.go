package value

import (
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

func initEndlessClosedRange() {
	EndlessClosedRangeClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	EndlessClosedRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("EndlessClosedRange", EndlessClosedRangeClass)
}
