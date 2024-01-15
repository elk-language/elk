package value

import (
	"fmt"

	"github.com/cespare/xxhash/v2"
)

var UInt8Class *Class // ::Std::UInt8

// Elk's UInt8 value
type UInt8 uint8

func (UInt8) Class() *Class {
	return UInt8Class
}

func (i UInt8) Copy() Value {
	return i
}

func (UInt8) DirectClass() *Class {
	return UInt64Class
}

func (UInt8) SingletonClass() *Class {
	return nil
}

// Convert to Elk String.
func (i UInt8) ToString() String {
	return String(fmt.Sprintf("%d", i))
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

func (i UInt8) InstanceVariables() SymbolMap {
	return nil
}

func (i UInt8) Hash() UInt64 {
	d := xxhash.New()
	d.Write([]byte{byte(i)})
	return UInt64(d.Sum64())
}

func initUInt8() {
	UInt8Class = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("UInt8", UInt8Class)
}
