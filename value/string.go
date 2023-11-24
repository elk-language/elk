package value

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/rivo/uniseg"
)

var StringClass *Class // ::Std::String

// Elk's String value
type String string

func (s String) Class() *Class {
	return StringClass
}

func (String) DirectClass() *Class {
	return StringClass
}

func (String) SingletonClass() *Class {
	return nil
}

func (s String) Inspect() string {
	return fmt.Sprintf("%q", s)
}

func (s String) InstanceVariables() SymbolMap {
	return nil
}

// Cmp compares x and y and returns:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
func (x String) Cmp(y String) int {
	return strings.Compare(string(x), string(y))
}

// Convert this String to a Char.
// Returns (Char, true) if the conversion was successful.
// Returns (0, false) if the conversion failed.
func (s String) ToChar() (Char, bool) {
	r, size := utf8.DecodeRuneInString(string(s))
	if size != len(s) {
		return 0, false
	}

	return Char(r), true
}

// Returns the number of bytes
// present in the string.
func (s String) ByteCount() int {
	return len(s)
}

// Returns the number of unicode chars
// present in the string.
func (s String) CharCount() int {
	return utf8.RuneCountInString(string(s))
}

// Returns the number of grapheme clusters
// present in the string.
func (s String) GraphemeCount() int {
	return uniseg.GraphemeClusterCount(string(s))
}

// Reverse the bytes of the string.
func (s String) ReverseBytes() String {
	a := []byte(s)
	for i, j := 0, len(s)-1; i < j; i++ {
		a[i], a[j] = a[j], a[i]
		j--
	}
	return String(a)
}

// Reverse the string while preserving the UTF-8 chars.
func (s String) ReverseChars() String {
	str := string(s)
	reversed := make([]byte, len(str))
	i := 0

	for len(str) > 0 {
		r, size := utf8.DecodeLastRuneInString(str)
		str = str[:len(str)-size]
		i += utf8.EncodeRune(reversed[i:], r)
	}

	return String(reversed)
}

// Reverse the string while preserving the grapheme clusters.
func (s String) ReverseGraphemes() String {
	return String(uniseg.ReverseString(string(s)))
}

// Concatenate another value with this string and return the result.
// If the operation is illegal an error will be returned.
func (s String) Concat(other Value) (String, *Error) {
	switch o := other.(type) {
	case Char:
		var buff strings.Builder
		buff.WriteString(string(s))
		buff.WriteRune(rune(o))
		return String(buff.String()), nil
	case String:
		return s + o, nil
	default:
		return "", Errorf(TypeErrorClass, "can't concat %s to string %s", other.Inspect(), s.Inspect())
	}
}

// Repeat the content of this string n times and return a new string containing the result.
// If the operation is illegal an error will be returned.
func (s String) Repeat(other Value) (String, *Error) {
	switch o := other.(type) {
	case SmallInt:
		return String(strings.Repeat(string(s), int(o))), nil
	case *BigInt:
		return "", Errorf(
			OutOfRangeErrorClass,
			"repeat count is too large %s",
			o.Inspect(),
		)
	default:
		return "", Errorf(TypeErrorClass, "can't repeat a string using %s", other.Inspect())
	}
}

// Return a copy of the string without the given suffix.
func (s String) RemoveSuffix(other Value) (String, *Error) {
	switch o := other.(type) {
	case Char:
		r, rLen := utf8.DecodeLastRuneInString(string(s))
		if len(s) > 0 && r == rune(o) {
			return s[0 : len(s)-rLen], nil
		}
		return s, nil
	case String:
		result, _ := strings.CutSuffix(string(s), string(o))
		return String(result), nil
	default:
		return "", NewCoerceError(s, other)
	}
}

// Check whether s is greater than other and return an error
// if something went wrong.
func (s String) GreaterThan(other Value) (Value, *Error) {
	switch o := other.(type) {
	case String:
		return ToElkBool(s.Cmp(o) == 1), nil
	case Char:
		return ToElkBool(s.Cmp(String(o)) == 1), nil
	default:
		return nil, NewCoerceError(s, other)
	}
}

// Check whether s is greater than or equal to other and return an error
// if something went wrong.
func (s String) GreaterThanEqual(other Value) (Value, *Error) {
	switch o := other.(type) {
	case String:
		return ToElkBool(s.Cmp(o) >= 0), nil
	case Char:
		return ToElkBool(s.Cmp(String(o)) >= 0), nil
	default:
		return nil, NewCoerceError(s, other)
	}
}

// Check whether s is less than other and return an error
// if something went wrong.
func (s String) LessThan(other Value) (Value, *Error) {
	switch o := other.(type) {
	case String:
		return ToElkBool(s.Cmp(o) == -1), nil
	case Char:
		return ToElkBool(s.Cmp(String(o)) == -1), nil
	default:
		return nil, NewCoerceError(s, other)
	}
}

// Check whether s is less than or equal to other and return an error
// if something went wrong.
func (s String) LessThanEqual(other Value) (Value, *Error) {
	switch o := other.(type) {
	case String:
		return ToElkBool(s.Cmp(o) <= 0), nil
	case Char:
		return ToElkBool(s.Cmp(String(o)) <= 0), nil
	default:
		return nil, NewCoerceError(s, other)
	}
}

// Check whether s is equal to other
func (s String) Equal(other Value) Value {
	switch o := other.(type) {
	case String:
		return ToElkBool(s == o)
	case Char:
		ch, ok := s.ToChar()
		if !ok {
			return False
		}

		return ToElkBool(ch == o)
	default:
		return False
	}
}

// Check whether s is strictly equal to other
func (s String) StrictEqual(other Value) Value {
	switch o := other.(type) {
	case String:
		return ToElkBool(s == o)
	default:
		return False
	}
}

// Convert the String to a Symbol
func (s String) ToSymbol() Symbol {
	return SymbolTable.Add(string(s))
}

func initString() {
	StringClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("String", StringClass)
}
