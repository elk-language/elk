package object

var NumericClass *Class // ::Std::Numeric

func initNumeric() {
	NumericClass = NewClass(ClassWithAbstract())
	StdModule.AddConstant("Numeric", NumericClass)
}
