package object

import "fmt"

var UInt8Class *Class // ::Std::UInt8

// Elk's UInt8 value
type UInt8 uint8

func (i UInt8) Class() *Class {
	return UInt8Class
}

func (i UInt8) IsFrozen() bool {
	return true
}

func (i UInt8) SetFrozen() {}

func (i UInt8) Inspect() string {
	return fmt.Sprintf("%du8", i)
}

func (i UInt8) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initUInt8() {
	UInt8Class = NewClass(
		ClassWithParent(NumericClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("UInt8", UInt8Class)
}
