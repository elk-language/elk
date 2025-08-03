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
			fBigFloat := NewBigFloat(float64(f))
			return Ref(fBigFloat.AddBigFloat(fBigFloat, o)), Undefined
		case *BigInt:
			oFloat, _ := o.ToGoBigInt().Float64()
			return (f + Float(oFloat)).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return (f + other.AsFloat()).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return (f + Float(other.AsSmallInt())).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// SubtractVal another value and return an error
// if something went wrong.
func (f Float) SubtractVal(other Value) (result, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			fBigFloat := NewBigFloat(float64(f))
			return Ref(fBigFloat.SubBigFloat(fBigFloat, o)), Undefined
		case *BigInt:
			oFloat, _ := o.ToGoBigInt().Float64()
			return (f - Float(oFloat)).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return (f - other.AsFloat()).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return (f - Float(other.AsSmallInt())).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Add another value and return an error
// if something went wrong.
func (f Float) MultiplyVal(other Value) (result, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			fBigFloat := NewBigFloat(float64(f))
			return Ref(fBigFloat.MulBigFloat(fBigFloat, o)), Undefined
		case *BigInt:
			oFloat, _ := o.ToGoBigInt().Float64()
			return (f * Float(oFloat)).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return (f * other.AsFloat()).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return (f * Float(other.AsSmallInt())).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// DivideVal by another value and return an error
// if something went wrong.
func (f Float) DivideVal(other Value) (result, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			fBigFloat := NewBigFloat(float64(f))
			return Ref(fBigFloat.DivBigFloat(fBigFloat, o)), Undefined
		case *BigInt:
			oFloat, _ := o.ToGoBigInt().Float64()
			return (f / Float(oFloat)).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return (f / other.AsFloat()).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return (f / Float(other.AsSmallInt())).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// ExponentiateVal by another value and return an error
// if something went wrong.
func (f Float) ExponentiateVal(other Value) (result, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			prec := max(o.Precision(), 53)
			fBigFloat := (&BigFloat{}).SetPrecision(prec).SetFloat(f)
			fBigFloat.ExpBigFloat(fBigFloat, o)
			return Ref(fBigFloat), Undefined
		case *BigInt:
			oFloat, _ := o.ToGoBigInt().Float64()
			return Float(math.Pow(float64(f), oFloat)).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return Float(math.Pow(float64(f), float64(other.AsFloat()))).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return Float(math.Pow(float64(f), float64(other.AsSmallInt()))).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
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
			prec := max(o.Precision(), 53)
			fBigFloat := (&BigFloat{}).SetPrecision(prec).SetFloat(f)
			return Ref(fBigFloat.Mod(fBigFloat, o)), Undefined
		case *BigInt:
			oFloat, _ := o.ToGoBigInt().Float64()
			return (f.Mod(Float(oFloat))).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return (f.Mod(other.AsFloat())).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return (f.Mod(Float(other.AsSmallInt()))).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Returns 1 if i is greater than other
// Returns 0 if both are equal.
// Returns -1 if i is less than other.
// Returns nil if the comparison was impossible (NaN)
func (f Float) CompareVal(other Value) (result, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			if f.IsNaN() || o.IsNaN() {
				return Nil, Undefined
			}
			iBigFloat := (&BigFloat{}).SetFloat(f)
			return SmallInt(iBigFloat.Cmp(o)).ToValue(), Undefined
		case *BigInt:
			if f.IsNaN() {
				return Nil, Undefined
			}
			return SmallInt(f.Cmp(o.ToFloat())).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		if f.IsNaN() || other.AsFloat().IsNaN() {
			return Nil, Undefined
		}
		return SmallInt(f.Cmp(other.AsFloat())).ToValue(), Undefined
	case SMALL_INT_FLAG:
		if f.IsNaN() {
			return Nil, Undefined
		}
		return SmallInt(f.Cmp(Float(other.AsSmallInt()))).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Check whether f is greater than other and return an error
// if something went wrong.
func (f Float) GreaterThanVal(other Value) (Value, Value) {
	result, err := f.GreaterThan(other)
	return ToElkBool(result), err
}

// Check whether f is greater than other and return an error
// if something went wrong.
func (f Float) GreaterThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			if f.IsNaN() || o.IsNaN() {
				return false, Undefined
			}
			fBigFloat := (&BigFloat{}).SetFloat(f)
			return fBigFloat.Cmp(o) == 1, Undefined
		case *BigInt:
			oFloat := o.ToFloat()
			return f > oFloat, Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		o := other.AsFloat()
		if f.IsNaN() || o.IsNaN() {
			return false, Undefined
		}
		return f > o, Undefined
	case SMALL_INT_FLAG:
		return f > Float(other.AsSmallInt()), Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Check whether f is greater than or equal to other and return an error
// if something went wrong.
func (f Float) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := f.GreaterThanEqual(other)
	return ToElkBool(result), err
}

// Check whether f is greater than or equal to other and return an error
// if something went wrong.
func (f Float) GreaterThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			if f.IsNaN() || o.IsNaN() {
				return false, Undefined
			}
			fBigFloat := (&BigFloat{}).SetFloat(f)
			return fBigFloat.Cmp(o) >= 0, Undefined
		case *BigInt:
			oFloat := o.ToFloat()
			return f >= oFloat, Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		o := other.AsFloat()
		if f.IsNaN() || o.IsNaN() {
			return false, Undefined
		}
		return f >= o, Undefined
	case SMALL_INT_FLAG:
		return f >= Float(other.AsSmallInt()), Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Check whether f is less than other and return an error
// if something went wrong.
func (f Float) LessThanVal(other Value) (Value, Value) {
	result, err := f.LessThan(other)
	return ToElkBool(result), err
}

// Check whether f is less than other and return an error
// if something went wrong.
func (f Float) LessThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			if f.IsNaN() || o.IsNaN() {
				return false, Undefined
			}
			fBigFloat := (&BigFloat{}).SetFloat(f)
			return fBigFloat.Cmp(o) == -1, Undefined
		case *BigInt:
			oFloat := o.ToFloat()
			return f < oFloat, Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		o := other.AsFloat()
		if f.IsNaN() || o.IsNaN() {
			return false, Undefined
		}
		return f < o, Undefined
	case SMALL_INT_FLAG:
		return f < Float(other.AsSmallInt()), Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Check whether f is less than or equal to other and return an error
// if something went wrong.
func (f Float) LessThanEqualVal(other Value) (Value, Value) {
	result, err := f.LessThanEqual(other)
	return ToElkBool(result), err
}

// Check whether f is less than or equal to other and return an error
// if something went wrong.
func (f Float) LessThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			if f.IsNaN() || o.IsNaN() {
				return false, Undefined
			}
			fBigFloat := (&BigFloat{}).SetFloat(f)
			return fBigFloat.Cmp(o) <= 0, Undefined
		case *BigInt:
			oFloat := o.ToFloat()
			return f <= oFloat, Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		o := other.AsFloat()
		if f.IsNaN() || o.IsNaN() {
			return false, Undefined
		}
		return f <= o, Undefined
	case SMALL_INT_FLAG:
		return f <= Float(other.AsSmallInt()), Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Check whether f is equal to other
func (f Float) LaxEqualVal(other Value) Value {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return ToElkBool(f == o.ToFloat())
		case *BigFloat:
			if f.IsNaN() || o.IsNaN() {
				return False
			}
			fBigFloat := (&BigFloat{}).SetFloat(f)
			return ToElkBool(fBigFloat.Cmp(o) == 0)
		case Int64:
			return ToElkBool(f == Float(o))
		case UInt64:
			return ToElkBool(f == Float(o))
		case Float64:
			return ToElkBool(float64(f) == float64(o))
		default:
			return False
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return ToElkBool(f == Float(other.AsSmallInt()))
	case FLOAT_FLAG:
		return ToElkBool(f == other.AsFloat())
	case INT64_FLAG:
		return ToElkBool(f == Float(other.AsInlineInt64()))
	case INT32_FLAG:
		return ToElkBool(f == Float(other.AsInt32()))
	case INT16_FLAG:
		return ToElkBool(f == Float(other.AsInt16()))
	case INT8_FLAG:
		return ToElkBool(f == Float(other.AsInt8()))
	case UINT64_FLAG:
		return ToElkBool(f == Float(other.AsInlineUInt64()))
	case UINT32_FLAG:
		return ToElkBool(f == Float(other.AsUInt32()))
	case UINT16_FLAG:
		return ToElkBool(f == Float(other.AsUInt16()))
	case UINT8_FLAG:
		return ToElkBool(f == Float(other.AsUInt8()))
	case FLOAT64_FLAG:
		return ToElkBool(float64(f) == float64(other.AsInlineFloat64()))
	case FLOAT32_FLAG:
		return ToElkBool(Float(f) == Float(other.AsFloat32()))
	default:
		return False
	}
}

// Check whether f is equal to other and return an error
// if something went wrong.
func (f Float) EqualVal(other Value) Value {
	return ToElkBool(f.Equal(other))
}

// Check whether f is equal to other and return an error
// if something went wrong.
func (f Float) Equal(other Value) bool {
	if other.IsFloat() {
		return f == other.AsFloat()
	}

	return false
}

// Check whether f is strictly equal to other and return an error
// if something went wrong.
func (f Float) StrictEqualVal(other Value) Value {
	return f.EqualVal(other)
}

func (f Float) Nanoseconds() Duration {
	return Duration(f)
}

func (f Float) Microseconds() Duration {
	return Duration(f * Float(Microsecond))
}

func (f Float) Milliseconds() Duration {
	return Duration(f * Float(Millisecond))
}

func (f Float) Seconds() Duration {
	return Duration(f * Float(Second))
}

func (f Float) Minutes() Duration {
	return Duration(f * Float(Minute))
}

func (f Float) Hours() Duration {
	return Duration(f * Float(Hour))
}

func (f Float) Days() Duration {
	return Duration(f * Float(Day))
}

func (f Float) Weeks() Duration {
	return Duration(f * Float(Week))
}

func (f Float) Years() Duration {
	return Duration(f * Float(Year))
}

func initFloat() {
	FloatClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	StdModule.AddConstantString("Float", Ref(FloatClass))
	FloatClass.AddConstantString("NAN", FloatNaN().ToValue())
	FloatClass.AddConstantString("INF", FloatInf().ToValue())
	FloatClass.AddConstantString("NEG_INF", FloatNegInf().ToValue())

	FloatClass.AddConstantString("Convertible", Ref(NewInterface()))
}
