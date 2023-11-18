package value

import (
	"fmt"
	"math"
	"math/big"
)

var SmallIntClass *Class // ::Std::SmallInt

// Elk's SmallInt value
type SmallInt int64

// Number of bits available for a small int.
const SmallIntBits = 64

// Max value of SmallInt
const MaxSmallInt = math.MaxInt64

// Min value of SmallInt
const MinSmallInt = math.MinInt64

func (SmallInt) Class() *Class {
	return SmallIntClass
}

func (SmallInt) DirectClass() *Class {
	return SmallIntClass
}

func (SmallInt) SingletonClass() *Class {
	return nil
}

func (SmallInt) IsFrozen() bool {
	return true
}

func (i SmallInt) SetFrozen() {}

func (i SmallInt) Inspect() string {
	return fmt.Sprintf("%d", i)
}

func (i SmallInt) InstanceVariables() SimpleSymbolMap {
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

// Add two small ints and check for overflow/underflow.
func (a SmallInt) AddOverflow(b SmallInt) (result SmallInt, ok bool) {
	c := a + b
	if (c > a) == (b > 0) {
		return c, true
	}
	return c, false
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
		return nil, NewCoerceError(i, other)
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
		return nil, NewCoerceError(i, other)
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
		return nil, NewCoerceError(i, other)
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
		return nil, NewCoerceError(i, other)
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
		return nil, NewCoerceError(i, other)
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
		return nil, NewCoerceError(i, other)
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
		return nil, NewCoerceError(i, other)
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
		return nil, NewCoerceError(i, other)
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
		return nil, NewCoerceError(i, other)
	}
}

// Check whether i is equal to other
func (i SmallInt) Equal(other Value) Value {
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
		return nil, NewCoerceError(i, other)
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
		return nil, NewCoerceError(i, other)
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
		return nil, NewCoerceError(i, other)
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
		return nil, NewCoerceError(i, other)
	}
}

func initSmallInt() {
	SmallIntClass = NewClassWithOptions(
		ClassWithParent(IntClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithSingleton(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("SmallInt", SmallIntClass)
}
