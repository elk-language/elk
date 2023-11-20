package value

import (
	"fmt"
	"math"

	"github.com/google/go-cmp/cmp"
)

var FloatClass *Class // ::Std::Float

// Elk's Float value
type Float float64

// Positive infinity
func FloatInf() Float {
	return Float(math.Inf(1))
}

// Negative infinity
func FloatNegInf() Float {
	return Float(math.Inf(-1))
}

// Not a number
func FloatNaN() Float {
	return Float(math.NaN())
}

func (Float) Class() *Class {
	return FloatClass
}

func (Float) DirectClass() *Class {
	return FloatClass
}

func (Float) SingletonClass() *Class {
	return nil
}

// IsNaN reports whether f is a “not-a-number” value.
func (f Float) IsNaN() bool {
	return math.IsNaN(float64(f))
}

// IsInf reports whether f is an infinity, according to sign.
// If sign > 0, IsInf reports whether f is positive infinity.
// If sign < 0, IsInf reports whether f is negative infinity.
// If sign == 0, IsInf reports whether f is either infinity.
func (f Float) IsInf(sign int) bool {
	return math.IsInf(float64(f), sign)
}

func (f Float) Inspect() string {
	if f.IsNaN() {
		return fmt.Sprintf("%s::NAN", f.Class().PrintableName())
	}
	if f.IsInf(1) {
		return fmt.Sprintf("%s::INF", f.Class().PrintableName())
	}
	if f.IsInf(-1) {
		return fmt.Sprintf("%s::NEG_INF", f.Class().PrintableName())
	}
	return fmt.Sprintf("%g", f)
}

func (f Float) InstanceVariables() SymbolMap {
	return nil
}

// Add another value and return an error
// if something went wrong.
func (f Float) Add(other Value) (Value, *Error) {
	switch o := other.(type) {
	case Float:
		return f + o, nil
	case *BigFloat:
		fBigFloat := NewBigFloat(float64(f))
		return fBigFloat.AddBigFloat(fBigFloat, o), nil
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
		fBigFloat := NewBigFloat(float64(f))
		return fBigFloat.SubBigFloat(fBigFloat, o), nil
	case SmallInt:
		return f - Float(o), nil
	case *BigInt:
		return f - o.ToFloat(), nil
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
		fBigFloat := NewBigFloat(float64(f))
		return fBigFloat.MulBigFloat(fBigFloat, o), nil
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
		fBigFloat := NewBigFloat(float64(f))
		return fBigFloat.DivBigFloat(fBigFloat, o), nil
	case SmallInt:
		return f / Float(o), nil
	case *BigInt:
		return f / o.ToFloat(), nil
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
		fBigFloat := (&BigFloat{}).SetPrecision(prec).SetFloat(f)
		fBigFloat.ExpBigFloat(fBigFloat, o)
		return fBigFloat, nil
	case SmallInt:
		return Float(math.Pow(float64(f), float64(o))), nil
	case *BigInt:
		oFloat, _ := o.ToGoBigInt().Float64()
		return Float(math.Pow(float64(f), oFloat)), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

func (a Float) Mod(b Float) Float {
	return Float(math.Mod(float64(a), float64(b)))
}

// Perform modulo by another numeric value and return an error
// if something went wrong.
func (f Float) Modulo(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return f.Mod(Float(o)), nil
	case *BigInt:
		return f.Mod(o.ToFloat()), nil
	case Float:
		return f.Mod(o), nil
	case *BigFloat:
		prec := max(o.Precision(), 53)
		fBigFloat := (&BigFloat{}).SetPrecision(prec).SetFloat(f)
		return fBigFloat.Mod(fBigFloat, o), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Check whether f is greater than other and return an error
// if something went wrong.
func (f Float) GreaterThan(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return ToElkBool(f > Float(o)), nil
	case *BigInt:
		return ToElkBool(f > o.ToFloat()), nil
	case Float:
		return ToElkBool(f > o), nil
	case *BigFloat:
		if f.IsNaN() || o.IsNaN() {
			return False, nil
		}
		iBigFloat := (&BigFloat{}).SetFloat(f)
		return ToElkBool(iBigFloat.Cmp(o) == 1), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Check whether f is greater than or equal to other and return an error
// if something went wrong.
func (f Float) GreaterThanEqual(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return ToElkBool(f >= Float(o)), nil
	case *BigInt:
		return ToElkBool(f >= o.ToFloat()), nil
	case Float:
		return ToElkBool(f >= o), nil
	case *BigFloat:
		if f.IsNaN() || o.IsNaN() {
			return False, nil
		}
		iBigFloat := (&BigFloat{}).SetFloat(f)
		return ToElkBool(iBigFloat.Cmp(o) >= 0), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Check whether f is less than other and return an error
// if something went wrong.
func (f Float) LessThan(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return ToElkBool(f < Float(o)), nil
	case *BigInt:
		return ToElkBool(f < o.ToFloat()), nil
	case Float:
		return ToElkBool(f < o), nil
	case *BigFloat:
		if f.IsNaN() || o.IsNaN() {
			return False, nil
		}
		iBigFloat := (&BigFloat{}).SetFloat(f)
		return ToElkBool(iBigFloat.Cmp(o) == -1), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Check whether f is less than or equal to other and return an error
// if something went wrong.
func (f Float) LessThanEqual(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return ToElkBool(f <= Float(o)), nil
	case *BigInt:
		return ToElkBool(f <= o.ToFloat()), nil
	case Float:
		return ToElkBool(f <= o), nil
	case *BigFloat:
		if f.IsNaN() || o.IsNaN() {
			return False, nil
		}
		iBigFloat := (&BigFloat{}).SetFloat(f)
		return ToElkBool(iBigFloat.Cmp(o) <= 0), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Check whether f is equal to other
func (f Float) Equal(other Value) Value {
	switch o := other.(type) {
	case SmallInt:
		return ToElkBool(f == Float(o))
	case *BigInt:
		return ToElkBool(f == o.ToFloat())
	case Float:
		return ToElkBool(f == o)
	case *BigFloat:
		if f.IsNaN() || o.IsNaN() {
			return False
		}
		fBigFloat := (&BigFloat{}).SetFloat(f)
		return ToElkBool(fBigFloat.Cmp(o) == 0)
	case Int64:
		return ToElkBool(f == Float(o))
	case Int32:
		return ToElkBool(f == Float(o))
	case Int16:
		return ToElkBool(f == Float(o))
	case Int8:
		return ToElkBool(f == Float(o))
	case UInt64:
		return ToElkBool(f == Float(o))
	case UInt32:
		return ToElkBool(f == Float(o))
	case UInt16:
		return ToElkBool(f == Float(o))
	case UInt8:
		return ToElkBool(f == Float(o))
	case Float64:
		return ToElkBool(float64(f) == float64(o))
	case Float32:
		return ToElkBool(float64(f) == float64(o))
	default:
		return False
	}
}

// Check whether f is strictly equal to other and return an error
// if something went wrong.
func (f Float) StrictEqual(other Value) Value {
	switch o := other.(type) {
	case Float:
		return ToElkBool(f == o)
	default:
		return False
	}
}

var FloatComparer cmp.Option

func initFloatComparer() {
	FloatComparer = cmp.Comparer(func(x, y Float) bool {
		if x.IsNaN() || y.IsNaN() {
			return x.IsNaN() && y.IsNaN()
		}
		return x == y
	})
}

func initFloat() {
	FloatClass = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	StdModule.AddConstantString("Float", FloatClass)
	FloatClass.AddConstantString("NAN", FloatNaN())
	FloatClass.AddConstantString("INF", FloatInf())
	FloatClass.AddConstantString("NEG_INF", FloatNegInf())
}
