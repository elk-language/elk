package object

import "fmt"

var UInt16Class *Class // ::Std::UInt16

// Elk's UInt16 value
type UInt16 uint16

func (i UInt16) Class() *Class {
	return UInt16Class
}

func (i UInt16) IsFrozen() bool {
	return true
}

func (i UInt16) SetFrozen() {}

func (i UInt16) Inspect() string {
	return fmt.Sprintf("%du16", i)
}

func (i UInt16) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initUInt16() {
	UInt16Class = NewClass(
		ClassWithParent(NumericClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("UInt16", UInt16Class)
}
