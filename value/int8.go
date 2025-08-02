package value

import (
	"fmt"
	"strconv"

	"github.com/cespare/xxhash/v2"
)

var Int8Class *Class // ::Std::Int8

// Elk's Int8 value
type Int8 int8

func (i Int8) ToValue() Value {
	return Value{
		flag: INT8_FLAG,
		data: uintptr(i),
	}
}

func (i Int8) Class() *Class {
	return Int8Class
}

func (Int8) DirectClass() *Class {
	return Int8Class
}

func (Int8) SingletonClass() *Class {
	return nil
}

// Convert to Elk String.
func (i Int8) ToString() String {
	return String(strconv.Itoa(int(i)))
}

// Convert to Elk SmallInt.
func (i Int8) ToSmallInt() SmallInt {
	return SmallInt(i)
}

// Convert to Elk Float.
func (i Int8) ToFloat() Float {
	return Float(i)
}

// Convert to Elk Float64.
func (i Int8) ToFloat64() Float64 {
	return Float64(i)
}

// Convert to Elk Float32.
func (i Int8) ToFloat32() Float32 {
	return Float32(i)
}

// Convert to Elk Int64.
func (i Int8) ToInt64() Int64 {
	return Int64(i)
}

// Convert to Elk Int32.
func (i Int8) ToInt32() Int32 {
	return Int32(i)
}

// Convert to Elk Int16.
func (i Int8) ToInt16() Int16 {
	return Int16(i)
}

// Convert to Elk UInt64.
func (i Int8) ToUInt64() UInt64 {
	return UInt64(i)
}

// Convert to Elk UInt32.
func (i Int8) ToUInt32() UInt32 {
	return UInt32(i)
}

// Convert to Elk UInt16.
func (i Int8) ToUInt16() UInt16 {
	return UInt16(i)
}

// Convert to Elk UInt8.
func (i Int8) ToUInt8() UInt8 {
	return UInt8(i)
}

func (i Int8) Inspect() string {
	return fmt.Sprintf("%di8", i)
}

func (i Int8) Error() string {
	return i.Inspect()
}

func (i Int8) InstanceVariables() *InstanceVariables {
	return nil
}

func (i Int8) Hash() UInt64 {
	d := xxhash.New()
	d.Write([]byte{byte(i)})
	return UInt64(d.Sum64())
}

func (i Int8) Add(other Value) (Int8, Value) {
	if !other.IsInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt8()
	return i + o, Undefined
}

// Perform a bitwise AND.
func (i Int8) BitwiseAnd(other Value) (Int8, Value) {
	if !other.IsInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt8()
	return i & o, Undefined
}

// Perform a bitwise AND NOT.
func (i Int8) BitwiseAndNot(other Value) (Int8, Value) {
	if !other.IsInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt8()
	return i &^ o, Undefined
}

// Perform a bitwise OR.
func (i Int8) BitwiseOr(other Value) (Int8, Value) {
	if !other.IsInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt8()
	return i | o, Undefined
}

// Perform a bitwise XOR.
func (i Int8) BitwiseXor(other Value) (Int8, Value) {
	if !other.IsInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt8()
	return i ^ o, Undefined
}

func (i Int8) ExponentiateVal(other Value) (Int8, Value) {
	if !other.IsInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt8()
	if o <= 0 {
		return 1, Undefined
	}
	result := i
	var j Int8
	for j = 2; j <= o; j++ {
		result *= i
	}
	return result, Undefined
}

func (i Int8) Subtract(other Value) (Int8, Value) {
	if !other.IsInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt8()
	return i - o, Undefined
}

func (i Int8) Multiply(other Value) (Int8, Value) {
	if !other.IsInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt8()
	return i * o, Undefined
}

func (i Int8) ModuloVal(other Value) (Int8, Value) {
	if !other.IsInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt8()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i % o, Undefined
}

func (i Int8) Divide(other Value) (Int8, Value) {
	if !other.IsInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsInt8()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i / o, Undefined
}

func (i Int8) CompareVal(other Value) (Value, Value) {
	if !other.IsInt8() {
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsInt8()

	if i > o {
		return SmallInt(1).ToValue(), Undefined
	}
	if i < o {
		return SmallInt(-1).ToValue(), Undefined
	}
	return SmallInt(0).ToValue(), Undefined
}

func (i Int8) GreaterThanVal(other Value) (Value, Value) {
	result, err := i.GreaterThan(other)
	return ToElkBool(result), err
}

func (i Int8) GreaterThan(other Value) (bool, Value) {
	if !other.IsInt8() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt8()
	return i > o, Undefined
}

func (i Int8) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := i.GreaterThanEqual(other)
	return ToElkBool(result), err
}

func (i Int8) GreaterThanEqual(other Value) (bool, Value) {
	if !other.IsInt8() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt8()
	return i >= o, Undefined
}

func (i Int8) LessThanVal(other Value) (Value, Value) {
	result, err := i.LessThan(other)
	return ToElkBool(result), err
}

func (i Int8) LessThan(other Value) (bool, Value) {
	if !other.IsInt8() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt8()
	return i < o, Undefined
}

func (i Int8) LessThanEqualVal(other Value) (Value, Value) {
	result, err := i.LessThanEqual(other)
	return ToElkBool(result), err
}

func (i Int8) LessThanEqual(other Value) (bool, Value) {
	if !other.IsInt8() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt8()
	return i <= o, Undefined
}

func (i Int8) EqualVal(other Value) Value {
	return ToElkBool(i.Equal(other))
}

func (i Int8) Equal(other Value) bool {
	if !other.IsInt8() {
		return false
	}

	o := other.AsInt8()
	return i == o
}

func (i Int8) StrictEqualVal(other Value) Value {
	return i.EqualVal(other)
}

func initInt8() {
	Int8Class = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	StdModule.AddConstantString("Int8", Ref(Int8Class))

	Int8Class.AddConstantString("Convertible", Ref(NewInterface()))
}
