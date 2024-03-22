package value

import (
	"strings"
)

// Represents an open range eg. `5<.<2`
type OpenRange struct {
	From Value // start value
	To   Value // end value
}

// Create a new open range class.
func NewOpenRange(from, to Value) *OpenRange {
	return &OpenRange{
		From: from,
		To:   to,
	}
}

func (r *OpenRange) Copy() Value {
	return r
}

func (*OpenRange) Class() *Class {
	return OpenRangeClass
}

func (*OpenRange) DirectClass() *Class {
	return OpenRangeClass
}

func (*OpenRange) SingletonClass() *Class {
	return nil
}

func (r *OpenRange) Inspect() string {
	var buff strings.Builder
	buff.WriteString(r.From.Inspect())
	buff.WriteString("<.<")
	buff.WriteString(r.To.Inspect())

	return buff.String()
}

func (r *OpenRange) InstanceVariables() SymbolMap {
	return nil
}

var OpenRangeClass *Class // ::Std::OpenRange

func initOpenRange() {
	OpenRangeClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	OpenRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("OpenRange", OpenRangeClass)
}
