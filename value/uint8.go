package value

import "fmt"

var UInt8Class *Class // ::Std::UInt8

// Elk's UInt8 value
type UInt8 uint8

func (UInt8) Class() *Class {
	return UInt8Class
}

func (UInt8) DirectClass() *Class {
	return UInt64Class
}

func (UInt8) SingletonClass() *Class {
	return nil
}

func (i UInt8) IsFrozen() bool {
	return true
}

func (i UInt8) SetFrozen() {}

func (i UInt8) Inspect() string {
	return fmt.Sprintf("%du8", i)
}

func (i UInt8) InstanceVariables() SymbolMap {
	return nil
}

func initUInt8() {
	UInt8Class = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("UInt8", UInt8Class)
}
