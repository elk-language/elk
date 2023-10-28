package value

var NilClass *Class // ::Std::Nil

type NilType struct{}

// Elk's Nil value
var Nil = NilType{}

func (NilType) Class() *Class {
	return NilClass
}

func (NilType) IsFrozen() bool {
	return true
}

func (NilType) SetFrozen() {}

func (NilType) Inspect() string {
	return "nil"
}

func (NilType) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initNil() {
	NilClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
		ClassWithImmutable(),
	)
	StdModule.AddConstantString("Nil", NilClass)
}
