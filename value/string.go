package value

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/cespare/xxhash/v2"
	"github.com/rivo/uniseg"
)

var StringClass *Class                 // ::Std::String
var StringCharIteratorClass *Class     // ::Std::String::CharIterator
var StringByteIteratorClass *Class     // ::Std::String::ByteIterator
var StringGraphemeIteratorClass *Class // ::Std::String::GraphemeIterator

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

func (s String) Copy() Reference {
	return s
}

func (s String) Error() string {
	return s.Inspect()
}

func (s String) String() string {
	return string(s)
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
		case '$':
			buffer.WriteString(`\$`)
		case '#':
			buffer.WriteString(`\#`)
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

// Convert this String to an Int.
func (s String) ToInt() (Value, Value) {
	return ParseInt(string(s), 0)
}

// Return a new string that has all characters turned to lowercase.
func (s String) Lowercase() String {
	return String(strings.ToLower(string(s)))
}

// Return a new string that has all characters turned to uppercase.
func (s String) Uppercase() String {
	return String(strings.ToUpper(string(s)))
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
func (s String) Concat(other Value) (result String, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			return s + o, err
		default:
			return result, Ref(Errorf(TypeErrorClass, "cannot concat %s to string %s", other.Inspect(), s.Inspect()))
		}
	}

	switch other.ValueFlag() {
	case CHAR_FLAG:
		ch := other.AsChar()
		var buff strings.Builder
		buff.WriteString(string(s))
		buff.WriteRune(rune(ch))
		return String(buff.String()), err
	default:
		return result, Ref(Errorf(TypeErrorClass, "cannot concat %s to string %s", other.Inspect(), s.Inspect()))
	}
}

// Repeat the content of this string n times and return a new string containing the result.
// If the operation is illegal an error will be returned.
func (s String) Repeat(other Value) (result String, err Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return "", Ref(Errorf(
				OutOfRangeErrorClass,
				"repeat count is too large %s",
				o.Inspect(),
			))
		default:
			return result, Ref(Errorf(TypeErrorClass, "cannot repeat a string using %s", other.Inspect()))
		}
	}
	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		o := other.AsSmallInt()
		if o < 0 {
			return result, Ref(Errorf(
				OutOfRangeErrorClass,
				"repeat count cannot be negative: %s",
				o.Inspect(),
			))
		}
		return String(strings.Repeat(string(s), int(o))), err
	default:
		return result, Ref(Errorf(TypeErrorClass, "cannot repeat a string using %s", other.Inspect()))
	}
}

// Return a copy of the string without the given suffix.
func (s String) RemoveSuffix(other Value) (String, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			result, _ := strings.CutSuffix(string(s), string(o))
			return String(result), Undefined
		default:
			return "", Ref(NewCoerceError(s.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case CHAR_FLAG:
		o := other.AsChar()
		r, rLen := utf8.DecodeLastRuneInString(string(s))
		if len(s) > 0 && r == rune(o) {
			return s[0 : len(s)-rLen], Undefined
		}
		return s, Undefined
	default:
		return "", Ref(NewCoerceError(s.Class(), other.Class()))
	}
}

// Returns 1 if i is greater than other
// Returns 0 if both are equal.
// Returns -1 if i is less than other.
// Returns nil if the comparison was impossible (NaN)
func (s String) Compare(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			return SmallInt(s.Cmp(o)).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(s.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case CHAR_FLAG:
		return SmallInt(s.Cmp(String(other.AsChar()))).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(s.Class(), other.Class()))
	}
}

// Check whether s is greater than other and return an error
// if something went wrong.
func (s String) GreaterThan(other Value) (Value, Value) {
	result, err := s.GreaterThanBool(other)
	return ToElkBool(result), err
}

// Check whether s is greater than other and return an error
// if something went wrong.
func (s String) GreaterThanBool(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			return s > o, Undefined
		default:
			return false, Ref(NewCoerceError(s.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case CHAR_FLAG:
		return s > String(other.AsChar()), Undefined
	default:
		return false, Ref(NewCoerceError(s.Class(), other.Class()))
	}
}

// Check whether s is greater than or equal to other and return an error
// if something went wrong.
func (s String) GreaterThanEqual(other Value) (Value, Value) {
	result, err := s.GreaterThanEqualBool(other)
	return ToElkBool(result), err
}

// Check whether s is greater than or equal to other and return an error
// if something went wrong.
func (s String) GreaterThanEqualBool(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			return s >= o, Undefined
		default:
			return false, Ref(NewCoerceError(s.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case CHAR_FLAG:
		return s >= String(other.AsChar()), Undefined
	default:
		return false, Ref(NewCoerceError(s.Class(), other.Class()))
	}
}

// Check whether s is less than other and return an error
// if something went wrong.
func (s String) LessThan(other Value) (Value, Value) {
	result, err := s.LessThanBool(other)
	return ToElkBool(result), err
}

// Check whether s is less than other and return an error
// if something went wrong.
func (s String) LessThanBool(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			return s < o, Undefined
		default:
			return false, Ref(NewCoerceError(s.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case CHAR_FLAG:
		return s < String(other.AsChar()), Undefined
	default:
		return false, Ref(NewCoerceError(s.Class(), other.Class()))
	}
}

// Check whether s is less than or equal to other and return an error
// if something went wrong.
func (s String) LessThanEqual(other Value) (Value, Value) {
	result, err := s.LessThanEqualBool(other)
	return ToElkBool(result), err
}

// Check whether s is less than or equal to other and return an error
// if something went wrong.
func (s String) LessThanEqualBool(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			return s <= o, Undefined
		default:
			return false, Ref(NewCoerceError(s.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case CHAR_FLAG:
		return s <= String(other.AsChar()), Undefined
	default:
		return false, Ref(NewCoerceError(s.Class(), other.Class()))
	}
}

// Check whether s is equal to other
func (s String) LaxEqual(other Value) Value {
	return ToElkBool(s.LaxEqualBool(other))
}

// Check whether s is equal to other
func (s String) LaxEqualBool(other Value) bool {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			return s == o
		default:
			return false
		}
	}

	switch other.ValueFlag() {
	case CHAR_FLAG:
		ch, ok := s.ToChar()
		if !ok {
			return false
		}

		return ch == other.AsChar()
	default:
		return false
	}
}

// Check whether s is equal to other
func (s String) Equal(other Value) Value {
	return ToElkBool(s.EqualBool(other))
}

// Check whether s is equal to other
func (s String) EqualBool(other Value) bool {
	if !other.IsReference() {
		return false
	}
	switch o := other.AsReference().(type) {
	case String:
		return s == o
	default:
		return false
	}
}

// Check whether s is strictly equal to other
func (s String) StrictEqual(other Value) Value {
	return s.Equal(other)
}

// Get an element under the given index.
func (s String) Get(index int) (Char, Value) {
	var i int
	if index < 0 {
		l := s.CharCount()
		i = l + index
		if i < 0 {
			return 0, Ref(NewIndexOutOfRangeError(fmt.Sprint(index), l))
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
			return 0, Ref(NewIndexOutOfRangeError(fmt.Sprint(index), j))
		}
		if result == utf8.RuneError && size == 1 {
			// invalid UTF-8 character
			result = rune(leftStr[0])
		}
		if j == i {
			// found the character
			return Char(result), Undefined
		}
		leftStr = leftStr[size:]
		j++
	}
}

// Get the character under the given index.
func (s String) Subscript(key Value) (Char, Value) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return 0, Ref(NewIndexOutOfRangeError(key.Inspect(), len(s)))
		}
		return 0, Ref(NewCoerceError(IntClass, key.Class()))
	}

	return s.Get(i)
}

// Get the byte under the given index.
func (s String) ByteAtInt(index int) (UInt8, Value) {
	l := len(s)
	if index >= l || index < -l {
		return 0, Ref(NewIndexOutOfRangeError(fmt.Sprint(index), l))
	}
	if index < 0 {
		index = l + index
	}
	return UInt8(s[index]), Undefined
}

// Get the byte under the given index.
func (s String) ByteAt(key Value) (UInt8, Value) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return 0, Ref(NewIndexOutOfRangeError(key.Inspect(), len(s)))
		}
		return 0, Ref(NewCoerceError(IntClass, key.Class()))
	}

	return s.ByteAtInt(i)
}

// Get the grapheme under the given index.
func (s String) GraphemeAtInt(index int) (String, Value) {
	var i int
	if index < 0 {
		l := s.GraphemeCount()
		i = l + index
		if i < 0 {
			return "", Ref(NewIndexOutOfRangeError(fmt.Sprint(index), l))
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
			return String(cluster), Undefined
		}
		j++
	}

	return "", Ref(NewIndexOutOfRangeError(fmt.Sprint(index), j))
}

// Get the grapheme under the given index.
func (s String) GraphemeAt(key Value) (String, Value) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return "", Ref(NewIndexOutOfRangeError(key.Inspect(), len(s)))
		}
		return "", Ref(NewCoerceError(IntClass, key.Class()))
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

func (s *StringCharIterator) Copy() Reference {
	return &StringCharIterator{
		String:     s.String,
		ByteOffset: s.ByteOffset,
	}
}

func (s *StringCharIterator) Inspect() string {
	return fmt.Sprintf("Std::String::CharIterator{string: %s, byte_offset: %d}", s.String.Inspect(), s.ByteOffset)
}

func (s *StringCharIterator) Error() string {
	return s.Inspect()
}

func (*StringCharIterator) InstanceVariables() SymbolMap {
	return nil
}

func (s *StringCharIterator) Next() (Value, Value) {
	if s.ByteOffset >= len(s.String) {
		return Undefined, stopIterationSymbol.ToValue()
	}
	run, size := utf8.DecodeRuneInString(string(s.String[s.ByteOffset:]))

	s.ByteOffset += size
	return Char(run).ToValue(), Undefined
}

func (s *StringCharIterator) Reset() {
	s.ByteOffset = 0
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

func (s *StringByteIterator) Error() string {
	return s.Inspect()
}

func (s *StringByteIterator) Copy() Reference {
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
		return Undefined, stopIterationSymbol.ToValue()
	}
	result := UInt8(s.String[s.ByteOffset])
	s.ByteOffset += 1
	return result.ToValue(), Undefined
}

func (s *StringByteIterator) Reset() {
	s.ByteOffset = 0
}

type StringGraphemeIterator struct {
	String String
	Rest   string
	State  int
}

func NewStringGraphemeIterator(str String) *StringGraphemeIterator {
	return &StringGraphemeIterator{
		String: str,
		Rest:   string(str),
		State:  -1,
	}
}

func NewStringGraphemeIteratorWithRestAndState(str String, rest string, state int) *StringGraphemeIterator {
	return &StringGraphemeIterator{
		String: str,
		Rest:   rest,
		State:  state,
	}
}

func (*StringGraphemeIterator) Class() *Class {
	return StringGraphemeIteratorClass
}

func (*StringGraphemeIterator) DirectClass() *Class {
	return StringGraphemeIteratorClass
}

func (*StringGraphemeIterator) SingletonClass() *Class {
	return nil
}

func (s *StringGraphemeIterator) Error() string {
	return s.Inspect()
}

func (s *StringGraphemeIterator) Copy() Reference {
	return &StringGraphemeIterator{
		String: s.String,
		Rest:   s.Rest,
		State:  s.State,
	}
}

func (s *StringGraphemeIterator) Inspect() string {
	return fmt.Sprintf("Std::String::GraphemeIterator{string: %s}", s.String.Inspect())
}

func (*StringGraphemeIterator) InstanceVariables() SymbolMap {
	return nil
}

func (s *StringGraphemeIterator) Next() (Value, Value) {
	if len(s.Rest) == 0 {
		return Undefined, stopIterationSymbol.ToValue()
	}

	var grapheme string
	grapheme, s.Rest, _, s.State = uniseg.FirstGraphemeClusterInString(s.Rest, s.State)
	return Ref(String(grapheme)), Undefined
}

func (s *StringGraphemeIterator) Reset() {
	s.Rest = string(s.String)
	s.State = -1
}

func initString() {
	StringClass = NewClass()
	StdModule.AddConstantString("String", Ref(StringClass))

	StringCharIteratorClass = NewClass()
	StringClass.AddConstantString("CharIterator", Ref(StringCharIteratorClass))

	StringByteIteratorClass = NewClass()
	StringClass.AddConstantString("ByteIterator", Ref(StringByteIteratorClass))

	StringGraphemeIteratorClass = NewClass()
	StringClass.AddConstantString("GraphemeIterator", Ref(StringGraphemeIteratorClass))
}
