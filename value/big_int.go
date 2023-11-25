package value

import (
	"math"
	"math/big"
)

// Elk's BigInt value
type BigInt big.Int

// Create a new BigInt with the specified value.
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

// Reports whether i is zero.
func (i *BigInt) IsZero() bool {
	return len(i.ToGoBigInt().Bits()) == 0
}

// Returns the SmallInt representation of i.
func (i *BigInt) ToSmallInt() SmallInt {
	return SmallInt(i.ToGoBigInt().Int64())
}

// Convert to Elk Int64
func (i *BigInt) ToInt64() Int64 {
	return i.ToSmallInt().ToInt64()
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
func (i *BigInt) Negate() Value {
	return ToElkBigInt((&big.Int{}).Neg(i.ToGoBigInt()))
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

// Add another value and return an error
// if something went wrong.
func (i *BigInt) Add(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := big.NewInt(int64(o))
		oBigInt.Add(i.ToGoBigInt(), oBigInt)
		result := ToElkBigInt(oBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return ToElkBigInt(oBigInt), nil
	case *BigInt:
		result := ToElkBigInt((&big.Int{}).Add(i.ToGoBigInt(), o.ToGoBigInt()))
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	case Float:
		return i.ToFloat() + o, nil
	case *BigFloat:
		prec := max(o.Precision(), uint(i.BitSize()), 64)
		iBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(i)
		iBigFloat.AddBigFloat(iBigFloat, o)
		return iBigFloat, nil
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
		result := ToElkBigInt(oBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return ToElkBigInt(oBigInt), nil
	case *BigInt:
		result := ToElkBigInt((&big.Int{}).Sub(i.ToGoBigInt(), o.ToGoBigInt()))
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	case Float:
		return i.ToFloat() - o, nil
	case *BigFloat:
		prec := max(o.Precision(), uint(i.BitSize()), 64)
		iBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(i)
		iBigFloat.SubBigFloat(iBigFloat, o)
		return iBigFloat, nil
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
		result := ToElkBigInt(oBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	case *BigInt:
		result := ToElkBigInt((&big.Int{}).Mul(i.ToGoBigInt(), o.ToGoBigInt()))
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	case Float:
		return i.ToFloat() * o, nil
	case *BigFloat:
		prec := max(o.Precision(), uint(i.BitSize()), 64)
		iBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(i)
		return iBigFloat.MulBigFloat(iBigFloat, o), nil
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
		result := ToElkBigInt(oBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return ToElkBigInt(oBigInt), nil
	case *BigInt:
		if o.IsZero() {
			return nil, NewZeroDivisionError()
		}
		result := ToElkBigInt((&big.Int{}).Div(i.ToGoBigInt(), o.ToGoBigInt()))
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	case Float:
		return i.ToFloat() / o, nil
	case *BigFloat:
		prec := max(o.Precision(), uint(i.BitSize()), 64)
		iBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(i)
		return iBigFloat.DivBigFloat(iBigFloat, o), nil
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
		result := ToElkBigInt(oBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return ToElkBigInt(oBigInt), nil
	case *BigInt:
		result := ToElkBigInt((&big.Int{}).Exp(i.ToGoBigInt(), o.ToGoBigInt(), nil))
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	case Float:
		iFloat, _ := i.ToGoBigInt().Float64()
		return Float(math.Pow(iFloat, float64(o))), nil
	case *BigFloat:
		prec := max(o.Precision(), uint(i.BitSize()), 64)
		iBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(i)
		iBigFloat.ExpBigFloat(iBigFloat, o)
		return iBigFloat, nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

// Perform modulo with another numeric value and return an error
// if something went wrong.
func (i *BigInt) Modulo(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		if o == 0 {
			return nil, NewZeroDivisionError()
		}
		iGo := i.ToGoBigInt()
		oBigInt := big.NewInt(int64(o))
		(&big.Int{}).QuoRem(iGo, oBigInt, oBigInt)
		result := ToElkBigInt(oBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	case *BigInt:
		if o.IsZero() {
			return nil, NewZeroDivisionError()
		}
		iGo := i.ToGoBigInt()
		oGo := o.ToGoBigInt()
		mod := &big.Int{}
		(&big.Int{}).QuoRem(iGo, oGo, mod)
		result := ToElkBigInt(mod)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	case Float:
		iFloat, _ := i.ToGoBigInt().Float64()
		return Float(math.Mod(iFloat, float64(o))), nil
	case *BigFloat:
		prec := max(o.Precision(), uint(i.BitSize()), 64)
		iBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(i)
		return iBigFloat.Mod(iBigFloat, o), nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

// Check whether i is greater than other and return an error
// if something went wrong.
func (i *BigInt) GreaterThan(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) == 1), nil
	case *BigInt:
		return ToElkBool(i.Cmp(o) == 1), nil
	case Float:
		return ToElkBool(i.ToFloat() > o), nil
	case *BigFloat:
		if o.IsNaN() {
			return False, nil
		}
		iBigFloat := (&BigFloat{}).SetBigInt(i)
		return ToElkBool(iBigFloat.Cmp(o) == 1), nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

// Check whether i is greater than or equal to other and return an error
// if something went wrong.
func (i *BigInt) GreaterThanEqual(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) >= 0), nil
	case *BigInt:
		return ToElkBool(i.Cmp(o) >= 0), nil
	case Float:
		return ToElkBool(i.ToFloat() >= o), nil
	case *BigFloat:
		if o.IsNaN() {
			return False, nil
		}
		iBigFloat := (&BigFloat{}).SetBigInt(i)
		return ToElkBool(iBigFloat.Cmp(o) >= 0), nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

// Check whether i is less than other and return an error
// if something went wrong.
func (i *BigInt) LessThan(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) == -1), nil
	case *BigInt:
		return ToElkBool(i.Cmp(o) == -1), nil
	case Float:
		return ToElkBool(i.ToFloat() < o), nil
	case *BigFloat:
		if o.IsNaN() {
			return False, nil
		}
		iBigFloat := (&BigFloat{}).SetBigInt(i)
		return ToElkBool(iBigFloat.Cmp(o) == -1), nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

// Check whether i is less than or equal to other and return an error
// if something went wrong.
func (i *BigInt) LessThanEqual(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) <= 0), nil
	case *BigInt:
		return ToElkBool(i.Cmp(o) <= 0), nil
	case Float:
		return ToElkBool(i.ToFloat() <= o), nil
	case *BigFloat:
		if o.IsNaN() {
			return False, nil
		}
		iBigFloat := (&BigFloat{}).SetBigInt(i)
		return ToElkBool(iBigFloat.Cmp(o) <= 0), nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

// Check whether i is equal to other
func (i *BigInt) Equal(other Value) Value {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case *BigInt:
		return ToElkBool(i.Cmp(o) == 0)
	case Float:
		return ToElkBool(i.ToFloat() == o)
	case *BigFloat:
		if o.IsNaN() {
			return False
		}
		iBigFloat := (&BigFloat{}).SetBigInt(i)
		return ToElkBool(iBigFloat.Cmp(o) == 0)
	case Int64:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case Int32:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case Int16:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case Int8:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case UInt64:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case UInt32:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case UInt16:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case UInt8:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case Float64:
		return ToElkBool(i.ToFloat() == Float(o))
	case Float32:
		return ToElkBool(i.ToFloat() == Float(o))
	default:
		return False
	}
}

// Check whether i is strictly equal to other
func (i *BigInt) StrictEqual(other Value) Value {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := NewBigInt(int64(o))
		return ToElkBool(i.Cmp(oBigInt) == 0)
	case *BigInt:
		return ToElkBool(i.Cmp(o) == 0)
	default:
		return False
	}
}

func rightBitshiftBigInt[T SimpleInt](i *BigInt, other T) Value {
	if other < 0 {
		return SmallInt(0)
	}
	iGo := i.ToGoBigInt()
	result := ToElkBigInt(iGo.Rsh(iGo, uint(other)))
	if result.IsSmallInt() {
		return result.ToSmallInt()
	}
	return result
}

// Bitshift to the right by another integer value and return an error
// if something went wrong.
func (i *BigInt) RightBitshift(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		if o < 0 {
			return leftBitshiftBigInt(i, -o), nil
		}
		return rightBitshiftBigInt(i, o), nil
	case Int64:
		if o < 0 {
			return leftBitshiftBigInt(i, -o), nil
		}
		return rightBitshiftBigInt(i, o), nil
	case Int32:
		if o < 0 {
			return leftBitshiftBigInt(i, -o), nil
		}
		return rightBitshiftBigInt(i, o), nil
	case Int16:
		if o < 0 {
			return leftBitshiftBigInt(i, -o), nil
		}
		return rightBitshiftBigInt(i, o), nil
	case Int8:
		if o < 0 {
			return leftBitshiftBigInt(i, -o), nil
		}
		return rightBitshiftBigInt(i, o), nil
	case UInt64:
		return rightBitshiftBigInt(i, o), nil
	case UInt32:
		return rightBitshiftBigInt(i, o), nil
	case UInt16:
		return rightBitshiftBigInt(i, o), nil
	case UInt8:
		return rightBitshiftBigInt(i, o), nil
	case *BigInt:
		if o.IsSmallInt() {
			oSmall := o.ToSmallInt()
			if oSmall < 0 {
				return leftBitshiftBigInt(i, -oSmall), nil
			}
			return rightBitshiftBigInt(i, oSmall), nil
		}
		return SmallInt(0), nil
	default:
		return nil, NewBitshiftOperandError(other)
	}
}

func leftBitshiftBigInt[T SimpleInt](i *BigInt, other T) Value {
	if other < 0 {
		return SmallInt(0)
	}
	iGo := i.ToGoBigInt()
	return ToElkBigInt(iGo.Lsh(iGo, uint(other)))
}

// Bitshift to the left by another integer value and return an error
// if something went wrong.
func (i *BigInt) LeftBitshift(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return leftBitshiftBigInt(i, o), nil
	case Int64:
		return leftBitshiftBigInt(i, o), nil
	case Int32:
		return leftBitshiftBigInt(i, o), nil
	case Int16:
		return leftBitshiftBigInt(i, o), nil
	case Int8:
		return leftBitshiftBigInt(i, o), nil
	case UInt64:
		return leftBitshiftBigInt(i, o), nil
	case UInt32:
		return leftBitshiftBigInt(i, o), nil
	case UInt16:
		return leftBitshiftBigInt(i, o), nil
	case UInt8:
		return leftBitshiftBigInt(i, o), nil
	case *BigInt:
		if o.IsSmallInt() {
			oSmall := o.ToSmallInt()
			return leftBitshiftBigInt(i, oSmall), nil
		}
		return SmallInt(0), nil
	default:
		return nil, NewBitshiftOperandError(other)
	}
}

// Perform bitwise AND with another value and return an error
// if something went wrong.
func (i *BigInt) BitwiseAnd(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := big.NewInt(int64(o))
		oBigInt.And(i.ToGoBigInt(), oBigInt)
		result := ToElkBigInt(oBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	case *BigInt:
		result := ToElkBigInt((&big.Int{}).And(i.ToGoBigInt(), o.ToGoBigInt()))
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

// Perform bitwise OR with another value and return an error
// if something went wrong.
func (i *BigInt) BitwiseOr(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := big.NewInt(int64(o))
		oBigInt.Or(i.ToGoBigInt(), oBigInt)
		result := ToElkBigInt(oBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	case *BigInt:
		result := ToElkBigInt((&big.Int{}).Or(i.ToGoBigInt(), o.ToGoBigInt()))
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	default:
		return nil, NewCoerceError(i, other)
	}
}

// Perform bitwise XOR with another value and return an error
// if something went wrong.
func (i *BigInt) BitwiseXor(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		oBigInt := big.NewInt(int64(o))
		oBigInt.Xor(i.ToGoBigInt(), oBigInt)
		result := ToElkBigInt(oBigInt)
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	case *BigInt:
		result := ToElkBigInt((&big.Int{}).Xor(i.ToGoBigInt(), o.ToGoBigInt()))
		if result.IsSmallInt() {
			return result.ToSmallInt(), nil
		}
		return result, nil
	default:
		return nil, NewCoerceError(i, other)
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

func (i *BigInt) InstanceVariables() SymbolMap {
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
