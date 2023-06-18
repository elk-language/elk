package object

import "fmt"

// Elk Value
type Value interface {
	Class() *Class  // Return the class of the object
	IsFrozen() bool // Whether the object is immutable
	SetFrozen()     // Freezes the object
}

// Return a string representation of the given object.
func Inspect(obj Value) string {
	switch o := obj.(type) {
	case SmallInt, Int64, UInt64, Int32, UInt32,
		Int16, UInt16, Int8, UInt8:
		return fmt.Sprintf("%d", o)
	case *Module:
		return o.Name
	case *Class:
		return o.Name
	case nil:
		return "undefined"
	default:
		return "<object>"
	}
}
