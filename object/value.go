package object

// BENCHMARK: self-implemented tagged union
// Elk Value
type Value interface {
	Class() *Class                      // Return the class of the object
	IsFrozen() bool                     // Whether the object is immutable
	SetFrozen()                         // Freezes the object
	Inspect() string                    // Returns the string representation of the value
	InstanceVariables() SimpleSymbolMap // Returns the map of instance vars of this object, nil if object doesn't support instance vars
}

// Converts an Elk Value to an Elk Bool.
func ToBool(val Value) Bool {
	switch val.(type) {
	case FalseType, NilType:
		return False
	default:
		return True
	}
}

// Converts an Elk Value to an Elk Bool
// and negates it.
func ToNotBool(val Value) Bool {
	switch val.(type) {
	case FalseType, NilType:
		return True
	default:
		return False
	}
}

// Returns true when the Elk value is truthy (works like true in boolean logic)
// otherwise returns false.
func Truthy(val Value) bool {
	switch val.(type) {
	case FalseType, NilType:
		return false
	default:
		return true
	}
}

// Returns true when the Elk value is falsy (works like false in boolean logic)
// otherwise returns false.
func Falsy(val Value) bool {
	switch val.(type) {
	case FalseType, NilType:
		return true
	default:
		return false
	}
}
