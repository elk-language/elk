package object

import (
	"fmt"
	"math/big"
)

var BigFloatClass *Class // ::Std::BigFloat

// Elk's BigFloat value
type BigFloat big.Float

// Convert Go's big.Float values to Elk's BigFloat values.
func ToElkBigFloat(f *big.Float) *BigFloat {
	return (*BigFloat)(f)
}

// Create a new BigFloat with the specified value.
func NewBigFloat(f float64) *BigFloat {
	return ToElkBigFloat(big.NewFloat(f))
}

// Convert Elk's BigFloat values to Go's big.Float values.
func (f *BigFloat) ToGoBigFloat() *big.Float {
	return (*big.Float)(f)
}

// Sets the f's precision to prec and possibly
// rounds the value.
func (f *BigFloat) SetPrecision(prec uint) *BigFloat {
	return ToElkBigFloat(f.ToGoBigFloat().SetPrec(prec))
}

// Negate the number and return the result.
func (f *BigFloat) Negate() *BigFloat {
	return ToElkBigFloat(
		(&big.Float{}).Neg(f.ToGoBigFloat()),
	)
}

func (*BigFloat) Class() *Class {
	return BigFloatClass
}

func (*BigFloat) IsFrozen() bool {
	return true
}

func (*BigFloat) SetFrozen() {}

func (f *BigFloat) Inspect() string {
	return fmt.Sprintf("%sbf", f.ToGoBigFloat().String())
}

func (f *BigFloat) InstanceVariables() SimpleSymbolMap {
	return nil
}

// Add another value and return an error
// if something went wrong.
func (f *BigFloat) Add(other Value) (Value, *Error) {
	switch o := other.(type) {
	case *BigFloat:
		result := ToElkBigFloat((&big.Float{}).Add(f.ToGoBigFloat(), o.ToGoBigFloat()))
		return result, nil
	case Float:
		otherBigFloat := big.NewFloat(float64(o))
		return ToElkBigFloat(otherBigFloat.Add(f.ToGoBigFloat(), otherBigFloat)), nil
	case SmallInt:
		otherBigFloat := (&big.Float{}).SetInt64(int64(o))
		result := ToElkBigFloat(otherBigFloat.Add(f.ToGoBigFloat(), otherBigFloat))
		return result, nil
	case *BigInt:
		otherBigFloat := (&big.Float{}).SetInt(o.ToGoBigInt())
		result := ToElkBigFloat(otherBigFloat.Add(f.ToGoBigFloat(), otherBigFloat))
		return result, nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Subtract another value and return an error
// if something went wrong.
func (f *BigFloat) Subtract(other Value) (Value, *Error) {
	switch o := other.(type) {
	case *BigFloat:
		result := ToElkBigFloat((&big.Float{}).Sub(f.ToGoBigFloat(), o.ToGoBigFloat()))
		return result, nil
	case Float:
		otherBigFloat := big.NewFloat(float64(o))
		return ToElkBigFloat(otherBigFloat.Sub(f.ToGoBigFloat(), otherBigFloat)), nil
	case SmallInt:
		otherBigFloat := (&big.Float{}).SetInt64(int64(o))
		result := ToElkBigFloat(otherBigFloat.Sub(f.ToGoBigFloat(), otherBigFloat))
		return result, nil
	case *BigInt:
		otherBigFloat := (&big.Float{}).SetInt(o.ToGoBigInt())
		result := ToElkBigFloat(otherBigFloat.Sub(f.ToGoBigFloat(), otherBigFloat))
		return result, nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Multiply by another value and return an error
// if something went wrong.
func (f *BigFloat) Multiply(other Value) (Value, *Error) {
	switch o := other.(type) {
	case *BigFloat:
		result := ToElkBigFloat((&big.Float{}).Mul(f.ToGoBigFloat(), o.ToGoBigFloat()))
		return result, nil
	case Float:
		otherBigFloat := big.NewFloat(float64(o))
		return ToElkBigFloat(otherBigFloat.Mul(f.ToGoBigFloat(), otherBigFloat)), nil
	case SmallInt:
		otherBigFloat := (&big.Float{}).SetInt64(int64(o))
		result := ToElkBigFloat(otherBigFloat.Mul(f.ToGoBigFloat(), otherBigFloat))
		return result, nil
	case *BigInt:
		otherBigFloat := (&big.Float{}).SetInt(o.ToGoBigInt())
		result := ToElkBigFloat(otherBigFloat.Mul(f.ToGoBigFloat(), otherBigFloat))
		return result, nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

func initBigFloat() {
	BigFloatClass = NewClass(
		ClassWithParent(NumericClass),
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	StdModule.AddConstant("BigFloat", BigFloatClass)
}
