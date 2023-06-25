package object

var TrueClass *Class // ::Std::True

type TrueType struct{}

// Elk's true value
var True = TrueType{}

func (TrueType) Class() *Class {
	return TrueClass
}

func (TrueType) IsFrozen() bool {
	return true
}

func (TrueType) SetFrozen() {}

func (TrueType) Inspect() string {
	return "true"
}

func (TrueType) InstanceVariables() SimpleSymbolMap {
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
