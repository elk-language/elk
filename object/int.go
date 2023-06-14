package object

// Elk's SmallInt value
type SmallInt int64

func (i SmallInt) Class() *Class {
	return IntClass
}

func (i SmallInt) IsFrozen() bool {
	return true
}

func (i SmallInt) SetFrozen() {}

var IntClass *Class      // ::Std::Int
var SmallIntClass *Class // ::Std::SmallInt
var BigIntClass *Class   // ::Std::BigInt

func initInt() {
	IntClass = NewClass()
	StdModule.AddConstant("Int", IntClass)

	SmallIntClass = NewClass(ClassWithParent(IntClass))
	StdModule.AddConstant("SmallInt", SmallIntClass)

	BigIntClass = NewClass(ClassWithParent(IntClass))
	StdModule.AddConstant("BigInt", BigIntClass)
}
