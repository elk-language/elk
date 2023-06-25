package object

var BoolClass *Class // ::Std::Bool

func initBool() {
	BoolClass = NewClass(
		ClassWithNoInstanceVariables(),
		ClassWithImmutable(),
		ClassWithSealed(),
	)
	StdModule.AddConstant("Bool", BoolClass)
}
