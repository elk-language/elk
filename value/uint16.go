package value

import "fmt"

var UInt16Class *Class // ::Std::UInt16

// Elk's UInt16 value
type UInt16 uint16

func (UInt16) Class() *Class {
	return UInt16Class
}

func (UInt16) DirectClass() *Class {
	return UInt16Class
}

func (UInt16) SingletonClass() *Class {
	return nil
}

func (UInt16) IsFrozen() bool {
	return true
}

func (UInt16) SetFrozen() {}

func (i UInt16) Inspect() string {
	return fmt.Sprintf("%du16", i)
}

func (i UInt16) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initUInt16() {
	UInt16Class = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("UInt16", UInt16Class)
}
