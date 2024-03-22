package value

import (
	"strings"
)

// Represents a beginless open range eg. `..<2`
type BeginlessOpenRange struct {
	To Value // end value
}

// Create a new beginless open range class.
func NewBeginlessOpenRange(to Value) *BeginlessOpenRange {
	return &BeginlessOpenRange{
		To: to,
	}
}

func (r *BeginlessOpenRange) Copy() Value {
	return r
}

func (*BeginlessOpenRange) Class() *Class {
	return BeginlessOpenRangeClass
}

func (*BeginlessOpenRange) DirectClass() *Class {
	return BeginlessOpenRangeClass
}

func (*BeginlessOpenRange) SingletonClass() *Class {
	return nil
}

func (r *BeginlessOpenRange) Inspect() string {
	var buff strings.Builder
	buff.WriteString("..<")
	buff.WriteString(r.To.Inspect())

	return buff.String()
}

func (r *BeginlessOpenRange) InstanceVariables() SymbolMap {
	return nil
}

var BeginlessOpenRangeClass *Class // ::Std::BeginlessOpenRange

func initBeginlessOpenRange() {
	BeginlessOpenRangeClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	BeginlessOpenRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("BeginlessOpenRange", BeginlessOpenRangeClass)
}
