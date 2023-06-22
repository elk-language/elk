package object

// BENCHMARK: self-implemented tagged union
// Elk Value
type Value interface {
	Class() *Class   // Return the class of the object
	IsFrozen() bool  // Whether the object is immutable
	SetFrozen()      // Freezes the object
	Inspect() string // Returns the string representation of the value
}
