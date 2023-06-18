package object

var UInt16Class *Class // ::Std::UInt16

// Elk's UInt32 value
type UInt16 uint32

func (i UInt16) Class() *Class {
	return UInt16Class
}

func (i UInt16) IsFrozen() bool {
	return true
}

func (i UInt16) SetFrozen() {}

func initUInt16() {
	UInt16Class = NewClass(ClassWithParent(NumericClass), ClassWithImmutable(), ClassWithSealed())
	StdModule.AddConstant("UInt16", UInt16Class)
}
