package value

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"unsafe"

	"github.com/cespare/xxhash/v2"
)

// Elk's SmallInt value
type SmallInt int

const PtrSize = 4 << (^uintptr(0) >> 63)

// Number of bits available for a small int.
const SmallIntBits = PtrSize * 8

// Max value of SmallInt
const MaxSmallInt = math.MaxInt

// Min value of SmallInt
const MinSmallInt = math.MinInt

func (i SmallInt) ToValue() Value {
	return Value{
		flag: SMALL_INT_FLAG,
		data: *(*uintptr)(unsafe.Pointer(&i)),
	}
}

func (SmallInt) Class() *Class {
	return IntClass
}

func (SmallInt) DirectClass() *Class {
	return IntClass
}

func (SmallInt) SingletonClass() *Class {
	return nil
}

func (i SmallInt) Inspect() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i SmallInt) Error() string {
	return i.Inspect()
}

func (i SmallInt) InstanceVariables() *InstanceVariables {
	return nil
}

// NegateVal the number and return the result.
func (i SmallInt) NegateVal() Value {
	if i == MinSmallInt {
		iBigInt := big.NewInt(int64(i))
		return Ref(ToElkBigInt(iBigInt.Neg(iBigInt)))
	}

	return (-i).ToValue()
}

// Increment the number and return the result.
func (i SmallInt) Increment() Value {
	result, ok := i.AddOverflow(1)
	if !ok {
		iBigInt := big.NewInt(int64(i))
		return Ref(ToElkBigInt(iBigInt.Add(iBigInt, big.NewInt(1))))
	}
	return result.ToValue()
}

// Decrement the number and return the result.
func (i SmallInt) Decrement() Value {
	result, ok := i.SubtractOverflow(1)
	if !ok {
		iBigInt := big.NewInt(int64(i))
		return Ref(ToElkBigInt(iBigInt.Sub(iBigInt, big.NewInt(1))))
	}
	return result.ToValue()
}

// Add two small ints and check for overflow/underflow.
func (a SmallInt) AddOverflow(b SmallInt) (result SmallInt, ok bool) {
	c := a + b
	if (c > a) == (b > 0) {
		return c, true
	}
	return c, false
}

// Convert to Elk String.
func (i SmallInt) ToString() String {
	return String(i.Inspect())
}

// Convert to Elk Float.
func (i SmallInt) ToFloat() Float {
	return Float(i)
}

// Convert to Elk Float64.
func (i SmallInt) ToFloat64() Float64 {
	return Float64(i)
}

// Convert to Elk Float32.
func (i SmallInt) ToFloat32() Float32 {
	return Float32(i)
}

// Convert to Elk Int64.
func (i SmallInt) ToInt64() Int64 {
	return Int64(i)
}

// Convert to Elk Int32.
func (i SmallInt) ToInt32() Int32 {
	return Int32(i)
}

// Convert to Elk Int16.
func (i SmallInt) ToInt16() Int16 {
	return Int16(i)
}

// Convert to Elk Int8.
func (i SmallInt) ToInt8() Int8 {
	return Int8(i)
}

// Convert to Elk UInt64.
func (i SmallInt) ToUInt64() UInt64 {
	return UInt64(i)
}

// Convert to Elk UInt32.
func (i SmallInt) ToUInt32() UInt32 {
	return UInt32(i)
}

// Convert to Elk UInt16.
func (i SmallInt) ToUInt16() UInt16 {
	return UInt16(i)
}

// Convert to Elk UInt8.
func (i SmallInt) ToUInt8() UInt8 {
	return UInt8(i)
}

// AddVal another value and return an error
// if something went wrong.
func (i SmallInt) AddVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := big.NewInt(int64(i))
			iBigInt.Add(iBigInt, o.ToGoBigInt())
			if iBigInt.IsInt64() {
				return SmallInt(iBigInt.Int64()).ToValue(), Undefined
			}
			return Ref(ToElkBigInt(iBigInt)), Undefined
		case *BigFloat:
			prec := max(o.Precision(), 64)
			iBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(i)
			iBigFloat.AddBigFloat(iBigFloat, o)
			return Ref(iBigFloat), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		result, ok := i.AddOverflow(o)
		if !ok {
			iBigInt := big.NewInt(int64(i))
			return Ref(ToElkBigInt(iBigInt.Add(iBigInt, big.NewInt(int64(o))))), Undefined
		}
		return result.ToValue(), Undefined
	case FLOAT_FLAG:
		return (Float(i) + other.AsFloat()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Subtract two small ints and check for overflow/underflow.
func (a SmallInt) SubtractOverflow(b SmallInt) (result SmallInt, ok bool) {
	c := a - b
	if (c < a) == (b > 0) {
		return c, true
	}
	return c, false
}

// Add another value and return an error
// if something went wrong.
func (i SmallInt) SubtractVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := big.NewInt(int64(i))
			iBigInt.Sub(iBigInt, o.ToGoBigInt())
			if iBigInt.IsInt64() {
				return SmallInt(iBigInt.Int64()).ToValue(), Undefined
			}
			return Ref(ToElkBigInt(iBigInt)), Undefined
		case *BigFloat:
			prec := max(o.Precision(), 64)
			iBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(i)
			iBigFloat.SubBigFloat(iBigFloat, o)
			return Ref(iBigFloat), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		result, ok := i.SubtractOverflow(o)
		if !ok {
			iBigInt := big.NewInt(int64(i))
			return Ref(ToElkBigInt(iBigInt.Sub(iBigInt, big.NewInt(int64(o))))), Undefined
		}
		return result.ToValue(), Undefined
	case FLOAT_FLAG:
		return (Float(i) - other.AsFloat()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Multiply two small ints and check for overflow/underflow.
func (a SmallInt) MultiplyOverflow(b SmallInt) (result SmallInt, ok bool) {
	if a == 0 || b == 0 {
		return 0, true
	}
	c := a * b
	if (c < 0) == ((a < 0) != (b < 0)) {
		if c/b == a {
			return c, true
		}
	}
	return c, false
}

// MultiplyVal another value and return an error
// if something went wrong.
func (i SmallInt) MultiplyVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := big.NewInt(int64(i))
			iBigInt.Mul(iBigInt, o.ToGoBigInt())
			if iBigInt.IsInt64() {
				return SmallInt(iBigInt.Int64()).ToValue(), Undefined
			}
			return Ref(ToElkBigInt(iBigInt)), Undefined
		case *BigFloat:
			prec := max(o.Precision(), 64)
			iBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(i)
			iBigFloat.MulBigFloat(iBigFloat, o)
			return Ref(iBigFloat), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		result, ok := i.MultiplyOverflow(o)
		if !ok {
			iBigInt := big.NewInt(int64(i))
			return Ref(ToElkBigInt(iBigInt.Mul(iBigInt, big.NewInt(int64(o))))), Undefined
		}
		return result.ToValue(), Undefined
	case FLOAT_FLAG:
		return (Float(i) * other.AsFloat()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Divide two small ints and check for overflow/underflow.
func (a SmallInt) DivideOverflow(b SmallInt) (result SmallInt, ok bool) {
	if b == 0 {
		return 0, false
	}
	c := a / b
	return c, (c < 0) == ((a < 0) != (b < 0))
}

// DivideVal another value and return an error
// if something went wrong.
func (i SmallInt) DivideVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			if o.IsZero() {
				return Undefined, Ref(NewZeroDivisionError())
			}
			iBigInt := big.NewInt(int64(i))
			iBigInt.Div(iBigInt, o.ToGoBigInt())
			if iBigInt.IsInt64() {
				return SmallInt(iBigInt.Int64()).ToValue(), Undefined
			}
			return Ref(ToElkBigInt(iBigInt)), Undefined
		case *BigFloat:
			prec := max(o.Precision(), 64)
			iBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(i)
			iBigFloat.DivBigFloat(iBigFloat, o)
			return Ref(iBigFloat), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		if o == 0 {
			return Undefined, Ref(NewZeroDivisionError())
		}
		result, ok := i.DivideOverflow(o)
		if !ok {
			iBigInt := big.NewInt(int64(i))
			return Ref(ToElkBigInt(iBigInt.Div(iBigInt, big.NewInt(int64(o))))), Undefined
		}
		return result.ToValue(), Undefined
	case FLOAT_FLAG:
		return (Float(i) / other.AsFloat()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// ExponentiateVal by another value and return an error
// if something went wrong.
func (i SmallInt) ExponentiateVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := big.NewInt(int64(i))
			iBigInt.Exp(iBigInt, o.ToGoBigInt(), nil)
			if iBigInt.IsInt64() {
				return SmallInt(iBigInt.Int64()).ToValue(), Undefined
			}
			return Ref(ToElkBigInt(iBigInt)), Undefined
		case *BigFloat:
			prec := max(o.Precision(), 64)
			iBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(i)
			iBigFloat.ExpBigFloat(iBigFloat, o)
			return Ref(iBigFloat), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		iBigInt := big.NewInt(int64(i))
		oBigInt := big.NewInt(int64(o))
		iBigInt.Exp(iBigInt, oBigInt, nil)
		if iBigInt.IsInt64() {
			return SmallInt(iBigInt.Int64()).ToValue(), Undefined
		}
		return Ref(ToElkBigInt(iBigInt)), Undefined
	case FLOAT_FLAG:
		return (Float(math.Pow(float64(i), float64(other.AsFloat())))).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Cmp compares x and y and returns:
//
//	  -1 if x <  y
//		 0 if x == y
//	  +1 if x >  y
func (x SmallInt) Cmp(y SmallInt) int {
	if x > y {
		return 1
	}
	if x < y {
		return -1
	}
	return 0
}

// Returns 1 if i is greater than other
// Returns 0 if both are equal.
// Returns -1 if i is less than other.
// Returns nil if the comparison was impossible (NaN)
func (i SmallInt) CompareVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := NewBigInt(int64(i))
			return SmallInt(iBigInt.Cmp(o)).ToValue(), Undefined
		case *BigFloat:
			if o.IsNaN() {
				return Nil, Undefined
			}
			iBigFloat := (&BigFloat{}).SetSmallInt(i)
			return SmallInt(iBigFloat.Cmp(o)).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return SmallInt(i.Cmp(other.AsSmallInt())).ToValue(), Undefined
	case FLOAT_FLAG:
		o := other.AsFloat()
		if o.IsNaN() {
			return Nil, Undefined
		}
		return SmallInt(i.ToFloat().Cmp(o)).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Check whether i is greater than other and return an error
// if something went wrong.
func (i SmallInt) GreaterThanVal(other Value) (Value, Value) {
	result, err := i.GreaterThan(other)
	return ToElkBool(result), err
}

// Check whether i is greater than other and return an error
// if something went wrong.
func (i SmallInt) GreaterThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := NewBigInt(int64(i))
			return iBigInt.Cmp(o) == 1, Undefined
		case *BigFloat:
			if o.IsNaN() {
				return false, Undefined
			}
			iBigFloat := (&BigFloat{}).SetSmallInt(i)
			return iBigFloat.Cmp(o) == 1, Undefined
		default:
			return false, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return i > other.AsSmallInt(), Undefined
	case FLOAT_FLAG:
		return Float(i) > other.AsFloat(), Undefined
	default:
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Check whether i is greater than or equal to other and return an error
// if something went wrong.
func (i SmallInt) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := i.GreaterThanEqual(other)
	return ToElkBool(result), err
}

// Check whether i is greater than or equal to other and return an error
// if something went wrong.
func (i SmallInt) GreaterThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := NewBigInt(int64(i))
			return iBigInt.Cmp(o) >= 0, Undefined
		case *BigFloat:
			if o.IsNaN() {
				return false, Undefined
			}
			iBigFloat := (&BigFloat{}).SetSmallInt(i)
			return iBigFloat.Cmp(o) >= 0, Undefined
		default:
			return false, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return i >= other.AsSmallInt(), Undefined
	case FLOAT_FLAG:
		return Float(i) >= other.AsFloat(), Undefined
	default:
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Check whether i is less than other and return an error
// if something went wrong.
func (i SmallInt) LessThanVal(other Value) (Value, Value) {
	result, err := i.LessThan(other)
	return ToElkBool(result), err
}

// Check whether i is less than other and return an error
// if something went wrong.
func (i SmallInt) LessThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := NewBigInt(int64(i))
			return iBigInt.Cmp(o) == -1, Undefined
		case *BigFloat:
			if o.IsNaN() {
				return false, Undefined
			}
			iBigFloat := (&BigFloat{}).SetSmallInt(i)
			return iBigFloat.Cmp(o) == -1, Undefined
		default:
			return false, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return i < other.AsSmallInt(), Undefined
	case FLOAT_FLAG:
		return Float(i) < other.AsFloat(), Undefined
	default:
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Check whether i is less than or equal to other and return an error
// if something went wrong.
func (i SmallInt) LessThanEqualVal(other Value) (Value, Value) {
	result, err := i.LessThanEqual(other)
	return ToElkBool(result), err
}

// Check whether i is less than or equal to other and return an error
// if something went wrong.
func (i SmallInt) LessThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := NewBigInt(int64(i))
			return iBigInt.Cmp(o) <= 0, Undefined
		case *BigFloat:
			if o.IsNaN() {
				return false, Undefined
			}
			iBigFloat := (&BigFloat{}).SetSmallInt(i)
			return iBigFloat.Cmp(o) <= 0, Undefined
		default:
			return false, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return i <= other.AsSmallInt(), Undefined
	case FLOAT_FLAG:
		return Float(i) <= other.AsFloat(), Undefined
	default:
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Check whether i is equal to other (with coercion)
func (i SmallInt) LaxEqualVal(other Value) Value {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := NewBigInt(int64(i))
			return ToElkBool(iBigInt.Cmp(o) == 0)
		case *BigFloat:
			if o.IsNaN() {
				return False
			}
			iBigFloat := (&BigFloat{}).SetSmallInt(i)
			return ToElkBool(iBigFloat.Cmp(o) == 0)
		case Int64:
			return ToElkBool(i == SmallInt(o))
		case UInt64:
			if o > MaxSmallInt {
				return False
			}
			return ToElkBool(i == SmallInt(o))
		default:
			return False
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return ToElkBool(i == other.AsSmallInt())
	case FLOAT_FLAG:
		return ToElkBool(Float(i) == other.AsFloat())
	case INT64_FLAG:
		return ToElkBool(i == SmallInt(other.AsInlineInt64()))
	case INT32_FLAG:
		return ToElkBool(i == SmallInt(other.AsInt32()))
	case INT16_FLAG:
		return ToElkBool(i == SmallInt(other.AsInt16()))
	case INT8_FLAG:
		return ToElkBool(i == SmallInt(other.AsInt8()))
	case UINT64_FLAG:
		o := other.AsInlineUInt64()
		if o > MaxSmallInt {
			return False
		}
		return ToElkBool(i == SmallInt(o))
	case UINT32_FLAG:
		return ToElkBool(i == SmallInt(other.AsUInt32()))
	case UINT16_FLAG:
		return ToElkBool(i == SmallInt(other.AsUInt16()))
	case UINT8_FLAG:
		return ToElkBool(i == SmallInt(other.AsUInt8()))
	case FLOAT64_FLAG:
		return ToElkBool(Float64(i) == other.AsInlineFloat64())
	case FLOAT32_FLAG:
		return ToElkBool(Float32(i) == other.AsFloat32())
	default:
		return False
	}
}

// Check whether i is equal to other
func (i SmallInt) EqualVal(other Value) Value {
	return i.StrictEqualVal(other)
}

// Check whether i is equal to other
func (i SmallInt) Equal(other Value) bool {
	return i.StrictEqualBool(other)
}

// Check whether i is strictly equal to other
func (i SmallInt) StrictEqualBool(other Value) bool {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := NewBigInt(int64(i))
			return iBigInt.Cmp(o) == 0
		default:
			return false
		}
	}
	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return i == other.AsSmallInt()
	default:
		return false
	}
}

// Check whether i is strictly equal to other
func (i SmallInt) StrictEqualVal(other Value) Value {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := NewBigInt(int64(i))
			return ToElkBool(iBigInt.Cmp(o) == 0)
		default:
			return False
		}
	}
	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return ToElkBool(i == other.AsSmallInt())
	default:
		return False
	}
}

func leftBitshiftSmallInt[T SimpleInt](i SmallInt, other T) Value {
	var bitsize T = SmallIntBits - 1
	if other < 0 {
		return SmallInt(0).ToValue()
	}
	complementaryShift := i >> (bitsize - other)
	if other > bitsize || (i < 0 && complementaryShift != -1) || (i > 0 && complementaryShift != 0) {
		// overflow
		iBig := big.NewInt(int64(i))
		iBig.Lsh(iBig, uint(other))
		return Ref(ToElkBigInt(iBig))
	}
	return (i << other).ToValue()
}

func rightBitshiftSmallInt[T SimpleInt](i SmallInt, other T) Value {
	if other < 0 {
		return SmallInt(0).ToValue()
	}
	return (i >> other).ToValue()
}

// Bitshift to the left by another integer value and return an error
// if something went wrong.
func (i SmallInt) LeftBitshiftVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case Int64:
			if o < 0 {
				return rightBitshiftSmallInt(i, -o), Undefined
			}
			return leftBitshiftSmallInt(i, o), Undefined
		case UInt64:
			return leftBitshiftSmallInt(i, o), Undefined
		case *BigInt:
			if o.IsSmallInt() {
				oSmall := o.ToSmallInt()
				if oSmall < 0 {
					return rightBitshiftSmallInt(i, -oSmall), Undefined
				}
				return leftBitshiftSmallInt(i, oSmall), Undefined
			}
			return SmallInt(0).ToValue(), Undefined
		default:
			return Undefined, Ref(NewBitshiftOperandError(other))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		if o < 0 {
			return rightBitshiftSmallInt(i, -o), Undefined
		}
		return leftBitshiftSmallInt(i, o), Undefined
	case INT64_FLAG:
		o := other.AsInlineInt64()
		if o < 0 {
			return rightBitshiftSmallInt(i, -o), Undefined
		}
		return leftBitshiftSmallInt(i, o), Undefined
	case INT32_FLAG:
		o := other.AsInt32()
		if o < 0 {
			return rightBitshiftSmallInt(i, -o), Undefined
		}
		return leftBitshiftSmallInt(i, o), Undefined
	case INT16_FLAG:
		o := other.AsInt16()
		if o < 0 {
			return rightBitshiftSmallInt(i, -o), Undefined
		}
		return leftBitshiftSmallInt(i, o), Undefined
	case INT8_FLAG:
		o := other.AsInt8()
		if o < 0 {
			return rightBitshiftSmallInt(i, -o), Undefined
		}
		return leftBitshiftSmallInt(i, o), Undefined
	case UINT64_FLAG:
		o := other.AsInlineUInt64()
		return leftBitshiftSmallInt(i, o), Undefined
	case UINT32_FLAG:
		o := other.AsUInt32()
		return leftBitshiftSmallInt(i, o), Undefined
	case UINT16_FLAG:
		o := other.AsUInt16()
		return leftBitshiftSmallInt(i, o), Undefined
	case UINT8_FLAG:
		o := other.AsUInt8()
		return leftBitshiftSmallInt(i, o), Undefined
	default:
		return Undefined, Ref(NewBitshiftOperandError(other))
	}
}

// Bitshift to the right by another integer value and return an error
// if something went wrong.
func (i SmallInt) RightBitshiftVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case Int64:
			if o < 0 {
				return leftBitshiftSmallInt(i, -o), Undefined
			}
			return (i >> o).ToValue(), Undefined
		case UInt64:
			return (i >> o).ToValue(), Undefined
		case *BigInt:
			if o.IsSmallInt() {
				oSmall := o.ToSmallInt()
				if oSmall < 0 {
					return leftBitshiftSmallInt(i, -oSmall), Undefined
				}
				return (i >> oSmall).ToValue(), Undefined
			}
			return SmallInt(0).ToValue(), Undefined
		default:
			return Undefined, Ref(NewBitshiftOperandError(other))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), Undefined
		}
		return (i >> o).ToValue(), Undefined
	case INT64_FLAG:
		o := other.AsInlineInt64()
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), Undefined
		}
		return (i >> o).ToValue(), Undefined
	case INT32_FLAG:
		o := other.AsInt32()
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), Undefined
		}
		return (i >> o).ToValue(), Undefined
	case INT16_FLAG:
		o := other.AsInt16()
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), Undefined
		}
		return (i >> o).ToValue(), Undefined
	case INT8_FLAG:
		o := other.AsInt8()
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), Undefined
		}
		return (i >> o).ToValue(), Undefined
	case UINT64_FLAG:
		o := other.AsInlineUInt64()
		return (i >> o).ToValue(), Undefined
	case UINT32_FLAG:
		o := other.AsUInt32()
		return (i >> o).ToValue(), Undefined
	case UINT16_FLAG:
		o := other.AsUInt16()
		return (i >> o).ToValue(), Undefined
	case UINT8_FLAG:
		o := other.AsUInt8()
		return (i >> o).ToValue(), Undefined
	default:
		return Undefined, Ref(NewBitshiftOperandError(other))
	}
}

// Perform a bitwise AND with another integer value and return an error
// if something went wrong.
func (i SmallInt) BitwiseAndVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := big.NewInt(int64(i))
			iBigInt.And(iBigInt, o.ToGoBigInt())
			result := ToElkBigInt(iBigInt)
			if result.IsSmallInt() {
				return result.ToSmallInt().ToValue(), Undefined
			}
			return Ref(result), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		return (i & o).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Perform bitwise NOT.
func (i SmallInt) BitwiseNot() SmallInt {
	return ^i
}

// Perform a bitwise AND NOT with another integer value and return an error
// if something went wrong.
func (i SmallInt) BitwiseAndNotVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := big.NewInt(int64(i))
			iBigInt.AndNot(iBigInt, o.ToGoBigInt())
			result := ToElkBigInt(iBigInt)
			if result.IsSmallInt() {
				return result.ToSmallInt().ToValue(), Undefined
			}
			return Ref(result), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		return (i &^ o).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Perform a bitwise OR with another integer value and return an error
// if something went wrong.
func (i SmallInt) BitwiseOrVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := big.NewInt(int64(i))
			iBigInt.Or(iBigInt, o.ToGoBigInt())
			result := ToElkBigInt(iBigInt)
			if result.IsSmallInt() {
				return result.ToSmallInt().ToValue(), Undefined
			}
			return Ref(result), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		return (i | o).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Perform a bitwise XOR with another integer value and return an error
// if something went wrong.
func (i SmallInt) BitwiseXorVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			iBigInt := big.NewInt(int64(i))
			iBigInt.Xor(iBigInt, o.ToGoBigInt())
			result := ToElkBigInt(iBigInt)
			if result.IsSmallInt() {
				return result.ToSmallInt().ToValue(), Undefined
			}
			return Ref(result), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		return (i ^ o).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Perform modulo by another numeric value and return an error
// if something went wrong.
func (i SmallInt) ModuloVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			if o.IsSmallInt() {
				oSmall := o.ToSmallInt()
				return (i % oSmall).ToValue(), Undefined
			}
			return i.ToValue(), Undefined
		case *BigFloat:
			prec := max(o.Precision(), 64)
			iBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(i)
			return Ref(iBigFloat.Mod(iBigFloat, o)), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return (i % other.AsSmallInt()).ToValue(), Undefined
	case FLOAT_FLAG:
		return Float(math.Mod(float64(i), float64(other.AsFloat()))).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

func (i SmallInt) Hash() UInt64 {
	d := xxhash.New()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	d.Write(b)
	return UInt64(d.Sum64())
}

func (i SmallInt) Nanoseconds() TimeSpan {
	return TimeSpan(i)
}

func (i SmallInt) Microseconds() TimeSpan {
	return TimeSpan(i) * Microsecond
}

func (i SmallInt) Milliseconds() TimeSpan {
	return TimeSpan(i) * Millisecond
}

func (i SmallInt) Seconds() TimeSpan {
	return TimeSpan(i) * Second
}

func (i SmallInt) Minutes() TimeSpan {
	return TimeSpan(i) * Minute
}

func (i SmallInt) Hours() TimeSpan {
	return TimeSpan(i) * Hour
}

func (i SmallInt) Days() DateSpan {
	return MakeDateSpan(0, 0, int(i))
}

func (i SmallInt) Weeks() DateSpan {
	return MakeDateSpan(0, 0, int(i)*7)
}

func (i SmallInt) Months() DateSpan {
	return MakeDateSpan(0, int(i), 0)
}

func (i SmallInt) Years() DateSpan {
	return MakeDateSpan(int(i), 0, 0)
}

func (i SmallInt) Centuries() DateSpan {
	return MakeDateSpan(int(i)*100, 0, 0)
}

func (i SmallInt) Millenia() DateSpan {
	return MakeDateSpan(int(i)*1000, 0, 0)
}

type SmallIntIterator struct {
	Int     SmallInt
	Counter SmallInt
}

func NewSmallIntIterator(i SmallInt) *SmallIntIterator {
	return &SmallIntIterator{
		Int: i,
	}
}

func NewSmallIntIteratorWithCounter(i SmallInt, index SmallInt) *SmallIntIterator {
	return &SmallIntIterator{
		Int:     i,
		Counter: index,
	}
}

func (*SmallIntIterator) Class() *Class {
	return IntIteratorClass
}

func (*SmallIntIterator) DirectClass() *Class {
	return IntIteratorClass
}

func (*SmallIntIterator) SingletonClass() *Class {
	return nil
}

func (l *SmallIntIterator) Copy() Reference {
	return &SmallIntIterator{
		Int:     l.Int,
		Counter: l.Counter,
	}
}

func (l *SmallIntIterator) Inspect() string {
	return fmt.Sprintf("Std::Int::Iterator{&: %p, int: %s, index: %d}", l, l.Int.Inspect(), l.Counter)
}

func (l *SmallIntIterator) Error() string {
	return l.Inspect()
}

func (*SmallIntIterator) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *SmallIntIterator) Next() (Value, Value) {
	if l.Counter >= l.Int {
		return Undefined, stopIterationSymbol.ToValue()
	}

	next := l.Counter
	l.Counter++
	return next.ToValue(), Undefined
}

func (l *SmallIntIterator) Reset() {
	l.Counter = 0
}
