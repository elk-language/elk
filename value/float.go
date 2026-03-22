package value

import (
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"

	"github.com/cespare/xxhash/v2"
)

var FloatClass *Class // ::Std::Float

// Elk's Float value

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

func (f Float) ToValue() Value {
	return Value{
		flag: FLOAT_FLAG,
		data: *(*uintptr)(unsafe.Pointer(&f)),
	}
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
		return bigInt.ToSmallInt().ToValue()
	}

	return Ref(bigInt)
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
func (f Float) ToUInt() UInt {
	return UInt(f)
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

func (f Float) Error() string {
	return f.Inspect()
}

func (f Float) InstanceVariables() *InstanceVariables {
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

// AddVal another value and return an error
// if something went wrong.
func (f Float) AddVal(other Value) (result, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return Ref(f.AddBigFloat(o)), Undefined
		case *BigInt:
			return f.AddBigInt(o).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.AddFloat(other.AsFloat()).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return f.AddSmallInt(other.AsSmallInt()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

func (f Float) AddBigFloat(other *BigFloat) *BigFloat {
	fBigFloat := NewBigFloat(float64(f))
	return fBigFloat.AddBigFloat(fBigFloat, other)
}

func (f Float) AddFloat(other Float) Float {
	return f + other
}

func (f Float) AddSmallInt(other SmallInt) Float {
	return f + Float(other)
}

func (f Float) AddInt(other Value) Float {
	if other.IsSmallInt() {
		return f.AddSmallInt(other.AsSmallInt())
	}
	return f.AddBigInt((*BigInt)(other.Pointer()))
}

func (f Float) AddBigInt(other *BigInt) Float {
	oFloat, _ := other.ToGoBigInt().Float64()
	return f + Float(oFloat)
}

// SubtractVal another value and return an error
// if something went wrong.
func (f Float) SubtractVal(other Value) (result, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return Ref(f.SubtractBigFloat(o)), Undefined
		case *BigInt:
			return f.SubtractBigInt(o).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return (f.SubtractFloat(other.AsFloat())).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return (f.SubtractSmallInt(other.AsSmallInt())).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

func (f Float) SubtractBigFloat(other *BigFloat) *BigFloat {
	fBigFloat := NewBigFloat(float64(f))
	return fBigFloat.SubBigFloat(fBigFloat, other)
}

func (f Float) SubtractFloat(other Float) Float {
	return f - other
}

func (f Float) SubtractSmallInt(other SmallInt) Float {
	return f - Float(other)
}

func (f Float) SubtractInt(other Value) Float {
	if other.IsSmallInt() {
		return f.SubtractSmallInt(other.AsSmallInt())
	}
	return f.SubtractBigInt((*BigInt)(other.Pointer()))
}

func (f Float) SubtractBigInt(other *BigInt) Float {
	oFloat, _ := other.ToGoBigInt().Float64()
	return f - Float(oFloat)
}

// Add another value and return an error
// if something went wrong.
func (f Float) MultiplyVal(other Value) (result, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return Ref(f.MultiplyBigFloat(o)), Undefined
		case *BigInt:
			return f.MultiplyBigInt(o).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return (f.MultiplyFloat(other.AsFloat())).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return (f.MultiplySmallInt(other.AsSmallInt())).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

func (f Float) MultiplyBigFloat(other *BigFloat) *BigFloat {
	fBigFloat := NewBigFloat(float64(f))
	return fBigFloat.MulBigFloat(fBigFloat, other)
}

func (f Float) MultiplyFloat(other Float) Float {
	return f * other
}

func (f Float) MultiplySmallInt(other SmallInt) Float {
	return f * Float(other)
}

func (f Float) MultiplyInt(other Value) Float {
	if other.IsSmallInt() {
		return f.MultiplySmallInt(other.AsSmallInt())
	}
	return f.MultiplyBigInt((*BigInt)(other.Pointer()))
}

func (f Float) MultiplyBigInt(other *BigInt) Float {
	oFloat, _ := other.ToGoBigInt().Float64()
	return f * Float(oFloat)
}

// DivideVal by another value and return an error
// if something went wrong.
func (f Float) DivideVal(other Value) (result, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return Ref(f.DivideBigFloat(o)), Undefined
		case *BigInt:
			return (f.DivideBigInt(o)).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return (f.DivideFloat(other.AsFloat())).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return (f.DivideSmallInt(other.AsSmallInt())).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

func (f Float) DivideBigFloat(other *BigFloat) *BigFloat {
	fBigFloat := NewBigFloat(float64(f))
	return fBigFloat.DivBigFloat(fBigFloat, other)
}

func (f Float) DivideFloat(other Float) Float {
	return f / other
}

func (f Float) DivideSmallInt(other SmallInt) Float {
	return f / Float(other)
}

func (f Float) DivideInt(other Value) Float {
	if other.IsSmallInt() {
		return f.DivideSmallInt(other.AsSmallInt())
	}
	return f.DivideBigInt((*BigInt)(other.Pointer()))
}

func (f Float) DivideBigInt(other *BigInt) Float {
	oFloat, _ := other.ToGoBigInt().Float64()
	return f / Float(oFloat)
}

// ExponentiateVal by another value and return an error
// if something went wrong.
func (f Float) ExponentiateVal(other Value) (result, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return Ref(f.ExponentiateBigFloat(o)), Undefined
		case *BigInt:
			return f.ExponentiateBigInt(o).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.ExponentiateFloat(other.AsFloat()).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return f.ExponentiateSmallInt(other.AsSmallInt()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

func (f Float) ExponentiateBigFloat(other *BigFloat) *BigFloat {
	prec := max(other.Precision(), 53)
	fBigFloat := (&BigFloat{}).SetPrecision(prec).SetFloat(f)
	fBigFloat.ExpBigFloat(fBigFloat, other)
	return fBigFloat
}

func (f Float) ExponentiateFloat(other Float) Float {
	return Float(math.Pow(float64(f), float64(other)))
}

func (f Float) ExponentiateSmallInt(other SmallInt) Float {
	return Float(math.Pow(float64(f), float64(other)))
}

func (f Float) ExponentiateInt(other Value) Float {
	if other.IsSmallInt() {
		return f.ExponentiateSmallInt(other.AsSmallInt())
	}
	return f.ExponentiateBigInt((*BigInt)(other.Pointer()))
}

func (f Float) ExponentiateBigInt(other *BigInt) Float {
	oFloat, _ := other.ToGoBigInt().Float64()
	return Float(math.Pow(float64(f), oFloat))
}

func (a Float) Mod(b Float) Float {
	return Float(math.Mod(float64(a), float64(b)))
}

// Perform modulo by another numeric value and return an error
// if something went wrong.
func (f Float) ModuloVal(other Value) (result, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return Ref(f.ModuloBigFloat(o)), Undefined
		case *BigInt:
			return f.ModuloBigInt(o).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return (f.ModuloFloat(other.AsFloat())).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return (f.ModuloSmallInt(other.AsSmallInt())).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

func (f Float) ModuloBigFloat(other *BigFloat) *BigFloat {
	prec := max(other.Precision(), 53)
	fBigFloat := (&BigFloat{}).SetPrecision(prec).SetFloat(f)
	return fBigFloat.Mod(fBigFloat, other)
}

func (f Float) ModuloFloat(other Float) Float {
	return f.Mod(other)
}

func (f Float) ModuloSmallInt(other SmallInt) Float {
	return f.Mod(Float(other))
}

func (f Float) ModuloInt(other Value) Float {
	if other.IsSmallInt() {
		return f.ModuloSmallInt(other.AsSmallInt())
	}
	return f.ModuloBigInt((*BigInt)(other.Pointer()))
}

func (f Float) ModuloBigInt(other *BigInt) Float {
	oFloat, _ := other.ToGoBigInt().Float64()
	return f.Mod(Float(oFloat))
}

// Returns 1 if i is greater than other
// Returns 0 if both are equal.
// Returns -1 if i is less than other.
// Returns nil if the comparison was impossible (NaN)
func (f Float) CompareVal(other Value) (result, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.CompareBigFloat(o), Undefined
		case *BigInt:
			return f.CompareBigInt(o), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.CompareFloat(other.AsFloat()), Undefined
	case SMALL_INT_FLAG:
		return f.CompareSmallInt(other.AsSmallInt()), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

func (f Float) CompareBigFloat(other *BigFloat) Value {
	if f.IsNaN() || other.IsNaN() {
		return Nil
	}
	iBigFloat := (&BigFloat{}).SetFloat(f)
	return SmallInt(iBigFloat.Cmp(other)).ToValue()
}

func (f Float) CompareFloat(other Float) Value {
	if f.IsNaN() || other.IsNaN() {
		return Nil
	}
	return SmallInt(f.Cmp(other)).ToValue()
}

func (f Float) CompareSmallInt(other SmallInt) Value {
	if f.IsNaN() {
		return Nil
	}
	return SmallInt(f.Cmp(Float(other))).ToValue()
}

func (f Float) CompareInt(other Value) Value {
	if other.IsSmallInt() {
		return f.CompareSmallInt(other.AsSmallInt())
	}
	return f.CompareBigInt((*BigInt)(other.Pointer()))
}

func (f Float) CompareBigInt(other *BigInt) Value {
	if f.IsNaN() {
		return Nil
	}
	return SmallInt(f.Cmp(other.ToFloat())).ToValue()
}

// Check whether f is greater than other and return an error
// if something went wrong.
func (f Float) GreaterThanVal(other Value) (Value, Value) {
	result, err := f.GreaterThan(other)
	return Bool(result).ToValue(), err
}

// Check whether f is greater than other and return an error
// if something went wrong.
func (f Float) GreaterThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.GreaterThanBigFloat(o), Undefined
		case *BigInt:
			return f.GreaterThanBigInt(o), Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.GreaterThanFloat(other.AsFloat()), Undefined
	case SMALL_INT_FLAG:
		return f.GreaterThanSmallInt(other.AsSmallInt()), Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

func (f Float) GreaterThanBigFloat(other *BigFloat) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}
	fBigFloat := (&BigFloat{}).SetFloat(f)
	return fBigFloat.Cmp(other) == 1
}

func (f Float) GreaterThanFloat(other Float) bool {
	return f > other
}

func (f Float) GreaterThanSmallInt(other SmallInt) bool {
	return f > Float(other)
}

func (f Float) GreaterThanInt(other Value) bool {
	if other.IsSmallInt() {
		return f.GreaterThanSmallInt(other.AsSmallInt())
	}
	return f.GreaterThanBigInt((*BigInt)(other.Pointer()))
}

func (f Float) GreaterThanBigInt(other *BigInt) bool {
	oFloat := other.ToFloat()
	return f > oFloat
}

// Check whether f is greater than or equal to other and return an error
// if something went wrong.
func (f Float) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := f.GreaterThanEqual(other)
	return Bool(result).ToValue(), err
}

// Check whether f is greater than or equal to other and return an error
// if something went wrong.
func (f Float) GreaterThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.GreaterThanEqualBigFloat(o), Undefined
		case *BigInt:
			return f.GreaterThanEqualBigInt(o), Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.GreaterThanEqualFloat(other.AsFloat()), Undefined
	case SMALL_INT_FLAG:
		return f.GreaterThanEqualSmallInt(other.AsSmallInt()), Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

func (f Float) GreaterThanEqualBigFloat(other *BigFloat) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}
	fBigFloat := (&BigFloat{}).SetFloat(f)
	return fBigFloat.Cmp(other) >= 0
}

func (f Float) GreaterThanEqualFloat(other Float) bool {
	return f >= other
}

func (f Float) GreaterThanEqualSmallInt(other SmallInt) bool {
	return f >= Float(other)
}

func (f Float) GreaterThanEqualInt(other Value) bool {
	if other.IsSmallInt() {
		return f.GreaterThanEqualSmallInt(other.AsSmallInt())
	}
	return f.GreaterThanEqualBigInt((*BigInt)(other.Pointer()))
}

func (f Float) GreaterThanEqualBigInt(other *BigInt) bool {
	oFloat := other.ToFloat()
	return f >= oFloat
}

// Check whether f is less than other and return an error
// if something went wrong.
func (f Float) LessThanVal(other Value) (Value, Value) {
	result, err := f.LessThan(other)
	return Bool(result).ToValue(), err
}

// Check whether f is less than other and return an error
// if something went wrong.
func (f Float) LessThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.LessThanBigFloat(o), Undefined
		case *BigInt:
			return f.LessThanBigInt(o), Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.LessThanFloat(other.AsFloat()), Undefined
	case SMALL_INT_FLAG:
		return f.LessThanSmallInt(other.AsSmallInt()), Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

func (f Float) LessThanBigFloat(other *BigFloat) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}
	fBigFloat := (&BigFloat{}).SetFloat(f)
	return fBigFloat.Cmp(other) == -1
}

func (f Float) LessThanFloat(other Float) bool {
	return f < other
}

func (f Float) LessThanSmallInt(other SmallInt) bool {
	return f < Float(other)
}

func (f Float) LessThanInt(other Value) bool {
	if other.IsSmallInt() {
		return f.LessThanSmallInt(other.AsSmallInt())
	}
	return f.LessThanBigInt((*BigInt)(other.Pointer()))
}

func (f Float) LessThanBigInt(other *BigInt) bool {
	oFloat := other.ToFloat()
	return f < oFloat
}

// Check whether f is less than or equal to other and return an error
// if something went wrong.
func (f Float) LessThanEqualVal(other Value) (Value, Value) {
	result, err := f.LessThanEqual(other)
	return Bool(result).ToValue(), err
}

// Check whether f is less than or equal to other and return an error
// if something went wrong.
func (f Float) LessThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.LessThanEqualBigFloat(o), Undefined
		case *BigInt:
			return f.LessThanEqualBigInt(o), Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.LessThanEqualFloat(other.AsFloat()), Undefined
	case SMALL_INT_FLAG:
		return f.LessThanEqualSmallInt(other.AsSmallInt()), Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

func (f Float) LessThanEqualBigFloat(other *BigFloat) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}
	fBigFloat := (&BigFloat{}).SetFloat(f)
	return fBigFloat.Cmp(other) <= 0
}

func (f Float) LessThanEqualFloat(other Float) bool {
	return f <= other
}

func (f Float) LessThanEqualSmallInt(other SmallInt) bool {
	return f <= Float(other)
}

func (f Float) LessThanEqualInt(other Value) bool {
	if other.IsSmallInt() {
		return f.LessThanEqualSmallInt(other.AsSmallInt())
	}
	return f.LessThanEqualBigInt((*BigInt)(other.Pointer()))
}

func (f Float) LessThanEqualBigInt(other *BigInt) bool {
	oFloat := other.ToFloat()
	return f <= oFloat
}

// Check whether f is equal to other
func (f Float) LaxEqualVal(other Value) Value {
	return BoolVal(f.LaxEqual(other))
}

func (f Float) LaxEqual(other Value) bool {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return f == o.ToFloat()
		case *BigFloat:
			if f.IsNaN() || o.IsNaN() {
				return false
			}
			fBigFloat := (&BigFloat{}).SetFloat(f)
			return fBigFloat.Cmp(o) == 0
		case Int64:
			return f == Float(o)
		case UInt64:
			return f == Float(o)
		case Float64:
			return float64(f) == float64(o)
		default:
			return false
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return f == Float(other.AsSmallInt())
	case FLOAT_FLAG:
		return f == other.AsFloat()
	case INT64_FLAG:
		return f == Float(other.AsInlineInt64())
	case INT32_FLAG:
		return f == Float(other.AsInt32())
	case INT16_FLAG:
		return f == Float(other.AsInt16())
	case INT8_FLAG:
		return f == Float(other.AsInt8())
	case UINT64_FLAG:
		return f == Float(other.AsInlineUInt64())
	case UINT32_FLAG:
		return f == Float(other.AsUInt32())
	case UINT16_FLAG:
		return f == Float(other.AsUInt16())
	case UINT8_FLAG:
		return f == Float(other.AsUInt8())
	case FLOAT64_FLAG:
		return float64(f) == float64(other.AsInlineFloat64())
	case FLOAT32_FLAG:
		return Float(f) == Float(other.AsFloat32())
	default:
		return false
	}
}

// Check whether f is equal to other and return an error
// if something went wrong.
func (f Float) EqualVal(other Value) Value {
	return Bool(f.Equal(other)).ToValue()
}

// Check whether f is equal to other and return an error
// if something went wrong.
func (f Float) Equal(other Value) bool {
	if other.IsFloat() {
		return f == other.AsFloat()
	}

	return false
}

func (f Float) EqualFloat(other Float) bool {
	return f == other
}

// Check whether f is strictly equal to other and return an error
// if something went wrong.
func (f Float) StrictEqualVal(other Value) Value {
	return f.EqualVal(other)
}

func (f Float) Nanoseconds() TimeSpan {
	return TimeSpan(f)
}

func (f Float) Microseconds() TimeSpan {
	return TimeSpan(f * Float(Microsecond))
}

func (f Float) Milliseconds() TimeSpan {
	return TimeSpan(f * Float(Millisecond))
}

func (f Float) Seconds() TimeSpan {
	return TimeSpan(f * Float(Second))
}

func (f Float) Minutes() TimeSpan {
	return TimeSpan(f * Float(Minute))
}

func (f Float) Hours() TimeSpan {
	return TimeSpan(f * Float(Hour))
}

func (f Float) Days() *DateTimeSpan {
	days, frac := math.Modf(float64(f))

	return NewDateTimeSpan(
		MakeDateSpan(0, 0, int(days)),
		TimeSpan(frac*float64(Day)),
	)
}

func (f Float) Weeks() *DateTimeSpan {
	days, frac := math.Modf(float64(f * 7))

	return NewDateTimeSpan(
		MakeDateSpan(0, 0, int(days)),
		TimeSpan(frac*float64(Day)),
	)
}

func (f Float) Months() *DateTimeSpan {
	months, frac := math.Modf(float64(f))
	days, frac := math.Modf(frac * MonthDays)

	return NewDateTimeSpan(
		MakeDateSpan(0, int(months), int(days)),
		TimeSpan(frac*float64(Day)),
	)
}

func (f Float) Years() *DateTimeSpan {
	years, frac := math.Modf(float64(f))
	months, frac := math.Modf(frac * 12)
	days, frac := math.Modf(frac * MonthDays)

	return NewDateTimeSpan(
		MakeDateSpan(int(years), int(months), int(days)),
		TimeSpan(frac*float64(Day)),
	)
}

func (f Float) Centuries() *DateTimeSpan {
	years, frac := math.Modf(float64(f) * 100)
	months, frac := math.Modf(frac * 12)
	days, frac := math.Modf(frac * MonthDays)

	return NewDateTimeSpan(
		MakeDateSpan(int(years), int(months), int(days)),
		TimeSpan(frac*float64(Day)),
	)
}

func (f Float) Millenia() *DateTimeSpan {
	years, frac := math.Modf(float64(f) * 1000)
	months, frac := math.Modf(frac * 12)
	days, frac := math.Modf(frac * MonthDays)

	return NewDateTimeSpan(
		MakeDateSpan(int(years), int(months), int(days)),
		TimeSpan(frac*float64(Day)),
	)
}

func initFloat() {
	FloatClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	StdModule.AddConstantString("Float", Ref(FloatClass))
	RegisterNativeClass("Std::Float", "value.FloatClass")

	FloatClass.AddConstantString("NAN", FloatNaN().ToValue())
	RegisterNativeConstant("Std::Float::NAN", "value.FloatNaN()", NewGoType("value.Float"))

	FloatClass.AddConstantString("INF", FloatInf().ToValue())
	RegisterNativeConstant("Std::Float::INF", "value.FloatInf()", NewGoType("value.Float"))

	FloatClass.AddConstantString("NEG_INF", FloatNegInf().ToValue())
	RegisterNativeConstant("Std::Float::NEG_INF", "value.FloatNegInf()", NewGoType("value.Float"))

	FloatClass.AddConstantString("Convertible", Ref(NewInterface()))
}
