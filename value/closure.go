package value

var ClosureClass *Class // ::Std::Closure

func initClosure() {
	ClosureClass = NewClassWithOptions(
		ClassWithSuperclass(FunctionClass),
	)
	StdModule.AddConstantString("Closure", Ref(ClosureClass))
	RegisterNativeClass("Std::Closure", "value.ClosureClass")
}
