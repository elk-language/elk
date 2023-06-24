package object

import (
	"fmt"
	"unicode/utf8"

	"github.com/rivo/uniseg"
)

var StringClass *Class // ::Std::String

// Elk's String value
type String string

func (s String) Class() *Class {
	return StringClass
}

func (s String) IsFrozen() bool {
	return true
}

func (s String) SetFrozen() {}

func (s String) Inspect() string {
	return fmt.Sprintf("%q", s)
}

func (s String) InstanceVariables() SimpleSymbolMap {
	return nil
}

// Returns the number of bytes
// present in the string.
func (s String) ByteLength() int {
	return len(s)
}

// Returns the number of unicode chars
// present in the string.
func (s String) CharLength() int {
	return utf8.RuneCountInString(string(s))
}

// Returns the number of grapheme clusters
// present in the string.
func (s String) GraphemeLength() int {
	return uniseg.GraphemeClusterCount(string(s))
}

// Concatenate another value with this string and return the result.
// If the operation is illegal an error will be returned.
func (s String) Concat(other Value) (String, *Error) {
	switch o := other.(type) {
	case String:
		return s + o, nil
	default:
		return "", Errorf(TypeErrorClass, "can't concat %s to string %s", other.Inspect(), s.Inspect())
	}
}

func initString() {
	StringClass = NewClass(
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("String", StringClass)
}
