package value

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/cespare/xxhash/v2"
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

// Cmp compares x and y and returns:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
//
// Panics if x or y are NaN.
func (x Float) Cmp(y Float) int {
	if x.IsNaN() || y.IsNaN() {
		panic("tried to compare NaN Float")
	}

	if x > y {
		return 1
	}
	if x < y {
		return -1
	}
	return 0
}

// Check if the float is an integer.
func (f Float) IsInt() bool {
	return f == Float(int(f))
}

// Convert to Elk Int.
func (f Float) ToInt() Value {
	bigInt := NewBigFloat(float64(f)).ToBigInt()
	if bigInt.IsSmallInt() {
		return bigInt.ToSmallInt()
	}

	return bigInt
}

// Convert to Elk Float64
func (f Float) ToFloat64() Float64 {
	return Float64(f)
}

// Convert to Elk Float32
func (f Float) ToFloat32() Float32 {
	return Float32(f)
}

// Convert to Elk Int64
func (f Float) ToInt64() Int64 {
	return Int64(f)
}

// Convert to Elk Int32
func (f Float) ToInt32() Int32 {
	return Int32(f)
}

// Convert to Elk Int16
func (f Float) ToInt16() Int16 {
	return Int16(f)
}

// Convert to Elk Int8
func (f Float) ToInt8() Int8 {
	return Int8(f)
}

// Convert to Elk UInt64
func (f Float) ToUInt64() UInt64 {
	return UInt64(f)
}

// Convert to Elk UInt32
func (f Float) ToUInt32() UInt32 {
	return UInt32(f)
}

// Convert to Elk UInt16
func (f Float) ToUInt16() UInt16 {
	return UInt16(f)
}

// Convert to Elk UInt8
func (f Float) ToUInt8() UInt8 {
	return UInt8(f)
}

func (f Float) Copy() Value {
	return f
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
	if f.IsInt() {
		return fmt.Sprintf("%.1f", f)
	}
	return fmt.Sprintf("%g", f)
}

func (f Float) InstanceVariables() SymbolMap {
	return nil
}

func (f Float) ToString() String {
	return String(fmt.Sprintf("%g", f))
}

func (f Float) Hash() UInt64 {
	d := xxhash.New()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(float64(f)))
	d.Write(b)
	return UInt64(d.Sum64())
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
		return nil, NewCoerceError(f.Class(), other.Class())
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
		return nil, NewCoerceError(f.Class(), other.Class())
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
		return nil, NewCoerceError(f.Class(), other.Class())
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
		return nil, NewCoerceError(f.Class(), other.Class())
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
		return nil, NewCoerceError(f.Class(), other.Class())
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
		return nil, NewCoerceError(f.Class(), other.Class())
	}
}

// Returns 1 if i is greater than other
// Returns 0 if both are equal.
// Returns -1 if i is less than other.
// Returns nil if the comparison was impossible (NaN)
func (f Float) Compare(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		if f.IsNaN() {
			return Nil, nil
		}
		return SmallInt(f.Cmp(Float(o))), nil
	case *BigInt:
		if f.IsNaN() {
			return Nil, nil
		}
		return SmallInt(f.Cmp(o.ToFloat())), nil
	case Float:
		if f.IsNaN() || o.IsNaN() {
			return Nil, nil
		}
		return SmallInt(f.Cmp(o)), nil
	case *BigFloat:
		if f.IsNaN() || o.IsNaN() {
			return Nil, nil
		}
		iBigFloat := (&BigFloat{}).SetFloat(f)
		return SmallInt(iBigFloat.Cmp(o)), nil
	default:
		return nil, NewCoerceError(f.Class(), other.Class())
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
		return nil, NewCoerceError(f.Class(), other.Class())
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
		return nil, NewCoerceError(f.Class(), other.Class())
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
		return nil, NewCoerceError(f.Class(), other.Class())
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
		return nil, NewCoerceError(f.Class(), other.Class())
	}
}

// Check whether f is equal to other
func (f Float) LaxEqual(other Value) Value {
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

// Check whether f is equal to other and return an error
// if something went wrong.
func (f Float) Equal(other Value) Value {
	switch o := other.(type) {
	case Float:
		return ToElkBool(f == o)
	default:
		return False
	}
}

// Check whether f is strictly equal to other and return an error
// if something went wrong.
func (f Float) StrictEqual(other Value) Value {
	return f.Equal(other)
}

func initFloat() {
	FloatClass = NewClass()
	StdModule.AddConstantString("Float", FloatClass)
	FloatClass.AddConstantString("NAN", FloatNaN())
	FloatClass.AddConstantString("INF", FloatInf())
	FloatClass.AddConstantString("NEG_INF", FloatNegInf())
}
