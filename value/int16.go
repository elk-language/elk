package value

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/cespare/xxhash/v2"
)

var Int16Class *Class // ::Std::Int16

// Elk's Int16 value
type Int16 int16

func (i Int16) ToValue() Value {
	return Value{
		flag: INT16_FLAG,
		data: uintptr(i),
	}
}

func (i Int16) Class() *Class {
	return Int16Class
}

func (Int16) DirectClass() *Class {
	return Int16Class
}

func (Int16) SingletonClass() *Class {
	return nil
}

// Convert to Elk String.
func (i Int16) ToString() String {
	return String(strconv.Itoa(int(i)))
}

// Convert to Elk SmallInt.
func (i Int16) ToSmallInt() SmallInt {
	return SmallInt(i)
}

// Convert to Elk Float.
func (i Int16) ToFloat() Float {
	return Float(i)
}

// Convert to Elk Float64.
func (i Int16) ToFloat64() Float64 {
	return Float64(i)
}

// Convert to Elk Float32.
func (i Int16) ToFloat32() Float32 {
	return Float32(i)
}

// Convert to Elk Int64.
func (i Int16) ToInt64() Int64 {
	return Int64(i)
}

// Convert to Elk Int32.
func (i Int16) ToInt32() Int32 {
	return Int32(i)
}

// Convert to Elk Int8.
func (i Int8) ToInt8() Int8 {
	return Int8(i)
}

// Convert to Elk UInt64.
func (i Int16) ToUInt64() UInt64 {
	return UInt64(i)
}

// Convert to Elk UInt32.
func (i Int16) ToUInt32() UInt32 {
	return UInt32(i)
}

// Convert to Elk UInt16.
func (i Int16) ToUInt16() UInt16 {
	return UInt16(i)
}

// Convert to Elk UInt8.
func (i Int16) ToUInt8() UInt8 {
	return UInt8(i)
}

func (i Int16) Inspect() string {
	return fmt.Sprintf("%di16", i)
}

func (i Int16) Error() string {
	return i.Inspect()
}

func (i Int16) InstanceVariables() SymbolMap {
	return nil
}

func (i Int16) Hash() UInt64 {
	d := xxhash.New()
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(i))
	d.Write(b)
	return UInt64(d.Sum64())
}

func (i Int16) Add(other Value) (Int16, Value) {
	if !other.IsInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt16()
	return i + o, Undefined
}

// Perform a bitwise AND.
func (i Int16) BitwiseAnd(other Value) (Int16, Value) {
	if !other.IsInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt16()
	return i & o, Undefined
}

// Perform a bitwise AND NOT.
func (i Int16) BitwiseAndNot(other Value) (Int16, Value) {
	if !other.IsInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt16()
	return i &^ o, Undefined
}

// Perform a bitwise OR.
func (i Int16) BitwiseOr(other Value) (Int16, Value) {
	if !other.IsInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt16()
	return i | o, Undefined
}

// Perform a bitwise XOR.
func (i Int16) BitwiseXor(other Value) (Int16, Value) {
	if !other.IsInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt16()
	return i ^ o, Undefined
}

func (i Int16) ExponentiateVal(other Value) (Int16, Value) {
	if !other.IsInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt16()
	if o <= 0 {
		return 1, Undefined
	}
	result := i
	var j Int16
	for j = 2; j <= o; j++ {
		result *= i
	}
	return result, Undefined
}

func (i Int16) Subtract(other Value) (Int16, Value) {
	if !other.IsInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt16()
	return i - o, Undefined
}

func (i Int16) Multiply(other Value) (Int16, Value) {
	if !other.IsInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt16()
	return i * o, Undefined
}

func (i Int16) ModuloVal(other Value) (Int16, Value) {
	if !other.IsInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt16()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i % o, Undefined
}

func (i Int16) Divide(other Value) (Int16, Value) {
	if !other.IsInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsInt16()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i / o, Undefined
}

func (i Int16) CompareVal(other Value) (Value, Value) {
	if !other.IsInt16() {
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsInt16()

	if i > o {
		return SmallInt(1).ToValue(), Undefined
	}
	if i < o {
		return SmallInt(-1).ToValue(), Undefined
	}
	return SmallInt(0).ToValue(), Undefined
}

func (i Int16) GreaterThanVal(other Value) (Value, Value) {
	result, err := i.GreaterThan(other)
	return ToElkBool(result), err
}

func (i Int16) GreaterThan(other Value) (bool, Value) {
	if !other.IsInt16() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt16()
	return i > o, Undefined
}

func (i Int16) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := i.GreaterThanEqual(other)
	return ToElkBool(result), err
}

func (i Int16) GreaterThanEqual(other Value) (bool, Value) {
	if !other.IsInt16() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt16()
	return i >= o, Undefined
}

func (i Int16) LessThanVal(other Value) (Value, Value) {
	result, err := i.LessThan(other)
	return ToElkBool(result), err
}

func (i Int16) LessThan(other Value) (bool, Value) {
	if !other.IsInt16() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt16()
	return i < o, Undefined
}

func (i Int16) LessThanEqualVal(other Value) (Value, Value) {
	result, err := i.LessThanEqual(other)
	return ToElkBool(result), err
}

func (i Int16) LessThanEqual(other Value) (bool, Value) {
	if !other.IsInt16() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt16()
	return i <= o, Undefined
}

func (i Int16) EqualVal(other Value) Value {
	return ToElkBool(i.Equal(other))
}

func (i Int16) Equal(other Value) bool {
	if !other.IsInt16() {
		return false
	}

	o := other.AsInt16()
	return i == o
}

func (i Int16) StrictEqualVal(other Value) Value {
	return i.EqualVal(other)
}

func initInt16() {
	Int16Class = NewClassWithOptions(ClassWithParent(ValueClass))
	StdModule.AddConstantString("Int16", Ref(Int16Class))

	Int16Class.AddConstantString("Convertible", Ref(NewInterface()))
}
