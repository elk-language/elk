// Package object contains definitions
// of Elk objects, classes, structs, modules etc.
package object

import "fmt"

// Elk Object
type Object interface {
	Class() *Class  // Return the class of the object
	IsFrozen() bool // Whether the object is immutable
	SetFrozen()     // Freezes the object
}

var ObjectClass *Class          // ::Std::Object
var PrimitiveObjectClass *Class // ::Std::PrimitiveObject

// Return a string representation of the given object.
func Inspect(obj Object) string {
	switch o := obj.(type) {
	case SmallInt:
		return fmt.Sprintf("%d", o)
	default:
		return "<object>"
	}
}
