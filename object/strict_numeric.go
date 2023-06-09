package object

import (
	"math"
	"strconv"
)

// Strict numerics are sized and can't be automatically coerced
// to other types.
type StrictNumeric interface {
	Float64 | Float32 | Int64 | Int32 | Int16 | Int8 | UInt64 | UInt32 | UInt16 | UInt8
	Value
}

// Add a strict numeric to another value and return the result.
// If the operation is illegal an error will be returned.
func StrictNumericAdd[T StrictNumeric](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, NewCoerceError(left, r)
	}

	return left + r, nil
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
