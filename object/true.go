package object

var TrueClass *Class // ::Std::True

// Elk's true value
type True struct{}

func (True) Class() *Class {
	return TrueClass
}

func (True) IsFrozen() bool {
	return true
}

func (True) SetFrozen() {}

func (True) Inspect() string {
	return "true"
}

func (True) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initTrue() {
	TrueClass = NewClass(
		ClassWithParent(BoolClass),
		ClassWithNoInstanceVariables(),
		ClassWithImmutable(),
		ClassWithSealed(),
	)
	StdModule.AddConstant("True", TrueClass)
}
