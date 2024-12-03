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

func initInt64() {
	Int64Class = NewClass()
	StdModule.AddConstantString("Int64", Int64Class)
}
