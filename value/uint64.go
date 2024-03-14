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

func (i UInt64) Copy() Value {
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

func initUInt64() {
	UInt64Class = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("UInt64", UInt64Class)
}
