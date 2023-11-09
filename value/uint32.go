package value

import "fmt"

var UInt32Class *Class // ::Std::UInt32

// Elk's UInt32 value
type UInt32 uint32

func (UInt32) Class() *Class {
	return UInt32Class
}

func (UInt32) DirectClass() *Class {
	return UInt32Class
}

func (UInt32) SingletonClass() *Class {
	return nil
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
	UInt32Class = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("UInt32", UInt32Class)
}
