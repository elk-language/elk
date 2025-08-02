package value

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/cespare/xxhash/v2"
)

// ::Std::Char
//
// Represents a single UTF-8 character.
// Takes up 4 bytes.
var CharClass *Class

// Elk's Char value
type Char rune

func (c Char) ToValue() Value {
	return Value{
		flag: CHAR_FLAG,
		data: uintptr(c),
	}
}

func (Char) Class() *Class {
	return CharClass
}

func (Char) DirectClass() *Class {
	return CharClass
}

func (Char) SingletonClass() *Class {
	return nil
}

func (c Char) Error() string {
	return c.Inspect()
}

func (c Char) Rune() rune {
	return rune(c)
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
			fmt.Fprintf(&buff, `\x%02x`, c.Rune())
		} else if c>>16 == 0 {
			fmt.Fprintf(&buff, `\u%04x`, c.Rune())
		} else {
			fmt.Fprintf(&buff, `\U%08X`, c.Rune())
		}
	}

	buff.WriteRune('`')
	return buff.String()
}

func (Char) InstanceVariables() *InstanceVariables {
	return nil
}

func (c Char) Hash() UInt64 {
	d := xxhash.New()
	d.WriteString(string(c))
	return UInt64(d.Sum64())
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
func (c Char) Concat(other Value) (String, Value) {
	if other.IsChar() {
		var buff strings.Builder
		buff.WriteRune(rune(c))
		buff.WriteRune(rune(other.AsChar()))
		return String(buff.String()), Undefined
	}

	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			var buff strings.Builder
			buff.WriteRune(rune(c))
			buff.WriteString(string(o))
			return String(buff.String()), Undefined
		}
	}

	return "", Ref(Errorf(TypeErrorClass, "cannot concat %s with char %s", other.Inspect(), c.Inspect()))
}

// Repeat this character n times and return a new string containing the result.
// If the operation is illegal an error will be returned.
func (c Char) Repeat(other Value) (String, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return "", Ref(Errorf(
				OutOfRangeErrorClass,
				"repeat count is too large %s",
				o.Inspect(),
			))
		default:
			return "", Ref(Errorf(TypeErrorClass, "cannot repeat a char using %s", other.Inspect()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		if o < 0 {
			return "", Ref(Errorf(
				OutOfRangeErrorClass,
				"repeat count cannot be negative: %s",
				o.Inspect(),
			))
		}
		var builder strings.Builder
		for i := 0; i < int(o); i++ {
			builder.WriteRune(rune(c))
		}
		return String(builder.String()), Undefined
	default:
		return "", Ref(Errorf(TypeErrorClass, "cannot repeat a char using %s", other.Inspect()))
	}
}

// Returns 1 if i is greater than other
// Returns 0 if both are equal.
// Returns -1 if i is less than other.
// Returns nil if the comparison was impossible (NaN)
func (c Char) CompareVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			return SmallInt(String(c).Cmp(o)).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(c.Class(), other.Class()))
		}
	}
	switch other.ValueFlag() {
	case CHAR_FLAG:
		return SmallInt(c.Cmp(other.AsChar())).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(c.Class(), other.Class()))
	}
}

// Check whether c is greater than other and return an error
// if something went wrong.
func (c Char) GreaterThanVal(other Value) (Value, Value) {
	result, err := c.GreaterThan(other)
	return ToElkBool(result), err
}

// Check whether c is greater than other and return an error
// if something went wrong.
func (c Char) GreaterThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			return String(c).Cmp(o) == 1, Undefined
		default:
			return false, Ref(NewCoerceError(c.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case CHAR_FLAG:
		return c > other.AsChar(), Undefined
	default:
		return false, Ref(NewCoerceError(c.Class(), other.Class()))
	}
}

// Check whether c is greater than or equal to other and return an error
// if something went wrong.
func (c Char) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := c.GreaterThanEqual(other)
	return ToElkBool(result), err
}

// Check whether c is greater than or equal to other and return an error
// if something went wrong.
func (c Char) GreaterThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			return String(c).Cmp(o) >= 0, Undefined
		default:
			return false, Ref(NewCoerceError(c.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case CHAR_FLAG:
		return c >= other.AsChar(), Undefined
	default:
		return false, Ref(NewCoerceError(c.Class(), other.Class()))
	}
}

// Check whether c is less than other and return an error
// if something went wrong.
func (c Char) LessThanVal(other Value) (Value, Value) {
	result, err := c.LessThan(other)
	return ToElkBool(result), err
}

// Check whether c is less than other and return an error
// if something went wrong.
func (c Char) LessThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			return String(c).Cmp(o) == -1, Undefined
		default:
			return false, Ref(NewCoerceError(c.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case CHAR_FLAG:
		return c < other.AsChar(), Undefined
	default:
		return false, Ref(NewCoerceError(c.Class(), other.Class()))
	}
}

// Check whether c is less than or equal to other and return an error
// if something went wrong.
func (c Char) LessThanEqualVal(other Value) (Value, Value) {
	result, err := c.LessThanEqual(other)
	return ToElkBool(result), err
}

// Check whether c is less than or equal to other and return an error
// if something went wrong.
func (c Char) LessThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			return String(c).Cmp(o) <= 0, Undefined
		default:
			return false, Ref(NewCoerceError(c.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case CHAR_FLAG:
		return c <= other.AsChar(), Undefined
	default:
		return false, Ref(NewCoerceError(c.Class(), other.Class()))
	}
}

// Check whether c is equal to other
func (c Char) LaxEqualVal(other Value) Value {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
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

	switch other.ValueFlag() {
	case CHAR_FLAG:
		return ToElkBool(c == other.AsChar())
	default:
		return False
	}
}

// Check whether s is equal to other
func (c Char) EqualVal(other Value) Value {
	return ToElkBool(c.Equal(other))
}

// Check whether s is equal to other
func (c Char) Equal(other Value) bool {
	if other.IsChar() {
		return c == other.AsChar()
	}

	return false
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
func (c Char) StrictEqualVal(other Value) Value {
	return c.EqualVal(other)
}

func initChar() {
	CharClass = NewClassWithOptions(ClassWithParent(ValueClass))
	StdModule.AddConstantString("Char", Ref(CharClass))
}
