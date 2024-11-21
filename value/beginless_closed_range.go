package value

import (
	"strings"
)

// Represents a beginless closed range eg. `...2`
type BeginlessClosedRange struct {
	End Value // end value
}

// Create a new beginless closed range class.
func NewBeginlessClosedRange(end Value) *BeginlessClosedRange {
	return &BeginlessClosedRange{
		End: end,
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
	buff.WriteString(r.End.Inspect())

	return buff.String()
}

func (r *BeginlessClosedRange) Error() string {
	return r.Inspect()
}

func (r *BeginlessClosedRange) InstanceVariables() SymbolMap {
	return nil
}

var BeginlessClosedRangeClass *Class // ::Std::BeginlessClosedRange

func initBeginlessClosedRange() {
	BeginlessClosedRangeClass = NewClass()
	BeginlessClosedRangeClass.IncludeMixin(RangeMixin)
	StdModule.AddConstantString("BeginlessClosedRange", BeginlessClosedRangeClass)
}
