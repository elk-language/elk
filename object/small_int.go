package object

import (
	"fmt"
	"math"
	"math/big"
)

var SmallIntClass *Class // ::Std::SmallInt

// Elk's SmallInt value
type SmallInt int64

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
