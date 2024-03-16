package value

import (
	"encoding/binary"
	"math"
	"math/big"
	"strconv"

	"github.com/cespare/xxhash/v2"
)

// Elk's SmallInt value
type SmallInt int64

// Number of bits available for a small int.
const SmallIntBits = 64

// Max value of SmallInt
const MaxSmallInt = math.MaxInt64

// Min value of SmallInt
const MinSmallInt = math.MinInt64

func (SmallInt) Class() *Class {
	return IntClass
}

func (SmallInt) DirectClass() *Class {
	return IntClass
}

func (SmallInt) SingletonClass() *Class {
	return nil
}

func (i SmallInt) Copy() Value {
	return i
}

func (i SmallInt) Inspect() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i SmallInt) InstanceVariables() SymbolMap {
	return nil
}

// Negate the number and return the result.
func (i SmallInt) Negate() Value {
	if i == math.MinInt64 {
		iBigInt := big.NewInt(int64(i))
		return ToElkBigInt(iBigInt.Neg(iBigInt))
	}

	return -i
}

// Increment the number and return the result.
func (i SmallInt) Increment() Value {
	result, ok := i.AddOverflow(1)
	if !ok {
		iBigInt := big.NewInt(int64(i))
		return ToElkBigInt(iBigInt.Add(iBigInt, big.NewInt(1)))
	}
	return result
}

// Decrement the number and return the result.
func (i SmallInt) Decrement() Value {
	result, ok := i.SubtractOverflow(1)
	if !ok {
		iBigInt := big.NewInt(int64(i))
		return ToElkBigInt(iBigInt.Sub(iBigInt, big.NewInt(1)))
	}
	return result
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

// Add another value and return an error
// if something went wrong.
func (i SmallInt) Add(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		result, ok := i.AddOverflow(o)
		if !ok {
			iBigInt := big.NewInt(int64(i))
			return ToElkBigInt(iBigInt.Add(iBigInt, big.NewInt(int64(o)))), nil
		}
		return result, nil
	case *BigInt:
		iBigInt := big.NewInt(int64(i))
		iBigInt.Add(iBigInt, o.ToGoBigInt())
		if iBigInt.IsInt64() {
			return SmallInt(iBigInt.Int64()), nil
		}
		return ToElkBigInt(iBigInt), nil
	case Float:
		return Float(i) + o, nil
	case *BigFloat:
		prec := max(o.Precision(), 64)
		iBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(i)
		iBigFloat.AddBigFloat(iBigFloat, o)
		return iBigFloat, nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
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
func (i SmallInt) Subtract(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		result, ok := i.SubtractOverflow(o)
		if !ok {
			iBigInt := big.NewInt(int64(i))
			return ToElkBigInt(iBigInt.Sub(iBigInt, big.NewInt(int64(o)))), nil
		}
		return result, nil
	case *BigInt:
		iBigInt := big.NewInt(int64(i))
		iBigInt.Sub(iBigInt, o.ToGoBigInt())
		if iBigInt.IsInt64() {
			return SmallInt(iBigInt.Int64()), nil
		}
		return ToElkBigInt(iBigInt), nil
	case Float:
		return Float(i) - o, nil
	case *BigFloat:
		prec := max(o.Precision(), 64)
		iBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(i)
		iBigFloat.SubBigFloat(iBigFloat, o)
		return iBigFloat, nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
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

// Multiply another value and return an error
// if something went wrong.
func (i SmallInt) Multiply(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		result, ok := i.MultiplyOverflow(o)
		if !ok {
			iBigInt := big.NewInt(int64(i))
			return ToElkBigInt(iBigInt.Mul(iBigInt, big.NewInt(int64(o)))), nil
		}
		return result, nil
	case *BigInt:
		iBigInt := big.NewInt(int64(i))
		iBigInt.Mul(iBigInt, o.ToGoBigInt())
		if iBigInt.IsInt64() {
			return SmallInt(iBigInt.Int64()), nil
		}
		return ToElkBigInt(iBigInt), nil
	case Float:
		return Float(i) * o, nil
	case *BigFloat:
		prec := max(o.Precision(), 64)
		iBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(i)
		return iBigFloat.MulBigFloat(iBigFloat, o), nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
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

// Divide another value and return an error
// if something went wrong.
func (i SmallInt) Divide(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		if o == 0 {
			return nil, NewZeroDivisionError()
		}
		result, ok := i.DivideOverflow(o)
		if !ok {
			iBigInt := big.NewInt(int64(i))
			return ToElkBigInt(iBigInt.Div(iBigInt, big.NewInt(int64(o)))), nil
		}
		return result, nil
	case *BigInt:
		if o.IsZero() {
			return nil, NewZeroDivisionError()
		}
		iBigInt := big.NewInt(int64(i))
		iBigInt.Div(iBigInt, o.ToGoBigInt())
		if iBigInt.IsInt64() {
			return SmallInt(iBigInt.Int64()), nil
		}
		return ToElkBigInt(iBigInt), nil
	case Float:
		return Float(i) / o, nil
	case *BigFloat:
		prec := max(o.Precision(), 64)
		iBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(i)
		iBigFloat.DivBigFloat(iBigFloat, o)
		return iBigFloat, nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
	}
}

// Exponentiate by another value and return an error
// if something went wrong.
func (i SmallInt) Exponentiate(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		iBigInt := big.NewInt(int64(i))
		oBigInt := big.NewInt(int64(o))
		iBigInt.Exp(iBigInt, oBigInt, nil)
		if iBigInt.IsInt64() {
			return SmallInt(iBigInt.Int64()), nil
		}
		return ToElkBigInt(iBigInt), nil
	case *BigInt:
		iBigInt := big.NewInt(int64(i))
		iBigInt.Exp(iBigInt, o.ToGoBigInt(), nil)
		if iBigInt.IsInt64() {
			return SmallInt(iBigInt.Int64()), nil
		}
		return ToElkBigInt(iBigInt), nil
	case Float:
		return Float(math.Pow(float64(i), float64(o))), nil
	case *BigFloat:
		prec := max(o.Precision(), 64)
		iBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(i)
		iBigFloat.ExpBigFloat(iBigFloat, o)
		return iBigFloat, nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
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
func (i SmallInt) Compare(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return SmallInt(i.Cmp(o)), nil
	case *BigInt:
		iBigInt := NewBigInt(int64(i))
		return SmallInt(iBigInt.Cmp(o)), nil
	case Float:
		if o.IsNaN() {
			return Nil, nil
		}
		return SmallInt(i.ToFloat().Cmp(o)), nil
	case *BigFloat:
		if o.IsNaN() {
			return Nil, nil
		}
		iBigFloat := (&BigFloat{}).SetSmallInt(i)
		return SmallInt(iBigFloat.Cmp(o)), nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
	}
}

// Check whether i is greater than other and return an error
// if something went wrong.
func (i SmallInt) GreaterThan(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return ToElkBool(i > o), nil
	case *BigInt:
		iBigInt := NewBigInt(int64(i))
		return ToElkBool(iBigInt.Cmp(o) == 1), nil
	case Float:
		return ToElkBool(Float(i) > o), nil
	case *BigFloat:
		if o.IsNaN() {
			return False, nil
		}
		iBigFloat := (&BigFloat{}).SetSmallInt(i)
		return ToElkBool(iBigFloat.Cmp(o) == 1), nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
	}
}

// Check whether i is greater than or equal to other and return an error
// if something went wrong.
func (i SmallInt) GreaterThanEqual(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return ToElkBool(i >= o), nil
	case *BigInt:
		iBigInt := NewBigInt(int64(i))
		return ToElkBool(iBigInt.Cmp(o) >= 0), nil
	case Float:
		return ToElkBool(Float(i) >= o), nil
	case *BigFloat:
		if o.IsNaN() {
			return False, nil
		}
		iBigFloat := (&BigFloat{}).SetSmallInt(i)
		return ToElkBool(iBigFloat.Cmp(o) >= 0), nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
	}
}

// Check whether i is less than other and return an error
// if something went wrong.
func (i SmallInt) LessThan(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return ToElkBool(i < o), nil
	case *BigInt:
		iBigInt := NewBigInt(int64(i))
		return ToElkBool(iBigInt.Cmp(o) == -1), nil
	case Float:
		return ToElkBool(Float(i) < o), nil
	case *BigFloat:
		if o.IsNaN() {
			return False, nil
		}
		iBigFloat := (&BigFloat{}).SetSmallInt(i)
		return ToElkBool(iBigFloat.Cmp(o) == -1), nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
	}
}

// Check whether i is less than or equal to other and return an error
// if something went wrong.
func (i SmallInt) LessThanEqual(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return ToElkBool(i <= o), nil
	case *BigInt:
		iBigInt := NewBigInt(int64(i))
		return ToElkBool(iBigInt.Cmp(o) <= 0), nil
	case Float:
		return ToElkBool(Float(i) <= o), nil
	case *BigFloat:
		if o.IsNaN() {
			return False, nil
		}
		iBigFloat := (&BigFloat{}).SetSmallInt(i)
		return ToElkBool(iBigFloat.Cmp(o) <= 0), nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
	}
}

// Check whether i is equal to other (with coercion)
func (i SmallInt) LaxEqual(other Value) Value {
	switch o := other.(type) {
	case SmallInt:
		return ToElkBool(i == o)
	case *BigInt:
		iBigInt := NewBigInt(int64(i))
		return ToElkBool(iBigInt.Cmp(o) == 0)
	case Float:
		return ToElkBool(Float(i) == o)
	case *BigFloat:
		if o.IsNaN() {
			return False
		}
		iBigFloat := (&BigFloat{}).SetSmallInt(i)
		return ToElkBool(iBigFloat.Cmp(o) == 0)
	case Int64:
		return ToElkBool(i == SmallInt(o))
	case Int32:
		return ToElkBool(i == SmallInt(o))
	case Int16:
		return ToElkBool(i == SmallInt(o))
	case Int8:
		return ToElkBool(i == SmallInt(o))
	case UInt64:
		if o > MaxSmallInt {
			return False
		}
		return ToElkBool(i == SmallInt(o))
	case UInt32:
		return ToElkBool(i == SmallInt(o))
	case UInt16:
		return ToElkBool(i == SmallInt(o))
	case UInt8:
		return ToElkBool(i == SmallInt(o))
	case Float64:
		return ToElkBool(Float64(i) == o)
	case Float32:
		return ToElkBool(Float32(i) == o)
	default:
		return False
	}
}

// Check whether i is equal to other
func (i SmallInt) Equal(other Value) Value {
	return i.StrictEqual(other)
}

// Check whether i is strictly equal to other
func (i SmallInt) StrictEqual(other Value) Value {
	switch o := other.(type) {
	case SmallInt:
		return ToElkBool(i == o)
	case *BigInt:
		iBigInt := NewBigInt(int64(i))
		return ToElkBool(iBigInt.Cmp(o) == 0)
	default:
		return False
	}
}

func leftBitshiftSmallInt[T SimpleInt](i SmallInt, other T) Value {
	var bitsize T = SmallIntBits - 1
	if other < 0 {
		return SmallInt(0)
	}
	complementaryShift := i >> (bitsize - other)
	if other > bitsize || (i < 0 && complementaryShift != -1) || (i > 0 && complementaryShift != 0) {
		// overflow
		iBig := big.NewInt(int64(i))
		iBig.Lsh(iBig, uint(other))
		return ToElkBigInt(iBig)
	}
	return i << other
}

func rightBitshiftSmallInt[T SimpleInt](i SmallInt, other T) Value {
	if other < 0 {
		return SmallInt(0)
	}
	return i >> other
}

// Bitshift to the left by another integer value and return an error
// if something went wrong.
func (i SmallInt) LeftBitshift(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		if o < 0 {
			return rightBitshiftSmallInt(i, -o), nil
		}
		return leftBitshiftSmallInt(i, o), nil
	case Int64:
		if o < 0 {
			return rightBitshiftSmallInt(i, -o), nil
		}
		return leftBitshiftSmallInt(i, o), nil
	case Int32:
		if o < 0 {
			return rightBitshiftSmallInt(i, -o), nil
		}
		return leftBitshiftSmallInt(i, o), nil
	case Int16:
		if o < 0 {
			return rightBitshiftSmallInt(i, -o), nil
		}
		return leftBitshiftSmallInt(i, o), nil
	case Int8:
		if o < 0 {
			return rightBitshiftSmallInt(i, -o), nil
		}
		return leftBitshiftSmallInt(i, o), nil
	case UInt64:
		return leftBitshiftSmallInt(i, o), nil
	case UInt32:
		return leftBitshiftSmallInt(i, o), nil
	case UInt16:
		return leftBitshiftSmallInt(i, o), nil
	case UInt8:
		return leftBitshiftSmallInt(i, o), nil
	case *BigInt:
		if o.IsSmallInt() {
			oSmall := o.ToSmallInt()
			if oSmall < 0 {
				return rightBitshiftSmallInt(i, -oSmall), nil
			}
			return leftBitshiftSmallInt(i, oSmall), nil
		}
		return SmallInt(0), nil
	default:
		return nil, NewBitshiftOperandError(other)
	}
}

// Bitshift to the right by another integer value and return an error
// if something went wrong.
func (i SmallInt) RightBitshift(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), nil
		}
		return i >> o, nil
	case Int64:
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), nil
		}
		return i >> o, nil
	case Int32:
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), nil
		}
		return i >> o, nil
	case Int16:
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), nil
		}
		return i >> o, nil
	case Int8:
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), nil
		}
		return i >> o, nil
	case UInt64:
		return i >> o, nil
	case UInt32:
		return i >> o, nil
	case UInt16:
		return i >> o, nil
	case UInt8:
		return i >> o, nil
	case *BigInt:
		if o.IsSmallInt() {
			oSmall := o.ToSmallInt()
			if oSmall < 0 {
				return leftBitshiftSmallInt(i, -oSmall), nil
			}
			return i >> oSmall, nil
		}
		return SmallInt(0), nil
	default:
		return nil, NewBitshiftOperandError(other)
	}
}

// Perform a bitwise AND with another integer value and return an error
// if something went wrong.
func (i SmallInt) BitwiseAnd(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return i & o, nil
	case *BigInt:
		iBigInt := big.NewInt(int64(i))
		iBigInt.And(iBigInt, o.ToGoBigInt())
		result := ToElkBigInt(iBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
	}
}

// Perform a bitwise AND NOT with another integer value and return an error
// if something went wrong.
func (i SmallInt) BitwiseAndNot(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return i &^ o, nil
	case *BigInt:
		iBigInt := big.NewInt(int64(i))
		iBigInt.AndNot(iBigInt, o.ToGoBigInt())
		result := ToElkBigInt(iBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
	}
}

// Perform a bitwise OR with another integer value and return an error
// if something went wrong.
func (i SmallInt) BitwiseOr(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return i | o, nil
	case *BigInt:
		iBigInt := big.NewInt(int64(i))
		iBigInt.Or(iBigInt, o.ToGoBigInt())
		result := ToElkBigInt(iBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
	}
}

// Perform a bitwise XOR with another integer value and return an error
// if something went wrong.
func (i SmallInt) BitwiseXor(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return i ^ o, nil
	case *BigInt:
		iBigInt := big.NewInt(int64(i))
		iBigInt.Xor(iBigInt, o.ToGoBigInt())
		result := ToElkBigInt(iBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
	}
}

// Perform modulo by another numeric value and return an error
// if something went wrong.
func (i SmallInt) Modulo(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return i % o, nil
	case *BigInt:
		if o.IsSmallInt() {
			oSmall := o.ToSmallInt()
			return i % oSmall, nil
		}
		return i, nil
	case Float:
		return Float(math.Mod(float64(i), float64(o))), nil
	case *BigFloat:
		prec := max(o.Precision(), 64)
		iBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(i)
		return iBigFloat.Mod(iBigFloat, o), nil
	default:
		return nil, NewCoerceError(i.Class(), other.Class())
	}
}

func (i SmallInt) Hash() UInt64 {
	d := xxhash.New()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	d.Write(b)
	return UInt64(d.Sum64())
}
