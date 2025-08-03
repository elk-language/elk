package value

var ClosureClass *Class // ::Std::Closure

func initClosure() {
	ClosureClass = NewClassWithOptions(
		ClassWithSuperclass(FunctionClass),
	)
	StdModule.AddConstantString("Closure", Ref(ClosureClass))
}
