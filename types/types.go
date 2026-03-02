// Package types contains the Elk type checker.
package types

// Get rid of nil or untyped
func Normalise(t Type) Type {
	if t == nil || IsUntyped(t) {
		return Void{}
	}

	return t
}
