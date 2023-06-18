package object

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

func initUInt8() {
	UInt8Class = NewClass(ClassWithParent(NumericClass), ClassWithImmutable(), ClassWithSealed())
	StdModule.AddConstant("UInt8", UInt8Class)
}
