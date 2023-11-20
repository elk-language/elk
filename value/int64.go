package value

import "fmt"

var Int64Class *Class // ::Std::Int64

// Elk's Int64 value
type Int64 int64

func (i Int64) Class() *Class {
	return Int64Class
}

func (Int64) DirectClass() *Class {
	return Int64Class
}

func (Int64) SingletonClass() *Class {
	return nil
}

func (i Int64) Inspect() string {
	return fmt.Sprintf("%di64", i)
}

func (i Int64) InstanceVariables() SymbolMap {
	return nil
}

func initInt64() {
	Int64Class = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Int64", Int64Class)
}
