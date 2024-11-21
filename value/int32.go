package value

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/cespare/xxhash/v2"
)

var Int32Class *Class // ::Std::Int32

// Elk's Int32 value
type Int32 int32

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

// Convert to Elk Int64.
func (i Int32) ToInt64() Int64 {
	return Int64(i)
}

// Convert to Elk Int16.
func (i Int32) ToInt16() Int16 {
	return Int16(i)
}

// Convert to Elk Int8.
func (i Int32) ToInt8() Int8 {
	return Int8(i)
}

// Convert to Elk UInt64.
func (i Int32) ToUInt64() UInt64 {
	return UInt64(i)
}

// Convert to Elk UInt32.
func (i Int32) ToUInt32() UInt32 {
	return UInt32(i)
}

// Convert to Elk UInt16.
func (i Int32) ToUInt16() UInt16 {
	return UInt16(i)
}

// Convert to Elk UInt8.
func (i Int32) ToUInt8() UInt8 {
	return UInt8(i)
}

func (i Int32) Copy() Value {
	return i
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

func initInt32() {
	Int32Class = NewClass()
	StdModule.AddConstantString("Int32", Int32Class)
}
