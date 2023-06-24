package object

// Strict numerics are sized and can't be automatically coerced
// to other types.
type StrictNumeric interface {
	Float64 | Float32 | Int64 | Int32 | Int16 | Int8 | UInt64 | UInt32 | UInt16 | UInt8
	Value
}

// Add a strict numeric to another value and return the result.
// If the operation is illegal an error will be returned.
func StrictNumericAdd[T StrictNumeric](left T, right Value) (T, *Error) {
	r, ok := right.(T)
	if !ok {
		return 0, Errorf(TypeErrorClass, "%s can't be coerced into %s", right.Class().PrintableName(), left.Class().PrintableName())
	}

	return left + r, nil
}
