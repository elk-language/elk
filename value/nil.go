package value

import "github.com/cespare/xxhash/v2"

var NilClass *Class // ::Std::Nil

type NilType struct{}

// Elk's Nil value
var Nil = NilType{}

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

func (n NilType) Copy() Value {
	return n
}

func (NilType) InstanceVariables() SymbolMap {
	return nil
}

func (NilType) Hash() UInt64 {
	d := xxhash.New()
	d.Write([]byte{2})
	return UInt64(d.Sum64())
}

func initNil() {
	NilClass = NewClass()
	StdModule.AddConstantString("Nil", NilClass)
}
