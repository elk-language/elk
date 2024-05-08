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

func (UInt16) Class() *Class {
	return UInt16Class
}

func (i UInt16) Copy() Value {
	return i
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

func initUInt16() {
	UInt16Class = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("UInt16", UInt16Class)
}
