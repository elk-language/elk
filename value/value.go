package value

// BENCHMARK: self-implemented tagged union
// Elk Value
type Value interface {
	Class() *Class                      // Return the class of the value
	IsFrozen() bool                     // Whether the value is immutable
	SetFrozen()                         // Freezes the value
	Inspect() string                    // Returns the string representation of the value
	InstanceVariables() SimpleSymbolMap // Returns the map of instance vars of this value, nil if value doesn't support instance vars
}

// Convert a Go bool value to Elk.
func ToElkBool(val bool) Bool {
	if val {
		return True
	}

	return False
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

// Add two values.
// When successful returns (result, nil, true).
// When an error occurred returns (nil, error, true).
// When there are no builtin addition functions for the given type returns (nil, nil, false).
func Add(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.Add(right)
	case *BigInt:
		result, err = l.Add(right)
	case Float:
		result, err = l.Add(right)
	case *BigFloat:
		result, err = l.Add(right)
	case Float64:
		result, err = StrictNumericAdd(l, right)
	case Float32:
		result, err = StrictNumericAdd(l, right)
	case Int64:
		result, err = StrictNumericAdd(l, right)
	case Int32:
		result, err = StrictNumericAdd(l, right)
	case Int16:
		result, err = StrictNumericAdd(l, right)
	case Int8:
		result, err = StrictNumericAdd(l, right)
	case UInt64:
		result, err = StrictNumericAdd(l, right)
	case UInt32:
		result, err = StrictNumericAdd(l, right)
	case UInt16:
		result, err = StrictNumericAdd(l, right)
	case UInt8:
		result, err = StrictNumericAdd(l, right)
	case String:
		result, err = l.Concat(right)
	case Char:
		result, err = l.Concat(right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Subtract two values
// When successful returns (result, nil, true).
// When an error occurred returns (nil, error, true).
// When there are no builtin addition functions for the given type returns (nil, nil, false).
func Subtract(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.Subtract(right)
	case *BigInt:
		result, err = l.Subtract(right)
	case Float:
		result, err = l.Subtract(right)
	case *BigFloat:
		result, err = l.Subtract(right)
	case String:
		result, err = l.RemoveSuffix(right)
	case Float64:
		result, err = StrictNumericSubtract(l, right)
	case Float32:
		result, err = StrictNumericSubtract(l, right)
	case Int64:
		result, err = StrictNumericSubtract(l, right)
	case Int32:
		result, err = StrictNumericSubtract(l, right)
	case Int16:
		result, err = StrictNumericSubtract(l, right)
	case Int8:
		result, err = StrictNumericSubtract(l, right)
	case UInt64:
		result, err = StrictNumericSubtract(l, right)
	case UInt32:
		result, err = StrictNumericSubtract(l, right)
	case UInt16:
		result, err = StrictNumericSubtract(l, right)
	case UInt8:
		result, err = StrictNumericSubtract(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Multiply two values
// When successful returns (result, nil, true).
// When an error occurred returns (nil, error, true).
// When there are no builtin addition functions for the given type returns (nil, nil, false).
func Multiply(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.Multiply(right)
	case *BigInt:
		result, err = l.Multiply(right)
	case Float:
		result, err = l.Multiply(right)
	case *BigFloat:
		result, err = l.Multiply(right)
	case Float64:
		result, err = StrictNumericMultiply(l, right)
	case Float32:
		result, err = StrictNumericMultiply(l, right)
	case Int64:
		result, err = StrictNumericMultiply(l, right)
	case Int32:
		result, err = StrictNumericMultiply(l, right)
	case Int16:
		result, err = StrictNumericMultiply(l, right)
	case Int8:
		result, err = StrictNumericMultiply(l, right)
	case UInt64:
		result, err = StrictNumericMultiply(l, right)
	case UInt32:
		result, err = StrictNumericMultiply(l, right)
	case UInt16:
		result, err = StrictNumericMultiply(l, right)
	case UInt8:
		result, err = StrictNumericMultiply(l, right)
	case String:
		result, err = l.Repeat(right)
	case Char:
		result, err = l.Repeat(right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Divide two values
// When successful returns (result, nil, true).
// When an error occurred returns (nil, error, true).
// When there are no builtin addition functions for the given type returns (nil, nil, false).
func Divide(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.Divide(right)
	case *BigInt:
		result, err = l.Divide(right)
	case Float:
		result, err = l.Divide(right)
	case *BigFloat:
		result, err = l.Divide(right)
	case Float64:
		result, err = StrictFloatDivide(l, right)
	case Float32:
		result, err = StrictFloatDivide(l, right)
	case Int64:
		result, err = StrictIntDivide(l, right)
	case Int32:
		result, err = StrictIntDivide(l, right)
	case Int16:
		result, err = StrictIntDivide(l, right)
	case Int8:
		result, err = StrictIntDivide(l, right)
	case UInt64:
		result, err = StrictIntDivide(l, right)
	case UInt32:
		result, err = StrictIntDivide(l, right)
	case UInt16:
		result, err = StrictIntDivide(l, right)
	case UInt8:
		result, err = StrictIntDivide(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Negate a value
// When successful returns (result, true).
// When there are no builtin negation functions for the given type returns (nil, false).
func Negate(operand Value) (Value, bool) {
	var result Value

	switch o := operand.(type) {
	case SmallInt:
		result = o.Negate()
	case *BigInt:
		result = o.Negate()
	case Float:
		result = -o
	case *BigFloat:
		result = o.Negate()
	case Float64:
		result = -o
	case Float32:
		result = -o
	case Int64:
		result = -o
	case Int32:
		result = -o
	case Int16:
		result = -o
	case Int8:
		result = -o
	case UInt64:
		result = -o
	case UInt32:
		result = -o
	case UInt16:
		result = -o
	case UInt8:
		result = -o
	default:
		return nil, false
	}

	return result, true
}

// Exponentiate two values
// When successful returns (result, nil, true).
// When an error occurred returns (nil, error, true).
// When there are no builtin addition functions for the given type returns (nil, nil, false).
func Exponentiate(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.Exponentiate(right)
	case *BigInt:
		result, err = l.Exponentiate(right)
	case Float:
		result, err = l.Exponentiate(right)
	case *BigFloat:
		result, err = l.Exponentiate(right)
	case Float64:
		result, err = StrictFloatExponentiate(l, right)
	case Float32:
		result, err = StrictFloatExponentiate(l, right)
	case Int64:
		result, err = StrictIntExponentiate(l, right)
	case Int32:
		result, err = StrictIntExponentiate(l, right)
	case Int16:
		result, err = StrictIntExponentiate(l, right)
	case Int8:
		result, err = StrictIntExponentiate(l, right)
	case UInt64:
		result, err = StrictIntExponentiate(l, right)
	case UInt32:
		result, err = StrictIntExponentiate(l, right)
	case UInt16:
		result, err = StrictIntExponentiate(l, right)
	case UInt8:
		result, err = StrictIntExponentiate(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Perform modulo on two values
// When successful returns (result, nil, true).
// When an error occurred returns (nil, error, true).
// When there are no builtin addition functions for the given type returns (nil, nil, false).
func Modulo(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.Modulo(right)
	case *BigInt:
		result, err = l.Modulo(right)
	case Float:
		result, err = l.Modulo(right)
	case *BigFloat:
		result, err = l.Modulo(right)
	case Float64:
		result, err = StrictFloatModulo(l, right)
	case Float32:
		result, err = StrictFloatModulo(l, right)
	case Int64:
		result, err = StrictIntModulo(l, right)
	case Int32:
		result, err = StrictIntModulo(l, right)
	case Int16:
		result, err = StrictIntModulo(l, right)
	case Int8:
		result, err = StrictIntModulo(l, right)
	case UInt64:
		result, err = StrictIntModulo(l, right)
	case UInt32:
		result, err = StrictIntModulo(l, right)
	case UInt16:
		result, err = StrictIntModulo(l, right)
	case UInt8:
		result, err = StrictIntModulo(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Check whether left is greater than right.
// When successful returns (result, nil, true).
// When an error occurred returns (nil, error, true).
// When there are no builtin addition functions for the given type returns (nil, nil, false).
func GreaterThan(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.GreaterThan(right)
	case *BigInt:
		result, err = l.GreaterThan(right)
	case Float:
		result, err = l.GreaterThan(right)
	case *BigFloat:
		result, err = l.GreaterThan(right)
	case String:
		result, err = l.GreaterThan(right)
	case Char:
		result, err = l.GreaterThan(right)
	case Float64:
		result, err = StrictNumericGreaterThan(l, right)
	case Float32:
		result, err = StrictNumericGreaterThan(l, right)
	case Int64:
		result, err = StrictNumericGreaterThan(l, right)
	case Int32:
		result, err = StrictNumericGreaterThan(l, right)
	case Int16:
		result, err = StrictNumericGreaterThan(l, right)
	case Int8:
		result, err = StrictNumericGreaterThan(l, right)
	case UInt64:
		result, err = StrictNumericGreaterThan(l, right)
	case UInt32:
		result, err = StrictNumericGreaterThan(l, right)
	case UInt16:
		result, err = StrictNumericGreaterThan(l, right)
	case UInt8:
		result, err = StrictNumericGreaterThan(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Check whether left is greater than or equal to right.
// When successful returns (result, nil, true).
// When an error occurred returns (nil, error, true).
// When there are no builtin addition functions for the given type returns (nil, nil, false).
func GreaterThanEqual(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.GreaterThanEqual(right)
	case *BigInt:
		result, err = l.GreaterThanEqual(right)
	case Float:
		result, err = l.GreaterThanEqual(right)
	case *BigFloat:
		result, err = l.GreaterThanEqual(right)
	case String:
		result, err = l.GreaterThanEqual(right)
	case Char:
		result, err = l.GreaterThanEqual(right)
	case Float64:
		result, err = StrictNumericGreaterThanEqual(l, right)
	case Float32:
		result, err = StrictNumericGreaterThanEqual(l, right)
	case Int64:
		result, err = StrictNumericGreaterThanEqual(l, right)
	case Int32:
		result, err = StrictNumericGreaterThanEqual(l, right)
	case Int16:
		result, err = StrictNumericGreaterThanEqual(l, right)
	case Int8:
		result, err = StrictNumericGreaterThanEqual(l, right)
	case UInt64:
		result, err = StrictNumericGreaterThanEqual(l, right)
	case UInt32:
		result, err = StrictNumericGreaterThanEqual(l, right)
	case UInt16:
		result, err = StrictNumericGreaterThanEqual(l, right)
	case UInt8:
		result, err = StrictNumericGreaterThanEqual(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Check whether left is less than right.
// When successful returns (result, nil, true).
// When an error occurred returns (nil, error, true).
// When there are no builtin addition functions for the given type returns (nil, nil, false).
func LessThan(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.LessThan(right)
	case *BigInt:
		result, err = l.LessThan(right)
	case Float:
		result, err = l.LessThan(right)
	case *BigFloat:
		result, err = l.LessThan(right)
	case String:
		result, err = l.LessThan(right)
	case Char:
		result, err = l.LessThan(right)
	case Float64:
		result, err = StrictNumericLessThan(l, right)
	case Float32:
		result, err = StrictNumericLessThan(l, right)
	case Int64:
		result, err = StrictNumericLessThan(l, right)
	case Int32:
		result, err = StrictNumericLessThan(l, right)
	case Int16:
		result, err = StrictNumericLessThan(l, right)
	case Int8:
		result, err = StrictNumericLessThan(l, right)
	case UInt64:
		result, err = StrictNumericLessThan(l, right)
	case UInt32:
		result, err = StrictNumericLessThan(l, right)
	case UInt16:
		result, err = StrictNumericLessThan(l, right)
	case UInt8:
		result, err = StrictNumericLessThan(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Check whether left is less than or equal to right.
// When successful returns (result, nil, true).
// When an error occurred returns (nil, error, true).
// When there are no builtin addition functions for the given type returns (nil, nil, false).
func LessThanEqual(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.LessThanEqual(right)
	case *BigInt:
		result, err = l.LessThanEqual(right)
	case Float:
		result, err = l.LessThanEqual(right)
	case *BigFloat:
		result, err = l.LessThanEqual(right)
	case String:
		result, err = l.LessThanEqual(right)
	case Char:
		result, err = l.LessThanEqual(right)
	case Float64:
		result, err = StrictNumericLessThanEqual(l, right)
	case Float32:
		result, err = StrictNumericLessThanEqual(l, right)
	case Int64:
		result, err = StrictNumericLessThanEqual(l, right)
	case Int32:
		result, err = StrictNumericLessThanEqual(l, right)
	case Int16:
		result, err = StrictNumericLessThanEqual(l, right)
	case Int8:
		result, err = StrictNumericLessThanEqual(l, right)
	case UInt64:
		result, err = StrictNumericLessThanEqual(l, right)
	case UInt32:
		result, err = StrictNumericLessThanEqual(l, right)
	case UInt16:
		result, err = StrictNumericLessThanEqual(l, right)
	case UInt8:
		result, err = StrictNumericLessThanEqual(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Check whether left is equal to right.
// When successful returns (result, true).
// When there are no builtin addition functions for the given type returns (nil, false).
func Equal(left, right Value) (Value, bool) {
	var result Value

	switch l := left.(type) {
	case SmallInt:
		result = l.Equal(right)
	case *BigInt:
		result = l.Equal(right)
	case Float:
		result = l.Equal(right)
	case *BigFloat:
		result = l.Equal(right)
	case String:
		result = l.Equal(right)
	case Char:
		result = l.Equal(right)
	case Float64:
		result = StrictFloatEqual(l, right)
	case Float32:
		result = StrictFloatEqual(l, right)
	case Int64:
		result = StrictSignedIntEqual(l, right)
	case Int32:
		result = StrictSignedIntEqual(l, right)
	case Int16:
		result = StrictSignedIntEqual(l, right)
	case Int8:
		result = StrictSignedIntEqual(l, right)
	case UInt64:
		result = StrictUnsignedIntEqual(l, right)
	case UInt32:
		result = StrictUnsignedIntEqual(l, right)
	case UInt16:
		result = StrictUnsignedIntEqual(l, right)
	case UInt8:
		result = StrictUnsignedIntEqual(l, right)
	default:
		return nil, false
	}

	return result, true
}

// Check whether left is strictly equal to right.
// When successful returns (result, true).
// When there are no builtin addition functions for the given type returns (nil, false).
func StrictEqual(left, right Value) (Value, bool) {
	var result Value

	switch l := left.(type) {
	case SmallInt:
		result = l.StrictEqual(right)
	case *BigInt:
		result = l.StrictEqual(right)
	case Float:
		result = l.StrictEqual(right)
	case *BigFloat:
		result = l.StrictEqual(right)
	case String:
		result = l.StrictEqual(right)
	case Char:
		result = l.StrictEqual(right)
	case Float64:
		result = StrictNumericStrictEqual(l, right)
	case Float32:
		result = StrictNumericStrictEqual(l, right)
	case Int64:
		result = StrictNumericStrictEqual(l, right)
	case Int32:
		result = StrictNumericStrictEqual(l, right)
	case Int16:
		result = StrictNumericStrictEqual(l, right)
	case Int8:
		result = StrictNumericStrictEqual(l, right)
	case UInt64:
		result = StrictNumericStrictEqual(l, right)
	case UInt32:
		result = StrictNumericStrictEqual(l, right)
	case UInt16:
		result = StrictNumericStrictEqual(l, right)
	case UInt8:
		result = StrictNumericStrictEqual(l, right)
	default:
		return nil, false
	}

	return result, true
}

// Execute a right bit shift >>.
// When successful returns (result, true).
// When there are no builtin negation functions for the given type returns (nil, false).
func RightBitshift(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.RightBitshift(right)
	case *BigInt:
		result, err = l.RightBitshift(right)
	case Int64:
		result, err = StrictIntRightBitshift(l, right)
	case Int32:
		result, err = StrictIntRightBitshift(l, right)
	case Int16:
		result, err = StrictIntRightBitshift(l, right)
	case Int8:
		result, err = StrictIntRightBitshift(l, right)
	case UInt64:
		result, err = StrictIntRightBitshift(l, right)
	case UInt32:
		result, err = StrictIntRightBitshift(l, right)
	case UInt16:
		result, err = StrictIntRightBitshift(l, right)
	case UInt8:
		result, err = StrictIntRightBitshift(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Execute a logical right bit shift >>>.
// When successful returns (result, true).
// When there are no builtin negation functions for the given type returns (nil, false).
func LogicalRightBitshift(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case Int64:
		result, err = StrictIntLogicalRightBitshift(l, right, logicalRightShift64)
	case Int32:
		result, err = StrictIntLogicalRightBitshift(l, right, logicalRightShift32)
	case Int16:
		result, err = StrictIntLogicalRightBitshift(l, right, logicalRightShift16)
	case Int8:
		result, err = StrictIntLogicalRightBitshift(l, right, logicalRightShift8)
	case UInt64:
		result, err = StrictIntRightBitshift(l, right)
	case UInt32:
		result, err = StrictIntRightBitshift(l, right)
	case UInt16:
		result, err = StrictIntRightBitshift(l, right)
	case UInt8:
		result, err = StrictIntRightBitshift(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Execute a left bit shift <<.
// When successful returns (result, true).
// When there are no builtin negation functions for the given type returns (nil, false).
func LeftBitshift(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.LeftBitshift(right)
	case *BigInt:
		result, err = l.LeftBitshift(right)
	case Int64:
		result, err = StrictIntLeftBitshift(l, right)
	case Int32:
		result, err = StrictIntLeftBitshift(l, right)
	case Int16:
		result, err = StrictIntLeftBitshift(l, right)
	case Int8:
		result, err = StrictIntLeftBitshift(l, right)
	case UInt64:
		result, err = StrictIntLeftBitshift(l, right)
	case UInt32:
		result, err = StrictIntLeftBitshift(l, right)
	case UInt16:
		result, err = StrictIntLeftBitshift(l, right)
	case UInt8:
		result, err = StrictIntLeftBitshift(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Execute a logical left bit shift <<<.
// When successful returns (result, true).
// When there are no builtin negation functions for the given type returns (nil, false).
func LogicalLeftBitshift(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case Int64:
		result, err = StrictIntLeftBitshift(l, right)
	case Int32:
		result, err = StrictIntLeftBitshift(l, right)
	case Int16:
		result, err = StrictIntLeftBitshift(l, right)
	case Int8:
		result, err = StrictIntLeftBitshift(l, right)
	case UInt64:
		result, err = StrictIntLeftBitshift(l, right)
	case UInt32:
		result, err = StrictIntLeftBitshift(l, right)
	case UInt16:
		result, err = StrictIntLeftBitshift(l, right)
	case UInt8:
		result, err = StrictIntLeftBitshift(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Execute a bitwise AND &.
// When successful returns (result, true).
// When there are no builtin negation functions for the given type returns (nil, false).
func BitwiseAnd(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.BitwiseAnd(right)
	case *BigInt:
		result, err = l.BitwiseAnd(right)
	case Int64:
		result, err = StrictIntBitwiseAnd(l, right)
	case Int32:
		result, err = StrictIntBitwiseAnd(l, right)
	case Int16:
		result, err = StrictIntBitwiseAnd(l, right)
	case Int8:
		result, err = StrictIntBitwiseAnd(l, right)
	case UInt64:
		result, err = StrictIntBitwiseAnd(l, right)
	case UInt32:
		result, err = StrictIntBitwiseAnd(l, right)
	case UInt16:
		result, err = StrictIntBitwiseAnd(l, right)
	case UInt8:
		result, err = StrictIntBitwiseAnd(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Execute a bitwise OR |.
// When successful returns (result, true).
// When there are no builtin negation functions for the given type returns (nil, false).
func BitwiseOr(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.BitwiseOr(right)
	case *BigInt:
		result, err = l.BitwiseOr(right)
	case Int64:
		result, err = StrictIntBitwiseOr(l, right)
	case Int32:
		result, err = StrictIntBitwiseOr(l, right)
	case Int16:
		result, err = StrictIntBitwiseOr(l, right)
	case Int8:
		result, err = StrictIntBitwiseOr(l, right)
	case UInt64:
		result, err = StrictIntBitwiseOr(l, right)
	case UInt32:
		result, err = StrictIntBitwiseOr(l, right)
	case UInt16:
		result, err = StrictIntBitwiseOr(l, right)
	case UInt8:
		result, err = StrictIntBitwiseOr(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}

// Execute a bitwise XOR ^.
// When successful returns (result, true).
// When there are no builtin negation functions for the given type returns (nil, false).
func BitwiseXor(left, right Value) (Value, *Error, bool) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.BitwiseXor(right)
	case *BigInt:
		result, err = l.BitwiseXor(right)
	case Int64:
		result, err = StrictIntBitwiseXor(l, right)
	case Int32:
		result, err = StrictIntBitwiseXor(l, right)
	case Int16:
		result, err = StrictIntBitwiseXor(l, right)
	case Int8:
		result, err = StrictIntBitwiseXor(l, right)
	case UInt64:
		result, err = StrictIntBitwiseXor(l, right)
	case UInt32:
		result, err = StrictIntBitwiseXor(l, right)
	case UInt16:
		result, err = StrictIntBitwiseXor(l, right)
	case UInt8:
		result, err = StrictIntBitwiseXor(l, right)
	default:
		return nil, nil, false
	}

	if err != nil {
		return nil, err, true
	}
	return result, nil, true
}
