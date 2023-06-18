package object

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

func initInt16() {
	Int16Class = NewClass(ClassWithParent(NumericClass), ClassWithImmutable(), ClassWithSealed())
	StdModule.AddConstant("Int16", Int16Class)
}
