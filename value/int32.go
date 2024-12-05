package value

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"unsafe"

	"github.com/cespare/xxhash/v2"
)

var Int32Class *Class // ::Std::Int32

// Elk's Int32 value
type Int32 int32

func (i Int32) ToValue() Value {
	return Value{
		data: unsafe.Pointer(uintptr(INT32_FLAG)),
		tab:  *(*uintptr)(unsafe.Pointer(&i)),
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

// Convert to Elk Int32.
func (i Int32) ToInt32() Int32 {
	return Int32(i)
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

func (i Int32) InstanceVariables() SymbolMap {
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
	return i + o, Nil
}

// Perform a bitwise AND.
func (i Int32) BitwiseAnd(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i & o, Nil
}

// Perform a bitwise AND NOT.
func (i Int32) BitwiseAndNot(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i &^ o, Nil
}

// Perform a bitwise OR.
func (i Int32) BitwiseOr(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i | o, Nil
}

// Perform a bitwise XOR.
func (i Int32) BitwiseXor(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i ^ o, Nil
}

func (i Int32) Exponentiate(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	if o <= 0 {
		return 1, Nil
	}
	result := i
	var j Int32
	for j = 2; j <= o; j++ {
		result *= i
	}
	return result, Nil
}

func (i Int32) Subtract(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i - o, Nil
}

func (i Int32) Multiply(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i * o, Nil
}

func (i Int32) Modulo(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return i % o, Nil
}

func (i Int32) Divide(other Value) (Int32, Value) {
	if !other.IsInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsInt32()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i / o, Nil
}

func (i Int32) Compare(other Value) (Value, Value) {
	if !other.IsInt32() {
		return Nil, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsInt32()

	if i > o {
		return SmallInt(1).ToValue(), Nil
	}
	if i < o {
		return SmallInt(-1).ToValue(), Nil
	}
	return SmallInt(0).ToValue(), Nil
}

func (i Int32) GreaterThan(other Value) (Value, Value) {
	if !other.IsInt32() {
		return Nil, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return ToElkBool(i > o), Nil
}

func (i Int32) GreaterThanEqual(other Value) (Value, Value) {
	if !other.IsInt32() {
		return Nil, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return ToElkBool(i >= o), Nil
}

func (i Int32) LessThan(other Value) (Value, Value) {
	if !other.IsInt32() {
		return Nil, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return ToElkBool(i < o), Nil
}

func (i Int32) LessThanEqual(other Value) (Value, Value) {
	if !other.IsInt32() {
		return Nil, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt32()
	return ToElkBool(i <= o), Nil
}

func (i Int32) Equal(other Value) Value {
	if !other.IsInt32() {
		return False
	}

	o := other.AsInt32()
	return ToElkBool(i == o)
}

func (i Int32) StrictEqual(other Value) Value {
	return i.Equal(other)
}

func initInt32() {
	Int32Class = NewClass()
	StdModule.AddConstantString("Int32", Ref(Int32Class))
}
