package object

import "fmt"

var Int32Class *Class // ::Std::Int32

// Elk's Int32 value
type Int32 int32

func (i Int32) Class() *Class {
	return Int32Class
}

func (i Int32) IsFrozen() bool {
	return true
}

func (i Int32) SetFrozen() {}

func (i Int32) Inspect() string {
	return fmt.Sprintf("%di32", i)
}

func initInt32() {
	Int32Class = NewClass(ClassWithParent(NumericClass), ClassWithImmutable(), ClassWithSealed())
	StdModule.AddConstant("Int32", Int32Class)
}
