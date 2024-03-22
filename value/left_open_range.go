package value

import (
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

func initLeftOpenRange() {
	LeftOpenRangeClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	LeftOpenRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("LeftOpenRange", LeftOpenRangeClass)
}
