package value

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/cespare/xxhash/v2"
)

var UInt16Class *Class // ::Std::UInt16

// Elk's UInt16 value
type UInt16 uint16

func (i UInt16) ToValue() Value {
	return Value{
		flag: UINT16_FLAG,
		data: uintptr(i),
	}
}

func (UInt16) Class() *Class {
	return UInt16Class
}

func (UInt16) DirectClass() *Class {
	return UInt16Class
}

func (UInt16) SingletonClass() *Class {
	return nil
}

// Convert to Elk String.
func (i UInt16) ToString() String {
	return String(strconv.FormatUint(uint64(i), 10))
}

// Convert to Elk SmallInt.
func (i UInt16) ToSmallInt() SmallInt {
	return SmallInt(i)
}

// Convert to Elk Float.
func (i UInt16) ToFloat() Float {
	return Float(i)
}

// Convert to Elk Float64.
func (i UInt16) ToFloat64() Float64 {
	return Float64(i)
}

// Convert to Elk Float32.
func (i UInt16) ToFloat32() Float32 {
	return Float32(i)
}

// Convert to Elk Int64.
func (i UInt16) ToInt64() Int64 {
	return Int64(i)
}

// Convert to Elk Int32.
func (i UInt16) ToInt32() Int32 {
	return Int32(i)
}

// Convert to Elk Int16.
func (i UInt16) ToInt16() Int16 {
	return Int16(i)
}

// Convert to Elk Int8.
func (i UInt16) ToInt8() Int8 {
	return Int8(i)
}

// Convert to Elk UInt64.
func (i UInt16) ToUInt64() UInt64 {
	return UInt64(i)
}

// Convert to Elk UInt32.
func (i UInt16) ToUInt32() UInt32 {
	return UInt32(i)
}

// Convert to Elk UInt8.
func (i UInt16) ToUInt8() UInt8 {
	return UInt8(i)
}

func (i UInt16) Inspect() string {
	return fmt.Sprintf("%du16", i)
}

func (i UInt16) Error() string {
	return i.Inspect()
}

func (i UInt16) InstanceVariables() SymbolMap {
	return nil
}

func (i UInt16) Hash() UInt64 {
	d := xxhash.New()
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(i))
	d.Write(b)
	return UInt64(d.Sum64())
}

func (i UInt16) Add(other Value) (UInt16, Value) {
	if !other.IsUInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt16()
	return i + o, Undefined
}

// Perform a bitwise AND.
func (i UInt16) BitwiseAnd(other Value) (UInt16, Value) {
	if !other.IsUInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt16()
	return i & o, Undefined
}

// Perform a bitwise AND NOT.
func (i UInt16) BitwiseAndNot(other Value) (UInt16, Value) {
	if !other.IsUInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt16()
	return i &^ o, Undefined
}

// Perform a bitwise OR.
func (i UInt16) BitwiseOr(other Value) (UInt16, Value) {
	if !other.IsUInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt16()
	return i | o, Undefined
}

// Perform a bitwise XOR.
func (i UInt16) BitwiseXor(other Value) (UInt16, Value) {
	if !other.IsUInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt16()
	return i ^ o, Undefined
}

func (i UInt16) ExponentiateVal(other Value) (UInt16, Value) {
	if !other.IsUInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt16()
	if o <= 0 {
		return 1, Undefined
	}
	result := i
	var j UInt16
	for j = 2; j <= o; j++ {
		result *= i
	}
	return result, Undefined
}

func (i UInt16) Subtract(other Value) (UInt16, Value) {
	if !other.IsUInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt16()
	return i - o, Undefined
}

func (i UInt16) Multiply(other Value) (UInt16, Value) {
	if !other.IsUInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt16()
	return i * o, Undefined
}

func (i UInt16) ModuloVal(other Value) (UInt16, Value) {
	if !other.IsUInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt16()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i % o, Undefined
}

func (i UInt16) Divide(other Value) (UInt16, Value) {
	if !other.IsUInt16() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsUInt16()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i / o, Undefined
}

func (i UInt16) CompareVal(other Value) (Value, Value) {
	if !other.IsUInt16() {
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsUInt16()

	if i > o {
		return SmallInt(1).ToValue(), Undefined
	}
	if i < o {
		return SmallInt(-1).ToValue(), Undefined
	}
	return SmallInt(0).ToValue(), Undefined
}

func (i UInt16) GreaterThanVal(other Value) (Value, Value) {
	result, err := i.GreaterThan(other)
	return ToElkBool(result), err
}

func (i UInt16) GreaterThan(other Value) (bool, Value) {
	if !other.IsUInt16() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt16()
	return i > o, Undefined
}

func (i UInt16) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := i.GreaterThanEqual(other)
	return ToElkBool(result), err
}

func (i UInt16) GreaterThanEqual(other Value) (bool, Value) {
	if !other.IsUInt16() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt16()
	return i >= o, Undefined
}

func (i UInt16) LessThanVal(other Value) (Value, Value) {
	result, err := i.LessThan(other)
	return ToElkBool(result), err
}

func (i UInt16) LessThan(other Value) (bool, Value) {
	if !other.IsUInt16() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt16()
	return i < o, Undefined
}

func (i UInt16) LessThanEqualVal(other Value) (Value, Value) {
	result, err := i.LessThanEqual(other)
	return ToElkBool(result), err
}

func (i UInt16) LessThanEqual(other Value) (bool, Value) {
	if !other.IsUInt16() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt16()
	return i <= o, Undefined
}

func (i UInt16) EqualVal(other Value) Value {
	return ToElkBool(i.Equal(other))
}

func (i UInt16) Equal(other Value) bool {
	if !other.IsUInt16() {
		return false
	}

	o := other.AsUInt16()
	return i == o
}

func (i UInt16) StrictEqualVal(other Value) Value {
	return i.EqualVal(other)
}

func initUInt16() {
	UInt16Class = NewClassWithOptions(ClassWithParent(ValueClass))
	StdModule.AddConstantString("UInt16", Ref(UInt16Class))
}
