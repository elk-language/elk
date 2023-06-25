package object

var FalseClass *Class // ::Std::False

type FalseType struct{}

// Elk's false value
var False = FalseType{}

func (FalseType) Class() *Class {
	return FalseClass
}

func (FalseType) IsFrozen() bool {
	return true
}

func (FalseType) SetFrozen() {}

func (FalseType) Inspect() string {
	return "false"
}

func (FalseType) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initFalse() {
	FalseClass = NewClass(
		ClassWithParent(BoolClass),
		ClassWithNoInstanceVariables(),
		ClassWithImmutable(),
		ClassWithSealed(),
	)
	StdModule.AddConstant("False", FalseClass)
}
