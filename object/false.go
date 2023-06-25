package object

var FalseClass *Class // ::Std::False

// Elk's false value
type False struct{}

func (False) Class() *Class {
	return FalseClass
}

func (False) IsFrozen() bool {
	return true
}

func (False) SetFrozen() {}

func (False) Inspect() string {
	return "false"
}

func (False) InstanceVariables() SimpleSymbolMap {
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
