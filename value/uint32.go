package value

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/cespare/xxhash/v2"
)

var UInt32Class *Class // ::Std::UInt32

// Elk's UInt32 value
type UInt32 uint32

func (UInt32) Class() *Class {
	return UInt32Class
}

func (i UInt32) Copy() Value {
	return i
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

func initUInt32() {
	UInt32Class = NewClass()
	StdModule.AddConstantString("UInt32", UInt32Class)
}
