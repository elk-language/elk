package object

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ALTree/bigfloat"
)

var SmallIntClass *Class // ::Std::SmallInt

// Elk's SmallInt value
type SmallInt int64

// Number of bits available for a small int.
const SmallIntBits = 64

func (i SmallInt) Class() *Class {
	return SmallIntClass
}

func (i SmallInt) IsFrozen() bool {
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
		iBigFloat := (&big.Float{}).SetPrec(prec).SetInt64(int64(i))
		iBigFloat.Add(iBigFloat, o.ToGoBigFloat())
		return ToElkBigFloat(iBigFloat), nil
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
		iBigFloat := (&big.Float{}).SetPrec(prec).SetInt64(int64(i))
		iBigFloat.Sub(iBigFloat, o.ToGoBigFloat())
		return ToElkBigFloat(iBigFloat), nil
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
		iBigFloat := (&big.Float{}).SetPrec(prec).SetInt64(int64(i))
		iBigFloat.Mul(iBigFloat, o.ToGoBigFloat())
		return ToElkBigFloat(iBigFloat), nil
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
		if len(o.ToGoBigInt().Bits()) == 0 {
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
		iBigFloat := (&big.Float{}).SetPrec(prec).SetInt64(int64(i))
		iBigFloat.Quo(iBigFloat, o.ToGoBigFloat())
		return ToElkBigFloat(iBigFloat), nil
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
		iBigFloat := (&big.Float{}).SetPrec(prec).SetInt64(int64(i))
		iBigFloat = bigfloat.Pow(iBigFloat, o.ToGoBigFloat())
		return ToElkBigFloat(iBigFloat), nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

func leftBitshiftSmallInt[T SimpleInt](i SmallInt, other T) Value {
	var bitsize T = SmallIntBits - 1
	if other < 0 {
		return SmallInt(0)
	}
	if other > bitsize || i>>(bitsize-other) != 0 {
		// overflow
		iBig := big.NewInt(int64(i))
		iBig.Lsh(iBig, uint(other))
		return ToElkBigInt(iBig)
	}
	return i << other
}

// Bitshift to the left by another integer value and return an error
// if something went wrong.
func (i SmallInt) LeftBitshift(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		if o < 0 {
			return i >> -o, nil
		}
		return leftBitshiftSmallInt(i, o), nil
	case Int64:
		if o < 0 {
			return i >> -o, nil
		}
		return leftBitshiftSmallInt(i, o), nil
	case Int32:
		if o < 0 {
			return i >> -o, nil
		}
		return leftBitshiftSmallInt(i, o), nil
	case Int16:
		if o < 0 {
			return i >> -o, nil
		}
		return leftBitshiftSmallInt(i, o), nil
	case Int8:
		if o < 0 {
			return i >> -o, nil
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
				return i >> -oSmall, nil
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

// Logically bitshift to the right by another integer value and return an error
// if something went wrong.
func (i SmallInt) LogicalRightBitshift(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), nil
		}
		return SmallInt(uint64(i) >> o), nil
	case Int64:
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), nil
		}
		return SmallInt(uint64(i) >> o), nil
	case Int32:
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), nil
		}
		return SmallInt(uint64(i) >> o), nil
	case Int16:
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), nil
		}
		return SmallInt(uint64(i) >> o), nil
	case Int8:
		if o < 0 {
			return leftBitshiftSmallInt(i, -o), nil
		}
		return SmallInt(uint64(i) >> o), nil
	case UInt64:
		return SmallInt(uint64(i) >> o), nil
	case UInt32:
		return SmallInt(uint64(i) >> o), nil
	case UInt16:
		return SmallInt(uint64(i) >> o), nil
	case UInt8:
		return SmallInt(uint64(i) >> o), nil
	case *BigInt:
		if o.IsSmallInt() {
			oSmall := o.ToSmallInt()
			if oSmall < 0 {
				return leftBitshiftSmallInt(i, -oSmall), nil
			}
			return SmallInt(uint64(i) >> oSmall), nil
		}
		return SmallInt(0), nil
	default:
		return nil, NewBitshiftOperandError(other)
	}
}

func initSmallInt() {
	SmallIntClass = NewClass(
		ClassWithParent(IntClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithSingleton(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("SmallInt", SmallIntClass)
}
