package value

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"unsafe"

	"github.com/google/go-cmp/cmp"
)

var ValueClass *Class // ::Std::Value

type iface struct {
	tab uintptr
	ptr unsafe.Pointer
}

const ValueSize = unsafe.Sizeof(Value{})

// `undefined` is the zero value of `Value`, it maps directly to Go `nil`
type Value struct {
	data uintptr
	ptr  unsafe.Pointer
	flag uint8
}

func MakeSentinelValue() Value {
	return Value{flag: SENTINEL_FLAG}
}

const (
	UNDEFINED_FLAG = iota
	TRUE_FLAG
	FALSE_FLAG
	NIL_FLAG
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
	DATE_FLAG

	// only 64 bit systems
	INT64_FLAG
	UINT64_FLAG
	FLOAT64_FLAG
	TIME_FLAG
	TIME_SPAN_FLAG
	DATE_SPAN_FLAG
	REFERENCE_FLAG

	SENTINEL_FLAG = 0xFF
)

// Convert a Reference to a Value
func Ref(ref Reference) Value {
	i := *(*iface)(unsafe.Pointer(&ref))

	return Value{
		data: i.tab,
		ptr:  i.ptr,
		flag: REFERENCE_FLAG,
	}
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
		return v.AsInlineSymbol().Inspect()
	case FLOAT64_FLAG:
		return v.AsInlineFloat64().Inspect()
	case FLOAT32_FLAG:
		return v.AsFloat32().Inspect()
	case INT8_FLAG:
		return v.AsInt8().Inspect()
	case INT16_FLAG:
		return v.AsInt16().Inspect()
	case INT32_FLAG:
		return v.AsInt32().Inspect()
	case INT64_FLAG:
		return v.AsInlineInt64().Inspect()
	case UINT8_FLAG:
		return v.AsUInt8().Inspect()
	case UINT16_FLAG:
		return v.AsUInt16().Inspect()
	case UINT32_FLAG:
		return v.AsUInt32().Inspect()
	case UINT64_FLAG:
		return v.AsInlineUInt64().Inspect()
	case CHAR_FLAG:
		return v.AsChar().Inspect()
	case TIME_SPAN_FLAG:
		return v.AsInlineTimeSpan().Inspect()
	case DATE_FLAG:
		return v.AsDate().Inspect()
	case DATE_SPAN_FLAG:
		return v.AsDateSpan().Inspect()
	case TIME_FLAG:
		return v.AsTime().Inspect()
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
		return v.AsInlineSymbol().Class()
	case FLOAT64_FLAG:
		return v.AsInlineFloat64().Class()
	case FLOAT32_FLAG:
		return v.AsFloat32().Class()
	case INT8_FLAG:
		return v.AsInt8().Class()
	case INT16_FLAG:
		return v.AsInt16().Class()
	case INT32_FLAG:
		return v.AsInt32().Class()
	case INT64_FLAG:
		return v.AsInlineInt64().Class()
	case UINT8_FLAG:
		return v.AsUInt8().Class()
	case UINT16_FLAG:
		return v.AsUInt16().Class()
	case UINT32_FLAG:
		return v.AsUInt32().Class()
	case UINT64_FLAG:
		return v.AsInlineUInt64().Class()
	case CHAR_FLAG:
		return v.AsChar().Class()
	case TIME_SPAN_FLAG:
		return v.AsInlineTimeSpan().Class()
	case DATE_FLAG:
		return v.AsDate().Class()
	case DATE_SPAN_FLAG:
		return v.AsDateSpan().Class()
	case TIME_FLAG:
		return v.AsTime().Class()
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
		return v.AsInlineSymbol().DirectClass()
	case FLOAT64_FLAG:
		return v.AsInlineFloat64().DirectClass()
	case FLOAT32_FLAG:
		return v.AsFloat32().DirectClass()
	case INT8_FLAG:
		return v.AsInt8().DirectClass()
	case INT16_FLAG:
		return v.AsInt16().DirectClass()
	case INT32_FLAG:
		return v.AsInt32().DirectClass()
	case INT64_FLAG:
		return v.AsInlineInt64().DirectClass()
	case UINT8_FLAG:
		return v.AsUInt8().DirectClass()
	case UINT16_FLAG:
		return v.AsUInt16().DirectClass()
	case UINT32_FLAG:
		return v.AsUInt32().DirectClass()
	case UINT64_FLAG:
		return v.AsInlineUInt64().DirectClass()
	case CHAR_FLAG:
		return v.AsChar().DirectClass()
	case TIME_SPAN_FLAG:
		return v.AsInlineTimeSpan().DirectClass()
	case DATE_FLAG:
		return v.AsDate().DirectClass()
	case DATE_SPAN_FLAG:
		return v.AsDateSpan().DirectClass()
	case TIME_FLAG:
		return v.AsTime().DirectClass()
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
		return v.AsInlineSymbol().SingletonClass()
	case FLOAT64_FLAG:
		return v.AsInlineFloat64().SingletonClass()
	case FLOAT32_FLAG:
		return v.AsFloat32().SingletonClass()
	case INT8_FLAG:
		return v.AsInt8().SingletonClass()
	case INT16_FLAG:
		return v.AsInt16().SingletonClass()
	case INT32_FLAG:
		return v.AsInt32().SingletonClass()
	case INT64_FLAG:
		return v.AsInlineInt64().SingletonClass()
	case UINT8_FLAG:
		return v.AsUInt8().SingletonClass()
	case UINT16_FLAG:
		return v.AsUInt16().SingletonClass()
	case UINT32_FLAG:
		return v.AsUInt32().SingletonClass()
	case UINT64_FLAG:
		return v.AsInlineUInt64().SingletonClass()
	case CHAR_FLAG:
		return v.AsChar().SingletonClass()
	case TIME_SPAN_FLAG:
		return v.AsInlineTimeSpan().SingletonClass()
	case DATE_FLAG:
		return v.AsDate().SingletonClass()
	case DATE_SPAN_FLAG:
		return v.AsDateSpan().SingletonClass()
	case TIME_FLAG:
		return v.AsTime().SingletonClass()
	default:
		panic(fmt.Sprintf("invalid inline value flag: %d", v.ValueFlag()))
	}
}

func (v Value) InstanceVariables() *InstanceVariables {
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
		return v.AsInlineSymbol().InstanceVariables()
	case FLOAT64_FLAG:
		return v.AsInlineFloat64().InstanceVariables()
	case FLOAT32_FLAG:
		return v.AsFloat32().InstanceVariables()
	case INT8_FLAG:
		return v.AsInt8().InstanceVariables()
	case INT16_FLAG:
		return v.AsInt16().InstanceVariables()
	case INT32_FLAG:
		return v.AsInt32().InstanceVariables()
	case INT64_FLAG:
		return v.AsInlineInt64().InstanceVariables()
	case UINT8_FLAG:
		return v.AsUInt8().InstanceVariables()
	case UINT16_FLAG:
		return v.AsUInt16().InstanceVariables()
	case UINT32_FLAG:
		return v.AsUInt32().InstanceVariables()
	case UINT64_FLAG:
		return v.AsInlineUInt64().InstanceVariables()
	case CHAR_FLAG:
		return v.AsChar().InstanceVariables()
	case TIME_SPAN_FLAG:
		return v.AsInlineTimeSpan().InstanceVariables()
	case DATE_FLAG:
		return v.AsDate().InstanceVariables()
	case DATE_SPAN_FLAG:
		return v.AsDateSpan().InstanceVariables()
	case TIME_FLAG:
		return v.AsTime().InstanceVariables()
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
		return v.AsInlineSymbol().Error()
	case FLOAT64_FLAG:
		return v.AsInlineFloat64().Error()
	case FLOAT32_FLAG:
		return v.AsFloat32().Error()
	case INT8_FLAG:
		return v.AsInt8().Error()
	case INT16_FLAG:
		return v.AsInt16().Error()
	case INT32_FLAG:
		return v.AsInt32().Error()
	case INT64_FLAG:
		return v.AsInlineInt64().Error()
	case UINT8_FLAG:
		return v.AsUInt8().Error()
	case UINT16_FLAG:
		return v.AsUInt16().Error()
	case UINT32_FLAG:
		return v.AsUInt32().Error()
	case UINT64_FLAG:
		return v.AsInlineUInt64().Error()
	case CHAR_FLAG:
		return v.AsChar().Error()
	case TIME_SPAN_FLAG:
		return v.AsInlineTimeSpan().Error()
	case DATE_FLAG:
		return v.AsDate().Error()
	case DATE_SPAN_FLAG:
		return v.AsDateSpan().Error()
	case TIME_FLAG:
		return v.AsTime().Error()
	default:
		panic(fmt.Sprintf("invalid inline value flag: %d", v.ValueFlag()))
	}
}

func (v Value) ValueFlag() uint8 {
	return v.flag
}

func (v Value) IsReference() bool {
	return v.flag == REFERENCE_FLAG
}

func (v Value) AsString() String {
	return v.AsReference().(String)
}

func (v Value) AsBigInt() *BigInt {
	return (*BigInt)(v.Pointer())
}

func (v Value) MustBigInt() *BigInt {
	return v.MustReference().(*BigInt)
}

// Returns `nil` when the value is not a reference
func (v Value) SafeAsReference() Reference {
	if !v.IsReference() {
		return nil
	}
	return v.AsReference()
}

func (v Value) AsReference() Reference {
	i := iface{
		tab: v.data,
		ptr: v.ptr,
	}
	return *(*Reference)(unsafe.Pointer(&i))
}

func (v Value) Pointer() unsafe.Pointer {
	return v.ptr
}

func (v Value) MustReference() Reference {
	if !v.IsReference() {
		panic(fmt.Sprintf("value `%s` is not a reference", v.Inspect()))
	}
	return v.AsReference()
}

func RefErr(ref Reference, err Value) (Value, Value) {
	if !err.IsUndefined() {
		return Undefined, err
	}
	return Ref(ref), Undefined
}

type ToValuer interface {
	ToValue() Value
}

func ToValueErr[T ToValuer](t T, err Value) (Value, Value) {
	if !err.IsUndefined() {
		return Undefined, err
	}
	return t.ToValue(), Undefined
}

func (v Value) IsInlineValue() bool {
	return v.flag != REFERENCE_FLAG
}

func (v Value) IsSmallInt() bool {
	return v.flag == SMALL_INT_FLAG
}

func (v Value) AsSmallInt() SmallInt {
	return SmallInt(v.data)
}

func (v Value) MustSmallInt() SmallInt {
	if !v.IsSmallInt() {
		panic(fmt.Sprintf("value `%s` is not a SmallInt", v.Inspect()))
	}
	return v.AsSmallInt()
}

func (v Value) AsInt() int {
	if v.IsReference() {
		return int(v.AsBigInt().ToSmallInt())
	}

	return int(v.AsSmallInt())
}

func (v Value) IsChar() bool {
	return v.flag == CHAR_FLAG
}

func (v Value) AsChar() Char {
	return Char(v.data)
}

func (v Value) MustChar() Char {
	if !v.IsChar() {
		panic(fmt.Sprintf("value `%s` is not a Char", v.Inspect()))
	}
	return v.AsChar()
}

func (v Value) IsFloat() bool {
	return v.flag == FLOAT_FLAG
}

func (v Value) AsFloat() Float {
	return *(*Float)(unsafe.Pointer(&v.data))
}

func (v Value) MustFloat() Float {
	if !v.IsFloat() {
		panic(fmt.Sprintf("value `%s` is not a Float", v.Inspect()))
	}
	return v.AsFloat()
}

func (v Value) IsFloat32() bool {
	return v.flag == FLOAT32_FLAG
}

func (v Value) AsFloat32() Float32 {
	return Float32(math.Float32frombits(uint32(v.data)))
}

func (v Value) MustFloat32() Float32 {
	if !v.IsFloat32() {
		panic(fmt.Sprintf("value `%s` is not a Float32", v.Inspect()))
	}
	return v.AsFloat32()
}

func (v Value) IsInlineFloat64() bool {
	return v.flag == FLOAT64_FLAG
}

func (v Value) AsInlineFloat64() Float64 {
	return Float64(math.Float64frombits(uint64(v.data)))
}

func (v Value) MustInlineFloat64() Float64 {
	if !v.IsInlineFloat64() {
		panic(fmt.Sprintf("value `%s` is not an inline Float64", v.Inspect()))
	}
	return v.AsInlineFloat64()
}

func (v Value) AsFloat64() Float64 {
	if v.IsReference() {
		return v.AsReference().(Float64)
	} else {
		return v.AsInlineFloat64()
	}
}

func (v Value) MustFloat64() Float64 {
	if v.IsReference() {
		return v.AsReference().(Float64)
	} else {
		return v.MustInlineFloat64()
	}
}

func (v Value) IsInt8() bool {
	return v.flag == INT8_FLAG
}

func (v Value) AsInt8() Int8 {
	return Int8(v.data)
}

func (v Value) MustInt8() Int8 {
	if !v.IsInt8() {
		panic(fmt.Sprintf("value `%s` is not a Int8", v.Inspect()))
	}
	return v.AsInt8()
}

func (v Value) IsUInt8() bool {
	return v.flag == UINT8_FLAG
}

func (v Value) AsUInt8() UInt8 {
	return UInt8(v.data)
}

func (v Value) MustUInt8() UInt8 {
	if !v.IsUInt8() {
		panic(fmt.Sprintf("value `%s` is not a UInt8", v.Inspect()))
	}
	return v.AsUInt8()
}

func (v Value) IsInt16() bool {
	return v.flag == INT16_FLAG
}

func (v Value) AsInt16() Int16 {
	return Int16(v.data)
}

func (v Value) MustInt16() Int16 {
	if !v.IsInt16() {
		panic(fmt.Sprintf("value `%s` is not a Int16", v.Inspect()))
	}
	return v.AsInt16()
}

func (v Value) IsUInt16() bool {
	return v.flag == UINT16_FLAG
}

func (v Value) AsUInt16() UInt16 {
	return UInt16(v.data)
}

func (v Value) MustUInt16() UInt16 {
	if !v.IsUInt16() {
		panic(fmt.Sprintf("value `%s` is not a UInt16", v.Inspect()))
	}
	return v.AsUInt16()
}

func (v Value) IsInt32() bool {
	return v.flag == INT32_FLAG
}

func (v Value) AsInt32() Int32 {
	return Int32(v.data)
}

func (v Value) MustInt32() Int32 {
	if !v.IsInt32() {
		panic(fmt.Sprintf("value `%s` is not a Int32", v.Inspect()))
	}
	return v.AsInt32()
}

func (v Value) IsUInt32() bool {
	return v.flag == UINT32_FLAG
}

func (v Value) AsUInt32() UInt32 {
	return UInt32(v.data)
}

func (v Value) MustUInt32() UInt32 {
	if !v.IsUInt32() {
		panic(fmt.Sprintf("value `%s` is not a UInt32", v.Inspect()))
	}
	return v.AsUInt32()
}

func (v Value) IsInlineInt64() bool {
	return v.flag == INT64_FLAG
}

func (v Value) AsInlineInt64() Int64 {
	return Int64(v.data)
}

func (v Value) MustInlineInt64() Int64 {
	if !v.IsInlineInt64() {
		panic(fmt.Sprintf("value `%s` is not an inline Int64", v.Inspect()))
	}
	return v.AsInlineInt64()
}

func (v Value) AsInt64() Int64 {
	if v.IsReference() {
		return v.AsReference().(Int64)
	} else {
		return v.AsInlineInt64()
	}
}

func (v Value) MustInt64() Int64 {
	if v.IsReference() {
		return v.AsReference().(Int64)
	} else {
		return v.MustInlineInt64()
	}
}

func (v Value) IsInlineUInt64() bool {
	return v.flag == UINT64_FLAG
}

func (v Value) AsInlineUInt64() UInt64 {
	return UInt64(v.data)
}

func (v Value) MustInlineUInt64() UInt64 {
	if !v.IsInlineUInt64() {
		panic(fmt.Sprintf("value `%s` is not an inline UInt64", v.Inspect()))
	}
	return v.AsInlineUInt64()
}

func (v Value) AsUInt64() UInt64 {
	if v.IsReference() {
		return v.AsReference().(UInt64)
	} else {
		return v.AsInlineUInt64()
	}
}

func (v Value) MustUInt64() UInt64 {
	if v.IsReference() {
		return v.AsReference().(UInt64)
	} else {
		return v.MustInlineUInt64()
	}
}

func (v Value) IsInlineTimeSpan() bool {
	return v.flag == TIME_SPAN_FLAG
}

func (v Value) AsInlineTimeSpan() TimeSpan {
	return TimeSpan(v.data)
}

func (v Value) AsTimeSpan() TimeSpan {
	if v.IsReference() {
		return v.AsReference().(TimeSpan)
	} else {
		return v.AsInlineTimeSpan()
	}
}

func (v Value) MustTimeSpan() TimeSpan {
	if v.IsReference() {
		return v.AsReference().(TimeSpan)
	} else {
		return v.MustInlineTimeSpan()
	}
}

func (v Value) MustInlineTimeSpan() TimeSpan {
	if !v.IsInlineTimeSpan() {
		panic(fmt.Sprintf("value `%s` is not an inline Time::Span", v.Inspect()))
	}
	return v.AsInlineTimeSpan()
}

func (v Value) IsDate() bool {
	return v.flag == DATE_FLAG
}

func (v Value) AsDate() Date {
	return Date{bits: uint32(v.data)}
}

func (v Value) MustDate() Date {
	if !v.IsDate() {
		panic(fmt.Sprintf("value `%s` is not a Date", v.Inspect()))
	}
	return v.AsDate()
}

func (v Value) IsInlineTime() bool {
	return v.flag == TIME_FLAG
}

func (v Value) AsInlineTime() Time {
	return Time{
		duration: TimeSpan(v.data),
	}
}

func (v Value) MustInlineTime() Time {
	if !v.IsInlineTime() {
		panic(fmt.Sprintf("value `%s` is not an inline Time", v.Inspect()))
	}
	return v.AsInlineTime()
}

func (v Value) AsTime() Time {
	if v.IsReference() {
		return v.AsReference().(Time)
	} else {
		return v.AsInlineTime()
	}
}

func (v Value) MustTime() Time {
	if v.IsReference() {
		return v.AsReference().(Time)
	} else {
		return v.MustInlineTime()
	}
}

func (v Value) IsInlineDateSpan() bool {
	return v.flag == DATE_SPAN_FLAG
}

func (v Value) AsInlineDateSpan() DateSpan {
	months := int32(v.data >> 32)
	days := int32(v.data & 0xFFFFFFFF)
	return DateSpan{months: months, days: days}
}

func (v Value) AsDateSpan() DateSpan {
	if v.IsReference() {
		return v.AsReference().(DateSpan)
	} else {
		return v.AsInlineDateSpan()
	}
}

func (v Value) MustDateSpan() DateSpan {
	if v.IsReference() {
		return v.AsReference().(DateSpan)
	} else {
		return v.MustInlineDateSpan()
	}
}

func (v Value) MustInlineDateSpan() DateSpan {
	if !v.IsInlineTimeSpan() {
		panic(fmt.Sprintf("value `%s` is not an inline Date::Span", v.Inspect()))
	}
	return v.AsInlineDateSpan()
}

func (v Value) IsInlineSymbol() bool {
	return v.flag == SYMBOL_FLAG
}

func (v Value) AsInlineSymbol() Symbol {
	return Symbol(v.data)
}

func (v Value) MustInlineSymbol() Symbol {
	if !v.IsInlineSymbol() {
		panic(fmt.Sprintf("value `%s` is not an inline Symbol", v.Inspect()))
	}
	return v.AsInlineSymbol()
}

func (v Value) IsTrue() bool {
	return v.flag == TRUE_FLAG
}

func (v Value) AsTrue() TrueType {
	return *(*TrueType)(unsafe.Pointer(&v.data))
}

func (v Value) MustTrue() TrueType {
	if !v.IsTrue() {
		panic(fmt.Sprintf("value `%s` is not True", v.Inspect()))
	}
	return v.AsTrue()
}

func (v Value) IsFalse() bool {
	return v.flag == FALSE_FLAG
}

func (v Value) AsFalse() FalseType {
	return *(*FalseType)(unsafe.Pointer(&v.data))
}

func (v Value) MustFalse() FalseType {
	if !v.IsFalse() {
		panic(fmt.Sprintf("value `%s` is not False", v.Inspect()))
	}
	return v.AsFalse()
}

func (v Value) IsNil() bool {
	return v.flag == NIL_FLAG
}

func (v Value) AsNil() NilType {
	return *(*NilType)(unsafe.Pointer(&v.data))
}

func (v Value) MustNil() NilType {
	if !v.IsNil() {
		panic(fmt.Sprintf("value `%s` is not Nil", v.Inspect()))
	}
	return v.AsNil()
}

func (v Value) IsUndefined() bool {
	return v.flag == UNDEFINED_FLAG
}

func (v Value) AsUndefined() UndefinedType {
	return *(*UndefinedType)(unsafe.Pointer(&v.data))
}

func (v Value) MustUndefined() UndefinedType {
	if !v.IsUndefined() {
		panic(fmt.Sprintf("value `%s` is not Undefined", v.Inspect()))
	}
	return v.AsUndefined()
}

// Set an object's instance variable with the given name to the given value
func SetInstanceVariableByName(object Value, name Symbol, val Value) (err Value) {
	class := object.DirectClass()
	ivars := object.InstanceVariables()
	if ivars == nil {
		return Ref(NewCantSetInstanceVariablesOnPrimitiveError(object.Inspect()))
	}

	ivarIndex := class.IvarIndices[name]
	ivars.Set(ivarIndex, val)
	return Undefined
}

// Get an object's instance variable with the given name
func GetInstanceVariableByName(object Value, name Symbol) (val, err Value) {
	class := object.DirectClass()
	ivars := object.InstanceVariables()
	if ivars == nil {
		return Undefined, Ref(NewCantSetInstanceVariablesOnPrimitiveError(object.Inspect()))
	}

	ivarIndex := class.IvarIndices[name]
	val = ivars.Get(ivarIndex)
	return val, Undefined
}

// Elk Reference Value
type Reference interface {
	Class() *Class                         // Return the class of the value
	DirectClass() *Class                   // Return the direct class of this value that will be searched for methods first
	SingletonClass() *Class                // Return the singleton class of this value that holds methods unique to this object
	InstanceVariables() *InstanceVariables // Returns a pointer to the slice of instance vars of this value, nil if value doesn't support instance vars
	Copy() Reference                       // Creates a shallow copy of the reference. If the value is immutable, no copying should be done, the same value should be returned.
	Inspect() string                       // Returns the string representation of the value
	Error() string                         // Implements the error interface
}

func IsMutableCollection(val Value) bool {
	if val.IsInlineValue() {
		return false
	}
	switch v := val.AsReference().(type) {
	case *ArrayList, *HashMap:
		return true
	case *ArrayTuple:
		if slices.ContainsFunc(*v, IsMutableCollection) {
			return true
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
		return int(val.AsInlineInt64()), true
	case UINT8_FLAG:
		return int(val.AsUInt8()), true
	case UINT16_FLAG:
		return int(val.AsUInt16()), true
	case UINT32_FLAG:
		return int(val.AsUInt32()), true
	case UINT64_FLAG:
		return int(val.AsInlineUInt64()), true
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
		v := val.AsInlineInt64()
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
		v := val.AsInlineUInt64()
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
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func SubscriptVal(collection, key Value) (result, err Value) {
	if !collection.IsReference() {
		return Undefined, Undefined
	}

	switch l := collection.AsReference().(type) {
	case *ArrayTuple:
		return l.Subscript(key)
	case *ArrayList:
		return l.Subscript(key)
	default:
		return Undefined, Undefined
	}
}

// Set an element under the given key.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func SubscriptSet(collection, key, val Value) (err Value) {
	if !collection.IsReference() {
		return Undefined
	}

	switch l := collection.AsReference().(type) {
	case *ArrayList:
		return l.SubscriptSet(key, val)
	case *ArrayTuple:
		return l.SubscriptSet(key, val)
	default:
		return Undefined
	}
}

// Calculate the hash of the value.
// When successful returns (result, undefined).
// When an error occurred returns (0, error).
// When there are no builtin addition functions for the given type returns (0, NotBuiltinError).
func Hash(key Value) (result UInt64, err Value) {
	if key.IsReference() {
		switch k := key.AsReference().(type) {
		case String:
			return k.Hash(), err
		case *BigInt:
			return k.Hash(), err
		case *BigFloat:
			return k.Hash(), err
		case Float64:
			return k.Hash(), err
		case Int64:
			return k.Hash(), err
		case UInt64:
			return k.Hash(), err
		default:
			return 0, Ref(NotBuiltinError)
		}
	}

	switch key.ValueFlag() {
	case CHAR_FLAG:
		k := key.AsChar()
		return k.Hash(), err
	case SYMBOL_FLAG:
		k := key.AsInlineSymbol()
		return k.Hash(), err
	case SMALL_INT_FLAG:
		k := key.AsSmallInt()
		return k.Hash(), err
	case FLOAT_FLAG:
		k := key.AsFloat()
		return k.Hash(), err
	case NIL_FLAG:
		k := key.AsNil()
		return k.Hash(), err
	case TRUE_FLAG:
		k := key.AsTrue()
		return k.Hash(), err
	case FALSE_FLAG:
		k := key.AsFalse()
		return k.Hash(), err
	case FLOAT64_FLAG:
		k := key.AsInlineFloat64()
		return k.Hash(), err
	case FLOAT32_FLAG:
		k := key.AsFloat32()
		return k.Hash(), err
	case INT64_FLAG:
		k := key.AsInlineInt64()
		return k.Hash(), err
	case INT32_FLAG:
		k := key.AsInt32()
		return k.Hash(), err
	case INT16_FLAG:
		k := key.AsInt16()
		return k.Hash(), err
	case INT8_FLAG:
		k := key.AsInt8()
		return k.Hash(), err
	case UINT64_FLAG:
		k := key.AsInlineUInt64()
		return k.Hash(), err
	case UINT32_FLAG:
		k := key.AsUInt32()
		return k.Hash(), err
	case UINT16_FLAG:
		k := key.AsUInt16()
		return k.Hash(), err
	case UINT8_FLAG:
		k := key.AsUInt8()
		return k.Hash(), err
	default:
		return 0, Ref(NotBuiltinError)
	}
}

// AddVal two values.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func AddVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.AddVal(right)
		case *BigFloat:
			return l.AddVal(right)
		case Float64:
			return ToValueErr(l.Add(right))
		case Int64:
			return ToValueErr(l.Add(right))
		case UInt64:
			return ToValueErr(l.Add(right))
		case String:
			return RefErr(l.Concat(right))
		case *Regex:
			return l.ConcatVal(right)
		case *ArrayList:
			return RefErr(l.Concat(right))
		case *ArrayTuple:
			return l.ConcatVal(right)
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.AddVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.AddVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return ToValueErr(l.Add(right))
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return ToValueErr(l.Add(right))
	case INT64_FLAG:
		l := left.AsInlineInt64()
		return ToValueErr(l.Add(right))
	case INT32_FLAG:
		l := left.AsInt32()
		return ToValueErr(l.Add(right))
	case INT16_FLAG:
		l := left.AsInt16()
		return ToValueErr(l.Add(right))
	case INT8_FLAG:
		l := left.AsInt8()
		return ToValueErr(l.Add(right))
	case UINT64_FLAG:
		l := left.AsInlineUInt64()
		return ToValueErr(l.Add(right))
	case UINT32_FLAG:
		l := left.AsUInt32()
		return ToValueErr(l.Add(right))
	case UINT16_FLAG:
		l := left.AsUInt16()
		return ToValueErr(l.Add(right))
	case UINT8_FLAG:
		l := left.AsUInt8()
		return ToValueErr(l.Add(right))
	case CHAR_FLAG:
		l := left.AsChar()
		return RefErr(l.Concat(right))
	default:
		return Undefined, Undefined
	}
}

// SubtractVal two values
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func SubtractVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.SubtractVal(right)
		case *BigFloat:
			return l.SubtractVal(right)
		case String:
			return RefErr(l.RemoveSuffix(right))
		case Float64:
			return ToValueErr(l.Subtract(right))
		case Int64:
			return ToValueErr(l.Subtract(right))
		case UInt64:
			return ToValueErr(l.Subtract(right))
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.SubtractVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.SubtractVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return ToValueErr(l.Subtract(right))
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return ToValueErr(l.Subtract(right))
	case INT64_FLAG:
		l := left.AsInlineInt64()
		return ToValueErr(l.Subtract(right))
	case INT32_FLAG:
		l := left.AsInt32()
		return ToValueErr(l.Subtract(right))
	case INT16_FLAG:
		l := left.AsInt16()
		return ToValueErr(l.Subtract(right))
	case INT8_FLAG:
		l := left.AsInt8()
		return ToValueErr(l.Subtract(right))
	case UINT64_FLAG:
		l := left.AsInlineUInt64()
		return ToValueErr(l.Subtract(right))
	case UINT32_FLAG:
		l := left.AsUInt32()
		return ToValueErr(l.Subtract(right))
	case UINT16_FLAG:
		l := left.AsUInt16()
		return ToValueErr(l.Subtract(right))
	case UINT8_FLAG:
		l := left.AsUInt8()
		return ToValueErr(l.Subtract(right))
	default:
		return Undefined, Undefined
	}
}

// MultiplyVal two values
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func MultiplyVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.MultiplyVal(right)
		case *BigFloat:
			return l.MultiplyVal(right)
		case Float64:
			return ToValueErr(l.Multiply(right))
		case Int64:
			return ToValueErr(l.Multiply(right))
		case UInt64:
			return ToValueErr(l.Multiply(right))
		case String:
			return RefErr(l.Repeat(right))
		case *Regex:
			return l.RepeatVal(right)
		case *ArrayList:
			return RefErr(l.Repeat(right))
		case *ArrayTuple:
			return RefErr(l.Repeat(right))
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.MultiplyVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.MultiplyVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return ToValueErr(l.Multiply(right))
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return ToValueErr(l.Multiply(right))
	case INT64_FLAG:
		l := left.AsInlineInt64()
		return ToValueErr(l.Multiply(right))
	case INT32_FLAG:
		l := left.AsInt32()
		return ToValueErr(l.Multiply(right))
	case INT16_FLAG:
		l := left.AsInt16()
		return ToValueErr(l.Multiply(right))
	case INT8_FLAG:
		l := left.AsInt8()
		return ToValueErr(l.Multiply(right))
	case UINT64_FLAG:
		l := left.AsInlineUInt64()
		return ToValueErr(l.Multiply(right))
	case UINT32_FLAG:
		l := left.AsUInt32()
		return ToValueErr(l.Multiply(right))
	case UINT16_FLAG:
		l := left.AsUInt16()
		return ToValueErr(l.Multiply(right))
	case UINT8_FLAG:
		l := left.AsUInt8()
		return ToValueErr(l.Multiply(right))
	case CHAR_FLAG:
		l := left.AsChar()
		return RefErr(l.Repeat(right))
	default:
		return Undefined, Undefined
	}
}

// DivideVal two values
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func DivideVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.DivideVal(right)
		case *BigFloat:
			return l.DivideVal(right)
		case Float64:
			return ToValueErr(l.Divide(right))
		case Int64:
			return ToValueErr(l.Divide(right))
		case UInt64:
			return ToValueErr(l.Divide(right))
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.DivideVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.DivideVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return ToValueErr(l.Divide(right))
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return ToValueErr(l.Divide(right))
	case INT64_FLAG:
		l := left.AsInlineInt64()
		return ToValueErr(l.Divide(right))
	case INT32_FLAG:
		l := left.AsInt32()
		return ToValueErr(l.Divide(right))
	case INT16_FLAG:
		l := left.AsInt16()
		return ToValueErr(l.Divide(right))
	case INT8_FLAG:
		l := left.AsInt8()
		return ToValueErr(l.Divide(right))
	case UINT64_FLAG:
		l := left.AsInlineUInt64()
		return ToValueErr(l.Divide(right))
	case UINT32_FLAG:
		l := left.AsUInt32()
		return ToValueErr(l.Divide(right))
	case UINT16_FLAG:
		l := left.AsUInt16()
		return ToValueErr(l.Divide(right))
	case UINT8_FLAG:
		l := left.AsUInt8()
		return ToValueErr(l.Divide(right))
	default:
		return Undefined, Undefined
	}
}

// NegateVal a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func NegateVal(operand Value) Value {
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
			return Undefined
		}
	}

	switch operand.ValueFlag() {
	case SMALL_INT_FLAG:
		o := operand.AsSmallInt()
		return o.NegateVal()
	case FLOAT_FLAG:
		o := operand.AsFloat()
		return (-o).ToValue()
	case FLOAT64_FLAG:
		o := operand.AsInlineFloat64()
		return (-o).ToValue()
	case FLOAT32_FLAG:
		o := operand.AsFloat32()
		return (-o).ToValue()
	case INT64_FLAG:
		o := operand.AsInlineInt64()
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
		o := operand.AsInlineUInt64()
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
		return Undefined
	}
}

// IncrementVal a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func IncrementVal(operand Value) Value {
	if operand.IsReference() {
		switch o := operand.AsReference().(type) {
		case *BigInt:
			return Ref(o.Increment())
		case Int64:
			return (o + 1).ToValue()
		case UInt64:
			return (o + 1).ToValue()
		default:
			return Undefined
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
		o := operand.AsInlineInt64()
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
		o := operand.AsInlineUInt64()
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
		return Undefined
	}
}

// DecrementVal a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func DecrementVal(operand Value) Value {
	if operand.IsReference() {
		switch o := operand.AsReference().(type) {
		case *BigInt:
			return o.DecrementVal()
		case Int64:
			return (o - 1).ToValue()
		case UInt64:
			return (o - 1).ToValue()
		default:
			return Undefined
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
		o := operand.AsInlineInt64()
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
		o := operand.AsInlineUInt64()
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
		return Undefined
	}
}

// Perform unary plus on a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func UnaryPlusVal(operand Value) Value {
	if operand.IsReference() {
		switch operand.AsReference().(type) {
		case *BigInt, *BigFloat,
			Float64, Int64,
			UInt64:
			return operand
		default:
			return Undefined
		}
	}

	switch operand.ValueFlag() {
	case SMALL_INT_FLAG, FLOAT_FLAG,
		FLOAT64_FLAG, FLOAT32_FLAG, INT64_FLAG, INT32_FLAG, INT16_FLAG, INT8_FLAG,
		UINT64_FLAG, UINT32_FLAG, UINT16_FLAG, UINT8_FLAG:
		return operand
	default:
		return Undefined
	}
}

// Perform bitwise not on a value
// When successful returns result.
// When there are no builtin negation functions for the given type returns nil.
func BitwiseNotVal(operand Value) Value {
	if operand.IsReference() {
		switch o := operand.AsReference().(type) {
		case *BigInt:
			return Ref(o.BitwiseNot())
		case Int64:
			return (^o).ToValue()
		case UInt64:
			return (^o).ToValue()
		default:
			return Undefined
		}
	}

	switch operand.ValueFlag() {
	case SMALL_INT_FLAG:
		o := operand.AsSmallInt()
		return (^o).ToValue()
	case INT64_FLAG:
		o := operand.AsInlineInt64()
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
		o := operand.AsInlineUInt64()
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
		return Undefined
	}
}

// ExponentiateVal two values
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func ExponentiateVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.ExponentiateVal(right)
		case *BigFloat:
			return l.ExponentiateVal(right)
		case Float64:
			return ToValueErr(l.ExponentiateVal(right))
		case Int64:
			return ToValueErr(l.ExponentiateVal(right))
		case UInt64:
			return ToValueErr(l.ExponentiateVal(right))
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.ExponentiateVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.ExponentiateVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return ToValueErr(l.ExponentiateVal(right))
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return ToValueErr(l.ExponentiateVal(right))
	case INT64_FLAG:
		l := left.AsInlineInt64()
		return ToValueErr(l.ExponentiateVal(right))
	case INT32_FLAG:
		l := left.AsInt32()
		return ToValueErr(l.ExponentiateVal(right))
	case INT16_FLAG:
		l := left.AsInt16()
		return ToValueErr(l.ExponentiateVal(right))
	case INT8_FLAG:
		l := left.AsInt8()
		return ToValueErr(l.ExponentiateVal(right))
	case UINT64_FLAG:
		l := left.AsInlineUInt64()
		return ToValueErr(l.ExponentiateVal(right))
	case UINT32_FLAG:
		l := left.AsUInt32()
		return ToValueErr(l.ExponentiateVal(right))
	case UINT16_FLAG:
		l := left.AsUInt16()
		return ToValueErr(l.ExponentiateVal(right))
	case UINT8_FLAG:
		l := left.AsUInt8()
		return ToValueErr(l.ExponentiateVal(right))
	default:
		return Undefined, Undefined
	}
}

// Perform modulo on two values
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func ModuloVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.ModuloVal(right)
		case *BigFloat:
			return l.ModuloVal(right)
		case Float64:
			return ToValueErr(l.ModuloVal(right))
		case Int64:
			return ToValueErr(l.ModuloVal(right))
		case UInt64:
			return ToValueErr(l.ModuloVal(right))
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.ModuloVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.ModuloVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return ToValueErr(l.ModuloVal(right))
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return ToValueErr(l.ModuloVal(right))
	case INT64_FLAG:
		l := left.AsInlineInt64()
		return ToValueErr(l.ModuloVal(right))
	case INT32_FLAG:
		l := left.AsInt32()
		return ToValueErr(l.ModuloVal(right))
	case INT16_FLAG:
		l := left.AsInt16()
		return ToValueErr(l.ModuloVal(right))
	case INT8_FLAG:
		l := left.AsInt8()
		return ToValueErr(l.ModuloVal(right))
	case UINT64_FLAG:
		l := left.AsInlineUInt64()
		return ToValueErr(l.ModuloVal(right))
	case UINT32_FLAG:
		l := left.AsUInt32()
		return ToValueErr(l.ModuloVal(right))
	case UINT16_FLAG:
		l := left.AsUInt16()
		return ToValueErr(l.ModuloVal(right))
	case UINT8_FLAG:
		l := left.AsUInt8()
		return ToValueErr(l.ModuloVal(right))
	default:
		return Undefined, Undefined
	}
}

// CompareVal two values.
// Returns 1 if left is greater than right.
// Returns 0 if both are equal.
// Returns -1 if left is less than right.
//
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func CompareVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.CompareVal(right)
		case *BigFloat:
			return l.CompareVal(right)
		case String:
			return l.CompareVal(right)
		case Float64:
			return l.CompareVal(right)
		case Int64:
			return l.CompareVal(right)
		case UInt64:
			return l.CompareVal(right)
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.CompareVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.CompareVal(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.CompareVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return l.CompareVal(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.CompareVal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
		return l.CompareVal(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.CompareVal(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.CompareVal(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.CompareVal(right)
	case UINT64_FLAG:
		l := left.AsInlineUInt64()
		return l.CompareVal(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.CompareVal(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.CompareVal(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.CompareVal(right)
	default:
		return Undefined, Undefined
	}
}

// Check whether left is greater than right.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func GreaterThanVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.GreaterThanVal(right)
		case *BigFloat:
			return l.GreaterThanVal(right)
		case String:
			return l.GreaterThanVal(right)
		case Float64:
			return l.GreaterThanVal(right)
		case Int64:
			return l.GreaterThanVal(right)
		case UInt64:
			return l.GreaterThanVal(right)
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.GreaterThanVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.GreaterThanVal(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.GreaterThanVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return l.GreaterThanVal(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.GreaterThanVal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
		return l.GreaterThanVal(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.GreaterThanVal(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.GreaterThanVal(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.GreaterThanVal(right)
	case UINT64_FLAG:
		l := left.AsInlineUInt64()
		return l.GreaterThanVal(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.GreaterThanVal(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.GreaterThanVal(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.GreaterThanVal(right)
	default:
		return Undefined, Undefined
	}
}

func GreaterThan(left, right Value) (result bool, err Value) {
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
			return false, Undefined
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
		l := left.AsInlineFloat64()
		return l.GreaterThan(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.GreaterThan(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return false, Undefined
	}
}

// Check whether left is greater than or equal to right.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func GreaterThanEqualVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.GreaterThanEqualVal(right)
		case *BigFloat:
			return l.GreaterThanEqualVal(right)
		case String:
			return l.GreaterThanEqualVal(right)
		case Float64:
			return l.GreaterThanEqualVal(right)
		case Int64:
			return l.GreaterThanEqualVal(right)
		case UInt64:
			return l.GreaterThanEqualVal(right)
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.GreaterThanEqualVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.GreaterThanEqualVal(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.GreaterThanEqualVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return l.GreaterThanEqualVal(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.GreaterThanEqualVal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
		return l.GreaterThanEqualVal(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.GreaterThanEqualVal(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.GreaterThanEqualVal(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.GreaterThanEqualVal(right)
	case UINT64_FLAG:
		l := left.AsInlineUInt64()
		return l.GreaterThanEqualVal(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.GreaterThanEqualVal(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.GreaterThanEqualVal(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.GreaterThanEqualVal(right)
	default:
		return Undefined, Undefined
	}
}

func GreaterThanEqual(left, right Value) (result bool, err Value) {
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
			return false, Undefined
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
		l := left.AsInlineFloat64()
		return l.GreaterThanEqual(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.GreaterThanEqual(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return false, Undefined
	}
}

// Check whether left is less than right.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func LessThanVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.LessThanVal(right)
		case *BigFloat:
			return l.LessThanVal(right)
		case String:
			return l.LessThanVal(right)
		case Float64:
			return l.LessThanVal(right)
		case Int64:
			return l.LessThanVal(right)
		case UInt64:
			return l.LessThanVal(right)
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.LessThanVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.LessThanVal(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.LessThanVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return l.LessThanVal(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.LessThanVal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
		return l.LessThanVal(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.LessThanVal(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.LessThanVal(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.LessThanVal(right)
	case UINT64_FLAG:
		l := left.AsInlineUInt64()
		return l.LessThanVal(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.LessThanVal(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.LessThanVal(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.LessThanVal(right)
	default:
		return Undefined, Undefined
	}
}

func LessThan(left, right Value) (result bool, err Value) {
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
			return false, Undefined
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
		l := left.AsInlineFloat64()
		return l.LessThan(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.LessThan(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return false, Undefined
	}
}

// Check whether left is less than or equal to right.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin functions for the given type returns (undefined, undefined).
func LessThanEqualVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.LessThanEqualVal(right)
		case *BigFloat:
			return l.LessThanEqualVal(right)
		case String:
			return l.LessThanEqualVal(right)
		case Float64:
			return l.LessThanEqualVal(right)
		case Int64:
			return l.LessThanEqualVal(right)
		case UInt64:
			return l.LessThanEqualVal(right)
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.LessThanEqualVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.LessThanEqualVal(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.LessThanEqualVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return l.LessThanEqualVal(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.LessThanEqualVal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
		return l.LessThanEqualVal(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.LessThanEqualVal(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.LessThanEqualVal(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.LessThanEqualVal(right)
	case UINT64_FLAG:
		l := left.AsInlineUInt64()
		return l.LessThanEqualVal(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.LessThanEqualVal(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.LessThanEqualVal(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.LessThanEqualVal(right)
	default:
		return Undefined, Undefined
	}
}

func LessThanEqual(left, right Value) (result bool, err Value) {
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
			return false, Undefined
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
		l := left.AsInlineFloat64()
		return l.LessThanEqual(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.LessThanEqual(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return false, Undefined
	}
}

// Check whether left is equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func LaxEqualVal(left, right Value) Value {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.LaxEqualVal(right)
		case *BigFloat:
			return l.LaxEqualVal(right)
		case String:
			return l.LaxEqualVal(right)
		case *Regex:
			return l.LaxEqualVal(right)
		case Float64:
			return StrictFloatLaxEqual(l, right)
		case Int64:
			return StrictSignedIntLaxEqual(l, right)
		case UInt64:
			return StrictUnsignedIntLaxEqual(l, right)
		default:
			return Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.LaxEqualVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.LaxEqualVal(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.LaxEqualVal(right)
	case SYMBOL_FLAG:
		l := left.AsInlineSymbol()
		return l.LaxEqualVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return StrictFloatLaxEqual(l, right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return StrictFloatLaxEqual(l, right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return Undefined
	}
}

// Check whether left is not equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func LaxNotEqualVal(left, right Value) Value {
	val := LaxEqualVal(left, right)
	if val.IsUndefined() {
		return Undefined
	}

	return ToNotBool(val)
}

// Check whether left is equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (undefined).
func EqualVal(left, right Value) Value {
	class := left.Class()
	if !IsA(right, class) {
		return False
	}

	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.EqualVal(right)
		case *BigFloat:
			return l.EqualVal(right)
		case String:
			return l.EqualVal(right)
		case Float64:
			return l.EqualVal(right)
		case Int64:
			return l.EqualVal(right)
		case UInt64:
			return l.EqualVal(right)
		default:
			return Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.EqualVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.EqualVal(right)
	case SYMBOL_FLAG:
		l := left.AsInlineSymbol()
		return l.EqualVal(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.EqualVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return l.EqualVal(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.EqualVal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
		return l.EqualVal(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.EqualVal(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.EqualVal(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.EqualVal(right)
	case UINT64_FLAG:
		l := left.AsInlineUInt64()
		return l.EqualVal(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.EqualVal(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.EqualVal(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.EqualVal(right)
	default:
		return Undefined
	}
}

func Equal(left, right Value) bool {
	class := left.Class()
	if !IsA(right, class) {
		return false
	}

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
			return false
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
		l := left.AsInlineSymbol()
		return l.Equal(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.Equal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return l.Equal(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.Equal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return false
	}
}

// Check whether left is not equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func NotEqualVal(left, right Value) Value {
	val := EqualVal(left, right)
	if val.IsUndefined() {
		return Undefined
	}

	return ToNotBool(val)
}

// Check whether left is strictly equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func StrictEqualVal(left, right Value) Value {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.StrictEqualVal(right)
		case *BigFloat:
			return l.StrictEqualVal(right)
		case String:
			return l.StrictEqualVal(right)
		case Float64:
			return l.StrictEqualVal(right)
		case Int64:
			return l.StrictEqualVal(right)
		case UInt64:
			return l.StrictEqualVal(right)
		default:
			return ToElkBool(left == right)
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.StrictEqualVal(right)
	case FLOAT_FLAG:
		l := left.AsFloat()
		return l.StrictEqualVal(right)
	case SYMBOL_FLAG:
		l := left.AsInlineSymbol()
		return l.StrictEqualVal(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.StrictEqualVal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return l.StrictEqualVal(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.StrictEqualVal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
		return l.StrictEqualVal(right)
	case INT32_FLAG:
		l := left.AsInt32()
		return l.StrictEqualVal(right)
	case INT16_FLAG:
		l := left.AsInt16()
		return l.StrictEqualVal(right)
	case INT8_FLAG:
		l := left.AsInt8()
		return l.StrictEqualVal(right)
	case UINT64_FLAG:
		l := left.AsInlineUInt64()
		return l.StrictEqualVal(right)
	case UINT32_FLAG:
		l := left.AsUInt32()
		return l.StrictEqualVal(right)
	case UINT16_FLAG:
		l := left.AsUInt16()
		return l.StrictEqualVal(right)
	case UINT8_FLAG:
		l := left.AsUInt8()
		return l.StrictEqualVal(right)
	default:
		return ToElkBool(left == right)
	}
}

func StrictEqual(left, right Value) bool {
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
			return left == right
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
		l := left.AsInlineSymbol()
		return l.Equal(right)
	case CHAR_FLAG:
		l := left.AsChar()
		return l.Equal(right)
	case FLOAT64_FLAG:
		l := left.AsInlineFloat64()
		return l.Equal(right)
	case FLOAT32_FLAG:
		l := left.AsFloat32()
		return l.Equal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return left == right
	}
}

// Check whether left is strictly not equal to right.
// When successful returns (result).
// When there are no builtin addition functions for the given type returns (nil).
func StrictNotEqualVal(left, right Value) Value {
	val := StrictEqual(left, right)

	return ToElkBool(!val)
}

// Execute a right bit shift >>.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func RightBitshiftVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.RightBitshiftVal(right)
		case Int64:
			r, err := StrictIntRightBitshift(l, right)
			return r.ToValue(), err
		case UInt64:
			r, err := StrictIntRightBitshift(l, right)
			return r.ToValue(), err
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.RightBitshiftVal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return Undefined, Undefined
	}
}

// Execute a logical right bit shift >>>.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func LogicalRightBitshiftVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case Int64:
			r, err := StrictIntLogicalRightBitshift(l, right, LogicalRightShift64)
			return r.ToValue(), err
		case UInt64:
			r, err := StrictIntRightBitshift(l, right)
			return r.ToValue(), err
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return Undefined, Undefined
	}
}

// Execute a left bit shift <<.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func LeftBitshiftVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.LeftBitshiftVal(right)
		case Int64:
			r, err := StrictIntLeftBitshift(l, right)
			return r.ToValue(), err
		case UInt64:
			r, err := StrictIntLeftBitshift(l, right)
			return r.ToValue(), err
		case *ArrayList:
			l.Append(right)
			return Ref(l), Undefined
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.LeftBitshiftVal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return Undefined, Undefined
	}
}

// Execute a logical left bit shift <<<.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func LogicalLeftBitshiftVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case Int64:
			r, err := StrictIntLogicalLeftBitshift(l, right, LogicalRightShift64)
			return r.ToValue(), err
		case UInt64:
			r, err := StrictIntLeftBitshift(l, right)
			return r.ToValue(), err
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return Undefined, Undefined
	}
}

// Execute a bitwise AND &.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func BitwiseAndVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.BitwiseAndVal(right)
		case Int64:
			result, err := l.BitwiseAnd(right)
			return result.ToValue(), err
		case UInt64:
			result, err := l.BitwiseAnd(right)
			return result.ToValue(), err
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.BitwiseAndVal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return Undefined, Undefined
	}
}

// Execute a bitwise AND NOT &^.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func BitwiseAndNotVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.BitwiseAndNotVal(right)
		case Int64:
			result, err := l.BitwiseAndNot(right)
			return result.ToValue(), err
		case UInt64:
			result, err := l.BitwiseAndNot(right)
			return result.ToValue(), err
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.BitwiseAndNotVal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return Undefined, Undefined
	}
}

// Execute a bitwise OR |.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func BitwiseOrVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.BitwiseOrVal(right)
		case Int64:
			result, err := l.BitwiseOr(right)
			return result.ToValue(), err
		case UInt64:
			result, err := l.BitwiseOr(right)
			return result.ToValue(), err
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.BitwiseOrVal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return Undefined, Undefined
	}
}

// Execute a bitwise XOR ^.
// When successful returns (result, undefined).
// When an error occurred returns (undefined, error).
// When there are no builtin addition functions for the given type returns (undefined, undefined).
func BitwiseXorVal(left, right Value) (result, err Value) {
	if left.IsReference() {
		switch l := left.AsReference().(type) {
		case *BigInt:
			return l.BitwiseXorVal(right)
		case Int64:
			result, err := l.BitwiseXor(right)
			return result.ToValue(), err
		case UInt64:
			result, err := l.BitwiseXor(right)
			return result.ToValue(), err
		default:
			return Undefined, Undefined
		}
	}

	switch left.ValueFlag() {
	case SMALL_INT_FLAG:
		l := left.AsSmallInt()
		return l.BitwiseXorVal(right)
	case INT64_FLAG:
		l := left.AsInlineInt64()
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
		l := left.AsInlineUInt64()
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
		return Undefined, Undefined
	}
}

// Call `next`
func Next(val Value) (result, err Value) {
	if !val.IsReference() {
		return Undefined, Undefined
	}

	switch v := val.AsReference().(type) {
	case *ArrayListIterator:
		return v.Next()
	case *ArrayTupleIterator:
		return v.Next()
	case *HashMapIterator:
		return v.Next()
	case *HashRecordIterator:
		return v.Next()
	case *HashSetIterator:
		return v.Next()
	case *StringCharIterator:
		return v.Next()
	case *StringByteIterator:
		return v.Next()
	case *StringGraphemeIterator:
		return v.Next()
	case *Channel:
		return v.Next()
	default:
		return Undefined, Undefined
	}
}

// Get the iterator of the value
func Iter(val Value) Value {
	if !val.IsReference() {
		return Undefined
	}

	switch v := val.AsReference().(type) {
	case *ArrayListIterator, *ArrayTupleIterator, *HashMapIterator,
		*HashRecordIterator, *HashSetIterator, *StringCharIterator,
		*StringByteIterator, *StringGraphemeIterator, *Channel:
		return val
	case String:
		return Ref(NewStringCharIterator(v))
	case *ArrayList:
		return Ref(NewArrayListIterator(v))
	case *ArrayTuple:
		return Ref(NewArrayTupleIterator(v))
	case *HashMap:
		return Ref(NewHashMapIterator(v))
	case *HashRecord:
		return Ref(NewHashRecordIterator(v))
	case *HashSet:
		return Ref(NewHashSetIterator(v))
	default:
		return Undefined
	}
}

func NewReferenceComparer() cmp.Option {
	filter := func(x Value, y Value) bool { return x.IsReference() && y.IsReference() }
	transformer := cmp.Transformer("ValueToReference", func(val Value) Reference { return val.AsReference() })

	return cmp.FilterValues(
		filter,
		transformer,
	)
}

func NewInlineValueComparer(opts *cmp.Options) cmp.Option {
	comparer := cmp.Comparer(func(x, y Value) bool {
		if x.IsReference() || y.IsReference() {
			return false
		}
		if x.ValueFlag() != y.ValueFlag() {
			return false
		}

		switch x.ValueFlag() {
		case FLOAT32_FLAG:
			x := x.AsFloat32()
			y := y.AsFloat32()
			if x.IsNaN() || y.IsNaN() {
				return x.IsNaN() && y.IsNaN()
			}
			return x == y
		case FLOAT64_FLAG:
			x := x.AsInlineFloat64()
			y := y.AsInlineFloat64()
			if x.IsNaN() || y.IsNaN() {
				return x.IsNaN() && y.IsNaN()
			}
			return x == y
		case FLOAT_FLAG:
			x := x.AsFloat()
			y := y.AsFloat()
			if x.IsNaN() || y.IsNaN() {
				return x.IsNaN() && y.IsNaN()
			}
			return x == y
		default:
			return x.data == y.data
		}
	})

	filter := func(x Value, y Value) bool { return !x.IsReference() || !y.IsReference() }

	return cmp.FilterValues(
		filter,
		comparer,
	)
}
