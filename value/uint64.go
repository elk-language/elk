package value

import "fmt"

var UInt64Class *Class // ::Std::UInt64

// Elk's UInt64 value
type UInt64 uint64

func (UInt64) Class() *Class {
	return UInt64Class
}

func (UInt64) DirectClass() *Class {
	return UInt64Class
}

func (UInt64) SingletonClass() *Class {
	return nil
}

func (i UInt64) IsFrozen() bool {
	return true
}

func (i UInt64) SetFrozen() {}

func (i UInt64) Inspect() string {
	return fmt.Sprintf("%du64", i)
}

func (i UInt64) InstanceVariables() SymbolMap {
	return nil
}

func initUInt64() {
	UInt64Class = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("UInt64", UInt64Class)
}
