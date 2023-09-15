package object

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ALTree/bigfloat"
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
		oFloat, _ := o.ToGoBigInt().Float64()
		return f + Float(oFloat), nil
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
		oFloat, _ := o.ToGoBigInt().Float64()
		return f - Float(oFloat), nil
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
		oFloat, _ := o.ToGoBigInt().Float64()
		return f * Float(oFloat), nil
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
		oFloat, _ := o.ToGoBigInt().Float64()
		return f / Float(oFloat), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Exponentiate by another value and return an error
// if something went wrong.
func (f Float) Exponentiate(other Value) (Value, *Error) {
	switch o := other.(type) {
	case Float:
		return Float(math.Pow(float64(f), float64(o))), nil
	case *BigFloat:
		prec := max(o.Precision(), 53)
		fBigFloat := (&big.Float{}).SetPrec(prec).SetFloat64(float64(f))
		fBigFloat = bigfloat.Pow(fBigFloat, o.ToGoBigFloat())
		return ToElkBigFloat(fBigFloat), nil
	case SmallInt:
		return Float(math.Pow(float64(f), float64(o))), nil
	case *BigInt:
		oFloat, _ := o.ToGoBigInt().Float64()
		return Float(math.Pow(float64(f), oFloat)), nil
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
