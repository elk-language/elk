package value

import "fmt"

var Int16Class *Class // ::Std::Int16

// Elk's Int16 value
type Int16 int16

func (i Int16) Class() *Class {
	return Int16Class
}

func (i Int16) IsFrozen() bool {
	return true
}

func (i Int16) SetFrozen() {}

func (i Int16) Inspect() string {
	return fmt.Sprintf("%di16", i)
}

func (i Int16) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initInt16() {
	Int16Class = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Int16", Int16Class)
}
