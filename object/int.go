package object

import "math/big"

var IntClass *Class // ::Std::Int

var SmallIntClass *Class // ::Std::SmallInt

// Elk's SmallInt value
type SmallInt int64

func (i SmallInt) Class() *Class {
	return SmallIntClass
}

func (i SmallInt) IsFrozen() bool {
	return true
}

func (i SmallInt) SetFrozen() {}

var BigIntClass *Class // ::Std::BigInt

// Elk's BigInt value
type BigInt struct {
	big.Int
}

func (i *BigInt) Class() *Class {
	return BigIntClass
}

func (i *BigInt) IsFrozen() bool {
	return true
}

func (i *BigInt) SetFrozen() {}

func initInt() {
	IntClass = NewClass(ClassWithParent(NumericClass))
	StdModule.AddConstant("Int", IntClass)

	SmallIntClass = NewClass(ClassWithParent(IntClass), ClassWithImmutable(), ClassWithSealed(), ClassWithSingleton())
	StdModule.AddConstant("SmallInt", SmallIntClass)

	BigIntClass = NewClass(ClassWithParent(IntClass), ClassWithImmutable(), ClassWithSealed(), ClassWithSingleton())
	StdModule.AddConstant("BigInt", BigIntClass)
}
