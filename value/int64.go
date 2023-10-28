package value

import "fmt"

var Int64Class *Class // ::Std::Int64

// Elk's Int64 value
type Int64 int64

func (i Int64) Class() *Class {
	return Int64Class
}

func (i Int64) IsFrozen() bool {
	return true
}

func (i Int64) SetFrozen() {}

func (i Int64) Inspect() string {
	return fmt.Sprintf("%di64", i)
}

func (i Int64) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initInt64() {
	Int64Class = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Int64", Int64Class)
}
