package value

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
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	StdModule.AddConstantString("True", TrueClass)
}
