package value

// Elk Method object
type Method interface {
	Reference
	Function
	// Name of the method
	Name() Symbol
}

var MethodClass *Class // ::Std::Method

func initMethod() {
	MethodClass = NewClassWithOptions(
		ClassWithParent(FunctionClass),
	)
	StdModule.AddConstantString("Method", Ref(MethodClass))
}
