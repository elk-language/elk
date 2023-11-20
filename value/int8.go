package value

import "fmt"

var Int8Class *Class // ::Std::Int8

// Elk's Int8 value
type Int8 int8

func (i Int8) Class() *Class {
	return Int8Class
}

func (Int8) DirectClass() *Class {
	return Int8Class
}

func (Int8) SingletonClass() *Class {
	return nil
}

func (i Int8) Inspect() string {
	return fmt.Sprintf("%di8", i)
}

func (i Int8) InstanceVariables() SymbolMap {
	return nil
}

func initInt8() {
	Int8Class = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Int8", Int8Class)
}
