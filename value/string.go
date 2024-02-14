package value

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/cespare/xxhash/v2"
	"github.com/rivo/uniseg"
)

var StringClass *Class             // ::Std::String
var StringCharIteratorClass *Class // ::Std::String::CharIterator
var StringByteIteratorClass *Class // ::Std::String::ByteIterator

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

func (s String) Copy() Value {
	return s
}

func (s String) Inspect() string {
	var buffer strings.Builder

	buffer.WriteString(`"`)
	leftStr := string(s)
	for {
		char, size := utf8.DecodeRuneInString(leftStr)
		if size == 0 {
			// reached the end of the string
			break
		}
		if char == utf8.RuneError && size == 1 {
			// invalid UTF-8 character
			char = rune(leftStr[0])
		}
		switch char {
		case '\\':
			buffer.WriteString(`\\`)
		case '\n':
			buffer.WriteString(`\n`)
		case '\t':
			buffer.WriteString(`\t`)
		case '"':
			buffer.WriteString(`\"`)
		case '\r':
			buffer.WriteString(`\r`)
		case '\a':
			buffer.WriteString(`\a`)
		case '\b':
			buffer.WriteString(`\b`)
		case '\v':
			buffer.WriteString(`\v`)
		case '\f':
			buffer.WriteString(`\f`)
		default:
			if unicode.IsGraphic(char) {
				buffer.WriteRune(char)
			} else if char>>8 == 0 {
				fmt.Fprintf(&buffer, `\x%02x`, char)
			} else if char>>16 == 0 {
				fmt.Fprintf(&buffer, `\u%04x`, char)
			} else {
				fmt.Fprintf(&buffer, `\U%08X`, char)
			}
		}
		leftStr = leftStr[size:]
	}

	buffer.WriteString(`"`)
	return buffer.String()
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

func (s String) IsEmpty() bool {
	return len(s) == 0
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
		return "", Errorf(TypeErrorClass, "cannot concat %s to string %s", other.Inspect(), s.Inspect())
	}
}

// Repeat the content of this string n times and return a new string containing the result.
// If the operation is illegal an error will be returned.
func (s String) Repeat(other Value) (String, *Error) {
	switch o := other.(type) {
	case SmallInt:
		if o < 0 {
			return "", Errorf(
				OutOfRangeErrorClass,
				"repeat count cannot be negative: %s",
				o.Inspect(),
			)
		}
		return String(strings.Repeat(string(s), int(o))), nil
	case *BigInt:
		return "", Errorf(
			OutOfRangeErrorClass,
			"repeat count is too large %s",
			o.Inspect(),
		)
	default:
		return "", Errorf(TypeErrorClass, "cannot repeat a string using %s", other.Inspect())
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
		return "", NewCoerceError(s.Class(), other.Class())
	}
}

// Returns 1 if i is greater than other
// Returns 0 if both are equal.
// Returns -1 if i is less than other.
// Returns nil if the comparison was impossible (NaN)
func (s String) Compare(other Value) (Value, *Error) {
	switch o := other.(type) {
	case String:
		return SmallInt(s.Cmp(o)), nil
	case Char:
		return SmallInt(s.Cmp(String(o))), nil
	default:
		return nil, NewCoerceError(s.Class(), other.Class())
	}
}

// Check whether s is greater than other and return an error
// if something went wrong.
func (s String) GreaterThan(other Value) (Value, *Error) {
	switch o := other.(type) {
	case String:
		return ToElkBool(s > o), nil
	case Char:
		return ToElkBool(s > String(o)), nil
	default:
		return nil, NewCoerceError(s.Class(), other.Class())
	}
}

// Check whether s is greater than or equal to other and return an error
// if something went wrong.
func (s String) GreaterThanEqual(other Value) (Value, *Error) {
	switch o := other.(type) {
	case String:
		return ToElkBool(s >= o), nil
	case Char:
		return ToElkBool(s >= String(o)), nil
	default:
		return nil, NewCoerceError(s.Class(), other.Class())
	}
}

// Check whether s is less than other and return an error
// if something went wrong.
func (s String) LessThan(other Value) (Value, *Error) {
	switch o := other.(type) {
	case String:
		return ToElkBool(s < o), nil
	case Char:
		return ToElkBool(s < String(o)), nil
	default:
		return nil, NewCoerceError(s.Class(), other.Class())
	}
}

// Check whether s is less than or equal to other and return an error
// if something went wrong.
func (s String) LessThanEqual(other Value) (Value, *Error) {
	switch o := other.(type) {
	case String:
		return ToElkBool(s <= o), nil
	case Char:
		return ToElkBool(s <= String(o)), nil
	default:
		return nil, NewCoerceError(s.Class(), other.Class())
	}
}

// Check whether s is equal to other
func (s String) LaxEqual(other Value) Value {
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

// Check whether s is equal to other
func (s String) Equal(other Value) Value {
	switch o := other.(type) {
	case String:
		return ToElkBool(s == o)
	default:
		return False
	}
}

// Check whether s is strictly equal to other
func (s String) StrictEqual(other Value) Value {
	return s.Equal(other)
}

// Get an element under the given index.
func (s String) Get(index int) (Char, *Error) {
	var i int
	if index < 0 {
		l := s.CharCount()
		i = l + index
		if i < 0 {
			return 0, NewIndexOutOfRangeError(fmt.Sprint(index), l)
		}
	} else {
		i = index
	}

	leftStr := string(s)
	var j int
	for {
		result, size := utf8.DecodeRuneInString(leftStr)
		if size == 0 {
			// reached the end of the string
			return 0, NewIndexOutOfRangeError(fmt.Sprint(index), j)
		}
		if result == utf8.RuneError && size == 1 {
			// invalid UTF-8 character
			result = rune(leftStr[0])
		}
		if j == i {
			// found the character
			return Char(result), nil
		}
		leftStr = leftStr[size:]
		j++
	}
}

// Get the character under the given index.
func (s String) Subscript(key Value) (Char, *Error) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return 0, NewIndexOutOfRangeError(key.Inspect(), len(s))
		}
		return 0, NewCoerceError(IntClass, key.Class())
	}

	return s.Get(i)
}

// Get the byte under the given index.
func (s String) ByteAtInt(index int) (UInt8, *Error) {
	l := len(s)
	if index >= l || index < -l {
		return 0, NewIndexOutOfRangeError(fmt.Sprint(index), l)
	}
	if index < 0 {
		index = l + index
	}
	return UInt8(s[index]), nil
}

// Get the byte under the given index.
func (s String) ByteAt(key Value) (UInt8, *Error) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return 0, NewIndexOutOfRangeError(key.Inspect(), len(s))
		}
		return 0, NewCoerceError(IntClass, key.Class())
	}

	return s.ByteAtInt(i)
}

// Get the grapheme under the given index.
func (s String) GraphemeAtInt(index int) (String, *Error) {
	var i int
	if index < 0 {
		l := s.GraphemeCount()
		i = l + index
		if i < 0 {
			return "", NewIndexOutOfRangeError(fmt.Sprint(index), l)
		}
	} else {
		i = index
	}

	str := string(s)
	state := -1
	var cluster string
	var j int
	for len(str) > 0 {
		cluster, str, _, state = uniseg.FirstGraphemeClusterInString(str, state)
		if j == i {
			// found the grapheme
			return String(cluster), nil
		}
		j++
	}

	return "", NewIndexOutOfRangeError(fmt.Sprint(index), j)
}

// Get the grapheme under the given index.
func (s String) GraphemeAt(key Value) (String, *Error) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return "", NewIndexOutOfRangeError(key.Inspect(), len(s))
		}
		return "", NewCoerceError(IntClass, key.Class())
	}

	return s.GraphemeAtInt(i)
}

func (s String) Hash() UInt64 {
	d := xxhash.New()
	d.WriteString(string(s))
	return UInt64(d.Sum64())
}

// Convert the String to a Symbol
func (s String) ToSymbol() Symbol {
	return SymbolTable.Add(string(s))
}

type StringCharIterator struct {
	String     String
	ByteOffset int
}

func NewStringCharIterator(str String) *StringCharIterator {
	return &StringCharIterator{
		String: str,
	}
}

func NewStringCharIteratorWithByteOffset(str String, offset int) *StringCharIterator {
	return &StringCharIterator{
		String:     str,
		ByteOffset: offset,
	}
}

func (*StringCharIterator) Class() *Class {
	return StringCharIteratorClass
}

func (*StringCharIterator) DirectClass() *Class {
	return StringCharIteratorClass
}

func (*StringCharIterator) SingletonClass() *Class {
	return nil
}

func (s *StringCharIterator) Copy() Value {
	return &StringCharIterator{
		String:     s.String,
		ByteOffset: s.ByteOffset,
	}
}

func (s *StringCharIterator) Inspect() string {
	return fmt.Sprintf("Std::String::CharIterator{string: %s, byte_offset: %d}", s.String.Inspect(), s.ByteOffset)
}

func (*StringCharIterator) InstanceVariables() SymbolMap {
	return nil
}

func (s *StringCharIterator) Next() (Value, Value) {
	if s.ByteOffset >= len(s.String) {
		return nil, stopIterationSymbol
	}
	run, size := utf8.DecodeRuneInString(string(s.String[s.ByteOffset:]))

	s.ByteOffset += size
	return Char(run), nil
}

type StringByteIterator struct {
	String     String
	ByteOffset int
}

func NewStringByteIterator(str String) *StringByteIterator {
	return &StringByteIterator{
		String: str,
	}
}

func NewStringByteIteratorWithByteOffset(str String, offset int) *StringByteIterator {
	return &StringByteIterator{
		String:     str,
		ByteOffset: offset,
	}
}

func (*StringByteIterator) Class() *Class {
	return StringByteIteratorClass
}

func (*StringByteIterator) DirectClass() *Class {
	return StringByteIteratorClass
}

func (*StringByteIterator) SingletonClass() *Class {
	return nil
}

func (s *StringByteIterator) Copy() Value {
	return &StringByteIterator{
		String:     s.String,
		ByteOffset: s.ByteOffset,
	}
}

func (s *StringByteIterator) Inspect() string {
	return fmt.Sprintf("Std::String::ByteIterator{string: %s, byte_offset: %d}", s.String.Inspect(), s.ByteOffset)
}

func (*StringByteIterator) InstanceVariables() SymbolMap {
	return nil
}

func (s *StringByteIterator) Next() (Value, Value) {
	if s.ByteOffset >= len(s.String) {
		return nil, stopIterationSymbol
	}
	result := UInt8(s.String[s.ByteOffset])
	s.ByteOffset += 1
	return result, nil
}

func initString() {
	StringClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("String", StringClass)

	StringCharIteratorClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StringClass.AddConstantString("CharIterator", StringCharIteratorClass)

	StringByteIteratorClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StringClass.AddConstantString("ByteIterator", StringByteIteratorClass)
}
