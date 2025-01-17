package value

// Elk function object
type Function interface {
	Reference
	// The number of parameters in total
	ParameterCount() int
	// The number of optional parameters with default values
	OptionalParameterCount() int
}

var FunctionClass *Class // ::Std::Function

func initFunction() {
	StdModule.AddConstantString("Function", Ref(FunctionClass))
}
