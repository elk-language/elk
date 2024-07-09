package value

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ::Std::Char
//
// Represents a single UTF-8 character.
// Takes up 4 bytes.
var CharClass *Class

// Elk's Char value
type Char rune

func (Char) Class() *Class {
	return CharClass
}

func (Char) DirectClass() *Class {
	return CharClass
}

func (Char) SingletonClass() *Class {
	return nil
}

func (c Char) Copy() Value {
	return c
}

func (c Char) Inspect() string {
	var buff strings.Builder
	buff.WriteRune('`')
	switch c {
	case '\\':
		buff.WriteString(`\\`)
	case '\n':
		buff.WriteString(`\n`)
	case '\t':
		buff.WriteString(`\t`)
	case '`':
		buff.WriteString("\\`")
	case '\r':
		buff.WriteString(`\r`)
	case '\a':
		buff.WriteString(`\a`)
	case '\b':
		buff.WriteString(`\b`)
	case '\v':
		buff.WriteString(`\v`)
	case '\f':
		buff.WriteString(`\f`)
	default:
		if unicode.IsGraphic(rune(c)) {
			buff.WriteRune(rune(c))
		} else if c>>8 == 0 {
			fmt.Fprintf(&buff, `\x%02x`, c)
		} else if c>>16 == 0 {
			fmt.Fprintf(&buff, `\u%04x`, c)
		} else {
			fmt.Fprintf(&buff, `\U%08X`, c)
		}
	}

	buff.WriteRune('`')
	return buff.String()
}

func (Char) InstanceVariables() SymbolMap {
	return nil
}

// Returns the number of bytes
// present in the character.
func (c Char) ByteCount() int {
	return utf8.RuneLen(rune(c))
}

func (Char) CharCount() int {
	return 1
}

func (Char) GraphemeCount() int {
	return 1
}

// Cmp compares x and y and returns:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
func (x Char) Cmp(y Char) int {
	if x > y {
		return 1
	}
	if x < y {
		return -1
	}

	return 0
}

// Concatenate another value with this character, creating a new string, and return the result.
// If the operation is illegal an error will be returned.
func (c Char) Concat(other Value) (String, *Error) {
	switch o := other.(type) {
	case Char:
		var buff strings.Builder
		buff.WriteRune(rune(c))
		buff.WriteRune(rune(o))
		return String(buff.String()), nil
	case String:
		var buff strings.Builder
		buff.WriteRune(rune(c))
		buff.WriteString(string(o))
		return String(buff.String()), nil
	default:
		return "", Errorf(TypeErrorClass, "cannot concat %s with char %s", other.Inspect(), c.Inspect())
	}
}

// Repeat this character n times and return a new string containing the result.
// If the operation is illegal an error will be returned.
func (c Char) Repeat(other Value) (String, *Error) {
	switch o := other.(type) {
	case SmallInt:
		if o < 0 {
			return "", Errorf(
				OutOfRangeErrorClass,
				"repeat count cannot be negative: %s",
				o.Inspect(),
			)
		}
		var builder strings.Builder
		for i := 0; i < int(o); i++ {
			builder.WriteRune(rune(c))
		}
		return String(builder.String()), nil
	case *BigInt:
		return "", Errorf(
			OutOfRangeErrorClass,
			"repeat count is too large %s",
			o.Inspect(),
		)
	default:
		return "", Errorf(TypeErrorClass, "cannot repeat a char using %s", other.Inspect())
	}
}

// Returns 1 if i is greater than other
// Returns 0 if both are equal.
// Returns -1 if i is less than other.
// Returns nil if the comparison was impossible (NaN)
func (c Char) Compare(other Value) (Value, *Error) {
	switch o := other.(type) {
	case Char:
		return SmallInt(c.Cmp(o)), nil
	case String:
		return SmallInt(String(c).Cmp(o)), nil
	default:
		return nil, NewCoerceError(c.Class(), other.Class())
	}
}

// Check whether c is greater than other and return an error
// if something went wrong.
func (c Char) GreaterThan(other Value) (Value, *Error) {
	switch o := other.(type) {
	case Char:
		return ToElkBool(c > o), nil
	case String:
		return ToElkBool(String(c).Cmp(o) == 1), nil
	default:
		return nil, NewCoerceError(c.Class(), other.Class())
	}
}

// Check whether c is greater than or equal to other and return an error
// if something went wrong.
func (c Char) GreaterThanEqual(other Value) (Value, *Error) {
	switch o := other.(type) {
	case Char:
		return ToElkBool(c >= o), nil
	case String:
		return ToElkBool(String(c).Cmp(o) >= 0), nil
	default:
		return nil, NewCoerceError(c.Class(), other.Class())
	}
}

// Check whether c is less than other and return an error
// if something went wrong.
func (c Char) LessThan(other Value) (Value, *Error) {
	switch o := other.(type) {
	case Char:
		return ToElkBool(c < o), nil
	case String:
		return ToElkBool(String(c).Cmp(o) == -1), nil
	default:
		return nil, NewCoerceError(c.Class(), other.Class())
	}
}

// Check whether c is less than or equal to other and return an error
// if something went wrong.
func (c Char) LessThanEqual(other Value) (Value, *Error) {
	switch o := other.(type) {
	case Char:
		return ToElkBool(c <= o), nil
	case String:
		return ToElkBool(String(c).Cmp(o) <= 0), nil
	default:
		return nil, NewCoerceError(c.Class(), other.Class())
	}
}

// Check whether c is equal to other
func (c Char) LaxEqual(other Value) Value {
	switch o := other.(type) {
	case Char:
		return ToElkBool(c == o)
	case String:
		ch, ok := o.ToChar()
		if !ok {
			return False
		}

		return ToElkBool(c == ch)
	default:
		return False
	}
}

// Check whether s is equal to other
func (c Char) Equal(other Value) Value {
	switch o := other.(type) {
	case Char:
		return ToElkBool(c == o)
	default:
		return False
	}
}

// Convert to uppercase
func (c Char) Uppercase() Char {
	return Char(unicode.ToUpper(rune(c)))
}

// Convert to lowercase
func (c Char) Lowercase() Char {
	return Char(unicode.ToLower(rune(c)))
}

// Check whether s is strictly equal to other
func (c Char) StrictEqual(other Value) Value {
	return c.Equal(other)
}

func initChar() {
	CharClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Char", CharClass)
}
