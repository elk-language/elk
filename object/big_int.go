package object

import (
	"math"
	"math/big"

	"github.com/ALTree/bigfloat"
)

var BigIntClass *Class // ::Std::BigInt

// Elk's BigInt value
type BigInt big.Int

// Create a new BigInt with teh specified value.
func NewBigInt(i int64) *BigInt {
	return ToElkBigInt(big.NewInt(i))
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

// Returns the SmallInt representation of i.
func (i *BigInt) ToSmallInt() SmallInt {
	return SmallInt(i.ToGoBigInt().Int64())
}

// Negate the number and return the result.
func (i *BigInt) Negate() Value {
	return ToElkBigInt((&big.Int{}).Neg(i.ToGoBigInt()))
}

// Add another value and return an error
// if something went wrong.
func (i *BigInt) Add(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := big.NewInt(int64(o))
		oBigInt.Add(i.ToGoBigInt(), oBigInt)
		if oBigInt.IsInt64() {
			return SmallInt(oBigInt.Int64()), nil
		}
		return ToElkBigInt(oBigInt), nil
	case *BigInt:
		result := (&big.Int{}).Add(i.ToGoBigInt(), o.ToGoBigInt())
		if result.IsInt64() {
			return SmallInt(result.Int64()), nil
		}
		return ToElkBigInt(result), nil
	case Float:
		iBigFloat := (&big.Float{}).SetInt(i.ToGoBigInt())
		oBigFloat := big.NewFloat(float64(o))
		result, _ := iBigFloat.Add(iBigFloat, oBigFloat).Float64()
		return Float(result), nil
	case *BigFloat:
		iGo := i.ToGoBigInt()
		prec := max(o.Precision(), uint(iGo.BitLen()), 64)
		iBigFloat := (&big.Float{}).SetPrec(prec).SetInt(iGo)
		iBigFloat.Add(iBigFloat, o.ToGoBigFloat())
		return ToElkBigFloat(iBigFloat), nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

// Subtract another value and return an error
// if something went wrong.
func (i *BigInt) Subtract(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := big.NewInt(int64(o))
		oBigInt.Sub(i.ToGoBigInt(), oBigInt)
		if oBigInt.IsInt64() {
			return SmallInt(oBigInt.Int64()), nil
		}
		return ToElkBigInt(oBigInt), nil
	case *BigInt:
		result := (&big.Int{}).Sub(i.ToGoBigInt(), o.ToGoBigInt())
		if result.IsInt64() {
			return SmallInt(result.Int64()), nil
		}
		return ToElkBigInt(result), nil
	case Float:
		iBigFloat := (&big.Float{}).SetInt(i.ToGoBigInt())
		oBigFloat := big.NewFloat(float64(o))
		result, _ := iBigFloat.Sub(iBigFloat, oBigFloat).Float64()
		return Float(result), nil
	case *BigFloat:
		iGo := i.ToGoBigInt()
		prec := max(o.Precision(), uint(iGo.BitLen()), 64)
		iBigFloat := (&big.Float{}).SetPrec(prec).SetInt(iGo)
		iBigFloat.Sub(iBigFloat, o.ToGoBigFloat())
		return ToElkBigFloat(iBigFloat), nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

// Multiply by another value and return an error
// if something went wrong.
func (i *BigInt) Multiply(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := big.NewInt(int64(o))
		oBigInt.Mul(i.ToGoBigInt(), oBigInt)
		if oBigInt.IsInt64() {
			return SmallInt(oBigInt.Int64()), nil
		}
		return ToElkBigInt(oBigInt), nil
	case *BigInt:
		result := (&big.Int{}).Mul(i.ToGoBigInt(), o.ToGoBigInt())
		if result.IsInt64() {
			return SmallInt(result.Int64()), nil
		}
		return ToElkBigInt(result), nil
	case Float:
		iBigFloat := (&big.Float{}).SetInt(i.ToGoBigInt())
		oBigFloat := big.NewFloat(float64(o))
		result, _ := iBigFloat.Mul(iBigFloat, oBigFloat).Float64()
		return Float(result), nil
	case *BigFloat:
		iGo := i.ToGoBigInt()
		prec := max(o.Precision(), uint(iGo.BitLen()), 64)
		iBigFloat := (&big.Float{}).SetPrec(prec).SetInt(iGo)
		iBigFloat.Mul(iBigFloat, o.ToGoBigFloat())
		return ToElkBigFloat(iBigFloat), nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

// Divide by another value and return an error
// if something went wrong.
func (i *BigInt) Divide(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		if o == 0 {
			return nil, NewZeroDivisionError()
		}
		oBigInt := big.NewInt(int64(o))
		oBigInt.Div(i.ToGoBigInt(), oBigInt)
		if oBigInt.IsInt64() {
			return SmallInt(oBigInt.Int64()), nil
		}
		return ToElkBigInt(oBigInt), nil
	case *BigInt:
		if len(o.ToGoBigInt().Bits()) == 0 {
			return nil, NewZeroDivisionError()
		}
		result := (&big.Int{}).Div(i.ToGoBigInt(), o.ToGoBigInt())
		if result.IsInt64() {
			return SmallInt(result.Int64()), nil
		}
		return ToElkBigInt(result), nil
	case Float:
		iBigFloat := (&big.Float{}).SetInt(i.ToGoBigInt())
		oBigFloat := big.NewFloat(float64(o))
		result, _ := iBigFloat.Quo(iBigFloat, oBigFloat).Float64()
		return Float(result), nil
	case *BigFloat:
		iGo := i.ToGoBigInt()
		prec := max(o.Precision(), uint(iGo.BitLen()), 64)
		iBigFloat := (&big.Float{}).SetPrec(prec).SetInt(iGo)
		iBigFloat.Quo(iBigFloat, o.ToGoBigFloat())
		return ToElkBigFloat(iBigFloat), nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

// Exponentiate by another value and return an error
// if something went wrong.
func (i *BigInt) Exponentiate(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := big.NewInt(int64(o))
		oBigInt.Exp(i.ToGoBigInt(), oBigInt, nil)
		if oBigInt.IsInt64() {
			return SmallInt(oBigInt.Int64()), nil
		}
		return ToElkBigInt(oBigInt), nil
	case *BigInt:
		result := (&big.Int{}).Exp(i.ToGoBigInt(), o.ToGoBigInt(), nil)
		if result.IsInt64() {
			return SmallInt(result.Int64()), nil
		}
		return ToElkBigInt(result), nil
	case Float:
		iFloat, _ := i.ToGoBigInt().Float64()
		return Float(math.Pow(iFloat, float64(o))), nil
	case *BigFloat:
		iGo := i.ToGoBigInt()
		prec := max(o.Precision(), uint(iGo.BitLen()), 64)
		iBigFloat := (&big.Float{}).SetPrec(prec).SetInt(iGo)
		result := bigfloat.Pow(iBigFloat, o.ToGoBigFloat())
		return ToElkBigFloat(result), nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

func (i *BigInt) Class() *Class {
	return BigIntClass
}

func (i *BigInt) IsFrozen() bool {
	return true
}

func (i *BigInt) SetFrozen() {}

func (i *BigInt) Inspect() string {
	return i.ToGoBigInt().String()
}

func (i *BigInt) InstanceVariables() SimpleSymbolMap {
	return nil
}

// Parses an unsigned big.Int from a string using Elk syntax.
func ParseUBigInt(s string, base int) (*BigInt, *Error) {
	if s == "" {
		return nil, Errorf(FormatErrorClass, "invalid integer format")
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
		return nil, Errorf(FormatErrorClass, "invalid integer base %d", base)
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
			return nil, Errorf(FormatErrorClass, "illegal characters in integer: %c", c)
		}

		if d >= byte(base) {
			return nil, Errorf(FormatErrorClass, "illegal characters in integer (base %d): %c", base, c)
		}

		n.Mul(n, big.NewInt(int64(base)))

		n.Add(n, big.NewInt(int64(d)))
	}

	return ToElkBigInt(n), nil
}

// Parses a signed big.Int from a string using Elk syntax.
func ParseBigInt(s string, base int) (*BigInt, *Error) {
	if s == "" {
		return nil, Errorf(FormatErrorClass, "invalid integer format")
	}

	// Pick off leading sign.
	neg := false
	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		neg = true
		s = s[1:]
	}

	// Convert unsigned and check range.
	u, err := ParseUBigInt(s, base)
	un := u.ToGoBigInt()

	if err != nil {
		return nil, err
	}

	if neg {
		un.Neg(un)
	}

	return ToElkBigInt(un), nil
}

// Same as [ParseBigInt] but panics on error.
func ParseBigIntPanic(s string, base int) *BigInt {
	result, err := ParseBigInt(s, base)
	if err != nil {
		panic(err)
	}

	return result
}

func initBigInt() {
	BigIntClass = NewClass(
		ClassWithParent(IntClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithSingleton(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("BigInt", BigIntClass)
}
