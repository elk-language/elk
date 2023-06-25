package object

var NilClass *Class // ::Std::Nil

// Elk's Nil value
type Nil struct{}

func (Nil) Class() *Class {
	return NilClass
}

func (Nil) IsFrozen() bool {
	return true
}

func (Nil) SetFrozen() {}

func (Nil) Inspect() string {
	return "nil"
}

func (Nil) InstanceVariables() SimpleSymbolMap {
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
