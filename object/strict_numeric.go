package object

import (
	"fmt"
)

// Strict numerics are sized and can't be automatically coerced
// to other types.
type StrictNumeric interface {
	Float64 | Float32 | Int64 | Int32 | Int16 | Int8 | UInt64 | UInt32 | UInt16 | UInt8
	Value
}

// Add a strict numeric to another value and return the result.
// If the operation is illegal an error will be returned.
func StrictNumericAdd[T StrictNumeric](left T, right Value) (T, error) {
	r, ok := right.(T)
	if !ok {
		return 0, fmt.Errorf("can't add %s to %s", left.Inspect(), right.Inspect())
	}

	return left + r, nil
}
