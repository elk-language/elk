package value

import (
	"strings"
)

// Represents an endless open range eg. `5<..`
type EndlessOpenRange struct {
	From Value // start value
}

// Create a new endless open range class.
func NewEndlessOpenRange(from Value) *EndlessOpenRange {
	return &EndlessOpenRange{
		From: from,
	}
}

func (r *EndlessOpenRange) Copy() Value {
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

func (r *EndlessOpenRange) Inspect() string {
	var buff strings.Builder
	buff.WriteString(r.From.Inspect())
	buff.WriteString("<..")

	return buff.String()
}

func (r *EndlessOpenRange) InstanceVariables() SymbolMap {
	return nil
}

var EndlessOpenRangeClass *Class // ::Std::EndlessOpenRange

func initEndlessOpenRange() {
	EndlessOpenRangeClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	EndlessOpenRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("EndlessOpenRange", EndlessOpenRangeClass)
}
