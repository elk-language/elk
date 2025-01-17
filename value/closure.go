package value

var ClosureClass *Class // ::Std::Closure

func initClosure() {
	ClosureClass = NewClassWithOptions(
		ClassWithParent(FunctionClass),
	)
	StdModule.AddConstantString("Closure", Ref(ClosureClass))
}
