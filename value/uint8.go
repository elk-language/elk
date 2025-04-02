package value

import (
	"fmt"
	"strconv"

	"github.com/cespare/xxhash/v2"
)

var UInt8Class *Class // ::Std::UInt8

// Elk's UInt8 value
type UInt8 uint8

func (i UInt8) ToValue() Value {
	return Value{
		flag: UINT8_FLAG,
		data: uintptr(i),
	}
}

func (UInt8) Class() *Class {
	return UInt8Class
}

func (UInt8) DirectClass() *Class {
	return UInt64Class
}

func (UInt8) SingletonClass() *Class {
	return nil
}

// Convert to Elk String.
func (i UInt8) ToString() String {
	return String(strconv.FormatUint(uint64(i), 10))
}

// Convert to Elk SmallInt.
func (i UInt8) ToSmallInt() SmallInt {
	return SmallInt(i)
}

// Convert to Elk Float.
func (i UInt8) ToFloat() Float {
	return Float(i)
}

// Convert to Elk Float64.
func (i UInt8) ToFloat64() Float64 {
	return Float64(i)
}

// Convert to Elk Float32.
func (i UInt8) ToFloat32() Float32 {
	return Float32(i)
}

// Convert to Elk Int64.
func (i UInt8) ToInt64() Int64 {
	return Int64(i)
}

// Convert to Elk Int32.
func (i UInt8) ToInt32() Int32 {
	return Int32(i)
}

// Convert to Elk Int16.
func (i UInt8) ToInt16() Int16 {
	return Int16(i)
}

// Convert to Elk Int8.
func (i UInt8) ToInt8() Int8 {
	return Int8(i)
}

// Convert to Elk UInt64.
func (i UInt8) ToUInt64() UInt64 {
	return UInt64(i)
}

// Convert to Elk UInt32.
func (i UInt8) ToUInt32() UInt32 {
	return UInt32(i)
}

// Convert to Elk UInt16.
func (i UInt8) ToUInt16() UInt16 {
	return UInt16(i)
}

func (i UInt8) Inspect() string {
	return fmt.Sprintf("%du8", i)
}

func (i UInt8) Error() string {
	return i.Inspect()
}

func (i UInt8) InstanceVariables() SymbolMap {
	return nil
}

func (i UInt8) Hash() UInt64 {
	d := xxhash.New()
	d.Write([]byte{byte(i)})
	return UInt64(d.Sum64())
}

func (i UInt8) Add(other Value) (UInt8, Value) {
	if !other.IsUInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt8()
	return i + o, Undefined
}

// Perform a bitwise AND.
func (i UInt8) BitwiseAnd(other Value) (UInt8, Value) {
	if !other.IsUInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt8()
	return i & o, Undefined
}

// Perform a bitwise AND NOT.
func (i UInt8) BitwiseAndNot(other Value) (UInt8, Value) {
	if !other.IsUInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt8()
	return i &^ o, Undefined
}

// Perform a bitwise OR.
func (i UInt8) BitwiseOr(other Value) (UInt8, Value) {
	if !other.IsUInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt8()
	return i | o, Undefined
}

// Perform a bitwise XOR.
func (i UInt8) BitwiseXor(other Value) (UInt8, Value) {
	if !other.IsUInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt8()
	return i ^ o, Undefined
}

func (i UInt8) ExponentiateVal(other Value) (UInt8, Value) {
	if !other.IsUInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt8()
	if o <= 0 {
		return 1, Undefined
	}
	result := i
	var j UInt8
	for j = 2; j <= o; j++ {
		result *= i
	}
	return result, Undefined
}

func (i UInt8) Subtract(other Value) (UInt8, Value) {
	if !other.IsUInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt8()
	return i - o, Undefined
}

func (i UInt8) Multiply(other Value) (UInt8, Value) {
	if !other.IsUInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt8()
	return i * o, Undefined
}

func (i UInt8) ModuloVal(other Value) (UInt8, Value) {
	if !other.IsUInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt8()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i % o, Undefined
}

func (i UInt8) Divide(other Value) (UInt8, Value) {
	if !other.IsUInt8() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsUInt8()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i / o, Undefined
}

func (i UInt8) CompareVal(other Value) (Value, Value) {
	if !other.IsUInt8() {
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsUInt8()

	if i > o {
		return SmallInt(1).ToValue(), Undefined
	}
	if i < o {
		return SmallInt(-1).ToValue(), Undefined
	}
	return SmallInt(0).ToValue(), Undefined
}

func (i UInt8) GreaterThanVal(other Value) (Value, Value) {
	result, err := i.GreaterThan(other)
	return ToElkBool(result), err
}

func (i UInt8) GreaterThan(other Value) (bool, Value) {
	if !other.IsUInt8() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt8()
	return i > o, Undefined
}

func (i UInt8) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := i.GreaterThanEqual(other)
	return ToElkBool(result), err
}

func (i UInt8) GreaterThanEqual(other Value) (bool, Value) {
	if !other.IsUInt8() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt8()
	return i >= o, Undefined
}

func (i UInt8) LessThanVal(other Value) (Value, Value) {
	result, err := i.LessThan(other)
	return ToElkBool(result), err
}

func (i UInt8) LessThan(other Value) (bool, Value) {
	if !other.IsUInt8() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt8()
	return i < o, Undefined
}

func (i UInt8) LessThanEqualVal(other Value) (Value, Value) {
	result, err := i.LessThanEqual(other)
	return ToElkBool(result), err
}

func (i UInt8) LessThanEqual(other Value) (bool, Value) {
	if !other.IsUInt8() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt8()
	return i <= o, Undefined
}

func (i UInt8) EqualVal(other Value) Value {
	return ToElkBool(i.Equal(other))
}

func (i UInt8) Equal(other Value) bool {
	if !other.IsUInt8() {
		return false
	}

	o := other.AsUInt8()
	return i == o
}

func (i UInt8) StrictEqualVal(other Value) Value {
	return i.EqualVal(other)
}

func initUInt8() {
	UInt8Class = NewClassWithOptions(ClassWithParent(ValueClass))
	StdModule.AddConstantString("UInt8", Ref(UInt8Class))
}
