package object

import (
	"fmt"
	"math"
	"math/big"
)

var IntClass *Class // ::Std::Int

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
		return ToElkBigInt(iBigInt.Add(iBigInt, o.ToGoBigInt())), nil
	case Float:
		return Float(i) + o, nil
	case *BigFloat:
		iBigFloat := (&big.Float{}).SetInt64(int64(i))
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
		return ToElkBigInt(iBigInt.Sub(iBigInt, o.ToGoBigInt())), nil
	case Float:
		return Float(i) - o, nil
	case *BigFloat:
		iBigFloat := (&big.Float{}).SetInt64(int64(i))
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
			return ToElkBigInt(iBigInt.Sub(iBigInt, big.NewInt(int64(o)))), nil
		}
		return result, nil
	case *BigInt:
		iBigInt := big.NewInt(int64(i))
		return ToElkBigInt(iBigInt.Mul(iBigInt, o.ToGoBigInt())), nil
	case Float:
		return Float(i) * o, nil
	case *BigFloat:
		iBigFloat := (&big.Float{}).SetInt64(int64(i))
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
		return ToElkBigInt(iBigInt.Div(iBigInt, o.ToGoBigInt())), nil
	case Float:
		return Float(i) / o, nil
	case *BigFloat:
		iBigFloat := (&big.Float{}).SetInt64(int64(i))
		iBigFloat.Quo(iBigFloat, o.ToGoBigFloat())
		return ToElkBigFloat(iBigFloat), nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

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
	result := (&big.Int{}).Neg(i.ToGoBigInt())
	if result.IsInt64() {
		return SmallInt(result.Int64())
	}
	return ToElkBigInt(result)
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

func initInt() {
	IntClass = NewClass(
		ClassWithParent(NumericClass),
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	StdModule.AddConstant("Int", IntClass)

	SmallIntClass = NewClass(
		ClassWithParent(IntClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithSingleton(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("SmallInt", SmallIntClass)

	BigIntClass = NewClass(
		ClassWithParent(IntClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithSingleton(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("BigInt", BigIntClass)
}
