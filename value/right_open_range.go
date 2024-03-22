package value

import (
	"strings"
)

// Represents a right open range eg. `5..<2`
type RightOpenRange struct {
	From Value // start value
	To   Value // end value
}

// Create a new right open range class.
func NewRightOpenRange(from, to Value) *RightOpenRange {
	return &RightOpenRange{
		From: from,
		To:   to,
	}
}

func (r *RightOpenRange) Copy() Value {
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
	buff.WriteString(r.From.Inspect())
	buff.WriteString("..<")
	buff.WriteString(r.To.Inspect())

	return buff.String()
}

func (r *RightOpenRange) InstanceVariables() SymbolMap {
	return nil
}

var RightOpenRangeClass *Class // ::Std::RightOpenRange

func initRightOpenRange() {
	RightOpenRangeClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	RightOpenRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("RightOpenRange", RightOpenRangeClass)
}
