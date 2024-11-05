package value

import (
	"fmt"
	"strconv"

	"github.com/cespare/xxhash/v2"
)

var Int8Class *Class // ::Std::Int8

// Elk's Int8 value
type Int8 int8

func (i Int8) Class() *Class {
	return Int8Class
}

func (Int8) DirectClass() *Class {
	return Int8Class
}

func (Int8) SingletonClass() *Class {
	return nil
}

// Convert to Elk String.
func (i Int8) ToString() String {
	return String(strconv.Itoa(int(i)))
}

// Convert to Elk SmallInt.
func (i Int8) ToSmallInt() SmallInt {
	return SmallInt(i)
}

// Convert to Elk Float.
func (i Int8) ToFloat() Float {
	return Float(i)
}

// Convert to Elk Float64.
func (i Int8) ToFloat64() Float64 {
	return Float64(i)
}

// Convert to Elk Float32.
func (i Int8) ToFloat32() Float32 {
	return Float32(i)
}

// Convert to Elk Int64.
func (i Int8) ToInt64() Int64 {
	return Int64(i)
}

// Convert to Elk Int32.
func (i Int8) ToInt32() Int32 {
	return Int32(i)
}

// Convert to Elk Int16.
func (i Int8) ToInt16() Int16 {
	return Int16(i)
}

// Convert to Elk UInt64.
func (i Int8) ToUInt64() UInt64 {
	return UInt64(i)
}

// Convert to Elk UInt32.
func (i Int8) ToUInt32() UInt32 {
	return UInt32(i)
}

// Convert to Elk UInt16.
func (i Int8) ToUInt16() UInt16 {
	return UInt16(i)
}

// Convert to Elk UInt8.
func (i Int8) ToUInt8() UInt8 {
	return UInt8(i)
}

func (i Int8) Copy() Value {
	return i
}

func (i Int8) Inspect() string {
	return fmt.Sprintf("%di8", i)
}

func (i Int8) InstanceVariables() SymbolMap {
	return nil
}

func (i Int8) Hash() UInt64 {
	d := xxhash.New()
	d.Write([]byte{byte(i)})
	return UInt64(d.Sum64())
}

func initInt8() {
	Int8Class = NewClass()
	StdModule.AddConstantString("Int8", Int8Class)
}
