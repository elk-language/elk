package value

import "github.com/cespare/xxhash/v2"

var TrueClass *Class // ::Std::True

type TrueType struct{}

// Elk's true value
var True = TrueType{}

func (TrueType) Class() *Class {
	return TrueClass
}

func (TrueType) DirectClass() *Class {
	return TrueClass
}

func (TrueType) Hash() UInt64 {
	d := xxhash.New()
	d.Write([]byte{1})
	return UInt64(d.Sum64())
}

func (TrueType) SingletonClass() *Class {
	return nil
}

func (t TrueType) Copy() Value {
	return t
}

func (TrueType) Inspect() string {
	return "true"
}

func (TrueType) InstanceVariables() SymbolMap {
	return nil
}

func initTrue() {
	TrueClass = NewClassWithOptions(
		ClassWithParent(BoolClass),
	)
	StdModule.AddConstantString("True", TrueClass)
}
