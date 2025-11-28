package value

import (
	"fmt"
	"math"
	"math/big"

	"github.com/cespare/xxhash/v2"
)

// Elk's BigInt value
type BigInt big.Int

// Create a new BigInt with the specified value.
func NewBigInt(i int64) *BigInt {
	return ToElkBigInt(big.NewInt(i))
}

// Convert Go int64 to Elk Int.
func ToElkInt(i int64) Value {
	if i > MaxSmallInt {
		bi := NewBigInt(i)
		return Ref(bi)
	}
	return SmallInt(i).ToValue()
}

// Convert Go big.Int value to Elk BigInt value.
func ToElkBigInt(i *big.Int) *BigInt {
	return (*BigInt)(i)
}

// Convert the Elk BigInt value to Go big.Int value.
func (i *BigInt) ToGoBigInt() *big.Int {
	return (*big.Int)(i)
}

// Reports whether i can be represented as a SmallInt.
func (i *BigInt) IsSmallInt() bool {
	return i.ToGoBigInt().IsInt64()
}

// Reports whether i is zero.
func (i *BigInt) IsZero() bool {
	return len(i.ToGoBigInt().Bits()) == 0
}

func (i *BigInt) IsEven() bool {
	bigInt := i.ToGoBigInt()
	result := bigInt.Mod(bigInt, big.NewInt(2))
	return len(result.Bits()) == 0
}

func (i *BigInt) IsOdd() bool {
	bigInt := i.ToGoBigInt()
	result := bigInt.Mod(bigInt, big.NewInt(2))
	return len(result.Bits()) != 0
}

// Returns the SmallInt representation of i.
func (i *BigInt) ToSmallInt() SmallInt {
	return SmallInt(i.ToGoBigInt().Int64())
}

// Convert to Elk Int64
func (i *BigInt) ToInt64() Int64 {
	return Int64(i.ToGoBigInt().Int64())
}

// Convert to Elk Int32
func (i *BigInt) ToInt32() Int32 {
	return i.ToSmallInt().ToInt32()
}

// Convert to Elk Int16
func (i *BigInt) ToInt16() Int16 {
	return i.ToSmallInt().ToInt16()
}

// Convert to Elk Int8
func (i *BigInt) ToInt8() Int8 {
	return i.ToSmallInt().ToInt8()
}

// Convert to Elk UInt64
func (i *BigInt) ToUInt64() UInt64 {
	return i.ToSmallInt().ToUInt64()
}

// Convert to Elk UInt32
func (i *BigInt) ToUInt32() UInt32 {
	return i.ToSmallInt().ToUInt32()
}

// Convert to Elk UInt16
func (i *BigInt) ToUInt16() UInt16 {
	return i.ToSmallInt().ToUInt16()
}

// Convert to Elk UInt8
func (i *BigInt) ToUInt8() UInt8 {
	return i.ToSmallInt().ToUInt8()
}

// Convert the Elk BigInt value to Elk String.
func (i *BigInt) ToString() String {
	return String(i.Inspect())
}

// Returns the Float representation of i.
func (i *BigInt) ToFloat() Float {
	f, _ := i.ToGoBigInt().Float64()
	return Float(f)
}

// Convert to Elk Float64
func (i *BigInt) ToFloat64() Float64 {
	return Float64(i.ToFloat())
}

// Convert to Elk Float32
func (i *BigInt) ToFloat32() Float32 {
	return Float32(i.ToFloat())
}

// Negate the number and return the result.
func (i *BigInt) Negate() *BigInt {
	return ToElkBigInt((&big.Int{}).Neg(i.ToGoBigInt()))
}

// Increment the number and return the result.
func (i *BigInt) Increment() *BigInt {
	oBigInt := big.NewInt(int64(1))
	oBigInt.Add(i.ToGoBigInt(), oBigInt)
	result := ToElkBigInt(oBigInt)
	return result
}

// DecrementVal the number and return the result.
func (i *BigInt) DecrementVal() Value {
	oBigInt := big.NewInt(int64(1))
	oBigInt.Sub(i.ToGoBigInt(), oBigInt)
	result := ToElkBigInt(oBigInt)
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue()
	}
	return Ref(result)
}

// Perform bitwise not on the number and return the result.
func (i *BigInt) BitwiseNot() *BigInt {
	return ToElkBigInt((&big.Int{}).Not(i.ToGoBigInt()))
}

// Number of bits required to represent this integer.
func (i *BigInt) BitSize() int {
	return i.ToGoBigInt().BitLen()
}

// Cmp compares x and y and returns:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
func (x *BigInt) Cmp(y *BigInt) int {
	return x.ToGoBigInt().Cmp(y.ToGoBigInt())
}

// Returns 1 if i is greater than other
// Returns 0 if both are equal.
// Returns -1 if i is less than other.
// Returns nil if the comparison was impossible (NaN)
func (i *BigInt) CompareVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return SmallInt(i.Cmp(o)).ToValue(), Undefined
		case *BigFloat:
			if o.IsNaN() {
				return Nil, Undefined
			}
			iBigFloat := (&BigFloat{}).SetBigInt(i)
			return SmallInt(iBigFloat.Cmp(o)).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		oBigInt := NewBigInt(int64(other.AsSmallInt()))
		return SmallInt(i.Cmp(oBigInt)).ToValue(), Undefined
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

// AddVal another value and return an error
// if something went wrong.
func (i *BigInt) AddVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return i.AddBigInt(o), Undefined
		case *BigFloat:
			return Ref(i.AddBigFloat(o)), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}
	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return i.AddSmallInt(other.AsSmallInt()), Undefined
	case FLOAT_FLAG:
		return i.AddFloat(other.AsFloat()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

func (i *BigInt) AddFloat(other Float) Float {
	return i.ToFloat() + other
}

func (i *BigInt) AddBigFloat(other *BigFloat) *BigFloat {
	prec := max(other.Precision(), uint(i.BitSize()), 64)
	iBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(i)
	iBigFloat.AddBigFloat(iBigFloat, other)
	return iBigFloat
}

func (i *BigInt) AddBigInt(other *BigInt) Value {
	result := ToElkBigInt((&big.Int{}).Add(i.ToGoBigInt(), other.ToGoBigInt()))
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue()
	}
	return Ref(result)
}

func (i *BigInt) AddSmallInt(other SmallInt) Value {
	oBigInt := big.NewInt(int64(other))
	oBigInt.Add(i.ToGoBigInt(), oBigInt)
	result := ToElkBigInt(oBigInt)
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue()
	}
	return Ref(ToElkBigInt(oBigInt))
}

// SubtractVal another value and return an error
// if something went wrong.
func (i *BigInt) SubtractVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return i.SubtractBigInt(o), Undefined
		case *BigFloat:
			return Ref(i.SubtractBigFloat(o)), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return i.SubtractSmallInt(other.AsSmallInt()), Undefined
	case FLOAT_FLAG:
		return i.SubtractFloat(other.AsFloat()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

func (i *BigInt) SubtractFloat(other Float) Float {
	return i.ToFloat() - other
}

func (i *BigInt) SubtractBigFloat(other *BigFloat) *BigFloat {
	prec := max(other.Precision(), uint(i.BitSize()), 64)
	iBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(i)
	iBigFloat.SubBigFloat(iBigFloat, other)
	return iBigFloat
}

func (i *BigInt) SubtractSmallInt(other SmallInt) Value {
	oBigInt := big.NewInt(int64(other))
	oBigInt.Sub(i.ToGoBigInt(), oBigInt)
	result := ToElkBigInt(oBigInt)
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue()
	}
	return Ref(ToElkBigInt(oBigInt))
}

func (i *BigInt) SubtractBigInt(other *BigInt) Value {
	result := ToElkBigInt((&big.Int{}).Sub(i.ToGoBigInt(), other.ToGoBigInt()))
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue()
	}
	return Ref(result)
}

// MultiplyVal by another value and return an error
// if something went wrong.
func (i *BigInt) MultiplyVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return i.MultiplyBigInt(o), Undefined
		case *BigFloat:
			return Ref(i.MultiplyBigFloat(o)), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return i.MultiplySmallInt(other.AsSmallInt()), Undefined
	case FLOAT_FLAG:
		return i.MultiplyFloat(other.AsFloat()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

func (i *BigInt) MultiplyBigFloat(other *BigFloat) *BigFloat {
	prec := max(other.Precision(), uint(i.BitSize()), 64)
	iBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(i)
	return iBigFloat.MulBigFloat(iBigFloat, other)
}

func (i *BigInt) MultiplyBigInt(other *BigInt) Value {
	result := ToElkBigInt((&big.Int{}).Mul(i.ToGoBigInt(), other.ToGoBigInt()))
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue()
	}
	return Ref(result)
}

func (i *BigInt) MultiplySmallInt(other SmallInt) Value {
	oBigInt := big.NewInt(int64(other))
	oBigInt.Mul(i.ToGoBigInt(), oBigInt)
	result := ToElkBigInt(oBigInt)
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue()
	}
	return Ref(result)
}

func (i *BigInt) MultiplyFloat(other Float) Float {
	return i.ToFloat() * other
}

// DivideVal by another value and return an error
// if something went wrong.
func (i *BigInt) DivideVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return i.DivideBigInt(o)
		case *BigFloat:
			return Ref(i.DivideBigFloat(o)), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return i.DivideSmallInt(other.AsSmallInt())
	case FLOAT_FLAG:
		return i.DivideFloat(other.AsFloat()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

func (i *BigInt) DivideBigInt(other *BigInt) (Value, Value) {
	if other.IsZero() {
		return Undefined, Ref(NewZeroDivisionError())
	}
	result := ToElkBigInt((&big.Int{}).Div(i.ToGoBigInt(), other.ToGoBigInt()))
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue(), Undefined
	}
	return Ref(result), Undefined
}

func (i *BigInt) DivideBigFloat(other *BigFloat) *BigFloat {
	prec := max(other.Precision(), uint(i.BitSize()), 64)
	iBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(i)
	return iBigFloat.DivBigFloat(iBigFloat, other)
}

func (i *BigInt) DivideSmallInt(other SmallInt) (Value, Value) {
	if other == 0 {
		return Undefined, Ref(NewZeroDivisionError())
	}
	oBigInt := big.NewInt(int64(other))
	oBigInt.Div(i.ToGoBigInt(), oBigInt)
	result := ToElkBigInt(oBigInt)
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue(), Undefined
	}
	return Ref(ToElkBigInt(oBigInt)), Undefined
}

func (i *BigInt) DivideFloat(other Float) Float {
	return i.ToFloat() / other
}

// ExponentiateVal by another value and return an error
// if something went wrong.
func (i *BigInt) ExponentiateVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return i.ExponentiateBigInt(o), Undefined
		case *BigFloat:
			return Ref(i.ExponentiateBigFloat(o)), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return i.ExponentiateSmallInt(other.AsSmallInt()), Undefined
	case FLOAT_FLAG:
		return i.ExponentiateFloat(other.AsFloat()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

func (i *BigInt) ExponentiateBigFloat(other *BigFloat) *BigFloat {
	prec := max(other.Precision(), uint(i.BitSize()), 64)
	iBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(i)
	iBigFloat.ExpBigFloat(iBigFloat, other)
	return iBigFloat
}

func (i *BigInt) ExponentiateFloat(other Float) Float {
	iFloat, _ := i.ToGoBigInt().Float64()
	return Float(math.Pow(iFloat, float64(other)))
}

func (i *BigInt) ExponentiateSmallInt(other SmallInt) Value {
	oBigInt := big.NewInt(int64(other))
	oBigInt.Exp(i.ToGoBigInt(), oBigInt, nil)
	result := ToElkBigInt(oBigInt)
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue()
	}
	return Ref(ToElkBigInt(oBigInt))
}

func (i *BigInt) ExponentiateBigInt(other *BigInt) Value {
	result := ToElkBigInt((&big.Int{}).Exp(i.ToGoBigInt(), other.ToGoBigInt(), nil))
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue()
	}
	return Ref(result)
}

// Perform modulo with another numeric value and return an error
// if something went wrong.
func (i *BigInt) ModuloVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return i.ModuloBigInt(o)
		case *BigFloat:
			return Ref(i.ModuloBigFloat(o)), Undefined
		default:
			return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return i.ModuloSmallInt(other.AsSmallInt())
	case FLOAT_FLAG:
		iFloat, _ := i.ToGoBigInt().Float64()
		return Float(math.Mod(iFloat, float64(other.AsFloat()))).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

func (i *BigInt) ModuloFloat(other Float) Float {
	iFloat, _ := i.ToGoBigInt().Float64()
	return Float(math.Mod(iFloat, float64(other)))
}

func (i *BigInt) ModuloBigFloat(other *BigFloat) *BigFloat {
	prec := max(other.Precision(), uint(i.BitSize()), 64)
	iBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(i)
	return iBigFloat.Mod(iBigFloat, other)
}

func (i *BigInt) ModuloBigInt(other *BigInt) (Value, Value) {
	if other.IsZero() {
		return Undefined, Ref(NewZeroDivisionError())
	}
	iGo := i.ToGoBigInt()
	oGo := other.ToGoBigInt()
	mod := &big.Int{}
	(&big.Int{}).QuoRem(iGo, oGo, mod)
	result := ToElkBigInt(mod)
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue(), Undefined
	}
	return Ref(result), Undefined
}

func (i *BigInt) ModuloSmallInt(other SmallInt) (Value, Value) {
	if other == 0 {
		return Undefined, Ref(NewZeroDivisionError())
	}
	iGo := i.ToGoBigInt()
	oBigInt := big.NewInt(int64(other))
	(&big.Int{}).QuoRem(iGo, oBigInt, oBigInt)
	result := ToElkBigInt(oBigInt)
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue(), Undefined
	}
	return Ref(result), Undefined
}

// Check whether i is greater than other and return an error
// if something went wrong.
func (i *BigInt) GreaterThanVal(other Value) (Value, Value) {
	result, err := i.GreaterThan(other)
	return ToElkBool(result), err
}

// Check whether i is greater than other and return an error
// if something went wrong.
func (i *BigInt) GreaterThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return i.Cmp(o) == 1, Undefined
		case *BigFloat:
			if o.IsNaN() {
				return false, Undefined
			}
			iBigFloat := (&BigFloat{}).SetBigInt(i)
			return iBigFloat.Cmp(o) == 1, Undefined
		default:
			return false, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}
	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		oBigInt := NewBigInt(int64(other.AsSmallInt()))
		return i.Cmp(oBigInt) == 1, Undefined
	case FLOAT_FLAG:
		return i.ToFloat() > other.AsFloat(), Undefined
	default:
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Check whether i is greater than or equal to other and return an error
// if something went wrong.
func (i *BigInt) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := i.GreaterThanEqual(other)
	return ToElkBool(result), err
}

// Check whether i is greater than or equal to other and return an error
// if something went wrong.
func (i *BigInt) GreaterThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return i.Cmp(o) >= 0, Undefined
		case *BigFloat:
			if o.IsNaN() {
				return false, Undefined
			}
			iBigFloat := (&BigFloat{}).SetBigInt(i)
			return iBigFloat.Cmp(o) >= 0, Undefined
		default:
			return false, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		oBigInt := NewBigInt(int64(other.AsSmallInt()))
		return i.Cmp(oBigInt) >= 0, Undefined
	case FLOAT_FLAG:
		return i.ToFloat() >= other.AsFloat(), Undefined
	default:
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Check whether i is less than other and return an error
// if something went wrong.
func (i *BigInt) LessThanVal(other Value) (Value, Value) {
	result, err := i.LessThan(other)
	return ToElkBool(result), err
}

// Check whether i is less than other and return an error
// if something went wrong.
func (i *BigInt) LessThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return i.Cmp(o) == -1, Undefined
		case *BigFloat:
			if o.IsNaN() {
				return false, Undefined
			}
			iBigFloat := (&BigFloat{}).SetBigInt(i)
			return iBigFloat.Cmp(o) == -1, Undefined
		default:
			return false, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		oBigInt := NewBigInt(int64(other.AsSmallInt()))
		return i.Cmp(oBigInt) == -1, Undefined
	case FLOAT_FLAG:
		return i.ToFloat() < other.AsFloat(), Undefined
	default:
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Check whether i is less than or equal to other and return an error
// if something went wrong.
func (i *BigInt) LessThanEqualVal(other Value) (Value, Value) {
	result, err := i.LessThanEqual(other)
	return ToElkBool(result), err
}

// Check whether i is less than or equal to other and return an error
// if something went wrong.
func (i *BigInt) LessThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return i.Cmp(o) <= 0, Undefined
		case *BigFloat:
			if o.IsNaN() {
				return false, Undefined
			}
			iBigFloat := (&BigFloat{}).SetBigInt(i)
			return iBigFloat.Cmp(o) <= 0, Undefined
		default:
			return false, Ref(NewCoerceError(i.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		oBigInt := NewBigInt(int64(other.AsSmallInt()))
		return i.Cmp(oBigInt) <= 0, Undefined
	case FLOAT_FLAG:
		return i.ToFloat() <= other.AsFloat(), Undefined
	default:
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Check whether i is equal to other (with coercion)
func (i *BigInt) LaxEqualVal(other Value) Value {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return ToElkBool(i.Cmp(o) == 0)
		case *BigFloat:
			if o.IsNaN() {
				return False
			}
			iBigFloat := (&BigFloat{}).SetBigInt(i)
			return ToElkBool(iBigFloat.Cmp(o) == 0)
		case Int64:
			oBigInt := NewBigInt(int64(o))
			return ToElkBool(i.Cmp(oBigInt) == 0)
		case UInt64:
			oBigInt := NewBigInt(int64(o))
			return ToElkBool(i.Cmp(oBigInt) == 0)
		case Float64:
			return ToElkBool(i.ToFloat() == Float(o))
		default:
			return False
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		oBigInt := NewBigInt(int64(other.AsSmallInt()))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case FLOAT_FLAG:
		return ToElkBool(i.ToFloat() == other.AsFloat())
	case INT64_FLAG:
		oBigInt := NewBigInt(int64(other.AsInlineInt64()))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case INT32_FLAG:
		oBigInt := NewBigInt(int64(other.AsInt32()))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case INT16_FLAG:
		oBigInt := NewBigInt(int64(other.AsInt16()))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case INT8_FLAG:
		oBigInt := NewBigInt(int64(other.AsInt8()))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case UINT64_FLAG:
		oBigInt := NewBigInt(int64(other.AsInlineUInt64()))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case UINT32_FLAG:
		oBigInt := NewBigInt(int64(other.AsUInt32()))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case UINT16_FLAG:
		oBigInt := NewBigInt(int64(other.AsUInt16()))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case UINT8_FLAG:
		oBigInt := NewBigInt(int64(other.AsUInt8()))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case FLOAT64_FLAG:
		return ToElkBool(i.ToFloat() == Float(other.AsInlineFloat64()))
	case FLOAT32_FLAG:
		return ToElkBool(i.ToFloat() == Float(other.AsFloat32()))
	default:
		return False
	}
}

// Check whether i is equal to other
func (i *BigInt) Equal(other Value) bool {
	if other.IsSmallInt() {
		oBigInt := NewBigInt(int64(other.AsSmallInt()))
		return i.Cmp(oBigInt) == 0
	}
	if !other.IsReference() {
		return false
	}

	switch o := other.AsReference().(type) {
	case *BigInt:
		return i.Cmp(o) == 0
	default:
		return false
	}
}

// Check whether i is equal to other
func (i *BigInt) EqualVal(other Value) Value {
	return ToElkBool(i.Equal(other))
}

// Check whether i is strictly equal to other
func (i *BigInt) StrictEqualVal(other Value) Value {
	return i.EqualVal(other)
}

func rightBitshiftBigInt[T SimpleInt](i *BigInt, other T) Value {
	if other < 0 {
		return SmallInt(0).ToValue()
	}
	iGo := i.ToGoBigInt()
	result := ToElkBigInt(iGo.Rsh(iGo, uint(other)))
	if result.IsSmallInt() {
		return result.ToSmallInt().ToValue()
	}
	return Ref(result)
}

// Bitshift to the right by another integer value and return an error
// if something went wrong.
func (i *BigInt) RightBitshiftVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case Int64:
			if o < 0 {
				return leftBitshiftBigInt(i, -o), Undefined
			}
			return rightBitshiftBigInt(i, o), Undefined
		case UInt64:
			return rightBitshiftBigInt(i, o), Undefined
		case *BigInt:
			if o.IsSmallInt() {
				oSmall := o.ToSmallInt()
				if oSmall < 0 {
					return leftBitshiftBigInt(i, -oSmall), Undefined
				}
				return rightBitshiftBigInt(i, oSmall), Undefined
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
			return leftBitshiftBigInt(i, -o), Undefined
		}
		return rightBitshiftBigInt(i, o), Undefined
	case INT64_FLAG:
		o := other.AsInlineInt64()
		if o < 0 {
			return leftBitshiftBigInt(i, -o), Undefined
		}
		return rightBitshiftBigInt(i, o), Undefined
	case INT32_FLAG:
		o := other.AsInt32()
		if o < 0 {
			return leftBitshiftBigInt(i, -o), Undefined
		}
		return rightBitshiftBigInt(i, o), Undefined
	case INT16_FLAG:
		o := other.AsInt16()
		if o < 0 {
			return leftBitshiftBigInt(i, -o), Undefined
		}
		return rightBitshiftBigInt(i, o), Undefined
	case INT8_FLAG:
		o := other.AsInt8()
		if o < 0 {
			return leftBitshiftBigInt(i, -o), Undefined
		}
		return rightBitshiftBigInt(i, o), Undefined
	case UINT64_FLAG:
		o := other.AsInlineUInt64()
		return rightBitshiftBigInt(i, o), Undefined
	case UINT32_FLAG:
		o := other.AsUInt32()
		return rightBitshiftBigInt(i, o), Undefined
	case UINT16_FLAG:
		o := other.AsUInt16()
		return rightBitshiftBigInt(i, o), Undefined
	case UINT8_FLAG:
		o := other.AsUInt8()
		return rightBitshiftBigInt(i, o), Undefined
	default:
		return Undefined, Ref(NewBitshiftOperandError(other))
	}
}

func leftBitshiftBigInt[T SimpleInt](i *BigInt, other T) Value {
	if other < 0 {
		return SmallInt(0).ToValue()
	}
	iGo := i.ToGoBigInt()
	return Ref(ToElkBigInt(iGo.Lsh(iGo, uint(other))))
}

// Bitshift to the left by another integer value and return an error
// if something went wrong.
func (i *BigInt) LeftBitshiftVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case Int64:
			return leftBitshiftBigInt(i, o), Undefined
		case UInt64:
			return leftBitshiftBigInt(i, o), Undefined
		case *BigInt:
			if o.IsSmallInt() {
				oSmall := o.ToSmallInt()
				return leftBitshiftBigInt(i, oSmall), Undefined
			}
			return SmallInt(0).ToValue(), Undefined
		default:
			return Undefined, Ref(NewBitshiftOperandError(other))
		}
	}
	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		return leftBitshiftBigInt(i, o), Undefined
	case INT64_FLAG:
		o := other.AsInlineInt64()
		return leftBitshiftBigInt(i, o), Undefined
	case INT32_FLAG:
		o := other.AsInt32()
		return leftBitshiftBigInt(i, o), Undefined
	case INT16_FLAG:
		o := other.AsInt16()
		return leftBitshiftBigInt(i, o), Undefined
	case INT8_FLAG:
		o := other.AsInt8()
		return leftBitshiftBigInt(i, o), Undefined
	case UINT64_FLAG:
		o := other.AsInlineUInt64()
		return leftBitshiftBigInt(i, o), Undefined
	case UINT32_FLAG:
		o := other.AsUInt32()
		return leftBitshiftBigInt(i, o), Undefined
	case UINT16_FLAG:
		o := other.AsUInt16()
		return leftBitshiftBigInt(i, o), Undefined
	case UINT8_FLAG:
		o := other.AsUInt8()
		return leftBitshiftBigInt(i, o), Undefined
	default:
		return Undefined, Ref(NewBitshiftOperandError(other))
	}
}

// Perform bitwise AND with another value and return an error
// if something went wrong.
func (i *BigInt) BitwiseAndVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			result := ToElkBigInt((&big.Int{}).And(i.ToGoBigInt(), o.ToGoBigInt()))
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
		oBigInt := big.NewInt(int64(other.AsSmallInt()))
		oBigInt.And(i.ToGoBigInt(), oBigInt)
		result := ToElkBigInt(oBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt().ToValue(), Undefined
		}
		return Ref(result), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Perform bitwise AND NOT with another value and return an error
// if something went wrong.
func (i *BigInt) BitwiseAndNotVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			result := ToElkBigInt((&big.Int{}).AndNot(i.ToGoBigInt(), o.ToGoBigInt()))
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
		oBigInt := big.NewInt(int64(other.AsSmallInt()))
		oBigInt.AndNot(i.ToGoBigInt(), oBigInt)
		result := ToElkBigInt(oBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt().ToValue(), Undefined
		}
		return Ref(result), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Perform bitwise OR with another value and return an error
// if something went wrong.
func (i *BigInt) BitwiseOrVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			result := ToElkBigInt((&big.Int{}).Or(i.ToGoBigInt(), o.ToGoBigInt()))
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
		oBigInt := big.NewInt(int64(other.AsSmallInt()))
		oBigInt.Or(i.ToGoBigInt(), oBigInt)
		result := ToElkBigInt(oBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt().ToValue(), Undefined
		}
		return Ref(result), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

// Perform bitwise XOR with another value and return an error
// if something went wrong.
func (i *BigInt) BitwiseXorVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			result := ToElkBigInt((&big.Int{}).Xor(i.ToGoBigInt(), o.ToGoBigInt()))
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
		oBigInt := big.NewInt(int64(other.AsSmallInt()))
		oBigInt.Xor(i.ToGoBigInt(), oBigInt)
		result := ToElkBigInt(oBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt().ToValue(), Undefined
		}
		return Ref(result), Undefined
	default:
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
}

func (*BigInt) Class() *Class {
	return IntClass
}

func (*BigInt) DirectClass() *Class {
	return IntClass
}

func (*BigInt) SingletonClass() *Class {
	return nil
}

func (i *BigInt) Inspect() string {
	return i.ToGoBigInt().String()
}

func (i *BigInt) Error() string {
	return i.Inspect()
}

func (i *BigInt) Copy() Reference {
	return i
}

func (i *BigInt) InstanceVariables() *InstanceVariables {
	return nil
}

func (i *BigInt) Hash() UInt64 {
	d := xxhash.New()
	d.Write(i.ToGoBigInt().Bytes())
	return UInt64(d.Sum64())
}

// Parses an unsigned big.Int from a string using Elk syntax.
func parseUBigInt(s string, base int, formatErr *Class) (*BigInt, Value) {
	if s == "" {
		return nil, Ref(NewError(formatErr, "invalid integer format"))
	}

	switch {
	case 2 <= base && base <= 36:
	case base == 0:
		// Look for binary, quaternary, octal, duodecimal, hex prefix.
		base = 10
		if s[0] == '0' {
			switch {
			case len(s) >= 3 && letterToLower(s[1]) == 'b':
				// binary int
				base = 2
				s = s[2:]
			case len(s) >= 3 && letterToLower(s[1]) == 'q':
				// quaternary int
				base = 4
				s = s[2:]
			case len(s) >= 3 && letterToLower(s[1]) == 'o':
				// octal int
				base = 8
				s = s[2:]
			case len(s) >= 3 && letterToLower(s[1]) == 'd':
				// duodecimal int
				base = 12
				s = s[2:]
			case len(s) >= 3 && letterToLower(s[1]) == 'x':
				// hexadecimal int
				base = 16
				s = s[2:]
			}
		}
	default:
		return nil, Ref(Errorf(formatErr, "invalid integer base %d", base))
	}

	n := &big.Int{}
	for _, c := range []byte(s) {
		var d byte
		switch {
		case c == '_':
			continue
		case '0' <= c && c <= '9':
			d = c - '0'
		case 'a' <= letterToLower(c) && letterToLower(c) <= 'z':
			d = letterToLower(c) - 'a' + 10
		default:
			return nil, Ref(Errorf(formatErr, "illegal characters in integer: %c", c))
		}

		if d >= byte(base) {
			return nil, Ref(Errorf(formatErr, "illegal characters in integer (base %d): %c", base, c))
		}

		n.Mul(n, big.NewInt(int64(base)))

		n.Add(n, big.NewInt(int64(d)))
	}

	return ToElkBigInt(n), Undefined
}

// Parses a signed big.Int from a string using Elk syntax.
func ParseBigInt(s string, base int) (*BigInt, Value) {
	return ParseBigIntWithErr(s, base, FormatErrorClass)
}

// Parses a signed big.Int from a string using Elk syntax.
func ParseBigIntWithErr(s string, base int, formatError *Class) (*BigInt, Value) {
	if s == "" {
		return nil, Ref(NewError(formatError, "invalid integer format"))
	}

	// Pick off leading sign.
	neg := false
	switch s[0] {
	case '+':
		s = s[1:]
	case '-':
		neg = true
		s = s[1:]
	}

	// Convert unsigned and check range.
	u, err := parseUBigInt(s, base, formatError)
	un := u.ToGoBigInt()

	if !err.IsUndefined() {
		return nil, err
	}

	if neg {
		un.Neg(un)
	}

	return ToElkBigInt(un), Undefined
}

func ParseIntWithErr(s string, base int, formatErr *Class) (Value, Value) {
	val, err := ParseBigIntWithErr(s, base, formatErr)
	if !err.IsUndefined() {
		return Undefined, err
	}

	if val.IsSmallInt() {
		return val.ToSmallInt().ToValue(), Undefined
	}

	return Ref(val), Undefined
}

func ParseInt(s string, base int) (Value, Value) {
	return ParseIntWithErr(s, base, FormatErrorClass)
}

// Same as [ParseBigInt] but panics on error.
func ParseBigIntPanic(s string, base int) *BigInt {
	result, err := ParseBigInt(s, base)
	if !err.IsUndefined() {
		panic(err)
	}

	return result
}

func (i *BigInt) Nanoseconds() TimeSpan {
	return TimeSpan(i.ToSmallInt())
}

func (i *BigInt) Microseconds() TimeSpan {
	oBigInt := big.NewInt(int64(Microsecond))
	oBigInt.Mul(i.ToGoBigInt(), oBigInt)
	return TimeSpan(i.ToSmallInt())
}

func (i *BigInt) Milliseconds() TimeSpan {
	oBigInt := big.NewInt(int64(Millisecond))
	oBigInt.Mul(i.ToGoBigInt(), oBigInt)
	return TimeSpan(i.ToSmallInt())
}

func (i *BigInt) Seconds() TimeSpan {
	oBigInt := big.NewInt(int64(Second))
	oBigInt.Mul(i.ToGoBigInt(), oBigInt)
	return TimeSpan(i.ToSmallInt())
}

func (i *BigInt) Minutes() TimeSpan {
	oBigInt := big.NewInt(int64(Minute))
	oBigInt.Mul(i.ToGoBigInt(), oBigInt)
	return TimeSpan(i.ToSmallInt())
}

func (i *BigInt) Hours() TimeSpan {
	oBigInt := big.NewInt(int64(Hour))
	oBigInt.Mul(i.ToGoBigInt(), oBigInt)
	return TimeSpan(i.ToSmallInt())
}

func (i *BigInt) Days() DateSpan {
	return MakeDateSpan(0, 0, int(i.ToSmallInt()))
}

func (i *BigInt) Weeks() DateSpan {
	oBigInt := big.NewInt(7)
	oBigInt.Mul(i.ToGoBigInt(), oBigInt)
	return MakeDateSpan(0, 0, int(oBigInt.Int64()))
}

func (i *BigInt) Months() DateSpan {
	return MakeDateSpan(0, int(i.ToSmallInt()), 0)
}

func (i *BigInt) Years() DateSpan {
	return MakeDateSpan(int(i.ToSmallInt()), 0, 0)
}

func (i *BigInt) Centuries() DateSpan {
	oBigInt := big.NewInt(100)
	oBigInt.Mul(i.ToGoBigInt(), oBigInt)
	return MakeDateSpan(int(oBigInt.Int64()), 0, 0)
}

func (i *BigInt) Millenia() DateSpan {
	oBigInt := big.NewInt(1000)
	oBigInt.Mul(i.ToGoBigInt(), oBigInt)
	return MakeDateSpan(int(oBigInt.Int64()), 0, 0)
}

type BigIntIterator struct {
	Int     *BigInt
	Counter Value
}

func NewBigIntIterator(i *BigInt) *BigIntIterator {
	return &BigIntIterator{
		Int:     i,
		Counter: SmallInt(0).ToValue(),
	}
}

func NewBigIntIteratorWithCounter(i *BigInt, counter Value) *BigIntIterator {
	return &BigIntIterator{
		Int:     i,
		Counter: counter,
	}
}

func (*BigIntIterator) Class() *Class {
	return IntIteratorClass
}

func (*BigIntIterator) DirectClass() *Class {
	return IntIteratorClass
}

func (*BigIntIterator) SingletonClass() *Class {
	return nil
}

func (l *BigIntIterator) Copy() Reference {
	return &BigIntIterator{
		Int:     l.Int,
		Counter: l.Counter,
	}
}

func (l *BigIntIterator) Inspect() string {
	return fmt.Sprintf("Std::Int::Iterator{&: %p, int: %s, counter: %s}", l, l.Int.Inspect(), l.Counter.Inspect())
}

func (l *BigIntIterator) Error() string {
	return l.Inspect()
}

func (*BigIntIterator) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *BigIntIterator) Next() (Value, Value) {
	stop, err := l.Int.LessThan(l.Counter)
	if !err.IsUndefined() {
		return Undefined, err
	}
	if stop {
		return Undefined, stopIterationSymbol.ToValue()
	}

	next := l.Counter
	l.Counter = IncrementVal(l.Counter)
	return next, Undefined
}

func (l *BigIntIterator) Reset() {
	l.Counter = SmallInt(0).ToValue()
}
