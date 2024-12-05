package value

import (
	"fmt"
	"math"
	"strings"
	"unsafe"
)

var ValueClass *Class // ::Std::Value

type Value struct {
	tab  uintptr
	data unsafe.Pointer
}

const (
	NIL_FLAG = iota
	TRUE_FLAG
	FALSE_FLAG
	UNDEFINED_FLAG
	SMALL_INT_FLAG
	FLOAT_FLAG
	FLOAT32_FLAG
	INT8_FLAG
	UINT8_FLAG
	INT16_FLAG
	UINT16_FLAG
	INT32_FLAG
	UINT32_FLAG
	CHAR_FLAG
	SYMBOL_FLAG

	// only 64 bit systems
	INT64_FLAG
	UINT64_FLAG
	FLOAT64_FLAG
	DURATION_FLAG
	INLINE_VALUE_FLAG
)

// Convert a Reference to a Value
func Ref(ref Reference) Value {
	return *(*Value)(unsafe.Pointer(&ref))
}

func (v Value) Inspect() string {
	if v.IsReference() {
		return v.AsReference().Inspect()
	}

	switch v.ValueFlag() {
	case TRUE_FLAG:
		return v.AsTrue().Inspect()
	case FALSE_FLAG:
		return v.AsFalse().Inspect()
	case NIL_FLAG:
		return v.AsNil().Inspect()
	case UNDEFINED_FLAG:
		return v.AsUndefined().Inspect()
	case SMALL_INT_FLAG:
		return v.AsSmallInt().Inspect()
	case FLOAT_FLAG:
		return v.AsFloat().Inspect()
	case SYMBOL_FLAG:
		return v.AsSymbol().Inspect()
	case FLOAT64_FLAG:
		return v.AsFloat64().Inspect()
	case FLOAT32_FLAG:
		return v.AsFloat32().Inspect()
	case INT8_FLAG:
		return v.AsInt8().Inspect()
	case INT16_FLAG:
		return v.AsInt16().Inspect()
	case INT32_FLAG:
		return v.AsInt32().Inspect()
	case INT64_FLAG:
		return v.AsInt64().Inspect()
	case UINT8_FLAG:
		return v.AsUInt8().Inspect()
	case UINT16_FLAG:
		return v.AsUInt16().Inspect()
	case UINT32_FLAG:
		return v.AsUInt32().Inspect()
	case UINT64_FLAG:
		return v.AsUInt64().Inspect()
	case CHAR_FLAG:
		return v.AsChar().Inspect()
	case DURATION_FLAG:
		return v.AsDuration().Inspect()
	default:
		panic(fmt.Sprintf("invalid inline value flag: %d", v.ValueFlag()))
	}
}

func (v Value) Copy() Value {
	if v.IsReference() {
		return Ref(v.AsReference().Copy())
	}

	return v
}

func (v Value) Class() *Class {
	if v.IsReference() {
		return v.AsReference().Class()
	}

	switch v.ValueFlag() {
	case TRUE_FLAG:
		return v.AsTrue().Class()
	case FALSE_FLAG:
		return v.AsFalse().Class()
	case NIL_FLAG:
		return v.AsNil().Class()
	case UNDEFINED_FLAG:
		return v.AsUndefined().Class()
	case SMALL_INT_FLAG:
		return v.AsSmallInt().Class()
	case FLOAT_FLAG:
		return v.AsFloat().Class()
	case SYMBOL_FLAG:
		return v.AsSymbol().Class()
	case FLOAT64_FLAG:
		return v.AsFloat64().Class()
	case FLOAT32_FLAG:
		return v.AsFloat32().Class()
	case INT8_FLAG:
		return v.AsInt8().Class()
	case INT16_FLAG:
		return v.AsInt16().Class()
	case INT32_FLAG:
		return v.AsInt32().Class()
	case INT64_FLAG:
		return v.AsInt64().Class()
	case UINT8_FLAG:
		return v.AsUInt8().Class()
	case UINT16_FLAG:
		return v.AsUInt16().Class()
	case UINT32_FLAG:
		return v.AsUInt32().Class()
	case UINT64_FLAG:
		return v.AsUInt64().Class()
	case CHAR_FLAG:
		return v.AsChar().Class()
	case DURATION_FLAG:
		return v.AsDuration().Class()
	default:
		panic(fmt.Sprintf("invalid inline value flag: %d", v.ValueFlag()))
	}
}

func (v Value) DirectClass() *Class {
	if v.IsReference() {
		return v.AsReference().DirectClass()
	}

	switch v.ValueFlag() {
	case TRUE_FLAG:
		return v.AsTrue().DirectClass()
	case FALSE_FLAG:
		return v.AsFalse().DirectClass()
	case NIL_FLAG:
		return v.AsNil().DirectClass()
	case UNDEFINED_FLAG:
		return v.AsUndefined().DirectClass()
	case SMALL_INT_FLAG:
		return v.AsSmallInt().DirectClass()
	case FLOAT_FLAG:
		return v.AsFloat().DirectClass()
	case SYMBOL_FLAG:
		return v.AsSymbol().DirectClass()
	case FLOAT64_FLAG:
		return v.AsFloat64().DirectClass()
	case FLOAT32_FLAG:
		return v.AsFloat32().DirectClass()
	case INT8_FLAG:
		return v.AsInt8().DirectClass()
	case INT16_FLAG:
		return v.AsInt16().DirectClass()
	case INT32_FLAG:
		return v.AsInt32().DirectClass()
	case INT64_FLAG:
		return v.AsInt64().DirectClass()
	case UINT8_FLAG:
		return v.AsUInt8().DirectClass()
	case UINT16_FLAG:
		return v.AsUInt16().DirectClass()
	case UINT32_FLAG:
		return v.AsUInt32().DirectClass()
	case UINT64_FLAG:
		return v.AsUInt64().DirectClass()
	case CHAR_FLAG:
		return v.AsChar().DirectClass()
	case DURATION_FLAG:
		return v.AsDuration().DirectClass()
	default:
		panic(fmt.Sprintf("invalid inline value flag: %d", v.ValueFlag()))
	}
}

func (v Value) SingletonClass() *Class {
	if v.IsReference() {
		return v.AsReference().SingletonClass()
	}

	switch v.ValueFlag() {
	case TRUE_FLAG:
		return v.AsTrue().SingletonClass()
	case FALSE_FLAG:
		return v.AsFalse().SingletonClass()
	case NIL_FLAG:
		return v.AsNil().SingletonClass()
	case UNDEFINED_FLAG:
		return v.AsUndefined().SingletonClass()
	case SMALL_INT_FLAG:
		return v.AsSmallInt().SingletonClass()
	case FLOAT_FLAG:
		return v.AsFloat().SingletonClass()
	case SYMBOL_FLAG:
		return v.AsSymbol().SingletonClass()
	case FLOAT64_FLAG:
		return v.AsFloat64().SingletonClass()
	case FLOAT32_FLAG:
		return v.AsFloat32().SingletonClass()
	case INT8_FLAG:
		return v.AsInt8().SingletonClass()
	case INT16_FLAG:
		return v.AsInt16().SingletonClass()
	case INT32_FLAG:
		return v.AsInt32().SingletonClass()
	case INT64_FLAG:
		return v.AsInt64().SingletonClass()
	case UINT8_FLAG:
		return v.AsUInt8().SingletonClass()
	case UINT16_FLAG:
		return v.AsUInt16().SingletonClass()
	case UINT32_FLAG:
		return v.AsUInt32().SingletonClass()
	case UINT64_FLAG:
		return v.AsUInt64().SingletonClass()
	case CHAR_FLAG:
		return v.AsChar().SingletonClass()
	case DURATION_FLAG:
		return v.AsDuration().SingletonClass()
	default:
		panic(fmt.Sprintf("invalid inline value flag: %d", v.ValueFlag()))
	}
}

func (v Value) InstanceVariables() SymbolMap {
	if v.IsReference() {
		return v.AsReference().InstanceVariables()
	}

	switch v.ValueFlag() {
	case TRUE_FLAG:
		return v.AsTrue().InstanceVariables()
	case FALSE_FLAG:
		return v.AsFalse().InstanceVariables()
	case NIL_FLAG:
		return v.AsNil().InstanceVariables()
	case UNDEFINED_FLAG:
		return v.AsUndefined().InstanceVariables()
	case SMALL_INT_FLAG:
		return v.AsSmallInt().InstanceVariables()
	case FLOAT_FLAG:
		return v.AsFloat().InstanceVariables()
	case SYMBOL_FLAG:
		return v.AsSymbol().InstanceVariables()
	case FLOAT64_FLAG:
		return v.AsFloat64().InstanceVariables()
	case FLOAT32_FLAG:
		return v.AsFloat32().InstanceVariables()
	case INT8_FLAG:
		return v.AsInt8().InstanceVariables()
	case INT16_FLAG:
		return v.AsInt16().InstanceVariables()
	case INT32_FLAG:
		return v.AsInt32().InstanceVariables()
	case INT64_FLAG:
		return v.AsInt64().InstanceVariables()
	case UINT8_FLAG:
		return v.AsUInt8().InstanceVariables()
	case UINT16_FLAG:
		return v.AsUInt16().InstanceVariables()
	case UINT32_FLAG:
		return v.AsUInt32().InstanceVariables()
	case UINT64_FLAG:
		return v.AsUInt64().InstanceVariables()
	case CHAR_FLAG:
		return v.AsChar().InstanceVariables()
	case DURATION_FLAG:
		return v.AsDuration().InstanceVariables()
	default:
		panic(fmt.Sprintf("invalid inline value flag: %d", v.ValueFlag()))
	}
}

func (v Value) Error() string {
	if v.IsReference() {
		return v.AsReference().Error()
	}

	switch v.ValueFlag() {
	case TRUE_FLAG:
		return v.AsTrue().Error()
	case FALSE_FLAG:
		return v.AsFalse().Error()
	case NIL_FLAG:
		return v.AsNil().Error()
	case UNDEFINED_FLAG:
		return v.AsUndefined().Error()
	case SMALL_INT_FLAG:
		return v.AsSmallInt().Error()
	case FLOAT_FLAG:
		return v.AsFloat().Error()
	case SYMBOL_FLAG:
		return v.AsSymbol().Error()
	case FLOAT64_FLAG:
		return v.AsFloat64().Error()
	case FLOAT32_FLAG:
		return v.AsFloat32().Error()
	case INT8_FLAG:
		return v.AsInt8().Error()
	case INT16_FLAG:
		return v.AsInt16().Error()
	case INT32_FLAG:
		return v.AsInt32().Error()
	case INT64_FLAG:
		return v.AsInt64().Error()
	case UINT8_FLAG:
		return v.AsUInt8().Error()
	case UINT16_FLAG:
		return v.AsUInt16().Error()
	case UINT32_FLAG:
		return v.AsUInt32().Error()
	case UINT64_FLAG:
		return v.AsUInt64().Error()
	case CHAR_FLAG:
		return v.AsChar().Error()
	case DURATION_FLAG:
		return v.AsDuration().Error()
	default:
		panic(fmt.Sprintf("invalid inline value flag: %d", v.ValueFlag()))
	}
}

func (v Value) ValueFlag() uint64 {
	return uint64(uintptr(v.data))
}

func (v Value) IsReference() bool {
	return uintptr(v.data) > INLINE_VALUE_FLAG
}

// Returns `nil` when the value is not a reference
func (v Value) SafeAsReference() Reference {
	if !v.IsReference() {
		return nil
	}
	return v.AsReference()
}

func (v Value) AsReference() Reference {
	return *(*Reference)(unsafe.Pointer(&v))
}

func (v Value) IsInlineValue() bool {
	return uintptr(v.data) < INLINE_VALUE_FLAG
}

func (v Value) IsSmallInt() bool {
	return uintptr(v.data) == SMALL_INT_FLAG
}

func (v Value) AsSmallInt() SmallInt {
	return *(*SmallInt)(unsafe.Pointer(&v.tab))
}

func (v Value) IsChar() bool {
	return uintptr(v.data) == CHAR_FLAG
}

func (v Value) AsChar() Char {
	return *(*Char)(unsafe.Pointer(&v.tab))
}

func (v Value) IsFloat() bool {
	return uintptr(v.data) == FLOAT_FLAG
}

func (v Value) AsFloat() Float {
	return *(*Float)(unsafe.Pointer(&v.tab))
}

func (v Value) IsFloat32() bool {
	return uintptr(v.data) == FLOAT32_FLAG
}

func (v Value) AsFloat32() Float32 {
	return *(*Float32)(unsafe.Pointer(&v.tab))
}

func (v Value) IsFloat64() bool {
	return uintptr(v.data) == FLOAT64_FLAG
}

func (v Value) AsFloat64() Float64 {
	return *(*Float64)(unsafe.Pointer(&v.tab))
}

func (v Value) IsInt8() bool {
	return uintptr(v.data) == INT8_FLAG
}

func (v Value) AsInt8() Int8 {
	return *(*Int8)(unsafe.Pointer(&v.tab))
}

func (v Value) IsUInt8() bool {
	return uintptr(v.data) == UINT8_FLAG
}

func (v Value) AsUInt8() UInt8 {
	return *(*UInt8)(unsafe.Pointer(&v.tab))
}

func (v Value) IsInt16() bool {
	return uintptr(v.data) == INT16_FLAG
}

func (v Value) AsInt16() Int16 {
	return *(*Int16)(unsafe.Pointer(&v.tab))
}

func (v Value) IsUInt16() bool {
	return uintptr(v.data) == UINT16_FLAG
}

func (v Value) AsUInt16() UInt16 {
	return *(*UInt16)(unsafe.Pointer(&v.tab))
}

func (v Value) IsInt32() bool {
	return uintptr(v.data) == INT32_FLAG
}

func (v Value) AsInt32() Int32 {
	return *(*Int32)(unsafe.Pointer(&v.tab))
}

func (v Value) IsUInt32() bool {
	return uintptr(v.data) == UINT32_FLAG
}

func (v Value) AsUInt32() UInt32 {
	return *(*UInt32)(unsafe.Pointer(&v.tab))
}

func (v Value) IsInt64() bool {
	return uintptr(v.data) == INT64_FLAG
}

func (v Value) AsInt64() Int64 {
	return *(*Int64)(unsafe.Pointer(&v.tab))
}

func (v Value) IsUInt64() bool {
	return uintptr(v.data) == UINT64_FLAG
}

func (v Value) AsUInt64() UInt64 {
	return *(*UInt64)(unsafe.Pointer(&v.tab))
}

func (v Value) IsDuration() bool {
	return uintptr(v.data) == DURATION_FLAG
}

func (v Value) AsDuration() Duration {
	return *(*Duration)(unsafe.Pointer(&v.tab))
}

func (v Value) IsTrue() bool {
	return uintptr(v.data) == TRUE_FLAG
}

func (v Value) AsTrue() TrueType {
	return *(*TrueType)(unsafe.Pointer(&v.tab))
}

func (v Value) IsSymbol() bool {
	return uintptr(v.data) == SYMBOL_FLAG
}

func (v Value) AsSymbol() Symbol {
	return *(*Symbol)(unsafe.Pointer(&v.tab))
}

func (v Value) IsFalse() bool {
	return uintptr(v.data) == FALSE_FLAG
}

func (v Value) AsFalse() FalseType {
	return *(*FalseType)(unsafe.Pointer(&v.tab))
}

func (v Value) IsNil() bool {
	return uintptr(v.data) == NIL_FLAG
}

func (v Value) AsNil() NilType {
	return *(*NilType)(unsafe.Pointer(&v.tab))
}

func (v Value) IsUndefined() bool {
	return uintptr(v.data) == UNDEFINED_FLAG
}

func (v Value) AsUndefined() UndefinedType {
	return *(*UndefinedType)(unsafe.Pointer(&v.tab))
}

// BENCHMARK: self-implemented tagged union
// Elk Value
type Reference interface {
	Class() *Class                // Return the class of the value
	DirectClass() *Class          // Return the direct class of this value that will be searched for methods first
	SingletonClass() *Class       // Return the singleton class of this value that holds methods unique to this object
	Inspect() string              // Returns the string representation of the value
	InstanceVariables() SymbolMap // Returns the map of instance vars of this value, nil if value doesn't support instance vars
	Copy() Reference              // Creates a shallow copy of the reference. If the value is immutable, no copying should be done, the same value should be returned.
	Error() string                // Implements the error interface
}

func IsMutableCollection(val Value) bool {
	if val.IsInlineValue() {
		return false
	}
	switch v := val.AsReference().(type) {
	case *ArrayList, *HashMap:
		return true
	case *ArrayTuple:
		for _, element := range *v {
			if IsMutableCollection(element) {
				return true
			}
		}
	case *HashRecord:
		for _, pair := range v.Table {
			if IsMutableCollection(pair.Key) || IsMutableCollection(pair.Value) {
				return true
			}
		}
	}

	return false
}

type Inspectable interface {
	Inspect() string
}

// Return the string representation of a slice
// of values.
func InspectSlice[T Inspectable](slice []T) string {
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

// Convert a Go bool value to Elk.
func ToElkBool(val bool) Value {
	if val {
		return True
	}

	return False
}

// Converts an Elk Value to an Elk Bool.
func ToBool(val Value) Value {
	if val.IsReference() {
		return True
	}

	switch val.ValueFlag() {
	case FALSE_FLAG, NIL_FLAG:
		return False
	default:
		return True
	}
}

// Converts an Elk Value to an Elk Bool
// and negates it.
func ToNotBool(val Value) Value {
	if val.IsReference() {
		return False
	}

	switch val.ValueFlag() {
	case FALSE_FLAG, NIL_FLAG:
		return True
	default:
		return False
	}
}

// Converts an Elk value strictly to Go int.
// Returns (0, false) when the value is incompatible.
// Returns (-1, false) when the value is a BigInt too large to be converted to int.
func IntToGoInt(val Value) (int, bool) {
	if val.IsReference() {
		switch v := val.AsReference().(type) {
		case *BigInt:
			if !v.IsSmallInt() {
				return -1, false
			}
			return int(v.ToSmallInt()), true
		}
		return 0, false
	}

	switch val.ValueFlag() {
	case SMALL_INT_FLAG:
		return int(val.AsSmallInt()), true
	}

	return 0, false
}

// Converts an Elk value to Go int.
// Returns (0, false) when the value is incompatible.
// Returns (-1, false) when the value is a BigInt too large to be converted to int.
func ToGoInt(val Value) (int, bool) {
	if val.IsReference() {
		switch v := val.AsReference().(type) {
		case *BigInt:
			if !v.IsSmallInt() {
				return -1, false
			}
			return int(v.ToSmallInt()), true
		case UInt64:
			return int(v), true
		}
		return 0, false
	}

	switch val.ValueFlag() {
	case SMALL_INT_FLAG:
		return int(val.AsSmallInt()), true
	case INT8_FLAG:
		return int(val.AsInt8()), true
	case INT16_FLAG:
		return int(val.AsInt16()), true
	case INT32_FLAG:
		return int(val.AsInt32()), true
	case INT64_FLAG:
		return int(val.AsInt64()), true
	case UINT8_FLAG:
		return int(val.AsUInt8()), true
	case UINT16_FLAG:
		return int(val.AsUInt16()), true
	case UINT32_FLAG:
		return int(val.AsUInt32()), true
	case UINT64_FLAG:
		return int(val.AsUInt64()), true
	}
	return 0, false
}

// Converts an Elk value to Go uint.
// Returns (0, false) when the value is incompatible, too large or negative.
func ToGoUInt(val Value) (uint, bool) {
	if val.IsReference() {
		switch v := val.AsReference().(type) {
		case *BigInt:
			if !v.IsSmallInt() {
				return 0, false
			}
			i := v.ToSmallInt()
			if i < 0 {
				return 0, false
			}
			if uint64(i) > math.MaxUint {
				return 0, false
			}
			return uint(i), true
		case Int64:
			if v < 0 {
				return 0, false
			}
			if uint64(v) > math.MaxUint {
				return 0, false
			}
			return uint(v), true
		case UInt64:
			if uint64(v) > math.MaxUint {
				return 0, false
			}
			return uint(v), true
		}
		return 0, false
	}

	switch val.ValueFlag() {
	case SMALL_INT_FLAG:
		v := val.AsSmallInt()
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	case INT8_FLAG:
		v := val.AsInt8()
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	case INT16_FLAG:
		v := val.AsInt16()
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	case INT32_FLAG:
		v := val.AsInt32()
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	case INT64_FLAG:
		v := val.AsInt64()
		if v < 0 {
			return 0, false
		}
		if uint64(v) > math.MaxUint {
			return 0, false
		}
		return uint(v), true
	case UINT8_FLAG:
		v := val.AsUInt8()
		return uint(v), true
	case UINT16_FLAG:
		v := val.AsUInt16()
		return uint(v), true
	case UINT32_FLAG:
		v := val.AsUInt32()
		return uint(v), true
	case UINT64_FLAG:
		v := val.AsUInt64()
		if uint64(v) > math.MaxUint {
			return 0, false
		}
		return uint(v), true
	}

	return 0, false
}

// Returns true when the Elk value is nil
// otherwise returns false.
func IsNil(val Value) bool {
	return val.IsNil()
}

// Returns true when the Elk value is truthy (works like true in boolean logic)
// otherwise returns false.
func Truthy(val Value) bool {
	if val.IsReference() {
		return true
	}
	switch val.ValueFlag() {
	case FALSE_FLAG, NIL_FLAG, UNDEFINED_FLAG:
		return false
	default:
		return true
	}
}

// Returns true when the Elk value is falsy (works like false in boolean logic)
// otherwise returns false.
func Falsy(val Value) bool {
	if val.IsReference() {
		return false
	}
	switch val.ValueFlag() {
	case FALSE_FLAG, NIL_FLAG, UNDEFINED_FLAG:
		return true
	default:
		return false
	}
}

// Check if the given value is an instance of the given class.
func InstanceOf(val Value, class *Class) bool {
	return class == val.Class()
}

func IsA(val Value, class *Class) bool {
	if class.IsMixin() {
		return mixinIsA(val, class)
	}

	return classIsA(val, class)
}

// Check if the given value is an instance of the given class or its subclasses.
func classIsA(val Value, class *Class) bool {
	currentClass := val.Class()
	for currentClass != nil {
		if currentClass == class {
			return true
		}

		currentClass = currentClass.Superclass()
	}

	return false
}

// Check if the given value is an instance of the classes that mix in the given mixin.
func mixinIsA(val Value, mixin *Mixin) bool {
	class := val.DirectClass()

	for parent := range class.Parents() {
		if parent == mixin {
			return true
		}
	}

	return false
}

// Get an element by key.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Subscript(collection, key Value) (result, err Value) {
	if !collection.IsReference() {
		return Nil, Nil
	}

	switch l := collection.AsReference().(type) {
	case *ArrayTuple:
		result, err = l.Subscript(key)
	case *ArrayList:
		result, err = l.Subscript(key)
	default:
		return Nil, Nil
	}

	if !err.IsNil() {
		return Nil, err
	}
	return result, Nil
}

// Set an element under the given key.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func SubscriptSet(collection, key, val Value) (result, err Value) {
	if !collection.IsReference() {
		return Nil, Nil
	}

	switch l := collection.AsReference().(type) {
	case *ArrayList:
		err = l.SubscriptSet(key, val)
	default:
		return Nil, Nil
	}

	if !err.IsNil() {
		return Nil, err
	}
	return val, Nil
}

// Calculate the hash of the value.
// When successful returns (result, nil).
// When an error occurred returns (0, error).
// When there are no builtin addition functions for the given type returns (0, NotBuiltinError).
func Hash(key Value) (UInt64, Value) {
	var result UInt64
	if key.IsReference() {
		switch k := key.AsReference().(type) {
		case String:
			result = k.Hash()
		case *BigInt:
			result = k.Hash()
		case *BigFloat:
			result = k.Hash()
		case Float64:
			result = k.Hash()
		case Int64:
			result = k.Hash()
		case UInt64:
			result = k.Hash()
		default:
			return 0, Ref(NotBuiltinError)
		}
		return result, Nil
	}

	switch key.ValueFlag() {
	case CHAR_FLAG:
		k := key.AsChar()
		result = k.Hash()
	case SYMBOL_FLAG:
		k := key.AsSymbol()
		result = k.Hash()
	case SMALL_INT_FLAG:
		k := key.AsSmallInt()
		result = k.Hash()
	case FLOAT_FLAG:
		k := key.AsFloat()
		result = k.Hash()
	case NIL_FLAG:
		k := key.AsNil()
		result = k.Hash()
	case TRUE_FLAG:
		k := key.AsTrue()
		result = k.Hash()
	case FALSE_FLAG:
		k := key.AsFalse()
		result = k.Hash()
	case FLOAT64_FLAG:
		k := key.AsFloat64()
		result = k.Hash()
	case FLOAT32_FLAG:
		k := key.AsFloat32()
		result = k.Hash()
	case INT64_FLAG:
		k := key.AsInt64()
		result = k.Hash()
	case INT32_FLAG:
		k := key.AsInt32()
		result = k.Hash()
	case INT16_FLAG:
		k := key.AsInt16()
		result = k.Hash()
	case INT8_FLAG:
		k := key.AsInt8()
		result = k.Hash()
	case UINT64_FLAG:
		k := key.AsUInt64()
		result = k.Hash()
	case UINT32_FLAG:
		k := key.AsUInt32()
		result = k.Hash()
	case UINT16_FLAG:
		k := key.AsUInt16()
		result = k.Hash()
	case UINT8_FLAG:
		k := key.AsUInt8()
		result = k.Hash()
	default:
		return 0, Ref(NotBuiltinError)
	}
	return result, Nil
}

// Add two values.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Add(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			result, err = l.Add(right)
		case *BigFloat:
			result, err = l.Add(right)
		case Float64:
			var r Float64
			r, err = l.Add(right)
			result = Ref(r)
		case Int64:
			var r Int64
			r, err = l.Add(right)
			result = Ref(r)
		case UInt64:
			var r UInt64
			r, err = l.Add(right)
			result = Ref(r)
		case String:
			var r String
			r, err = l.Concat(right)
			result = Ref(r)
		case *Regex:
			result, err = l.Concat(right)
		case *ArrayList:
			var r *ArrayList
			r, err = l.Concat(right)
			result = Ref(r)
		case *ArrayTuple:
			result, err = l.Concat(right)
		default:
			return Nil, Nil
		}

		if !err.IsNil() {
			return Nil, err
		}
		return result, Nil
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		result, err = l.Add(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		result, err = l.Add(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		var r Float64
		r, err = l.Add(right)
		result = r.ToValue()
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		var r Float32
		r, err = l.Add(right)
		result = r.ToValue()
	case INT64_FLAG:
		l := left.AsInt64()
		var r Int64
		r, err = l.Add(right)
		result = r.ToValue()
	case INT32_FLAG:
		l := left.AsInt32()
		var r Int32
		r, err = l.Add(right)
		result = r.ToValue()
	case INT16_FLAG:
		l := left.AsInt16()
		var r Int16
		r, err = l.Add(right)
		result = r.ToValue()
	case INT8_FLAG:
		l := left.AsInt8()
		var r Int8
		r, err = l.Add(right)
		result = r.ToValue()
	case UINT64_FLAG:
		l := left.AsUInt64()
		var r UInt64
		r, err = l.Add(right)
		result = r.ToValue()
	case UINT32_FLAG:
		l := left.AsUInt32()
		var r UInt32
		r, err = l.Add(right)
		result = r.ToValue()
	case UINT16_FLAG:
		l := left.AsUInt16()
		var r UInt16
		r, err = l.Add(right)
		result = r.ToValue()
	case UINT8_FLAG:
		l := left.AsUInt8()
		var r UInt8
		r, err = l.Add(right)
		result = r.ToValue()
	case CHAR_FLAG:
		l := left.AsChar()
		var r String
		r, err = l.Concat(right)
		result = Ref(r)
	default:
		return Nil, Nil
	}

	if !err.IsNil() {
		return Nil, err
	}
	return result, Nil
}

// Subtract two values
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Subtract(left, right Value) (result, err Value) {
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
func Multiply(left, right Value) (result, err Value) {
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
		var r Float64
		r, err = l.Multiply(right)
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
	case *Regex:
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
func Divide(left, right Value) (result, err Value) {
	if right.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			result, err = l.Divide(right)
		case *BigFloat:
			result, err = l.Divide(right)
		case Float64:
			var r Float64
			r, err = l.Divide(right)
			result = r.ToValue()
		case Int64:
			var r Int64
			r, err = l.Divide(right)
			result = r.ToValue()
		case UInt64:
			var r UInt64
			r, err = l.Divide(right)
			result = r.ToValue()
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		result, err = l.Divide(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		result, err = l.Divide(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		var r Float64
		r, err = l.Divide(right)
		result = r.ToValue()
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		var r Float32
		r, err = l.Divide(right)
		result = r.ToValue()
	case INT64_FLAG:
		l := left.AsInt64()
		var r Int64
		r, err = l.Divide(right)
		result = r.ToValue()
	case INT32_FLAG:
		l := left.AsInt32()
		var r Int32
		r, err = l.Divide(right)
		result = r.ToValue()
	case INT16_FLAG:
		l := left.AsInt16()
		var r Int16
		r, err = l.Divide(right)
		result = r.ToValue()
	case INT8_FLAG:
		l := left.AsInt8()
		var r Int8
		r, err = l.Divide(right)
		result = r.ToValue()
	case UINT64_FLAG:
		l := left.AsUInt64()
		var r UInt64
		r, err = l.Divide(right)
		result = r.ToValue()
	case UINT32_FLAG:
		l := left.AsUInt32()
		var r UInt32
		r, err = l.Divide(right)
		result = r.ToValue()
	case UINT16_FLAG:
		l := left.AsUInt16()
		var r UInt16
		r, err = l.Divide(right)
		result = r.ToValue()
	case UINT8_FLAG:
		l := left.AsUInt8()
		var r UInt8
		r, err = l.Divide(right)
		result = r.ToValue()
	default:
		return Nil, Nil
	}

	if !err.IsNil() {
		return Nil, err
	}
	return result, Nil
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

// Increment a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func Increment(operand Value) Value {
	switch o := operand.(type) {
	case SmallInt:
		return o.Increment()
	case *BigInt:
		return o.Increment()
	case Char:
		return o + 1
	case Int64:
		return o + 1
	case Int32:
		return o + 1
	case Int16:
		return o + 1
	case Int8:
		return o + 1
	case UInt64:
		return o + 1
	case UInt32:
		return o + 1
	case UInt16:
		return o + 1
	case UInt8:
		return o + 1
	default:
		return nil
	}
}

// Decrement a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func Decrement(operand Value) Value {
	switch o := operand.(type) {
	case SmallInt:
		return o.Decrement()
	case *BigInt:
		return o.Decrement()
	case Char:
		return o - 1
	case Int64:
		return o - 1
	case Int32:
		return o - 1
	case Int16:
		return o - 1
	case Int8:
		return o - 1
	case UInt64:
		return o - 1
	case UInt32:
		return o - 1
	case UInt16:
		return o - 1
	case UInt8:
		return o - 1
	default:
		return nil
	}
}

// Perform unary plus on a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func UnaryPlus(operand Value) Value {
	switch o := operand.(type) {
	case SmallInt, *BigInt, Float, *BigFloat,
		Float64, Float32, Int64, Int32, Int16, Int8,
		UInt64, UInt32, UInt16, UInt8:
		return o
	default:
		return nil
	}
}

// Perform bitwise not on a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func BitwiseNot(operand Value) Value {
	switch o := operand.(type) {
	case SmallInt:
		return ^o
	case *BigInt:
		return o.BitwiseNot()
	case Int64:
		return ^o
	case Int32:
		return ^o
	case Int16:
		return ^o
	case Int8:
		return ^o
	case UInt64:
		return ^o
	case UInt32:
		return ^o
	case UInt16:
		return ^o
	case UInt8:
		return ^o
	default:
		return nil
	}
}

// Exponentiate two values
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Exponentiate(left, right Value) (result, err Value) {
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
func Modulo(left, right Value) (result, err Value) {
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
func Compare(left, right Value) (result, err Value) {
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
func GreaterThan(left, right Value) (result, err Value) {
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
func GreaterThanEqual(left, right Value) (result, err Value) {
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
func LessThan(left, right Value) (result, err Value) {
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
func LessThanEqual(left, right Value) (result, err Value) {
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
func LaxEqual(left, right Value) Value {
	var result Value

	switch l := left.(type) {
	case SmallInt:
		result = l.LaxEqual(right)
	case *BigInt:
		result = l.LaxEqual(right)
	case Float:
		result = l.LaxEqual(right)
	case *BigFloat:
		result = l.LaxEqual(right)
	case String:
		result = l.LaxEqual(right)
	case *Regex:
		result = l.LaxEqual(right)
	case Char:
		result = l.LaxEqual(right)
	case Symbol:
		result = l.LaxEqual(right)
	case Float64:
		result = StrictFloatLaxEqual(l, right)
	case Float32:
		result = StrictFloatLaxEqual(l, right)
	case Int64:
		result = StrictSignedIntLaxEqual(l, right)
	case Int32:
		result = StrictSignedIntLaxEqual(l, right)
	case Int16:
		result = StrictSignedIntLaxEqual(l, right)
	case Int8:
		result = StrictSignedIntLaxEqual(l, right)
	case UInt64:
		result = StrictUnsignedIntLaxEqual(l, right)
	case UInt32:
		result = StrictUnsignedIntLaxEqual(l, right)
	case UInt16:
		result = StrictUnsignedIntLaxEqual(l, right)
	case UInt8:
		result = StrictUnsignedIntLaxEqual(l, right)
	default:
		return nil
	}

	return result
}

// Check whether left is not equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func LaxNotEqual(left, right Value) Value {
	val := LaxEqual(left, right)
	if val == nil {
		return nil
	}

	return ToNotBool(val)
}

// Check whether left is equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func Equal(left, right Value) Value {
	class := left.Class()
	if !IsA(right, class) {
		return False
	}

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
	case *Regex:
		result = l.Equal(right)
	case Symbol:
		result = l.Equal(right)
	case Char:
		result = l.Equal(right)
	case Float64:
		result = StrictNumericEqual(l, right)
	case Float32:
		result = StrictNumericEqual(l, right)
	case Int64:
		result = StrictNumericEqual(l, right)
	case Int32:
		result = StrictNumericEqual(l, right)
	case Int16:
		result = StrictNumericEqual(l, right)
	case Int8:
		result = StrictNumericEqual(l, right)
	case UInt64:
		result = StrictNumericEqual(l, right)
	case UInt32:
		result = StrictNumericEqual(l, right)
	case UInt16:
		result = StrictNumericEqual(l, right)
	case UInt8:
		result = StrictNumericEqual(l, right)
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
		return ToElkBool(left == right)
	}

	return result
}

// Check whether left is strictly not equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func StrictNotEqual(left, right Value) Value {
	val := StrictEqual(left, right)

	return ToNotBool(val)
}

// Execute a right bit shift >>.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func RightBitshift(left, right Value) (result, err Value) {
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
func LogicalRightBitshift(left, right Value) (result, err Value) {
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
func LeftBitshift(left, right Value) (result, err Value) {
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
func LogicalLeftBitshift(left, right Value) (result, err Value) {
	switch l := left.(type) {
	case Int64:
		result, err = StrictIntLogicalLeftBitshift(l, right, LogicalRightShift64)
	case Int32:
		result, err = StrictIntLogicalLeftBitshift(l, right, LogicalRightShift32)
	case Int16:
		result, err = StrictIntLogicalLeftBitshift(l, right, LogicalRightShift16)
	case Int8:
		result, err = StrictIntLogicalLeftBitshift(l, right, LogicalRightShift8)
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
func BitwiseAnd(left, right Value) (result, err Value) {
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

// Execute a bitwise AND NOT &^.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func BitwiseAndNot(left, right Value) (result, err Value) {
	switch l := left.(type) {
	case SmallInt:
		result, err = l.BitwiseAndNot(right)
	case *BigInt:
		result, err = l.BitwiseAndNot(right)
	case Int64:
		result, err = StrictIntBitwiseAndNot(l, right)
	case Int32:
		result, err = StrictIntBitwiseAndNot(l, right)
	case Int16:
		result, err = StrictIntBitwiseAndNot(l, right)
	case Int8:
		result, err = StrictIntBitwiseAndNot(l, right)
	case UInt64:
		result, err = StrictIntBitwiseAndNot(l, right)
	case UInt32:
		result, err = StrictIntBitwiseAndNot(l, right)
	case UInt16:
		result, err = StrictIntBitwiseAndNot(l, right)
	case UInt8:
		result, err = StrictIntBitwiseAndNot(l, right)
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
func BitwiseOr(left, right Value) (result, err Value) {
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
func BitwiseXor(left, right Value) (result, err Value) {
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
