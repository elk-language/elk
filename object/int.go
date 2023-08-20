package object

var IntClass *Class // ::Std::Int

func initInt() {
	IntClass = NewClass(
		ClassWithParent(NumericClass),
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	StdModule.AddConstant("Int", IntClass)
}
