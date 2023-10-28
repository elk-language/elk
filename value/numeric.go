package value

var NumericClass *Class // ::Std::Numeric

func initNumeric() {
	NumericClass = NewClassWithOptions(ClassWithAbstract())
	StdModule.AddConstantString("Numeric", NumericClass)
}
