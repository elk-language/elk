package value

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"unsafe"

	"github.com/cespare/xxhash/v2"
)

var UInt32Class *Class // ::Std::UInt32

// Elk's UInt32 value
type UInt32 uint32

func (i UInt32) ToValue() Value {
	return Value{
		flag: UINT32_FLAG,
		data: *(*uintptr)(unsafe.Pointer(&i)),
	}
}

func (UInt32) Class() *Class {
	return UInt32Class
}

func (UInt32) DirectClass() *Class {
	return UInt32Class
}

func (UInt32) SingletonClass() *Class {
	return nil
}

// Convert to Elk String.
func (i UInt32) ToString() String {
	return String(strconv.FormatUint(uint64(i), 10))
}

// Convert to Elk SmallInt.
func (i UInt32) ToSmallInt() SmallInt {
	return SmallInt(i)
}

// Convert to Elk Float.
func (i UInt32) ToFloat() Float {
	return Float(i)
}

// Convert to Elk Float64.
func (i UInt32) ToFloat64() Float64 {
	return Float64(i)
}

// Convert to Elk Float32.
func (i UInt32) ToFloat32() Float32 {
	return Float32(i)
}

// Convert to Elk Int64.
func (i UInt32) ToInt64() Int64 {
	return Int64(i)
}

// Convert to Elk Int32.
func (i UInt32) ToInt32() Int32 {
	return Int32(i)
}

// Convert to Elk Int16.
func (i UInt32) ToInt16() Int16 {
	return Int16(i)
}

// Convert to Elk Int8.
func (i UInt32) ToInt8() Int8 {
	return Int8(i)
}

// Convert to Elk UInt64.
func (i UInt32) ToUInt64() UInt64 {
	return UInt64(i)
}

// Convert to Elk UInt16.
func (i UInt32) ToUInt16() UInt16 {
	return UInt16(i)
}

// Convert to Elk UInt8.
func (i UInt32) ToUInt8() UInt8 {
	return UInt8(i)
}

func (i UInt32) Inspect() string {
	return fmt.Sprintf("%du32", i)
}

func (i UInt32) Error() string {
	return i.Inspect()
}

func (i UInt32) InstanceVariables() SymbolMap {
	return nil
}

func (i UInt32) Hash() UInt64 {
	d := xxhash.New()
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	d.Write(b)
	return UInt64(d.Sum64())
}

func (i UInt32) Add(other Value) (UInt32, Value) {
	if !other.IsUInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt32()
	return i + o, Undefined
}

// Perform a bitwise AND.
func (i UInt32) BitwiseAnd(other Value) (UInt32, Value) {
	if !other.IsUInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt32()
	return i & o, Undefined
}

// Perform a bitwise AND NOT.
func (i UInt32) BitwiseAndNot(other Value) (UInt32, Value) {
	if !other.IsUInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt32()
	return i &^ o, Undefined
}

// Perform a bitwise OR.
func (i UInt32) BitwiseOr(other Value) (UInt32, Value) {
	if !other.IsUInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt32()
	return i | o, Undefined
}

// Perform a bitwise XOR.
func (i UInt32) BitwiseXor(other Value) (UInt32, Value) {
	if !other.IsUInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt32()
	return i ^ o, Undefined
}

func (i UInt32) Exponentiate(other Value) (UInt32, Value) {
	if !other.IsUInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt32()
	if o <= 0 {
		return 1, Undefined
	}
	result := i
	var j UInt32
	for j = 2; j <= o; j++ {
		result *= i
	}
	return result, Undefined
}

func (i UInt32) Subtract(other Value) (UInt32, Value) {
	if !other.IsUInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt32()
	return i - o, Undefined
}

func (i UInt32) Multiply(other Value) (UInt32, Value) {
	if !other.IsUInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt32()
	return i * o, Undefined
}

func (i UInt32) Modulo(other Value) (UInt32, Value) {
	if !other.IsUInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt32()
	return i % o, Undefined
}

func (i UInt32) Divide(other Value) (UInt32, Value) {
	if !other.IsUInt32() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsUInt32()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i / o, Undefined
}

func (i UInt32) Compare(other Value) (Value, Value) {
	if !other.IsUInt32() {
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsUInt32()

	if i > o {
		return SmallInt(1).ToValue(), Undefined
	}
	if i < o {
		return SmallInt(-1).ToValue(), Undefined
	}
	return SmallInt(0).ToValue(), Undefined
}

func (i UInt32) GreaterThan(other Value) (Value, Value) {
	if !other.IsUInt32() {
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt32()
	return ToElkBool(i > o), Undefined
}

func (i UInt32) GreaterThanEqual(other Value) (Value, Value) {
	if !other.IsUInt32() {
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt32()
	return ToElkBool(i >= o), Undefined
}

func (i UInt32) LessThan(other Value) (Value, Value) {
	if !other.IsUInt32() {
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt32()
	return ToElkBool(i < o), Undefined
}

func (i UInt32) LessThanEqual(other Value) (Value, Value) {
	if !other.IsUInt32() {
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt32()
	return ToElkBool(i <= o), Undefined
}

func (i UInt32) Equal(other Value) Value {
	if !other.IsUInt32() {
		return False
	}

	o := other.AsUInt32()
	return ToElkBool(i == o)
}

func (i UInt32) StrictEqual(other Value) Value {
	return i.Equal(other)
}

func initUInt32() {
	UInt32Class = NewClass()
	StdModule.AddConstantString("UInt32", Ref(UInt32Class))
}
