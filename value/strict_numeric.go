package value

import (
	"math"
	"math/big"
	"strconv"
)

// Strict numerics are sized and can't be automatically coerced
// to other types.
type StrictNumeric interface {
	Float64 | Float32 | Int64 | Int32 | Int16 | Int8 | UInt64 | UInt32 | UInt16 | UInt8
	Value
}

// Strict integers are sized and can't be automatically coerced
// to other types.
type StrictInt interface {
	Int64 | Int32 | Int16 | Int8 | UInt64 | UInt32 | UInt16 | UInt8
	Value
}

// Strict unsigned integers are sized and can't be automatically coerced
// to other types.
type StrictUnsignedInt interface {
	UInt64 | UInt32 | UInt16 | UInt8
	Value
}

// Strict signed integers are sized and can't be automatically coerced
// to other types.
type StrictSignedInt interface {
	Int64 | Int32 | Int16 | Int8
	Value
}

// Strict floats are sized and can't be automatically coerced
// to other types.
type StrictFloat interface {
	Float64 | Float32
	Value
}

func logicalRightShift64[L SimpleInt](left L, right uint64) L {
	return L(uint64(left) >> right)
}

func logicalRightShift32[L SimpleInt](left L, right uint64) L {
	return L(uint32(left) >> right)
}

func logicalRightShift16[L SimpleInt](left L, right uint64) L {
	return L(uint16(left) >> right)
}

func logicalRightShift8[L SimpleInt](left L, right uint64) L {
	return L(uint8(left) >> right)
}

type logicalShiftFunc[L SimpleInt] func(left L, right uint64) L

// Logically bitshift a strict int to the right.
func StrictIntLogicalRightBitshift[T StrictInt](left T, right Value, shiftFunc logicalShiftFunc[T]) (T, *Error) {
	switch r := right.(type) {
	case SmallInt:
		if r < 0 {
			return left << -r, nil
		}
		return shiftFunc(left, uint64(r)), nil
	case Int64:
		if r < 0 {
			return left << -r, nil
		}
		return shiftFunc(left, uint64(r)), nil
	case Int32:
		if r < 0 {
			return left << -r, nil
		}
		return shiftFunc(left, uint64(r)), nil
	case Int16:
		if r < 0 {
			return left << -r, nil
		}
		return shiftFunc(left, uint64(r)), nil
	case Int8:
		if r < 0 {
			return left << -r, nil
		}
		return shiftFunc(left, uint64(r)), nil
	case UInt64:
		return shiftFunc(left, uint64(r)), nil
	case UInt32:
		return shiftFunc(left, uint64(r)), nil
	case UInt16:
		return shiftFunc(left, uint64(r)), nil
	case UInt8:
		return shiftFunc(left, uint64(r)), nil
	case *BigInt:
		if r.IsSmallInt() {
			rSmall := r.ToSmallInt()
			if rSmall < 0 {
				return left << -rSmall, nil
			}
			return shiftFunc(left, uint64(rSmall)), nil
		}

		return 0, nil
	default:
		return 0, NewBitshiftOperandError(right)
	}
}

// Bitshift a strict int to the right.
func StrictIntRightBitshift[T StrictInt](left T, right Value) (T, *Error) {
	switch r := right.(type) {
	case SmallInt:
		if r < 0 {
			return left << -r, nil
		}
		return left >> r, nil
	case Int64:
		if r < 0 {
			return left << -r, nil
		}
		return left >> r, nil
	case Int32:
		if r < 0 {
			return left << -r, nil
		}
		return left >> r, nil
	case Int16:
		if r < 0 {
			return left << -r, nil
		}
		return left >> r, nil
	case Int8:
		if r < 0 {
			return left << -r, nil
		}
		return left >> r, nil
	case UInt64:
		return left >> r, nil
	case UInt32:
		return left >> r, nil
	case UInt16:
		return left >> r, nil
	case UInt8:
		return left >> r, nil
	case *BigInt:
		if r.IsSmallInt() {
			rSmall := r.ToSmallInt()
			if rSmall < 0 {
				return left << -rSmall, nil
			}
			return left >> rSmall, nil
		}

		return 0, nil
	default:
		return 0, NewBitshiftOperandError(right)
	}
}

// Bitshift a strict int to the left.
func StrictIntLeftBitshift[T StrictInt](left T, right Value) (T, *Error) {
	switch r := right.(type) {
	case SmallInt:
		if r < 0 {
			return left >> -r, nil
		}
		return left << r, nil
	case Int64:
		if r < 0 {
			return left >> -r, nil
		}
		return left << r, nil
	case Int32:
		if r < 0 {
			return left >> -r, nil
		}
		return left << r, nil
	case Int16:
		if r < 0 {
			return left >> -r, nil
		}
		return left << r, nil
	case Int8:
		if r < 0 {
			return left >> -r, nil
		}
		return left << r, nil
	case UInt64:
		return left << r, nil
	case UInt32:
		return left << r, nil
	case UInt16:
		return left << r, nil
	case UInt8:
		return left << r, nil
	case *BigInt:
		if r.IsSmallInt() {
			rSmall := r.ToSmallInt()
			if rSmall < 0 {
				return left >> -rSmall, nil
			}
			return left << rSmall, nil
		}

		return 0, nil
	default:
		return 0, NewBitshiftOperandError(right)
	}
}

// Perform a bitwise AND.
func StrictIntBitwiseAnd[T StrictInt](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, NewCoerceError(left, right)
	}

	return left & r, nil
}

// Perform a bitwise OR.
func StrictIntBitwiseOr[T StrictInt](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, NewCoerceError(left, right)
	}

	return left | r, nil
}

// Perform a bitwise XOR.
func StrictIntBitwiseXor[T StrictInt](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, NewCoerceError(left, right)
	}

	return left ^ r, nil
}

// Exponentiate a strict int by the right value.
func StrictFloatExponentiate[T StrictFloat](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, NewCoerceError(left, right)
	}

	return T(math.Pow(float64(left), float64(r))), nil
}

// Exponentiate a strict int by the right value.
func StrictIntExponentiate[T StrictInt](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, NewCoerceError(left, right)
	}

	if r <= 0 {
		return 1, nil
	}
	result := left
	var i T
	for i = 2; i <= r; i++ {
		result *= left
	}
	return result, nil
}

// Add a strict numeric to another value and return the result.
// If the operation is illegal an error will be returned.
func StrictNumericAdd[T StrictNumeric](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, NewCoerceError(left, right)
	}

	return left + r, nil
}

// Subtract a strict numeric from another value and return the result.
// If the operation is illegal an error will be returned.
func StrictNumericSubtract[T StrictNumeric](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, NewCoerceError(left, right)
	}

	return left - r, nil
}

// Multiply a strict numeric by another value and return the result.
// If the operation is illegal an error will be returned.
func StrictNumericMultiply[T StrictNumeric](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, NewCoerceError(left, right)
	}

	return left * r, nil
}

// Perform modulo on a strict integer and return the result.
// If the operation is illegal an error will be returned.
func StrictIntModulo[T StrictInt](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, NewCoerceError(left, right)
	}
	if r == 0 {
		return 0, NewZeroDivisionError()
	}

	return left % r, nil
}

// Perform modulo on a strict integer and return the result.
// If the operation is illegal an error will be returned.
func StrictFloatModulo[T StrictFloat](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, NewCoerceError(left, right)
	}

	return T(math.Mod(float64(left), float64(r))), nil
}

// Divide a strict float by another value and return the result.
// If the operation is illegal an error will be returned.
func StrictFloatDivide[T StrictFloat](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, NewCoerceError(left, right)
	}

	return left / r, nil
}

// Divide a strict int by another value and return the result.
// If the operation is illegal an error will be returned.
func StrictIntDivide[T StrictInt](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, NewCoerceError(left, right)
	}
	if r == 0 {
		return 0, NewZeroDivisionError()
	}

	return left / r, nil
}

// Check whether left is greater than right.
// If the operation is illegal an error will be returned.
func StrictNumericGreaterThan[T StrictNumeric](left T, right Value) (Bool, *Error) {
	r, ok := right.(T)
	if !ok {
		return nil, NewCoerceError(left, right)
	}

	return ToElkBool(left > r), nil
}

// Check whether left is greater than or equal to right.
// If the operation is illegal an error will be returned.
func StrictNumericGreaterThanEqual[T StrictNumeric](left T, right Value) (Bool, *Error) {
	r, ok := right.(T)
	if !ok {
		return nil, NewCoerceError(left, right)
	}

	return ToElkBool(left >= r), nil
}

// Check whether left is less than right.
// If the operation is illegal an error will be returned.
func StrictNumericLessThan[T StrictNumeric](left T, right Value) (Bool, *Error) {
	r, ok := right.(T)
	if !ok {
		return nil, NewCoerceError(left, right)
	}

	return ToElkBool(left < r), nil
}

// Check whether left is less than or equal to right.
// If the operation is illegal an error will be returned.
func StrictNumericLessThanEqual[T StrictNumeric](left T, right Value) (Bool, *Error) {
	r, ok := right.(T)
	if !ok {
		return nil, NewCoerceError(left, right)
	}

	return ToElkBool(left <= r), nil
}

// Check whether the left float is equal to right.
func StrictFloatEqual[T StrictFloat](left T, right Value) Bool {
	switch r := right.(type) {
	case T:
		return ToElkBool(left == r)
	case SmallInt:
		return ToElkBool(float64(left) == float64(r))
	case *BigInt:
		return ToElkBool(float64(left) == float64(r.ToFloat()))
	case Float:
		return ToElkBool(float64(left) == float64(r))
	case *BigFloat:
		if r.IsNaN() {
			return False
		}
		iBigFloat := (&big.Float{}).SetFloat64(float64(left))
		return ToElkBool(iBigFloat.Cmp(r.AsGoBigFloat()) == 0)
	case Int64:
		return ToElkBool(float64(left) == float64(r))
	case Int32:
		return ToElkBool(float64(left) == float64(r))
	case Int16:
		return ToElkBool(float64(left) == float64(r))
	case Int8:
		return ToElkBool(float64(left) == float64(r))
	case UInt64:
		return ToElkBool(float64(left) == float64(r))
	case UInt32:
		return ToElkBool(float64(left) == float64(r))
	case UInt16:
		return ToElkBool(float64(left) == float64(r))
	case UInt8:
		return ToElkBool(float64(left) == float64(r))
	case Float64:
		return ToElkBool(float64(left) == float64(r))
	case Float32:
		return ToElkBool(float64(left) == float64(r))
	default:
		return False
	}
}

// Check whether the left signed integer is equal to right.
func StrictSignedIntEqual[T StrictSignedInt](left T, right Value) Bool {
	switch r := right.(type) {
	case T:
		return ToElkBool(left == r)
	case SmallInt:
		return ToElkBool(int64(left) == int64(r))
	case *BigInt:
		iBigInt := big.NewInt(int64(left))
		return ToElkBool(iBigInt.Cmp(r.ToGoBigInt()) == 0)
	case Float:
		return ToElkBool(float64(left) == float64(r))
	case *BigFloat:
		if r.IsNaN() {
			return False
		}
		iBigFloat := (&big.Float{}).SetInt64(int64(left))
		return ToElkBool(iBigFloat.Cmp(r.AsGoBigFloat()) == 0)
	case Int64:
		return ToElkBool(int64(left) == int64(r))
	case Int32:
		return ToElkBool(int64(left) == int64(r))
	case Int16:
		return ToElkBool(int64(left) == int64(r))
	case Int8:
		return ToElkBool(left == T(r))
	case UInt64:
		if r > math.MaxInt64 {
			return False
		}
		return ToElkBool(int64(left) == int64(r))
	case UInt32:
		return ToElkBool(int64(left) == int64(r))
	case UInt16:
		return ToElkBool(int64(left) == int64(r))
	case UInt8:
		return ToElkBool(int64(left) == int64(r))
	case Float64:
		return ToElkBool(float64(left) == float64(r))
	case Float32:
		return ToElkBool(float64(left) == float64(r))
	default:
		return False
	}
}

// Check whether the left unsigned integer is equal to right.
func StrictUnsignedIntEqual[T StrictUnsignedInt](left T, right Value) Bool {
	switch r := right.(type) {
	case T:
		return ToElkBool(left == r)
	case SmallInt:
		if uint64(left) > math.MaxInt64 {
			return False
		}
		return ToElkBool(int64(left) == int64(r))
	case *BigInt:
		iBigInt := (&big.Int{}).SetUint64(uint64(left))
		return ToElkBool(iBigInt.Cmp(r.ToGoBigInt()) == 0)
	case Float:
		return ToElkBool(float64(left) == float64(r))
	case *BigFloat:
		if r.IsNaN() {
			return False
		}
		iBigFloat := (&big.Float{}).SetUint64(uint64(left))
		return ToElkBool(iBigFloat.Cmp(r.AsGoBigFloat()) == 0)
	case Int64:
		if uint64(left) > math.MaxInt64 {
			return False
		}
		return ToElkBool(int64(left) == int64(r))
	case Int32:
		if uint64(left) > math.MaxInt64 {
			return False
		}
		return ToElkBool(int64(left) == int64(r))
	case Int16:
		if uint64(left) > math.MaxInt64 {
			return False
		}
		return ToElkBool(int64(left) == int64(r))
	case Int8:
		if uint64(left) > math.MaxInt64 {
			return False
		}
		return ToElkBool(int64(left) == int64(r))
	case UInt64:
		return ToElkBool(uint64(left) == uint64(r))
	case UInt32:
		return ToElkBool(uint64(left) == uint64(r))
	case UInt16:
		return ToElkBool(uint64(left) == uint64(r))
	case UInt8:
		return ToElkBool(left == T(r))
	case Float64:
		return ToElkBool(float64(left) == float64(r))
	case Float32:
		return ToElkBool(float64(left) == float64(r))
	default:
		return False
	}
}

// Check whether left is strictly equal to right.
func StrictNumericStrictEqual[T StrictNumeric](left T, right Value) Bool {
	r, ok := right.(T)
	if !ok {
		return False
	}

	return ToElkBool(left == r)
}

// Parses an unsigned strict integer from a string using Elk syntax.
func StrictParseUint(s string, base int, bitSize int) (uint64, *Error) {
	if s == "" {
		return 0, Errorf(FormatErrorClass, "invalid integer format")
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
		return 0, Errorf(FormatErrorClass, "invalid integer base %d", base)
	}

	if bitSize == 0 {
		bitSize = strconv.IntSize
	} else if bitSize < 0 || bitSize > 64 {
		return 0, Errorf(FormatErrorClass, "invalid integer bit size %d", bitSize)
	}

	// Cutoff is the smallest number such that cutoff*base > math.MaxUint64.
	// Use compile-time constants for common cases.
	var cutoff uint64
	switch base {
	case 2:
		cutoff = math.MaxUint64/2 + 1
	case 8:
		cutoff = math.MaxUint64/8 + 1
	case 10:
		cutoff = math.MaxUint64/10 + 1
	case 16:
		cutoff = math.MaxUint64/16 + 1
	default:
		cutoff = math.MaxUint64/uint64(base) + 1
	}

	maxVal := uint64(1)<<uint(bitSize) - 1

	var n uint64
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
			return 0, Errorf(FormatErrorClass, "illegal characters in integer: %c", c)
		}

		if d >= byte(base) {
			return 0, Errorf(FormatErrorClass, "illegal characters in integer (base %d): %c", base, c)
		}

		if n >= cutoff {
			// n*base overflows
			return maxVal, Errorf(FormatErrorClass, "value overflows")
		}
		n *= uint64(base)

		n1 := n + uint64(d)
		if n1 < n || n1 > maxVal {
			// n+d overflows
			return maxVal, Errorf(FormatErrorClass, "value overflows")
		}
		n = n1
	}

	return n, nil
}

// Parses a signed strict integer from a string using Elk syntax.
func StrictParseInt(s string, base int, bitSize int) (int64, *Error) {
	if s == "" {
		return 0, Errorf(FormatErrorClass, "invalid integer format")
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
	var un uint64
	un, err := StrictParseUint(s, base, bitSize)
	if err != nil {
		return 0, err
	}

	if bitSize == 0 {
		bitSize = strconv.IntSize
	}

	cutoff := uint64(1 << uint(bitSize-1))
	if !neg && un >= cutoff {
		return int64(cutoff - 1), Errorf(FormatErrorClass, "value overflows")
	}
	if neg && un > cutoff {
		return -int64(cutoff), Errorf(FormatErrorClass, "value overflows")
	}
	n := int64(un)
	if neg {
		n = -n
	}
	return n, nil
}

// Converts letters to lowercase.
func letterToLower(c byte) byte {
	return c | ('x' - 'X')
}
