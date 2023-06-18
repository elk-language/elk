package object

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

func initUInt64() {
	UInt64Class = NewClass(ClassWithParent(NumericClass), ClassWithImmutable(), ClassWithSealed())
	StdModule.AddConstant("UInt64", UInt64Class)
}
