package object

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
	NilClass = NewClass(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
		ClassWithImmutable(),
	)
	StdModule.AddConstant("Nil", NilClass)
}
