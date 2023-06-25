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

// Returns true when the Elk value is truthy (works like true in boolean logic)
// otherwise returns false.
func Truthy(val Value) bool {
	switch val.(type) {
	case False, Nil:
		return false
	default:
		return true
	}
}

// Returns true when the Elk value is falsy (works like false in boolean logic)
// otherwise returns false.
func Falsy(val Value) bool {
	switch val.(type) {
	case False, Nil:
		return true
	default:
		return false
	}
}
