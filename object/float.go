package object

import (
	"fmt"
	"math/big"
)

var FloatClass *Class // ::Std::Float

// Elk's Float value
type Float float64

func (Float) Class() *Class {
	return FloatClass
}

func (Float) IsFrozen() bool {
	return true
}

func (Float) SetFrozen() {}

func (f Float) Inspect() string {
	return fmt.Sprintf("%g", f)
}

func (f Float) InstanceVariables() SimpleSymbolMap {
	return nil
}

// Add another value and return an error
// if something went wrong.
func (f Float) Add(other Value) (Value, *Error) {
	switch o := other.(type) {
	case Float:
		return f + o, nil
	case *BigFloat:
		fBigFloat := big.NewFloat(float64(f))
		return ToElkBigFloat(fBigFloat.Add(fBigFloat, o.ToGoBigFloat())), nil
	case SmallInt:
		return f + Float(o), nil
	case *BigInt:
		fBigFloat := big.NewFloat(float64(f))
		otherBigFloat := (&big.Float{}).SetInt(o.ToGoBigInt())
		result, _ := fBigFloat.Add(fBigFloat, otherBigFloat).Float64()
		return Float(result), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Subtract another value and return an error
// if something went wrong.
func (f Float) Subtract(other Value) (Value, *Error) {
	switch o := other.(type) {
	case Float:
		return f - o, nil
	case *BigFloat:
		fBigFloat := big.NewFloat(float64(f))
		return ToElkBigFloat(fBigFloat.Sub(fBigFloat, o.ToGoBigFloat())), nil
	case SmallInt:
		return f - Float(o), nil
	case *BigInt:
		fBigFloat := big.NewFloat(float64(f))
		otherBigFloat := (&big.Float{}).SetInt(o.ToGoBigInt())
		result, _ := fBigFloat.Sub(fBigFloat, otherBigFloat).Float64()
		return Float(result), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Add another value and return an error
// if something went wrong.
func (f Float) Multiply(other Value) (Value, *Error) {
	switch o := other.(type) {
	case Float:
		return f * o, nil
	case *BigFloat:
		fBigFloat := big.NewFloat(float64(f))
		return ToElkBigFloat(fBigFloat.Mul(fBigFloat, o.ToGoBigFloat())), nil
	case SmallInt:
		return f * Float(o), nil
	case *BigInt:
		fBigFloat := big.NewFloat(float64(f))
		otherBigFloat := (&big.Float{}).SetInt(o.ToGoBigInt())
		result, _ := fBigFloat.Mul(fBigFloat, otherBigFloat).Float64()
		return Float(result), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Divide by another value and return an error
// if something went wrong.
func (f Float) Divide(other Value) (Value, *Error) {
	switch o := other.(type) {
	case Float:
		return f / o, nil
	case *BigFloat:
		fBigFloat := big.NewFloat(float64(f))
		return ToElkBigFloat(fBigFloat.Quo(fBigFloat, o.ToGoBigFloat())), nil
	case SmallInt:
		return f / Float(o), nil
	case *BigInt:
		fBigFloat := big.NewFloat(float64(f))
		otherBigFloat := (&big.Float{}).SetInt(o.ToGoBigInt())
		result, _ := fBigFloat.Quo(fBigFloat, otherBigFloat).Float64()
		return Float(result), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

func initFloat() {
	FloatClass = NewClass(
		ClassWithParent(NumericClass),
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	StdModule.AddConstant("Float", FloatClass)
}
