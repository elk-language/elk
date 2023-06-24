package object

import (
	"fmt"
	"math/big"
)

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

func (i SmallInt) Inspect() string {
	return fmt.Sprintf("%d", i)
}

func (i SmallInt) InstanceVariables() SimpleSymbolMap {
	return nil
}

var BigIntClass *Class // ::Std::BigInt

// Elk's BigInt value
type BigInt big.Int

// Convert Go big.Int value to Elk BigInt value.
func ToElkBigInt(i *big.Int) *BigInt {
	return (*BigInt)(i)
}

// Convert the Elk BigInt value to Go big.Int value.
func (i *BigInt) ToGoBigInt() *big.Int {
	return (*big.Int)(i)
}

// Negate the number and return the result.
func (i *BigInt) Neg() *BigInt {
	return ToElkBigInt(
		(&big.Int{}).Neg(i.ToGoBigInt()),
	)
}

func (i *BigInt) Class() *Class {
	return BigIntClass
}

func (i *BigInt) IsFrozen() bool {
	return true
}

func (i *BigInt) SetFrozen() {}

func (i *BigInt) Inspect() string {
	return i.ToGoBigInt().String()
}

func (i *BigInt) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initInt() {
	IntClass = NewClass(
		ClassWithParent(NumericClass),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("Int", IntClass)

	SmallIntClass = NewClass(
		ClassWithParent(IntClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithSingleton(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("SmallInt", SmallIntClass)

	BigIntClass = NewClass(
		ClassWithParent(IntClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithSingleton(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("BigInt", BigIntClass)
}
