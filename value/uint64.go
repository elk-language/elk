package value

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/cespare/xxhash/v2"
)

var UInt64Class *Class // ::Std::UInt64

// Elk's UInt64 value
type UInt64 uint64

func (UInt64) Class() *Class {
	return UInt64Class
}

func (i UInt64) Copy() Reference {
	return i
}

func (UInt64) DirectClass() *Class {
	return UInt64Class
}

func (UInt64) SingletonClass() *Class {
	return nil
}

// Convert to Elk String.
func (i UInt64) ToString() String {
	return String(strconv.FormatUint(uint64(i), 10))
}

// Convert to Elk SmallInt.
func (i UInt64) ToSmallInt() SmallInt {
	return SmallInt(i)
}

// Convert to Elk Float.
func (i UInt64) ToFloat() Float {
	return Float(i)
}

// Convert to Elk Float64.
func (i UInt64) ToFloat64() Float64 {
	return Float64(i)
}

// Convert to Elk Float32.
func (i UInt64) ToFloat32() Float32 {
	return Float32(i)
}

// Convert to Elk Int64.
func (i UInt64) ToInt64() Int64 {
	return Int64(i)
}

// Convert to Elk Int32.
func (i UInt64) ToInt32() Int32 {
	return Int32(i)
}

// Convert to Elk Int16.
func (i UInt64) ToInt16() Int16 {
	return Int16(i)
}

// Convert to Elk Int8.
func (i UInt64) ToInt8() Int8 {
	return Int8(i)
}

// Convert to Elk UInt32.
func (i UInt64) ToUInt32() UInt32 {
	return UInt32(i)
}

// Convert to Elk UInt16.
func (i UInt64) ToUInt16() UInt16 {
	return UInt16(i)
}

// Convert to Elk UInt8.
func (i UInt64) ToUInt8() UInt8 {
	return UInt8(i)
}

func (i UInt64) Inspect() string {
	return fmt.Sprintf("%du64", i)
}

func (i UInt64) Error() string {
	return i.Inspect()
}

func (i UInt64) InstanceVariables() SymbolMap {
	return nil
}

func (i UInt64) Hash() UInt64 {
	d := xxhash.New()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	d.Write(b)
	return UInt64(d.Sum64())
}

func (i UInt64) Add(other Value) (UInt64, Value) {
	if !other.IsUInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt64()
	return i + o, Undefined
}

// Perform a bitwise AND.
func (i UInt64) BitwiseAnd(other Value) (UInt64, Value) {
	if !other.IsUInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt64()
	return i & o, Undefined
}

// Perform a bitwise AND NOT.
func (i UInt64) BitwiseAndNot(other Value) (UInt64, Value) {
	if !other.IsUInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt64()
	return i &^ o, Undefined
}

// Perform a bitwise OR.
func (i UInt64) BitwiseOr(other Value) (UInt64, Value) {
	if !other.IsUInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt64()
	return i | o, Undefined
}

// Perform a bitwise XOR.
func (i UInt64) BitwiseXor(other Value) (UInt64, Value) {
	if !other.IsUInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt64()
	return i ^ o, Undefined
}

func (i UInt64) Exponentiate(other Value) (UInt64, Value) {
	if !other.IsUInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt64()
	if o <= 0 {
		return 1, Undefined
	}
	result := i
	var j UInt64
	for j = 2; j <= o; j++ {
		result *= i
	}
	return result, Undefined
}

func (i UInt64) Subtract(other Value) (UInt64, Value) {
	if !other.IsUInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt64()
	return i - o, Undefined
}

func (i UInt64) Multiply(other Value) (UInt64, Value) {
	if !other.IsUInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt64()
	return i * o, Undefined
}

func (i UInt64) Modulo(other Value) (UInt64, Value) {
	if !other.IsUInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt64()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i % o, Undefined
}

func (i UInt64) Divide(other Value) (UInt64, Value) {
	if !other.IsUInt64() {
		return 0, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsUInt64()
	if o == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return i / o, Undefined
}

func (i UInt64) Compare(other Value) (Value, Value) {
	if !other.IsUInt64() {
		return Undefined, Ref(NewCoerceError(i.Class(), other.Class()))
	}
	o := other.AsUInt64()

	if i > o {
		return SmallInt(1).ToValue(), Undefined
	}
	if i < o {
		return SmallInt(-1).ToValue(), Undefined
	}
	return SmallInt(0).ToValue(), Undefined
}

func (i UInt64) GreaterThan(other Value) (Value, Value) {
	result, err := i.GreaterThanBool(other)
	return ToElkBool(result), err
}

func (i UInt64) GreaterThanBool(other Value) (bool, Value) {
	if !other.IsUInt64() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt64()
	return i > o, Undefined
}

func (i UInt64) GreaterThanEqual(other Value) (Value, Value) {
	result, err := i.GreaterThanEqualBool(other)
	return ToElkBool(result), err
}

func (i UInt64) GreaterThanEqualBool(other Value) (bool, Value) {
	if !other.IsUInt64() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt64()
	return i >= o, Undefined
}

func (i UInt64) LessThan(other Value) (Value, Value) {
	result, err := i.LessThanBool(other)
	return ToElkBool(result), err
}

func (i UInt64) LessThanBool(other Value) (bool, Value) {
	if !other.IsUInt64() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt64()
	return i < o, Undefined
}

func (i UInt64) LessThanEqual(other Value) (Value, Value) {
	result, err := i.LessThanEqualBool(other)
	return ToElkBool(result), err
}

func (i UInt64) LessThanEqualBool(other Value) (bool, Value) {
	if !other.IsUInt64() {
		return false, Ref(NewCoerceError(i.Class(), other.Class()))
	}

	o := other.AsUInt64()
	return i <= o, Undefined
}

func (i UInt64) Equal(other Value) Value {
	return ToElkBool(i.EqualBool(other))
}

func (i UInt64) EqualBool(other Value) bool {
	if !other.IsUInt64() {
		return false
	}

	o := other.AsUInt64()
	return i == o
}

func (i UInt64) StrictEqual(other Value) Value {
	return i.Equal(other)
}

func initUInt64() {
	UInt64Class = NewClassWithOptions(ClassWithParent(ValueClass))
	StdModule.AddConstantString("UInt64", Ref(UInt64Class))
}
