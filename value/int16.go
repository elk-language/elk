package value

import (
	"encoding/binary"
	"fmt"

	"github.com/cespare/xxhash/v2"
)

var Int16Class *Class // ::Std::Int16

// Elk's Int16 value
type Int16 int16

func (i Int16) Class() *Class {
	return Int16Class
}

func (Int16) DirectClass() *Class {
	return Int64Class
}

func (Int16) SingletonClass() *Class {
	return nil
}

// Convert to Elk String.
func (i Int16) ToString() String {
	return String(fmt.Sprintf("%d", i))
}

// Convert to Elk SmallInt.
func (i Int16) ToSmallInt() SmallInt {
	return SmallInt(i)
}

// Convert to Elk Float.
func (i Int16) ToFloat() Float {
	return Float(i)
}

// Convert to Elk Float64.
func (i Int16) ToFloat64() Float64 {
	return Float64(i)
}

// Convert to Elk Float32.
func (i Int16) ToFloat32() Float32 {
	return Float32(i)
}

// Convert to Elk Int64.
func (i Int16) ToInt64() Int64 {
	return Int64(i)
}

// Convert to Elk Int32.
func (i Int16) ToInt32() Int32 {
	return Int32(i)
}

// Convert to Elk Int8.
func (i Int8) ToInt8() Int8 {
	return Int8(i)
}

// Convert to Elk UInt64.
func (i Int16) ToUInt64() UInt64 {
	return UInt64(i)
}

// Convert to Elk UInt32.
func (i Int16) ToUInt32() UInt32 {
	return UInt32(i)
}

// Convert to Elk UInt16.
func (i Int16) ToUInt16() UInt16 {
	return UInt16(i)
}

// Convert to Elk UInt8.
func (i Int16) ToUInt8() UInt8 {
	return UInt8(i)
}

func (i Int16) Copy() Value {
	return i
}

func (i Int16) Inspect() string {
	return fmt.Sprintf("%di16", i)
}

func (i Int16) InstanceVariables() SymbolMap {
	return nil
}

func (i Int16) Hash() UInt64 {
	d := xxhash.New()
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(i))
	d.Write(b)
	return UInt64(d.Sum64())
}

func initInt16() {
	Int16Class = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Int16", Int16Class)
}
