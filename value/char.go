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

func (c Char) Inspect() string {
	var content string
	switch c {
	case '\\':
		content = `\\`
	case '\n':
		content = `\n`
	case '\t':
		content = `\t`
	case '"':
		content = `\"`
	case '\r':
		content = `\r`
	case '\a':
		content = `\a`
	case '\b':
		content = `\b`
	case '\v':
		content = `\v`
	case '\f':
		content = `\f`
	default:
		if unicode.IsGraphic(rune(c)) {
			content = string(c)
		} else if utf8.RuneLen(rune(c)) == 1 {
			content = fmt.Sprintf(`\x%02x`, c)
		} else {
			content = fmt.Sprintf(`\U%08X`, c)
		}
	}

	return fmt.Sprintf(`c"%s"`, content)
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
		return "", Errorf(TypeErrorClass, "can't concat %s with char %s", other.Inspect(), c.Inspect())
	}
}

// Repeat this character n times and return a new string containing the result.
// If the operation is illegal an error will be returned.
func (c Char) Repeat(other Value) (String, *Error) {
	switch o := other.(type) {
	case SmallInt:
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
		return "", Errorf(TypeErrorClass, "can't repeat a char using %s", other.Inspect())
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
		return nil, NewCoerceError(c, other)
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
		return nil, NewCoerceError(c, other)
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
		return nil, NewCoerceError(c, other)
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
		return nil, NewCoerceError(c, other)
	}
}

// Check whether c is equal to other
func (c Char) Equal(other Value) Value {
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

// Check whether s is strictly equal to other
func (c Char) StrictEqual(other Value) Value {
	switch o := other.(type) {
	case Char:
		return ToElkBool(c == o)
	default:
		return False
	}
}

func initChar() {
	CharClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Char", CharClass)
}
