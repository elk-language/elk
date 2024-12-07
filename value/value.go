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

func (v Value) MustReference() Reference {
	if !v.IsReference() {
		panic(fmt.Sprintf("value `%s` is not a reference", v.Inspect()))
	}
	return v.AsReference()
}

func RefErr(ref Reference, err Value) (Value, Value) {
	if !err.IsNil() {
		return Nil, err
	}
	return Ref(ref), Nil
}

type ToValuer interface {
	ToValue() Value
}

func ToValueErr[T ToValuer](t T, err Value) (Value, Value) {
	if !err.IsNil() {
		return Nil, err
	}
	return t.ToValue(), Nil
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

func (v Value) MustSmallInt() SmallInt {
	if !v.IsSmallInt() {
		panic(fmt.Sprintf("value `%s` is not a SmallInt", v.Inspect()))
	}
	return v.AsSmallInt()
}

func (v Value) IsChar() bool {
	return uintptr(v.data) == CHAR_FLAG
}

func (v Value) AsChar() Char {
	return *(*Char)(unsafe.Pointer(&v.tab))
}

func (v Value) MustChar() Char {
	if !v.IsChar() {
		panic(fmt.Sprintf("value `%s` is not a Char", v.Inspect()))
	}
	return v.AsChar()
}

func (v Value) IsFloat() bool {
	return uintptr(v.data) == FLOAT_FLAG
}

func (v Value) AsFloat() Float {
	return *(*Float)(unsafe.Pointer(&v.tab))
}

func (v Value) MustFloat() Float {
	if !v.IsFloat() {
		panic(fmt.Sprintf("value `%s` is not a Float", v.Inspect()))
	}
	return v.AsFloat()
}

func (v Value) IsFloat32() bool {
	return uintptr(v.data) == FLOAT32_FLAG
}

func (v Value) AsFloat32() Float32 {
	return *(*Float32)(unsafe.Pointer(&v.tab))
}

func (v Value) MustFloat32() Float32 {
	if !v.IsFloat32() {
		panic(fmt.Sprintf("value `%s` is not a Float32", v.Inspect()))
	}
	return v.AsFloat32()
}

func (v Value) IsFloat64() bool {
	return uintptr(v.data) == FLOAT64_FLAG
}

func (v Value) AsFloat64() Float64 {
	return *(*Float64)(unsafe.Pointer(&v.tab))
}

func (v Value) MustFloat64() Float64 {
	if !v.IsFloat64() {
		panic(fmt.Sprintf("value `%s` is not a Float64", v.Inspect()))
	}
	return v.AsFloat64()
}

func (v Value) IsInt8() bool {
	return uintptr(v.data) == INT8_FLAG
}

func (v Value) AsInt8() Int8 {
	return *(*Int8)(unsafe.Pointer(&v.tab))
}

func (v Value) MustInt8() Int8 {
	if !v.IsInt8() {
		panic(fmt.Sprintf("value `%s` is not an Int8", v.Inspect()))
	}
	return v.AsInt8()
}

func (v Value) IsUInt8() bool {
	return uintptr(v.data) == UINT8_FLAG
}

func (v Value) AsUInt8() UInt8 {
	return *(*UInt8)(unsafe.Pointer(&v.tab))
}

func (v Value) MustUInt8() UInt8 {
	if !v.IsUInt8() {
		panic(fmt.Sprintf("value `%s` is not a UInt8", v.Inspect()))
	}
	return v.AsUInt8()
}

func (v Value) IsInt16() bool {
	return uintptr(v.data) == INT16_FLAG
}

func (v Value) AsInt16() Int16 {
	return *(*Int16)(unsafe.Pointer(&v.tab))
}

func (v Value) MustInt16() Int16 {
	if !v.IsInt16() {
		panic(fmt.Sprintf("value `%s` is not an Int16", v.Inspect()))
	}
	return v.AsInt16()
}

func (v Value) IsUInt16() bool {
	return uintptr(v.data) == UINT16_FLAG
}

func (v Value) AsUInt16() UInt16 {
	return *(*UInt16)(unsafe.Pointer(&v.tab))
}

func (v Value) MustUInt16() UInt16 {
	if !v.IsUInt16() {
		panic(fmt.Sprintf("value `%s` is not a UInt16", v.Inspect()))
	}
	return v.AsUInt16()
}

func (v Value) IsInt32() bool {
	return uintptr(v.data) == INT32_FLAG
}

func (v Value) AsInt32() Int32 {
	return *(*Int32)(unsafe.Pointer(&v.tab))
}

func (v Value) MustInt32() Int32 {
	if !v.IsInt32() {
		panic(fmt.Sprintf("value `%s` is not an Int32", v.Inspect()))
	}
	return v.AsInt32()
}

func (v Value) IsUInt32() bool {
	return uintptr(v.data) == UINT32_FLAG
}

func (v Value) AsUInt32() UInt32 {
	return *(*UInt32)(unsafe.Pointer(&v.tab))
}

func (v Value) MustUInt32() UInt32 {
	if !v.IsUInt32() {
		panic(fmt.Sprintf("value `%s` is not a UInt32", v.Inspect()))
	}
	return v.AsUInt32()
}

func (v Value) IsInt64() bool {
	return uintptr(v.data) == INT64_FLAG
}

func (v Value) AsInt64() Int64 {
	return *(*Int64)(unsafe.Pointer(&v.tab))
}

func (v Value) MustInt64() Int64 {
	if !v.IsInt64() {
		panic(fmt.Sprintf("value `%s` is not an Int64", v.Inspect()))
	}
	return v.AsInt64()
}

func (v Value) IsUInt64() bool {
	return uintptr(v.data) == UINT64_FLAG
}

func (v Value) AsUInt64() UInt64 {
	return *(*UInt64)(unsafe.Pointer(&v.tab))
}

func (v Value) MustUInt64() UInt64 {
	if !v.IsUInt64() {
		panic(fmt.Sprintf("value `%s` is not a UInt64", v.Inspect()))
	}
	return v.AsUInt64()
}

func (v Value) IsDuration() bool {
	return uintptr(v.data) == DURATION_FLAG
}

func (v Value) AsDuration() Duration {
	return *(*Duration)(unsafe.Pointer(&v.tab))
}

func (v Value) MustDuration() Duration {
	if !v.IsDuration() {
		panic(fmt.Sprintf("value `%s` is not a Duration", v.Inspect()))
	}
	return v.AsDuration()
}

func (v Value) IsTrue() bool {
	return uintptr(v.data) == TRUE_FLAG
}

func (v Value) AsTrue() TrueType {
	return *(*TrueType)(unsafe.Pointer(&v.tab))
}

func (v Value) MustTrue() TrueType {
	if !v.IsTrue() {
		panic(fmt.Sprintf("value `%s` is not True", v.Inspect()))
	}
	return v.AsTrue()
}

func (v Value) IsSymbol() bool {
	return uintptr(v.data) == SYMBOL_FLAG
}

func (v Value) AsSymbol() Symbol {
	return *(*Symbol)(unsafe.Pointer(&v.tab))
}

func (v Value) MustSymbol() Symbol {
	if !v.IsSymbol() {
		panic(fmt.Sprintf("value `%s` is not a Symbol", v.Inspect()))
	}
	return v.AsSymbol()
}

func (v Value) IsFalse() bool {
	return uintptr(v.data) == FALSE_FLAG
}

func (v Value) AsFalse() FalseType {
	return *(*FalseType)(unsafe.Pointer(&v.tab))
}

func (v Value) MustFalse() FalseType {
	if !v.IsFalse() {
		panic(fmt.Sprintf("value `%s` is not False", v.Inspect()))
	}
	return v.AsFalse()
}

func (v Value) IsNil() bool {
	return uintptr(v.data) == NIL_FLAG
}

func (v Value) AsNil() NilType {
	return *(*NilType)(unsafe.Pointer(&v.tab))
}

func (v Value) MustNil() NilType {
	if !v.IsNil() {
		panic(fmt.Sprintf("value `%s` is not Nil", v.Inspect()))
	}
	return v.AsNil()
}

func (v Value) IsUndefined() bool {
	return uintptr(v.data) == UNDEFINED_FLAG
}

func (v Value) AsUndefined() UndefinedType {
	return *(*UndefinedType)(unsafe.Pointer(&v.tab))
}

func (v Value) MustUndefined() UndefinedType {
	if !v.IsUndefined() {
		panic(fmt.Sprintf("value `%s` is not Undefined", v.Inspect()))
	}
	return v.AsUndefined()
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
	if right.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			result, err = l.Subtract(right)
		case *BigFloat:
			result, err = l.Subtract(right)
		case String:
			var r String
			r, err = l.RemoveSuffix(right)
			result = Ref(r)
		case Float64:
			var r Float64
			r, err = l.Subtract(right)
			result = r.ToValue()
		case Int64:
			var r Int64
			r, err = l.Subtract(right)
			result = r.ToValue()
		case UInt64:
			var r UInt64
			r, err = l.Subtract(right)
			result = r.ToValue()
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
		result, err = l.Subtract(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		result, err = l.Subtract(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		var r Float64
		r, err = l.Subtract(right)
		result = r.ToValue()
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		var r Float32
		r, err = l.Subtract(right)
		result = r.ToValue()
	case INT64_FLAG:
		l := left.AsInt64()
		var r Int64
		r, err = l.Subtract(right)
		result = r.ToValue()
	case INT32_FLAG:
		l := left.AsInt32()
		var r Int32
		r, err = l.Subtract(right)
		result = r.ToValue()
	case INT16_FLAG:
		l := left.AsInt16()
		var r Int16
		r, err = l.Subtract(right)
		result = r.ToValue()
	case INT8_FLAG:
		l := left.AsInt8()
		var r Int8
		r, err = l.Subtract(right)
		result = r.ToValue()
	case UINT64_FLAG:
		l := left.AsUInt64()
		var r UInt64
		r, err = l.Subtract(right)
		result = r.ToValue()
	case UINT32_FLAG:
		l := left.AsUInt32()
		var r UInt32
		r, err = l.Subtract(right)
		result = r.ToValue()
	case UINT16_FLAG:
		l := left.AsUInt16()
		var r UInt16
		r, err = l.Subtract(right)
		result = r.ToValue()
	case UINT8_FLAG:
		l := left.AsUInt8()
		var r UInt8
		r, err = l.Subtract(right)
		result = r.ToValue()
	default:
		return Nil, Nil
	}

	if !err.IsNil() {
		return Nil, err
	}
	return result, Nil
}

// Multiply two values
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Multiply(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			result, err = l.Multiply(right)
		case *BigFloat:
			result, err = l.Multiply(right)
		case Float64:
			var r Float64
			r, err = l.Multiply(right)
			result = r.ToValue()
		case Int64:
			var r Int64
			r, err = l.Multiply(right)
			result = r.ToValue()
		case UInt64:
			var r UInt64
			r, err = l.Multiply(right)
			result = r.ToValue()
		case String:
			var r String
			r, err = l.Repeat(right)
			result = Ref(r)
		case *Regex:
			result, err = l.Repeat(right)
		case *ArrayList:
			var r *ArrayList
			r, err = l.Repeat(right)
			result = Ref(r)
		case *ArrayTuple:
			var r *ArrayTuple
			r, err = l.Repeat(right)
			result = Ref(r)
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		result, err = l.Multiply(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		result, err = l.Multiply(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		var r Float64
		r, err = l.Multiply(right)
		result = r.ToValue()
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		var r Float32
		r, err = l.Multiply(right)
		result = r.ToValue()
	case INT64_FLAG:
		l := left.AsInt64()
		var r Int64
		r, err = l.Multiply(right)
		result = r.ToValue()
	case INT32_FLAG:
		l := left.AsInt32()
		var r Int32
		r, err = l.Multiply(right)
		result = r.ToValue()
	case INT16_FLAG:
		l := left.AsInt16()
		var r Int16
		r, err = l.Multiply(right)
		result = r.ToValue()
	case INT8_FLAG:
		l := left.AsInt8()
		var r Int8
		r, err = l.Multiply(right)
		result = r.ToValue()
	case UINT64_FLAG:
		l := left.AsUInt64()
		var r UInt64
		r, err = l.Multiply(right)
		result = r.ToValue()
	case UINT32_FLAG:
		l := left.AsUInt32()
		var r UInt32
		r, err = l.Multiply(right)
		result = r.ToValue()
	case UINT16_FLAG:
		l := left.AsUInt16()
		var r UInt16
		r, err = l.Multiply(right)
		result = r.ToValue()
	case UINT8_FLAG:
		l := left.AsUInt8()
		var r UInt8
		r, err = l.Multiply(right)
		result = r.ToValue()
	case CHAR_FLAG:
		l := left.AsChar()
		var r String
		r, err = l.Repeat(right)
		result = Ref(r)
	default:
		return Nil, Nil
	}

	if !err.IsNil() {
		return Nil, err
	}
	return result, Nil
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
	if operand.IsReference() {
		switch o := operand.AsReference().(type) {
		case *BigInt:
			return Ref(o.Negate())
		case *BigFloat:
			return Ref(o.Negate())
		case Float64:
			return (-o).ToValue()
		case Int64:
			return (-o).ToValue()
		case UInt64:
			return (-o).ToValue()
		default:
			return Nil
		}
	}

	switch operand.ValueFlag() {
	case SMALL_INT_FLAG:
		o := operand.AsSmallInt()
		return o.Negate()
	case FLOAT_FLAG:
		o := operand.AsFloat()
		return (-o).ToValue()
	case FLOAT64_FLAG:
		o := operand.AsFloat64()
		return (-o).ToValue()
	case FLOAT32_FLAG:
		o := operand.AsFloat32()
		return (-o).ToValue()
	case INT64_FLAG:
		o := operand.AsInt64()
		return (-o).ToValue()
	case INT32_FLAG:
		o := operand.AsInt32()
		return (-o).ToValue()
	case INT16_FLAG:
		o := operand.AsInt16()
		return (-o).ToValue()
	case INT8_FLAG:
		o := operand.AsInt8()
		return (-o).ToValue()
	case UINT64_FLAG:
		o := operand.AsUInt64()
		return (-o).ToValue()
	case UINT32_FLAG:
		o := operand.AsUInt32()
		return (-o).ToValue()
	case UINT16_FLAG:
		o := operand.AsUInt16()
		return (-o).ToValue()
	case UINT8_FLAG:
		o := operand.AsUInt8()
		return (-o).ToValue()
	default:
		return Nil
	}
}

// Increment a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func Increment(operand Value) Value {
	if operand.IsReference() {
		switch o := operand.AsReference().(type) {
		case *BigInt:
			return Ref(o.Increment())
		case Int64:
			return (o + 1).ToValue()
		case UInt64:
			return (o + 1).ToValue()
		default:
			return Nil
		}
	}

	switch operand.ValueFlag() {
	case SMALL_INT_FLAG:
		o := operand.AsSmallInt()
		return o.Increment()
	case CHAR_FLAG:
		o := operand.AsChar()
		return (o + 1).ToValue()
	case INT64_FLAG:
		o := operand.AsInt64()
		return (o + 1).ToValue()
	case INT32_FLAG:
		o := operand.AsInt32()
		return (o + 1).ToValue()
	case INT16_FLAG:
		o := operand.AsInt16()
		return (o + 1).ToValue()
	case INT8_FLAG:
		o := operand.AsInt8()
		return (o + 1).ToValue()
	case UINT64_FLAG:
		o := operand.AsUInt64()
		return (o + 1).ToValue()
	case UINT32_FLAG:
		o := operand.AsUInt32()
		return (o + 1).ToValue()
	case UINT16_FLAG:
		o := operand.AsUInt16()
		return (o + 1).ToValue()
	case UINT8_FLAG:
		o := operand.AsUInt8()
		return (o + 1).ToValue()
	default:
		return Nil
	}
}

// Decrement a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func Decrement(operand Value) Value {
	if operand.IsReference() {
		switch o := operand.AsReference().(type) {
		case *BigInt:
			return o.Decrement()
		case Int64:
			return (o - 1).ToValue()
		case UInt64:
			return (o - 1).ToValue()
		default:
			return Nil
		}
	}

	switch operand.ValueFlag() {
	case SMALL_INT_FLAG:
		o := operand.AsSmallInt()
		return o.Decrement()
	case CHAR_FLAG:
		o := operand.AsChar()
		return (o - 1).ToValue()
	case INT64_FLAG:
		o := operand.AsInt64()
		return (o - 1).ToValue()
	case INT32_FLAG:
		o := operand.AsInt32()
		return (o - 1).ToValue()
	case INT16_FLAG:
		o := operand.AsInt16()
		return (o - 1).ToValue()
	case INT8_FLAG:
		o := operand.AsInt8()
		return (o - 1).ToValue()
	case UINT64_FLAG:
		o := operand.AsUInt64()
		return (o - 1).ToValue()
	case UINT32_FLAG:
		o := operand.AsUInt32()
		return (o - 1).ToValue()
	case UINT16_FLAG:
		o := operand.AsUInt16()
		return (o - 1).ToValue()
	case UINT8_FLAG:
		o := operand.AsUInt8()
		return (o - 1).ToValue()
	default:
		return Nil
	}
}

// Perform unary plus on a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func UnaryPlus(operand Value) Value {
	if operand.IsReference() {
		switch operand.AsReference().(type) {
		case *BigInt, *BigFloat,
			Float64, Int64,
			UInt64:
			return operand
		default:
			return Nil
		}
	}

	switch operand.ValueFlag() {
	case SMALL_INT_FLAG, FLOAT_FLAG,
		FLOAT64_FLAG, FLOAT32_FLAG, INT64_FLAG, INT32_FLAG, INT16_FLAG, INT8_FLAG,
		UINT64_FLAG, UINT32_FLAG, UINT16_FLAG, UINT8_FLAG:
		return operand
	default:
		return Nil
	}
}

// Perform bitwise not on a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func BitwiseNot(operand Value) Value {
	if operand.IsReference() {
		switch o := operand.AsReference().(type) {
		case *BigInt:
			return Ref(o.BitwiseNot())
		case Int64:
			return (^o).ToValue()
		case UInt64:
			return (^o).ToValue()
		default:
			return Nil
		}
	}

	switch operand.ValueFlag() {
	case SMALL_INT_FLAG:
		o := operand.AsSmallInt()
		return (^o).ToValue()
	case INT64_FLAG:
		o := operand.AsInt64()
		return (^o).ToValue()
	case INT32_FLAG:
		o := operand.AsInt32()
		return (^o).ToValue()
	case INT16_FLAG:
		o := operand.AsInt16()
		return (^o).ToValue()
	case INT8_FLAG:
		o := operand.AsInt8()
		return (^o).ToValue()
	case UINT64_FLAG:
		o := operand.AsUInt64()
		return (^o).ToValue()
	case UINT32_FLAG:
		o := operand.AsUInt32()
		return (^o).ToValue()
	case UINT16_FLAG:
		o := operand.AsUInt16()
		return (^o).ToValue()
	case UINT8_FLAG:
		o := operand.AsUInt8()
		return (^o).ToValue()
	default:
		return Nil
	}
}

// Exponentiate two values
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Exponentiate(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			result, err = l.Exponentiate(right)
		case *BigFloat:
			result, err = l.Exponentiate(right)
		case Float64:
			var r Float64
			r, err = l.Exponentiate(right)
			result = r.ToValue()
		case Int64:
			var r Int64
			r, err = l.Exponentiate(right)
			result = r.ToValue()
		case UInt64:
			var r UInt64
			r, err = l.Exponentiate(right)
			result = r.ToValue()
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
		result, err = l.Exponentiate(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		result, err = l.Exponentiate(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		var r Float64
		r, err = l.Exponentiate(right)
		result = r.ToValue()
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		var r Float32
		r, err = l.Exponentiate(right)
		result = r.ToValue()
	case INT64_FLAG:
		l := left.AsInt64()
		var r Int64
		r, err = l.Exponentiate(right)
		result = r.ToValue()
	case INT32_FLAG:
		l := left.AsInt32()
		var r Int32
		r, err = l.Exponentiate(right)
		result = r.ToValue()
	case INT16_FLAG:
		l := left.AsInt16()
		var r Int16
		r, err = l.Exponentiate(right)
		result = r.ToValue()
	case INT8_FLAG:
		l := left.AsInt8()
		var r Int8
		r, err = l.Exponentiate(right)
		result = r.ToValue()
	case UINT64_FLAG:
		l := left.AsUInt64()
		var r UInt64
		r, err = l.Exponentiate(right)
		result = r.ToValue()
	case UINT32_FLAG:
		l := left.AsUInt32()
		var r UInt32
		r, err = l.Exponentiate(right)
		result = r.ToValue()
	case UINT16_FLAG:
		l := left.AsUInt16()
		var r UInt16
		r, err = l.Exponentiate(right)
		result = r.ToValue()
	case UINT8_FLAG:
		l := left.AsUInt8()
		var r UInt8
		r, err = l.Exponentiate(right)
		result = r.ToValue()
	default:
		return Nil, Nil
	}
	if !err.IsNil() {
		return Nil, err
	}
	return result, Nil
}

// Perform modulo on two values
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func Modulo(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			result, err = l.Modulo(right)
		case *BigFloat:
			result, err = l.Modulo(right)
		case Float64:
			var r Float64
			r, err = l.Modulo(right)
			result = r.ToValue()
		case Int64:
			var r Int64
			r, err = l.Modulo(right)
			result = r.ToValue()
		case UInt64:
			var r UInt64
			r, err = l.Modulo(right)
			result = r.ToValue()
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
		result, err = l.Modulo(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		result, err = l.Modulo(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		var r Float64
		r, err = l.Modulo(right)
		result = r.ToValue()
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		var r Float32
		r, err = l.Modulo(right)
		result = r.ToValue()
	case INT64_FLAG:
		l := left.AsInt64()
		var r Int64
		r, err = l.Modulo(right)
		result = r.ToValue()
	case INT32_FLAG:
		l := left.AsInt32()
		var r Int32
		r, err = l.Modulo(right)
		result = r.ToValue()
	case INT16_FLAG:
		l := left.AsInt16()
		var r Int16
		r, err = l.Modulo(right)
		result = r.ToValue()
	case INT8_FLAG:
		l := left.AsInt8()
		var r Int8
		r, err = l.Modulo(right)
		result = r.ToValue()
	case UINT64_FLAG:
		l := left.AsUInt64()
		var r UInt64
		r, err = l.Modulo(right)
		result = r.ToValue()
	case UINT32_FLAG:
		l := left.AsUInt32()
		var r UInt32
		r, err = l.Modulo(right)
		result = r.ToValue()
	case UINT16_FLAG:
		l := left.AsUInt16()
		var r UInt16
		r, err = l.Modulo(right)
		result = r.ToValue()
	case UINT8_FLAG:
		l := left.AsUInt8()
		var r UInt8
		r, err = l.Modulo(right)
		result = r.ToValue()
	default:
		return Nil, Nil
	}
	if !err.IsNil() {
		return Nil, err
	}
	return result, Nil
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
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.Compare(right)
		case *BigFloat:
			return l.Compare(right)
		case String:
			return l.Compare(right)
		case Float64:
			return l.Compare(right)
		case Int64:
			return l.Compare(right)
		case UInt64:
			return l.Compare(right)
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.Compare(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.Compare(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.Compare(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		return l.Compare(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.Compare(right)
	case INT64_FLAG:
		l := left.AsInt64()
		return l.Compare(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.Compare(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.Compare(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.Compare(right)
	case UINT64_FLAG:
		l := left.AsUInt64()
		return l.Compare(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.Compare(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.Compare(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.Compare(right)
	default:
		return Nil, Nil
	}
}

// Check whether left is greater than right.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func GreaterThan(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.GreaterThan(right)
		case *BigFloat:
			return l.GreaterThan(right)
		case String:
			return l.GreaterThan(right)
		case Float64:
			return l.GreaterThan(right)
		case Int64:
			return l.GreaterThan(right)
		case UInt64:
			return l.GreaterThan(right)
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.GreaterThan(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.GreaterThan(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.GreaterThan(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		return l.GreaterThan(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.GreaterThan(right)
	case INT64_FLAG:
		l := left.AsInt64()
		return l.GreaterThan(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.GreaterThan(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.GreaterThan(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.GreaterThan(right)
	case UINT64_FLAG:
		l := left.AsUInt64()
		return l.GreaterThan(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.GreaterThan(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.GreaterThan(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.GreaterThan(right)
	default:
		return Nil, Nil
	}
}

// Check whether left is greater than or equal to right.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func GreaterThanEqual(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.GreaterThanEqual(right)
		case *BigFloat:
			return l.GreaterThanEqual(right)
		case String:
			return l.GreaterThanEqual(right)
		case Float64:
			return l.GreaterThanEqual(right)
		case Int64:
			return l.GreaterThanEqual(right)
		case UInt64:
			return l.GreaterThanEqual(right)
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.GreaterThanEqual(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.GreaterThanEqual(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.GreaterThanEqual(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		return l.GreaterThanEqual(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.GreaterThanEqual(right)
	case INT64_FLAG:
		l := left.AsInt64()
		return l.GreaterThanEqual(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.GreaterThanEqual(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.GreaterThanEqual(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.GreaterThanEqual(right)
	case UINT64_FLAG:
		l := left.AsUInt64()
		return l.GreaterThanEqual(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.GreaterThanEqual(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.GreaterThanEqual(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.GreaterThanEqual(right)
	default:
		return Nil, Nil
	}
}

// Check whether left is less than right.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func LessThan(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.LessThan(right)
		case *BigFloat:
			return l.LessThan(right)
		case String:
			return l.LessThan(right)
		case Float64:
			return l.LessThan(right)
		case Int64:
			return l.LessThan(right)
		case UInt64:
			return l.LessThan(right)
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.LessThan(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.LessThan(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.LessThan(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		return l.LessThan(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.LessThan(right)
	case INT64_FLAG:
		l := left.AsInt64()
		return l.LessThan(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.LessThan(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.LessThan(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.LessThan(right)
	case UINT64_FLAG:
		l := left.AsUInt64()
		return l.LessThan(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.LessThan(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.LessThan(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.LessThan(right)
	default:
		return Nil, Nil
	}
}

// Check whether left is less than or equal to right.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func LessThanEqual(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.LessThanEqual(right)
		case *BigFloat:
			return l.LessThanEqual(right)
		case String:
			return l.LessThanEqual(right)
		case Float64:
			return l.LessThanEqual(right)
		case Int64:
			return l.LessThanEqual(right)
		case UInt64:
			return l.LessThanEqual(right)
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.LessThanEqual(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.LessThanEqual(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.LessThanEqual(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		return l.LessThanEqual(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.LessThanEqual(right)
	case INT64_FLAG:
		l := left.AsInt64()
		return l.LessThanEqual(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.LessThanEqual(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.LessThanEqual(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.LessThanEqual(right)
	case UINT64_FLAG:
		l := left.AsUInt64()
		return l.LessThanEqual(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.LessThanEqual(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.LessThanEqual(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.LessThanEqual(right)
	default:
		return Nil, Nil
	}
}

// Check whether left is equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func LaxEqual(left, right Value) Value {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.LaxEqual(right)
		case *BigFloat:
			return l.LaxEqual(right)
		case String:
			return l.LaxEqual(right)
		case *Regex:
			return l.LaxEqual(right)
		case Float64:
			return StrictFloatLaxEqual(l, right)
		case Int64:
			return StrictSignedIntLaxEqual(l, right)
		case UInt64:
			return StrictUnsignedIntLaxEqual(l, right)
		default:
			return Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.LaxEqual(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.LaxEqual(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.LaxEqual(right)
	case SYMBOL_FLAG:
		l := left.AsSymbol()
		return l.LaxEqual(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		return StrictFloatLaxEqual(l, right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return StrictFloatLaxEqual(l, right)
	case INT64_FLAG:
		l := left.AsInt64()
		return StrictSignedIntLaxEqual(l, right)
	case INT32_FLAG:
		l := left.AsInt32()
		return StrictSignedIntLaxEqual(l, right)
	case INT16_FLAG:
		l := left.AsInt16()
		return StrictSignedIntLaxEqual(l, right)
	case INT8_FLAG:
		l := left.AsInt8()
		return StrictSignedIntLaxEqual(l, right)
	case UINT64_FLAG:
		l := left.AsUInt64()
		return StrictUnsignedIntLaxEqual(l, right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return StrictUnsignedIntLaxEqual(l, right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return StrictUnsignedIntLaxEqual(l, right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return StrictUnsignedIntLaxEqual(l, right)
	default:
		return Nil
	}
}

// Check whether left is not equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func LaxNotEqual(left, right Value) Value {
	val := LaxEqual(left, right)
	if val.IsNil() {
		return Nil
	}

	return ToNotBool(val)
}

// Check whether left is equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func Equal(left, right Value) Value {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.Equal(right)
		case *BigFloat:
			return l.Equal(right)
		case String:
			return l.Equal(right)
		case Float64:
			return l.Equal(right)
		case Int64:
			return l.Equal(right)
		case UInt64:
			return l.Equal(right)
		default:
			return ToElkBool(left == right)
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.Equal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.Equal(right)
	case SYMBOL_FLAG:
		l := left.AsSymbol()
		return l.Equal(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.Equal(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		return l.Equal(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.Equal(right)
	case INT64_FLAG:
		l := left.AsInt64()
		return l.Equal(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.Equal(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.Equal(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.Equal(right)
	case UINT64_FLAG:
		l := left.AsUInt64()
		return l.Equal(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.Equal(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.Equal(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.Equal(right)
	default:
		return ToElkBool(left == right)
	}
}

// Check whether left is not equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func NotEqual(left, right Value) Value {
	val := Equal(left, right)
	if val.IsNil() {
		return Nil
	}

	return ToNotBool(val)
}

// Check whether left is strictly equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func StrictEqual(left, right Value) Value {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.StrictEqual(right)
		case *BigFloat:
			return l.StrictEqual(right)
		case String:
			return l.StrictEqual(right)
		case Float64:
			return l.StrictEqual(right)
		case Int64:
			return l.StrictEqual(right)
		case UInt64:
			return l.StrictEqual(right)
		default:
			return ToElkBool(left == right)
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.StrictEqual(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.StrictEqual(right)
	case SYMBOL_FLAG:
		l := left.AsSymbol()
		return l.StrictEqual(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.StrictEqual(right)
	case FLOAT64_FLAG:
		l := left.AsFloat64()
		return l.StrictEqual(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.StrictEqual(right)
	case INT64_FLAG:
		l := left.AsInt64()
		return l.StrictEqual(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.StrictEqual(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.StrictEqual(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.StrictEqual(right)
	case UINT64_FLAG:
		l := left.AsUInt64()
		return l.StrictEqual(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.StrictEqual(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.StrictEqual(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.StrictEqual(right)
	default:
		return ToElkBool(left == right)
	}
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
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.RightBitshift(right)
		case Int64:
			r, err := StrictIntRightBitshift(l, right)
			return r.ToValue(), err
		case UInt64:
			r, err := StrictIntRightBitshift(l, right)
			return r.ToValue(), err
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.RightBitshift(right)
	case INT64_FLAG:
		l := left.AsInt64()
		result, err := StrictIntRightBitshift(l, right)
		return result.ToValue(), err
	case INT32_FLAG:
		l := left.AsInt32()
		result, err := StrictIntRightBitshift(l, right)
		return result.ToValue(), err
	case INT16_FLAG:
		l := left.AsInt16()
		result, err := StrictIntRightBitshift(l, right)
		return result.ToValue(), err
	case INT8_FLAG:
		l := left.AsInt8()
		result, err := StrictIntRightBitshift(l, right)
		return result.ToValue(), err
	case UINT64_FLAG:
		l := left.AsUInt64()
		result, err := StrictIntRightBitshift(l, right)
		return result.ToValue(), err
	case UINT32_FLAG:
		l := left.AsUInt32()
		result, err := StrictIntRightBitshift(l, right)
		return result.ToValue(), err
	case UINT16_FLAG:
		l := left.AsUInt16()
		result, err := StrictIntRightBitshift(l, right)
		return result.ToValue(), err
	case UINT8_FLAG:
		l := left.AsUInt8()
		result, err := StrictIntRightBitshift(l, right)
		return result.ToValue(), err
	default:
		return Nil, Nil
	}
}

// Execute a logical right bit shift >>>.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func LogicalRightBitshift(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case Int64:
			r, err := StrictIntLogicalRightBitshift(l, right, LogicalRightShift64)
			return r.ToValue(), err
		case UInt64:
			r, err := StrictIntRightBitshift(l, right)
			return r.ToValue(), err
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case INT64_FLAG:
		l := left.AsInt64()
		r, err := StrictIntLogicalRightBitshift(l, right, LogicalRightShift64)
		return r.ToValue(), err
	case INT32_FLAG:
		l := left.AsInt32()
		r, err := StrictIntLogicalRightBitshift(l, right, LogicalRightShift32)
		return r.ToValue(), err
	case INT16_FLAG:
		l := left.AsInt16()
		r, err := StrictIntLogicalRightBitshift(l, right, LogicalRightShift16)
		return r.ToValue(), err
	case INT8_FLAG:
		l := left.AsInt8()
		r, err := StrictIntLogicalRightBitshift(l, right, LogicalRightShift8)
		return r.ToValue(), err
	case UINT64_FLAG:
		l := left.AsUInt64()
		r, err := StrictIntRightBitshift(l, right)
		return r.ToValue(), err
	case UINT32_FLAG:
		l := left.AsUInt32()
		r, err := StrictIntRightBitshift(l, right)
		return r.ToValue(), err
	case UINT16_FLAG:
		l := left.AsUInt16()
		r, err := StrictIntRightBitshift(l, right)
		return r.ToValue(), err
	case UINT8_FLAG:
		l := left.AsUInt8()
		r, err := StrictIntRightBitshift(l, right)
		return r.ToValue(), err
	default:
		return Nil, Nil
	}
}

// Execute a left bit shift <<.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func LeftBitshift(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.LeftBitshift(right)
		case Int64:
			r, err := StrictIntLeftBitshift(l, right)
			return r.ToValue(), err
		case UInt64:
			r, err := StrictIntLeftBitshift(l, right)
			return r.ToValue(), err
		case *ArrayList:
			l.Append(right)
			return Ref(l), Nil
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.LeftBitshift(right)
	case INT64_FLAG:
		l := left.AsInt64()
		result, err := StrictIntLeftBitshift(l, right)
		return result.ToValue(), err
	case INT32_FLAG:
		l := left.AsInt32()
		result, err := StrictIntLeftBitshift(l, right)
		return result.ToValue(), err
	case INT16_FLAG:
		l := left.AsInt16()
		result, err := StrictIntLeftBitshift(l, right)
		return result.ToValue(), err
	case INT8_FLAG:
		l := left.AsInt8()
		result, err := StrictIntLeftBitshift(l, right)
		return result.ToValue(), err
	case UINT64_FLAG:
		l := left.AsUInt64()
		result, err := StrictIntLeftBitshift(l, right)
		return result.ToValue(), err
	case UINT32_FLAG:
		l := left.AsUInt32()
		result, err := StrictIntLeftBitshift(l, right)
		return result.ToValue(), err
	case UINT16_FLAG:
		l := left.AsUInt16()
		result, err := StrictIntLeftBitshift(l, right)
		return result.ToValue(), err
	case UINT8_FLAG:
		l := left.AsUInt8()
		result, err := StrictIntLeftBitshift(l, right)
		return result.ToValue(), err
	default:
		return Nil, Nil
	}
}

// Execute a logical left bit shift <<<.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func LogicalLeftBitshift(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case Int64:
			r, err := StrictIntLogicalLeftBitshift(l, right, LogicalRightShift64)
			return r.ToValue(), err
		case UInt64:
			r, err := StrictIntLeftBitshift(l, right)
			return r.ToValue(), err
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case INT64_FLAG:
		l := left.AsInt64()
		r, err := StrictIntLogicalLeftBitshift(l, right, LogicalRightShift64)
		return r.ToValue(), err
	case INT32_FLAG:
		l := left.AsInt32()
		r, err := StrictIntLogicalLeftBitshift(l, right, LogicalRightShift32)
		return r.ToValue(), err
	case INT16_FLAG:
		l := left.AsInt16()
		r, err := StrictIntLogicalLeftBitshift(l, right, LogicalRightShift16)
		return r.ToValue(), err
	case INT8_FLAG:
		l := left.AsInt8()
		r, err := StrictIntLogicalLeftBitshift(l, right, LogicalRightShift8)
		return r.ToValue(), err
	case UINT64_FLAG:
		l := left.AsUInt64()
		r, err := StrictIntLeftBitshift(l, right)
		return r.ToValue(), err
	case UINT32_FLAG:
		l := left.AsUInt32()
		r, err := StrictIntLeftBitshift(l, right)
		return r.ToValue(), err
	case UINT16_FLAG:
		l := left.AsUInt16()
		r, err := StrictIntLeftBitshift(l, right)
		return r.ToValue(), err
	case UINT8_FLAG:
		l := left.AsUInt8()
		r, err := StrictIntLeftBitshift(l, right)
		return r.ToValue(), err
	default:
		return Nil, Nil
	}
}

// Execute a bitwise AND &.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func BitwiseAnd(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.BitwiseAnd(right)
		case Int64:
			result, err := l.BitwiseAnd(right)
			return result.ToValue(), err
		case UInt64:
			result, err := l.BitwiseAnd(right)
			return result.ToValue(), err
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.BitwiseAnd(right)
	case INT64_FLAG:
		l := left.AsInt64()
		result, err := l.BitwiseAnd(right)
		return result.ToValue(), err
	case INT32_FLAG:
		l := left.AsInt32()
		result, err := l.BitwiseAnd(right)
		return result.ToValue(), err
	case INT16_FLAG:
		l := left.AsInt16()
		result, err := l.BitwiseAnd(right)
		return result.ToValue(), err
	case INT8_FLAG:
		l := left.AsInt8()
		result, err := l.BitwiseAnd(right)
		return result.ToValue(), err
	case UINT64_FLAG:
		l := left.AsUInt64()
		result, err := l.BitwiseAnd(right)
		return result.ToValue(), err
	case UINT32_FLAG:
		l := left.AsUInt32()
		result, err := l.BitwiseAnd(right)
		return result.ToValue(), err
	case UINT16_FLAG:
		l := left.AsUInt16()
		result, err := l.BitwiseAnd(right)
		return result.ToValue(), err
	case UINT8_FLAG:
		l := left.AsUInt8()
		result, err := l.BitwiseAnd(right)
		return result.ToValue(), err
	default:
		return Nil, Nil
	}
}

// Execute a bitwise AND NOT &^.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func BitwiseAndNot(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.BitwiseAndNot(right)
		case Int64:
			result, err := l.BitwiseAndNot(right)
			return result.ToValue(), err
		case UInt64:
			result, err := l.BitwiseAndNot(right)
			return result.ToValue(), err
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.BitwiseAndNot(right)
	case INT64_FLAG:
		l := left.AsInt64()
		result, err := l.BitwiseAndNot(right)
		return result.ToValue(), err
	case INT32_FLAG:
		l := left.AsInt32()
		result, err := l.BitwiseAndNot(right)
		return result.ToValue(), err
	case INT16_FLAG:
		l := left.AsInt16()
		result, err := l.BitwiseAndNot(right)
		return result.ToValue(), err
	case INT8_FLAG:
		l := left.AsInt8()
		result, err := l.BitwiseAndNot(right)
		return result.ToValue(), err
	case UINT64_FLAG:
		l := left.AsUInt64()
		result, err := l.BitwiseAndNot(right)
		return result.ToValue(), err
	case UINT32_FLAG:
		l := left.AsUInt32()
		result, err := l.BitwiseAndNot(right)
		return result.ToValue(), err
	case UINT16_FLAG:
		l := left.AsUInt16()
		result, err := l.BitwiseAndNot(right)
		return result.ToValue(), err
	case UINT8_FLAG:
		l := left.AsUInt8()
		result, err := l.BitwiseAndNot(right)
		return result.ToValue(), err
	default:
		return Nil, Nil
	}
}

// Execute a bitwise OR |.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func BitwiseOr(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.BitwiseOr(right)
		case Int64:
			result, err := l.BitwiseOr(right)
			return result.ToValue(), err
		case UInt64:
			result, err := l.BitwiseOr(right)
			return result.ToValue(), err
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.BitwiseOr(right)
	case INT64_FLAG:
		l := left.AsInt64()
		result, err := l.BitwiseOr(right)
		return result.ToValue(), err
	case INT32_FLAG:
		l := left.AsInt32()
		result, err := l.BitwiseOr(right)
		return result.ToValue(), err
	case INT16_FLAG:
		l := left.AsInt16()
		result, err := l.BitwiseOr(right)
		return result.ToValue(), err
	case INT8_FLAG:
		l := left.AsInt8()
		result, err := l.BitwiseOr(right)
		return result.ToValue(), err
	case UINT64_FLAG:
		l := left.AsUInt64()
		result, err := l.BitwiseOr(right)
		return result.ToValue(), err
	case UINT32_FLAG:
		l := left.AsUInt32()
		result, err := l.BitwiseOr(right)
		return result.ToValue(), err
	case UINT16_FLAG:
		l := left.AsUInt16()
		result, err := l.BitwiseOr(right)
		return result.ToValue(), err
	case UINT8_FLAG:
		l := left.AsUInt8()
		result, err := l.BitwiseOr(right)
		return result.ToValue(), err
	default:
		return Nil, Nil
	}
}

// Execute a bitwise XOR ^.
// When successful returns (result, nil).
// When an error occurred returns (nil, error).
// When there are no builtin addition functions for the given type returns (nil, nil).
func BitwiseXor(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.BitwiseXor(right)
		case Int64:
			result, err := l.BitwiseXor(right)
			return result.ToValue(), err
		case UInt64:
			result, err := l.BitwiseXor(right)
			return result.ToValue(), err
		default:
			return Nil, Nil
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.BitwiseXor(right)
	case INT64_FLAG:
		l := left.AsInt64()
		result, err := l.BitwiseXor(right)
		return result.ToValue(), err
	case INT32_FLAG:
		l := left.AsInt32()
		result, err := l.BitwiseXor(right)
		return result.ToValue(), err
	case INT16_FLAG:
		l := left.AsInt16()
		result, err := l.BitwiseXor(right)
		return result.ToValue(), err
	case INT8_FLAG:
		l := left.AsInt8()
		result, err := l.BitwiseXor(right)
		return result.ToValue(), err
	case UINT64_FLAG:
		l := left.AsUInt64()
		result, err := l.BitwiseXor(right)
		return result.ToValue(), err
	case UINT32_FLAG:
		l := left.AsUInt32()
		result, err := l.BitwiseXor(right)
		return result.ToValue(), err
	case UINT16_FLAG:
		l := left.AsUInt16()
		result, err := l.BitwiseXor(right)
		return result.ToValue(), err
	case UINT8_FLAG:
		l := left.AsUInt8()
		result, err := l.BitwiseXor(right)
		return result.ToValue(), err
	default:
		return Nil, Nil
	}
}
