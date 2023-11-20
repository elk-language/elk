package value

import "fmt"

var Int32Class *Class // ::Std::Int32

// Elk's Int32 value
type Int32 int32

func (i Int32) Class() *Class {
	return Int32Class
}

func (Int32) DirectClass() *Class {
	return Int32Class
}

func (Int32) SingletonClass() *Class {
	return nil
}

func (i Int32) Inspect() string {
	return fmt.Sprintf("%di32", i)
}

func (i Int32) InstanceVariables() SymbolMap {
	return nil
}

func initInt32() {
	Int32Class = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Int32", Int32Class)
}
