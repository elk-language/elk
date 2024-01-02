package value

var FalseClass *Class // ::Std::False

type FalseType struct{}

// Elk's false value
var False = FalseType{}

func (FalseType) Class() *Class {
	return FalseClass
}

func (FalseType) DirectClass() *Class {
	return FalseClass
}

func (FalseType) SingletonClass() *Class {
	return nil
}

func (f FalseType) Copy() Value {
	return f
}

func (FalseType) Inspect() string {
	return "false"
}

func (FalseType) InstanceVariables() SymbolMap {
	return nil
}

func initFalse() {
	FalseClass = NewClassWithOptions(
		ClassWithParent(BoolClass),
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	StdModule.AddConstantString("False", FalseClass)
}
