package object

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

func initUInt32() {
	UInt32Class = NewClass(ClassWithParent(NumericClass), ClassWithImmutable(), ClassWithSealed())
	StdModule.AddConstant("UInt32", UInt32Class)
}
