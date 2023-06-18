package object

var Int64Class *Class // ::Std::Int64

// Elk's Int64 value
type Int64 int64

func (i Int64) Class() *Class {
	return Int64Class
}

func (i Int64) IsFrozen() bool {
	return true
}

func (i Int64) SetFrozen() {}

func initInt64() {
	Int64Class = NewClass(ClassWithParent(NumericClass), ClassWithImmutable(), ClassWithSealed())
	StdModule.AddConstant("Int64", Int64Class)
}
