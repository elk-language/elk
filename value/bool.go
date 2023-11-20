package value

var BoolClass *Class // ::Std::Bool

type Bool interface {
	Value
	bool()
}

func (TrueType) bool()  {}
func (FalseType) bool() {}

func initBool() {
	BoolClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	StdModule.AddConstantString("Bool", BoolClass)
}
