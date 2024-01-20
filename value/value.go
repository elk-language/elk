package value

import (
	"strings"
)

var ValueClass *Class // ::Std::Value

// BENCHMARK: self-implemented tagged union
// Elk Value
type Value interface {
	Class() *Class                // Return the class of the value
	DirectClass() *Class          // Return the direct class of this value that will be searched for methods first
	SingletonClass() *Class       // Return the singleton class of this value that holds methods unique to this object
	Inspect() string              // Returns the string representation of the value
	InstanceVariables() SymbolMap // Returns the map of instance vars of this value, nil if value doesn't support instance vars
	Copy() Value                  // Creates a shallow copy of the value. If the value is immutable, no copying should be done, the same value should be returned.
}

func IsMutableCollection(val Value) bool {
	switch val.(type) {
	case *ArrayList, *HashMap:
		return true
	}

	return false
}

// Return the string representation of a slice
// of values.
func InspectSlice[T Value](slice []T) string {
	var builder strings.Builder

	builder.WriteString("[")

	for i, element := range slice {
		if i != 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(element.Inspect())
	}

	builder.WriteString("]")
	return builder.String()
}

// Convert a pair of (Value, *Error) return values
// to (Value, Value)
func ToValueErr(val Value, err *Error) (Value, Value) {
	if err != nil {
		return nil, err
	}
	return val, nil
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

// Converts an Elk value strictly to Go int.
// Returns (0, false) when the value is incompatible.
// Returns (-1, false) when the value is a BigInt too large to be converted to int.
func IntToGoInt(val Value) (int, bool) {
	switch v := val.(type) {
	case SmallInt:
		return int(v), true
	case *BigInt:
		if !v.IsSmallInt() {
			return -1, false
		}
		return int(v.ToSmallInt()), true
	}

	return 0, false
}

// Converts an Elk value to Go int.
// Returns (0, false) when the value is incompatible.
// Returns (-1, false) when the value is a BigInt too large to be converted to int.
func ToGoInt(val Value) (int, bool) {
	switch v := val.(type) {
	case SmallInt:
		return int(v), true
	case *BigInt:
		if !v.IsSmallInt() {
			return -1, false
		}
		return int(v.ToSmallInt()), true
	case Int8:
		return int(v), true
	case Int16:
		return int(v), true
	case Int32:
		return int(v), true
	case Int64:
		return int(v), true
	case UInt8:
		return int(v), true
	case UInt16:
		return int(v), true
	case UInt32:
		return int(v), true
	case UInt64:
		return int(v), true
	}

	return 0, false
}

// Returns true when the Elk value is nil
// otherwise returns false.
func IsNil(val Value) bool {
	switch val.(type) {
	case NilType:
		return true
	default:
		return false
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

// Get an element by key.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Subscript(collection, key Value) (Value, *Error) {
	var result Value
	var err *Error

	switch l := collection.(type) {
	case *ArrayTuple:
		result, err = l.Subscript(key)
	case *ArrayList:
		result, err = l.Subscript(key)
	default:
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Set an element under the given key.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func SubscriptSet(collection, key, val Value) (Value, *Error) {
	var err *Error

	switch l := collection.(type) {
	case *ArrayList:
		err = l.SubscriptSet(key, val)
	default:
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return val, nil
}

// Calculate the hash of the value.
// When successful returns (result, nil).
// When an error occurred returns (0, error).
// When there are no builtin addition functions for the given type returns (0, NotBuiltinError).
func Hash(key Value) (UInt64, *Error) {
	var result UInt64
	var err *Error

	switch k := key.(type) {
	case String:
		result = k.Hash()
	case Symbol:
		result = k.Hash()
	case SmallInt:
		result = k.Hash()
	case *BigInt:
		result = k.Hash()
	case Float:
		result = k.Hash()
	case NilType:
		result = k.Hash()
	case Float64:
		result = k.Hash()
	case Float32:
		result = k.Hash()
	case Int64:
		result = k.Hash()
	case Int32:
		result = k.Hash()
	case Int16:
		result = k.Hash()
	case Int8:
		result = k.Hash()
	case UInt64:
		result = k.Hash()
	case UInt32:
		result = k.Hash()
	case UInt16:
		result = k.Hash()
	case UInt8:
		result = k.Hash()
	default:
		return 0, NotBuiltinError
	}

	if err != nil {
		return 0, err
	}
	return result, nil
}

// Add two values.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Add(left, right Value) (Value, *Error) {
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
	case *ArrayList:
		result, err = l.Concat(right)
	case *ArrayTuple:
		result, err = l.Concat(right)
	default:
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Subtract two values
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Subtract(left, right Value) (Value, *Error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Multiply two values
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Multiply(left, right Value) (Value, *Error) {
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
	case *ArrayList:
		result, err = l.Repeat(right)
	case *ArrayTuple:
		result, err = l.Repeat(right)
	default:
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Divide two values
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Divide(left, right Value) (Value, *Error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Negate a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func Negate(operand Value) Value {
	switch o := operand.(type) {
	case SmallInt:
		return o.Negate()
	case *BigInt:
		return o.Negate()
	case Float:
		return -o
	case *BigFloat:
		return o.Negate()
	case Float64:
		return -o
	case Float32:
		return -o
	case Int64:
		return -o
	case Int32:
		return -o
	case Int16:
		return -o
	case Int8:
		return -o
	case UInt64:
		return -o
	case UInt32:
		return -o
	case UInt16:
		return -o
	case UInt8:
		return -o
	default:
		return nil
	}
}

// Exponentiate two values
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Exponentiate(left, right Value) (Value, *Error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Perform modulo on two values
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Modulo(left, right Value) (Value, *Error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Compare two values.
// Returns 1 if left is greater than right.
// Returns 0 if both are equal.
// Returns -1 if left is less than right.
//
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Compare(left, right Value) (Value, *Error) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case SmallInt:
		result, err = l.Compare(right)
	case *BigInt:
		result, err = l.Compare(right)
	case Float:
		result, err = l.Compare(right)
	case *BigFloat:
		result, err = l.Compare(right)
	case String:
		result, err = l.Compare(right)
	case Char:
		result, err = l.Compare(right)
	case Float64:
		result, err = StrictFloatCompare(l, right)
	case Float32:
		result, err = StrictFloatCompare(l, right)
	case Int64:
		result, err = StrictIntCompare(l, right)
	case Int32:
		result, err = StrictIntCompare(l, right)
	case Int16:
		result, err = StrictIntCompare(l, right)
	case Int8:
		result, err = StrictIntCompare(l, right)
	case UInt64:
		result, err = StrictIntCompare(l, right)
	case UInt32:
		result, err = StrictIntCompare(l, right)
	case UInt16:
		result, err = StrictIntCompare(l, right)
	case UInt8:
		result, err = StrictIntCompare(l, right)
	default:
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Check whether left is greater than right.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func GreaterThan(left, right Value) (Value, *Error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Check whether left is greater than or equal to right.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func GreaterThanEqual(left, right Value) (Value, *Error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Check whether left is less than right.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func LessThan(left, right Value) (Value, *Error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Check whether left is less than or equal to right.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func LessThanEqual(left, right Value) (Value, *Error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Check whether left is equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func Equal(left, right Value) Value {
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
		return nil
	}

	return result
}

// Check whether left is not equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func NotEqual(left, right Value) Value {
	val := Equal(left, right)
	if val == nil {
		return nil
	}

	return ToNotBool(val)
}

// Check whether left is strictly equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func StrictEqual(left, right Value) Value {
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
	case Symbol:
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
		return nil
	}

	return result
}

// Check whether left is strictly not equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func StrictNotEqual(left, right Value) Value {
	val := StrictEqual(left, right)
	if val == nil {
		return nil
	}

	return ToNotBool(val)
}

// Execute a right bit shift >>.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func RightBitshift(left, right Value) (Value, *Error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Execute a logical right bit shift >>>.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func LogicalRightBitshift(left, right Value) (Value, *Error) {
	var result Value
	var err *Error

	switch l := left.(type) {
	case Int64:
		result, err = StrictIntLogicalRightBitshift(l, right, LogicalRightShift64)
	case Int32:
		result, err = StrictIntLogicalRightBitshift(l, right, LogicalRightShift32)
	case Int16:
		result, err = StrictIntLogicalRightBitshift(l, right, LogicalRightShift16)
	case Int8:
		result, err = StrictIntLogicalRightBitshift(l, right, LogicalRightShift8)
	case UInt64:
		result, err = StrictIntRightBitshift(l, right)
	case UInt32:
		result, err = StrictIntRightBitshift(l, right)
	case UInt16:
		result, err = StrictIntRightBitshift(l, right)
	case UInt8:
		result, err = StrictIntRightBitshift(l, right)
	default:
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Execute a left bit shift <<.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func LeftBitshift(left, right Value) (Value, *Error) {
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
	case *ArrayList:
		l.Append(right)
		result = l
	default:
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Execute a logical left bit shift <<<.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func LogicalLeftBitshift(left, right Value) (Value, *Error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Execute a bitwise AND &.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func BitwiseAnd(left, right Value) (Value, *Error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Execute a bitwise OR |.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func BitwiseOr(left, right Value) (Value, *Error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Execute a bitwise XOR ^.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func BitwiseXor(left, right Value) (Value, *Error) {
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}
