package object

import "fmt"

var UInt64Class *Class // ::Std::UInt64

// Elk's Int64 value
type UInt64 int64

func (i UInt64) Class() *Class {
	return UInt64Class
}

func (i UInt64) IsFrozen() bool {
	return true
}

func (i UInt64) SetFrozen() {}

func (i UInt64) Inspect() string {
	return fmt.Sprintf("%du64", i)
}

func (i UInt64) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initUInt64() {
	UInt64Class = NewClass(
		ClassWithParent(NumericClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("UInt64", UInt64Class)
}
