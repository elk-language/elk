package value

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/cespare/xxhash/v2"
)

var Int64Class *Class // ::Std::Int64

// Elk's Int64 value
type Int64 int64

func (i Int64) Class() *Class {
	return Int64Class
}

func (Int64) DirectClass() *Class {
	return Int64Class
}

func (Int64) SingletonClass() *Class {
	return nil
}

// Convert to Elk String.
func (i Int64) ToString() String {
	return String(strconv.FormatInt(int64(i), 10))
}

// Convert to Elk SmallInt.
func (i Int64) ToSmallInt() SmallInt {
	return SmallInt(i)
}

// Convert to Elk Float.
func (i Int64) ToFloat() Float {
	return Float(i)
}

// Convert to Elk Float64.
func (i Int64) ToFloat64() Float64 {
	return Float64(i)
}

// Convert to Elk Float32.
func (i Int64) ToFloat32() Float32 {
	return Float32(i)
}

// Convert to Elk Int32.
func (i Int64) ToInt32() Int32 {
	return Int32(i)
}

// Convert to Elk Int16.
func (i Int64) ToInt16() Int16 {
	return Int16(i)
}

// Convert to Elk Int8.
func (i Int64) ToInt8() Int8 {
	return Int8(i)
}

// Convert to Elk UInt64.
func (i Int64) ToUInt64() UInt64 {
	return UInt64(i)
}

// Convert to Elk UInt32.
func (i Int64) ToUInt32() UInt32 {
	return UInt32(i)
}

// Convert to Elk UInt16.
func (i Int64) ToUInt16() UInt16 {
	return UInt16(i)
}

// Convert to Elk UInt8.
func (i Int64) ToUInt8() UInt8 {
	return UInt8(i)
}

func (i Int64) Copy() Reference {
	return i
}

func (i Int64) Inspect() string {
	return fmt.Sprintf("%di64", i)
}

func (i Int64) Error() string {
	return i.Inspect()
}

func (i Int64) InstanceVariables() SymbolMap {
	return nil
}

func (i Int64) Hash() UInt64 {
	d := xxhash.New()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	d.Write(b)
	return UInt64(d.Sum64())
}

func (i Int64) Add(other Value) (Int64, Value) {
	if !other.IsInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt64()
	return i + o, Undefined
}

// Perform a bitwise AND.
func (i Int64) BitwiseAnd(other Value) (Int64, Value) {
	if !other.IsInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt64()
	return i & o, Undefined
}

// Perform a bitwise AND NOT.
func (i Int64) BitwiseAndNot(other Value) (Int64, Value) {
	if !other.IsInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt64()
	return i &^ o, Undefined
}

// Perform a bitwise OR.
func (i Int64) BitwiseOr(other Value) (Int64, Value) {
	if !other.IsInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt64()
	return i | o, Undefined
}

// Perform a bitwise XOR.
func (i Int64) BitwiseXor(other Value) (Int64, Value) {
	if !other.IsInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt64()
	return i ^ o, Undefined
}

func (i Int64) ExponentiateVal(other Value) (Int64, Value) {
	if !other.IsInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt64()
	if o <= 0 {
		return 1, Undefined
	}
	result := i
	var j Int64
	for j = 2; j <= o; j++ {
		result *= i
	}
	return result, Undefined
}

func (i Int64) Subtract(other Value) (Int64, Value) {
	if !other.IsInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt64()
	return i - o, Undefined
}

func (i Int64) Multiply(other Value) (Int64, Value) {
	if !other.IsInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt64()
	return i * o, Undefined
}

func (i Int64) ModuloVal(other Value) (Int64, Value) {
	if !other.IsInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt64()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i % o, Undefined
}

func (i Int64) Divide(other Value) (Int64, Value) {
	if !other.IsInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsInt64()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i / o, Undefined
}

func (i Int64) CompareVal(other Value) (Value, Value) {
	if !other.IsInt64() {
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsInt64()

	if i > o {
		return SmallInt(1).ToValue(), Undefined
	}
	if i < o {
		return SmallInt(-1).ToValue(), Undefined
	}
	return SmallInt(0).ToValue(), Undefined
}

func (i Int64) GreaterThanVal(other Value) (Value, Value) {
	result, err := i.GreaterThan(other)
	return ToElkBool(result), err
}

func (i Int64) GreaterThan(other Value) (bool, Value) {
	if !other.IsInt64() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt64()
	return i > o, Undefined
}

func (i Int64) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := i.GreaterThanEqual(other)
	return ToElkBool(result), err
}

func (i Int64) GreaterThanEqual(other Value) (bool, Value) {
	if !other.IsInt64() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt64()
	return i >= o, Undefined
}

func (i Int64) LessThanVal(other Value) (Value, Value) {
	result, err := i.LessThan(other)
	return ToElkBool(result), err
}

func (i Int64) LessThan(other Value) (bool, Value) {
	if !other.IsInt64() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt64()
	return i < o, Undefined
}

func (i Int64) LessThanEqualVal(other Value) (Value, Value) {
	result, err := i.LessThanEqual(other)
	return ToElkBool(result), err
}

func (i Int64) LessThanEqual(other Value) (bool, Value) {
	if !other.IsInt64() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsInt64()
	return i <= o, Undefined
}

func (i Int64) EqualVal(other Value) Value {
	return ToElkBool(i.Equal(other))
}

func (i Int64) Equal(other Value) bool {
	if !other.IsInt64() {
		return false
	}

	o := other.AsInt64()
	return i == o
}

func (i Int64) StrictEqualVal(other Value) Value {
	return i.EqualVal(other)
}

func initInt64() {
	Int64Class = NewClassWithOptions(ClassWithParent(ValueClass))
	StdModule.AddConstantString("Int64", Ref(Int64Class))
}
