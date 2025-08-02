package value

import (
	"github.com/cespare/xxhash/v2"
)

var NilClass *Class // ::Std::Nil

type NilType struct{}

// Elk's Nil value
var Nil Value = NilType{}.ToValue()

func (n NilType) ToValue() Value {
	return Value{
		flag: NIL_FLAG,
	}
}

func (NilType) Class() *Class {
	return NilClass
}

func (NilType) DirectClass() *Class {
	return NilClass
}

func (NilType) SingletonClass() *Class {
	return nil
}

func (NilType) Inspect() string {
	return "nil"
}

func (NilType) ToString() String {
	return ""
}

func (NilType) ToChar() Char {
	return 0
}

func (NilType) ToSmallInt() SmallInt {
	return 0
}

func (NilType) ToInt() Value {
	return SmallInt(0).ToValue()
}

func (NilType) ToInt64() Int64 {
	return 0
}

func (NilType) ToInt32() Int32 {
	return 0
}

func (NilType) ToInt16() Int16 {
	return 0
}

func (NilType) ToInt8() Int8 {
	return 0
}

func (NilType) ToUInt64() UInt64 {
	return 0
}

func (NilType) ToUInt32() UInt32 {
	return 0
}

func (NilType) ToUInt16() UInt16 {
	return 0
}

func (NilType) ToUInt8() UInt8 {
	return 0
}

func (NilType) ToFloat() Float {
	return 0
}

func (NilType) ToFloat64() Float64 {
	return 0
}

func (NilType) ToFloat32() Float32 {
	return 0
}

func (NilType) ToBigFloat() *BigFloat {
	return NewBigFloat(0)
}

func (n NilType) Error() string {
	return n.Inspect()
}

func (NilType) InstanceVariables() *InstanceVariables {
	return nil
}

func (NilType) Hash() UInt64 {
	d := xxhash.New()
	d.Write([]byte{2})
	return UInt64(d.Sum64())
}

func initNil() {
	NilClass = NewClassWithOptions(ClassWithParent(ValueClass))
	StdModule.AddConstantString("Nil", Ref(NilClass))
}
