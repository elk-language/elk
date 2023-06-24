package object

import "fmt"

var UInt32Class *Class // ::Std::UInt32

// Elk's UInt32 value
type UInt32 uint32

func (i UInt32) Class() *Class {
	return UInt32Class
}

func (i UInt32) IsFrozen() bool {
	return true
}

func (i UInt32) SetFrozen() {}

func (i UInt32) Inspect() string {
	return fmt.Sprintf("%du32", i)
}

func (i UInt32) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initUInt32() {
	UInt32Class = NewClass(
		ClassWithParent(NumericClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("UInt32", UInt32Class)
}
