package value

// Elk Method object
type Method interface {
	Value
	Function
	IsSealed() bool // Whether the method is non-overridable
	SetSealed()
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
