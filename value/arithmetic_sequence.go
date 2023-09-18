package value

import (
	"fmt"
)

// Represents a slice of an arithmetic sequence.
type ArithmeticSequence struct {
	From      Value // start value
	To        Value // end value
	Step      Value // difference between elements
	Exclusive bool  // whether the slice is exclusive or inclusive
}

// Create a new arithmetic sequence.
func NewArithmeticSequence(from, to, step Value, exclusive bool) *ArithmeticSequence {
	return &ArithmeticSequence{
		From:      from,
		To:        to,
		Step:      step,
		Exclusive: exclusive,
	}
}

func (a *ArithmeticSequence) Class() *Class {
	return ArithmeticSequenceClass
}

func (*ArithmeticSequence) IsFrozen() bool {
	return true
}

func (*ArithmeticSequence) SetFrozen() {}

func (a *ArithmeticSequence) Inspect() string {
	var op, to string
	if a.Exclusive {
		op = exclusiveRangeOp
	} else {
		op = inclusiveRangeOp
	}

	if a.To != nil {
		to = a.To.Inspect()
	}
	return fmt.Sprintf("%s%s%s:%s", a.From.Inspect(), op, to, a.Step.Inspect())
}

func (a *ArithmeticSequence) InstanceVariables() SimpleSymbolMap {
	return nil
}

var ArithmeticSequenceClass *Class // ::Std::ArithmeticSequence

func initArithmeticSequence() {
	ArithmeticSequenceClass = NewClass(
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("ArithmeticSequence", ArithmeticSequenceClass)
}
