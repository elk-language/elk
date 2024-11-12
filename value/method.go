package value

// Elk Method object
type Method interface {
	Value
	Function
	// Name of the method
	Name() Symbol
}

var MethodClass *Class // ::Std::Method

func initMethod() {
	MethodClass = NewClassWithOptions(
		ClassWithParent(FunctionClass),
	)
	StdModule.AddConstantString("Method", MethodClass)
}
