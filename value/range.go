package value

import (
	"fmt"
)

// Represents a range.
type Range struct {
	From      Value // start value
	To        Value // end value
	Exclusive bool  // whether the range is exclusive or inclusive
}

// Create a new class.
func NewRange(from, to Value, exclusive bool) *Range {
	return &Range{
		From:      from,
		To:        to,
		Exclusive: exclusive,
	}
}

func (*Range) Class() *Class {
	return RangeClass
}

func (*Range) DirectClass() *Class {
	return RangeClass
}

func (*Range) SingletonClass() *Class {
	return nil
}

func (r *Range) IsFrozen() bool {
	return true
}

func (r *Range) SetFrozen() {}

const (
	inclusiveRangeOp = ".."
	exclusiveRangeOp = "..."
)

func (r *Range) Inspect() string {
	var from, op, to string
	if r.Exclusive {
		op = exclusiveRangeOp
	} else {
		op = inclusiveRangeOp
	}

	if r.From != nil {
		from = r.From.Inspect()
	}
	if r.To != nil {
		to = r.To.Inspect()
	}
	return fmt.Sprintf("%s%s%s", from, op, to)
}

func (r *Range) InstanceVariables() SimpleSymbolMap {
	return nil
}

var RangeClass *Class // ::Std::Range

func initRange() {
	RangeClass = NewClassWithOptions(
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Range", RangeClass)
}
