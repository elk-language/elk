package value

import (
	"strings"
)

// Represents a beginless closed range eg. `...2`
type BeginlessClosedRange struct {
	To Value // end value
}

// Create a new beginless closed range class.
func NewBeginlessClosedRange(to Value) *BeginlessClosedRange {
	return &BeginlessClosedRange{
		To: to,
	}
}

func (r *BeginlessClosedRange) Copy() Value {
	return r
}

func (*BeginlessClosedRange) Class() *Class {
	return BeginlessClosedRangeClass
}

func (*BeginlessClosedRange) DirectClass() *Class {
	return BeginlessClosedRangeClass
}

func (*BeginlessClosedRange) SingletonClass() *Class {
	return nil
}

func (r *BeginlessClosedRange) Inspect() string {
	var buff strings.Builder
	buff.WriteString("...")
	buff.WriteString(r.To.Inspect())

	return buff.String()
}

func (r *BeginlessClosedRange) InstanceVariables() SymbolMap {
	return nil
}

var BeginlessClosedRangeClass *Class // ::Std::BeginlessClosedRange

func initBeginlessClosedRange() {
	BeginlessClosedRangeClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	BeginlessClosedRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("BeginlessClosedRange", BeginlessClosedRangeClass)
}
