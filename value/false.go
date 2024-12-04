package value

import (
	"unsafe"

	"github.com/cespare/xxhash/v2"
)

var FalseClass *Class // ::Std::False

type FalseType struct{}

// Elk's false value
var False = FalseType{}.ToValue()

func (FalseType) ToValue() Value {
	return Value{
		data: unsafe.Pointer(uintptr(FALSE_FLAG)),
	}
}

func (FalseType) Class() *Class {
	return FalseClass
}

func (FalseType) DirectClass() *Class {
	return FalseClass
}

func (FalseType) SingletonClass() *Class {
	return nil
}

func (FalseType) Hash() UInt64 {
	d := xxhash.New()
	d.Write([]byte{0})
	return UInt64(d.Sum64())
}

func (FalseType) Inspect() string {
	return "false"
}

func (FalseType) InstanceVariables() SymbolMap {
	return nil
}

func (f FalseType) Error() string {
	return f.Inspect()
}

func initFalse() {
	FalseClass = NewClassWithOptions(
		ClassWithParent(BoolClass),
	)
	StdModule.AddConstantString("False", Ref(FalseClass))
}
