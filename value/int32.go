package value

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/cespare/xxhash/v2"
)

var Int32Class *Class // ::Std::Int32

// Elk's Int32 value
type Int32 int32

func (i Int32) ToValue() Value {
	return Value{
		flag: INT32_FLAG,
		data: uintptr(i),
	}
}

func (i Int32) Class() *Class {
	return Int32Class
}

func (Int32) DirectClass() *Class {
	return Int32Class
}

func (Int32) SingletonClass() *Class {
	return nil
}

// Convert to Elk String.
func (i Int32) ToString() String {
	return String(strconv.Itoa(int(i)))
}

// Convert to Elk SmallInt.
func (i Int32) ToSmallInt() SmallInt {
	return SmallInt(i)
}

// Convert to Elk Float.
func (i Int32) ToFloat() Float {
	return Float(i)
}

// Convert to Elk Float64.
func (i Int32) ToFloat64() Float64 {
	return Float64(i)
}

// Convert to Elk Float32.
func (i Int32) ToFloat32() Float32 {
	return Float32(i)
}

// Convert to Elk Int64.
func (i Int32) ToInt64() Int64 {
	return Int64(i)
}

// Convert to Elk Int16.
func (i Int32) ToInt16() Int16 {
	return Int16(i)
}

// Convert to Elk Int8.
func (i Int32) ToInt8() Int8 {
	return Int8(i)
}

// Convert to Elk UInt32.
func (i Int32) ToUInt32() UInt32 {
	return UInt32(i)
}

// Convert to Elk UInt.
func (i Int32) ToUInt() UInt {
	return UInt(i)
}

// Convert to Elk UInt64.
func (i Int32) ToUInt64() UInt64 {
	return UInt64(i)
}

// Convert to Elk UInt16.
func (i Int32) ToUInt16() UInt16 {
	return UInt16(i)
}

// Convert to Elk UInt8.
func (i Int32) ToUInt8() UInt8 {
	return UInt8(i)
}

func (i Int32) Inspect() string {
	return fmt.Sprintf("%di32", i)
}

func (i Int32) Error() string {
	return i.Inspect()
}

func (i Int32) InstanceVariables() *InstanceVariables {
	return nil
}

func (i Int32) Hash() UInt64 {
	d := xxhash.New()
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	d.Write(b)
	return UInt64(d.Sum64())
}

func (i Int32) Add(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i + o, Undefined
}

// Perform a bitwise AND.
func (i Int32) BitwiseAnd(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i & o, Undefined
}

// Perform a bitwise AND NOT.
func (i Int32) BitwiseAndNot(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i &^ o, Undefined
}

// Perform a bitwise OR.
func (i Int32) BitwiseOr(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i | o, Undefined
}

// Perform a bitwise XOR.
func (i Int32) BitwiseXor(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i ^ o, Undefined
}

func (i Int32) ExponentiateVal(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	if o <= 0 {
		return 1, Undefined
	}
	result := i
	var j Int32
	for j = 2; j <= o; j++ {
		result *= i
	}
	return result, Undefined
}

func (i Int32) Subtract(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i - o, Undefined
}

func (i Int32) Multiply(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i * o, Undefined
}

func (i Int32) ModuloVal(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i % o, Undefined
}

func (i Int32) Divide(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsInt32()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i / o, Undefined
}

func (i Int32) CompareVal(other Value) (Value, Value) {
	if !other.IsInt32() {
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsInt32()

	if i > o {
		return SmallInt(1).ToValue(), Undefined
	}
	if i < o {
		return SmallInt(-1).ToValue(), Undefined
	}
	return SmallInt(0).ToValue(), Undefined
}

func (i Int32) GreaterThanVal(other Value) (Value, Value) {
	result, err := i.GreaterThan(other)
	return ToElkBool(result), err
}

func (i Int32) GreaterThan(other Value) (bool, Value) {
	if !other.IsInt32() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i > o, Undefined
}

func (i Int32) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := i.GreaterThanEqual(other)
	return ToElkBool(result), err
}

func (i Int32) GreaterThanEqual(other Value) (bool, Value) {
	if !other.IsInt32() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i >= o, Undefined
}

func (i Int32) LessThanVal(other Value) (Value, Value) {
	result, err := i.LessThan(other)
	return ToElkBool(result), err
}

func (i Int32) LessThan(other Value) (bool, Value) {
	if !other.IsInt32() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i < o, Undefined
}

func (i Int32) LessThanEqualVal(other Value) (Value, Value) {
	result, err := i.LessThanEqual(other)
	return ToElkBool(result), err
}

func (i Int32) LessThanEqual(other Value) (bool, Value) {
	if !other.IsInt32() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i <= o, Undefined
}
func (i Int32) EqualVal(other Value) Value {
	return ToElkBool(i.Equal(other))
}

func (i Int32) Equal(other Value) bool {
	if !other.IsInt32() {
		return false
	}

	o := other.AsInt32()
	return i == o
}

func (i Int32) StrictEqualVal(other Value) Value {
	return i.EqualVal(other)
}

func initInt32() {
	Int32Class = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	StdModule.AddConstantString("Int32", Ref(Int32Class))

	Int32Class.AddConstantString("Convertible", Ref(NewInterface()))
}
