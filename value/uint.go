package value

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/cespare/xxhash/v2"
)

var UIntClass *Class // ::Std::UInt

func (i UInt) ToValue() Value {
	return Value{
		flag: UINT_FLAG,
		data: uintptr(i),
	}
}

func (UInt) Class() *Class {
	return UIntClass
}

func (UInt) DirectClass() *Class {
	return UIntClass
}

func (UInt) SingletonClass() *Class {
	return nil
}

// Convert to Elk String.
func (i UInt) ToString() String {
	return String(strconv.FormatUint(uint64(i), 10))
}

// Convert to Elk SmallInt.
func (i UInt) ToSmallInt() SmallInt {
	return SmallInt(i)
}

// Convert to Elk Float.
func (i UInt) ToFloat() Float {
	return Float(i)
}

// Convert to Elk Float64.
func (i UInt) ToFloat64() Float64 {
	return Float64(i)
}

// Convert to Elk Float32.
func (i UInt) ToFloat32() Float32 {
	return Float32(i)
}

// Convert to Elk Int64.
func (i UInt) ToInt64() Int64 {
	return Int64(i)
}

// Convert to Elk Int32.
func (i UInt) ToInt32() Int32 {
	return Int32(i)
}

// Convert to Elk Int16.
func (i UInt) ToInt16() Int16 {
	return Int16(i)
}

// Convert to Elk Int8.
func (i UInt) ToInt8() Int8 {
	return Int8(i)
}

// Convert to Elk UInt64.
func (i UInt) ToUInt64() UInt64 {
	return UInt64(i)
}

// Convert to Elk UInt32.
func (i UInt) ToUInt32() UInt32 {
	return UInt32(i)
}

// Convert to Elk UInt16.
func (i UInt) ToUInt16() UInt16 {
	return UInt16(i)
}

// Convert to Elk UInt8.
func (i UInt) ToUInt8() UInt8 {
	return UInt8(i)
}

func (i UInt) Inspect() string {
	return fmt.Sprintf("%du", i)
}

func (i UInt) Error() string {
	return i.Inspect()
}

func (i UInt) InstanceVariables() *InstanceVariables {
	return nil
}

func (i UInt) Hash() UInt64 {
	d := xxhash.New()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	d.Write(b)
	return UInt64(d.Sum64())
}

func (i UInt) Add(other Value) (UInt, Value) {
	if !other.IsUInt() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt()
	return i + o, Undefined
}

// Perform a bitwise AND.
func (i UInt) BitwiseAnd(other Value) (UInt, Value) {
	if !other.IsUInt() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt()
	return i & o, Undefined
}

// Perform a bitwise AND NOT.
func (i UInt) BitwiseAndNot(other Value) (UInt, Value) {
	if !other.IsUInt() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt()
	return i &^ o, Undefined
}

// Perform a bitwise OR.
func (i UInt) BitwiseOr(other Value) (UInt, Value) {
	if !other.IsUInt() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt()
	return i | o, Undefined
}

// Perform a bitwise XOR.
func (i UInt) BitwiseXor(other Value) (UInt, Value) {
	if !other.IsUInt() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt()
	return i ^ o, Undefined
}

func (i UInt) ExponentiateVal(other Value) (UInt, Value) {
	if !other.IsUInt() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt()
	if o <= 0 {
		return 1, Undefined
	}
	result := i
	var j UInt
	for j = 2; j <= o; j++ {
		result *= i
	}
	return result, Undefined
}

func (i UInt) Subtract(other Value) (UInt, Value) {
	if !other.IsUInt() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt()
	return i - o, Undefined
}

func (i UInt) Multiply(other Value) (UInt, Value) {
	if !other.IsUInt() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt()
	return i * o, Undefined
}

func (i UInt) ModuloVal(other Value) (UInt, Value) {
	if !other.IsUInt() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i % o, Undefined
}

func (i UInt) Divide(other Value) (UInt, Value) {
	if !other.IsUInt() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsUInt()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i / o, Undefined
}

func (i UInt) CompareVal(other Value) (Value, Value) {
	if !other.IsUInt() {
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsUInt()

	if i > o {
		return SmallInt(1).ToValue(), Undefined
	}
	if i < o {
		return SmallInt(-1).ToValue(), Undefined
	}
	return SmallInt(0).ToValue(), Undefined
}

func (i UInt) GreaterThanVal(other Value) (Value, Value) {
	result, err := i.GreaterThan(other)
	return ToElkBool(result), err
}

func (i UInt) GreaterThan(other Value) (bool, Value) {
	if !other.IsUInt() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt()
	return i > o, Undefined
}

func (i UInt) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := i.GreaterThanEqual(other)
	return ToElkBool(result), err
}

func (i UInt) GreaterThanEqual(other Value) (bool, Value) {
	if !other.IsUInt() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt()
	return i >= o, Undefined
}

func (i UInt) LessThanVal(other Value) (Value, Value) {
	result, err := i.LessThan(other)
	return ToElkBool(result), err
}

func (i UInt) LessThan(other Value) (bool, Value) {
	if !other.IsUInt() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt()
	return i < o, Undefined
}

func (i UInt) LessThanEqualVal(other Value) (Value, Value) {
	result, err := i.LessThanEqual(other)
	return ToElkBool(result), err
}

func (i UInt) LessThanEqual(other Value) (bool, Value) {
	if !other.IsUInt() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt()
	return i <= o, Undefined
}

func (i UInt) EqualVal(other Value) Value {
	return ToElkBool(i.Equal(other))
}

func (i UInt) Equal(other Value) bool {
	if !other.IsUInt() {
		return false
	}

	o := other.AsUInt()
	return i == o
}

func (i UInt) StrictEqualVal(other Value) Value {
	return i.EqualVal(other)
}

func initUInt() {
	UIntClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	StdModule.AddConstantString("UInt", Ref(UIntClass))

	UIntClass.AddConstantString("Convertible", Ref(NewInterface()))
}
