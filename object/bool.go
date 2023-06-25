package object

var BoolClass *Class // ::Std::Bool

type Bool interface {
	Value
	bool()
}

func (TrueType) bool()  {}
func (FalseType) bool() {}

func initBool() {
	BoolClass = NewClass(
		ClassWithNoInstanceVariables(),
		ClassWithImmutable(),
		ClassWithSealed(),
	)
	StdModule.AddConstant("Bool", BoolClass)
}
