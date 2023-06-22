package object

import "fmt"

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

// Concatenate another value with this string and return the result.
// If the operation is illegal an error will be returned.
func (s String) Concat(other Value) (String, error) {
	switch o := other.(type) {
	case String:
		return s + o, nil
	default:
		return "", fmt.Errorf("can't concat %s to string %s", other.Inspect(), s.Inspect())
	}
}

func initString() {
	StringClass = NewClass(ClassWithImmutable(), ClassWithSealed())
	StdModule.AddConstant("String", StringClass)
}
