package value

import (
	"math"
	"math/big"
	"strconv"
)

// Strict numerics are sized and cannot be automatically coerced
// to other types.
type StrictNumeric interface {
	Float64 | Float32 | Int64 | Int32 | Int16 | Int8 | UInt64 | UInt32 | UInt16 | UInt8
}

// Strict integers are sized and cannot be automatically coerced
// to other types.
type StrictInt interface {
	Int64 | Int32 | Int16 | Int8 | UInt64 | UInt32 | UInt16 | UInt8
}

// Strict unsigned integers are sized and cannot be automatically coerced
// to other types.
type StrictUnsignedInt interface {
	UInt64 | UInt32 | UInt16 | UInt8
}

// Strict signed integers are sized and cannot be automatically coerced
// to other types.
type StrictSignedInt interface {
	Int64 | Int32 | Int16 | Int8
}

// Strict floats are sized and cannot be automatically coerced
// to other types.
type StrictFloat interface {
	Float64 | Float32
}

func LogicalRightShift64[L SimpleInt](left L, right uint64) L {
	return L(uint64(left) >> right)
}

func LogicalRightShift32[L SimpleInt](left L, right uint64) L {
	return L(uint32(left) >> right)
}

func LogicalRightShift16[L SimpleInt](left L, right uint64) L {
	return L(uint16(left) >> right)
}

func LogicalRightShift8[L SimpleInt](left L, right uint64) L {
	return L(uint8(left) >> right)
}

type logicalShiftFunc[L SimpleInt] func(left L, right uint64) L

// Bitshift a strict int to the left.
func StrictIntLogicalLeftBitshift[T StrictInt](left T, right Value, shiftFunc logicalShiftFunc[T]) (T, Value) {
	if right.IsReference() {
		switch r := right.AsReference().(type) {
		case Int64:
			if r < 0 {
				return shiftFunc(left, uint64(-r)), Undefined
			}
			return left << r, Undefined
		case UInt64:
			return left << r, Undefined
		case *BigInt:
			if r.IsSmallInt() {
				rSmall := r.ToSmallInt()
				if rSmall < 0 {
					return left >> -rSmall, Undefined
				}
				return left << rSmall, Undefined
			}

			return 0, Undefined
		default:
			return 0, Ref(NewBitshiftOperandError(right))
		}
	}

	switch right.ValueFlag() {
	case SMALL_INT_FLAG:
		r := right.AsSmallInt()
		if r < 0 {
			return shiftFunc(left, uint64(-r)), Undefined
		}
		return left << r, Undefined
	case INT64_FLAG:
		r := right.AsInlineInt64()
		if r < 0 {
			return shiftFunc(left, uint64(-r)), Undefined
		}
		return left << r, Undefined
	case INT32_FLAG:
		r := right.AsInt32()
		if r < 0 {
			return shiftFunc(left, uint64(-r)), Undefined
		}
		return left << r, Undefined
	case INT16_FLAG:
		r := right.AsInt16()
		if r < 0 {
			return shiftFunc(left, uint64(-r)), Undefined
		}
		return left << r, Undefined
	case INT8_FLAG:
		r := right.AsInt8()
		if r < 0 {
			return shiftFunc(left, uint64(-r)), Undefined
		}
		return left << r, Undefined
	case UINT64_FLAG:
		r := right.AsInlineUInt64()
		return left << r, Undefined
	case UINT32_FLAG:
		r := right.AsUInt32()
		return left << r, Undefined
	case UINT16_FLAG:
		r := right.AsUInt16()
		return left << r, Undefined
	case UINT8_FLAG:
		r := right.AsUInt8()
		return left << r, Undefined
	default:
		return 0, Ref(NewBitshiftOperandError(right))
	}
}

// Logically bitshift a strict int to the right.
func StrictIntLogicalRightBitshift[T StrictInt](left T, right Value, shiftFunc logicalShiftFunc[T]) (T, Value) {
	if right.IsReference() {
		switch r := right.AsReference().(type) {
		case Int64:
			if r < 0 {
				return left << -r, Undefined
			}
			return shiftFunc(left, uint64(r)), Undefined
		case UInt64:
			return shiftFunc(left, uint64(r)), Undefined
		case *BigInt:
			if r.IsSmallInt() {
				rSmall := r.ToSmallInt()
				if rSmall < 0 {
					return left << -rSmall, Undefined
				}
				return shiftFunc(left, uint64(rSmall)), Undefined
			}

			return 0, Undefined
		default:
			return 0, Ref(NewBitshiftOperandError(right))
		}
	}

	switch right.ValueFlag() {
	case SMALL_INT_FLAG:
		r := right.AsSmallInt()
		if r < 0 {
			return left << -r, Undefined
		}
		return shiftFunc(left, uint64(r)), Undefined
	case INT64_FLAG:
		r := right.AsInlineInt64()
		if r < 0 {
			return left << -r, Undefined
		}
		return shiftFunc(left, uint64(r)), Undefined
	case INT32_FLAG:
		r := right.AsInt32()
		if r < 0 {
			return left << -r, Undefined
		}
		return shiftFunc(left, uint64(r)), Undefined
	case INT16_FLAG:
		r := right.AsInt16()
		if r < 0 {
			return left << -r, Undefined
		}
		return shiftFunc(left, uint64(r)), Undefined
	case INT8_FLAG:
		r := right.AsInt8()
		if r < 0 {
			return left << -r, Undefined
		}
		return shiftFunc(left, uint64(r)), Undefined
	case UINT64_FLAG:
		r := right.AsInlineUInt64()
		return shiftFunc(left, uint64(r)), Undefined
	case UINT32_FLAG:
		r := right.AsUInt32()
		return shiftFunc(left, uint64(r)), Undefined
	case UINT16_FLAG:
		r := right.AsUInt16()
		return shiftFunc(left, uint64(r)), Undefined
	case UINT8_FLAG:
		r := right.AsUInt8()
		return shiftFunc(left, uint64(r)), Undefined
	default:
		return 0, Ref(NewBitshiftOperandError(right))
	}
}

// Bitshift a strict int to the right.
func StrictIntRightBitshift[T StrictInt](left T, right Value) (T, Value) {
	if right.IsReference() {
		switch r := right.AsReference().(type) {
		case Int64:
			if r < 0 {
				return left << -r, Undefined
			}
			return left >> r, Undefined
		case UInt64:
			return left >> r, Undefined
		case *BigInt:
			if r.IsSmallInt() {
				rSmall := r.ToSmallInt()
				if rSmall < 0 {
					return left << -rSmall, Undefined
				}
				return left >> rSmall, Undefined
			}

			return 0, Undefined
		default:
			return 0, Ref(NewBitshiftOperandError(right))
		}
	}

	switch right.ValueFlag() {
	case SMALL_INT_FLAG:
		r := right.AsSmallInt()
		if r < 0 {
			return left << -r, Undefined
		}
		return left >> r, Undefined
	case INT64_FLAG:
		r := right.AsInlineInt64()
		if r < 0 {
			return left << -r, Undefined
		}
		return left >> r, Undefined
	case INT32_FLAG:
		r := right.AsInt32()
		if r < 0 {
			return left << -r, Undefined
		}
		return left >> r, Undefined
	case INT16_FLAG:
		r := right.AsInt16()
		if r < 0 {
			return left << -r, Undefined
		}
		return left >> r, Undefined
	case INT8_FLAG:
		r := right.AsInt8()
		if r < 0 {
			return left << -r, Undefined
		}
		return left >> r, Undefined
	case UINT64_FLAG:
		r := right.AsInlineUInt64()
		return left >> r, Undefined
	case UINT32_FLAG:
		r := right.AsUInt32()
		return left >> r, Undefined
	case UINT16_FLAG:
		r := right.AsUInt16()
		return left >> r, Undefined
	case UINT8_FLAG:
		r := right.AsUInt8()
		return left >> r, Undefined
	default:
		return 0, Ref(NewBitshiftOperandError(right))
	}
}

// Bitshift a strict int to the left.
func StrictIntLeftBitshift[T StrictInt](left T, right Value) (T, Value) {
	if right.IsReference() {
		switch r := right.AsReference().(type) {
		case Int64:
			if r < 0 {
				return left >> -r, Undefined
			}
			return left << r, Undefined
		case UInt64:
			return left << r, Undefined
		case *BigInt:
			if r.IsSmallInt() {
				rSmall := r.ToSmallInt()
				if rSmall < 0 {
					return left >> -rSmall, Undefined
				}
				return left << rSmall, Undefined
			}

			return 0, Undefined
		default:
			return 0, Ref(NewBitshiftOperandError(right))
		}
	}

	switch right.ValueFlag() {
	case SMALL_INT_FLAG:
		r := right.AsSmallInt()
		if r < 0 {
			return left >> -r, Undefined
		}
		return left << r, Undefined
	case INT64_FLAG:
		r := right.AsInlineInt64()
		if r < 0 {
			return left >> -r, Undefined
		}
		return left << r, Undefined
	case INT32_FLAG:
		r := right.AsInt32()
		if r < 0 {
			return left >> -r, Undefined
		}
		return left << r, Undefined
	case INT16_FLAG:
		r := right.AsInt16()
		if r < 0 {
			return left >> -r, Undefined
		}
		return left << r, Undefined
	case INT8_FLAG:
		r := right.AsInt8()
		if r < 0 {
			return left >> -r, Undefined
		}
		return left << r, Undefined
	case UINT64_FLAG:
		r := right.AsInlineUInt64()
		return left << r, Undefined
	case UINT32_FLAG:
		r := right.AsUInt32()
		return left << r, Undefined
	case UINT16_FLAG:
		r := right.AsUInt16()
		return left << r, Undefined
	case UINT8_FLAG:
		r := right.AsUInt8()
		return left << r, Undefined
	default:
		return 0, Ref(NewBitshiftOperandError(right))
	}
}

// Check whether the left float is equal to right.
func StrictFloatLaxEqual[T StrictFloat](left T, right Value) Value {
	if right.IsReference() {
		switch r := right.AsReference().(type) {
		case *BigInt:
			return ToElkBool(T(left) == T(r.ToFloat()))
		case *BigFloat:
			if r.IsNaN() {
				return False
			}
			iBigFloat := (&big.Float{}).SetFloat64(float64(left))
			return ToElkBool(iBigFloat.Cmp(r.AsGoBigFloat()) == 0)
		case Int64:
			return ToElkBool(T(left) == T(r))
		case UInt64:
			return ToElkBool(T(left) == T(r))
		case Float64:
			return ToElkBool(float64(left) == float64(r))
		default:
			return False
		}
	}

	switch right.ValueFlag() {
	case SMALL_INT_FLAG:
		r := right.AsSmallInt()
		return ToElkBool(T(left) == T(r))
	case FLOAT_FLAG:
		r := right.AsFloat()
		return ToElkBool(float64(left) == float64(r))
	case INT64_FLAG:
		r := right.AsInlineInt64()
		return ToElkBool(T(left) == T(r))
	case INT32_FLAG:
		r := right.AsInt32()
		return ToElkBool(T(left) == T(r))
	case INT16_FLAG:
		r := right.AsInt16()
		return ToElkBool(T(left) == T(r))
	case INT8_FLAG:
		r := right.AsInt8()
		return ToElkBool(T(left) == T(r))
	case UINT64_FLAG:
		r := right.AsInlineUInt64()
		return ToElkBool(T(left) == T(r))
	case UINT32_FLAG:
		r := right.AsUInt32()
		return ToElkBool(T(left) == T(r))
	case UINT16_FLAG:
		r := right.AsUInt16()
		return ToElkBool(T(left) == T(r))
	case UINT8_FLAG:
		r := right.AsUInt8()
		return ToElkBool(T(left) == T(r))
	case FLOAT64_FLAG:
		r := right.AsInlineFloat64()
		return ToElkBool(float64(left) == float64(r))
	case FLOAT32_FLAG:
		r := right.AsFloat32()
		return ToElkBool(T(left) == T(r))
	default:
		return False
	}
}

// Check whether the left signed integer is equal to right.
func StrictSignedIntLaxEqual[T StrictSignedInt](left T, right Value) Value {
	if right.IsReference() {
		switch r := right.AsReference().(type) {
		case *BigInt:
			iBigInt := big.NewInt(int64(left))
			return ToElkBool(iBigInt.Cmp(r.ToGoBigInt()) == 0)
		case *BigFloat:
			if r.IsNaN() {
				return False
			}
			iBigFloat := (&big.Float{}).SetInt64(int64(left))
			return ToElkBool(iBigFloat.Cmp(r.AsGoBigFloat()) == 0)
		case Int64:
			return ToElkBool(int64(left) == int64(r))
		case UInt64:
			if r > math.MaxInt64 {
				return False
			}
			return ToElkBool(int64(left) == int64(r))
		case Float64:
			return ToElkBool(float64(left) == float64(r))
		default:
			return False
		}
	}

	switch right.ValueFlag() {
	case SMALL_INT_FLAG:
		r := right.AsSmallInt()
		return ToElkBool(int64(left) == int64(r))
	case FLOAT_FLAG:
		r := right.AsFloat()
		return ToElkBool(float64(left) == float64(r))
	case INT64_FLAG:
		r := right.AsInlineInt64()
		return ToElkBool(int64(left) == int64(r))
	case INT32_FLAG:
		r := right.AsInt32()
		return ToElkBool(int64(left) == int64(r))
	case INT16_FLAG:
		r := right.AsInt16()
		return ToElkBool(int64(left) == int64(r))
	case INT8_FLAG:
		r := right.AsInt8()
		return ToElkBool(left == T(r))
	case UINT64_FLAG:
		r := right.AsInlineUInt64()
		if r > math.MaxInt64 {
			return False
		}
		return ToElkBool(int64(left) == int64(r))
	case UINT32_FLAG:
		r := right.AsUInt32()
		return ToElkBool(int64(left) == int64(r))
	case UINT16_FLAG:
		r := right.AsUInt16()
		return ToElkBool(int64(left) == int64(r))
	case UINT8_FLAG:
		r := right.AsUInt8()
		return ToElkBool(int64(left) == int64(r))
	case FLOAT64_FLAG:
		r := right.AsInlineFloat64()
		return ToElkBool(float64(left) == float64(r))
	case FLOAT32_FLAG:
		r := right.AsFloat32()
		return ToElkBool(float64(left) == float64(r))
	default:
		return False
	}
}

// Check whether the left unsigned integer is equal to right.
func StrictUnsignedIntLaxEqual[T StrictUnsignedInt](left T, right Value) Value {
	if right.IsReference() {
		switch r := right.AsReference().(type) {
		case *BigInt:
			iBigInt := (&big.Int{}).SetUint64(uint64(left))
			return ToElkBool(iBigInt.Cmp(r.ToGoBigInt()) == 0)
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
		case UInt64:
			return ToElkBool(uint64(left) == uint64(r))
		case Float64:
			return ToElkBool(float64(left) == float64(r))
		default:
			return False
		}
	}

	switch right.ValueFlag() {
	case SMALL_INT_FLAG:
		r := right.AsSmallInt()
		if uint64(left) > math.MaxInt64 {
			return False
		}
		return ToElkBool(int64(left) == int64(r))
	case FLOAT_FLAG:
		r := right.AsFloat()
		return ToElkBool(float64(left) == float64(r))
	case INT64_FLAG:
		r := right.AsInlineInt64()
		if uint64(left) > math.MaxInt64 {
			return False
		}
		return ToElkBool(int64(left) == int64(r))
	case INT32_FLAG:
		r := right.AsInt32()
		if uint64(left) > math.MaxInt64 {
			return False
		}
		return ToElkBool(int64(left) == int64(r))
	case INT16_FLAG:
		r := right.AsInt16()
		if uint64(left) > math.MaxInt64 {
			return False
		}
		return ToElkBool(int64(left) == int64(r))
	case INT8_FLAG:
		r := right.AsInt8()
		if uint64(left) > math.MaxInt64 {
			return False
		}
		return ToElkBool(int64(left) == int64(r))
	case UINT64_FLAG:
		r := right.AsInlineUInt64()
		return ToElkBool(uint64(left) == uint64(r))
	case UINT32_FLAG:
		r := right.AsUInt32()
		return ToElkBool(uint64(left) == uint64(r))
	case UINT16_FLAG:
		r := right.AsUInt16()
		return ToElkBool(uint64(left) == uint64(r))
	case UINT8_FLAG:
		r := right.AsUInt8()
		return ToElkBool(left == T(r))
	case FLOAT64_FLAG:
		r := right.AsInlineFloat64()
		return ToElkBool(float64(left) == float64(r))
	case FLOAT32_FLAG:
		r := right.AsFloat32()
		return ToElkBool(float64(left) == float64(r))
	default:
		return False
	}
}

// Parses an unsigned strict integer from a string using Elk syntax.
func StrictParseUintWithErr(s string, base int, bitSize int, formatErr *Class) (uint64, Value) {
	if s == "" {
		return 0, Ref(Errorf(formatErr, "invalid integer format"))
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
		return 0, Ref(Errorf(formatErr, "invalid integer base %d", base))
	}

	if bitSize == 0 {
		bitSize = strconv.IntSize
	} else if bitSize < 0 || bitSize > 64 {
		return 0, Ref(Errorf(formatErr, "invalid integer bit size %d", bitSize))
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
			return 0, Ref(Errorf(formatErr, "illegal characters in integer: %c", c))
		}

		if d >= byte(base) {
			return 0, Ref(Errorf(formatErr, "illegal characters in integer (base %d): %c", base, c))
		}

		if n >= cutoff {
			// n*base overflows
			return maxVal, Ref(Errorf(formatErr, "value overflows"))
		}
		n *= uint64(base)

		n1 := n + uint64(d)
		if n1 < n || n1 > maxVal {
			// n+d overflows
			return maxVal, Ref(Errorf(formatErr, "value overflows"))
		}
		n = n1
	}

	return n, Undefined
}

// Parses an unsigned strict integer from a string using Elk syntax.
func StrictParseUint(s string, base int, bitSize int) (uint64, Value) {
	return StrictParseUintWithErr(s, base, bitSize, FormatErrorClass)
}

// Parses a signed strict integer from a string using Elk syntax.
func StrictParseIntWithErr(s string, base int, bitSize int, formatErr *Class) (int64, Value) {
	if s == "" {
		return 0, Ref(Errorf(formatErr, "invalid integer format"))
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
	var un uint64
	un, err := StrictParseUint(s, base, bitSize)
	if !err.IsUndefined() {
		return 0, err
	}

	if bitSize == 0 {
		bitSize = strconv.IntSize
	}

	cutoff := uint64(1 << uint(bitSize-1))
	if !neg && un >= cutoff {
		return int64(cutoff - 1), Ref(Errorf(formatErr, "value overflows"))
	}
	if neg && un > cutoff {
		return -int64(cutoff), Ref(Errorf(formatErr, "value overflows"))
	}
	n := int64(un)
	if neg {
		n = -n
	}
	return n, Undefined
}

// Parses a signed strict integer from a string using Elk syntax.
func StrictParseInt(s string, base int, bitSize int) (int64, Value) {
	return StrictParseIntWithErr(s, base, bitSize, FormatErrorClass)
}

// Converts letters to lowercase.
func letterToLower(c byte) byte {
	return c | ('x' - 'X')
}
